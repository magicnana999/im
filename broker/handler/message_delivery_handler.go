package handler

import (
	"context"
	"github.com/magicnana999/im/pb"
)

var (
	DefaultMessageDeliveryHandler = &MessageDeliveryHandler{}
)

type MessageDeliveryHandler struct {
}

func (m *MessageDeliveryHandler) HandlePacket(ctx context.Context, packet *pb.Packet) (*pb.Packet, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MessageDeliveryHandler) IsSupport(ctx context.Context, packetType int32) bool {
	//TODO implement me
	panic("implement me")
}

func (m *MessageDeliveryHandler) InitHandler() error {
	//TODO implement me
	panic("implement me")
}

func (m *MessageDeliveryHandler) ReceiveACK(p *pb.Packet) (*pb.Packet, error) {
	panic("ss")
}
