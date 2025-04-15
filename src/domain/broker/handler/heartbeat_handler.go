package broker

import (
	"context"
	"github.com/magicnana999/im/api/kitex_gen/api"
	dto "github.com/magicnana999/im/broker/domain"
	"github.com/magicnana999/im/domain/broker"
	"github.com/magicnana999/im/domain/broker/state"
	"github.com/magicnana999/im/errors"
	"github.com/magicnana999/im/global"
	"github.com/magicnana999/im/pb"
	"github.com/magicnana999/im/pkg/logger"
	"github.com/panjf2000/gnet/v2"
	goPool "github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
	"sync"
	"time"
)

var defaultHeartbeatHandler *heartbeatHandler
var hhOnce sync.Once

type heartbeatHandler struct {
	heartbeatPool *goPool.Pool
	userState     *state.UserState
	interval      time.Duration
	m             sync.Map
}

func initHeartbeatHandler() *heartbeatHandler {

	hhOnce.Do(func() {

		interval := global.GetBroker().HeartbeatInterval
		if interval <= 0 {
			interval = 30
		}
		defaultHeartbeatHandler = &heartbeatHandler{
			heartbeatPool: goPool.Default(),
			userState:     state.InitUserState(),
			interval:      time.Duration(interval) * time.Second,
		}
	})

	return defaultHeartbeatHandler
}

func (h *heartbeatHandler) handlePacket(ctx context.Context, packet *api.Packet) (*api.Packet, error) {
	uc, err := broker.UserFromCtx(ctx)
	if err != nil {
		return nil, errors.HeartbeatError.SetDetail(err)
	}
	h.setLastHeartbeat(uc)

	return api.NewHeartbeat(int32(1)), nil
}

func (h *heartbeatHandler) isSupport(ctx context.Context, packetType int32) bool {
	return packetType == pb.TypeHeartbeat
}

func (h *heartbeatHandler) startTicker(ctx context.Context, c gnet.Conn, uc *dto.UserConnection) error {

	task := &heartbeatTask{fd: c.Fd()}

	_, exist := h.m.LoadOrStore(task.fd, task)
	if exist {
		task = nil
		return nil
	}

	ct, cancel := context.WithCancel(ctx)

	task.uc = uc
	task.fd = c.Fd()
	task.remoteAddr = c.RemoteAddr().String()
	task.lastHeartbeat = time.Now().UnixMilli()
	task.ctx = ct
	task.cancel = cancel
	task.ticker = time.NewTicker(h.interval)
	task.c = c

	err := h.heartbeatPool.Submit(func() {

		defer task.ticker.Stop()

		for {
			select {
			case <-ct.Done():
				return

			case <-task.ticker.C:

				now := time.Now()
				if now.UnixMilli()-task.lastHeartbeat > h.interval.Milliseconds() {

					logger.Errorf("[%s#%s] heartbeatTask timeout,now:%s,last:%s,interval:%d(ms)",
						task.remoteAddr,
						task.uc.Label(),
						format(now),
						format(time.UnixMilli(task.lastHeartbeat)),
						h.interval.Milliseconds())

					h.stopTicker(task.c)
				} else {
					h.userState.RefreshUser(task.ctx, task.uc)
				}
			}
		}
	})

	if err != nil {

		errorMap := make(map[string]any)
		errorMap["ucLabel"] = task.uc.Label()
		errorMap["remoteAddr"] = task.remoteAddr
		errorMap["heartbeatPoolCap"] = h.heartbeatPool.Cap()
		errorMap["heartbeatPoolRunning"] = h.heartbeatPool.Running()
		errorMap["heartbeatPoolFree"] = h.heartbeatPool.Free()
		errorMap["error"] = err.Error()

		return errors.HeartbeatTaskError.SetDetail(errorMap)
	}

	return nil
}

func format(time time.Time) string {
	return time.Format("2006-01-02 15:04:05.045")
}

func (h *heartbeatHandler) stopTickerAll() {
	h.m.Range(func(key, value interface{}) bool {
		task := value.(*heartbeatTask)
		h.stopTicker(task.c)
		return true
	})
}

func (h *heartbeatHandler) stopTicker(c gnet.Conn) error {
	defer func() {
		if c != nil {
			c.Close()
		}
	}()

	tsk, flag := h.m.Load(c.Fd())

	if !flag {
		return nil
	}

	t := tsk.(*heartbeatTask)
	t.cancel()

	h.m.Delete(c.Fd())

	if c != nil {
		c.Close()
	}
	return nil
}

func (h *heartbeatHandler) setLastHeartbeat(c *dto.UserConnection) {
	task, _ := h.m.Load(c.Fd)
	if task == nil {
		return
	}

	task.(*heartbeatTask).setLastHeartbeat()
}

func (h *heartbeatHandler) isRunning(fd int) bool {
	_, exist := h.m.Load(fd)
	return exist
}

type heartbeatTask struct {
	uc            *dto.UserConnection
	fd            int
	remoteAddr    string
	lastHeartbeat int64
	ctx           context.Context
	cancel        context.CancelFunc
	ticker        *time.Ticker
	c             gnet.Conn
	//closeFunc     func(c gnet.Conn) error
}

func (t *heartbeatTask) setLastHeartbeat() {
	t.lastHeartbeat = time.Now().UnixMilli()
}
