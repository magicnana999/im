package infra

import (
	"context"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/cloudwego/kitex/pkg/registry"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/magicnana999/im/api/kitex_gen/api/businessservice"
	"github.com/magicnana999/im/api/kitex_gen/api/routerservice"
	"github.com/magicnana999/im/global"
	"github.com/magicnana999/im/pkg/logger"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"time"
)

const (
	Etcd           = "etcd"
	BusinessClient = "businessclient"
	RouterClient   = "routerclient"
	BrokerClient   = "brokerclient"
)

type EtcdConfig struct {
	clientv3.Config
}

func NewEtcdRegistry(lc fx.Lifecycle) registry.Registry {

	c := global.GetEtcd()

	if c == nil {
		logger.Fatal("etcd configuration not found",
			zap.String(logger.SCOPE, Etcd),
			zap.String(logger.OP, Init))
	}

	reg, err := etcd.NewEtcdRegistry(c.Endpoints, etcd.WithDialTimeoutOpt(c.DialTimeout))
	if err != nil {
		logger.Fatal("etcd could not be open",
			zap.String(logger.SCOPE, Etcd),
			zap.String(logger.OP, Init),
			zap.Error(err))
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("etcd connection established",
				zap.String(logger.SCOPE, Etcd),
				zap.String(logger.OP, Init))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if e := reg.Deregister(nil); e != nil {
				logger.Error("etcd could not close",
					zap.String(logger.SCOPE, Etcd),
					zap.String(logger.OP, Close),
					zap.Error(e))
				return e
			} else {
				logger.Info("etcd closed",
					zap.String(logger.SCOPE, Etcd),
					zap.String(logger.OP, Close))
				return nil
			}
		},
	})

	return reg
}

func NewEtcdResolver(lc fx.Lifecycle) discovery.Resolver {

	c := global.GetEtcd()

	if c == nil {
		logger.Fatal("etcd configuration not found",
			zap.String(logger.SCOPE, Etcd),
			zap.String(logger.OP, Init))
	}

	reg, err := etcd.NewEtcdResolver(c.Endpoints, etcd.WithDialTimeoutOpt(c.DialTimeout))
	if err != nil {
		logger.Fatal("etcd could not be open",
			zap.String(logger.SCOPE, Etcd),
			zap.String(logger.OP, Init),
			zap.Error(err))
	}

	return reg
}

func NewBusinessCli(resolver discovery.Resolver, lc fx.Lifecycle) businessservice.Client {
	cli, err := businessservice.NewClient(
		"im.business",
		client.WithResolver(resolver),
		client.WithMuxConnection(2),
		client.WithRPCTimeout(3*time.Second),
	)
	if err != nil {
		logger.Fatal("business client could not be open",
			zap.String(logger.SCOPE, BusinessClient),
			zap.String(logger.OP, Init),
			zap.Error(err))
	}

	return cli
}

func NewRouterCli(resolver discovery.Resolver, lc fx.Lifecycle) routerservice.Client {
	cli, err := routerservice.NewClient(
		"im.router",
		client.WithResolver(resolver),
		client.WithMuxConnection(2),
		client.WithRPCTimeout(3*time.Second),
	)
	if err != nil {
		logger.Fatal("router client could not be open",
			zap.String(logger.SCOPE, RouterClient),
			zap.String(logger.OP, Init),
			zap.Error(err))
	}

	return cli
}
