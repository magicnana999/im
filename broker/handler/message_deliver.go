package handler

import (
	"context"
	"github.com/magicnana999/im/broker/core"
	"github.com/magicnana999/im/broker/state"
	"github.com/magicnana999/im/conf"
	"github.com/magicnana999/im/kafka"
	"github.com/magicnana999/im/pb"
	"sync"
)

var (
	DefaultMessageDeliver = &MessageDeliver{}
	mdLock                sync.Mutex
)

type MessageDeliver struct {
	outsideDeliver *kafka.Consumer
	userState      *state.UserState
	codec          core.Codec
}

func InitMessageDeliver() *MessageDeliver {
	mdLock.Lock()
	defer mdLock.Unlock()

	if DefaultMessageDeliver.outsideDeliver != nil {
		return DefaultMessageDeliver
	}

	addr := []string{conf.Global.Kafka.String()}
	topic := kafka.Deliver
	DefaultMessageDeliver.outsideDeliver = kafka.InitConsumer(addr, topic, DefaultMessageDeliver)
	DefaultMessageDeliver.userState = state.InitUserState()
	DefaultMessageDeliver.codec = core.InitCodec()
	return DefaultMessageDeliver
}

func (m *MessageDeliver) deliver(ctx context.Context, message *pb.MessageBody) error {
	ucs := m.userState.LoadLocalUser(message.AppId, message.UserId)
	for _, uc := range ucs {
		buffer, e := m.codec.Encode(uc.C, message.Wrap())
		if e != nil {
			return e
		}
		uc.C.Write(buffer.Bytes())
	}
	return nil
}

func (m *MessageDeliver) Consume(ctx context.Context, msg *pb.MQMessage) error {
	message := msg.GetMessage()
	return m.deliver(ctx, message)
}
