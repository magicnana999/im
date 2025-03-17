package infrastructure

import (
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/cloudwego/kitex/pkg/registry"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/magicnana999/im/api/kitex_gen/api/routerservice"
	"github.com/magicnana999/im/api/kitex_gen/api/serverservice"
	"github.com/magicnana999/im/global"
	"github.com/magicnana999/im/pkg/logger"
	"sync"
)

var (
	registryOnce sync.Once
	resolverOnce sync.Once

	Register registry.Registry
	Resolver discovery.Resolver

	serverCliOnce sync.Once
	routerCliOnce sync.Once

	ServerCli serverservice.Client
	RouterCli routerservice.Client
)

func InitEtcdRegistry() registry.Registry {
	registryOnce.Do(func() {
		register, err := etcd.NewEtcdRegistry(global.GetMicroService().EtcdAddr)
		if err != nil {
			logger.Fatalf("init micro service fail, %v", err)
		}
		Register = register
	})
	return Register
}

func InitEtcdResolver() discovery.Resolver {
	resolverOnce.Do(func() {
		resolver, err := etcd.NewEtcdResolver(global.GetMicroService().EtcdAddr)
		if err != nil {
			logger.Fatalf("init micro service fail, %v", err)
		}

		Resolver = resolver
	})
	return Resolver
}

func InitServerCli() serverservice.Client {
	resolver := InitEtcdResolver()

	serverCliOnce.Do(func() {
		cli, err := serverservice.NewClient(
			global.GetMicroService().ServerName,
			client.WithResolver(resolver))

		if err != nil {
			logger.Fatalf("init server client failed, %v", err)
		}

		ServerCli = cli
	})

	return ServerCli
}

func InitRouterCli() routerservice.Client {
	resolver := InitEtcdResolver()

	routerCliOnce.Do(func() {
		cli, err := routerservice.NewClient(
			global.GetMicroService().ServerName,
			client.WithResolver(resolver))

		if err != nil {
			logger.Fatalf("init router client failed, %v", err)
		}

		RouterCli = cli
	})

	return RouterCli
}
