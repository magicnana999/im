package router

import (
	"context"
	"github.com/magicnana999/im/api/kitex_gen/api"
	"github.com/magicnana999/im/errors"
	"github.com/magicnana999/im/infra"
	"github.com/magicnana999/im/router/vo"
)

type DeliveryService struct {
	us  *UserService
	bcr *infra.BrokerClientResolver
}

func NewDeliveryService(us *UserService, bcr *infra.BrokerClientResolver) *DeliveryService {
	return &DeliveryService{us: us, bcr: bcr}
}

func (s *DeliveryService) deliverToUser(ctx context.Context, m *api.Message) ([]vo.DeliverFail, error) {

	ret := make([]vo.DeliverFail, 0)

	ucs, err := s.us.GetUserClients(ctx, m.AppId, m.To)
	if err != nil {
		return []vo.DeliverFail{{M: m, UserId: m.To}}, errors.RouteErr.SetDetail(err.Error())
	}

	eachBrokerMsg := make(map[string]*api.DeliverRequest)

	for _, v := range ucs {
		request := eachBrokerMsg[v.BrokerAddr]
		if request == nil {
			request = &api.DeliverRequest{
				MessageId:  m.MessageId,
				Message:    m,
				UserLabels: make([]string, 0),
			}
			eachBrokerMsg[v.BrokerAddr] = request
		}

		request.UserLabels = append(request.UserLabels, v.Label)
	}

	if len(eachBrokerMsg) == 0 {
		return []vo.DeliverFail{{M: m, UserId: m.To}}, errors.RouteErr.SetDetail("no user clients online")

	}

	for k, v := range eachBrokerMsg {
		cli, err := s.bcr.Client(ctx, k)
		if err != nil {
			ret = append(ret, vo.DeliverFail{M: m, UserId: m.To, Label: v.UserLabels})
			continue
		}

		rep, err := cli.Deliver(ctx, v)
		if err != nil {
			ret = append(ret, vo.DeliverFail{M: m, UserId: m.To, Label: v.UserLabels})
			continue
		}

		if rep != nil && rep.Code != 0 {
			ret = append(ret, vo.DeliverFail{M: m, UserId: m.To, Label: v.UserLabels})
			continue
		}
	}

	if len(ret) != 0 {
		return ret, errors.RouteErr.SetDetail("some connection delivery fail")
	}

	return nil, nil
}
