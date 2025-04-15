package broker

import (
	"context"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/magicnana999/im/api/kitex_gen/api"
	"github.com/magicnana999/im/api/kitex_gen/api/brokerservice"
	"github.com/magicnana999/im/infra"
	"github.com/magicnana999/im/pkg/logger"
	"go.uber.org/zap"
	"net"
	"sync"
)

var (
	rpcInstance *RpcBrokerServer
	rpcOnce     sync.Once
)

type RpcBrokerServer struct {
	svr server.Server
}

func NewRpcBrokerServer() *RpcBrokerServer {

	rpcOnce.Do(func() {
		registry := infra.NewEtcdRegistry()
		rpcServer := &RpcBrokerServer{}

		addr, _ := net.ResolveTCPAddr("tcp", "0.0.0.0:5075")
		svr := brokerservice.NewServer(rpcServer,
			server.WithServiceAddr(addr),
			server.WithRegistry(registry),
			server.WithServerBasicInfo(
				&rpcinfo.EndpointBasicInfo{
					ServiceName: "im.broker",
				},
			),
		)

		rpcServer.svr = svr
		rpcInstance = rpcServer
	})

	return rpcInstance
}

func (s *RpcBrokerServer) start(ctx context.Context) error {
	go func() {
		if err := s.svr.Run(); err != nil {
			logger.Fatal("Failed to start rpc broker server", zap.Error(err))
		}
	}()
	return nil
}

func (s *RpcBrokerServer) Deliver(ctx context.Context, req *api.DeliverRequest) (res *api.DeliverReply, err error) {
	//TODO implement me
	panic("implement me")
}
