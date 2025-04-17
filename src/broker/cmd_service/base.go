package cmd_service

import (
	"github.com/magicnana999/im/api/kitex_gen/api/businessservice"
	"github.com/magicnana999/im/infra"
	"github.com/magicnana999/im/pkg/singleton"
)

var baseSingleton = singleton.NewSingleton[*base]

type base struct {
	businessClient businessservice.Client
}

func newBase() *base {
	f := func() *base {
		return &base{
			businessClient: infra.NewBusinessClient(),
		}
	}
	return baseSingleton().Get(f)
}
