package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/kitex/server"
	"log"
	"main/kitex_gen/api"
	"main/kitex_gen/api/hello"
	"net"
)

type HelloImpl struct{}

func (h HelloImpl) Echo(ctx context.Context, req *api.Request) (r *api.Response, err error) {
	fmt.Println("Hello in server", req)
	return &api.Response{Message: "你好啊"}, nil
}

func main() {

	addr, _ := net.ResolveTCPAddr("tcp", "0.0.0.0:8080")
	svr := hello.NewServer(new(HelloImpl), server.WithServiceAddr(addr))

	err := svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
