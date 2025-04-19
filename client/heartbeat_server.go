package main

import (
	"context"
	"time"

	"github.com/magicnana999/im/pkg/timewheel"
)

type HeartbeatFunc func(now time.Time) timewheel.TaskResult

func (f HeartbeatFunc) Execute(now time.Time) timewheel.TaskResult {
	return f(now)
}

type HeartbeatServer struct {
	tw       *timewheel.Timewheel // tw schedules heartbeat tasks using a timewheel.
	interval time.Duration
}

func NewHeartbeatServer(interval time.Duration) (*HeartbeatServer, error) {
	twc := &timewheel.Config{
		Tick:                time.Second,
		SlotCount:           60,
		MaxLengthOfEachSlot: 100_0000,
	}

	tw, err := timewheel.NewTimewheel(twc, nil, nil)
	if err != nil {
		return nil, err
	}

	hs := &HeartbeatServer{
		tw:       tw,
		interval: interval,
	}

	return hs, nil
}

func (s *HeartbeatServer) Start(ctx context.Context) error {
	go s.tw.Start(context.Background())
	return nil
}

func (s *HeartbeatServer) Stop(ctx context.Context) error {
	s.tw.Stop()
	return nil
}

func (s *HeartbeatServer) StartTicking(fun HeartbeatFunc, interval time.Duration) (int, error) {
	return s.tw.Submit(fun, interval)
}
