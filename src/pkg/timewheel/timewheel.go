package timewheel

import (
	"context"
	"github.com/magicnana999/im/pkg/queue"
	"github.com/panjf2000/ants/v2"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

const (
	maxLengthEachSlot = 1_000_000
)

type Task interface {
	Execute(now int64) error
}

type TimeWheel struct {
	slots  []*queue.LockFreeQueue[Task]
	tick   time.Duration
	cancel context.CancelFunc
	once   sync.Once
	index  int64
	second int64
	pool   *ants.Pool
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
					tw.advance(t.Unix())
				}
			}
		}()
	})
}
func (tw *TimeWheel) Stop() {
	if tw.cancel != nil {
		tw.cancel()
		tw.pool.Release()
	}
}

func (tw *TimeWheel) advance(now int64) {
	slot := now % int64(len(tw.slots))

	atomic.StoreInt64(&tw.second, now)
	atomic.StoreInt64(&tw.index, slot)

	for {
		task, err := tw.slots[slot].Dequeue()
		if err != nil {
			break
		}

		if task == nil {
			continue
		}

		tw.pool.Submit(func() {
			if e := task.Execute(now); e == nil {
				tw.slots[slot].Enqueue(task)
			}
		})
	}
}

// Submit adds a task to the TimeWheel. The slot is calculated based on the current time,
// which may not be precise. Business logic should handle timing in Execute.
func (tw *TimeWheel) Submit(task Task) (int64, error) {
	slot := (time.Now().Unix() - atomic.LoadInt64(&tw.second)) % int64(len(tw.slots))
	if slot < 0 {
		slot = 0
	}
	return slot, tw.slots[slot].Enqueue(task)
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

func (tw *TimeWheel) CurrentMaxQueueLength() int64 {
	var max int64
	for _, slot := range tw.slots {
		if slot.Len() > max {
			max = slot.Len()
		}
	}
	return max
}
