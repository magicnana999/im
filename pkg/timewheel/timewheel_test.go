package timewheel

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/magicnana999/im/global"
	"github.com/magicnana999/im/pkg/logger"
	"go.uber.org/atomic"
	"sync"
	"testing"
	"time"
)

type task func(now time.Time) TaskResult

func (t task) Execute(now time.Time) TaskResult {
	return t(now)
}

func TestTimewheelSecond(t *testing.T) {

	var wg sync.WaitGroup

	logger.Init(nil)
	defer logger.Close()

	c := &Config{
		Tick:      time.Second,
		SlotCount: 90,
	}

	config := global.RedisConfig{
		Addr:    "127.0.0.1:6379",
		DB:      0,
		Timeout: time.Second,
	}

	o := config.ToOptions()

	rds := redis.NewClient(o)

	f := func(now time.Time) TaskResult {
		rds.Incr(context.Background(), now.Format(time.DateTime))
		return Break
	}

	eachSlotTasks := 10000

	tw, _ := NewTimewheel(c, logger.Named("timewheel"), nil)

	var tasks atomic.Int64

	for i := 0; i < c.SlotCount; i++ {
		time.Sleep(time.Millisecond * 10)
		go func() {
			wg.Add(1)
			defer wg.Done()
			for j := 0; j < eachSlotTasks; j++ {
				time.Sleep(time.Millisecond * 10)
				tw.Submit(task(f), c.Tick*time.Duration(c.SlotCount))
				tasks.Inc()
			}
		}()
	}
	wg.Wait()

	fmt.Println("submit ok --------------------------------", tasks.Load())
	wg.Add(1)
	tw.Start(context.Background())
	defer tw.Stop()
	wg.Wait()
}

func TestTimewheelMilli(t *testing.T) {

	logger.Init(nil)
	defer logger.Close()

	c := &Config{
		Tick:      time.Millisecond * 100,
		SlotCount: 20,
	}

	//totalTask := c.SlotCount * 2

	tw, _ := NewTimewheel(c, logger.Named("timewheel"), nil)
	tw.Start(context.Background())
	defer tw.Stop()

	var counter int64
	//
	//for i := 0; i < totalTask; i++ {
	//	if _, err := tw.Submit(&MyTask{&counter}, time.Second*2); err != nil {
	//		t.Fatalf("Failed to submit task: %v", err)
	//	}
	//}

	time.Sleep(c.Tick*time.Duration(c.SlotCount) + 1)
	fmt.Println(counter)
}
