package router

import (
	"context"
	"fmt"
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/magicnana999/im/api/kitex_gen/api"
	"github.com/magicnana999/im/api/kitex_gen/api/routerservice"
	"github.com/magicnana999/im/errors"
	"github.com/magicnana999/im/global"
	"github.com/magicnana999/im/infra"
	"go.uber.org/fx"
	"net"
)

type RpcRouterServer struct {
	cfg      *global.RRSConfig
	registry registry.Registry
	server   server.Server
	ds       *DeliveryService
}

func getOrDefaultRBSConfig(g *global.Config) (*global.RRSConfig, error) {
	c := &global.RRSConfig{}
	if g != nil && g.RRS != nil {
		*c = *g.RRS
	}

	if c.Network == "" {
		c.Network = "tcp"
	}

	if c.Addr == "" {
		c.Addr = ":5075"
	}

	return c, nil
}
func NewRpcRouterServer(
	registry registry.Registry,
	g *global.Config,
	us *UserService,
	bcr *infra.BrokerClientResolver,
	lc fx.Lifecycle) (*RpcRouterServer, error) {

	c, err := getOrDefaultRBSConfig(g)
	if err != nil {
		return nil, err
	}

	s := &RpcRouterServer{
		cfg:      c,
		registry: registry,
		ds:       NewDeliveryService(us, bcr),
	}

	addr, _ := net.ResolveTCPAddr(c.Network, c.Addr)
	svr := routerservice.NewServer(s,
		server.WithServiceAddr(addr),
		server.WithRegistry(registry),
		server.WithServerBasicInfo(
			&rpcinfo.EndpointBasicInfo{
				ServiceName: "im.broker",
			},
		),
	)

	s.server = svr

	//lc.Append(fx.Hook{
	//	OnStart: func(ctx context.Context) error {
	//		return s.Start(ctx)
	//	},
	//	OnStop: func(ctx context.Context) error {
	//		return s.Stop(ctx)
	//	},
	//})
	return s, nil
}

//func (s *RpcRouterServer) Start(ctx context.Context) error {
//	go func() {
//		err := s.server.Run()
//		s.logger.SrvInfo("rpc server start", SrvLifecycle, err)
//
//	}()
//	return nil
//}
//
//func (s *RpcRouterServer) Stop(ctx context.Context) error {
//	err := s.server.Stop()
//	s.logger.SrvInfo("rpc server stop", SrvLifecycle, err)
//
//	return err
//}

func (s *RpcRouterServer) Route(ctx context.Context, m *api.Message) (res *api.RouteReply, err error) {

	//TODO ..route
	//TODO ..save
	//TODO ..update conversation

	if m.GetGroupId() == 0 {
		//TODO ... 群组消息
	} else {
		err := m.Validate()
		if err != nil {
			return nil, errors.RouteErr.SetDetail(err.Error())
		}

		fail, err := s.ds.deliverToUser(ctx, m)
		if err != nil {
			//TODO... setToOffline
			fmt.Println(fail)
		}
	}

	return nil, nil
}
