package broker

import (
	"context"
	"github.com/magicnana999/im/broker/domain"
	"github.com/magicnana999/im/broker/holder"
	"github.com/magicnana999/im/errors"
	"github.com/magicnana999/im/global"
	"github.com/magicnana999/im/pkg/logger"
	"github.com/magicnana999/im/pkg/timewheel"
	"github.com/panjf2000/gnet/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"sync/atomic"
	"time"
)

const (
	HTS   = "hts"
	Init  = "init"
	Close = "close"
)

type HeartbeatServer struct {
	userHolder *holder.UserHolder
	tw         *timewheel.Timewheel
	interval   int64
	timeout    int64
}

func getOrDefaultHTSConfig(g *global.Config) *global.HTSConfig {
	c := &global.HTSConfig{}
	if g != nil && g.HTS != nil {
		*c = *g.HTS
	}

	if c.Interval == 0 {
		c.Interval = time.Second * 30
	}

	if c.Timeout == 0 {
		c.Timeout = time.Second * 60
	}

	return c
}

func NewHeartbeatServer(g *global.Config, uh *holder.UserHolder, lc fx.Lifecycle) (*HeartbeatServer, error) {

	log := logger.Named("hts")

	c := getOrDefaultHTSConfig(g)

	tw, err := timewheel.NewTimeWheel(time.Second, 60, nil)
	if err != nil {
		logger.Fatal("timewheel could not be open",
			zap.String(logger.SCOPE, HTS),
			zap.String(logger.OP, Init),
			zap.Error(err))
	}

	hs := &HeartbeatServer{
		userHolder: uh,
		tw:         tw,
		interval:   c.Interval,
		timeout:    c.Timeout,
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

func (s *HeartbeatServer) Start(ctx context.Context) {
	go s.tw.Start(ctx)
}

func (s *HeartbeatServer) Stop() error {
	s.tw.Stop()
	return nil
}

func (s *HeartbeatServer) Submit(ctx context.Context, c gnet.Conn, uc *domain.UserConnection) error {
	task := &heartbeatTask{
		c:             c,
		uc:            uc,
		lastHeartbeat: uc.LastHeartbeat,
		interval:      s.interval,
		timeout:       s.timeout,
		userHolder:    s.userHolder,
		ctx:           ctx,
	}

	_, err := s.tw.Submit(task, s.interval)
	return err

}

type heartbeatTask struct {
	c             gnet.Conn
	uc            *domain.UserConnection
	lastHeartbeat atomic.Int64
	interval      int64
	timeout       int64
	userHolder    *holder.UserHolder
	ctx           context.Context
}

func (t *heartbeatTask) Execute(nowSecond int64) (timewheel.TaskResult, error) {
	if nowSecond-t.lastHeartbeat.Load() > (t.interval * 2) {
		t.c.Close()
		return timewheel.Break, errors.HeartbeatTimeout
	} else {
		t.userHolder.RefreshUser(t.ctx, t.uc)
		return timewheel.Retry, nil
	}
}
