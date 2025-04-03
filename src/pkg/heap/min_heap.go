package heap

import (
	"container/heap"
	"errors"
	"sync"
)

var (
	EmptyHeap = errors.New("heap is empty")
	FullHeap  = errors.New("heap is full")
)

// MinHeap 通用的工程化小顶堆
type MinHeap[T any] struct {
	data    []T
	less    func(i, j T) bool // 比较函数
	mu      sync.RWMutex      // 线程安全锁
	maxSize int               // 最大容量
}

func NewMinHeap[T any](less func(i, j T) bool, maxSize int) *MinHeap[T] {

	h := &MinHeap[T]{
		data:    make([]T, 0),
		less:    less,
		maxSize: maxSize,
	}
	heap.Init(h)
	return h
}

// implement heap.Interface begin...

func (h *MinHeap[T]) Len() int {
	return len(h.data)
}

func (h *MinHeap[T]) Less(i, j int) bool {
	return h.less(h.data[i], h.data[j])
}

func (h *MinHeap[T]) Swap(i, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
}

func (h *MinHeap[T]) Push(x interface{}) {
	h.data = append(h.data, x.(T))
}

func (h *MinHeap[T]) Pop() interface{} {
	n := len(h.data)
	x := h.data[n-1]
	h.data = h.data[0 : n-1]
	return x
}

// implement heap.Interface end...

// Add 添加元素
func (h *MinHeap[T]) Add(item T) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.maxSize > 0 && len(h.data) >= h.maxSize {
		return FullHeap
	}

	heap.Push(h, item)
	return nil
}

// Remove 移除并返回堆顶
func (h *MinHeap[T]) Remove() (T, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if len(h.data) == 0 {
		var zero T
		return zero, EmptyHeap
	}

	item := heap.Pop(h).(T)
	return item, nil
}

// Top 查看堆顶（不移除）
func (h *MinHeap[T]) Top() (T, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.data) == 0 {
		var zero T
		return zero, EmptyHeap
	}
	return h.data[0], nil
}

// Count 返回当前大小
func (h *MinHeap[T]) Count() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.data)
}

func (h *MinHeap[T]) Iterate() []T {
	h.mu.RLock()
	dataCopy := make([]T, len(h.data))
	copy(dataCopy, h.data)
	h.mu.RUnlock()

	// 创建临时堆
	tempHeap := &MinHeap[T]{
		data:    dataCopy,
		less:    h.less,
		maxSize: 0, // 无需容量限制
	}
	heap.Init(tempHeap)

	// 反复移除堆顶，构建有序切片
	result := make([]T, 0, len(dataCopy))
	for tempHeap.Len() > 0 {
		item := heap.Pop(tempHeap).(T)
		result = append(result, item)
	}
	return result
}
