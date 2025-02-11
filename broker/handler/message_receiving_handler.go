package handler

import (
	"context"
	"github.com/magicnana999/im/kafka"
	"github.com/magicnana999/im/pb"
)

var (
	DefaultMessageReceivingHandler = &MessageReceivingHandler{}
)

type MessageReceivingHandler struct {
	mqProducer *kafka.Producer
}

func (m *MessageReceivingHandler) HandlePacket(ctx context.Context, packet *pb.Packet) (*pb.Packet, error) {
	return m.ReceivePacket(ctx, packet)
}

func (m *MessageReceivingHandler) ReceivePacket(ctx context.Context, packet *pb.Packet) (*pb.Packet, error) {
	if err := m.mqProducer.SendMessageRoute(ctx, packet); err != nil {
		return nil, err
	}

	if packet.GetMessageBody().NeedAck == pb.YES {
		return packet, nil
	}

	return nil, nil
}

func (m *MessageReceivingHandler) IsSupport(ctx context.Context, packetType int32) bool {
	//TODO implement me
	panic("implement me")
}

func (m *MessageReceivingHandler) InitHandler() error {
	m.mqProducer = kafka.InitProducer([]string{""})
	return nil
}
