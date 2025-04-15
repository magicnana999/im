package broker

import (
	"context"
	"github.com/magicnana999/im/api/kitex_gen/api"
	"github.com/magicnana999/im/broker/domain"
	"github.com/magicnana999/im/errors"
	"github.com/magicnana999/im/global"
	"github.com/magicnana999/im/pkg/logger"
	"github.com/magicnana999/im/pkg/timewheel"
	"github.com/panjf2000/gnet/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"sync"
	"time"
)

const (
	MSS = "mss"
)

type MSSConfig struct {
	Interval int64 `yaml:"interval"` // 秒
	Timeout  int64 `yaml:"timeout"`  // 秒
}

type MessageSendServer struct {
	m        sync.Map
	tw       *timewheel.Timewheel
	interval int64
	timeout  int64
}

func NewMessageSendServer(lc fx.Lifecycle) *MessageSendServer {

	c := global.GetMSS()

	if c == nil {
		logger.Fatal("mss configuration not found",
			zap.String(logger.SCOPE, MSS),
			zap.String(logger.OP, Init))
	}

	tw, err := timewheel.NewTimeWheel(time.Millisecond*100, 30, nil)
	if err != nil {
		logger.Fatal("timewheel could not be open",
			zap.String(logger.SCOPE, MSS),
			zap.String(logger.OP, Init),
			zap.Error(err))
	}

	hs := &MessageSendServer{
		m:        sync.Map{},
		tw:       tw,
		interval: c.Interval,
		timeout:  c.Timeout,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("heartbeat server is ready",
				zap.String(logger.SCOPE, HTS),
				zap.String(logger.OP, Init))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if e := hs.Stop(); e != nil {
				logger.Error("heartbeat server could not close",
					zap.String(logger.SCOPE, HTS),
					zap.String(logger.OP, Close),
					zap.Error(e))
				return e
			} else {
				logger.Info("heartbeat server closed",
					zap.String(logger.SCOPE, HTS),
					zap.String(logger.OP, Close))
				return nil
			}
		},
	})

	return hs
}

func (s *MessageSendServer) Start(ctx context.Context) {
	go s.tw.Start(ctx)
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
	interval int64
	timeout  int64
}

func (t *messageSendTask) Execute(nowSecond int64) (timewheel.TaskResult, error) {
	if nowSecond-t.lastHeartbeat.Load() > (t.interval * 2) {
		t.c.Close()
		return timewheel.Break, errors.HeartbeatTimeout
	} else {
		t.userHolder.RefreshUser(t.ctx, t.uc)
		return timewheel.Retry, nil
	}
}
