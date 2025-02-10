package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/magicnana999/im/common/pb"
	"github.com/magicnana999/im/common/protocol"
	"github.com/magicnana999/im/errors"
	"github.com/magicnana999/im/logger"
	"github.com/magicnana999/im/util/id"
	"github.com/panjf2000/ants/v2"
	goPool "github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"io"
	"math"
	"net"
	"os"
	"sync"
	"time"
)

const (
	userSig = "cukpovu1a37hpofg6sj0"
	//userSig           = "cuf5ofe1a37nfi3p4b6g"

	serverAddress     = "127.0.0.1:7539" // IM 服务器地址和端口
	heartbeatInterval = 5 * time.Second  // 心跳间隔
)

var (
	CurrentUserId        int64
	CurrentAppId         string
	lastHeartbeatReceive int64 = 0
)

type resendTask struct {
	ctx      context.Context
	cancel   context.CancelFunc
	id       string
	interval int
	packet   *pb.Packet
	ticker   *time.Ticker
	conn     net.Conn
}

func newResendTask(ctx context.Context, cancel context.CancelFunc, interval int, packet *pb.Packet, conn net.Conn) *resendTask {
	return &resendTask{
		ctx:      ctx,
		cancel:   cancel,
		id:       packet.Id,
		interval: interval,
		packet:   packet,
		ticker:   time.NewTicker(time.Duration(interval) * time.Second),
		conn:     conn,
	}
}

type sender struct {
	ctx      context.Context
	conn     net.Conn
	packets  chan *pb.Packet
	executor *goPool.Pool
	m        map[string]*resendTask
	lock     sync.RWMutex
}

func initSender(conn net.Conn, ctx context.Context) *sender {
	pool, err := ants.NewPool(100)
	if err != nil {
		panic(err)
	}

	return &sender{
		ctx:      ctx,
		conn:     conn,
		packets:  make(chan *pb.Packet),
		executor: pool,
		m:        make(map[string]*resendTask),
	}
}

func (s *sender) closeAll() {
	s.lock.Lock()
	defer s.lock.Unlock()

	for id, _ := range s.m {
		s.close(id)
	}
}

func (s *sender) send(packet *pb.Packet) {

	if packet.BType != pb.BTypeHeartbeat {
		fmt.Println("发送 ", packet.Id, packet.BType)
	}

	s.packets <- packet
}

func (s *sender) close(id string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	task := s.m[id]
	if task != nil {
		task.cancel()
		delete(s.m, id)
	}
}

func (s *sender) sendPacket(packet *pb.Packet) {
	write(s.conn, packet)
}

func (s *sender) sendMessage(packet *pb.Packet) {
	s.lock.RLock()
	_, exist := s.m[packet.Id]
	s.lock.RUnlock()
	if !exist {
		subCtx, cancel := context.WithCancel(s.ctx)
		task := newResendTask(subCtx, cancel, 1, packet, s.conn)

		s.lock.Lock()
		if _, doubleCheck := s.m[packet.Id]; doubleCheck { // Double-check 防止并发问题
			s.lock.Unlock()
			return
		}
		s.m[task.id] = task
		s.lock.Unlock()

		s.executor.Submit(func() {
			for {
				select {
				case <-task.ctx.Done():
					return

				case <-task.ticker.C:
					write(task.conn, packet)
					next := fibonacci(task.interval)
					if next >= 8 {

						fmt.Println("重试超过限制，关闭连接:", packet.Id)
						s.close(packet.Id)
						s.conn.Close()
						return

					}
					task.interval = next
					task.ticker.Reset(time.Duration(next) * time.Second)
				}
			}
		})
	}
}

func (s *sender) start() {
	for {
		select {
		case <-s.ctx.Done():
			s.closeAll()
			return
		case packet, ok := <-s.packets:
			if !ok {
				return // 退出 Goroutine
			}

			if pb.IsResponse(packet) {
				s.close(packet.Id)
			}

			if pb.IsMessage(packet) {
				s.sendMessage(packet)
			} else {
				s.sendPacket(packet)
			}
		}
	}
}

func main() {
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		logger.FatalF("无法连接到服务器: %v", err)
	}
	defer conn.Close()

	var wg sync.WaitGroup

	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())

	sender := initSender(conn, ctx)
	go sender.start()

	go startHeartbeat(ctx, sender, cancel)

	go startRead(ctx, conn, sender)

	go sendMessage(ctx, cancel, conn, &wg)

	login(sender, userSig)

	wg.Wait()
	fmt.Println("exit")
}

func sendMessage(ctx context.Context, cancel context.CancelFunc, conn net.Conn, wg *sync.WaitGroup) {

	defer wg.Done()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Println("请输入消息: ")
			if !scanner.Scan() {
				continue
			}
			message := scanner.Text()
			if message == "exit" {
				cancel()
			} else if message == "logout" {

			}
		}
	}
}

func login(sender *sender, userSig string) {
	loginRequest := pb.LoginRequest{
		AppId:        "19860220",
		UserSig:      userSig,
		Version:      "1.0.0",
		Os:           pb.OSType_OSIos,
		PushDeviceId: id.GenerateXId(),
	}

	request, err := pb.NewCommandRequest(0, pb.CTypeUserLogin, &loginRequest)
	if err != nil {
		panic(err)
	}

	sender.send(request)
}

func encode(p *pb.Packet) (buffer *bytes.Buffer, err error) {

	if pb.IsHeartbeat(p) {

		var hb wrapperspb.UInt32Value
		if err := p.Body.UnmarshalTo(&hb); err != nil {
			return nil, errors.ConnectionEncodeError.Fill("failed to unmarshal length field," + err.Error())
		}

		if hb.Value < 0 || hb.Value >= math.MaxInt32 {
			return nil, errors.ConnectionEncodeError.Fill("invalid heartbeat value")
		}

		buffer := new(bytes.Buffer)
		binary.Write(buffer, binary.BigEndian, uint32(4))
		binary.Write(buffer, binary.BigEndian, hb.Value)
		return buffer, nil
	} else {

		var err error

		bs := make([][]byte, 2)
		bs[0] = make([]byte, 4)
		bs[1], err = proto.Marshal(p)

		if err != nil {
			return nil, errors.ConnectionDecodeError.Fill("failed to marshal packet," + err.Error())
		}

		if len(bs[1]) <= 0 || len(bs[1]) >= math.MaxInt32 {
			return nil, errors.ConnectionEncodeError.Fill("invalid packet length")
		}

		buffer := new(bytes.Buffer)
		binary.Write(buffer, binary.BigEndian, uint32(len(bs[1])))
		binary.Write(buffer, binary.BigEndian, bs[1])

		return buffer, nil
	}
}

func startHeartbeat(ctx context.Context, sender *sender, stop context.CancelFunc) {

	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:

			if lastHeartbeatReceive > 0 && time.Now().UnixMilli()-lastHeartbeatReceive >= time.Minute.Milliseconds() {
				stop()
				return
			}

			packet, _ := pb.NewHeartbeatRequest(12)
			sender.send(packet)

		case <-ctx.Done():
			return
		}
	}
}

func startRead(ctx context.Context, conn net.Conn, sender *sender) {

	reader := bufio.NewReader(conn)
	for {

		select {
		case <-ctx.Done():
			return
		default:

			bs := make([]byte, 4)
			l, err := io.ReadFull(reader, bs)
			if err != nil || l != 4 {
				panic(err)
			}

			length := int(binary.BigEndian.Uint32(bs))

			if length == 4 && reader.Buffered() >= length {

				hb := make([]byte, length)
				hbl, hberr := io.ReadFull(reader, hb)
				if hberr != nil || hbl != 4 {
					panic(err)
				}

				heartbeat := binary.BigEndian.Uint32(hb)
				p, _ := pb.NewHeartbeatRequest(int32(heartbeat))
				handleMessage(ctx, p, sender)

			}

			if length > 4 && reader.Buffered() >= length {

				bb := make([]byte, length)
				bbl, bberr := io.ReadFull(reader, bb)
				if bberr != nil || bbl != length {
					panic(bberr)
				}

				var p pb.Packet
				if e4 := proto.Unmarshal(bb, &p); err != nil {
					panic(e4)
				}

				handleMessage(ctx, &p, sender)
			}
		}
	}
}

func handleMessage(ctx context.Context, packet *pb.Packet, sender *sender) {
	switch packet.BType {
	case pb.BTypeHeartbeat:

		var hb wrapperspb.UInt32Value
		if err := packet.Body.UnmarshalTo(&hb); err != nil {
			panic(err)
		}

		lastHeartbeatReceive = time.Now().UnixMilli()
		return
	case pb.BTypeMessage:
		if pb.IsResponse(packet) {
			sender.close(packet.Id)
		} else {
			receiveMessage(ctx, packet, sender)
		}
	case pb.BTypeCommand:

		if pb.IsResponse(packet) {
			receiveCommand(ctx, packet, sender)
		} else {
			panic(errors.New(0, "为什么收到一个request command"))
		}
	default:
		fmt.Printf("不知道啥Type: %d\n", packet.BType)
	}
}

func receiveCommand(ctx context.Context, packet *pb.Packet, s *sender) {
	p, _ := pb.RevertPacket(packet)

	if p.BType == pb.BTypeCommand {

		if body, ok := p.Body.(*protocol.CommandBody); ok {
			handleCommandResponse(packet, body)
		} else {
			panic(errors.New(0, "怎么不是一个commandBody呢"))
		}

	} else {
		panic(errors.New(0, "收到的BType不是command"))
	}
}

func handleCommandResponse(packet *pb.Packet, body *protocol.CommandBody) {
	switch body.CType {
	case pb.CTypeUserLogin:
		if reply, ok := body.Reply.(*protocol.LoginReply); ok && packet.Status.Code == 0 {
			CurrentAppId = reply.AppId
			CurrentUserId = reply.UserId
			fmt.Println("登录成功")
		} else {
			js, _ := json.Marshal(body)
			fmt.Println("登录失败,", string(js))
			panic(errors.New(0, "登录失败"))
		}
	default:
		return
	}
}

func receiveMessage(ctx context.Context, packet *pb.Packet, s *sender) {
	p, _ := pb.RevertPacket(packet)
	js, _ := json.Marshal(p)
	fmt.Println(string(js))
}

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

func write(conn net.Conn, packet *pb.Packet) (int, error) {

	buffer, err := encode(packet)
	if err != nil {
		panic(err)
	}

	total := buffer.Len()
	sent := 0
	for sent < total {
		n, err := conn.Write(buffer.Bytes()[sent:])
		if err != nil {
			return 0, err
		}
		sent += n
	}

	return total, nil
}
