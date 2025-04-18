package infra

import (
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/cloudwego/kitex/pkg/registry"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/magicnana999/im/api/kitex_gen/api/businessservice"
	"github.com/magicnana999/im/api/kitex_gen/api/routerservice"
	"github.com/magicnana999/im/global"
	"github.com/magicnana999/im/pkg/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"time"
)

// getOrDefaultEtcdConfig 返回 Etcd 配置，优先使用全局配置，缺失时应用默认值。
// 不会修改输入的 global.Config。
func getOrDefaultEtcdConfig(g *global.Config) *global.EtcdConfig {

	c := &global.EtcdConfig{}
	if g != nil && g.Etcd != nil {
		*c = *g.Etcd
	}

	if c.Endpoints == nil || len(c.Endpoints) == 0 {
		logger.Named("kitex").Warn("etcd endpoints is empty")
		c.Endpoints = []string{"127.0.0.1:2379"}
	}

	if c.DialTimeout == 0 {
		c.DialTimeout = 5 * time.Second
	}
	return c
}

func NewEtcdRegistry(g *global.Config, lc fx.Lifecycle) (registry.Registry, error) {

	log := logger.Named("kitex")

	c := getOrDefaultEtcdConfig(g)

	reg, err := etcd.NewEtcdRegistry(c.Endpoints, etcd.WithDialTimeoutOpt(c.DialTimeout))
	if err != nil {
		log.Error("etcd could not be open",
			zap.Error(err))
		return nil, err
	}

	return reg, nil
}

func NewEtcdResolver(g *global.Config, lc fx.Lifecycle) (discovery.Resolver, error) {

	log := logger.Named("kitex")

	c := getOrDefaultEtcdConfig(g)

	reg, err := etcd.NewEtcdResolver(c.Endpoints, etcd.WithDialTimeoutOpt(c.DialTimeout))

	if err != nil {
		log.Error("etcd could not be open",
			zap.Error(err))
		return nil, err
	}

	return reg, nil
}

func NewBusinessClient(resolver discovery.Resolver, lc fx.Lifecycle) (businessservice.Client, error) {

	log := logger.Named("kitex")

	cli, err := businessservice.NewClient(
		"im.business",
		client.WithResolver(resolver),
		client.WithMuxConnection(2),
		client.WithRPCTimeout(3*time.Second),
	)
	if err != nil {
		log.Error("business client could not be open", zap.Error(err))
		return nil, err
	}

	return cli, nil
}

func NewRouterClient(resolver discovery.Resolver, lc fx.Lifecycle) (routerservice.Client, error) {
	log := logger.Named("kitex")

	cli, err := routerservice.NewClient(
		"im.router",
		client.WithResolver(resolver),
		client.WithMuxConnection(2),
		client.WithRPCTimeout(3*time.Second),
	)
	if err != nil {
		log.Error("router client could not be open", zap.Error(err))
		return nil, err
	}

	return cli, nil
}

// 以后再说
//func loggerMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
//	return func(ctx context.Context, req, resp interface{}) (err error) {
//		ri := rpcinfo.GetRPCInfo(ctx)
//		// get real request
//		klog.Infof("real request: %+v\n", req.(args).GetFirstArgument())
//		// get local cmd_service information
//		klog.Infof("local cmd_service name: %v\n", ri.From().ServiceName())
//		// get remote cmd_service information
//		klog.Infof("remote cmd_service name: %v, remote method: %v\n", ri.To().ServiceName(), ri.To().Method())
//		if err := next(ctx, req, resp); err != nil {
//			return err
//		}
//		// get real response
//		klog.Infof("real response: %+v\n", resp.(result).GetResult())
//		return nil
//	}
//}
//
//func ClientMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
//	return func(ctx context.Context, req, resp interface{}) (err error) {
//		ri := rpcinfo.GetRPCInfo(ctx)
//		// get timeout information
//		klog.Infof("rpc timeout: %v, readwrite timeout: %v\n", ri.Config().RPCTimeout(), ri.Config().ConnectTimeout())
//		if err := next(ctx, req, resp); err != nil {
//			return err
//		}
//		// get server information
//		klog.Infof("server address: %v\n", ri.To().Address())
//		return nil
//	}
//}
//
//func ServerMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
//	return func(ctx context.Context, req, resp interface{}) (err error) {
//		ri := rpcinfo.GetRPCInfo(ctx)
//		// get client information
//		klog.Infof("client address: %v\n", ri.From().Address())
//		if err := next(ctx, req, resp); err != nil {
//			return err
//		}
//		return nil
//	}
//}
