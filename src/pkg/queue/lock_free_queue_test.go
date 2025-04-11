package queue

import (
	"fmt"
	"github.com/magicnana999/im/pkg/heap"
	"sync"
	"testing"
	"time"
)

func TestLogic(t *testing.T) {
	q := NewLockFreeQueue[string](9999)

	// 入队 3 个元素
	values := []string{"A", "B", "C", "D", "E", "F"}
	for _, v := range values {
		q.Enqueue(v)
	}

	// 验证长度
	if q.Len() != int64(len(values)) {
		t.Errorf("expected length 3, got %d", q.Len())
	}

	// 出队并验证顺序
	for i, expected := range values {
		value, err := q.Dequeue()
		if err != nil {
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
}

func TestLoginConcurrent(t *testing.T) {
	q := NewLockFreeQueue[string](100_0000)
	var wg sync.WaitGroup
	// 多线程入队
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				q.Enqueue(fmt.Sprintf("%d%d", i, j))
			}
		}()
	}
	wg.Wait()

	if q.Len() != 100 {
		t.Errorf("expected length %d, got %d", 100, q.Len())
	}

	for i := 0; i < 100; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				fmt.Println(q.Dequeue())
				if v, err := q.Dequeue(); err != nil {
					return
				} else {
					fmt.Println(v)
				}
			}
		}()
	}

	time.Sleep(time.Second * 1)

}

func TestMaxLength(t *testing.T) {

	q := NewLockFreeQueue[int64](8000_0000)
	for {
		if err := q.Enqueue(time.Now().UnixNano()); err != nil {
			break
		}
	}

	fmt.Println("ok")
	fmt.Println(q.Len())
}

func TestMaxLengthConcurrent(t *testing.T) {

	var wg sync.WaitGroup

	q := NewLockFreeQueue[int64](8000_0000)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				if err := q.Enqueue(time.Now().UnixNano()); err != nil {
					return
				}
			}
		}()
	}

	wg.Wait()
	fmt.Println(q.Len())
}

// go test -bench=BenchmarkConcurrentEnqueue -benchtime=5s -v -run=^$ -benchmem -race
func BenchmarkConcurrentEnqueue(b *testing.B) {
	q := NewLockFreeQueue[int](1_0000_0000)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			q.Enqueue(100)
		}
	})
}

// go test -bench=BenchmarkConcurrentDequeue -benchtime=5s -v -run=^$ -benchmem -race
func BenchmarkConcurrentDequeue(b *testing.B) {
	q := NewLockFreeQueue[int](100_0000_0000)

	for i := 0; i < b.N; i++ {
		q.Enqueue(100)
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := q.Dequeue(); err != nil {
				return
			}
		}
	})
}

// go test -bench=BenchmarkConcurrentMixEnqueueDequeue -benchtime=5s -v -run=^$ -race

// go test -bench=BenchmarkConcurrentDequeue -benchmem -run=^$ -cpuprofile=cpu.prof
// go tool pprof cpu.prof
func BenchmarkConcurrentMixEnqueueDequeue(b *testing.B) {

	type Item struct {
		value int
	}

	f := func(i, j Item) bool {
		return i.value < j.value
	}

	const numItems = 10000000
	q := NewLockFreeQueue[Item](100_0000_0000)
	h := heap.NewMinHeap[Item](f, 0)
	for i := 0; i < numItems; i++ {
		h.Push(Item{i})
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
