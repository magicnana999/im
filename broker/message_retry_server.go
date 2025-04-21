package broker

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
	"sync/atomic"
	"time"

	"github.com/magicnana999/im/api/kitex_gen/api"
	"github.com/magicnana999/im/broker/domain"
	"github.com/magicnana999/im/global"
	"github.com/magicnana999/im/pkg/jsonext"
	"github.com/magicnana999/im/pkg/logger"
	"github.com/magicnana999/im/pkg/timewheel"
)

func getOrDefaultMRSConfig(g *global.Config) *global.MRSConfig {
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
		s := fmt.Sprintf("timeout[%s] is less than interval[%s] ,set as %s",
			c.Timeout.String(),
			c.Interval.String(),
			c.Interval.String())
		logger.Named("mrs").Warn(s)
		c.Timeout = c.Interval
	}

	return c
}

type MessageRetryServer struct {
	tasks  sync.Map             // tasks stores message retry tasks by message ID.
	tw     *timewheel.Timewheel // tw schedules retry tasks.
	cfg    *global.MRSConfig    // cfg holds retry configuration.
	logger *Logger              // logger records server events.
	mw     *MessageWriter
	mr     *MessageResaver
}

func NewMessageRetryServer(g *global.Config, lc fx.Lifecycle) (*MessageRetryServer, error) {
	c := getOrDefaultMRSConfig(g)

	log := NewLogger("mrs", c.DebugMode)
	log.SrvInfo(string(jsonext.MarshalNoErr(c)), SrvLifecycle, nil)

	twc := &timewheel.Config{
		Tick:                time.Millisecond * 100,
		SlotCount:           30,
		MaxLengthOfEachSlot: 100_0000,
	}
	log.SrvInfo(string(jsonext.MarshalNoErr(twc)), SrvLifecycle, nil)

	twLogger := logger.NameWithOptions("timewheel-mrs", zap.IncreaseLevel(zapcore.InfoLevel))
	tw, err := timewheel.NewTimewheel(twc, twLogger, nil)
	if err != nil {
		log.SrvInfo("failed to init timewheel", SrvLifecycle, nil)
		return nil, err
	}

	s := &MessageRetryServer{
		tasks:  sync.Map{},
		tw:     tw,
		cfg:    c,
		logger: log,
		mw:     NewMessageWriter(NewCodec(), log),
		mr:     NewMessageResaver(),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return s.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return s.Stop(ctx)
		},
	})

	return s, nil
}

func (s *MessageRetryServer) Start(ctx context.Context) error {
	go func() {
		s.tw.Start(context.Background())
		s.logger.SrvInfo("timewheel-mrs started", SrvLifecycle, nil)
	}()
	return nil
}

func (s *MessageRetryServer) Stop(ctx context.Context) error {
	s.tw.Stop()
	s.logger.SrvInfo("timewheel-mrs stopped", SrvLifecycle, nil)
	s.logger.SrvInfo("resave the remaining message", SrvLifecycle, nil)
	go s.resaveMessages()
	return nil
}

func (s *MessageRetryServer) resaveMessages() {

	s.tasks.Range(func(key, value interface{}) bool {
		task, ok := value.(*messageRetryTask)
		if ok && task != nil {
			s.mr.Resave(task.m, task.uc)
		}
		return true
	})
}

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
	_, err := s.tw.Submit(task, s.cfg.Interval)
	return err
}

func (s *MessageRetryServer) Ack(messageID string) error {
	if messageID == "" {
		return errors.New("invalid message ID")
	}
	value, ok := s.tasks.LoadAndDelete(messageID)
	if ok {
		if task, ok := value.(*messageRetryTask); ok {
			task.isAckOK.Store(true)
		}
	}
	return nil
}

type messageRetryTask struct {
	uc              *domain.UserConn
	m               *api.Message
	isAckOK         atomic.Bool                                //是否已收到ACK
	firstSendSecond atomic.Int64                               //首次发送时间
	interval        time.Duration                              //重发间隔
	timeout         time.Duration                              //重发超时时间，超多此时间将不在重发
	writeFunc       func(*api.Message, *domain.UserConn) error //消息写入方法
	resaveFunc      func(*api.Message, *domain.UserConn)       //消息保存方法
}

func (t *messageRetryTask) Execute(now time.Time) timewheel.TaskResult {
	if t.isAckOK.Load() {
		return timewheel.Break
	}

	if time.Since(time.Unix(t.firstSendSecond.Load(), 0)) >= t.timeout {
		t.resaveFunc(t.m, t.uc)
		return timewheel.Break
	}

	if err := t.writeFunc(t.m, t.uc); err != nil {
		t.resaveFunc(t.m, t.uc)
		return timewheel.Break
	}

	return timewheel.Retry
}

func (s *MessageRetryServer) write(m *api.Message, uc *domain.UserConn) error {
	return s.mw.Write(m, uc)
}

func (s *MessageRetryServer) resave(m *api.Message, uc *domain.UserConn) {
	s.mr.Resave(m, uc)
}
