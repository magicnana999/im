package ms

import (
	"context"
	"github.com/magicnana999/im/api/kitex_gen/api"
)

var DefaultBrokerService = &BrokerServiceImpl{}

func InitServiceImpl() *BrokerServiceImpl {
	return DefaultBrokerService
}

type BrokerServiceImpl struct {
	id int
}

func (b BrokerServiceImpl) Deliver(ctx context.Context, req *api.DeliverRequest) (res *api.DeliverReply, err error) {
	panic("implement me")
}
