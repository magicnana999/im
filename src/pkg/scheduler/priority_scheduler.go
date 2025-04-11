package scheduler

import (
	"github.com/magicnana999/im/pkg/queue"
	"github.com/panjf2000/ants/v2"
)

type Priority string

type PriorityRunnable interface {
	GetPriority() Priority
	Run() error
}

type PriorityScheduler struct {
	highQueues *queue.LockFreeQueue[PriorityRunnable]
	midQueues  *queue.LockFreeQueue[PriorityRunnable]
	lowQueues  *queue.LockFreeQueue[PriorityRunnable]
	highPool   *ants.Pool
	midPool    *ants.Pool
	lowPool    *ants.Pool
}
