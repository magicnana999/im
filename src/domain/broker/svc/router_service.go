package svc

import (
	"github.com/magicnana999/im/api/kitex_gen/api/routerservice"
	"github.com/magicnana999/im/infrastructure"
	"sync"
)

var (
	rsvcOnce         sync.Once
	DefaultRouterSvc *RouterSvc
)

type RouterSvc struct {
	routerCli routerservice.Client
}

func InitRouterSvc() *RouterSvc {
	rsvcOnce.Do(func() {
		DefaultRouterSvc = &RouterSvc{
			routerCli: infrastructure.InitRouterCli(),
		}
	})

	return DefaultRouterSvc
}
