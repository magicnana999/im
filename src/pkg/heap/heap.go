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
	MaxSize    int
	OnOverflow func(t T) bool
	Logger     *log.Logger
}

func NewMinHeap[T any](less func(i, j T) bool, opts ...Option[T]) *MinHeap[T] {
	opt := Option[T]{}
	if len(opts) > 0 {
		opt = opts[0]
	}
	if opt.MaxSize <= 0 {
		opt.MaxSize = 0
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
	heap.Init(h) // 无需锁，创建时单线程且 data 为空
	return h
}

// --- heap.Interface 方法，实现 container/heap 所需接口 ---

func (h *MinHeap[T]) Len() int           { return len(h.data) }
func (h *MinHeap[T]) Less(i, j int) bool { return h.less(h.data[i], h.data[j]) }
func (h *MinHeap[T]) Swap(i, j int)      { h.data[i], h.data[j] = h.data[j], h.data[i] }
func (h *MinHeap[T]) Push(x interface{}) { h.data = append(h.data, x.(T)) }
func (h *MinHeap[T]) Pop() interface{} {
	n := len(h.data)
	x := h.data[n-1]
	h.data = h.data[0 : n-1]
	return x
}

// --- 业务方法，提供工程化功能 ---

// Add 添加元素
func (h *MinHeap[T]) Add(item T) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.maxSize > 0 && len(h.data) >= h.maxSize {
		h.logger.Printf("WARN: Heap size %d exceeds maxSize %d", len(h.data), h.maxSize)
		if !h.onOverflow(item) {
			return errors.New("heap is full and overflow handler rejected item")
		}
		return nil
	}

	heap.Push(h, item)
	h.logger.Printf("DEBUG: Added item, heap size: %d", len(h.data))
	return nil
}

// Remove 移除并返回堆顶
func (h *MinHeap[T]) Remove() (T, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if len(h.data) == 0 {
		return *new(T), errors.New("heap is empty")
	}

	item := heap.Pop(h).(T)
	h.logger.Printf("DEBUG: Removed item, heap size: %d", len(h.data))
	return item, nil
}

// Top 查看堆顶（不移除）
func (h *MinHeap[T]) Top() (T, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.data) == 0 {
		return *new(T), errors.New("heap is empty")
	}
	return h.data[0], nil
}

// Count 返回当前大小
func (h *MinHeap[T]) Count() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.data)
}
