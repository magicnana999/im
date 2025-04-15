package timewheel

import (
	"context"
	"fmt"
	"github.com/magicnana999/im/pkg/logger"
	"github.com/magicnana999/im/pkg/queue"
	"github.com/panjf2000/ants/v2"
	"go.uber.org/zap"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

const (
	maxLengthEachSlot = 0

	SCOPE        = "timewheel"
	OperaStart   = "start"
	OperaSubmit  = "submit"
	OperaStop    = "stop"
	OperaAdvance = "advance"
)

type TaskResult string

const (
	Retry TaskResult = "retry"
	Break TaskResult = "break"
)

type Task interface {
	Execute(now int64) (TaskResult, error)
}

type TimeWheel struct {
	slots  []*queue.LockFreeQueue[Task]
	tick   time.Duration
	cancel context.CancelFunc
	once   sync.Once
	index  int64
	second int64
	pool   *ants.Pool
	logger logger.Interface
}

func NewTimeWheel(tick time.Duration, size int, pool *ants.Pool) (*TimeWheel, error) {
	tw := &TimeWheel{
		slots: make([]*queue.LockFreeQueue[Task], size),
		tick:  tick,
		pool:  pool,
	}

	for i := 0; i < size; i++ {
		tw.slots[i] = queue.NewLockFreeQueue[Task](maxLengthEachSlot)
	}

	if tw.pool == nil {
		p, err := ants.NewPool(runtime.NumCPU(), ants.WithPreAlloc(true))
		if err != nil {
			return nil, err
		}
		tw.pool = p
	}

	if tw.logger == nil {
		tw.logger = logger.Init(nil)
	}
	return tw, nil
}

func (tw *TimeWheel) Start(ctx context.Context) {
	tw.once.Do(func() {

		ticker := time.NewTicker(tw.tick)
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

		if tw.logger != nil && tw.logger.IsDebugEnabled() {
			tw.logger.Debug("start",
				zap.String(logger.SCOPE, SCOPE),
				zap.String(logger.OP, OperaStart))
		}
	})
}
func (tw *TimeWheel) Stop() {
	if tw.cancel != nil {
		tw.cancel()
		tw.pool.Release()
	}

	if tw.pool != nil {
		tw.pool.Release()
	}

	if tw.logger != nil {
		tw.logger.Sync()
	}

	if tw.logger != nil && tw.logger.IsDebugEnabled() {
		tw.logger.Debug("stop",
			zap.String(logger.SCOPE, SCOPE),
			zap.String(logger.OP, OperaStop))
	}

}

func (tw *TimeWheel) advance(now time.Time) {
	slot := now.Unix() % int64(len(tw.slots))

	atomic.StoreInt64(&tw.second, now.Unix())
	atomic.StoreInt64(&tw.index, slot)

	if tw.logger != nil && tw.logger.IsDebugEnabled() {
		msg := fmt.Sprintf("ticking slot[%d],now:%s,len:%d", slot, now.Format(time.DateTime), tw.slots[slot].Len())
		tw.logger.Debug(msg, zap.String(logger.SCOPE, SCOPE), zap.String(logger.OP, OperaAdvance))
	}

	for {
		task, err := tw.slots[slot].Dequeue()
		if err != nil {
			break
		}

		if task == nil {
			continue
		}

		tw.pool.Submit(func() {
			if tr, e := task.Execute(now.Unix()); e == nil && tr == Retry {
				tw.slots[slot].Enqueue(task)
			}
		})
	}
}

// Submit adds a task to the TimeWheel. The slot is calculated based on the current time,
// which may not be precise. Business logic should handle timing in Execute.
func (tw *TimeWheel) Submit(task Task, delay int64) (int64, error) {

	if delay < 0 {
		delay = 0
	}

	slotCount := int64(len(tw.slots))
	currentSlot := atomic.LoadInt64(&tw.index)
	slot := (currentSlot + delay) % slotCount
	if slot < 0 {
		slot = 0
	}

	err := tw.slots[slot].Enqueue(task)
	if tw.logger != nil && err != nil {
		tw.logger.Error(err.Error(),
			zap.String(logger.SCOPE, SCOPE),
			zap.String(logger.OP, OperaSubmit),
		)
	} else {
		if tw.logger != nil && tw.logger.IsDebugEnabled() {
			tw.logger.Debug(fmt.Sprintf("submit into slots[%d],len:%d", slot, tw.slots[slot].Len()),
				zap.String(logger.SCOPE, SCOPE),
				zap.String(logger.OP, OperaSubmit),
			)
		}
	}
	return slot, err
}

func (tw *TimeWheel) CurrentSlotIndex() int64 {
	return atomic.LoadInt64(&tw.index)
}

func (tw *TimeWheel) CurrentTimeSecond() int64 {
	return atomic.LoadInt64(&tw.second)
}

func (tw *TimeWheel) CurrentSlotsLen() int64 {
	return int64(len(tw.slots))
}

func (tw *TimeWheel) CurrentQueueLen() []int {
	ret := make([]int, len(tw.slots))
	for i, slot := range tw.slots {
		ret[i] = int(slot.Len())
	}
	return ret
}
