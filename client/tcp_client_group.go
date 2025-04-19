package main

import (
	"github.com/magicnana999/im/broker"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
)

// 实现 gnet 的事件处理器
type clientHandler struct {
	gnet.EventHandler
	codec *broker.Codec
}

func (h *clientHandler) OnBoot(eng gnet.Engine) gnet.Action {
	logging.Infof("Client booted")
	return gnet.None
}

func (h *clientHandler) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	logging.Infof("Connected to server: %s", c.RemoteAddr().String())
	// 发送初始消息
	out = []byte("Hello, gnet!\n")
	return
}

func (h *clientHandler) OnTraffic(c gnet.Conn) gnet.Action {
	return gnet.None
}

func (h *clientHandler) OnClose(c gnet.Conn, err error) gnet.Action {
	logging.Infof("Connection closed: %v", err)
	return gnet.None
}

type TcpClientGroup struct {
	clients *gnet.Client
	codec   *broker.Codec
}

func NewTcpClientGroup() *TcpClientGroup {

	codec := broker.NewCodec()
	client, err := gnet.NewClient(&clientHandler{codec: codec},
		gnet.WithLogger(logging.GetDefaultLogger()),
		gnet.WithReadBufferCap(1024),
		gnet.WithWriteBufferCap(1024))
	if err != nil {
		logging.Fatalf("Failed to create client: %v", err)
	}

	return &TcpClientGroup{client, broker.NewCodec()}
}

func (t *TcpClientGroup) Start() {
	t.clients.Start()
}

func (t *TcpClientGroup) Stop() {
	t.clients.Stop()
}

func (t *TcpClientGroup) NewClient() (gnet.Conn, error) {
	return t.clients.Dial("tcp", "127.0.0.1:5075")
}
