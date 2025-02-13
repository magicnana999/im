package broker

import (
	"context"
	"github.com/magicnana999/im/conf"
	"github.com/magicnana999/im/domain"
	"github.com/magicnana999/im/enum"
	"github.com/magicnana999/im/logger"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	bb "github.com/panjf2000/gnet/v2/pkg/pool/bytebuffer"
	//bbPool "github.com/panjf2000/gnet/v2/pkg/pool/bytebuffer"
	goPool "github.com/panjf2000/gnet/v2/pkg/pool/goroutine"

	"time"
)

var defaultInstance *Instance

type Instance struct {
	*gnet.BuiltinEventEngine
	eng              gnet.Engine
	addr             string
	multicore        bool
	async            bool
	writev           bool
	nclients         int
	started          int32
	connected        int32
	connectionActive int32
	workerPool       *goPool.Pool
	heartbeatPool    *goPool.Pool
	ctx              context.Context
	cancel           context.CancelFunc
	interval         int
	handler          packetHandler
	heartbeatHandler *heartbeatHandler
	brokerState      *brokerState
	codec            codec
	deliver          *deliver
}

func Start(ctx context.Context, cancel context.CancelFunc) {

	i := conf.Global.Broker.ServerInterval
	if i <= 0 {
		i = 30
	}

	ts := &Instance{
		addr:             conf.Global.Broker.Addr,
		multicore:        true,
		async:            true,
		writev:           true,
		nclients:         2,
		workerPool:       goPool.Default(),
		heartbeatPool:    goPool.Default(),
		ctx:              ctx,
		cancel:           cancel,
		interval:         i,
		handler:          initHandler(),
		brokerState:      initBrokerState(),
		heartbeatHandler: initHeartbeatHandler(),
		codec:            initCodec(),
		deliver:          initDeliver(ctx, initCodec()),
	}

	defaultInstance = ts

	err := gnet.Run(ts, ts.addr,
		gnet.WithMulticore(true),
		gnet.WithLockOSThread(true),
		gnet.WithReadBufferCap(4096),
		gnet.WithWriteBufferCap(4096),
		gnet.WithLoadBalancing(gnet.RoundRobin),
		gnet.WithNumEventLoop(1),
		gnet.WithReuseAddr(true),
		gnet.WithReusePort(true),
		gnet.WithTCPKeepAlive(time.Minute),
		gnet.WithTCPNoDelay(gnet.TCPNoDelay),
		gnet.WithSocketRecvBuffer(4096),
		gnet.WithSocketSendBuffer(4096),
		gnet.WithTicker(true),
		gnet.WithLogger(logger.Logger),
		gnet.WithLogLevel(logging.DebugLevel),
		gnet.WithEdgeTriggeredIO(true),
		gnet.WithEdgeTriggeredIOChunk(0))
	if err != nil {
		logger.FatalF("Start BrokerInstance error: %v", err)
	}

	ts.deliver.start()
}

func (s *Instance) OnBoot(eng gnet.Engine) (action gnet.Action) {
	logger.Info("BrokerInstance started")

	s.eng = eng

	brokerInfo := domain.BrokerInfo{Addr: s.addr, StartAt: time.Now().UnixMilli()}
	if _, e := s.brokerState.StoreBroker(s.ctx, brokerInfo); e != nil {

		logger.FatalF("failed to store broker info: %v", e)
	}

	return gnet.None
}

func (s *Instance) OnShutdown(eng gnet.Engine) {
	s.cancel()
	s.heartbeatHandler.stopTickerAll()
}

func (s *Instance) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	uc := &domain.UserConnection{
		Fd:          c.Fd(),
		AppId:       "",
		UserId:      0,
		ClientAddr:  c.RemoteAddr().String(),
		BrokerAddr:  c.LocalAddr().String(),
		OS:          enum.OSType(0),
		ConnectTime: time.Now().UnixMilli(),
		C:           c,
	}

	subCtx := context.WithValue(s.ctx, currentUserKey, uc)
	c.SetContext(subCtx)

	logger.InfoF("[%s#%s] Connection open", c.RemoteAddr().String(), uc.Label())

	if err := s.heartbeatHandler.startTicker(subCtx, c, uc); err != nil {
		logger.ErrorF("Connection open error: %v", err)
		s.heartbeatHandler.stopTicker(c)
	}
	return nil, gnet.None
}

func (s *Instance) OnClose(c gnet.Conn, err error) (action gnet.Action) {

	s.heartbeatHandler.stopTicker(c)
	return
}

func (s *Instance) OnTraffic(c gnet.Conn) (action gnet.Action) {

	ctx := c.Context().(context.Context)

	user, err := currentUserFromConn(c)
	if err != nil {
		logger.ErrorF("Connection traffic error: %v", err)
		s.heartbeatHandler.stopTicker(c)
		return gnet.None
	}

	logger.DebugF("[%s#%s] Connection traffic", c.RemoteAddr().String(), user.Label())

	packets, e := s.codec.decode(c)
	if e != nil {

		logger.ErrorF("[%s#%s] Connection traffic error:%v",
			c.RemoteAddr().String(),
			user.Label(),
			e)
		s.heartbeatHandler.stopTicker(c)
		return gnet.None
	}

	if packets != nil {
		for _, packet := range packets {
			response, err11 := s.handler.handlePacket(ctx, packet)
			if err11 != nil {
				logger.ErrorF("[%s#%s] Connection traffic error:%v", c.RemoteAddr().String(), user.Label(), err11)
				s.heartbeatHandler.stopTicker(c)
				return gnet.None
			}

			if response == nil {
				continue
			}

			bs, err12 := s.codec.encode(c, response)
			if err12 != nil {
				logger.ErrorF("[%s#%s] Connection traffic error:%v", c.RemoteAddr().String(), user.Label(), err12)
				s.heartbeatHandler.stopTicker(c)
				return gnet.None
			}

			c.AsyncWrite(bs.Bytes(), func(c gnet.Conn, err error) error {
				if err != nil {
					logging.Fatalf("[%s#%s] Connection traffic,write error:%v", c.RemoteAddr().String(), user.Label(), err)
					s.heartbeatHandler.stopTicker(c)
				}
				return err
			})

			bb.Put(bs)
		}
	}

	return gnet.None
}

func (s *Instance) OnTick() (delay time.Duration, action gnet.Action) {

	broker := domain.BrokerInfo{
		Addr:    s.addr,
		StartAt: time.Now().UnixMilli(),
	}
	_, err := s.brokerState.RefreshBroker(s.ctx, broker)
	if err != nil {
		logger.FatalF("BrokerInstance ticking error: %v", err)
	}

	return time.Duration(s.interval) * time.Second, gnet.None
}
