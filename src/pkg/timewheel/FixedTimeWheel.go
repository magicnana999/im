package timewheel

import (
	"sync"
	"time"
)

type TimeWheelQueue interface {
}

type TimeWheel struct {
	slots   [60]*LockFreeQueue     // 60 个 slot
	current int                    // 当前槽索引
	workers int                    // 工作协程数
	ticker  *time.Ticker           // 每秒触发
	stopCh  chan struct{}          // 停止信号
	wg      sync.WaitGroup         // 同步工作协程
	conns   map[string]*Connection // 连接池
	connMu  sync.RWMutex           // 连接池锁
}
