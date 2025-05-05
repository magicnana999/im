package infra

import (
	"context"
	"errors"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/cloudwego/kitex/pkg/registry"
	jsoniter "github.com/json-iterator/go"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/magicnana999/im/api/kitex_gen/api"
	"github.com/magicnana999/im/api/kitex_gen/api/brokerservice"
	"github.com/magicnana999/im/api/kitex_gen/api/businessservice"
	"github.com/magicnana999/im/api/kitex_gen/api/routerservice"
	"github.com/magicnana999/im/global"
	"github.com/magicnana999/im/pkg/logger"
	"go.etcd.io/etcd/api/v3/mvccpb"
	etcdclient "go.etcd.io/etcd/client/v3"
	"go.uber.org/atomic"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"strings"
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

var (
	BrokerIsDown = errors.New("broker is down")
)

type LockableBrokerClient struct {
	isShutdown   atomic.Bool
	brokerClient brokerservice.Client
}

func (l *LockableBrokerClient) Deliver(ctx context.Context, req *api.DeliverRequest) (res *api.DeliverReply, err error) {
	if l.isShutdown.Load() {
		return nil, BrokerIsDown
	}
	return l.brokerClient.Deliver(ctx, req)
}

type BrokerClientResolver struct {
	endpoints map[string]*LockableBrokerClient
	client    *etcdclient.Client
	logger    *logger.Logger
	lock      sync.RWMutex
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

	srv := &BrokerClientResolver{
		endpoints: make(map[string]*LockableBrokerClient),
		client:    cli,
		logger:    logger.Named("kitex"),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return srv.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return srv.Stop(ctx)
		},
	})

	return srv, nil
}

func (s *BrokerClientResolver) Deliver(ctx context.Context, req *api.DeliverRequest) (res *api.DeliverReply, err error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.endpoints["127.0.0.1:5075"].Deliver(ctx, req)
}

func (s *BrokerClientResolver) Stop(ctx context.Context) error {
	if err := s.client.Close(); err != nil {
		s.logger.Error("Failed to close etcd client", zap.Error(err))
	}

	s.lock.Lock()
	defer s.lock.Unlock()
	for k, v := range s.endpoints {
		v.isShutdown.Store(true)
		delete(s.endpoints, k)
	}
	return nil
}

func (s *BrokerClientResolver) Start(ctx context.Context) error {
	// 获取现有服务
	resp, err := s.client.Get(ctx, "kitex/registry-etcd/im.broker", etcdclient.WithPrefix())
	if err != nil {
		s.logger.Error("Failed to get initial brokers", zap.Error(err))
		return err
	}
	for _, kv := range resp.Kvs {
		addr := jsoniter.Get(kv.Value, "addr").ToString()
		for i := 0; i < 3; i++ {
			cli, err := brokerservice.NewClient("im.broker", client.WithHostPorts(addr))
			if err == nil {
				s.lock.Lock()
				s.endpoints[addr] = &LockableBrokerClient{brokerClient: cli}
				s.lock.Unlock()
				s.logger.Info("broker client,ok", zap.String("addr", addr), zap.String("event", "INIT"), zap.String("key", string(kv.Key)))
				break
			}
			s.logger.Error("broker client,error", zap.Error(err), zap.Int("attempt", i+1))
			time.Sleep(time.Duration(i*100) * time.Millisecond)
		}
	}

	// 启动 watch
	go func() {
		watchChan := s.client.Watch(ctx, "kitex/registry-etcd/im.broker", etcdclient.WithPrefix())
		for watchResp := range watchChan {
			if err := watchResp.Err(); err != nil {
				s.logger.Error("Watch error", zap.Error(err))
				continue
			}
			for _, event := range watchResp.Events {
				if event.Type == mvccpb.PUT {
					addr := jsoniter.Get(event.Kv.Value, "addr").ToString()
					for i := 0; i < 3; i++ {
						cli, err := brokerservice.NewClient("im.broker", client.WithHostPorts(addr))
						if err == nil {
							s.lock.Lock()
							s.endpoints[addr] = &LockableBrokerClient{brokerClient: cli}
							s.lock.Unlock()
							s.logger.Info("broker client,ok", zap.String("addr", addr), zap.String("event", "PUT"), zap.String("key", string(event.Kv.Key)))
							break
						}
						s.logger.Error("broker client,error", zap.Error(err), zap.Int("attempt", i+1))
						time.Sleep(time.Duration(i*100) * time.Millisecond)
					}
					continue
				}
				if event.Type == mvccpb.DELETE {
					key := string(event.Kv.Key)
					if !strings.HasPrefix(key, "kitex/registry-etcd/im.broker/") {
						s.logger.Warn("Invalid key for DELETE event", zap.String("key", key))
						continue
					}
					addr := key[len("kitex/registry-etcd/im.broker/"):]
					s.lock.Lock()
					c, ok := s.endpoints[addr]
					if ok && c != nil {
						c.isShutdown.Store(true)
						delete(s.endpoints, addr)
						s.logger.Info("broker client,deleted", zap.String("addr", addr), zap.String("event", "DELETE"), zap.String("key", key))
					}
					s.lock.Unlock()
				}
			}
		}
		s.logger.Info("Watch stopped", zap.Error(ctx.Err()))
	}()
	return nil
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
