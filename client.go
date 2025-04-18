package main

//
//import (
//	"bufio"
//	"context"
//	"encoding/binary"
//	"fmt"
//	"github.com/magicnana999/im/api/kitex_gen/api"
//	"github.com/magicnana999/im/pkg/id"
//	"github.com/panjf2000/ants/v2"
//	bb "github.com/panjf2000/gnet/v2/pkg/pool/bytebuffer"
//	goPool "github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
//	"google.golang.org/protobuf/encoding/protojson"
//	"google.golang.org/protobuf/proto"
//	"io"
//	"net"
//	"os"
//	"strings"
//	"sync"
//	"time"
//)
//
//type User struct {
//	c      net.Conn
//	appId  string
//	userId int64
//}
//
//type Client struct {
//	BrokerAddr           string
//	HeartbeatInterval    int // 心跳间隔
//	userId               int64
//	appId                string
//	LastHeartbeatReceive int64
//	Sender               *sender
//}
//
//func (c *Client) Start() {
//
//	conn, err := net.Dial("tcp", c.ServerAddress)
//	if err != nil {
//		fmt.Errorf("无法连接到服务器: %v", err)
//	}
//	defer conn.Close()
//
//	c.C = conn
//	c.Wg.Add(1)
//
//	ctx, cancel := context.WithCancel(context.Background())
//
//	c.Sender = initSender(conn, ctx)
//	go c.Sender.start()
//
//	go c.startHeartbeat(ctx, c.Sender, cancel)
//
//	go c.startRead(ctx, conn, c.Sender)
//
//	go c.getMessageFromScan(ctx, cancel, conn, &c.Wg)
//
//	login(c.Sender, c.UserSig)
//
//	c.Wg.Wait()
//	fmt.Println("exit")
//}
//
//func (c *Client) getMessageFromScan(ctx context.Context, cancel context.CancelFunc, conn net.Conn, wg *sync.WaitGroup) {
//
//	defer wg.Done()
//
//	scanner := bufio.NewScanner(os.Stdin)
//
//	for {
//		select {
//		case <-ctx.Done():
//			return
//		default:
//			if !scanner.Scan() {
//				continue
//			}
//			message := scanner.Text()
//			if message == "exit" {
//				cancel()
//			} else {
//				c.saySomething(message)
//			}
//		}
//	}
//}
//
//func (c *Client) saySomething(t string) {
//
//	text := &api.Text{Text: t}
//
//	packet := api.NewMessage(c.UserId, c.To, 0, 100, "19860220", "1-1", text).Wrap()
//	c.Sender.send(packet)
//}
//
//func (c *Client) startHeartbeat(ctx context.Context, sender *sender, stop context.CancelFunc) {
//
//	ticker := time.NewTicker(time.Duration(c.HeartbeatInterval) * time.Second)
//	defer ticker.Stop()
//
//	for {
//		select {
//		case <-ticker.C:
//
//			if c.LastHeartbeatReceive > 0 && time.Now().UnixMilli()-c.LastHeartbeatReceive >= time.Minute.Milliseconds() {
//				stop()
//				return
//			}
//
//			packet := api.NewHeartbeat(12).Wrap()
//			sender.send(packet)
//
//		case <-ctx.Done():
//			return
//		}
//	}
//}
//
//func (c *Client) startRead(ctx context.Context, conn net.Conn, sender *sender) {
//	for {
//		select {
//		case <-ctx.Done():
//			return
//		default:
//			c.decode(ctx, conn, sender)
//		}
//	}
//}
//
//func (c *Client) decode(ctx context.Context, conn net.Conn, sender *sender) error {
//
//	var length int32
//	binary.Read(conn, binary.BigEndian, &length)
//
//	if length == 4 {
//
//		var heartbeat int32
//		binary.Read(conn, binary.BigEndian, &heartbeat)
//
//		p := api.NewHeartbeat(heartbeat).Wrap()
//		c.handlePacket(ctx, p, sender)
//		return nil
//	}
//
//	if length > 4 {
//
//		bs := make([]byte, length)
//		l, e := io.ReadFull(conn, bs)
//		if e != nil || l != int(length) {
//			panic(e)
//		}
//
//		var p api.Packet
//		if e4 := proto.Unmarshal(bs, &p); e4 != nil {
//			panic(e4)
//		}
//
//		bss, _ := protojson.Marshal(&p)
//		fmt.Printf("收到：%s\n\n", string(bss))
//
//		c.handlePacket(ctx, &p, sender)
//		return nil
//	}
//
//	return nil
//}
//
//func (c *Client) handlePacket(ctx context.Context, packet *api.Packet, sender *sender) {
//	switch packet.Type {
//	case api.TypeHeartbeat:
//		c.LastHeartbeatReceive = time.Now().UnixMilli()
//		return
//	case api.TypeMessage:
//		mb := packet.GetMessage()
//		if mb.IsResponse() {
//			sender.close(mb.MessageId)
//		} else {
//			c.receiveMessage(ctx, mb, sender)
//		}
//	case api.TypeCommand:
//		c.receiveCommand(ctx, packet, sender)
//	default:
//		fmt.Printf("不知道啥Type: %d\n", packet.Type)
//	}
//}
//
//func (c *Client) receiveCommand(ctx context.Context, packet *api.Packet, s *sender) {
//
//	switch packet.GetCommand().CommandType {
//	case api.CommandTypeUserLogin:
//		reply := packet.GetCommand().GetLoginReply()
//		if packet.GetCommand().Code == 0 {
//			c.CurrentAppId = reply.AppId
//			c.CurrentUserId = reply.UserId
//			fmt.Println("登录成功")
//		} else {
//			fmt.Println("登录失败,", packet.GetCommand().Message)
//
//		}
//	default:
//		return
//	}
//}
//
//func (c *Client) receiveMessage(ctx context.Context, mb *api.Message, s *sender) {
//	write(c.C, mb.Success(nil).Wrap())
//}
//
//func fibonacci(n int) int {
//	if n <= 1 {
//		return n
//	}
//	a, b := 0, 1
//	for i := 2; i <= n; i++ {
//		a, b = b, a+b
//	}
//	return b
//}
//
//type resendTask struct {
//	ctx      context.Context
//	cancel   context.CancelFunc
//	id       string
//	interval int
//	packet   *api.Packet
//	ticker   *time.Ticker
//	conn     net.Conn
//}
//
//func newResendTask(ctx context.Context, cancel context.CancelFunc, interval int, cmd *api.Packet, conn net.Conn) *resendTask {
//	return &resendTask{
//		ctx:      ctx,
//		cancel:   cancel,
//		id:       cmd.GetMessage().MessageId,
//		interval: interval,
//		packet:   cmd,
//		ticker:   time.NewTicker(time.Duration(interval) * time.Second),
//		conn:     conn,
//	}
//}
//
//type sender struct {
//	ctx      context.Context
//	conn     net.Conn
//	packets  chan *api.Packet
//	executor *goPool.Pool
//	m        map[string]*resendTask
//	lock     sync.RWMutex
//}
//
//func initSender(conn net.Conn, ctx context.Context) *sender {
//	pool, err := ants.NewPool(100)
//	if err != nil {
//		panic(err)
//	}
//
//	return &sender{
//		ctx:      ctx,
//		conn:     conn,
//		packets:  make(chan *api.Packet),
//		executor: pool,
//		m:        make(map[string]*resendTask),
//	}
//}
//
//func (s *sender) closeAll() {
//	s.lock.Lock()
//	defer s.lock.Unlock()
//
//	for id, _ := range s.m {
//		s.close(id)
//	}
//}
//
//func (s *sender) send(packet *api.Packet) {
//
//	s.packets <- packet
//}
//
//func (s *sender) close(id string) {
//	s.lock.Lock()
//	defer s.lock.Unlock()
//	task := s.m[id]
//	if task != nil {
//		task.cancel()
//		delete(s.m, id)
//	}
//}
//
//func (s *sender) start() {
//	for {
//		select {
//		case <-s.ctx.Done():
//			s.closeAll()
//			return
//		case packet, ok := <-s.packets:
//			if !ok {
//				return // 退出 Goroutine
//			}
//
//			switch packet.Type {
//			case api.TypeHeartbeat:
//				s.sendHeartbeat(packet)
//			case api.TypeCommand:
//				s.sendCommand(packet)
//			case api.TypeMessage:
//				s.sendMessage(packet)
//			}
//		}
//	}
//}
//
//func (s *sender) sendHeartbeat(packet *api.Packet) {
//	write(s.conn, packet)
//}
//
//func (s *sender) sendCommand(packet *api.Packet) {
//	write(s.conn, packet)
//}
//
//func (s *sender) sendMessage(packet *api.Packet) {
//	s.lock.RLock()
//	mb := packet.GetMessage()
//	_, exist := s.m[mb.MessageId]
//	s.lock.RUnlock()
//	if !exist {
//		subCtx, cancel := context.WithCancel(s.ctx)
//		task := newResendTask(subCtx, cancel, 1, packet, s.conn)
//
//		s.lock.Lock()
//		if _, doubleCheck := s.m[packet.GetMessage().GetMessageId()]; doubleCheck { // Double-check 防止并发问题
//			s.lock.Unlock()
//			return
//		}
//		s.m[task.id] = task
//		s.lock.Unlock()
//
//		s.executor.Submit(func() {
//			for {
//				select {
//				case <-task.ctx.Done():
//					return
//
//				case <-task.ticker.C:
//					write(task.conn, packet)
//					next := fibonacci(task.interval)
//					if next >= 8 {
//
//						fmt.Println("重试超过限制，关闭连接:", packet.GetMessage().MessageId)
//						s.close(packet.GetMessage().MessageId)
//						s.conn.Close()
//						return
//
//					}
//					task.interval = next
//					task.ticker.Reset(time.Duration(next) * time.Second)
//				}
//			}
//		})
//	}
//}
//
//func write(conn net.Conn, packet *api.Packet) (int, error) {
//
//	buffer, err := encode(packet)
//	defer bb.Put(buffer)
//
//	if err != nil {
//		panic(err)
//	}
//
//	total := buffer.Len()
//	sent := 0
//	for sent < total {
//		n, err := conn.Write(buffer.Bytes()[sent:])
//		if err != nil {
//			return 0, err
//		}
//		sent += n
//	}
//
//	return total, nil
//}
//
//func encode(p *api.Packet) (*bb.ByteBuffer, error) {
//
//	if p.Type != api.TypeHeartbeat {
//		bbs, _ := protojson.Marshal(p)
//		fmt.Printf("发送：%s\n\n", string(bbs))
//	}
//
//	buffer := bb.Get()
//
//	if p.IsHeartbeat() {
//		binary.Write(buffer, binary.BigEndian, uint32(4))
//		binary.Write(buffer, binary.BigEndian, p.GetHeartbeat().Value)
//	} else {
//
//		bs, e := proto.Marshal(p)
//		if e != nil {
//			panic(e)
//		}
//		binary.Write(buffer, binary.BigEndian, uint32(len(bs)))
//		binary.Write(buffer, binary.BigEndian, bs)
//	}
//	return buffer, nil
//}
//
//func login(sender *sender, userSig string) {
//	loginRequest := api.LoginRequest{
//		AppId:    "",
//		UserSig:  userSig,
//		Version:  "1.0.0",
//		Os:       "iOS",
//		DeviceId: strings.ToLower(id.GenerateXId()),
//	}
//
//	request := api.NewCommand(&loginRequest)
//
//	sender.send(request)
//}
