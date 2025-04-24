package broker

import (
	"context"
	"fmt"
	"github.com/magicnana999/im/api/kitex_gen/api"
	brokerctx "github.com/magicnana999/im/broker/ctx"
	"github.com/magicnana999/im/broker/domain"
	"github.com/magicnana999/im/broker/handler"
	"github.com/magicnana999/im/broker/holder"
	"github.com/magicnana999/im/global"
	"github.com/magicnana999/im/pkg/jsonext"
	"github.com/magicnana999/im/pkg/timewheel"
	"github.com/panjf2000/ants/v2"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	bb "github.com/panjf2000/gnet/v2/pkg/pool/bytebuffer"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"runtime"
	"time"
)

type TcpServer struct {
	*gnet.BuiltinEventEngine
	eng            gnet.Engine
	cfg            *global.TCPConfig
	interval       time.Duration
	hts            *HeartbeatServer
	mrs            *MessageRetryServer
	commandHandler *handler.CommandHandler
	messageHandler *handler.MessageHandler
	brokerHolder   *holder.BrokerHolder
	userHolder     *holder.UserHolder
	codec          *Codec
	ctx            context.Context
	worker         *ants.Pool
	logger         *Logger
}

func getOrDefaultTCPConfig(g *global.Config) *global.TCPConfig {
	c := &global.TCPConfig{Worker: &global.TcpWorkerConfig{}, Heartbeat: &global.TcpHeartbeatConfig{}}
	if g != nil && g.TCP != nil {
		*c = *g.TCP
	}

	if c.Addr == "" {
		c.Addr = "0.0.0.0:5075"
	}

	if c.Interval <= 0 {
		c.Interval = time.Minute
	}

	if c.Heartbeat.Interval <= 0 {
		c.Heartbeat.Interval = time.Second * 30
	}

	if c.Heartbeat.Timeout <= 0 {
		c.Heartbeat.Timeout = time.Second * 30
	}

	if c.Worker.Size <= 0 {
		c.Worker.Size = runtime.NumCPU() * 1000
	}

	if c.Worker.ExpireDuration <= 0 {
		c.Worker.ExpireDuration = time.Minute
	}

	if c.Worker.MaxBlockingTasks <= 0 {
		c.Worker.MaxBlockingTasks = 100_000
	}

	return c
}

func newWorkerPool(c *global.TcpWorkerConfig) (*ants.Pool, error) {

	logger := NewLogger("tcp", true)

	panicHandler := func(panicErr interface{}) {
		logger.Error("tcp worker panic",
			zap.Any("error", panicErr),
			zap.Stack("stack"),
		)
	}

	return ants.NewPool(c.Size,
		ants.WithExpiryDuration(c.ExpireDuration),
		ants.WithMaxBlockingTasks(c.MaxBlockingTasks),
		ants.WithNonblocking(false),
		ants.WithPanicHandler(panicHandler),
		ants.WithLogger(logger),
		ants.WithPreAlloc(false),
		ants.WithDisablePurge(false),
	)
}

func NewTcpServer(
	conf *global.Config,
	hts *HeartbeatServer,
	mrs *MessageRetryServer,
	ch *handler.CommandHandler,
	mh *handler.MessageHandler,
	bh *holder.BrokerHolder,
	uh *holder.UserHolder,
	lc fx.Lifecycle) (*TcpServer, error) {

	logger := NewLogger("tcp", true)

	c := getOrDefaultTCPConfig(conf)

	worker, err := newWorkerPool(c.Worker)
	if err != nil {
		return nil, err
	}

	ts := &TcpServer{
		cfg:            c,
		hts:            hts,
		mrs:            mrs,
		commandHandler: ch,
		messageHandler: mh,
		brokerHolder:   bh,
		userHolder:     uh,
		codec:          NewCodec(),
		interval:       time.Second * 30,
		logger:         logger,
		worker:         worker,
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

func (s *TcpServer) Stop(ctx context.Context) error {

	s.userHolder.RangeAllUserConn(func(conn *domain.UserConn) bool {
		s.closeConnFD(conn.Conn, conn, "tcp server shutting down")
		return true
	})

	if !s.worker.IsClosed() {
		s.worker.ReleaseTimeout(time.Millisecond * 500)
	}
	return s.eng.Stop(ctx)
}

func (s *TcpServer) Start(ctx context.Context) error {
	s.ctx = ctx
	go func() {
		err := gnet.Run(s,
			fmt.Sprintf("tcp://%s", s.cfg.Addr),
			gnet.WithMulticore(true),
			gnet.WithLockOSThread(true),
			gnet.WithReadBufferCap(4096),
			gnet.WithWriteBufferCap(4096),
			gnet.WithLoadBalancing(gnet.RoundRobin),
			gnet.WithNumEventLoop(runtime.NumCPU()),
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

		if err != nil {
			s.logger.Fatal("tcp start failed", zap.Error(err))
		}
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
	uc := domain.NewUserConn(c)
	err := s.openConn(c, uc)
	s.logger.ConnDebug("connect", uc.Desc(), ConnLifecycle, err, zap.String("uc", string(jsonext.MarshalNoErr(uc))))
	if err != nil {
		//s.closeConnFD(c, uc, err.Error())
		return nil, gnet.Close
	}

	return nil, gnet.None
}

func (s *TcpServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {

	ctx := s.getContext(c)
	uc, err := brokerctx.GetCurUserConn(ctx)
	if err != nil {
		s.closeConn(ctx, c, uc)
	}
	s.logger.ConnDebug("close", uc.Desc(), ConnLifecycle, err)
	uc = nil

	return gnet.None
}

func (s *TcpServer) OnTraffic(c gnet.Conn) (action gnet.Action) {

	ctx := s.getContext(c)
	uc, err := brokerctx.GetCurUserConn(ctx)
	if err != nil {
		s.logger.ConnDebug("read", "", ConnLifecycle, err)
		//s.closeConnFD(c, uc, err.Error())
		return gnet.Close
	}

	s.RefreshUser(ctx, uc)

	packets, err := s.codec.Decode(c)

	if err != nil {
		s.logger.ConnDebug("decode", uc.Desc(), ConnLifecycle, err)
		//s.closeConnFD(c, uc, err.Error())
		return gnet.Close
	}

	err = s.worker.Submit(func() {
		for _, packet := range packets {
			resp := s.processPacket(ctx, c, uc, packet)
			if resp != nil {
				s.response(c, uc, resp)
			}
		}
	})

	if err != nil {
		s.logger.ConnDebug("submit decode failed", uc.Desc(), ConnLifecycle, err)
		//s.closeConnFD(c, uc, "server busy")
		return gnet.Close
	}
	return gnet.None
}

func (s *TcpServer) processPacket(ctx context.Context, c gnet.Conn, uc *domain.UserConn, packet *api.Packet) *api.Packet {
	if packet.IsHeartbeat() {
		s.response(c, uc, api.NewHeartbeatPacket(int32(1)))
		return nil
	}

	s.logger.PktDebug("read", uc.Desc(), packet.GetPacketId(), string(jsonext.PbMarshalNoErr(packet)), PacketTracking, nil)

	if packet.IsCommand() {
		return s.processCommand(ctx, c, uc, packet)
	}

	if packet.IsMessage() {
		if packet.GetMessage().IsRequest() {
			return s.processMessage(ctx, c, uc, packet)
		} else {
			s.mrs.Ack(packet.GetMessage().MessageId)
			s.logger.PktDebug("ack process", uc.Desc(), packet.GetPacketId(), string(jsonext.PbMarshalNoErr(packet)), PacketTracking, nil)
			return nil
		}
	}

	return nil
}

func (s *TcpServer) processMessage(ctx context.Context, c gnet.Conn, uc *domain.UserConn, packet *api.Packet) *api.Packet {
	message := packet.GetMessage()
	if message.IsRequest() {
		ret, err := s.messageHandler.HandlePacket(ctx, packet)
		s.logger.PktDebug("message process", uc.Desc(), packet.GetPacketId(), string(jsonext.PbMarshalNoErr(packet)), PacketTracking, err)
		return ret
	}

	return nil
}

func (s *TcpServer) processCommand(ctx context.Context, c gnet.Conn, uc *domain.UserConn, packet *api.Packet) *api.Packet {
	ret, err := s.commandHandler.HandlePacket(ctx, packet)
	s.logger.PktDebug("command process", uc.Desc(), packet.GetPacketId(), string(jsonext.PbMarshalNoErr(packet)), PacketTracking, err)
	if err == nil && packet.GetCommand().CommandType == api.CommandTypeUserLogin {
		s.OnUserLogin(ctx, uc, packet.GetCommand().GetLoginRequest(), ret.GetCommand().GetLoginReply())
	}
	return ret
}

func (s *TcpServer) OnTick() (delay time.Duration, action gnet.Action) {
	fmt.Println("--------- gnet.conns:", s.eng.CountConnections(), ",worker.running:", s.worker.Running())
	return time.Second, gnet.None
}

func (s *TcpServer) RefreshUser(ctx context.Context, uc *domain.UserConn) {

	if uc.IsClosed.Load() {
		return
	}
	uc.Refresh(time.Now())

	if uc.IsLogin.Load() {
		s.userHolder.RefreshUserConn(ctx, uc)
	}
}

func (s *TcpServer) OnUserLogin(
	ctx context.Context,
	uc *domain.UserConn,
	req *api.LoginRequest,
	rep *api.LoginReply) {

	if !uc.Login(rep.GetAppId(), rep.GetUserId(), req.Os) {
		return
	}

	s.userHolder.HoldUserConn(uc)
	s.userHolder.StoreUserConn(ctx, uc)
	s.userHolder.StoreUserClients(ctx, uc)
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

		if !uc.IsLogin.Load() {
			s.logger.ConnDebug("heartbeat not login", uc.Desc(), ConnLifecycle, nil)
			s.closeConnFD(c, uc, "heartbeat not login")
			return timewheel.Break
		}

		if uc.IsClosed.Load() {
			s.logger.ConnDebug("heartbeat break,conn is closed", uc.Desc(), ConnLifecycle, nil)
			s.closeConnFD(c, uc, "heartbeat,conn has been closed")
			return timewheel.Break
		}

		if time.Since(uc.LastHeartbeat.Load()) >= s.cfg.Heartbeat.Timeout {
			s.logger.ConnDebug("heartbeat timeout", uc.Desc(), ConnLifecycle, nil, zap.Time("lastHeartbeat", uc.LastHeartbeat.Load()))
			s.closeConnFD(c, uc, "heartbeat timeout")
			return timewheel.Break
		}

		return timewheel.Retry
	}

	if _, err := s.hts.Ticking(fun, s.cfg.Heartbeat.Interval); err != nil {
		return err
	}

	return nil
}

func (s *TcpServer) closeConn(ctx context.Context, c gnet.Conn, uc *domain.UserConn) {
	uc.Close()
	s.userHolder.RemoveUserConn(uc)
	s.userHolder.DeleteUserConn(ctx, uc)
	s.userHolder.DeleteUserClient(ctx, uc)
	s.delContext(c)
}

// 打开链接：删除ctx、删除uc到本地、停止心跳
func (s *TcpServer) closeConnFD(c gnet.Conn, uc *domain.UserConn, reason string) {

	err := c.CloseWithCallback(func(c gnet.Conn, err error) error {
		s.logger.ConnDebug("close completed", uc.Desc(), ConnLifecycle, err, zap.String("reason", reason))
		return nil
	})

	if err != nil {
		s.logger.ConnDebug("close error", uc.Desc(), ConnLifecycle, err)
	}
}

func (s *TcpServer) response(c gnet.Conn, uc *domain.UserConn, packet *api.Packet) {

	bs, err := s.codec.Encode(packet)
	defer bb.Put(bs)

	if err != nil {
		s.logger.PktDebug("encode error", uc.Desc(), packet.GetPacketId(), "", PacketTracking, err)
		s.closeConnFD(c, uc, err.Error())
		return
	}

	if err := c.AsyncWrite(bs.Bytes(), func(c gnet.Conn, err error) error {
		if !packet.IsHeartbeat() {
			s.logger.PktDebug("write completed", uc.Desc(), packet.GetPacketId(), "", PacketTracking, err)
		}
		return err
	}); err != nil {
		s.logger.PktDebug("write error", uc.Desc(), packet.GetPacketId(), "", PacketTracking, err)
		s.closeConnFD(c, uc, err.Error())
		return
	}

}
