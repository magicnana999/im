package broker

import (
	"context"
	"github.com/magicnana999/im/pkg/jsonext"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"runtime"
	"time"

	"github.com/magicnana999/im/global"
	"github.com/magicnana999/im/pkg/logger"
	"github.com/magicnana999/im/pkg/timewheel"
	"go.uber.org/fx"
)

type HeartbeatFunc func(now time.Time) timewheel.TaskResult

func (f HeartbeatFunc) Execute(now time.Time) timewheel.TaskResult {
	return f(now)
}

type HeartbeatServer struct {
	tw     *timewheel.Timewheel
	cfg    *global.TcpHeartbeatConfig
	logger *Logger
}

func getOrDefaultHTSConfig(g *global.Config) *global.TcpHeartbeatConfig {
	c := &global.TcpHeartbeatConfig{}
	if g != nil && g.TCP != nil && g.TCP.Heartbeat != nil {
		*c = *g.TCP.Heartbeat
	}

	if c.SlotTick <= 0 {
		c.SlotTick = time.Second
	}

	if c.SlotCount <= 0 {
		c.SlotCount = 30
	}

	if c.SlotMaxLength <= 0 {
		c.SlotMaxLength = 20000
	}

	if c.WorkerCount <= 0 {
		c.WorkerCount = runtime.NumCPU() * 50
	}

	if c.WorkerExpire <= 0 {
		c.WorkerExpire = time.Second * 10
	}

	return c
}

func NewHeartbeatServer(g *global.Config, lc fx.Lifecycle) (*HeartbeatServer, error) {
	c := getOrDefaultHTSConfig(g)

	log := NewLogger("hts", true)
	log.SrvInfo(string(jsonext.MarshalNoErr(c)), SrvLifecycle, nil)

	twc := &timewheel.Config{
		SlotTick:          c.SlotTick,
		SlotCount:         c.SlotCount,
		SlotMaxLength:     c.SlotMaxLength,
		WorkerCount:       c.WorkerCount,
		WorkerExpire:      c.WorkerExpire,
		WorkerNonBlocking: c.WorkerNonBlocking,
		WorkerPreAlloc:    c.WorkerPreAlloc,
	}

	twLogger := logger.NameWithOptions("hts", zap.IncreaseLevel(zapcore.InfoLevel))
	tw, err := timewheel.NewTimewheel(twc, twLogger, nil)
	if err != nil {
		return nil, err
	}

	hs := &HeartbeatServer{
		tw:     tw,
		cfg:    c,
		logger: log,
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
	go func() {
		s.tw.Start(context.Background())
		s.logger.SrvInfo("timewheel start", SrvLifecycle, nil)
	}()
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

func (s *HeartbeatServer) Ticking(fun HeartbeatFunc) (int, int64, error) {
	return s.tw.Submit(fun)
}
