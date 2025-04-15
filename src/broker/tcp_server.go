package broker

import (
	"context"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/magicnana999/im/broker/domain"
	"github.com/magicnana999/im/broker/handler"
	"github.com/magicnana999/im/broker/holder"
	"github.com/magicnana999/im/pkg/singleton"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	bb "github.com/panjf2000/gnet/v2/pkg/pool/bytebuffer"
	"google.golang.org/protobuf/encoding/protojson"

	"time"
)

var tcpServerSingleton = singleton.NewSingleton[*TcpServer]()

type TcpServer struct {
	*gnet.BuiltinEventEngine
	eng              gnet.Engine
	interval         int
	commandHandler   *handler.CommandHandler
	heartbeatHandler *handler.HeartbeatHandler
	messageHandler   *handler.MessageHandler
	brokerHolder     *holder.BrokerHolder
	codec            codec
	ctx              context.Context
}

func NewTcpServer(mss *MessageSendServer, hts *HeartbeatServer) *TcpServer {
	return tcpServerSingleton.Get(func() *TcpServer {
		return &TcpServer{
			interval:         30,
			commandHandler:   handler.NewCommandHandler(),
			messageHandler:   handler.NewMessageHandler(mss),
			heartbeatHandler: handler.NewHeartbeatHandler(hts),
			brokerHolder:     holder.NewBrokerHolder(),
			codec:            newCodec(),
		}
	})
}

func (s *TcpServer) Start(ctx context.Context) error {
	s.ctx = ctx
	return gnet.Run(s,
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
}

func (s *TcpServer) OnBoot(eng gnet.Engine) (action gnet.Action) {
	s.eng = eng

	brokerInfo := domain.BrokerInfo{Addr: "", StartAt: time.Now().UnixMilli()}
	if _, e := s.brokerHolder.StoreBroker(s.ctx, brokerInfo); e != nil {
		logger.Fatalf("failed to store broker info: %v", e)
	}

	return gnet.None
}

func (s *TcpServer) OnShutdown(eng gnet.Engine) {
}

func (s *TcpServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	uc := &domain.UserConnection{
		Fd:          c.Fd(),
		AppId:       "",
		UserId:      0,
		ClientAddr:  c.RemoteAddr().String(),
		BrokerAddr:  c.LocalAddr().String(),
		ConnectTime: time.Now().UnixMilli(),
		C:           c,
	}

	subCtx := context.WithValue(s.ctx, currentUserKey, uc)
	c.SetContext(subCtx)

	if err := s.heartbeatHandler.StartHeartbeat(subCtx, c, uc); err != nil {
		logger.Errorf("failed to open connection: %v", err)
		s.eng.Stop(subCtx)
	}
	return nil, gnet.None
}

func (s *TcpServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	s.heartbeatHandler.stopTicker(c)
	return
}

func (s *TcpServer) OnTraffic(c gnet.Conn) (action gnet.Action) {

	ctx := c.Context().(context.Context)

	ctx = logger.NewSpan(ctx, "traffic")
	defer logger.EndSpan(ctx)

	packets, e := s.codec.decode(c)
	if e != nil {
		t, a := traffic(ctx, c, "decode error:%v", e)
		logger.Errorf(t, a...)

		s.heartbeatHandler.stopTicker(c)
		return gnet.None
	}

	if packets != nil {
		for _, packet := range packets {

			if !packet.IsHeartbeat() && logger.IsDebugEnable() {
				js, _ := protojson.Marshal(packet)
				t, a := traffic(ctx, c, "decode:%s", string(js))
				logger.Debugf(t, a...)
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
				t, a := traffic(ctx, c, "encode:%s", string(js))
				logger.Debugf(t, a...)
			}

			bs, err12 := s.codec.encode(c, response)
			if err12 != nil {
				t, a := traffic(ctx, c, "encode error:%v", err12)
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
		Addr:    s.addr,
		StartAt: time.Now().UnixMilli(),
	}
	_, err := s.brokerState.RefreshBroker(s.ctx, broker)
	if err != nil {
		logger.Fatalf("failed to ticking error: %v", err)
	}

	return time.Duration(s.interval) * time.Second, gnet.None
}
