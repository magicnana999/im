package core

import (
	"context"
	"fmt"
	"github.com/magicnana999/im/broker/brokerstate"
	"github.com/magicnana999/im/logger"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	bbPool "github.com/panjf2000/gnet/v2/pkg/pool/bytebuffer"
	goPool "github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
	"github.com/timandy/routine"
	"sync/atomic"
	"time"
)

type BrokerServer struct {
	*gnet.BuiltinEventEngine
	eng          gnet.Engine
	name         string
	multicore    bool
	async        bool
	writev       bool
	nclients     int
	started      int32
	connected    int32
	disconnected int32
	clientActive int32
	workerPool   *goPool.Pool
	ctx          context.Context
	interval     time.Duration
}

func Start(ctx context.Context, option *Option) {
	ts := &BrokerServer{
		multicore:  true,
		async:      true,
		writev:     true,
		nclients:   2,
		workerPool: goPool.Default(),
		ctx:        ctx,
		interval:   option.Interval,
		name:       option.Name,
	}
	err := gnet.Run(ts, fmt.Sprintf("tcp://0.0.0.0:%s", DefaultPort),
		gnet.WithMulticore(true),
		gnet.WithLockOSThread(true),
		gnet.WithReadBufferCap(4096),
		gnet.WithWriteBufferCap(4096),
		gnet.WithLoadBalancing(gnet.RoundRobin),
		gnet.WithNumEventLoop(1),
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
		logger.FatalF("Start broker server error: %v", err)
	}
}

func (s *BrokerServer) OnBoot(eng gnet.Engine) (action gnet.Action) {
	logger.InfoF("[%d] Broker instance started", routine.Goid())

	s.eng = eng

	broker := &brokerstate.BrokerInfo{
		Addr:    s.name,
		StartAt: time.Now().UnixMilli(),
	}
	brokerstate.SetBroker(s.ctx, broker)
	return gnet.None
}

func (s *BrokerServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	ctx := context.Background()
	c.SetContext(ctx)
	atomic.AddInt32(&s.connected, 1)
	out = []byte("sweetness\r\n")

	logger.InfoF("[%d] Open connect", routine.Goid())

	return nil, gnet.None
}

func (s *BrokerServer) OnShutdown(eng gnet.Engine) {
	fd, err := s.eng.Dup()
	fmt.Println(fd, err)
	logger.InfoF("[%d] Broker instance shutdown", routine.Goid())

}

func (s *BrokerServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	if err != nil {
		logging.Debugf("error occurred on closed, %v\n", err)
	}
	atomic.AddInt32(&s.disconnected, 1)

	logger.InfoF("[%d] Cloud connect", routine.Goid())

	return
}

func (s *BrokerServer) OnTraffic(c gnet.Conn) (action gnet.Action) {

	logger.InfoF("[%d] Traffic connect", routine.Goid())

	if s.async {
		buf := bbPool.Get()

		n1, _ := c.WriteTo(buf)
		logger.InfoF("read ok %d %d %d", routine.Goid(), n1, buf.Len())

		if c.LocalAddr().Network() == "tcp" || c.LocalAddr().Network() == "unix" {
			// just for test
			n2 := c.InboundBuffered()
			n3 := c.OutboundBuffered()
			n4, _ := c.Discard(1)
			logger.InfoF("%d,inbound:%d outbound:%d discard:%d", routine.Goid(), n2, n3, n4)

			_ = s.workerPool.Submit(
				func() {
					if s.writev {
						mid := buf.Len() / 2
						bs := make([][]byte, 2)
						bs[0] = buf.B[:mid]
						bs[1] = buf.B[mid:]
						_ = c.AsyncWritev(bs, func(c gnet.Conn, err error) error {
							if c.RemoteAddr() != nil {
								logging.Debugf("conn=%s done writev: %v", c.RemoteAddr().String(), err)
							}
							bbPool.Put(buf)
							return nil
						})
					} else {
						_ = c.AsyncWrite(buf.Bytes(), func(c gnet.Conn, err error) error {
							if c.RemoteAddr() != nil {
								logging.Debugf("conn=%s done write: %v", c.RemoteAddr().String(), err)
							}
							bbPool.Put(buf)
							return nil
						})
					}
				})
			return
		} else if c.LocalAddr().Network() == "udp" {
			_ = s.workerPool.Submit(
				func() {
					_ = c.AsyncWrite(buf.Bytes(), nil)
				})
			return
		}
		return
	}

	buf, _ := c.Next(-1)
	if s.writev {
		mid := len(buf) / 2
		_, _ = c.Writev([][]byte{buf[:mid], buf[mid:]})
	} else {
		_, _ = c.Write(buf)
	}

	return
}

func (s *BrokerServer) OnTick() (delay time.Duration, action gnet.Action) {
	logger.InfoF("[%d] tick", routine.Goid())

	broker := &brokerstate.BrokerInfo{
		Addr:    s.name,
		StartAt: time.Now().UnixMilli(),
	}
	brokerstate.RefreshBroker(s.ctx, broker)
	return s.interval, gnet.None
}
