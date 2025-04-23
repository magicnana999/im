package infra

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/cloudwego/kitex/pkg/registry"
	jsoniter "github.com/json-iterator/go"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/magicnana999/im/api/kitex_gen/api/brokerservice"
	"github.com/magicnana999/im/api/kitex_gen/api/businessservice"
	"github.com/magicnana999/im/api/kitex_gen/api/routerservice"
	"github.com/magicnana999/im/global"
	"github.com/magicnana999/im/pkg/logger"
	etcdclient "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"sync"
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

type BrokerClientResolver struct {
	Endpoints map[string]brokerservice.Client
	lock      sync.RWMutex
	client    *etcdclient.Client
	logger    *logger.Logger
}

func NewBrokerClientResolver(g *global.Config, lc fx.Lifecycle) (*BrokerClientResolver, error) {
	c := getOrDefaultEtcdConfig(g)

	cli, err := etcdclient.New(etcdclient.Config{
		Endpoints:   c.Endpoints,
		DialTimeout: c.DialTimeout,
	})

	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return cli.Close()
		},
	})

	return &BrokerClientResolver{
		Endpoints: make(map[string]brokerservice.Client),
		client:    cli,
		logger:    logger.Named("kitex"),
	}, nil
}

func (r *BrokerClientResolver) Client(ctx context.Context, addr string) (brokerservice.Client, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	if v, ok := r.Endpoints[addr]; ok && v != nil {
		return v, nil
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	if v, ok := r.Endpoints[addr]; ok && v != nil {
		return v, nil
	}

	if err := r.fetch(ctx); err != nil {
		return nil, err
	}

	if v, ok := r.Endpoints[addr]; ok && v != nil {
		return v, nil
	}

	r.logger.Error("no registered broker service client")
	return nil, errors.New("no registered broker service client")
}

func (r *BrokerClientResolver) fetch(ctx context.Context) error {
	resp, err := r.client.Get(ctx, "kitex/registry-etcd/im.broker", etcdclient.WithPrefix())
	if err != nil {
		return err
	}

	for _, v := range resp.Kvs {
		addr := jsoniter.Get(v.Value, "addr").ToString()
		errMsg := fmt.Sprintf("invalid and ignore json:%s", string(v.Value))
		r.logger.Error(errMsg)
		if addr == "" {
			continue
		}
		cli, err := brokerservice.NewClient("im.broker", client.WithHostPorts(addr))
		if err != nil {
			return err
		}
		r.Endpoints[addr] = cli
	}
	return nil
}

func (r *BrokerClientResolver) Close() {
	r.client.Close()
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
