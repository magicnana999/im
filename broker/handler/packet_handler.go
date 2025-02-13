package handler

import (
	"context"
	"github.com/magicnana999/im/errors"
	"github.com/magicnana999/im/logger"
	"github.com/magicnana999/im/pb"
	"sync"
)

var handlerMu sync.Mutex

type PacketHandler interface {
	HandlePacket(ctx context.Context, packet *pb.Packet) (*pb.Packet, error)
	IsSupport(ctx context.Context, packetType int32) bool
}

type PacketHandlerImpl struct {
	handlers []PacketHandler
}

var DefaultHandler = &PacketHandlerImpl{
	handlers: make([]PacketHandler, 0),
}

func (p *PacketHandlerImpl) HandlePacket(ctx context.Context, packet *pb.Packet) (*pb.Packet, error) {
	for _, handler := range p.handlers {
		if handler.IsSupport(ctx, packet.Type) {
			return handler.HandlePacket(ctx, packet)
		}
	}
	return nil, errors.PacketProcessError
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

func (p *PacketHandlerImpl) HeartbeatHandler() *HeartbeatHandler {
	return DefaultHeartbeatHandler
}

func InitHandler() *PacketHandlerImpl {

	handlerMu.Lock()
	defer handlerMu.Unlock()

	if len(DefaultHandler.handlers) != 0 {
		return DefaultHandler
	}
	DefaultHandler.handlers = append(DefaultHandler.handlers, InitHeartbeatHandler())
	DefaultHandler.handlers = append(DefaultHandler.handlers, InitCommandHandler())
	DefaultHandler.handlers = append(DefaultHandler.handlers, InitMessageHandler())

	return DefaultHandler
}
