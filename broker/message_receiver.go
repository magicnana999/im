package broker

import (
	"context"
	"github.com/magicnana999/im/conf"
	"github.com/magicnana999/im/errors"
	"github.com/magicnana999/im/kafka"
	"github.com/magicnana999/im/logger"
	"github.com/magicnana999/im/pb"
	"sync"
)

var (
	defaultMessageReceiver = &messageReceiver{}
	mrLock                 sync.Mutex
)

type messageReceiver struct {
	mqProducer *kafka.Producer
}

func initMessageReceiver() *messageReceiver {

	mrLock.Lock()
	defer mrLock.Unlock()

	if defaultMessageReceiver.mqProducer != nil {
		return defaultMessageReceiver
	}

	defaultMessageReceiver.mqProducer = kafka.InitProducer([]string{conf.Global.Kafka.String()})

	logger.DebugF("messageReceiver init")

	return defaultMessageReceiver
}

func (m *messageReceiver) receive(ctx context.Context, message *pb.MessageBody) (*pb.Packet, error) {

	if err := m.mqProducer.SendRoute(ctx, message, 1); err != nil {
		return nil, errors.MsgMQProduceError.Detail(err)
	}
	return message.Success(nil).Wrap(), nil
}
