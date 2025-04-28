package main

import (
	"context"
	"encoding/json"
	"fmt"
	console "github.com/asynkron/goconsole"
	"github.com/magicnana999/im/api/kitex_gen/api"
	"github.com/magicnana999/im/broker"
	"github.com/magicnana999/im/pkg/id"
	"github.com/magicnana999/im/pkg/logger"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"go.uber.org/zap"
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

	h.htsHandler.StartTicking(&SendHeartbeat{user: user, handle: h.handler})

	//登录
	req := &api.LoginRequest{
		AppId:    "1201",
		UserSig:  "",
		Os:       "iOS",
		DeviceId: "",
	}
	packet := api.NewCommand(req)
	h.write(packet, user)

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

		if !packet.IsHeartbeat() {
			logging.Infof("%d read: %s", user.UserID, toJson(packet))
		}
		ret := h.handler.Handle(packet, user)
		if ret != nil {
			h.handler.Write(ret, user)
		}

		if packet.IsCommand() && packet.GetCommand().CommandType == api.CommandTypeUserLogin {
			h.curUsers.Store(user.UserID, user)
		}
	}

	return gnet.None
}

func (h *TcpServer) OnClose(c gnet.Conn, err error) gnet.Action {

	user := GetUser(c)
	if user != nil {
		user.IsClosed.Store(true)
		h.curUsers.Delete(user.UserID)

		loginUser := h.GetLoginUser()
		if loginUser != nil && loginUser.UserID == user.UserID {
			atomic.StorePointer(&h.curUser, nil)
			logging.Infof("%d current user closed", user.UserID)
		}
	}

	logging.Infof("%d Connection closed: %v", user.UserID, err)

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

func (h *TcpServer) connect(size int) {
	for i := 0; i < size; i++ {
		h.agent.Dial("tcp", "127.0.0.1:5075")
	}
}

func (h *TcpServer) login(size int) {
	for i := 0; i < size; i++ {

		h.connect(1)

		req := &api.LoginRequest{
			AppId:    "1201",
			UserSig:  "",
			Os:       "iOS",
			DeviceId: "",
		}
		packet := api.NewCommand(req)
		h.write(packet, h.GetLoginUser())
	}
}

func (h *TcpServer) write(ht *api.Packet, user *User) {
	h.handler.Write(ht, user)
}

func NewTcpServer(handler *PacketHandler, hts *HeartbeatServer) *TcpServer {

	codec := broker.NewCodec()
	eh := &TcpServer{codec: codec, handler: handler, htsHandler: hts}
	l := logger.NameWithOptions("tcp-client", zap.AddCallerSkip(2))
	client, err := gnet.NewClient(eh,
		gnet.WithLogger(&log{log: l}),
		gnet.WithReadBufferCap(8192),
		gnet.WithWriteBufferCap(8192))
	if err != nil {
		logging.Fatalf("Failed to create client: %v", err)
	}

	c := console.NewConsole(func(s string) {
	})

	c.Command("connslowly", func(text string) {

		s := strings.Trim(text, " ")
		if s != "" {
			if size, err := strconv.Atoi(s); err == nil {

				for i := 0; i < 10; i++ {
					time.Sleep(time.Second)
					eh.connect(size)
				}

			}
		} else {
			eh.connect(1)
		}

	})

	c.Command("conn", func(text string) {

		s := strings.Trim(text, " ")
		if s != "" {
			if size, err := strconv.Atoi(s); err == nil {
				eh.connect(size)
			}
		} else {
			eh.connect(1)
		}

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

type log struct {
	log *logger.Logger
}

func (l log) Debugf(format string, args ...any) {
	l.log.Debug(fmt.Sprintf(format, args...))
}

func (l log) Infof(format string, args ...any) {
	l.log.Info(fmt.Sprintf(format, args...))

}

func (l log) Warnf(format string, args ...any) {
	l.log.Warn(fmt.Sprintf(format, args...))

}

func (l log) Errorf(format string, args ...any) {
	l.log.Error(fmt.Sprintf(format, args...))

}

func (l log) Fatalf(format string, args ...any) {
	l.log.Fatal(fmt.Sprintf(format, args...))

}
