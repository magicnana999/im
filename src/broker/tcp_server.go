package broker

import (
	"context"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/magicnana999/im/api/kitex_gen/api"
	"github.com/magicnana999/im/broker/domain"
	"github.com/magicnana999/im/broker/handler"
	"github.com/magicnana999/im/broker/holder"
	"github.com/magicnana999/im/define"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	bb "github.com/panjf2000/gnet/v2/pkg/pool/bytebuffer"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"

	"time"
)

type TcpServer struct {
	*gnet.BuiltinEventEngine
	eng              gnet.Engine
	interval         time.Duration
	commandHandler   *handler.CommandHandler
	heartbeatHandler *handler.HeartbeatHandler
	messageHandler   *handler.MessageHandler
	brokerHolder     *holder.BrokerHolder
	codec            *Codec
	ctx              context.Context
	logger           *Logger
}

func NewTcpServer(ch *handler.CommandHandler, hh *handler.HeartbeatHandler, mh *handler.MessageHandler, bh *holder.BrokerHolder, lc fx.Lifecycle) (*TcpServer, error) {

	logger := NewLogger("tcp", true)
	return &TcpServer{
		commandHandler:   ch,
		messageHandler:   mh,
		heartbeatHandler: hh,
		brokerHolder:     bh,
		codec:            NewCodec(),
		interval:         time.Second * 30,
		logger:           logger,
	}, nil
}

func (s *TcpServer) Start(ctx context.Context) error {
	s.ctx = ctx
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

	s.logger.SrvInfo("tcp start", SrvLifecycle, err)
	return err
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

	err := s.heartbeatHandler.StartHeartbeat(subCtx, uc)
	if err != nil {
		s.logger.ConnDebug("start heartbeat", uc.Desc(), ConnLifecycle, err)
		s.Close(c, uc, err.Error())
	}

	return nil, gnet.None
}

func (s *TcpServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	ctx := GetContext(c)
	uc := GetCurUserConn(ctx)
	s.logger.ConnDebug("close", uc.Desc(), ConnLifecycle, err)
	return gnet.None
}

func (s *TcpServer) OnTraffic(c gnet.Conn) (action gnet.Action) {

	ctx := GetContext(c)
	uc, err := GetCurUserConn(ctx)
	if err != nil {
		s.logger.ConnDebug("current user not exist", "", ConnLifecycle, err)
		s.Close(c, uc, err.Error())
	}
	if uc == nil {
		s.logger.ConnDebug("current user is nil", "", ConnLifecycle, UserConnIsNil)
		s.Close(c, uc, err.Error())
	}

	packets, err := s.codec.Decode(c)
	if err != nil {
		s.logger.ConnDebug("decode error", uc.Desc(), ConnLifecycle, e)
		s.heartbeatHandler.StopHeartbeat(ctx, uc)
		s.Close(c, uc, e.Error())
		return gnet.None
	}

	if packets != nil {
		for _, packet := range packets {

			var response *api.Packet

			switch packet.Type {
			case api.TypeHeartbeat:
				res, err := s.heartbeatHandler.HandlePacket(ctx, uc, packet)
				if err != nil {
					s.heartbeatHandler.StopHeartbeat(ctx, uc)
					c.Close()
				}
				response = res
				break
			case api.TypeCommand:
				res, err := s.commandHandler.HandlePacket(ctx, packet)
			}

			if packet.IsHeartbeat() {
				response, err := s.heartbeatHandler.HandlePacket(ctx, packet)
				if err != nil {
					s.heartbeatHandler.StopHeartbeat(ctx, uc)
					c.Close()
				}
			}

			response, err11 := s.handler.handlePacket(ctx, packet)
			if err11 != nil {
				t, a := traffic(ctx, c, "handle error:%v", err11)
				logger.Errorf(t, a...)

				s.heartbeatHandler.stopTicker(c)
				return gnet.None
			}

			if response == nil {
				continue
			}

			if !response.IsHeartbeat() && logger.IsDebugEnable() {
				js, _ := protojson.Marshal(response)
				t, a := traffic(ctx, c, "Encode:%s", string(js))
				logger.Debugf(t, a...)
			}

			bs, err12 := s.codec.Encode(c, response)
			if err12 != nil {
				t, a := traffic(ctx, c, "Encode error:%v", err12)
				logger.Errorf(t, a...)
				s.heartbeatHandler.stopTicker(c)
				return gnet.None
			}

			c.AsyncWrite(bs.Bytes(), func(c gnet.Conn, err error) error {
				if err != nil {
					t, a := traffic(ctx, c, "write error:%v", err12)
					logger.Errorf(t, a...)
					s.heartbeatHandler.stopTicker(c)
				}
				return err
			})

			bb.Put(bs)
		}
	}

	return gnet.None
}

func (s *TcpServer) OnTick() (delay time.Duration, action gnet.Action) {

	broker := domain.BrokerInfo{
		Addr:    "hello",
		StartAt: time.Now().UnixMilli(),
	}
	_, err := s.brokerHolder.RefreshBroker(s.ctx, broker)
	if err != nil {
		s.logger.DebugOrError("tcp server could not tick", "", define.OpTicking, "", err)
	}

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

func DelContext(c gnet.Conn) {
	c.SetContext(nil)
}

func (s *TcpServer) Close(c gnet.Conn, uc *domain.UserConn, reason string) {
	uc.Close()
	c.CloseWithCallback(func(c gnet.Conn, err error) error {
		s.logger.ConnDebug("close completed", uc.Desc(), ConnLifecycle, err, zap.String("reason", reason))
		DelContext(c)
		return nil
	})
}
