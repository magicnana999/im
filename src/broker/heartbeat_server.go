package broker

import (
	"context"
	"errors"
	"fmt"
	"github.com/magicnana999/im/pkg/jsonext"
	"time"

	"github.com/magicnana999/im/broker/holder"
	"github.com/magicnana999/im/global"
	"github.com/magicnana999/im/pkg/logger"
	"github.com/magicnana999/im/pkg/timewheel"
	"go.uber.org/fx"
)

var (
	errConnClosedAlready = errors.New("conn already closed")

	errHeartbeatTimeout = errors.New("heartbeat timeout")

	errUCIsNil = errors.New("uc is nil")
)

type HeartbeatFunc func(now time.Time) timewheel.TaskResult

func (f HeartbeatFunc) Execute(now time.Time) timewheel.TaskResult {
	return f(now)
}

type HeartbeatServer struct {
	userHolder *holder.UserHolder   // userHolder manages all client connections.
	tw         *timewheel.Timewheel // tw schedules heartbeat tasks using a timewheel.
	cfg        *global.HTSConfig    // cfg holds the heartbeat cmd_service configuration.
	logger     *Logger              // logger records cmd_service events.
}

func getOrDefaultHTSConfig(g *global.Config) *global.HTSConfig {
	c := &global.HTSConfig{}
	if g != nil && g.HTS != nil {
		*c = *g.HTS
	}

	if c.Interval <= 0 {
		c.Interval = time.Second * 30 // Default heartbeat interval is 30 seconds.
	}

	if c.Timeout <= 0 {
		c.Timeout = time.Second * 60 // Default timeout is 60 seconds.
	}

	if c.Timeout < c.Interval {
		c.Timeout = c.Interval
		s := fmt.Sprintf("invalid timeout,set as %d", c.Timeout)
		logger.Named("hts").Warn(s)
	}

	return c
}

func NewHeartbeatServer(g *global.Config, uh *holder.UserHolder, lc fx.Lifecycle) (*HeartbeatServer, error) {
	c := getOrDefaultHTSConfig(g)

	log := NewLogger("hts", c.DebugMode)
	log.SrvInfo(string(jsonext.MarshalNoErr(c)), SrvLifecycle, nil)

	twc := &timewheel.Config{
		Tick:                time.Second,
		SlotCount:           60,
		MaxLengthOfEachSlot: 100_0000,
	}
	log.SrvInfo(string(jsonext.MarshalNoErr(twc)), SrvLifecycle, nil)

	tw, err := timewheel.NewTimewheel(twc, logger.Named("timewheel-hts"), nil)
	if err != nil {
		return nil, err
	}

	hs := &HeartbeatServer{
		userHolder: uh,
		tw:         tw,
		cfg:        c,
		logger:     log,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return hs.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return hs.Stop(ctx)
		},
	})

	return hs, nil
}

func (s *HeartbeatServer) Start(ctx context.Context) error {
	go s.tw.Start(ctx)
	s.logger.SrvInfo("timewheel start", SrvLifecycle, nil)
	return nil
}

func (s *HeartbeatServer) Stop(ctx context.Context) error {
	s.tw.Stop()
	s.logger.SrvInfo("timewheel stop", SrvLifecycle, nil)

	if s.logger != nil {
		s.logger.Close()
	}
	return nil
}

func (s *HeartbeatServer) StartTicking(fun HeartbeatFunc, interval time.Duration) (int, error) {
	return s.tw.Submit(fun, interval)
}
