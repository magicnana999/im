package handler

import (
	"context"
	"github.com/magicnana999/im/pb"
)

var DefaultMessageHandler = &MessageHandler{}

type MessageHandler struct {
	receiver *MessageReceiver
	deliver  *MessageDeliver
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

func InitMessageHandler() *MessageHandler {
	DefaultMessageHandler.deliver = InitMessageDeliver()
	DefaultMessageHandler.receiver = InitMessageReceiver()
	return DefaultMessageHandler
}
