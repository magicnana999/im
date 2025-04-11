package heap

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

type Item struct {
	Value int
}

// lessFunc 定义比较函数，小顶堆按 Value 升序
func lessFunc(i, j Item) bool {
	return i.Value < j.Value
}

// go test -v -race -run TestHeap
func TestHeapLogic(t *testing.T) {

	f := func(i, j *Item) bool {
		if i == nil || j == nil {
			return true
		}
		return i.Value < j.Value
	}

	heap := NewMinHeap[*Item](f, 10)

	heap.Push(&Item{500})
	heap.Push(&Item{200})
	heap.Push(&Item{300})
	heap.Push(&Item{600})
	heap.Push(&Item{100})
	heap.Push(&Item{400})
	for i, v := range heap.Iterate() {
		fmt.Println(i, v)
	}
}

// go test -v -race -run TestHeapMaxSize -benchmem
func TestHeapMaxSize(t *testing.T) {
	f := func(i, j *Item) bool {
		if i == nil || j == nil {
			return true
		}
		return i.Value < j.Value
	}

	size := 100_0000_00
	heap := NewMinHeap[*Item](f, size)

	for i := 0; i < size; i++ {
		heap.Push(&Item{time.Now().Nanosecond()})
	}
	fmt.Println(heap.Len())
}

// BenchmarkMinHeapAdd 测试 Push 操作性能（单线程）
func BenchmarkMinHeapAdd(b *testing.B) {
	h := NewMinHeap(lessFunc, 10000) // 容量 10000

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Push(Item{Value: rand.Intn(10000)})
	}
}

// BenchmarkMinHeapRemove 测试 Pop 操作性能（单线程）
func BenchmarkMinHeapRemove(b *testing.B) {
	h := NewMinHeap(lessFunc, 10000)
	// 预填充堆
	for i := 0; i < 10000; i++ {
		h.Push(Item{Value: rand.Intn(10000)})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if h.Len() == 0 {
			// 如果堆空，重置堆
			for j := 0; j < 10000; j++ {
				h.Push(Item{Value: rand.Intn(10000)})
			}
		}
		h.Pop()
	}
}

// BenchmarkMinHeapTop 测试 Top 操作性能（单线程）
func BenchmarkMinHeapTop(b *testing.B) {
	h := NewMinHeap(lessFunc, 10000)
	// 预填充堆
	for i := 0; i < 10000; i++ {
		h.Push(Item{Value: rand.Intn(10000)})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Top()
	}
}

// BenchmarkMinHeapConcurrentAdd 测试 Push 操作性能（并发）
// go test -bench=BenchmarkMinHeapConcurrentAdd -benchmem
// go test -bench=BenchmarkMinHeapConcurrentAdd -benchtime=5s -v -benchmem
func BenchmarkMinHeapConcurrentAdd(b *testing.B) {
	h := NewMinHeap(lessFunc, 10000)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			h.Push(Item{Value: rand.Intn(10000)})
		}
	})
}

// BenchmarkMinHeapConcurrentRemove 测试 Pop 操作性能（并发）
func BenchmarkMinHeapConcurrentRemove(b *testing.B) {
	h := NewMinHeap(lessFunc, 10000)
	// 预填充堆
	for i := 0; i < 10000; i++ {
		h.Push(Item{Value: rand.Intn(10000)})
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if h.Len() == 0 {
				// 避免空堆，单线程重置
				if h.Len() == 0 {
					for j := 0; j < 10000; j++ {
						h.Push(Item{Value: rand.Intn(10000)})
					}
				}
			}
			h.Pop()
		}
	})
}

// BenchmarkMinHeapConcurrentTop 测试 Top 操作性能（并发）
// go test -bench=BenchmarkMinHeapConcurrentTop -benchtime=5s -v -benchmem
func BenchmarkMinHeapConcurrentTop(b *testing.B) {
	h := NewMinHeap(lessFunc, 10000)
	// 预填充堆
	for i := 0; i < 10000; i++ {
		h.Push(Item{Value: rand.Intn(10000)})
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			h.Top()
		}
	})
}
