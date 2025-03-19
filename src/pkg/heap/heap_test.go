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

func TestMinHeapConcurrent(t *testing.T) {
	logger := log.New(os.Stderr, "", log.LstdFlags)
	less := func(i, j Task) bool { return i.expireTime < j.expireTime }
	onOverflow := func(t Task) bool {
		// 简化溢出处理，避免日志阻塞
		return true
	}

	h := NewMinHeap(less, Option[Task]{
		MaxSize:    100,
		OnOverflow: onOverflow,
		Logger:     logger,
	})

	const (
		numWorkers     = 50
		tasksPerWorker = 10
	)

	t.Run("ConcurrentPushPop", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var wg sync.WaitGroup
		// 推送任务
		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()
				for j := 0; j < tasksPerWorker; j++ {
					select {
					case <-ctx.Done():
						t.Logf("Worker%d timeout", workerID)
						return
					default:
						task := Task{
							clientID:   fmt.Sprintf("worker%d-task%d", workerID, j),
							expireTime: time.Now().UnixMilli() + int64((workerID+j)*100),
						}
						if err := h.PushTask(task); err != nil {
							t.Errorf("PushTask failed: %v", err)
						}
					}
				}
			}(i)
		}

		// 弹出任务
		for i := 0; i < numWorkers/2; i++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()
				for j := 0; j < tasksPerWorker/2; j++ {
					select {
					case <-ctx.Done():
						t.Logf("Worker%d timeout", workerID)
						return
					default:
						if _, err := h.PopTask(); err != nil && err.Error() != "heap is empty" {
							t.Errorf("PopTask failed: %v", err)
						}
					}
				}
			}(i)
		}

		// 等待或超时
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			// 正常完成
		case <-ctx.Done():
			t.Fatal("Test timed out after 5 seconds")
		}

		size := h.Size()
		if size < 0 {
			t.Errorf("Invalid heap size: %d", size)
		}
		t.Logf("Final heap size: %d", size)
	})
}
