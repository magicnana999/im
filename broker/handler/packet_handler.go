package handler

import (
	"context"
	"github.com/magicnana999/im/common/pb"
	"github.com/magicnana999/im/errors"
	"github.com/magicnana999/im/logger"
)

type PacketHandler interface {
	HandlePacket(ctx context.Context, packet *pb.Packet) (*pb.Packet, error)
	IsSupport(ctx context.Context, packetType int32) bool
	InitHandler() error
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
	return nil, errors.HandlerNoSupportError
}

func (p *PacketHandlerImpl) IsSupport(ctx context.Context, packetType int32) bool {
	for _, handler := range p.handlers {
		if handler.IsSupport(ctx, packetType) {
			return true
		}
	}

	logger.ErrorF("no handler support %d", packetType)
	return false
}

func (p *PacketHandlerImpl) InitHandler() error {
	for _, handler := range p.handlers {
		err := handler.InitHandler()
		if err != nil {
			logger.FatalF("init handler error: %v", err)
		}
	}

	return nil
}
