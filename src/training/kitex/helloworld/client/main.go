package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/transport"
	"main/kitex_gen/api"
	"main/kitex_gen/api/hello"
)

func main() {
	c, e := hello.NewClient(
		"hello",
		client.WithHostPorts("0.0.0.0:8080"),
		client.WithTransportProtocol(transport.TTHeader),
	)
	if e != nil {
		panic(e)
	}

	ret, err := c.Echo(context.Background(), &api.Request{Message: "你好"})
	if err != nil {
		panic(err)
	}

	fmt.Println(ret.Message)
}
