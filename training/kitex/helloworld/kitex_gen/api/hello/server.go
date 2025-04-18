// Code generated by Kitex v0.12.3. DO NOT EDIT.
package hello

import (
	server "github.com/cloudwego/kitex/server"
	api "main/kitex_gen/api"
)

// NewServer creates a server.Server with the given handler and options.
func NewServer(handler api.Hello, opts ...server.Option) server.Server {
	var options []server.Option

	options = append(options, opts...)
	options = append(options, server.WithCompatibleMiddlewareForUnary())

	svr := server.NewServer(options...)
	if err := svr.RegisterService(serviceInfo(), handler); err != nil {
		panic(err)
	}
	return svr
}

func RegisterService(svr server.Server, handler api.Hello, opts ...server.RegisterOption) error {
	return svr.RegisterService(serviceInfo(), handler, opts...)
}
