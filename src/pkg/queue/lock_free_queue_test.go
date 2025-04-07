package queue

import (
	"fmt"
	"github.com/magicnana999/im/pkg/heap"
	"sync"
	"testing"
	"time"
)

func TestEnqueueDequeueLogic(t *testing.T) {
	q := NewLockFreeQueue[string]()

	// 入队 3 个元素
	values := []string{"A", "B", "C"}
	for _, v := range values {
		q.Enqueue(v)
	}

	// 验证长度
	if q.Len() != 3 {
		t.Errorf("expected length 3, got %d", q.Len())
	}

	// 出队并验证顺序
	for i, expected := range values {
		value, ok := q.Dequeue()
		if !ok {
			t.Errorf("expected ok=true at index %d, got false", i)
		}
		if value != expected {
			t.Errorf("expected value %s at index %d, got %v", expected, i, value)
		}
	}

	// 验证队列为空
	if q.Len() != 0 {
		t.Errorf("expected length 0 after dequeue, got %d", q.Len())
	}
	value, ok := q.Dequeue()
	if ok || value != "" {
		t.Errorf("expected nil, false for empty queue, got %v, %v", value, ok)
	}
}

func TestConcurrentCorrectness(t *testing.T) {
	q := NewLockFreeQueue[int]()
	var wg sync.WaitGroup
	// 多线程入队
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				q.Enqueue(time.Now().Nanosecond())
			}
		}()
	}
	wg.Wait()

	if q.Len() != 100 {
		t.Errorf("expected length %d, got %d", 100, q.Len())
	}

	for i := 0; i < 5; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				fmt.Println(q.Dequeue())
			}
		}()
	}

	time.Sleep(time.Second * 10)

}

// go test -bench=BenchmarkConcurrentEnqueue -benchtime=5s -v -run=^$ -race
func BenchmarkConcurrentEnqueue(b *testing.B) {
	q := NewLockFreeQueue[int]()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			q.Enqueue(100)
		}
	})
}

// go test -bench=BenchmarkSingleThreadDequeue -benchtime=5s -v -race -run=^$
func BenchmarkSingleThreadDequeue(b *testing.B) {
	q := NewLockFreeQueue[int]()
	for i := 0; i < b.N; i++ {
		q.Enqueue(100)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Dequeue()
	}
}

// go test -bench=BenchmarkConcurrentDequeue -benchtime=5s -v -run=^$ -race
func BenchmarkConcurrentDequeue(b *testing.B) {
	q := NewLockFreeQueue[int]()

	for i := 0; i < b.N; i++ {
		q.Enqueue(100)
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, ok := q.Dequeue(); !ok {
				return
			}
		}
	})
}

// go test -bench=BenchmarkConcurrentMixEnqueueDequeue -benchtime=5s -v -run=^$ -race
func BenchmarkConcurrentMixEnqueueDequeue(b *testing.B) {

	type Item struct {
		value int
	}

	f := func(i, j Item) bool {
		return i.value < j.value
	}

	const numItems = 10000000
	q := NewLockFreeQueue[Item]()
	h := heap.NewMinHeap[Item](f, 0)
	for i := 0; i < numItems; i++ {
		h.Add(Item{i})
	}

	for _, item := range h.Iterate() {
		q.Enqueue(item)
	}

	b.ResetTimer()

	var wg sync.WaitGroup
	b.RunParallel(func(pb *testing.PB) {
		wg.Add(1)
		// 每个 goroutine 随机入队或出队
		for pb.Next() {
			q.Dequeue()
		}
		wg.Done()
	})

	wg.Wait()
}
