package timewheel

import (
	"container/heap"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

type Task struct {
	clientID   string
	expireTime int64
	index      int
}

type MinHeap []*Task

func (h *MinHeap) Len() int           { return len(*h) }
func (h *MinHeap) Less(i, j int) bool { return (*h)[i].expireTime < (*h)[j].expireTime }
func (h *MinHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
	(*h)[i].index = i
	(*h)[j].index = j
}
func (h *MinHeap) Push(x interface{}) {
	task := x.(*Task)
	task.index = len(*h)
	*h = append(*h, task)
}
func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	task := old[n-1]
	*h = old[0 : n-1]
	task.index = -1
	return task
}

type Heap struct {
	data       MinHeap
	mu         sync.RWMutex
	maxTasks   int
	threshold  int
	onOverflow func(*Task) bool
}

type Option struct {
	MaxTasks   int
	Threshold  int
	OnOverflow func(*Task) bool
}

func NewHeap(opt Option) *Heap {
	if opt.MaxTasks <= 0 {
		opt.MaxTasks = 10000
	}
	if opt.Threshold <= 0 || opt.Threshold > opt.MaxTasks {
		opt.Threshold = opt.MaxTasks / 2
	}
	if opt.OnOverflow == nil {
		opt.OnOverflow = func(*Task) bool { return false }
	}

	h := &Heap{
		data:       make(MinHeap, 0),
		maxTasks:   opt.MaxTasks,
		threshold:  opt.Threshold,
		onOverflow: opt.OnOverflow,
	}
	heap.Init(&h.data)
	return h
}

func (h *Heap) PushTask(clientID string, expireTime int64) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.data.Len() >= h.maxTasks {
		return errors.New("heap is full")
	}

	task := &Task{clientID: clientID, expireTime: expireTime}
	if h.data.Len() >= h.threshold {
		log.Printf("Heap size %d exceeds threshold %d, triggering overflow", h.data.Len(), h.threshold)
		if !h.onOverflow(task) {
			return errors.New("task rejected due to overflow")
		}
		return nil
	}

	heap.Push(&h.data, task)
	log.Printf("Pushed %s:%d, heap size: %d", clientID, expireTime, h.data.Len())
	return nil
}

func (h *Heap) PopTask() (*Task, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.data.Len() == 0 {
		return nil, errors.New("heap is empty")
	}

	task := heap.Pop(&h.data).(*Task)
	log.Printf("Popped %s:%d, heap size: %d", task.clientID, task.expireTime, h.data.Len())
	return task, nil
}

func (h *Heap) Peek() (*Task, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if h.data.Len() == 0 {
		return nil, errors.New("heap is empty")
	}
	return h.data[0], nil
}

func (h *Heap) Size() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.data.Len()
}

func main() {
	onOverflow := func(t *Task) bool {
		log.Printf("Overflow: moving %s:%d to another slot", t.clientID, t.expireTime)
		return true
	}

	h := NewHeap(Option{
		MaxTasks:   100,
		Threshold:  50,
		OnOverflow: onOverflow,
	})

	for i := 0; i < 60; i++ {
		clientID := fmt.Sprintf("client%d", i)
		expireTime := time.Now().UnixMilli() + int64(i*1000)
		err := h.PushTask(clientID, expireTime)
		if err != nil {
			log.Printf("Error: %v", err)
		}
	}

	if top, err := h.Peek(); err == nil {
		log.Printf("Top: %s:%d", top.clientID, top.expireTime)
	}

	for i := 0; i < 5; i++ {
		if task, err := h.PopTask(); err == nil {
			log.Printf("Popped: %s:%d", task.clientID, task.expireTime)
		}
	}

	log.Printf("Final size: %d", h.Size())
}
