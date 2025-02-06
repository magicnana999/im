package handler

import (
	"context"
	"errors"
	"github.com/magicnana999/im/common/pb"
)

type PacketHandler interface {
	HandlePacket(ctx context.Context, packet *pb.Packet) (*pb.Packet, error)
	IsSupport(ctx context.Context, packetType int32) bool
	InitHandler()
}

type PacketHandlerImpl struct {
	handlers []PacketHandler
}

var DefaultHandler = &PacketHandlerImpl{
	handlers: make([]PacketHandler, 0),
}

func init() {

	DefaultHandler.handlers = append(DefaultHandler.handlers,
		DefaultHeartbeatHandler, DefaultCommandHandler)
	DefaultHandler.InitHandler()
}

func (p *PacketHandlerImpl) HandlePacket(ctx context.Context, packet *pb.Packet) (*pb.Packet, error) {
	for _, handler := range p.handlers {
		if handler.IsSupport(ctx, packet.BType) {
			return handler.HandlePacket(ctx, packet)
		}
	}
	return nil, errors.New("not support")
}

func (p *PacketHandlerImpl) IsSupport(ctx context.Context, packetType int32) bool {
	for _, handler := range p.handlers {
		if handler.IsSupport(ctx, packetType) {
			return true
		}
	}
	return false
}

func (p *PacketHandlerImpl) InitHandler() {
	for _, handler := range p.handlers {
		handler.InitHandler()
	}
}
