package broker

import (
	"context"
	"errors"
	"fmt"
	"github.com/magicnana999/im/pkg/jsonext"
	"time"

	"github.com/magicnana999/im/broker/domain"
	"github.com/magicnana999/im/broker/holder"
	"github.com/magicnana999/im/define"
	"github.com/magicnana999/im/global"
	"github.com/magicnana999/im/pkg/logger"
	"github.com/magicnana999/im/pkg/timewheel"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	// errConnClosedAlready indicates that the client connection is already closed.
	errConnClosedAlready = errors.New("conn already closed")

	// errHeartbeatTimeout indicates that the heartbeat has timed out, requiring the connection to be closed.
	errHeartbeatTimeout = errors.New("heartbeat timeout")

	// errUCIsNil indicates that the user connection provided to Submit is nil.
	errUCIsNil = errors.New("uc is nil")
)

// HeartbeatServer manages client heartbeat detection and timeout handling.
// It runs in a separate goroutine and must be stopped explicitly using Stop.
type HeartbeatServer struct {
	userHolder *holder.UserHolder   // userHolder manages all client connections.
	tw         *timewheel.Timewheel // tw schedules heartbeat tasks using a timewheel.
	cfg        *global.HTSConfig    // cfg holds the heartbeat cmd_service configuration.
	logger     *Logger              // logger records cmd_service events.
}

// getOrDefaultHTSConfig returns the HTS configuration, prioritizing global configuration and applying defaults if necessary.
// It does not modify the input global.Config.
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
		logger.Named("hts").Warn(s, zap.String(define.OP, define.OpInit))
	}

	return c
}

// NewHeartbeatServer initializes a new HeartbeatServer.
// It uses global.Config for configuration and fx.Lifecycle for lifecycle management.
// It returns the configured HeartbeatServer instance and an error if initialization fails (e.g., timewheel creation fails).
func NewHeartbeatServer(g *global.Config, uh *holder.UserHolder, lc fx.Lifecycle) (*HeartbeatServer, error) {
	c := getOrDefaultHTSConfig(g)

	log := NewLogger("hts", c.DebugMode)
	log.InfoOrError(string(jsonext.MarshalNoErr(c)), "", define.OpInit, "", nil)

	twc := &timewheel.Config{
		Tick:                time.Second,
		SlotCount:           60,
		MaxLengthOfEachSlot: 100_0000,
	}
	log.InfoOrError(string(jsonext.MarshalNoErr(twc)), "", define.OpInit, "", nil)

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

// Start launches the heartbeat cmd_service in a separate goroutine.
// It returns an error if the timewheel fails to start.
func (s *HeartbeatServer) Start(ctx context.Context) error {
	go s.tw.Start(ctx)
	s.logger.InfoOrError("timewheel start", "", define.OpStart, "", nil)
	return nil
}

// Stop shuts down the heartbeat cmd_service.
// It stops the timewheel and returns an error if the shutdown fails.
func (s *HeartbeatServer) Stop(ctx context.Context) error {
	s.tw.Stop()
	s.logger.InfoOrError("timewheel stop", "", define.OpStart, "", nil)

	if s.logger != nil {
		s.logger.Close()
	}
	return nil
}

// Submit schedules a heartbeat task for the given user connection.
// It returns an error if the user connection is nil or the task submission fails.
func (s *HeartbeatServer) Submit(ctx context.Context, uc *domain.UserConn) error {
	if uc == nil {
		s.logger.Error("failed to submit", zap.Error(errUCIsNil))
		return errUCIsNil
	}

	if uc.IsClosed.Load() {
		s.logger.Error("failed to submit", zap.Error(errConnClosedAlready))
		return errConnClosedAlready
	}

	close := func() {
		s.userHolder.Close(ctx, uc)
		s.logger.DebugOrError("close", uc.Desc(), define.OpTicking, "", nil)
	}

	refresh := func() {
		err := s.userHolder.RefreshUser(ctx, uc)
		if err != nil {
			s.logger.DebugOrError("failed to refresh user,to closing", uc.Desc(), define.OpTicking, "", err)
			s.userHolder.Close(ctx, uc)
		} else {
			s.logger.DebugOrError("refresh", uc.Desc(), define.OpTicking, "", nil)
		}
	}

	task := &heartbeatTask{
		uc:          uc,
		interval:    s.cfg.Interval,
		timeout:     s.cfg.Timeout,
		closeFunc:   close,
		refreshFunc: refresh,
		logger:      s.logger,
	}

	_, err := s.tw.Submit(task, task.interval)
	s.logger.DebugOrError("failed to submit", uc.Desc(), define.OpSubmit, "", err)
	return err
}

// heartbeatTask represents a single heartbeat task for a user connection.
// It is executed periodically by the timewheel to check connection status.
type heartbeatTask struct {
	uc          *domain.UserConn // uc is the user connection.
	interval    time.Duration    // interval is the heartbeat interval.
	timeout     time.Duration    // timeout is the maximum allowed inactivity duration.
	closeFunc   func()           // closeFunc closes the connection on timeout.
	refreshFunc func()           // refreshFunc updates the connection's online status.
	logger      *Logger          // logger records task events.
}

// Execute is called by the timewheel to process the heartbeat task at the specified time.
// It checks if the connection is closed or timed out, refreshes the connection status if active,
// and returns the task result (Retry or Break) along with any error encountered.
func (t *heartbeatTask) Execute(now time.Time) (timewheel.TaskResult, error) {
	// Connection is already closed.

	if t.uc == nil {
		t.logger.DebugOrError("uc is nil", "", define.OpTicking, "", errUCIsNil)
		return timewheel.Break, errUCIsNil
	}

	if t.uc.IsClosed.Load() {
		t.logger.DebugOrError("uc is closed", t.uc.Desc(), define.OpTicking, "", errConnClosedAlready)
		return timewheel.Break, errConnClosedAlready
	}

	// Connection has timed out.
	if time.Since(time.Unix(t.uc.LastHeartbeat.Load(), 0)) >= t.timeout {
		t.logger.DebugOrError("timeout", t.uc.Desc(), define.OpTicking, "", nil)
		t.closeFunc()
		return timewheel.Break, errHeartbeatTimeout
	}

	// Refresh connection status.
	t.refreshFunc()

	// Continue scheduling the task.
	return timewheel.Retry, nil
}
