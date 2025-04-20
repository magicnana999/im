package main

import (
	"context"
	"github.com/magicnana999/im/api/kitex_gen/api"
	"github.com/magicnana999/im/broker"
	"github.com/magicnana999/im/pkg/timewheel"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"time"
)

var key string = "USER"

type TcpClientListener interface {
	Connect(*User)
	DisConnect(*User)
}

type TcpServer struct {
	agent      *gnet.Client
	codec      *broker.Codec
	handler    *PacketHandler
	htsHandler *HeartbeatServer
	listener   TcpClientListener
}

func (h *TcpServer) OnShutdown(eng gnet.Engine) {}
func (h *TcpServer) OnTick() (t time.Duration, action gnet.Action) {
	return time.Second * 20, gnet.None
}

func (h *TcpServer) OnBoot(eng gnet.Engine) gnet.Action {
	logging.Infof("Client booted")
	return gnet.None
}

func GetUser(c gnet.Conn) *User {
	ctx, ok := c.Context().(context.Context)
	if ok && ctx != nil {
		u, ok := ctx.Value(key).(*User)
		if ok && u != nil {
			return u
		}
	}

	return nil
}

func (h *TcpServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	logging.Infof("Connected to server: %s", c.RemoteAddr().String())

	user := &User{
		FD:     c.Fd(),
		Writer: c,
		Reader: c,
	}

	ctx := context.WithValue(context.Background(), key, user)
	c.SetContext(ctx)

	h.htsHandler.StartTicking(func(now time.Time) timewheel.TaskResult {
		if user.IsClosed.Load() {
			return timewheel.Break
		}
		ht := api.NewHeartbeatPacket(100)
		h.handler.Write(ht, c)
		return timewheel.Retry
	}, time.Second*30)

	h.listener.Connect(user)
	return
}

func (h *TcpServer) OnTraffic(c gnet.Conn) gnet.Action {

	user := GetUser(c)
	if user == nil {
		c.Close()
		return gnet.None
	}

	ps, err := h.codec.Decode(c)
	if err != nil {
		c.Close()
		return gnet.None
	}

	for _, packet := range ps {
		if packet.IsHeartbeat() {
			user.LastHTS.Store(time.Now().Unix())
		}

		ret := h.handler.Handle(packet, user)
		if ret != nil {
			h.handler.Write(ret, c)
		}
	}

	return gnet.None
}

func (h *TcpServer) OnClose(c gnet.Conn, err error) gnet.Action {
	logging.Infof("Connection closed: %v", err)
	user := GetUser(c)
	h.listener.DisConnect(user)
	return gnet.None
}

func (h *TcpServer) Start() {
	h.agent.Start()
}

func (h *TcpServer) Stop() {
	h.agent.Stop()
}

func (h *TcpServer) Connect() (gnet.Conn, error) {
	return h.agent.Dial("tcp", "127.0.0.1:5075")
}

func NewTcpServer(handler *PacketHandler, hts *HeartbeatServer, l TcpClientListener) *TcpServer {

	codec := broker.NewCodec()
	eh := &TcpServer{codec: codec, handler: handler, htsHandler: hts, listener: l}
	client, err := gnet.NewClient(eh,
		gnet.WithLogger(logging.GetDefaultLogger()),
		gnet.WithReadBufferCap(1024),
		gnet.WithWriteBufferCap(1024))
	if err != nil {
		logging.Fatalf("Failed to create client: %v", err)
	}

	eh.agent = client
	return eh
}
