package broker

import (
	"context"
	"github.com/magicnana999/im/api/kitex_gen/api"
	"github.com/magicnana999/im/broker/domain"
	"github.com/magicnana999/im/define"
	"github.com/magicnana999/im/errors"
	"github.com/magicnana999/im/global"
	"github.com/magicnana999/im/pkg/logger"
	"github.com/magicnana999/im/pkg/timewheel"
	"github.com/panjf2000/gnet/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"sync"
	"sync/atomic"
	"time"
)

func getOrDefaultMSSConfig(g *global.Config) *global.MSSConfig {
	c := &global.MSSConfig{}
	if g != nil && g.MSS != nil {
		*c = *g.MSS
	}

	if c.Interval <= 0 {
		c.Interval = time.Second * 1
	}

	if c.Timeout <= 0 {
		c.Timeout = time.Second * 2
	}

	return c
}

type MessageSendServer struct {
	m      sync.Map
	tw     *timewheel.Timewheel
	cfg    *global.MSSConfig
	logger *logger.Logger
}

func NewMessageSendServer(g *global.Config, lc fx.Lifecycle) (*MessageSendServer, error) {

	log := logger.Named("mss")

	c := getOrDefaultMSSConfig(g)

	twc := &timewheel.Config{
		Tick:                time.Millisecond * 100,
		SlotCount:           30,
		MaxLengthOfEachSlot: 100_0000,
	}

	tw, err := timewheel.NewTimeWheel(twc, log, nil)
	if err != nil {
		log.Error("failed to init timewheel",
			zap.String(define.OP, define.OpInit),
			zap.Error(err))
		return nil, err
	}

	hs := &MessageSendServer{
		m:      sync.Map{},
		tw:     tw,
		cfg:    c,
		logger: log,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if e := hs.Start(ctx); e != nil {
				log.Error("mss start failed", zap.String(define.OP, define.OpStart), zap.Error(e))
				return e
			}
			log.Info("mss is started", zap.String(define.OP, define.OpStart))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if e := hs.Stop(); e != nil {
				log.Error("mss cloud not close",
					zap.String(define.OP, define.OpClose),
					zap.Error(e))
				return e
			} else {
				log.Info("mss is closed",
					zap.String(define.OP, define.OpClose))
				return nil
			}
		},
	})

	return hs, nil
}

func (s *MessageSendServer) Start(ctx context.Context) error {
	go s.tw.Start(ctx)
	return nil
}

func (s *MessageSendServer) Stop() error {
	s.tw.Stop()
	return nil
}

func (s *MessageSendServer) Send(ctx context.Context, c gnet.Conn, uc *domain.UserConnection, m *api.Message) error {
	task := &messageSendTask{
		c:        c,
		uc:       uc,
		msg:      m,
		interval: s.interval,
		timeout:  s.timeout,
	}

	_, err := s.tw.Submit(task, s.interval)
	return err

}

type messageSendTask struct {
	c        gnet.Conn
	uc       *domain.UserConnection
	msg      *api.Message
	isOK     atomic.Bool
	lastSend atomic.Int64
	interval time.Duration
	timeout  time.Duration
}

func (t *messageSendTask) Execute(now time.Time) (timewheel.TaskResult, error) {

	if t.isOK.Load() {
		return timewheel.Break, nil
	}

	t.c.Write()
	if now.UnixMilli()-t.lastSend.Load() > (t.timeout.Milliseconds()) {
		t.c.Close()
		return timewheel.Break, errors.HeartbeatTimeout
	} else {
		t.userHolder.RefreshUser(t.ctx, t.uc)
		return timewheel.Retry, nil
	}
}
