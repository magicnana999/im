package broker

import (
	"context"
	"errors"
	"fmt"
	"github.com/magicnana999/im/broker/msg_service"
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
	"go.uber.org/zap"
)

var (
	retryTimeout = errors.New("retry timeout")
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
		c.Timeout = c.Interval
		s := fmt.Sprintf("invalid timeout,set as %d", c.Timeout)
		logger.Named("mrs").Warn(s, zap.String(define.OP, define.OpInit))
	}

	return c
}

type MessageRetryServer struct {
	tasks  sync.Map             // tasks stores message retry tasks by message ID.
	tw     *timewheel.Timewheel // tw schedules retry tasks.
	cfg    *global.MRSConfig    // cfg holds retry configuration.
	logger *Logger              // logger records server events.
	mw     *msg_service.MessageWriter
	mr     *msg_service.MessageResaver
}

func NewMessageRetryServer(g *global.Config) (*MessageRetryServer, error) {
	c := getOrDefaultMRSConfig(g)

	log := NewLogger("mrs", c.DebugMode)
	log.InfoOrError(string(jsonext.MarshalNoErr(c)), "", define.OpInit, "", nil)

	twc := &timewheel.Config{
		Tick:                time.Millisecond * 100,
		SlotCount:           30,
		MaxLengthOfEachSlot: 100_0000,
	}
	log.InfoOrError(string(jsonext.MarshalNoErr(twc)), "", define.OpInit, "", nil)

	tw, err := timewheel.NewTimewheel(twc, logger.Named("timewheel-mrs"), nil)
	if err != nil {
		log.InfoOrError("failed to new timewheel", "", define.OpInit, "", err)
		return nil, err
	}

	return &MessageRetryServer{
		tasks:  sync.Map{},
		tw:     tw,
		cfg:    c,
		logger: log,
		mw:     msg_service.NewMessageWriter(NewCodec(), log),
		mr:     msg_service.NewMessageResaver(log),
	}, nil
}

func (s *MessageRetryServer) Start(ctx context.Context) error {
	go s.tw.Start(ctx)
	s.logger.InfoOrError("timewheel started", "", define.OpStart, "", nil)
	return nil
}

func (s *MessageRetryServer) Stop(ctx context.Context) error {
	s.tw.Stop()
	s.logger.InfoOrError("timewheel stopped", "", define.OpStop, "", nil)

	txt := fmt.Sprintf("remaining messages: could not be know")
	s.logger.InfoOrError(txt, "", define.OpStop, "", nil)
	s.resaveMessages()
	return nil
}

func (s *MessageRetryServer) resaveMessages() {

	s.tasks.Range(func(key, value interface{}) bool {
		messageId := key.(string)
		a, ok := s.tasks.Load(messageId)
		if ok && a != nil {
			task, ok := a.(*messageRetryTask)
			if ok && task != nil {
				s.mr.Resave(task.m, task.uc)
			}
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

func (s *MessageRetryServer) retryNoMore(messageID string) error {
	if messageID == "" {
		return errors.New("invalid message ID")
	}
	value, ok := s.tasks.LoadAndDelete(messageID)
	if ok && value != nil {
		if task, ok := value.(*messageRetryTask); ok && task != nil {
			task.IsRetryOK.Store(false)
		}
	}
	return nil
}

type messageRetryTask struct {
	uc              *domain.UserConn
	m               *api.Message
	isAckOK         atomic.Bool                                //是否已收到ACK
	IsRetryOK       atomic.Bool                                //是否需要再次重试
	firstSendSecond atomic.Int64                               //首次发送时间
	interval        time.Duration                              //重发间隔
	timeout         time.Duration                              //重发超时时间，超多此时间将不在重发
	writeFunc       func(*api.Message, *domain.UserConn) error //消息写入方法
	resaveFunc      func(*api.Message, *domain.UserConn)       //消息保存方法
}

func (t *messageRetryTask) Execute(now time.Time) (timewheel.TaskResult, error) {
	if t.isAckOK.Load() {
		return timewheel.Break, nil
	}

	if !t.IsRetryOK.Load() {
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

func (s *MessageRetryServer) write(m *api.Message, uc *domain.UserConn) error {
	return s.mw.Write(m, uc)
}

func (s *MessageRetryServer) resave(m *api.Message, uc *domain.UserConn) {
	s.mr.Resave(m, uc)
}
