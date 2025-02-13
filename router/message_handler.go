package router

import (
	"context"
	"github.com/magicnana999/im/conf"
	"github.com/magicnana999/im/kafka"
	"github.com/magicnana999/im/pb"
	"github.com/magicnana999/im/redis"
	"runtime"
)

var DefaultRouteMessageHandler = &RouteMessageHandler{}

func InitRouteMessageHandler(ctx context.Context) *RouteMessageHandler {

	cpu := runtime.NumCPU()
	brokers := []string{conf.Global.Kafka.String()}
	topic := kafka.Route

	c := kafka.InitConsumer(brokers, cpu, topic, DefaultRouteMessageHandler)
	c.Start(ctx)
	DefaultRouteMessageHandler.mqConsumer = c
	DefaultRouteMessageHandler.userStorage = redis.InitUserStorage()

	return DefaultRouteMessageHandler
}

type RouteMessageHandler struct {
	mqConsumer  *kafka.Consumer
	userStorage *redis.UserStorage
}

func (m *RouteMessageHandler) Consume(ctx context.Context, msg *pb.MQMessage) error {
	mb := msg.GetMessage()

	if !mb.IsToGroup() {
		m, e := m.userStorage.LoadUserConn(ctx, mb.AppId, mb.To)
		if e != nil {
			return e
		}

		for k, v := range m {

		}
	}

	return nil
}
