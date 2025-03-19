package heap

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

type Task struct {
	clientID   string
	expireTime int64
}

func TestHeap(t *testing.T) {
	logger := log.New(os.Stderr, "", log.LstdFlags)
	less := func(i, j Task) bool { return i.expireTime < j.expireTime }
	onOverflow := func(t Task) bool {
		fmt.Println(t)
		return true
	}

	heap := NewMinHeap(less, Option[Task]{
		MaxSize:    100,
		OnOverflow: onOverflow,
		Logger:     logger,
	})

	for i := 0; i < 100; i++ {
		now := time.Now().UnixNano()
		fmt.Sprintf("%d", now)
		heap.Add(Task{clientID: fmt.Sprintf("%d", now), expireTime: now})
	}

	fmt.Printf("heap size: %d\n", heap.Count())

}

var counter int

// go test -v -race -run TestRaceCondition
func TestRaceCondition(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter++ // 未加锁的并发写
		}()
	}
	wg.Wait()
}

// go test -v -race -run TestMinHeapConcurrent
func TestMinHeapConcurrent(t *testing.T) {
	logger := log.New(os.Stderr, "", log.LstdFlags)
	less := func(i, j Task) bool { return i.expireTime < j.expireTime }
	onOverflow := func(t Task) bool {
		logger.Printf("INFO: Overflow: moving %s:%d to another slot", t.clientID, t.expireTime)
		return true
	}

	h := NewMinHeap(less, Option[Task]{
		MaxSize:    50,
		OnOverflow: onOverflow,
		Logger:     logger,
	})

	const (
		numWorkers     = 20
		tasksPerWorker = 5
		popWorkers     = 10
		popsPerWorker  = 3
	)

	t.Run("ConcurrentOperations", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var wg sync.WaitGroup
		var lastRemoved int64 = -1 // 记录上一次移除的 expireTime
		var mu sync.Mutex          // 保护 lastRemoved

		// 并发推送任务
		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()
				for j := 0; j < tasksPerWorker; j++ {
					select {
					case <-ctx.Done():
						t.Logf("Push worker%d timeout, tasks completed: %d/%d", workerID, j, tasksPerWorker)
						return
					default:
						task := Task{
							clientID:   fmt.Sprintf("worker%d-task%d", workerID, j),
							expireTime: time.Now().UnixMilli() + int64((workerID+j)*100),
						}
						if err := h.Add(task); err != nil {
							t.Errorf("Worker%d Add failed: %v", workerID, err)
						}
					}
				}
			}(i)
		}

		// 并发弹出任务
		for i := 0; i < popWorkers; i++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()
				for j := 0; j < popsPerWorker; j++ {
					select {
					case <-ctx.Done():
						t.Logf("Pop worker%d timeout, tasks completed: %d/%d", workerID, j, popsPerWorker)
						return
					default:
						if task, err := h.Remove(); err != nil {
							if err.Error() != "heap is empty" {
								t.Errorf("Worker%d Remove failed: %v", workerID, err)
							}
						} else {
							mu.Lock()
							if task.expireTime < lastRemoved {
								t.Errorf("Worker%d popped invalid task: %s:%d < last %d", workerID, task.clientID, task.expireTime, lastRemoved)
							}
							lastRemoved = task.expireTime
							mu.Unlock()
						}
					}
				}
			}(i)
		}

		// 并发查看堆顶和计数
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()
				for j := 0; j < 2; j++ {
					select {
					case <-ctx.Done():
						t.Logf("Read worker%d timeout", workerID)
						return
					default:
						if top, err := h.Top(); err == nil {
							t.Logf("Worker%d Top: %s:%d", workerID, top.clientID, top.expireTime)
						}
						count := h.Count()
						if count < 0 {
							t.Errorf("Worker%d Count invalid: %d", workerID, count)
						}
					}
				}
			}(i)
		}

		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			t.Log("All goroutines completed successfully")
		case <-ctx.Done():
			t.Fatal("Test timed out after 5 seconds")
		}

		finalSize := h.Count()
		if finalSize < 0 {
			t.Errorf("Final heap size invalid: %d", finalSize)
		}
		t.Logf("Final heap size: %d", finalSize)

		if top, err := h.Top(); err == nil {
			t.Logf("Final top: %s:%d", top.clientID, top.expireTime)
		} else if err.Error() != "heap is empty" {
			t.Errorf("Final Top failed: %v", err)
		}
	})
}
