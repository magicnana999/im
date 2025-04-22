package main

import (
	"context"
	"encoding/json"
	"fmt"
	console "github.com/asynkron/goconsole"
	"github.com/magicnana999/im/api/kitex_gen/api"
	"github.com/magicnana999/im/broker"
	"github.com/magicnana999/im/pkg/id"
	"github.com/magicnana999/im/pkg/timewheel"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

var key string = "USER"

type TcpServer struct {
	agent      *gnet.Client
	codec      *broker.Codec
	handler    *PacketHandler
	htsHandler *HeartbeatServer
	curUser    unsafe.Pointer
	curUsers   sync.Map
	console    *console.Console
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

func (h *TcpServer) GetLoginUser() *User {
	v := atomic.LoadPointer(&h.curUser)
	if v != nil {
		return (*User)(v)
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

	atomic.StorePointer(&h.curUser, unsafe.Pointer(user))

	ctx := context.WithValue(context.Background(), key, user)
	c.SetContext(ctx)

	h.htsHandler.StartTicking(func(now time.Time) timewheel.TaskResult {
		if user.IsClosed.Load() {
			return timewheel.Break
		}
		ht := api.NewHeartbeatPacket(100)
		h.handler.Write(ht, c, user)
		return timewheel.Retry
	}, time.Second*20)

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

		logging.Infof("%d read: %s", user.UserID, toJson(packet))
		ret := h.handler.Handle(packet, user)
		if ret != nil {
			h.handler.Write(ret, c, user)
		}

		if packet.IsCommand() && packet.GetCommand().CommandType == api.CommandTypeUserLogin {
			h.curUsers.Store(user.UserID, user)
		}
	}

	return gnet.None
}

func (h *TcpServer) OnClose(c gnet.Conn, err error) gnet.Action {
	logging.Infof("Connection closed: %v", err)
	return gnet.None
}

func (h *TcpServer) Start() {
	go h.agent.Start()
	go h.htsHandler.Start()
	h.console.Run()

}

func (h *TcpServer) Stop() {
	h.agent.Stop()
	h.htsHandler.Stop()
}

func (h *TcpServer) connect() (gnet.Conn, error) {
	return h.agent.Dial("tcp", "127.0.0.1:5075")
}

func (h *TcpServer) write(ht *api.Packet, user *User) {
	h.handler.Write(ht, user.Writer, user)
}

func NewTcpServer(handler *PacketHandler, hts *HeartbeatServer) *TcpServer {

	codec := broker.NewCodec()
	eh := &TcpServer{codec: codec, handler: handler, htsHandler: hts}
	client, err := gnet.NewClient(eh,
		gnet.WithLogger(logging.GetDefaultLogger()),
		gnet.WithReadBufferCap(1024),
		gnet.WithWriteBufferCap(1024))
	if err != nil {
		logging.Fatalf("Failed to create client: %v", err)
	}

	c := console.NewConsole(func(s string) {
	})

	c.Command("conn", func(text string) {
		eh.connect()
	})

	c.Command("login", func(text string) {
		req := &api.LoginRequest{
			AppId:    "1201",
			UserSig:  "",
			Os:       "iOS",
			DeviceId: "",
		}

		packet := api.NewCommand(req)
		eh.write(packet, eh.GetLoginUser())
	})

	c.Command("show", func(text string) {
		if text == "users" {
			eh.curUsers.Range(func(key, value interface{}) bool {
				fmt.Println(key, MarshalNoError(value))
				return true
			})
		}

		if text == "user" {
			fmt.Println(MarshalNoError(eh.GetLoginUser()))
		}
	})

	c.Command("changeuser", func(text string) {

		id, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			fmt.Errorf("change current user failed: %v", err)
		}

		v, ok := eh.curUsers.Load(id)
		if !ok || v == nil {
			fmt.Errorf("target user not found: %s", text)
		}

		u, ok := v.(*User)
		if !ok || u == nil {
			fmt.Errorf("target user is invalid: %s", text)
		}

		atomic.StorePointer(&eh.curUser, unsafe.Pointer(u))
	})

	c.Command("sendtouser", func(text string) {
		s := strings.Split(text, "|")
		if len(s) != 2 {
			return
		}

		toUserId, err := strconv.ParseInt(s[0], 10, 64)
		if err != nil {
			return
		}

		appId := (*User)(atomic.LoadPointer(&eh.curUser)).AppID

		textBody := &api.Text{
			Text: s[1],
		}
		m := api.NewMessage(
			(*User)(atomic.LoadPointer(&eh.curUser)).UserID,
			toUserId,
			0,
			id.SnowflakeID(),
			appId,
			"",
			textBody)

		eh.write(m.Wrap(), (*User)(atomic.LoadPointer(&eh.curUser)))
	})

	eh.console = c
	eh.agent = client
	return eh
}

func MarshalNoError(any any) string {
	bs, err := json.Marshal(any)
	if err != nil {
		return ""
	}
	return string(bs)
}
