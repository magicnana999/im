package main

import (
	"context"
	etcd "github.com/api-contrib/registry-etcd"
	console "github.com/asynkron/goconsole"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"log"
	"main/kitex_gen/api"
	"main/kitex_gen/api/brokerservice"
	"main/kitex_gen/api/routerservice"
	"net"
	"time"
)

type RouterServiceImpl struct{}

func (r2 RouterServiceImpl) Route(ctx context.Context, req *api.RouteRequest) (r *api.RouteResponse, err error) {
	log.Printf("route request: %v", req)
	return &api.RouteResponse{Message: "OK"}, nil
}

func StartService() {

	registry, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Fatal(err)
	}

	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8888")
	svr := routerservice.NewServer(new(RouterServiceImpl),
		server.WithServiceAddr(addr),
		server.WithRegistry(registry),
		server.WithServerBasicInfo(
			&rpcinfo.EndpointBasicInfo{
				ServiceName: "im.router",
			},
		),
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}

func GetBrokeServiceClient1() (brokerservice.Client, error) {
	//r, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//return brokerservice.NewClient("im.broker", client.WithResolver(r))
	return brokerservice.NewClient("im.broker", client.WithHostPorts("127.0.0.1:5075"))
}

func GetBrokeServiceClient2() (brokerservice.Client, error) {
	//r, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//return brokerservice.NewClient("im.broker", client.WithResolver(r))
	return brokerservice.NewClient("im.broker", client.WithHostPorts("127.0.0.1:5076"))
}

func main() {
	go StartService()
	time.Sleep(time.Second)
	client1, _ := GetBrokeServiceClient1()
	client2, _ := GetBrokeServiceClient2()

	_, _ = console.ReadLine()
	client1.Deliver(context.Background(), &api.DeliverRequest{Message: "haha 1111"})
	client2.Deliver(context.Background(), &api.DeliverRequest{Message: "haha 2222"})

	_, _ = console.ReadLine()

}
