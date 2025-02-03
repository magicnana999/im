package core

import (
	"context"
	"fmt"
	"github.com/magicnana999/im/broker/handler"
	"github.com/magicnana999/im/broker/state"
	"github.com/magicnana999/im/logger"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	//bbPool "github.com/panjf2000/gnet/v2/pkg/pool/bytebuffer"
	goPool "github.com/panjf2000/gnet/v2/pkg/pool/goroutine"

	"time"
)

type BrokerServer struct {
	*gnet.BuiltinEventEngine
	eng               gnet.Engine
	name              string
	multicore         bool
	async             bool
	writev            bool
	nclients          int
	started           int32
	connected         int32
	connectionActive  int32
	workerPool        *goPool.Pool
	heartbeatPool     *goPool.Pool
	ctx               context.Context
	interval          time.Duration
	heartbeatInterval time.Duration
}

func Start(ctx context.Context, option *Option) {
	ts := &BrokerServer{
		name:              option.Name,
		multicore:         true,
		async:             true,
		writev:            true,
		nclients:          2,
		workerPool:        goPool.Default(),
		heartbeatPool:     goPool.Default(),
		ctx:               ctx,
		interval:          option.ServerInterval,
		heartbeatInterval: option.HeartbeatInterval,
	}
	err := gnet.Run(ts, fmt.Sprintf("tcp://0.0.0.0:%s", DefaultPort),
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
}

func (s *BrokerServer) OnBoot(eng gnet.Engine) (action gnet.Action) {
	logger.Info("BrokerInstance started")

	s.eng = eng

	broker := state.BrokerInfo{
		Addr:    s.name,
		StartAt: time.Now().UnixMilli(),
	}
	state.SetBroker(s.ctx, broker)
	return gnet.None
}

func (s *BrokerServer) OnShutdown(eng gnet.Engine) {
	logger.Info("BrokerInstance shutdown")
}

func (s *BrokerServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	uc := state.EmptyUserConnection(c)
	subCtx := context.WithValue(context.Background(), state.CurrentUser, &uc)
	c.SetContext(subCtx)

	logger.InfoF("[%s#%s] Connection open", c.RemoteAddr().String(), uc.Label())

	handler.DefaultHeartbeatHandler.StartTicker(subCtx, c, &uc)
	return nil, gnet.None
}

func (s *BrokerServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	user, err := state.CurrentUserFromConn(c)
	if err != nil {
		logger.Error(err)
	}

	handler.DefaultHeartbeatHandler.StopTicker(c)

	logger.InfoF("[%s#%s] Connection close", c.RemoteAddr().String(), user.Label())
	return
}

func (s *BrokerServer) OnTraffic(c gnet.Conn) (action gnet.Action) {

	ctx := c.Context().(context.Context)

	user, err := state.CurrentUserFromConn(c)
	if err != nil {
		logger.Error(err)
	}

	logger.InfoF("[%s#%s] Connection traffic", c.RemoteAddr().String(), user.Label())

	packets, e := defaultCodec.Decode(s, c)
	if e != nil {

		logger.FatalF("[%s#%s] Connection traffic error:%v",
			c.RemoteAddr().String(),
			user.Label(),
			e)
	}

	if packets != nil {
		for _, packet := range packets {
			response, err11 := handler.DefaultHandler.HandlePacket(ctx, packet)
			if err11 != nil {
				logger.FatalF("[%s#%s] Connection traffic error:%v", c.RemoteAddr().String(), user.Label(), err11)
			}

			bs, err12 := defaultCodec.Encode(response)
			if err12 != nil {
				logger.FatalF("[%s#%s] Connection traffic error:%v", c.RemoteAddr().String(), user.Label(), err12)
			}

			_ = c.AsyncWritev(bs, func(c gnet.Conn, err error) error {
				if c.RemoteAddr() != nil {
					logging.Debugf("[%s#%s] Connection traffic,write done", c.RemoteAddr().String(), err)
				}
				return nil
			})
		}
	}

	return gnet.None
}

func (s *BrokerServer) OnTick() (delay time.Duration, action gnet.Action) {

	broker := state.BrokerInfo{
		Addr:    s.name,
		StartAt: time.Now().UnixMilli(),
	}
	state.RefreshBroker(s.ctx, broker)

	return s.interval, gnet.None
}
