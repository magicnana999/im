package broker

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/magicnana999/im/api/kitex_gen/api"
	"github.com/magicnana999/im/broker/domain"
	"github.com/magicnana999/im/define"
	"github.com/magicnana999/im/global"
	"github.com/magicnana999/im/pkg/jsonext"
	"github.com/magicnana999/im/pkg/logger"
	"github.com/magicnana999/im/pkg/timewheel"
	bb "github.com/panjf2000/gnet/v2/pkg/pool/bytebuffer"
	"go.uber.org/zap"
)

// Error definitions for MessageRetryServer.
var (
	invalidTimeout = errors.New("timeout must be greater than or equal to interval")
	connIsClosed   = errors.New("connection is closed")
	retryTimeout   = errors.New("retry timeout")
	writeError     = errors.New("write error")
)

// getOrDefaultMRSConfig returns the MRS configuration, prioritizing global configuration and applying defaults if necessary.
// It returns an error if Timeout is less than Interval.
func getOrDefaultMRSConfig(g *global.Config) (*global.MRSConfig, error) {
	c := &global.MRSConfig{}
	if g != nil && g.MRS != nil {
		*c = *g.MRS
	}

	if c.Interval <= 0 {
		c.Interval = time.Second
	}
	if c.Timeout <= 0 {
		c.Timeout = time.Second * 2
	}

	if c.Timeout < c.Interval {
		logger.Named("mrs").Error("invalid timeout",
			zap.String(define.OP, define.OpInit),
			zap.Error(invalidTimeout))
		return nil, invalidTimeout
	}

	return c, nil
}

// MessageRetryServer manages message retries using a time wheel.
// It stores tasks in a concurrent map and supports submitting, acknowledging, and resaving messages.
type MessageRetryServer struct {
	tasks  sync.Map                             // tasks stores message retry tasks by message ID.
	codec  codec                                // codec encodes messages for writing.
	tw     *timewheel.Timewheel                 // tw schedules retry tasks.
	cfg    *global.MRSConfig                    // cfg holds retry configuration.
	logger *Logger                              // logger records server events.
	resave func(*api.Message, *domain.UserConn) // resave handles failed messages.
}

// NewMessageRetryServer creates a new MessageRetryServer with the specified configuration and resave function.
// It initializes a time wheel for scheduling retries and returns an error if configuration or time wheel creation fails.
func NewMessageRetryServer(g *global.Config, resave func(*api.Message, *domain.UserConn)) (*MessageRetryServer, error) {
	c, err := getOrDefaultMRSConfig(g)
	if err != nil {
		return nil, err
	}

	log := NewLogger("mrs", c.DebugMode)
	log.Info(string(jsonext.MarshalNoErr(c)), "", define.OpInit, "", nil)

	twc := &timewheel.Config{
		Tick:                time.Millisecond * 100,
		SlotCount:           30,
		MaxLengthOfEachSlot: 100_0000,
	}
	log.Info(string(jsonext.MarshalNoErr(twc)), "", define.OpInit, "", nil)

	tw, err := timewheel.NewTimewheel(twc, logger.Named("timewheel-mrs"), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create time wheel: %w", err)
	}

	return &MessageRetryServer{
		tasks:  sync.Map{},
		codec:  newCodec(),
		tw:     tw,
		cfg:    c,
		logger: log,
		resave: resave,
	}, nil
}

// Start launches the time wheel in a separate goroutine.
// It returns nil, as the time wheel handles its own errors.
func (s *MessageRetryServer) Start(ctx context.Context) error {
	go s.tw.Start(ctx)
	s.logger.Info("timewheel started", "", define.OpStart, "", nil)
	return nil
}

// Stop shuts down the time wheel.
// It returns nil, as the time wheel handles its own cleanup.
func (s *MessageRetryServer) Stop(ctx context.Context) error {
	s.tw.Stop()
	s.logger.Info("timewheel stopped", "", define.OpStop, "", nil)

	txt := fmt.Sprintf("remaining messages: could not be know")
	s.logger.Info(txt, "", define.OpStop, "", nil)
	s.ResaveMessages()
	return nil
}

// ResaveMessages processes unsent messages in the channel.
// It calls onCloseFunc for each message and logs errors.
func (s *MessageRetryServer) ResaveMessages() {

	s.tasks.Range(func(key, value interface{}) bool {
		messageId := key.(string)
		a, ok := s.tasks.Load(messageId)
		if ok && a != nil {
			task, ok := a.(*messageRetryTask)
			if ok && task != nil {
				s.resave(task.m, task.uc)
			}
		}
		return true

	})
}

// Submit adds a message retry task to the server.
// It stores the task in the map and schedules it with the time wheel.
// It returns an error if the message or connection is invalid.
func (s *MessageRetryServer) Submit(m *api.Message, uc *domain.UserConn, firstSend int64) error {
	if m == nil || m.MessageId == "" {
		return errors.New("invalid message")
	}
	if uc == nil || uc.IsClosed.Load() {
		return errors.New("invalid or closed connection")
	}

	task := &messageRetryTask{
		uc:         uc,
		m:          m,
		interval:   s.cfg.Interval,
		timeout:    s.cfg.Timeout,
		writeFunc:  s.write,
		resaveFunc: s.resave,
	}
	task.firstSendSecond.Store(firstSend)

	s.tasks.Store(m.MessageId, task)
	s.tw.Submit(task, s.cfg.Interval)
	return nil
}

// Ack acknowledges a message, marking it as successfully processed.
// It removes the task from the map and sets its acknowledgment flag.
func (s *MessageRetryServer) Ack(messageId string) error {
	if messageId == "" {
		return errors.New("invalid message ID")
	}
	value, ok := s.tasks.LoadAndDelete(messageId)
	if ok {
		if task, ok := value.(*messageRetryTask); ok {
			task.isAckOK.Store(true)
		}
	}
	return nil
}

// messageRetryTask represents a retry task for a message.
// It tracks the message, connection, and retry state.
type messageRetryTask struct {
	uc              *domain.UserConn                           // uc is the user connection.
	m               *api.Message                               // m is the message to retry.
	isAckOK         atomic.Bool                                // isAckOK indicates if the message is acknowledged.
	firstSendSecond atomic.Int64                               // firstSendSecond is the Unix timestamp of the first send.
	interval        time.Duration                              // interval is the retry interval.
	timeout         time.Duration                              // timeout is the maximum retry duration.
	writeFunc       func(*api.Message, *domain.UserConn) error // writeFunc writes the message.
	resaveFunc      func(*api.Message, *domain.UserConn)       // resaveFunc handles failed messages.
}

// Execute runs the retry task at the specified time.
// It returns Break if the task is acknowledged, timed out, or failed; otherwise, it returns Retry.
// It returns an error for timeout or write failures.
func (t *messageRetryTask) Execute(now time.Time) (timewheel.TaskResult, error) {
	if t.isAckOK.Load() {
		return timewheel.Break, nil
	}

	if time.Since(time.Unix(t.firstSendSecond.Load(), 0)) >= t.timeout {
		t.resaveFunc(t.m, t.uc)
		return timewheel.Break, retryTimeout
	}

	if err := t.writeFunc(t.m, t.uc); err != nil {
		t.resaveFunc(t.m, t.uc)
		return timewheel.Break, err
	}

	return timewheel.Retry, nil
}

// write encodes and writes a message to the user connection.
// It returns an error if the connection is closed or writing fails.
func (s *MessageRetryServer) write(m *api.Message, uc *domain.UserConn) error {
	if uc.IsClosed.Load() {
		return connIsClosed
	}

	buffer, err := s.codec.encode(m.Wrap())
	defer bb.Put(buffer)

	if err != nil {
		return err
	}

	total := buffer.Len()
	sent := 0
	for sent < total {
		n, err := uc.Writer.Write(buffer.Bytes()[sent:])
		if err != nil {
			return err
		}
		sent += n
	}
	s.logger.Debug("write ok", uc.Desc(), "write", m.MessageId, nil)
	return nil
}
