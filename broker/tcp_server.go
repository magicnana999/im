package broker

import (
	"context"
	"fmt"
	"github.com/magicnana999/im/api/kitex_gen/api"
	brokerctx "github.com/magicnana999/im/broker/ctx"
	"github.com/magicnana999/im/broker/domain"
	"github.com/magicnana999/im/broker/handler"
	"github.com/magicnana999/im/broker/holder"
	"github.com/magicnana999/im/pkg/jsonext"
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
	fmt.Println("open", c.Fd())
	uc := domain.NewUserConn(c)

	err := s.openConn(c, uc)
	s.logger.ConnDebug("connect", uc.Desc(), ConnLifecycle, err)
	if err != nil {
		s.closeConn(c, uc, err.Error())
	}

	return nil, gnet.None
}

func (s *TcpServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	ctx := s.getContext(c)
	uc, _ := brokerctx.GetCurUserConn(ctx)
	s.logger.ConnDebug("close", uc.Desc(), ConnLifecycle, err)
	return gnet.None
}

func (s *TcpServer) OnTraffic(c gnet.Conn) (action gnet.Action) {

	fmt.Println("traffic", c.Fd())

	ctx := s.getContext(c)
	uc, err := brokerctx.GetCurUserConn(ctx)
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

			if packet.IsHeartbeat() {
				s.response(c, uc, api.NewHeartbeatPacket(int32(1)))
				uc.Refresh(time.Now().UnixMilli())
				continue
			}

			s.logger.PktDebug("read", uc.Desc(), string(jsonext.PbMarshalNoErr(packet)), packet.GetPacketId(), PacketTracking, nil)

			if packet.IsCommand() {
				ret, err := s.commandHandler.HandlePacket(ctx, packet)
				s.logger.PktDebug("command process", uc.Desc(), packet.GetPacketId(), string(jsonext.PbMarshalNoErr(packet)), PacketTracking, err)
				if err != nil {
					s.closeConn(c, uc, err.Error())
					continue
				}
				s.response(c, uc, ret)
				continue
			}

			if packet.IsMessage() {
				message := packet.GetMessage()
				if message.IsRequest() {
					ret, err := s.messageHandler.HandlePacket(ctx, packet)
					s.logger.PktDebug("message process", uc.Desc(), packet.GetPacketId(), string(jsonext.PbMarshalNoErr(packet)), PacketTracking, err)

					if err != nil {
						s.closeConn(c, uc, err.Error())
						continue
					}

					s.response(c, uc, ret)
					continue
				} else {
					err := s.mrs.Ack(message.MessageId)
					s.logger.PktDebug("ack process", uc.Desc(), packet.GetPacketId(), string(jsonext.PbMarshalNoErr(packet)), PacketTracking, err)
					if err != nil {
						s.closeConn(c, uc, err.Error())
						continue
					}
					continue
				}
			}
		}
	}

	return gnet.None
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

// initContext 新连接到来时，初始化ctx
func (s *TcpServer) initContext(c gnet.Conn, uc *domain.UserConn) context.Context {
	subCtx := context.WithValue(s.ctx, brokerctx.CurrentUserKey, uc)
	c.SetContext(subCtx)
	return subCtx
}

// 删除ctx
func (s *TcpServer) delContext(c gnet.Conn) {
	c.SetContext(nil)
}

// 获取ctx
func (s *TcpServer) getContext(c gnet.Conn) context.Context {
	if c == nil {
		return nil
	}

	if ctx, o := c.Context().(context.Context); o {
		return ctx
	}
	return nil
}

// 打开链接：初始化ctx、保存uc到本地、启动心跳
func (s *TcpServer) openConn(c gnet.Conn, uc *domain.UserConn) error {

	s.initContext(c, uc)

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

		s.logger.ConnDebug("heartbeat retry", uc.Desc(), ConnLifecycle, nil)
		return timewheel.Retry
	}

	if _, err := s.hts.StartTicking(fun, time.Second*30); err != nil {
		return err
	}

	s.logger.ConnDebug("open connection", uc.Desc(), ConnLifecycle, nil)
	return nil
}

// 打开链接：删除ctx、删除uc到本地、停止心跳
func (s *TcpServer) closeConn(c gnet.Conn, uc *domain.UserConn, reason string) {

	if !uc.Close() {
		return
	}

	s.userHolder.RemoveUserConn(uc)

	//停止心跳不会从hts中删除任务，把uc的IsClose = true，在hts中自己根据IsClose退出
	uc.Close()
	ucd := ""
	if uc != nil {
		ucd = uc.Desc()
	}

	err := c.CloseWithCallback(func(c gnet.Conn, err error) error {
		s.logger.ConnDebug("close completed", ucd, ConnLifecycle, err, zap.String("reason", reason))
		s.delContext(c)
		s.logger.ConnDebug("close connection", ucd, ConnLifecycle, nil)
		return nil
	})

	if err != nil {
		s.logger.ConnDebug("close error", ucd, ConnLifecycle, err)
	}
}

func (s *TcpServer) response(c gnet.Conn, uc *domain.UserConn, packet *api.Packet) {
	s.logger.PktDebug("write", uc.Desc(), packet.GetPacketId(), string(jsonext.PbMarshalNoErr(packet)), PacketTracking, nil)
	bs, err := s.codec.Encode(packet)
	defer bb.Put(bs)

	if err != nil {
		s.logger.PktDebug("encode error", uc.Desc(), packet.GetPacketId(), "", PacketTracking, err)
		s.closeConn(c, uc, err.Error())
		return
	}

	if err := c.AsyncWrite(bs.Bytes(), func(c gnet.Conn, err error) error {
		s.logger.PktDebug("write completed", uc.Desc(), packet.GetPacketId(), "", PacketTracking, err)
		return err
	}); err != nil {
		s.logger.PktDebug("write error", uc.Desc(), packet.GetPacketId(), "", PacketTracking, err)
		s.closeConn(c, uc, err.Error())
		return
	}

}
