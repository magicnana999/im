package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/magicnana999/im/broker/handler"
	"github.com/magicnana999/im/logger"
	"github.com/magicnana999/im/state"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	goPool "github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
	"time"
)

const (
	CurrentUser string = `CurrentUser`
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
		interval:          option.TickInterval,
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
	subCtx := context.WithValue(s.ctx, CurrentUser, &uc)
	c.SetContext(subCtx)

	logger.InfoF("[%s#%s] Connection open", c.RemoteAddr().String(), uc.Label())

	handler.DefaultHeartbeatHandler.StartTicker(subCtx, c, &uc)
	return nil, gnet.None
}

func (s *BrokerServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	user, err := currentUserConnection(c)
	if err != nil {
		logger.Error(err)
	}

	logger.InfoF("[%s#%s] Connection close", c.RemoteAddr().String(), user.Label())
	return
}

func (s *BrokerServer) OnTraffic(c gnet.Conn) (action gnet.Action) {

	user, err := currentUserConnection(c)
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
			handler.DefaultHandler.HandlePacket(c, packet)
		}
	}

	//if s.async {
	//	buf := bbPool.Get()
	//
	//	n1, _ := c.WriteTo(buf)
	//	logger.InfoF("read ok %d %d %d", routine.Goid(), n1, buf.Len())
	//
	//	if c.LocalAddr().Network() == "tcp" || c.LocalAddr().Network() == "unix" {
	//		// just for test
	//		n2 := c.InboundBuffered()
	//		n3 := c.OutboundBuffered()
	//		n4, _ := c.Discard(1)
	//		logger.InfoF("%d,inbound:%d outbound:%d discard:%d", routine.Goid(), n2, n3, n4)
	//
	//		_ = s.workerPool.Submit(
	//			func() {
	//				if s.writev {
	//					mid := buf.Len() / 2
	//					bs := make([][]byte, 2)
	//					bs[0] = buf.B[:mid]
	//					bs[1] = buf.B[mid:]
	//					_ = c.AsyncWritev(bs, func(c gnet.Conn, err error) error {
	//						if c.RemoteAddr() != nil {
	//							logging.Debugf("conn=%s done writev: %v", c.RemoteAddr().String(), err)
	//						}
	//						bbPool.Put(buf)
	//						return nil
	//					})
	//				} else {
	//					_ = c.AsyncWrite(buf.Bytes(), func(c gnet.Conn, err error) error {
	//						if c.RemoteAddr() != nil {
	//							logging.Debugf("conn=%s done write: %v", c.RemoteAddr().String(), err)
	//						}
	//						bbPool.Put(buf)
	//						return nil
	//					})
	//				}
	//			})
	//		return
	//	} else if c.LocalAddr().Network() == "udp" {
	//		_ = s.workerPool.Submit(
	//			func() {
	//				_ = c.AsyncWrite(buf.Bytes(), nil)
	//			})
	//		return
	//	}
	//	return
	//}
	//
	//buf, _ := c.Next(-1)
	//if s.writev {
	//	mid := len(buf) / 2
	//	_, _ = c.Writev([][]byte{buf[:mid], buf[mid:]})
	//} else {
	//	_, _ = c.Write(buf)
	//}
	//
	//return
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

func currentUserConnection(c gnet.Conn) (*state.UserConnection, error) {
	if ctx, o := c.Context().(context.Context); o {
		if u, ok := ctx.Value(CurrentUser).(*state.UserConnection); ok {
			return u, nil
		}
	}

	return nil, errors.New("not found user")
}
