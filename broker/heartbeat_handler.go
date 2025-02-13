package broker

import (
	"context"
	"github.com/magicnana999/im/conf"
	"github.com/magicnana999/im/domain"
	"github.com/magicnana999/im/errors"
	"github.com/magicnana999/im/logger"
	"github.com/magicnana999/im/pb"
	"github.com/panjf2000/ants/v2"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	goPool "github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
	"sync"
	"time"
)

var heartbeatMu sync.RWMutex
var defaultHeartbeatHandler = &heartbeatHandler{}

type heartbeatHandler struct {
	heartbeatPool *goPool.Pool
	userState     *userState
	interval      time.Duration
	m             sync.Map
}

func (h *heartbeatHandler) handlePacket(ctx context.Context, packet *pb.Packet) (*pb.Packet, error) {
	uc, err := currentUserFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	h.setLastHeartbeat(uc)

	return pb.NewHeartbeat(int32(1)), nil
}

func (h *heartbeatHandler) isSupport(ctx context.Context, packetType int32) bool {
	return packetType == pb.TypeHeartbeat
}

func initHeartbeatHandler() *heartbeatHandler {

	heartbeatMu.Lock()
	defer heartbeatMu.Unlock()

	if defaultHeartbeatHandler.heartbeatPool != nil {
		return defaultHeartbeatHandler
	}

	var PoolSize = 1 << 18
	var ExpiryDuration = 10 * time.Second
	var Nonblocking = true

	options := ants.Options{
		ExpiryDuration: ExpiryDuration,
		Nonblocking:    Nonblocking,
		Logger:         logger.Logger,
		PanicHandler: func(a any) {
			logging.Errorf("goroutine pool panic: %v", a)
		},
	}
	defaultAntsPool, err := ants.NewPool(PoolSize, ants.WithOptions(options))

	if err != nil {
		logger.FatalF("init heartbeat ants pool error: %v", err)
	}

	interval := conf.Global.Broker.HeartbeatInterval
	if interval <= 0 {
		interval = 30
	}

	defaultHeartbeatHandler.heartbeatPool = defaultAntsPool
	defaultHeartbeatHandler.userState = initUserState()
	defaultHeartbeatHandler.interval = time.Duration(interval) * time.Second

	logger.DebugF("heartbeatHandler init")

	return defaultHeartbeatHandler
}

func (h *heartbeatHandler) startTicker(ctx context.Context, c gnet.Conn, uc *domain.UserConnection) error {

	task := &heartbeatTask{fd: c.Fd()}

	_, exist := h.m.LoadOrStore(task.fd, task)
	if exist {
		task = nil
		logger.WarnF("heartbeat deliverTask already exist,skipping deliverTask creation, fd:%d,remote:%s",
			c.Fd(),
			c.RemoteAddr().String())
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

				logger.InfoF("[%s#%s] heartbeatTask done",
					task.remoteAddr,
					task.uc.Label())
				return

			case <-task.ticker.C:

				now := time.Now()
				if now.UnixMilli()-task.lastHeartbeat > h.interval.Milliseconds() {
					logger.ErrorF("[%s#%s] heartbeatTask timeout,now:%s,last:%s,interval:%d(ms)",
						task.remoteAddr,
						task.uc.Label(),
						format(now),
						format(time.UnixMilli(task.lastHeartbeat)),
						h.interval.Milliseconds())
					h.stopTicker(task.c)
				} else {
					h.userState.refreshUser(task.ctx, task.uc)
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

		return errors.HeartbeatError.DetailJson(errorMap)
	}

	logger.DebugF("[%s#%s] heartbeatTask started", task.remoteAddr, task.uc.Label())

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
	logger.InfoF("[%s#%s] heartbeatTask closed", t.remoteAddr, t.uc.Label())

	return nil
}

func (h *heartbeatHandler) setLastHeartbeat(c *domain.UserConnection) {
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
	uc            *domain.UserConnection
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
