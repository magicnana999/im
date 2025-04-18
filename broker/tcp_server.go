package broker

import (
	"context"
	"github.com/magicnana999/im/api/kitex_gen/api"
	"github.com/magicnana999/im/broker/domain"
	"github.com/magicnana999/im/broker/handler"
	"github.com/magicnana999/im/broker/holder"
	"github.com/magicnana999/im/pkg/timewheel"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	bb "github.com/panjf2000/gnet/v2/pkg/pool/bytebuffer"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"time"
)

type TcpServer struct {
	*gnet.BuiltinEventEngine
	eng            gnet.Engine
	interval       time.Duration
	hts            *HeartbeatServer
	mrs            *MessageRetryServer
	commandHandler *handler.CommandHandler
	messageHandler *handler.MessageHandler
	brokerHolder   *holder.BrokerHolder
	userHolder     *holder.UserHolder
	codec          *Codec
	ctx            context.Context
	logger         *Logger
}

func NewTcpServer(
	hts *HeartbeatServer,
	mrs *MessageRetryServer,
	ch *handler.CommandHandler,
	mh *handler.MessageHandler,
	bh *holder.BrokerHolder,
	uh *holder.UserHolder,
	lc fx.Lifecycle) (*TcpServer, error) {

	logger := NewLogger("tcp", true)
	ts := &TcpServer{
		hts:            hts,
		mrs:            mrs,
		commandHandler: ch,
		messageHandler: mh,
		brokerHolder:   bh,
		userHolder:     uh,
		codec:          NewCodec(),
		interval:       time.Second * 30,
		logger:         logger,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return ts.Start(context.Background())
		},
		OnStop: func(ctx context.Context) error {
			return ts.eng.Stop(ctx)
		},
	})

	return ts, nil
}

func (s *TcpServer) Start(ctx context.Context) error {
	s.ctx = ctx
	go func() {
		err := gnet.Run(s,
			"tcp://:5075",
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
			gnet.WithLogLevel(logging.DebugLevel),
			gnet.WithEdgeTriggeredIO(true),
			gnet.WithEdgeTriggeredIOChunk(0))

		s.logger.SrvInfo("tcp starting", SrvLifecycle, err)
	}()
	return nil
}

func (s *TcpServer) OnBoot(eng gnet.Engine) (action gnet.Action) {
	s.eng = eng
	s.logger.SrvInfo("tcp started", SrvLifecycle, nil)

	//brokerInfo := domain.BrokerInfo{Addr: "", StartAt: time.Now().UnixMilli()}
	//if _, e := s.brokerHolder.StoreBroker(s.ctx, brokerInfo); e != nil {
	//	logger.Fatalf("failed to store broker info: %v", e)
	//}

	return gnet.None
}

func (s *TcpServer) OnShutdown(eng gnet.Engine) {
	s.logger.SrvInfo("tcp shutdown", SrvLifecycle, nil)

}

func (s *TcpServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	uc := &domain.UserConn{
		Fd:          c.Fd(),
		AppId:       "",
		UserId:      0,
		ClientAddr:  c.RemoteAddr().String(),
		BrokerAddr:  c.LocalAddr().String(),
		ConnectTime: time.Now().UnixMilli(),
		Reader:      c,
		Writer:      c,
	}

	subCtx := context.WithValue(s.ctx, currentUserKey, uc)
	c.SetContext(subCtx)

	if err := s.openConn(subCtx, c, uc); err != nil {
		s.logger.ConnDebug("connect", uc.Desc(), ConnLifecycle, err)
		s.closeConn(c, uc, err.Error())
	}

	return nil, gnet.None
}

func (s *TcpServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	ctx := GetContext(c)
	uc, _ := GetCurUserConn(ctx)
	s.logger.ConnDebug("close", uc.Desc(), ConnLifecycle, err)
	return gnet.None
}

func (s *TcpServer) OnTraffic(c gnet.Conn) (action gnet.Action) {

	ctx := GetContext(c)
	uc, err := GetCurUserConn(ctx)
	if err != nil {
		s.logger.ConnDebug("read", "", ConnLifecycle, err)
		s.closeConn(c, uc, err.Error())
		return gnet.None
	}

	packets, err := s.codec.Decode(c)
	if err != nil {
		s.logger.ConnDebug("decode", uc.Desc(), ConnLifecycle, err)
		s.closeConn(c, uc, err.Error())
		return gnet.None
	}

	if packets != nil {
		for _, packet := range packets {

			ret, err := s.handle(ctx, packet)
			if err != nil {
				s.logger.ConnDebug("process packet", uc.Desc(), ConnLifecycle, err)
				s.closeConn(c, uc, err.Error())
			}

			if packet.IsHeartbeat() {
				ret = api.NewHeartbeatPacket(int32(1))
				uc.RefreshHeartbeat(time.Now().UnixMilli())
			}

			if packet.IsMessage() && packet.GetMessage().IsResponse() {
				s.mrs.Ack(packet.GetMessage().MessageId)
			}

			if ret != nil {
				s.response(c, uc, ret)
			}

		}
	}

	return gnet.None
}

func (s *TcpServer) handle(ctx context.Context, packet *api.Packet) (*api.Packet, error) {

	switch packet.Type {
	case api.TypeCommand:
		return s.commandHandler.HandlePacket(ctx, packet)
	case api.TypeMessage:
		return s.messageHandler.HandlePacket(ctx, packet)
	}
	return nil, nil
}
func (s *TcpServer) OnTick() (delay time.Duration, action gnet.Action) {

	//broker := domain.BrokerInfo{
	//	Addr:    "hello",
	//	StartAt: time.Now().UnixMilli(),
	//}
	//_, err := s.brokerHolder.RefreshBroker(s.ctx, broker)
	//if err != nil {
	//	s.logger.DebugOrError("tcp server could not tick", "", define.OpTicking, "", err)
	//}

	return time.Duration(s.interval) * time.Second, gnet.None
}

func GetContext(c gnet.Conn) context.Context {

	if c == nil {
		return nil
	}

	if ctx, o := c.Context().(context.Context); o {
		return ctx
	}
	return nil
}

func (s *TcpServer) DelContext(c gnet.Conn) {
	c.SetContext(nil)
}

func (s *TcpServer) openConn(ctx context.Context, c gnet.Conn, uc *domain.UserConn) error {
	s.userHolder.StoreTransient(uc)

	fun := func(now time.Time) timewheel.TaskResult {

		if uc.IsClosed.Load() {
			s.logger.ConnDebug("heartbeat break,conn is closed", uc.Desc(), ConnLifecycle, nil)
			s.closeConn(c, uc, "heartbeat,conn has been closed")
			return timewheel.Break
		}

		if time.Since(time.Unix(uc.LastHeartbeat.Load(), 0)) >= time.Second*60 {
			s.logger.ConnDebug("heartbeat break,timeout", uc.Desc(), ConnLifecycle, nil)
			s.closeConn(c, uc, "heartbeat timeout")
			return timewheel.Break
		}

		return timewheel.Retry
	}

	if _, err := s.hts.StartTicking(fun, time.Second*30); err != nil {
		return err
	}

	return nil
}

func (s *TcpServer) closeConn(c gnet.Conn, uc *domain.UserConn, reason string) {
	s.userHolder.RemoveUserConn(uc)
	uc.Close()
	ucd := ""
	if uc != nil {
		ucd = uc.Desc()
	}

	c.CloseWithCallback(func(c gnet.Conn, err error) error {
		s.logger.ConnDebug("close completed", ucd, ConnLifecycle, err, zap.String("reason", reason))
		s.DelContext(c)
		return nil
	})
}

func (s *TcpServer) response(c gnet.Conn, uc *domain.UserConn, packet *api.Packet) {
	bs, err := s.codec.Encode(packet)
	defer bb.Put(bs)

	if err != nil {
		s.logger.ConnDebug("encode", uc.Desc(), ConnLifecycle, err)
		s.closeConn(c, uc, err.Error())
		return
	}

	c.AsyncWrite(bs.Bytes(), func(c gnet.Conn, err error) error {
		s.logger.ConnDebug("write ack", uc.Desc(), ConnLifecycle, err)
		return err
	})

}
