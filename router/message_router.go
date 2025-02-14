package router

import (
	"context"
	"github.com/goccy/go-json"
	"github.com/magicnana999/im/conf"
	"github.com/magicnana999/im/domain"
	"github.com/magicnana999/im/kafka"
	"github.com/magicnana999/im/pb"
	"github.com/magicnana999/im/redis"
	"github.com/magicnana999/im/svc"
	"sync"
)

var defaultMessageRouter *messageRouter
var mrLock sync.Mutex

type messageRouter struct {
	groupMemberSvc *svc.GroupMemberSvc
	userStorage    *redis.UserStorage
	mqProducer     *kafka.Producer
}

func initMessageRouter() *messageRouter {
	mrLock.Lock()
	defer mrLock.Unlock()
	if defaultMessageRouter != nil {
		return defaultMessageRouter
	}
	defaultMessageRouter = &messageRouter{
		groupMemberSvc: svc.InitGroupMemberSvc(),
		userStorage:    redis.InitUserStorage(),
		mqProducer:     kafka.InitProducer([]string{conf.Global.Kafka.String()}),
	}

	return defaultMessageRouter
}

func (s *messageRouter) routeMessage(ctx context.Context, delivery *pb.MQMessage) error {
	message := delivery.Message
	if message == nil {
		return nil
	}

	var userIds []int64
	var err error

	if message.IsToGroup() {
		userIds, err = s.groupMemberSvc.LoadAndFetch(ctx, message.AppId, message.GroupId)
		if err != nil {
			return err
		}
	} else {
		userIds = append(userIds, message.UserId)
	}

	online := make(map[string]*domain.UserConnection)
	offline := make([]int64, 0)
	for _, userId := range userIds {
		m, e := s.userStorage.LoadUserConn(ctx, message.AppId, userId)
		if e != nil {
			return e
		}

		if m == nil || len(m) == 0 {
			offline = append(offline, userId)
			continue
		}

		for k, v := range m {
			var uc *domain.UserConnection
			ee := json.Unmarshal([]byte(v), uc)
			if ee != nil {
				return e
			}
			online[k] = uc
		}
	}

	brokerMap := make(map[string][]string)
	for _, uc := range online {
		if brokerMap[uc.BrokerAddr] == nil {
			brokerMap[uc.BrokerAddr] = make([]string, 0)
		}

		brokerMap[uc.BrokerAddr] = append(brokerMap[uc.BrokerAddr], uc.Label())
	}

	for key, v := range brokerMap {
		s.mqProducer.SendDeliver(ctx, key, message, v)
	}
	s.mqProducer.SendStore(ctx, message)
	s.mqProducer.SendOffline(ctx, message, offline)

	return nil

}
