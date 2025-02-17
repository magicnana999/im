package client

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/magicnana999/im/enum"
	"github.com/magicnana999/im/logger"
	"github.com/magicnana999/im/pb"
	"github.com/magicnana999/im/util/id"
	"github.com/panjf2000/ants/v2"
	bb "github.com/panjf2000/gnet/v2/pkg/pool/bytebuffer"
	goPool "github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type Client struct {
	UserSig              string
	UserId               int64
	To                   int64
	ServerAddress        string
	HeartbeatInterval    int // 心跳间隔
	CurrentUserId        int64
	CurrentAppId         string
	LastHeartbeatReceive int64
	C                    net.Conn
	Wg                   sync.WaitGroup
	Sender               *sender
}

func (c *Client) Start() {

	logger.InitLogger("debug")

	conn, err := net.Dial("tcp", c.ServerAddress)
	if err != nil {
		logger.FatalF("无法连接到服务器: %v", err)
	}
	defer conn.Close()

	c.C = conn
	c.Wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())

	c.Sender = initSender(conn, ctx)
	go c.Sender.start()

	go c.startHeartbeat(ctx, c.Sender, cancel)

	go c.startRead(ctx, conn, c.Sender)

	go c.getMessageFromScan(ctx, cancel, conn, &c.Wg)

	login(c.Sender, c.UserSig)

	c.Wg.Wait()
	fmt.Println("exit")
}

func (c *Client) getMessageFromScan(ctx context.Context, cancel context.CancelFunc, conn net.Conn, wg *sync.WaitGroup) {

	defer wg.Done()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if !scanner.Scan() {
				continue
			}
			message := scanner.Text()
			if message == "exit" {
				cancel()
			} else {
				c.saySomething(message)
			}
		}
	}
}

func (c *Client) saySomething(t string) {

	text := &pb.TextContent{Text: t}

	packet := pb.NewMessage(c.UserId, c.To, 0, 100, "19860220", "1-1", text).Wrap()
	c.Sender.send(packet)
}

func (c *Client) startHeartbeat(ctx context.Context, sender *sender, stop context.CancelFunc) {

	ticker := time.NewTicker(time.Duration(c.HeartbeatInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:

			if c.LastHeartbeatReceive > 0 && time.Now().UnixMilli()-c.LastHeartbeatReceive >= time.Minute.Milliseconds() {
				stop()
				return
			}

			packet := pb.NewHeartbeat(12)
			sender.send(packet)

		case <-ctx.Done():
			return
		}
	}
}

func (c *Client) startRead(ctx context.Context, conn net.Conn, sender *sender) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			c.decode(ctx, conn, sender)
		}
	}
}

func (c *Client) decode(ctx context.Context, conn net.Conn, sender *sender) error {

	var length int32
	binary.Read(conn, binary.BigEndian, &length)

	if length == 4 {

		var heartbeat int32
		binary.Read(conn, binary.BigEndian, &heartbeat)

		p := pb.NewHeartbeat(heartbeat)
		c.handlePacket(ctx, p, sender)
		return nil
	}

	if length > 4 {

		bs := make([]byte, length)
		l, e := io.ReadFull(conn, bs)
		if e != nil || l != int(length) {
			panic(e)
		}

		var p pb.Packet
		if e4 := proto.Unmarshal(bs, &p); e4 != nil {
			panic(e4)
		}

		bss, _ := protojson.Marshal(&p)
		logger.DebugF("收到：%s\n\n", string(bss))

		c.handlePacket(ctx, &p, sender)
		return nil
	}

	return nil
}

func (c *Client) handlePacket(ctx context.Context, packet *pb.Packet, sender *sender) {
	switch packet.Type {
	case pb.TypeHeartbeat:
		c.LastHeartbeatReceive = time.Now().UnixMilli()
		return
	case pb.TypeMessage:
		mb := packet.GetMessageBody()
		if mb.IsResponse() {
			sender.close(mb.Id)
		} else {
			c.receiveMessage(ctx, mb, sender)
		}
	case pb.TypeCommand:
		c.receiveCommand(ctx, packet, sender)
	default:
		fmt.Printf("不知道啥Type: %d\n", packet.Type)
	}
}

func (c *Client) receiveCommand(ctx context.Context, packet *pb.Packet, s *sender) {

	switch packet.GetCommandBody().CType {
	case pb.CTypeUserLogin:
		reply := packet.GetCommandBody().GetLoginReply()
		if packet.GetCommandBody().Code == 0 {
			c.CurrentAppId = reply.AppId
			c.CurrentUserId = reply.UserId
			fmt.Println("登录成功")
		} else {
			fmt.Println("登录失败,", packet.GetCommandBody().Message)

		}
	default:
		return
	}
}

func (c *Client) receiveMessage(ctx context.Context, mb *pb.MessageBody, s *sender) {
	write(c.C, mb.Success(nil).Wrap())
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

type resendTask struct {
	ctx      context.Context
	cancel   context.CancelFunc
	id       string
	interval int
	packet   *pb.Packet
	ticker   *time.Ticker
	conn     net.Conn
}

func newResendTask(ctx context.Context, cancel context.CancelFunc, interval int, cmd *pb.Packet, conn net.Conn) *resendTask {
	return &resendTask{
		ctx:      ctx,
		cancel:   cancel,
		id:       cmd.GetMessageBody().Id,
		interval: interval,
		packet:   cmd,
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

			switch packet.Type {
			case pb.TypeHeartbeat:
				s.sendHeartbeat(packet)
			case pb.TypeCommand:
				s.sendCommand(packet)
			case pb.TypeMessage:
				s.sendMessage(packet)
			}
		}
	}
}

func (s *sender) sendHeartbeat(packet *pb.Packet) {
	write(s.conn, packet)
}

func (s *sender) sendCommand(packet *pb.Packet) {
	write(s.conn, packet)
}

func (s *sender) sendMessage(packet *pb.Packet) {
	s.lock.RLock()
	mb := packet.GetMessageBody()
	_, exist := s.m[mb.Id]
	s.lock.RUnlock()
	if !exist {
		subCtx, cancel := context.WithCancel(s.ctx)
		task := newResendTask(subCtx, cancel, 1, packet, s.conn)

		s.lock.Lock()
		if _, doubleCheck := s.m[packet.GetMessageBody().GetId()]; doubleCheck { // Double-check 防止并发问题
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

						fmt.Println("重试超过限制，关闭连接:", packet.GetMessageBody().GetId())
						s.close(packet.GetMessageBody().GetId())
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

func write(conn net.Conn, packet *pb.Packet) (int, error) {

	buffer, err := encode(packet)
	defer bb.Put(buffer)

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

func encode(p *pb.Packet) (*bb.ByteBuffer, error) {

	if p.Type != pb.TypeHeartbeat {
		bbs, _ := protojson.Marshal(p)
		logger.DebugF("发送：%s\n\n", string(bbs))
	}

	buffer := bb.Get()

	if p.IsHeartbeat() {
		binary.Write(buffer, binary.BigEndian, uint32(4))
		binary.Write(buffer, binary.BigEndian, p.GetHeartbeatBody().Value)
	} else {

		bs, e := proto.Marshal(p)
		if e != nil {
			panic(e)
		}
		binary.Write(buffer, binary.BigEndian, uint32(len(bs)))
		binary.Write(buffer, binary.BigEndian, bs)
	}
	return buffer, nil
}

func login(sender *sender, userSig string) {
	loginRequest := pb.LoginRequest{
		AppId:        "19860220",
		UserSig:      userSig,
		Version:      "1.0.0",
		Os:           int32(enum.Ios),
		PushDeviceId: strings.ToLower(id.GenerateXId()),
	}

	request := pb.NewCommand(&loginRequest)

	sender.send(request)
}
