package timewheel

import (
	"context"
	"fmt"
	"github.com/magicnana999/im/define"
	"github.com/magicnana999/im/pkg/logger"
	"github.com/magicnana999/im/pkg/queue"
	"github.com/panjf2000/ants/v2"
	"go.uber.org/zap"
	"runtime"
	"sync"
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
	Execute(now time.Time) (TaskResult, error)
}

type Config struct {
	Tick                time.Duration `yaml:"tick"`                //时间间隔
	SlotCount           int           `yaml:"slotCount"`           //槽数量
	MaxLengthOfEachSlot int           `yaml:"maxLengthOfEachSlot"` //每个槽的最大长度，0：无限大
	maxIntervalOfSubmit time.Duration //提交任务时最大间隔
}

// Timewheel 时间轮
type Timewheel struct {
	slots        []*queue.LockFreeQueue[Task] //槽
	cancel       context.CancelFunc           //停止方法
	once         sync.Once                    //用来保证start一次
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

	if c.Tick == 0 {
		c.Tick = time.Second
	}

	if c.SlotCount == 0 {
		c.SlotCount = 60
	}

	c.maxIntervalOfSubmit = c.Tick * time.Duration(c.SlotCount)

	return c
}
func NewTimeWheel(config *Config, logger logger.Interface, pool *ants.Pool) (*Timewheel, error) {

	c := getOrDefaultConfig(config)

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
		p, err := ants.NewPool(runtime.NumCPU(), ants.WithPreAlloc(true))
		if err != nil {
			return nil, err
		}
		tw.pool = p
	}

	return tw, nil
}

func (tw *Timewheel) Start(ctx context.Context) {
	tw.once.Do(func() {

		ticker := time.NewTicker(tw.cfg.Tick)
		c, cancel := context.WithCancel(ctx)

		tw.cancel = cancel

		go func() {
			defer ticker.Stop()
			for {
				select {
				case <-c.Done():
					return
				case t := <-ticker.C:
					tw.advance(t)
				}
			}
		}()

		if tw.logger != nil {
			tw.logger.Info("started", zap.String(define.OP, define.OpStart))
		}
	})
}
func (tw *Timewheel) Stop() {
	if tw.cancel != nil {
		tw.cancel()
	}

	if tw.pool != nil {
		tw.pool.Release()
	}

	if tw.logger != nil {
		tw.logger.Close()
	}

	if tw.logger != nil {
		tw.logger.Info("closed", zap.String(define.OP, define.OpClose))
	}

}

func (tw *Timewheel) advance(now time.Time) {
	// 计算槽索引：毫秒时间戳除以 tick（毫秒）
	tickMs := int64(tw.cfg.Tick / time.Millisecond)
	slot := (now.UnixMilli() / tickMs) % int64(tw.cfg.SlotCount)

	// 更新状态
	atomic.StoreInt64(&tw.currentMilli, now.UnixMilli())
	atomic.StoreInt64(&tw.currentIndex, slot)

	// 记录日志
	if tw.logger != nil && tw.logger.IsDebugEnabled() {
		msg := fmt.Sprintf("ticking slot[%d], now:%s, len:%d", slot, now.Format(time.StampMilli), tw.slots[slot].Len())
		tw.logger.Debug(msg,
			zap.String(define.OP, define.OpAdvance),
			zap.Duration("tick", tw.cfg.Tick))
	}

	// 处理任务
	for {
		task, err := tw.slots[slot].Dequeue()
		if err != nil {
			break
		}
		if task == nil {
			continue
		}

		tw.pool.Submit(func() {
			if tr, e := task.Execute(now); e == nil && tr == Retry {
				tw.slots[slot].Enqueue(task)
			}
		})
	}
}

// Submit adds a task to the Timewheel. The slot is calculated based on the current time,
// which may not be precise. Business logic should handle timing in Execute.
func (tw *Timewheel) Submit(task Task, delay time.Duration) (int, error) {

	if delay < 0 {
		delay = 0
	}

	// 检查最大延迟
	if delay > tw.cfg.maxIntervalOfSubmit {
		err := fmt.Errorf("delay %v exceeds max interval %v", delay, tw.cfg.maxIntervalOfSubmit)
		if tw.logger != nil {
			tw.logger.Error(err.Error(),
				zap.String(define.OP, define.OpSubmit),
			)
		}
		return -1, err
	}

	// 计算目标槽
	slotsToAdvance := int64(delay / tw.cfg.Tick)
	currentSlot := atomic.LoadInt64(&tw.currentIndex)
	slot := int(currentSlot+slotsToAdvance) % tw.cfg.SlotCount

	err := tw.slots[slot].Enqueue(task)
	if tw.logger != nil && err != nil {
		tw.logger.Error(err.Error(),
			zap.String(define.OP, define.OpSubmit),
		)
	} else {
		if tw.logger != nil && tw.logger.IsDebugEnabled() {
			tw.logger.Debug(fmt.Sprintf("submit into slots[%d],len:%d", slot, tw.slots[slot].Len()),
				zap.String(define.OP, define.OpSubmit),
			)
		}
	}
	return slot, err
}
