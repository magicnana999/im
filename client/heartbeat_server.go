package main

import (
	"context"
	"github.com/magicnana999/im/api/kitex_gen/api"
	"github.com/magicnana999/im/pkg/logger"
	"github.com/magicnana999/im/pkg/timewheel"
	"go.uber.org/zap"
	"runtime"
	"time"
)

type SendHeartbeat struct {
	user   *User
	handle *PacketHandler
}

func (s *SendHeartbeat) Execute(now time.Time) timewheel.TaskResult {
	if s.user.IsClosed.Load() {
		return timewheel.Break
	}
	ht := api.NewHeartbeatPacket(100)
	s.handle.Write(ht, s.user)
	return timewheel.Retry
}

type HeartbeatServer struct {
	tw       *timewheel.Timewheel
	interval time.Duration
}

func NewHeartbeatServer(interval time.Duration) *HeartbeatServer {
	twc := &timewheel.Config{
		SlotTick:      time.Millisecond * 500,
		SlotCount:     60,
		SlotMaxLength: 1_000,
		WorkerCount:   runtime.NumCPU() * 10,
	}

	log := logger.NameWithOptions("hts-client", zap.IncreaseLevel(zap.InfoLevel))

	tw, err := timewheel.NewTimewheel(twc, log, nil)
	if err != nil {
		logger.Named("hts").Fatal("hts start fail")
		return nil
	}

	hs := &HeartbeatServer{
		tw:       tw,
		interval: interval,
	}

	return hs
}

func (s *HeartbeatServer) Start() error {
	go s.tw.Start(context.Background())
	return nil
}

func (s *HeartbeatServer) Stop() error {
	s.tw.Stop()
	return nil
}

func (s *HeartbeatServer) StartTicking(fun *SendHeartbeat) (int, int64, error) {
	return s.tw.Submit(fun)
}
