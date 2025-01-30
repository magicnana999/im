package handler

import (
	"context"
	"github.com/magicnana999/im/broker/protocol"
	"github.com/magicnana999/im/logger"
	"github.com/magicnana999/im/state"
	"github.com/panjf2000/ants/v2"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	goPool "github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
	"sync"
	"time"
)

const (
	heartbeatInterval = 10 * time.Second
)

var DefaultHeartbeatHandler = &HeartbeatHandler{}

type HeartbeatHandler struct {
	heartbeatPool *goPool.Pool
	m             map[int]*HeartbeatTask
	mu            sync.RWMutex
}

func (h *HeartbeatHandler) HandlePacket(c gnet.Conn, packet *protocol.Packet) error {
	h.SetLastHeartbeat(c)
	return nil
}

func (h *HeartbeatHandler) IsSupport(packetType int32) bool {
	return packetType == protocol.TypeHeartbeat
}

func (h *HeartbeatHandler) InitHandler() {
	var (
		DefaultAntsPoolSize = 1 << 18
		ExpiryDuration      = 10 * time.Second
		Nonblocking         = true
	)

	options := ants.Options{
		ExpiryDuration: ExpiryDuration,
		Nonblocking:    Nonblocking,
		Logger:         logger.Logger,
		PanicHandler: func(a any) {
			logging.Errorf("goroutine pool panic: %v", a)
		},
	}
	defaultAntsPool, err := ants.NewPool(DefaultAntsPoolSize, ants.WithOptions(options))

	if err != nil {
		logger.FatalF("Init default ants pool error: %v", err)
	}

	DefaultHeartbeatHandler.heartbeatPool = defaultAntsPool
	DefaultHeartbeatHandler.m = make(map[int]*HeartbeatTask)

	logger.DebugF("HeartbeatHandler init")
}

func (h *HeartbeatHandler) Count() int {
	return len(h.m)
}

func (h *HeartbeatHandler) StartTicker(ctx context.Context, c gnet.Conn, uc *state.UserConnection) error {
	ct, cancel := context.WithCancel(ctx)

	task := &HeartbeatTask{
		uc:            uc,
		fd:            c.Fd(),
		remoteAddr:    c.RemoteAddr().String(),
		lastHeartbeat: time.Now().UnixMilli(),
		ctx:           ct,
		cancel:        cancel,
		ticker:        time.NewTicker(heartbeatInterval),
		c:             c,
		callback: func(c gnet.Conn) error {
			return c.Close()
		},
	}

	h.mu.Lock()
	h.m[c.Fd()] = task
	h.mu.Unlock()

	h.heartbeatPool.Submit(func() {

		defer task.ticker.Stop()

		for {
			select {
			case <-ct.Done():

				logger.InfoF("[%s#%s] HeartbeatTask done",
					task.remoteAddr,
					task.uc.Label())
				return

			case <-task.ticker.C:

				now := time.Now()
				if now.UnixMilli()-task.lastHeartbeat > heartbeatInterval.Milliseconds() {
					logger.InfoF("[%s#%s] HeartbeatTask timeout,now:%s,last:%s,interval:%d(ms)",
						task.remoteAddr,
						task.uc.Label(),
						format(now),
						format(time.UnixMilli(task.lastHeartbeat)),
						heartbeatInterval.Milliseconds())
					h.StopTicker(task.c)
				}
			}
		}
	})

	logger.InfoF("[%s#%s] HeartbeatTask started", task.remoteAddr, task.uc.Label())

	return nil
}

func format(time time.Time) string {
	return time.Format("2006-01-02 15:04:05.045")
}

func (h *HeartbeatHandler) StopTicker(c gnet.Conn) error {
	c.Close()

	h.mu.RLock()
	task, flag := h.m[c.Fd()]
	h.mu.RUnlock()

	if !flag {
		return nil
	}

	task.cancel()

	h.mu.Lock()
	delete(h.m, c.Fd())
	h.mu.Unlock()

	logger.InfoF("[%s#%s] HeartbeatTask closed", task.remoteAddr, task.uc.Label())

	return nil
}

func (h *HeartbeatHandler) SetLastHeartbeat(c gnet.Conn) error {
	h.mu.RLock()
	task, flag := h.m[c.Fd()]
	h.mu.RUnlock()
	if !flag {
		return nil
	}

	task.setLastHeartbeat()

	logger.InfoF("[%s#%s] HeartbeatTask setLastHeartbeat", task.remoteAddr, task.uc.Label())

	return nil
}

type HeartbeatTask struct {
	uc            *state.UserConnection
	fd            int
	remoteAddr    string
	lastHeartbeat int64
	ctx           context.Context
	cancel        context.CancelFunc
	ticker        *time.Ticker
	c             gnet.Conn
	callback      func(c gnet.Conn) error
}

func (t *HeartbeatTask) setLastHeartbeat() {
	t.lastHeartbeat = time.Now().UnixMilli()
}
