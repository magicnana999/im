package ms

import (
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/magicnana999/im/api/kitex_gen/api/brokerservice"
	"github.com/magicnana999/im/global"
	"github.com/magicnana999/im/infrastructure"
	"net"
	"sync"
)

var (
	DefaultMS *BrokerMS
	dmsOnce   sync.Once
)

type BrokerMS struct {
	registry  registry.Registry
	brokerSrv server.Server
}

func InitMicroService() *BrokerMS {
	dmsOnce.Do(func() {

		register := infrastructure.InitEtcdRegistry()

		addr, _ := net.ResolveTCPAddr("tcp", global.GetMicroService().BrokerAddr)
		svr := brokerservice.NewServer(&BrokerServiceImpl{},
			server.WithServiceAddr(addr),
			server.WithRegistry(register),
			server.WithServerBasicInfo(
				&rpcinfo.EndpointBasicInfo{
					ServiceName: global.GetMicroService().BrokerName,
				},
			),
		)

		DefaultMS = &BrokerMS{
			registry:  register,
			brokerSrv: svr,
		}

	})

	return DefaultMS
}
