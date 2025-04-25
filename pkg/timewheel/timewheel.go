package timewheel

import (
	"context"
	"errors"
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

const (
	DefSlotTick  = time.Second
	DefSlotCount = 60
	DefMaxLength = int64(1_000_000)
	DefExpire    = time.Second * 10
)

type Task interface {
	// Execute do not return any error,return Break instead
	Execute(now time.Time) TaskResult
}

type Config struct {
	SlotTick          time.Duration `yaml:"slotTick",json:"slotTick"`           //时间间隔
	SlotCount         int           `yaml:"slotCount",json:"slotCount"`         //槽数量
	SlotMaxLength     int64         `yaml:"slotMaxLength",json:"slotMaxLength"` //每个槽的最大长度，0：无限大
	WorkerCount       int           `yaml:"workerCount",json:"workerCount"`
	WorkerNonBlocking bool          `yaml:"workerNonBlocking",json:"workerNonBlocking"`
	WorkerExpire      time.Duration `yaml:"workerExpire",json:"workerExpire"`
	WorkerPreAlloc    bool          `yaml:"workerPreAlloc",json:"workerPreAlloc"`
	TaskInterval      time.Duration `yaml:"taskInterval",json:"taskInterval"`
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

	if c.SlotTick <= 0 {
		c.SlotTick = DefSlotTick
	}

	if c.SlotCount <= 0 {
		c.SlotCount = DefSlotCount
	}

	if c.TaskInterval > c.SlotTick*time.Duration(c.SlotCount) || c.TaskInterval <= 0 {
		c.TaskInterval = c.SlotTick * time.Duration(c.SlotCount)
	}

	if c.SlotMaxLength <= 0 {
		c.SlotMaxLength = DefMaxLength
	}

	if c.WorkerCount <= 0 {
		c.WorkerCount = runtime.NumCPU() * 100
	}

	if c.WorkerExpire <= 0 {
		c.WorkerExpire = DefExpire
	}

	return c
}
func NewTimewheel(config *Config, log logger.Interface, pool *ants.Pool) (*Timewheel, error) {

	c := getOrDefaultConfig(config)

	tw := &Timewheel{
		slots:  make([]*queue.LockFreeQueue[Task], c.SlotCount),
		pool:   pool,
		logger: log,
		cfg:    c,
	}

	for i := 0; i < c.SlotCount; i++ {
		tw.slots[i] = queue.NewLockFreeQueue[Task](c.SlotMaxLength)
	}

	if tw.logger == nil {
		logger.Init(nil)
		tw.logger = logger.Named("timewheel")
	}

	if tw.logger.IsDebugEnabled() {
		tw.logger.Debug(string(jsonext.MarshalNoErr(c)))
	}

	if tw.pool == nil {
		p, err := ants.NewPool(
			c.WorkerCount,
			ants.WithNonblocking(c.WorkerNonBlocking),
			ants.WithExpiryDuration(c.WorkerExpire),
			ants.WithPreAlloc(c.WorkerPreAlloc))
		if err != nil {
			tw.logger.Error("new ants pool err", zap.Error(err))
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

		ticker := time.NewTicker(tw.cfg.SlotTick)
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
					now := baseTime.Add(time.Duration(tickCount) * tw.cfg.SlotTick)
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
	tickMs := int64(tw.cfg.SlotTick / time.Millisecond)
	slot := (now.UnixMilli() / tickMs) % int64(tw.cfg.SlotCount)

	// 更新状态
	atomic.StoreInt64(&tw.currentMilli, now.UnixMilli())
	atomic.StoreInt64(&tw.currentIndex, slot)

	// 记录日志
	if tw.logger != nil && tw.logger.IsDebugEnabled() {
		tw.logger.Debug("ticking",
			zap.Int64("slot", slot),
			zap.Time("now", now),
			zap.Int("running", tw.pool.Running()),
			zap.Int("free", tw.pool.Free()),
			zap.Int64("slotsLen", tw.slots[slot].Len()))
	}

	tasks, err := tw.slots[slot].BatchDequeue(int(tw.cfg.SlotMaxLength))
	if err != nil && !errors.Is(err, queue.ErrQueueEmpty) {
		tw.logger.Error("batch dequeue failed", zap.Error(err))
	}

	batchSize := len(tasks) / (tw.cfg.WorkerCount / 20) // 更细粒度分批
	if batchSize < 100 {
		batchSize = 100
	} else if batchSize > 2000 {
		batchSize = 2000
	}

	if len(tasks) > 0 {

		for i := 0; i < len(tasks); i += batchSize {
			end := i + batchSize
			if end > len(tasks) {
				end = len(tasks)
			}
			batch := tasks[i:end]

			err := tw.pool.Submit(func() {
				retryTasks := make([]Task, 0, len(batch))
				for _, task := range batch {
					if tr := task.Execute(now); tr == Retry {
						retryTasks = append(retryTasks, task)
					}
				}
				if len(retryTasks) > 0 {
					if err := tw.slots[slot].BatchEnqueue(retryTasks); err != nil {
						tw.logger.Error("batch enqueue failed", zap.Error(err))
						for _, task := range retryTasks {
							tw.slots[slot].Enqueue(task) // 回退逐个入队
						}
					}
				}
			})

			if err != nil {
				tw.logger.Error("pool submit failed", zap.Error(err))
				// 直接执行，防止任务丢失
				for _, task := range batch {
					if tr := task.Execute(now); tr == Retry {
						tw.slots[slot].Enqueue(task)
					}
				}
			}
		}
	}

	return nil
}

func (tw *Timewheel) Put(task Task, index int) error {
	return tw.slots[index].Enqueue(task)
}

func (tw *Timewheel) Submit(task Task) (int, int64, error) {

	// 计算目标槽基于当前时间
	tickMs := int64(tw.cfg.SlotTick / time.Millisecond)
	currentMilli := atomic.LoadInt64(&tw.currentMilli)
	if currentMilli == 0 {
		currentMilli = time.Now().UnixMilli()
	}
	targetMilli := currentMilli + int64(tw.cfg.TaskInterval/time.Millisecond)
	slot := int((targetMilli / tickMs) % int64(tw.cfg.SlotCount))

	err := tw.slots[slot].Enqueue(task)
	if tw.logger != nil && err != nil {
		tw.logger.Error("failed to submit", zap.Error(err))
	} else {
		//if tw.logger != nil && tw.logger.IsDebugEnabled() {
		//	tw.logger.Debug("submit into slot",
		//		zap.Int("slot", slot),
		//		zap.Int64("len", tw.slots[slot].Len()))
		//}
	}
	return slot, tw.slots[slot].Len(), err
}
