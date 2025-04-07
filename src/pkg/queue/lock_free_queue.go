package queue

import (
	"errors"
	"sync/atomic"
	"time"
	"unsafe"
)

var ErrQueueFull = errors.New("queue is full")

// LockFreeNode 泛型节点
type LockFreeNode[T any] struct {
	value T
	next  unsafe.Pointer
}

// LockFreeQueue 泛型无锁队列
type LockFreeQueue[T any] struct {
	head    unsafe.Pointer
	tail    unsafe.Pointer
	length  int64 // 任务计数
	maxSize int64
}

// NewLockFreeQueue 创建泛型无锁队列
func NewLockFreeQueue[T any]() *LockFreeQueue[T] {
	dummy := unsafe.Pointer(&LockFreeNode[T]{})
	return &LockFreeQueue[T]{head: dummy, tail: dummy}
}

// Enqueue 入队
func (q *LockFreeQueue[T]) Enqueue(value T) error {

	if q.maxSize > 0 && atomic.LoadInt64(&q.length) >= q.maxSize {
		return ErrQueueFull
	}

	node := &LockFreeNode[T]{value: value}
	retries := 0
	for {
		tail := atomic.LoadPointer(&q.tail)
		next := atomic.LoadPointer(&(*LockFreeNode[T])(tail).next)
		if tail != atomic.LoadPointer(&q.tail) {
			// 快速失败
			if retries > 5 {
				time.Sleep(time.Duration(retries) * time.Nanosecond) // 初始轻退避
			}
			retries++
			continue
		}
		if next == nil {
			if atomic.CompareAndSwapPointer(&(*LockFreeNode[T])(tail).next, nil, unsafe.Pointer(node)) {
				atomic.AddInt64(&q.length, 1)
				atomic.CompareAndSwapPointer(&q.tail, tail, unsafe.Pointer(node))
				return nil
			}
		} else {
			atomic.CompareAndSwapPointer(&q.tail, tail, next)
		}
		// 统一线性退避
		if retries > 0 {
			delay := retries * 10
			if delay > 100 { // 最大 100 ns
				delay = 100
			}
			time.Sleep(time.Duration(delay) * time.Nanosecond)
		}
		retries++
	}
}

// Dequeue 出队
func (q *LockFreeQueue[T]) Dequeue() (T, bool) {
	var zero T // 默认零值
	retries := 0
	for {
		head := atomic.LoadPointer(&q.head)
		tail := atomic.LoadPointer(&q.tail)
		next := atomic.LoadPointer(&(*LockFreeNode[T])(head).next)
		if head != atomic.LoadPointer(&q.head) {
			// 快速失败，重试
			if retries > 5 {
				time.Sleep(time.Duration(retries) * time.Nanosecond) // 线性退避
			}
			retries++
			continue
		}
		if head == tail {
			if next == nil {
				return zero, false // 队列为空
			}
			// 帮助移动 tail
			atomic.CompareAndSwapPointer(&q.tail, tail, next)
		} else {
			value := (*LockFreeNode[T])(next).value
			if atomic.CompareAndSwapPointer(&q.head, head, next) {
				atomic.AddInt64(&q.length, -1)
				return value, true
			}
		}
		// 退避优化：线性增长，限制上限
		if retries > 0 {
			delay := retries * 10
			if delay > 100 { // 最大 100 ns
				delay = 100
			}
			time.Sleep(time.Duration(delay) * time.Nanosecond)
		}
		retries++
	}
}

// Len 返回队列长度
func (q *LockFreeQueue[T]) Len() int64 {
	return atomic.LoadInt64(&q.length)
}
