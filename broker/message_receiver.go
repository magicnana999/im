package broker

import (
	"context"
	"github.com/magicnana999/im/conf"
	"github.com/magicnana999/im/errors"
	"github.com/magicnana999/im/kafka"
	"github.com/magicnana999/im/pb"
	"sync"
)

var (
	defaultMessageReceiver *messageReceiver
	mrOnce                 sync.Once
)

type messageReceiver struct {
	mqProducer *kafka.Producer
}

func initMessageReceiver() *messageReceiver {

	mrOnce.Do(func() {
		defaultMessageReceiver = &messageReceiver{}
		defaultMessageReceiver.mqProducer = kafka.InitProducer([]string{conf.Global.Kafka.String()})

	})

	return defaultMessageReceiver
}

func (m *messageReceiver) receive(ctx context.Context, message *pb.MessageBody) (*pb.Packet, error) {

	if err := m.mqProducer.SendRoute(ctx, message, 1); err != nil {
		return nil, errors.MsgMQProduceError.Detail(err)
	}
	return message.Success(nil).Wrap(), nil
}
