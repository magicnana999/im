package main

import (
	"context"
	"fmt"
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

type BrokerServiceImpl struct {
	id int
}

func (b BrokerServiceImpl) Deliver(ctx context.Context, req *api.DeliverRequest) (r *api.DeliverResponse, err error) {
	log.Printf("deliver request: %d %v", b.id, req)
	return &api.DeliverResponse{Message: "OK"}, nil
}

func StartService(id int) {

	registry, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Fatal(err)
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

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}

func GetRouterServiceClient() (routerservice.Client, error) {
	r, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Fatal(err)
	}
	return routerservice.NewClient("im.router", client.WithResolver(r))
}

func main() {
	go StartService(0)
	go StartService(1)
	time.Sleep(time.Second)
	client, err := GetRouterServiceClient()
	if err != nil {
		log.Fatal(err)
	}

	_, _ = console.ReadLine()
	ret, er := client.Route(context.Background(), &api.RouteRequest{Message: "haha"})
	fmt.Println(ret, er)

	_, _ = console.ReadLine()

}
