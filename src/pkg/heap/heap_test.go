package heap

import (
	"fmt"
	"math/rand"
	"runtime"
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
func TestHeap(t *testing.T) {

	f := func(i, j *Item) bool {
		if i == nil || j == nil {
			return true
		}
		return i.Value < j.Value
	}

	heap := NewMinHeap[*Item](f, 10)

	for i := 0; i < runtime.NumCPU(); i++ {
		heap.Add(&Item{time.Now().Nanosecond()})
	}
	fmt.Println("haha")
	for i := range heap.Iterate() {
		fmt.Println(i)
	}

}

// BenchmarkMinHeapAdd 测试 Add 操作性能（单线程）
func BenchmarkMinHeapAdd(b *testing.B) {
	h := NewMinHeap(lessFunc, 10000) // 容量 10000

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Add(Item{Value: rand.Intn(10000)})
	}
}

// BenchmarkMinHeapRemove 测试 Remove 操作性能（单线程）
func BenchmarkMinHeapRemove(b *testing.B) {
	h := NewMinHeap(lessFunc, 10000)
	// 预填充堆
	for i := 0; i < 10000; i++ {
		h.Add(Item{Value: rand.Intn(10000)})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if h.Count() == 0 {
			// 如果堆空，重置堆
			for j := 0; j < 10000; j++ {
				h.Add(Item{Value: rand.Intn(10000)})
			}
		}
		h.Remove()
	}
}

// BenchmarkMinHeapTop 测试 Top 操作性能（单线程）
func BenchmarkMinHeapTop(b *testing.B) {
	h := NewMinHeap(lessFunc, 10000)
	// 预填充堆
	for i := 0; i < 10000; i++ {
		h.Add(Item{Value: rand.Intn(10000)})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Top()
	}
}

// BenchmarkMinHeapConcurrentAdd 测试 Add 操作性能（并发）
func BenchmarkMinHeapConcurrentAdd(b *testing.B) {
	h := NewMinHeap(lessFunc, 10000)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			h.Add(Item{Value: rand.Intn(10000)})
		}
	})
}

// BenchmarkMinHeapConcurrentRemove 测试 Remove 操作性能（并发）
func BenchmarkMinHeapConcurrentRemove(b *testing.B) {
	h := NewMinHeap(lessFunc, 10000)
	// 预填充堆
	for i := 0; i < 10000; i++ {
		h.Add(Item{Value: rand.Intn(10000)})
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if h.Count() == 0 {
				// 避免空堆，单线程重置
				if h.Count() == 0 {
					for j := 0; j < 10000; j++ {
						h.Add(Item{Value: rand.Intn(10000)})
					}
				}
			}
			h.Remove()
		}
	})
}

// BenchmarkMinHeapConcurrentTop 测试 Top 操作性能（并发）
func BenchmarkMinHeapConcurrentTop(b *testing.B) {
	h := NewMinHeap(lessFunc, 10000)
	// 预填充堆
	for i := 0; i < 10000; i++ {
		h.Add(Item{Value: rand.Intn(10000)})
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			h.Top()
		}
	})
}
