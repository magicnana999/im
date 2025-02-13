package core

import (
	"context"
	"fmt"
	"github.com/magicnana999/im/broker/handler"
	"github.com/magicnana999/im/broker/state"
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

const (
	DefaultPort = "7539"
)

type Option struct {
	Name              string        `json:"name"`
	ServerInterval    time.Duration `json:"serverInterval"`
	HeartbeatInterval time.Duration `json:"heartbeatInterval"`
}

type BrokerServer struct {
	*gnet.BuiltinEventEngine
	eng               gnet.Engine
	addr              string
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
	handler           handler.PacketHandler
	heartbeatHandler  *handler.HeartbeatHandler
	brokerState       *state.BrokerState
	codec             Codec
}

func Start(ctx context.Context, option *Option) {
	ts := &BrokerServer{
		addr:              option.Name,
		multicore:         true,
		async:             true,
		writev:            true,
		nclients:          2,
		workerPool:        goPool.Default(),
		heartbeatPool:     goPool.Default(),
		ctx:               ctx,
		interval:          option.ServerInterval,
		heartbeatInterval: option.HeartbeatInterval,
		handler:           handler.InitHandler(),
		brokerState:       state.InitBrokerState(),
		heartbeatHandler:  handler.InitHeartbeatHandler(),
		codec:             InitCodec(),
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

	brokerInfo := domain.BrokerInfo{Addr: s.addr, StartAt: time.Now().UnixMilli()}
	if _, e := s.brokerState.StoreBroker(s.ctx, brokerInfo); e != nil {

		logger.FatalF("BrokerInstance start error: %v", e)
	}

	return gnet.None
}

func (s *BrokerServer) OnShutdown(eng gnet.Engine) {
	logger.Info("BrokerInstance shutdown")
}

func (s *BrokerServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
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

	subCtx := context.WithValue(s.ctx, state.CurrentUser, uc)
	c.SetContext(subCtx)

	logger.InfoF("[%s#%s] Connection open", c.RemoteAddr().String(), uc.Label())

	if err := s.heartbeatHandler.StartTicker(subCtx, c, uc); err != nil {
		logger.ErrorF("Connection open error: %v", err)
		s.heartbeatHandler.StopTicker(c)
	}
	return nil, gnet.None
}

func (s *BrokerServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {

	user, err := state.CurrentUserFromConn(c)
	if err != nil {
		logger.ErrorF("Connection close error: %v", err)
	}

	s.heartbeatHandler.StopTicker(c)

	logger.InfoF("[%s#%s] Connection close", c.RemoteAddr().String(), user.Label())
	return
}

func (s *BrokerServer) OnTraffic(c gnet.Conn) (action gnet.Action) {

	ctx := c.Context().(context.Context)

	user, err := state.CurrentUserFromConn(c)
	if err != nil {
		logger.ErrorF("Connection traffic error: %v", err)
		s.heartbeatHandler.StopTicker(c)
		return gnet.None
	}

	logger.DebugF("[%s#%s] Connection traffic", c.RemoteAddr().String(), user.Label())

	packets, e := s.codec.Decode(c)
	if e != nil {

		logger.ErrorF("[%s#%s] Connection traffic error:%v",
			c.RemoteAddr().String(),
			user.Label(),
			e)
		s.heartbeatHandler.StopTicker(c)
		return gnet.None
	}

	if packets != nil {
		for _, packet := range packets {
			response, err11 := s.handler.HandlePacket(ctx, packet)
			if err11 != nil {
				logger.ErrorF("[%s#%s] Connection traffic error:%v", c.RemoteAddr().String(), user.Label(), err11)
				s.heartbeatHandler.StopTicker(c)
				return gnet.None
			}

			if response == nil {
				continue
			}

			bs, err12 := s.codec.Encode(c, response)
			if err12 != nil {
				logger.ErrorF("[%s#%s] Connection traffic error:%v", c.RemoteAddr().String(), user.Label(), err12)
				s.heartbeatHandler.StopTicker(c)
				return gnet.None
			}

			c.AsyncWrite(bs.Bytes(), func(c gnet.Conn, err error) error {
				if err != nil {
					logging.Fatalf("[%s#%s] Connection traffic,write error:%v", c.RemoteAddr().String(), user.Label(), err)
					s.heartbeatHandler.StopTicker(c)
				}
				return err
			})

			bb.Put(bs)
		}
	}

	return gnet.None
}

func (s *BrokerServer) OnTick() (delay time.Duration, action gnet.Action) {

	broker := domain.BrokerInfo{
		Addr:    s.addr,
		StartAt: time.Now().UnixMilli(),
	}
	_, err := s.brokerState.RefreshBroker(s.ctx, broker)
	if err != nil {
		logger.FatalF("BrokerInstance ticking error: %v", err)
	}

	return s.interval, gnet.None
}
