package timewheel

import (
	"context"
	"fmt"
	"github.com/magicnana999/im/pkg/jsonext"
	"github.com/magicnana999/im/pkg/logger"
	"github.com/magicnana999/im/pkg/queue"
	"github.com/panjf2000/ants/v2"
	"go.uber.org/zap"
	"runtime"
	"sync/atomic"
	"time"
)

type TaskResult string

const (
	Retry TaskResult = "retry" //需要再次调度任务
	Break TaskResult = "break" //退出，不需要再次调度
)

// Task 任务接口，所有被加入的任务必须实现此接口
type Task interface {
	Execute(now time.Time) TaskResult //不能有error，内部处理所有可能的error
}

type Config struct {
	Tick                time.Duration `yaml:"tick",json:"tick"`                               //时间间隔
	SlotCount           int           `yaml:"slotCount",json:"slotCount"`                     //槽数量
	MaxLengthOfEachSlot int           `yaml:"maxLengthOfEachSlot",json:"maxLengthOfEachSlot"` //每个槽的最大长度，0：无限大
	MaxIntervalOfSubmit time.Duration `yaml:"maxIntervalOfSubmit",json:"maxIntervalOfSubmit"` //提交任务时最大间隔
}

// Timewheel 时间轮
type Timewheel struct {
	slots        []*queue.LockFreeQueue[Task] //槽
	cancel       context.CancelFunc           //停止方法
	IsRunning    atomic.Bool                  // consurrent
	currentIndex int64                        //当前的slot索引
	currentMilli int64                        //毫秒
	pool         *ants.Pool                   //worker池
	logger       logger.Interface             //日志
	cfg          *Config                      //配置
}

func getOrDefaultConfig(cf *Config) *Config {
	c := &Config{}
	if cf != nil {
		*c = *cf
	}

	if c.Tick <= 0 {
		c.Tick = time.Second
	}

	if c.SlotCount <= 0 {
		c.SlotCount = 60
	}

	if c.MaxLengthOfEachSlot <= 0 {
		c.MaxLengthOfEachSlot = 0
	}

	if c.MaxIntervalOfSubmit <= 0 {
		c.MaxIntervalOfSubmit = c.Tick * time.Duration(c.SlotCount)
	}

	return c
}
func NewTimewheel(config *Config, logger logger.Interface, pool *ants.Pool) (*Timewheel, error) {

	c := getOrDefaultConfig(config)

	if c.MaxIntervalOfSubmit > c.Tick*time.Duration(c.SlotCount) {
		c.MaxIntervalOfSubmit = c.Tick * time.Duration(c.SlotCount)
		if logger != nil {
			logger.Warn("invalid MaxIntervalOfSubmit,reset as default value")
		}
	}

	if logger != nil && logger.IsDebugEnabled() {
		logger.Debug(string(jsonext.MarshalNoErr(c)))
	}

	tw := &Timewheel{
		slots:  make([]*queue.LockFreeQueue[Task], c.SlotCount),
		pool:   pool,
		logger: logger,
		cfg:    c,
	}

	for i := 0; i < c.SlotCount; i++ {
		tw.slots[i] = queue.NewLockFreeQueue[Task](int64(c.MaxLengthOfEachSlot))
	}

	if tw.pool == nil {
		p, err := ants.NewPool(
			runtime.NumCPU()*100,
			ants.WithNonblocking(true),
			ants.WithExpiryDuration(10*time.Second),
			ants.WithPreAlloc(true))
		if err != nil {
			if logger != nil {
				logger.Error("new ants pool err", zap.Error(err))
			}
			return nil, err
		}
		tw.pool = p
	}

	return tw, nil
}

func (tw *Timewheel) Start(ctx context.Context) {

	if tw.IsRunning.CompareAndSwap(false, true) {

		if tw.logger != nil {
			tw.logger.Info("ticker start")
		}

		ticker := time.NewTicker(tw.cfg.Tick)
		c, cancel := context.WithCancel(ctx)

		tw.cancel = cancel

		baseTime := time.Now()
		tickCount := int64(0)

		go func() {
			defer ticker.Stop()
			for {
				select {
				case <-c.Done():
					if tw.logger != nil && tw.logger.IsDebugEnabled() {
						tw.logger.Debug("ticker done", zap.Error(c.Err()))
					}
					return
				case <-ticker.C:
					tickCount++
					now := baseTime.Add(time.Duration(tickCount) * tw.cfg.Tick)
					if err := tw.advance(now); err != nil && tw.logger != nil {
						tw.logger.Error("advance failed", zap.Error(err))
					}
				}
			}
		}()
	}
}

func (tw *Timewheel) Stop() {

	if tw.IsRunning.CompareAndSwap(true, false) {

		if tw.logger != nil {
			tw.logger.Info("ticker stop")
		}

		if tw.cancel != nil {
			tw.cancel()
		}

		if tw.pool != nil {
			tw.pool.ReleaseTimeout(time.Second * 2)
		}

		if tw.logger != nil {
			tw.logger.Close()
		}
	}
}

func (tw *Timewheel) advance(now time.Time) error {
	// 计算槽索引：毫秒时间戳除以 tick（毫秒）
	tickMs := int64(tw.cfg.Tick / time.Millisecond)
	slot := (now.UnixMilli() / tickMs) % int64(tw.cfg.SlotCount)

	// 更新状态
	atomic.StoreInt64(&tw.currentMilli, now.UnixMilli())
	atomic.StoreInt64(&tw.currentIndex, slot)

	// 记录日志
	if tw.logger != nil && tw.logger.IsDebugEnabled() {
		tw.logger.Debug("ticking",
			zap.Int64("slot", slot),
			zap.Time("now", now),
			zap.Int64("slotsLen", tw.slots[slot].Len()))
	}

	if tw.pool != nil {
		err := tw.pool.Submit(func() {

			for {
				task, err := tw.slots[slot].Dequeue()
				if err != nil {
					break
				}
				if task == nil {
					continue
				}

				if tr := task.Execute(now); tr == Retry {
					if err := tw.slots[slot].Enqueue(task); err != nil && tw.logger != nil {
						tw.logger.Error("executed,but enqueue failed", zap.Error(err))
					}
				}

			}
		})
		if err != nil {
			tw.logger.Error("pool submit failed,call it directly", zap.Error(err))
			return err
		}
	}
	return nil
}

func (tw *Timewheel) Submit(task Task, delay time.Duration) (int, error) {
	if delay < 0 {
		err := fmt.Errorf("negative delay %v is not allowed", delay)
		tw.logger.Error("failed to submit", zap.Error(err))
		return -1, err
	}

	if delay > tw.cfg.MaxIntervalOfSubmit {
		err := fmt.Errorf("delay %v exceeds max interval %v", delay, tw.cfg.MaxIntervalOfSubmit)
		if tw.logger != nil {
			tw.logger.Error("failed to submit", zap.Error(err))
		}
		return -1, err
	}

	// 计算目标槽基于当前时间
	tickMs := int64(tw.cfg.Tick / time.Millisecond)
	currentMilli := atomic.LoadInt64(&tw.currentMilli)
	if currentMilli == 0 {
		currentMilli = time.Now().UnixMilli()
	}
	targetMilli := currentMilli + int64(delay/time.Millisecond)
	slot := int((targetMilli / tickMs) % int64(tw.cfg.SlotCount))

	err := tw.slots[slot].Enqueue(task)
	if tw.logger != nil && err != nil {
		tw.logger.Error("failed to submit", zap.Error(err))
	} else {
		if tw.logger != nil && tw.logger.IsDebugEnabled() {
			tw.logger.Debug("submit into slot",
				zap.Int("slot", slot),
				zap.Int64("len", tw.slots[slot].Len()))
		}
	}
	return slot, err
}
