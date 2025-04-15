package timewheel

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

type MyTask struct {
	Counter *int64
}

func (t *MyTask) Execute(now int64) error {

	if time.Now().Unix()-now > 1 {
		fmt.Println("error1")
	}
	atomic.AddInt64(t.Counter, 1)

	if time.Now().Unix()-now > 1 {
		fmt.Println("error2")
	}
	return errors.New("task is completed")
}

// TestTimeWheelFunctionality verifies the functional correctness of TimeWheel.
func TestTimeWheelFunctionality(t *testing.T) {

	slot := 3
	totalTask := 21
	tick := time.Second

	tw, _ := NewTimeWheel(tick, slot, nil)
	tw.Start(context.Background())
	defer tw.Stop()

	var counter int64

	for i := 0; i < totalTask; i++ {
		if _, err := tw.Submit(&MyTask{&counter}, 30); err != nil {
			t.Fatalf("Failed to submit task: %v", err)
		}
		time.Sleep(time.Second)
	}

	time.Sleep(time.Second * 2)
	fmt.Println(counter)
}

//// TestTimeWheelRace verifies there are no race conditions in TimeWheel.
//func TestTimeWheelRace(t *testing.T) {
//	tw, err := NewTimeWheel(100*time.Millisecond, 10, nil)
//	if err != nil {
//		t.Fatalf("Failed to create TimeWheel: %v", err)
//	}
//	tw.Start(context.Background())
//	defer tw.Stop()
//
//	var wg sync.WaitGroup
//	var counter int64
//	const numTasks = 1000
//
//	// Concurrently submit tasks
//	for i := 0; i < numTasks; i++ {
//		wg.Add(1)
//		go func(id int) {
//			defer wg.Done()
//			task := &MyTask{
//				ID:        int64(id),
//				ExecuteAt: time.Now().Unix(),
//				Counter:   &counter,
//			}
//			if err := tw.Submit(task); err != nil {
//				t.Errorf("Failed to submit task %d: %v", id, err)
//			}
//		}(i)
//	}
//
//	wg.Wait()
//	time.Sleep(1 * time.Second) // Allow tasks to execute
//
//	if atomic.LoadInt64(&counter) < numTasks {
//		t.Errorf("Expected at least %d executions, got %d", numTasks, counter)
//	}
//}
//
//func BenchmarkTimeWheelSubmit(b *testing.B) {
//	tw, err := NewTimeWheel(100*time.Millisecond, 60, nil)
//	if err != nil {
//		b.Fatalf("Failed to create TimeWheel: %v", err)
//	}
//	tw.Start(context.Background())
//	defer tw.Stop()
//
//	var counter int64
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		task := &MyTask{
//			ID:        int64(i),
//			ExecuteAt: time.Now().Unix(),
//			Counter:   &counter,
//		}
//		if err := tw.Submit(task); err != nil {
//			b.Fatalf("Failed to submit task %d: %v", i, err)
//		}
//	}
//}
//
//// BenchmarkTimeWheelExecute measures the performance of executing tasks in TimeWheel.
//func BenchmarkTimeWheelExecute(b *testing.B) {
//	tw, err := NewTimeWheel(100*time.Millisecond, 60, nil)
//	if err != nil {
//		b.Fatalf("Failed to create TimeWheel: %v", err)
//	}
//	tw.Start(context.Background())
//	defer tw.Stop()
//
//	var counter int64
//	const numTasks = 100_000 // Preload tasks
//	for i := 0; i < numTasks; i++ {
//		task := &MyTask{
//			ID:        int64(i),
//			ExecuteAt: time.Now().Unix(),
//			Counter:   &counter,
//		}
//		tw.Submit(task)
//	}
//
//	// Reset counter and timer before execution benchmark
//	atomic.StoreInt64(&counter, 0)
//	b.ResetTimer()
//
//	// Run for b.N iterations, each iteration processes the preloaded tasks
//	for i := 0; i < b.N; i++ {
//		time.Sleep(100 * time.Millisecond) // Simulate one tick
//	}
//
//	b.StopTimer()
//	executed := atomic.LoadInt64(&counter)
//	b.Logf("Executed %d tasks in %v iterations", executed, b.N)
//}
