package timewheel

import (
	"context"
	"fmt"
	"github.com/magicnana999/im/pkg/logger"
	"sync/atomic"
	"testing"
	"time"
)

type MyTask struct {
	Counter *int64
}

func (t *MyTask) Execute(now time.Time) (TaskResult, error) {
	atomic.AddInt64(t.Counter, 1)
	return Break, nil
}

func TestTimewheelSecond(t *testing.T) {

	logger.Init(nil)
	defer logger.Close()

	c := &Config{
		Tick:      time.Second,
		SlotCount: 2,
	}

	totalTask := c.SlotCount * 10

	tw, _ := NewTimeWheel(c, logger.Named("timewheel"), nil)
	tw.Start(context.Background())
	defer tw.Stop()

	var counter int64

	for i := 0; i < totalTask; i++ {
		if _, err := tw.Submit(&MyTask{&counter}, time.Second*2); err != nil {
			t.Fatalf("Failed to submit task: %v", err)
		}
	}

	time.Sleep(c.Tick*time.Duration(c.SlotCount) + 1)
	fmt.Println(counter)
}

func TestTimewheelMilli(t *testing.T) {

	logger.Init(nil)
	defer logger.Close()

	c := &Config{
		Tick:      time.Millisecond * 100,
		SlotCount: 20,
	}

	totalTask := c.SlotCount * 2

	tw, _ := NewTimeWheel(c, logger.Named("timewheel"), nil)
	tw.Start(context.Background())
	defer tw.Stop()

	var counter int64

	for i := 0; i < totalTask; i++ {
		if _, err := tw.Submit(&MyTask{&counter}, time.Second*2); err != nil {
			t.Fatalf("Failed to submit task: %v", err)
		}
	}

	time.Sleep(c.Tick*time.Duration(c.SlotCount) + 1)
	fmt.Println(counter)
}
