package heap

import (
	"container/heap"
	"errors"
	"sync"
)

var (
	ErrEmptyHeap = errors.New("heap is empty")
	ErrFullHeap  = errors.New("heap is full")
)

type minHeapContainer[T any] struct {
	data []T
	less func(i, j T) bool
}

func (c *minHeapContainer[T]) Len() int {
	return len(c.data)
}

func (c *minHeapContainer[T]) Less(i, j int) bool {
	return c.less(c.data[i], c.data[j])
}

func (c *minHeapContainer[T]) Swap(i, j int) {
	c.data[i], c.data[j] = c.data[j], c.data[i]
}

func (c *minHeapContainer[T]) Push(x interface{}) {
	c.data = append(c.data, x.(T))
}

func (c *minHeapContainer[T]) Pop() interface{} {
	n := len(c.data)
	x := c.data[n-1]
	c.data = c.data[0 : n-1]
	return x
}

type MinHeap[T any] struct {
	c       minHeapContainer[T]
	mu      sync.RWMutex
	maxSize int
}

func NewMinHeap[T any](less func(i, j T) bool, maxSize int) *MinHeap[T] {

	h := &MinHeap[T]{
		c:       minHeapContainer[T]{data: make([]T, 0), less: less},
		maxSize: maxSize,
	}
	heap.Init(&h.c)
	return h
}

func (h *MinHeap[T]) Push(item T) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.maxSize > 0 && h.c.Len() >= h.maxSize {
		return ErrFullHeap
	}

	heap.Push(&h.c, item)
	return nil
}

func (h *MinHeap[T]) Pop() (T, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.c.Len() == 0 {
		var zero T
		return zero, ErrEmptyHeap
	}

	item := heap.Pop(&h.c).(T)
	return item, nil
}

func (h *MinHeap[T]) Top() (T, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if h.c.Len() == 0 {
		var zero T
		return zero, ErrEmptyHeap
	}
	return h.c.data[0], nil
}

func (h *MinHeap[T]) Len() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.c.Len()
}

func (h *MinHeap[T]) Iterate() []T {
	h.mu.RLock()
	dataCopy := make([]T, len(h.c.data))
	copy(dataCopy, h.c.data)
	h.mu.RUnlock()

	tempHeap := &minHeapContainer[T]{data: dataCopy, less: h.c.less}
	heap.Init(tempHeap)
	result := make([]T, 0, len(dataCopy))
	for tempHeap.Len() > 0 {
		item := heap.Pop(tempHeap).(T)
		result = append(result, item)
	}
	return result
}
