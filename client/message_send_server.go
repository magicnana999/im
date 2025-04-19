package main

import (
	"context"
	"errors"
	"github.com/magicnana999/im/api/kitex_gen/api"
	"sync/atomic"
)

var (
	invalidMessage  = errors.New("invalid message")
	invalidUserConn = errors.New("invalid user conn")
	userConnClosed  = errors.New("user conn closed")
	mssNotRunning   = errors.New("server not running")
	channelFull     = errors.New("channel full")
)

type messageSending struct {
	message *api.Message
	user    *User
}

// MessageSendServer 消息投递服务，用于主动给客户端投递消息，不用于收到客户端消息后的ack
type MessageSendServer struct {
	isRunning atomic.Bool
	cancel    context.CancelFunc
	ch        chan *messageSending //消息投递队列
	mw        *MessageWriter       //消息写入服务
}

func NewMessageSendServer() (*MessageSendServer, error) {

	mss := &MessageSendServer{
		ch: make(chan *messageSending, 10000),
		mw: NewMessageWriter(),
	}

	return mss, nil
}

func (mss *MessageSendServer) Start(ctx context.Context) error {
	if mss.isRunning.CompareAndSwap(false, true) {
		ctx, cancel := context.WithCancel(context.Background())
		mss.cancel = cancel
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case ms := <-mss.ch:
					mss.write(ms.message, ms.user)
				}
			}
		}()
	}
	return nil
}

// Stop 停服后，把没有收到ack的消息保存到离线
func (mss *MessageSendServer) Stop(ctx context.Context) error {
	if mss.isRunning.CompareAndSwap(true, false) {
		mss.cancel()
	}
	return nil
}

// Send 外部调用
func (mss *MessageSendServer) Send(m *api.Message, user *User) error {
	if user.IsClosed.Load() {
		return userConnClosed
	}
	if !mss.isRunning.Load() {
		return mssNotRunning
	}

	select {
	case mss.ch <- &messageSending{message: m, user: user}:
		return nil
	default:
		return channelFull
	}
}

// write 成功后开始消息重发逻辑，失败后直接写入离线
func (mss *MessageSendServer) write(m *api.Message, uc *User) {
	mss.mw.Write(m, uc)
}
