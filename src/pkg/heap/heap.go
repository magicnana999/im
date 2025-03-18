package heap

import (
	"container/heap"
	"errors"
	"log"
	"sync"
)

// MinHeap 通用的工程化小顶堆
type MinHeap[T any] struct {
	data       []T
	less       func(i, j T) bool // 比较函数
	mu         sync.RWMutex      // 线程安全锁
	maxSize    int               // 最大容量
	onOverflow func(t T) bool    // 溢出处理回调
	logger     *log.Logger       // 标准日志
}

// Option 配置选项
type Option[T any] struct {
	MaxSize    int            // 最大容量，默认无限制
	OnOverflow func(t T) bool // 溢出处理回调，默认拒绝
	Logger     *log.Logger    // 日志，默认标准输出
}

func NewMinHeap[T any](less func(i, j T) bool, opts ...Option[T]) *MinHeap[T] {
	opt := Option[T]{}
	if len(opts) > 0 {
		opt = opts[0]
	}
	if opt.MaxSize <= 0 {
		opt.MaxSize = 0 // 0 表示无限制
	}
	if opt.OnOverflow == nil {
		opt.OnOverflow = func(t T) bool { return false }
	}
	if opt.Logger == nil {
		opt.Logger = log.Default()
	}

	h := &MinHeap[T]{
		data:       make([]T, 0),
		less:       less,
		maxSize:    opt.MaxSize,
		onOverflow: opt.OnOverflow,
		logger:     opt.Logger,
	}
	heap.Init(h)
	return h
}

func (h *MinHeap[T]) Len() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.data)
}

func (h *MinHeap[T]) Less(i, j int) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.less(h.data[i], h.data[j])
}

func (h *MinHeap[T]) Swap(i, j int) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.data[i], h.data[j] = h.data[j], h.data[i]
}

func (h *MinHeap[T]) Push(x interface{}) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.data = append(h.data, x.(T))
}

func (h *MinHeap[T]) Pop() interface{} {
	h.mu.Lock()
	defer h.mu.Unlock()
	n := len(h.data)
	x := h.data[n-1]
	h.data = h.data[0 : n-1]
	return x
}

// PushTask 添加元素
func (h *MinHeap[T]) PushTask(item T) error {
	// 先检查容量，不锁
	if h.maxSize > 0 {
		h.mu.RLock()
		size := len(h.data)
		h.mu.RUnlock()
		if size >= h.maxSize {
			h.logger.Printf("WARN: Heap size %d exceeds maxSize %d", size, h.maxSize)
			if !h.onOverflow(item) {
				return errors.New("heap is full and overflow handler rejected item")
			}
			return nil
		}
	}

	// 只在修改时加锁
	h.mu.Lock()
	heap.Push(h, item)
	size := len(h.data)
	h.mu.Unlock()

	h.logger.Printf("DEBUG: Pushed item, heap size: %d", size)
	return nil
}

// PopTask 移除并返回堆顶
func (h *MinHeap[T]) PopTask() (T, error) {
	h.mu.Lock()
	if len(h.data) == 0 {
		h.mu.Unlock()
		return *new(T), errors.New("heap is empty")
	}
	item := heap.Pop(h).(T)
	size := len(h.data)
	h.mu.Unlock()

	h.logger.Printf("DEBUG: Popped item, heap size: %d", size)
	return item, nil
}

// Peek 查看堆顶（不移除）
func (h *MinHeap[T]) Peek() (T, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.data) == 0 {
		return *new(T), errors.New("heap is empty")
	}
	return h.data[0], nil
}

// Size 返回当前大小
func (h *MinHeap[T]) Size() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.data)
}
