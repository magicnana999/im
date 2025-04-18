package queue

import (
	"errors"
	"runtime"
	"sync/atomic"
	"unsafe"
)

var ErrQueueFull = errors.New("queue is full")
var ErrQueueEmpty = errors.New("queue is empty")

// lockFreeNode 泛型节点
type lockFreeNode[T any] struct {
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
func NewLockFreeQueue[T any](maxSize int64) *LockFreeQueue[T] {
	dummy := unsafe.Pointer(&lockFreeNode[T]{})
	return &LockFreeQueue[T]{head: dummy, tail: dummy, maxSize: maxSize}
}

// Enqueue 入队
func (q *LockFreeQueue[T]) Enqueue(value T) error {

	if q.maxSize > 0 && atomic.LoadInt64(&q.length) >= q.maxSize {
		return ErrQueueFull
	}

	node := &lockFreeNode[T]{value: value}

	retries := 0

	for {
		tail := atomic.LoadPointer(&q.tail)
		next := atomic.LoadPointer(&(*lockFreeNode[T])(tail).next)
		if tail != atomic.LoadPointer(&q.tail) {
			if retries > 100 {
				runtime.Gosched() // 高竞争时让出 CPU
			}
			retries++
			continue
		}
		if next == nil {
			if atomic.CompareAndSwapPointer(&(*lockFreeNode[T])(tail).next, nil, unsafe.Pointer(node)) {
				atomic.AddInt64(&q.length, 1)
				atomic.CompareAndSwapPointer(&q.tail, tail, unsafe.Pointer(node))
				return nil
			}
		} else {
			atomic.CompareAndSwapPointer(&q.tail, tail, next)
		}
		retries++

	}
}

// Dequeue 出队
func (q *LockFreeQueue[T]) Dequeue() (T, error) {

	var zero T // 默认零值

	if q.Len() == 0 {
		return zero, ErrQueueEmpty
	}

	retries := 0

	for {
		head := atomic.LoadPointer(&q.head)
		tail := atomic.LoadPointer(&q.tail)
		next := atomic.LoadPointer(&(*lockFreeNode[T])(head).next)
		if head != atomic.LoadPointer(&q.head) {

			if retries > 100 {
				runtime.Gosched() // 高竞争时让出 CPU
			}
			retries++
			continue
		}
		if head == tail {
			if next == nil {
				return zero, ErrQueueEmpty // 队列为空
			}
			atomic.CompareAndSwapPointer(&q.tail, tail, next)
		} else {
			value := (*lockFreeNode[T])(next).value
			if atomic.CompareAndSwapPointer(&q.head, head, next) {
				atomic.AddInt64(&q.length, -1)
				return value, nil
			}
		}
		retries++
	}
}

// Len 返回队列长度
func (q *LockFreeQueue[T]) Len() int64 {
	return atomic.LoadInt64(&q.length)
}
