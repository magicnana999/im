package handler

import (
	"context"
	"github.com/magicnana999/im/conf"
	"github.com/magicnana999/im/kafka"
	"github.com/magicnana999/im/pb"
	"sync"
)

var (
	DefaultMessageReceiver = &MessageReceiver{}
	mrLock                 sync.Mutex
)

type MessageReceiver struct {
	mqProducer *kafka.Producer
}

func InitMessageReceiver() *MessageReceiver {

	mrLock.Lock()
	defer mrLock.Unlock()

	if DefaultMessageReceiver.mqProducer != nil {
		return DefaultMessageReceiver
	}

	DefaultMessageReceiver.mqProducer = kafka.InitProducer([]string{conf.Global.Kafka.String()})
	return DefaultMessageReceiver
}

func (m *MessageReceiver) Receive(ctx context.Context, message *pb.MessageBody) (*pb.Packet, error) {
	if err := m.mqProducer.SendRoute(ctx, message, 1); err != nil {
		return nil, err
	}
	return message.Success(nil).Wrap(), nil
}
