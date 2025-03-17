package broker

import (
	"fmt"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/magicnana999/im/api/kitex_gen/api/brokerservice"
	"github.com/magicnana999/im/global"
	"github.com/magicnana999/im/pkg/logger"
	"net"
	"sync"
)

var (
	DefaultMS *MicroService
	dmsOnce   sync.Once
)

type MicroService struct {
	registry registry.Registry
	server   server.Server
	cli      []client.Client
}

func InitMicroService() {
	dmsOnce.Do(func() {

		r, err := etcd.NewEtcdRegistry(global.GetMicroService().EtcdAddr)
		if err != nil {
			logger.Fatalf("init micro service fail, %v", err)
		}

		hp := fmt.Sprintf("127.0.0.1:507%d", 5+id)
		addr, _ := net.ResolveTCPAddr("tcp", hp)
		svr := brokerservice.NewServer(&BrokerServiceImpl{id: id},
			server.WithServiceAddr(addr),
			server.WithRegistry(registry),
			server.WithServerBasicInfo(
				&rpcinfo.EndpointBasicInfo{
					ServiceName: "im.broker",
				},
			),
		)

		DefaultMS = &MicroService{
			registry: r,
			services: make([]server.Server, 0),
		}
	})
}
