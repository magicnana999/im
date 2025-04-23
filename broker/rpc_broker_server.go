package broker

import (
	"context"
	"fmt"
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/magicnana999/im/api/kitex_gen/api"
	"github.com/magicnana999/im/api/kitex_gen/api/brokerservice"
	"github.com/magicnana999/im/broker/holder"
	"github.com/magicnana999/im/global"
	"go.uber.org/fx"
	"net"
)

type RpcBrokerServer struct {
	cfg        *global.RBSConfig
	registry   registry.Registry
	server     server.Server
	mss        *MessageSendServer
	userHolder *holder.UserHolder
	logger     *Logger
}

func getOrDefaultRBSConfig(g *global.Config) (*global.RBSConfig, error) {
	c := &global.RBSConfig{}
	if g != nil && g.RBS != nil {
		*c = *g.RBS
	}

	if c.Network == "" {
		c.Network = "tcp"
	}

	if c.Addr == "" {
		c.Addr = ":5075"
	}

	return c, nil
}
func NewRpcBrokerServer(
	registry registry.Registry,
	mss *MessageSendServer,
	userHolder *holder.UserHolder,
	g *global.Config,
	lc fx.Lifecycle) (*RpcBrokerServer, error) {

	c, err := getOrDefaultRBSConfig(g)
	if err != nil {
		return nil, err
	}

	logger := NewLogger("rbs", c.DebugMode)

	s := &RpcBrokerServer{
		cfg:        c,
		registry:   registry,
		mss:        mss,
		userHolder: userHolder,
		logger:     logger,
	}

	addr, _ := net.ResolveTCPAddr(c.Network, c.Addr)
	svr := brokerservice.NewServer(s,
		server.WithServiceAddr(addr),
		server.WithRegistry(registry),
		server.WithServerBasicInfo(
			&rpcinfo.EndpointBasicInfo{
				ServiceName: "im.broker",
			},
		),
	)

	s.server = svr

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return s.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return s.Stop(ctx)
		},
	})
	return s, nil
}

func (s *RpcBrokerServer) Start(ctx context.Context) error {
	go func() {
		fmt.Println("---------------", s.server.GetServiceInfos())

		err := s.server.Run()
		s.logger.SrvInfo("rpc server start", SrvLifecycle, err)
		fmt.Println("---------------", s.server.GetServiceInfos())

	}()
	return nil
}

func (s *RpcBrokerServer) Stop(ctx context.Context) error {
	err := s.server.Stop()
	s.logger.SrvInfo("rpc server stop", SrvLifecycle, err)

	return err
}

func (s *RpcBrokerServer) Deliver(ctx context.Context, req *api.DeliverRequest) (res *api.DeliverReply, err error) {
	for _, label := range req.UserLabels {
		uc := s.userHolder.GetUserConn(label)
		s.mss.Send(req.Message, uc)
	}

	return &api.DeliverReply{
		MessageId: req.MessageId,
		Code:      0,
		Message:   "",
	}, nil
}
