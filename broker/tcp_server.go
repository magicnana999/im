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

const (
	// DefaultTCPAddress 定义 TCP 服务器的默认监听地址和端口。
	// 默认值为 "0.0.0.0:5075"，表示监听所有网络接口的 5075 端口。
	// 可通过 global.TCPConfig.Addr 覆盖，生产环境建议指定具体 IP 以提高安全性。
	DefaultTCPAddress = "0.0.0.0:5075"

	// DefaultTickInterval 定义服务器定时任务（OnTick）的默认执行间隔。
	// 默认值为 1 秒，用于定期检查服务器状态或执行清理任务。
	// 可通过 global.TCPConfig.Interval 调整，需根据实际需求平衡性能和频率。
	DefaultTickInterval = time.Second

	// DefWorkerSizeOfEachCPU 定义服务器Worker的数量，cpu的倍数。
	// 默认值为 cpu的个数 * 1000
	// 可通过 global.TCPConfig.DefWorkerSizeOfEachCPU 调整。
	DefWorkerSizeOfEachCPU = 1000

	// DefWorkerExpireDuration 定义协程池（ants.Pool）中空闲协程的默认过期时间。
	// 默认值为 1 分钟，过期后空闲协程将被回收以释放资源。
	// 可通过 global.TcpWorkerConfig.ExpireDuration 调整，需根据负载情况优化。
	DefWorkerExpireDuration = time.Minute

	// DefWorkerMaxTask 定义协程池的最大任务队列长度。
	// 默认值为 100,000，表示最多可排队的任务数，防止任务堆积。
	// 可通过 global.TcpWorkerConfig.MaxBlockingTasks 调整，需通过压力测试确定。
	DefWorkerMaxTask = 100_000

	// DefWorkerReleaseWaiting 定义协程池关闭时的默认等待时间。
	// 默认值为 500 毫秒，确保关闭时所有任务有足够时间完成。
	// 可通过配置调整，需平衡关闭速度和任务完成需求。
	DefWorkerReleaseWaiting = time.Millisecond * 500

	// DefMulticore 控制 gnet 是否启用多核处理。
	// 默认值为 true，利用多核 CPU 提高并发性能。
	// 通常无需修改，除非在特定场景下需要单核调试。
	DefMulticore = true

	// DefLockOSThread 控制 gnet 是否将事件循环绑定到特定 OS 线程。
	// 默认值为 true，减少上下文切换，提高性能。
	// 通常无需修改，除非在特定环境中需要释放线程调度。
	DefLockOSThread = true

	// DefReadBufferCap 定义 gnet 连接的默认读缓冲区大小。
	// 默认值为 4096 字节，平衡内存使用和读取效率。
	// 可根据消息大小和并发量调整，需通过性能测试优化。
	DefReadBufferCap = 4096

	// DefWriteBufferCap 定义 gnet 连接的默认写缓冲区大小。
	// 默认值为 4096 字节，平衡内存使用和写入效率。
	// 可根据消息大小和并发量调整，需通过性能测试优化。
	DefWriteBufferCap = 4096

	// DefLoadBalancing 定义 gnet 的事件循环负载均衡策略。
	// 默认值为 gnet.RoundRobin，表示轮询分配连接到事件循环。
	// 通常无需修改，除非需要其他策略（如最小连接数）。
	DefLoadBalancing = gnet.RoundRobin

	// DefReuseAddr 控制 gnet 是否启用 SO_REUSEADDR 选项。
	// 默认值为 true，允许快速重启服务器，复用同一端口。
	// 生产环境中建议保持启用，避免端口占用问题。
	DefReuseAddr = true

	// DefReusePort 控制 gnet 是否启用 SO_REUSEPORT 选项。
	// 默认值为 true，允许多个进程绑定同一端口，提高并发性能。
	// 生产环境中建议启用，尤其在多进程或容器化部署时。
	DefReusePort = true

	// DefTcpKeepAlive 定义 TCP 连接的默认保活时间。
	// 默认值为 1 分钟，定期发送保活探针以检测连接状态。
	// 可根据网络环境调整，需平衡保活频率和资源开销。
	DefTcpKeepAlive = time.Minute

	// DefTcpNoDelay 控制 TCP 是否启用 Nagle 算法。
	// 默认值为 gnet.TCPNoDelay，表示禁用 Nagle 算法，减少延迟。
	// 适合实时性要求高的 IM 系统，通常无需修改。
	DefTcpNoDelay = gnet.TCPNoDelay

	// DefSocketRecvBuffer 定义 TCP 套接字的默认接收缓冲区大小。
	// 默认值为 8192 字节，影响接收数据的吞吐量。
	// 可根据网络带宽和消息大小调整，需通过性能测试优化。
	DefSocketRecvBuffer = 8192

	// DefSocketSendBuffer 定义 TCP 套接字的默认发送缓冲区大小。
	// 默认值为 8192 字节，影响发送数据的吞吐量。
	// 可根据网络带宽和消息大小调整，需通过性能测试优化。
	DefSocketSendBuffer = 8192

	// DefTicker 控制 gnet 是否启用定时器。
	// 默认值为 true，启用 OnTick 定时任务，用于心跳检测等。
	// 通常无需修改，除非不需要定时任务。
	DefTicker = true

	// DefLogLevel 定义 gnet 的默认日志级别。
	// 默认值为 logging.DebugLevel，记录详细日志，适合开发调试。
	// 生产环境中建议调整为 Info 或更高级别，减少日志开销。
	DefLogLevel = logging.DebugLevel

	// DefEdgeTriggeredIO 控制 gnet 是否启用边沿触发 I/O。
	// 默认值为 true，优化高并发场景下的 I/O 性能。
	// 通常无需修改，除非在特定场景下需要水平触发。
	DefEdgeTriggeredIO = true

	// DefEdgeTriggeredIOChunk 定义边沿触发 I/O 的数据块大小。
	// 默认值为 0，表示由 gnet 自动管理。
	// 通常无需修改，除非需要特定优化。
	DefEdgeTriggeredIOChunk = 0
)

type TcpServer struct {
	*gnet.BuiltinEventEngine
	eng            gnet.Engine
	cfg            *global.TCPConfig
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
		c.Addr = DefaultTCPAddress
	}

	if c.Interval <= 0 {
		c.Interval = DefaultTickInterval
	}

	if c.Worker.Size <= 0 {
		c.Worker.Size = runtime.NumCPU() * DefWorkerSizeOfEachCPU
	}

	if c.Worker.ExpireDuration <= 0 {
		c.Worker.ExpireDuration = DefWorkerExpireDuration
	}

	if c.Worker.MaxBlockingTasks <= 0 {
		c.Worker.MaxBlockingTasks = DefWorkerMaxTask
	}

	return c
}

func newWorkerPool(c *global.TcpWorkerConfig) (*ants.Pool, error) {

	logger := NewLogger("tcp")

	panicHandler := func(panicErr interface{}) {
		logger.Error("tcp worker panic",
			zap.Any("error", panicErr),
			zap.Stack("stack"),
		)
	}

	return ants.NewPool(c.Size,
		ants.WithExpiryDuration(c.ExpireDuration),
		ants.WithMaxBlockingTasks(c.MaxBlockingTasks),
		ants.WithNonblocking(c.Nonblocking),
		ants.WithPanicHandler(panicHandler),
		ants.WithLogger(logger),
		ants.WithPreAlloc(c.PreAlloc),
		ants.WithDisablePurge(c.DisablePurge),
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

	logger := NewLogger("tcp")

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

// Stop 停止
func (s *TcpServer) Stop(ctx context.Context) error {

	s.userHolder.RangeAllUserConn(func(conn *domain.UserConn) bool {
		s.closeConnFD(conn.Conn, conn, "tcp server shutting down")
		return true
	})

	if !s.worker.IsClosed() {
		s.worker.ReleaseTimeout(DefWorkerReleaseWaiting)
	}
	return s.eng.Stop(ctx)
}

// Start 启动
func (s *TcpServer) Start(ctx context.Context) error {
	s.ctx = ctx
	go func() {
		err := gnet.Run(s,
			fmt.Sprintf("tcp://%s", s.cfg.Addr),
			gnet.WithMulticore(DefMulticore),
			gnet.WithLockOSThread(DefLockOSThread),
			gnet.WithReadBufferCap(DefReadBufferCap),
			gnet.WithWriteBufferCap(DefWriteBufferCap),
			gnet.WithLoadBalancing(DefLoadBalancing),
			gnet.WithNumEventLoop(runtime.NumCPU()),
			gnet.WithReuseAddr(DefReuseAddr),
			gnet.WithReusePort(DefReusePort),
			gnet.WithTCPKeepAlive(DefTcpKeepAlive),
			gnet.WithTCPNoDelay(DefTcpNoDelay),
			gnet.WithSocketRecvBuffer(DefSocketRecvBuffer),
			gnet.WithSocketSendBuffer(DefSocketSendBuffer),
			gnet.WithTicker(DefTicker),
			gnet.WithLogger(s.logger),
			gnet.WithLogLevel(DefLogLevel),
			gnet.WithEdgeTriggeredIO(DefEdgeTriggeredIO),
			gnet.WithEdgeTriggeredIOChunk(DefEdgeTriggeredIOChunk))

		s.logger.SrvInfo("tcp starting", SrvLifecycle, err)

		if err != nil {
			s.logger.Fatal("tcp start failed", zap.Error(err))
		}
	}()
	return nil
}

// OnBoot 启动回调
func (s *TcpServer) OnBoot(eng gnet.Engine) (action gnet.Action) {
	s.eng = eng
	s.logger.SrvInfo("tcp started", SrvLifecycle, nil)

	return gnet.None
}

// OnShutdown 停止回调
func (s *TcpServer) OnShutdown(eng gnet.Engine) {
	s.logger.SrvInfo("tcp shutdown", SrvLifecycle, nil)
}

// OnOpen 新链接后回调
func (s *TcpServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	uc := domain.NewUserConn(c)
	err := s.openConn(c, uc)
	s.logger.ConnDebug("connect", uc.Desc(), ConnLifecycle, err, zap.String("uc", string(jsonext.MarshalNoErr(uc))))
	if err != nil {
		return nil, gnet.Close
	}

	return nil, gnet.None
}

// OnClose 关闭时回调
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

// OnTraffic 收到消息
func (s *TcpServer) OnTraffic(c gnet.Conn) (action gnet.Action) {

	ctx := s.getContext(c)
	uc, err := brokerctx.GetCurUserConn(ctx)
	if err != nil {
		s.logger.ConnDebug("read", "", ConnLifecycle, err)
		return gnet.Close
	}

	s.RefreshUser(ctx, uc)

	packets, err := s.codec.Decode(c)

	if err != nil {
		s.logger.ConnDebug("decode", uc.Desc(), ConnLifecycle, err)
		return gnet.Close
	}

	err = s.worker.Submit(func() {
		for _, packet := range packets {
			resp := s.processPacket(ctx, c, uc, packet)
			err := s.response(resp, uc)
			if err != nil {
				s.closeConnFD(c, uc, err.Error())
				return
			}
		}
	})

	if err != nil {
		s.logger.ConnDebug("submit decode failed", uc.Desc(), ConnLifecycle, err)
		return gnet.Close
	}
	return gnet.None
}

// processPacket 处理客户端发来的Packet，heartbeat；command；message
func (s *TcpServer) processPacket(ctx context.Context, c gnet.Conn, uc *domain.UserConn, packet *api.Packet) *api.Packet {

	//心跳
	if packet.IsHeartbeat() {
		return api.HeartbeatACK
	}

	s.logger.PktDebug("read", uc.Desc(), packet.GetPacketId(), packet, PacketTracking, nil)

	//command
	if packet.IsCommand() {
		return s.processCommand(ctx, c, uc, packet)
	}

	//message
	if packet.IsMessage() {
		if packet.GetMessage().IsRequest() {
			return s.processMessage(ctx, c, uc, packet)
		} else {
			s.mrs.Ack(packet.GetMessage().MessageId)
			s.logger.PktDebug("ack process", uc.Desc(), packet.GetPacketId(), packet, PacketTracking, nil)
			return nil
		}
	}

	return nil
}

// 处理message
func (s *TcpServer) processMessage(ctx context.Context, c gnet.Conn, uc *domain.UserConn, packet *api.Packet) *api.Packet {
	message := packet.GetMessage()
	if message.IsRequest() {
		ret, err := s.messageHandler.HandlePacket(ctx, packet)
		s.logger.PktDebug("message process", uc.Desc(), packet.GetPacketId(), packet, PacketTracking, err)
		return ret
	}

	return nil
}

// 处理command
func (s *TcpServer) processCommand(ctx context.Context, c gnet.Conn, uc *domain.UserConn, packet *api.Packet) *api.Packet {
	ret, err := s.commandHandler.HandlePacket(ctx, packet)
	s.logger.PktDebug("command process", uc.Desc(), packet.GetPacketId(), packet, PacketTracking, err)
	if err == nil && packet.GetCommand().CommandType == api.CommandTypeUserLogin {
		s.OnUserLogin(ctx, uc, packet.GetCommand().GetLoginRequest(), ret.GetCommand().GetLoginReply())
	}
	return ret
}

// OnTick gnet ticker
func (s *TcpServer) OnTick() (delay time.Duration, action gnet.Action) {
	return s.cfg.Interval, gnet.None
}

// RefreshUser 刷新用户状态，如果用户已登陆，说明缓存里有值，也需要刷新
func (s *TcpServer) RefreshUser(ctx context.Context, uc *domain.UserConn) {

	if uc.IsClosed.Load() {
		return
	}
	uc.Refresh(time.Now())

	if uc.IsLogin.Load() {
		s.userHolder.RefreshUserConn(ctx, uc)
	}
}

// OnUserLogin 登录成功后处理本地map和redis
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

		if time.Since(uc.LastHeartbeat.Load()) > s.cfg.Heartbeat.Timeout {
			s.logger.ConnDebug("heartbeat timeout", uc.Desc(), ConnLifecycle, nil, zap.Time("lastHeartbeat", uc.LastHeartbeat.Load()))
			s.closeConnFD(c, uc, "heartbeat timeout")
			return timewheel.Break
		}

		return timewheel.Retry
	}

	if _, _, err := s.hts.Ticking(fun); err != nil {
		return err
	}

	return nil
}

// 断开连接时，清理本地map和redis
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

// 回复消息给客户端
func (s *TcpServer) response(packet *api.Packet, uc *domain.UserConn) error {

	if packet == nil {
		return nil
	}

	buffer, err := s.codec.Encode(packet)
	if err != nil {
		return err
	}
	defer bb.Put(buffer)

	if packet.IsHeartbeat() {
		_, err := uc.Conn.Write(buffer.Bytes())
		return err
	}

	err = uc.Conn.AsyncWrite(buffer.Bytes(), func(c gnet.Conn, err error) error {
		s.logger.PktDebug("write completed", uc.Desc(), packet.GetPacketId(), nil, PacketTracking, err)
		return nil
	})

	if err != nil {
		s.logger.PktDebug("write error", uc.Desc(), packet.GetPacketId(), nil, PacketTracking, err)
		return err
	}
	return err
}
