package broker

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/magicnana999/im/api/kitex_gen/api"
	"github.com/magicnana999/im/broker/domain"
	"github.com/magicnana999/im/global"
	"github.com/magicnana999/im/pkg/jsonext"
	"go.uber.org/fx"
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
	uc      *domain.UserConn
}

// MessageSendServer 消息投递服务，用于主动给客户端投递消息，不用于收到客户端消息后的ack
type MessageSendServer struct {
	isRunning atomic.Bool
	cancel    context.CancelFunc
	ch        chan *messageSending //消息投递队列
	logger    *Logger
	mrs       *MessageRetryServer //消息重发服务
	mw        *PacketWriter       //消息写入服务
}

func getOrDefaultMSSConfig(g *global.Config) *global.MSSConfig {
	c := &global.MSSConfig{}
	if g != nil && g.MSS != nil {
		*c = *g.MSS
	}

	if c.MaxRemaining <= 0 || c.MaxRemaining > 100000 {
		c.MaxRemaining = 100000
	}

	return c
}

func NewMessageSendServer(g *global.Config, mrs *MessageRetryServer, lc fx.Lifecycle) (*MessageSendServer, error) {
	c := getOrDefaultMSSConfig(g)

	log := NewLogger("mss")
	log.SrvInfo(string(jsonext.MarshalNoErr(c)), SrvLifecycle, nil)

	mss := &MessageSendServer{
		ch:     make(chan *messageSending, c.MaxRemaining),
		logger: log,
		mrs:    mrs,
		mw:     NewPacketWriter(NewCodec(), log),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return mss.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return mss.Stop(ctx)
		},
	})
	return mss, nil
}

func (mss *MessageSendServer) Start(ctx context.Context) error {
	if mss.isRunning.CompareAndSwap(false, true) {
		ctx, cancel := context.WithCancel(context.Background())
		mss.cancel = cancel
		go func() {
			mss.logger.SrvInfo("message send loop started", SrvLifecycle, nil)

			for {
				select {
				case <-ctx.Done():
					return
				case ms := <-mss.ch:
					mss.write(ms.message, ms.uc)
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
		close(mss.ch)
		mss.logger.SrvInfo("message send loop stopped", SrvLifecycle, nil)

		mss.logger.SrvInfo("resave the remaining message", SrvLifecycle, nil)
		go mss.resaveMessages()
	}
	return nil
}

func (mss *MessageSendServer) resaveMessages() {
	for ms := range mss.ch {
		mss.resave(ms.message, ms.uc)
	}
}

// Send 外部调用
func (mss *MessageSendServer) Send(m *api.Message, uc *domain.UserConn) error {
	if m == nil {
		return invalidMessage
	}
	if uc == nil {
		return invalidUserConn
	}
	if uc.IsClosed.Load() {
		return userConnClosed
	}
	if !mss.isRunning.Load() {
		return mssNotRunning
	}

	select {
	case mss.ch <- &messageSending{message: m, uc: uc}:
		return nil
	default:
		return channelFull
	}
}

// write 成功后开始消息重发逻辑，失败后直接写入离线
func (mss *MessageSendServer) write(m *api.Message, uc *domain.UserConn) {
	if err := mss.mw.Write(m.Wrap(), uc); err != nil {
		mss.resave(m, uc)
	} else {
		mss.submit(m, uc)
	}
}

func (mss *MessageSendServer) submit(ms *api.Message, uc *domain.UserConn) {
	if err := mss.mrs.Submit(ms, uc, time.Now().Unix()); err != nil {
		mss.logger.PktDebug("failed to submit mrs,resave it", uc.Desc(), ms.MessageId, nil, PacketTracking, nil)
		mss.resave(ms, uc)
	}
}

func (mss *MessageSendServer) resave(ms *api.Message, uc *domain.UserConn) {
	mss.logger.PktDebug("resave message", uc.Desc(), ms.MessageId, nil, PacketTracking, nil)
	fmt.Print(jsonext.PbMarshalNoErr(ms))
}
