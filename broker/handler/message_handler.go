package handler

import (
	"context"
	"github.com/magicnana999/im/pb"
)

var DefaultMessageHandler = &MessageHandler{}

type MessageHandler struct {
	receiver *MessageReceivingHandler
	deliver  *MessageDeliveryHandler
}

func (m *MessageHandler) HandlePacket(ctx context.Context, p *pb.Packet) (*pb.Packet, error) {
	if p.IsRequest() {
		return m.receiver.ReceivePacket(ctx, p)
	}

	if p.IsResponse() {
		return m.deliver.ReceiveACK(p)
	}

	return nil, nil
}

func (m *MessageHandler) IsSupport(ctx context.Context, packetType int32) bool {
	return packetType == pb.TypeMessage
}

func (m *MessageHandler) InitHandler() error {
	//m.deliver = DefaultMessageDeliveryHandler
	//m.receiver = DefaultMessageReceivingHandler
	//
	//m.deliver.InitHandler()
	//m.receiver.InitHandler()

	return nil
}
