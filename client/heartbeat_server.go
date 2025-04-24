package main

import (
	"context"
	"github.com/magicnana999/im/pkg/logger"
	"github.com/magicnana999/im/pkg/timewheel"
	"time"
)

type HeartbeatFunc func(now time.Time) timewheel.TaskResult

func (f HeartbeatFunc) Execute(now time.Time) timewheel.TaskResult {
	return f(now)
}

type HeartbeatServer struct {
	tw       *timewheel.Timewheel
	interval time.Duration
}

func NewHeartbeatServer(interval time.Duration) *HeartbeatServer {
	twc := &timewheel.Config{
		Tick:                time.Second,
		SlotCount:           10,
		MaxLengthOfEachSlot: 1_000_000,
	}

	tw, err := timewheel.NewTimewheel(twc, nil, nil)
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

func (s *HeartbeatServer) StartTicking(fun HeartbeatFunc, interval time.Duration) (int, error) {
	return s.tw.Submit(fun, interval)
}
