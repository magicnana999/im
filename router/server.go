package router

import (
	"context"
	"github.com/magicnana999/im/conf"
	"github.com/magicnana999/im/kafka"
	"github.com/magicnana999/im/pb"
	goPool "github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
	"sync"
)

var DefaultServer *Server
var sOnce sync.Once

type Server struct {
	mqProducer *kafka.Producer
	mqConsumer *kafka.Consumer
	msgRouter  *messageRouter
	executor   *goPool.Pool
}

func Start(ctx context.Context) *Server {

	sOnce.Do(func() {

		s := &Server{}

		s.mqConsumer = kafka.InitConsumer([]string{conf.Global.Kafka.String()}, kafka.Route, s)
		s.mqProducer = kafka.InitProducer([]string{conf.Global.Kafka.String()})
		s.executor = goPool.Default()
		s.msgRouter = initMessageRouter()
		s.mqConsumer.Start(ctx)

		DefaultServer = s

		initLogger()

	})
	return DefaultServer
}

func (s *Server) Consume(ctx context.Context, msg *pb.MQMessage) error {

	if msg.Count > 3 {
		s.mqProducer.SendRouteDLQ(ctx, msg.Message)
	}

	s.executor.Submit(func() {
		if err := s.msgRouter.routeMessage(ctx, msg); err != nil {
			s.mqProducer.SendRoute(ctx, msg.Message, msg.Count+1)
		}
	})

	return nil
}
