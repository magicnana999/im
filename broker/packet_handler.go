package broker

import (
	"context"
	"github.com/magicnana999/im/errors"
	"github.com/magicnana999/im/logger"
	"github.com/magicnana999/im/pb"
	"sync"
)

var handlerMu sync.Mutex

type packetHandler interface {
	handlePacket(ctx context.Context, packet *pb.Packet) (*pb.Packet, error)
	isSupport(ctx context.Context, packetType int32) bool
}

type packetHandlerImpl struct {
	handlers []packetHandler
}

var defaultHandler = &packetHandlerImpl{handlers: make([]packetHandler, 0)}

func (p *packetHandlerImpl) handlePacket(ctx context.Context, packet *pb.Packet) (*pb.Packet, error) {
	for _, handler := range p.handlers {
		if handler.isSupport(ctx, packet.Type) {
			return handler.handlePacket(ctx, packet)
		}
	}
	return nil, errors.PacketProcessError
}

func (p *packetHandlerImpl) isSupport(ctx context.Context, packetType int32) bool {
	for _, handler := range p.handlers {
		if handler.isSupport(ctx, packetType) {
			return true
		}
	}

	logger.ErrorF("no handler support %d", packetType)
	return false
}

func (p *packetHandlerImpl) HeartbeatHandler() *heartbeatHandler {
	return defaultHeartbeatHandler
}

func InitHandler() *packetHandlerImpl {

	handlerMu.Lock()
	defer handlerMu.Unlock()

	if len(defaultHandler.handlers) != 0 {
		return defaultHandler
	}
	defaultHandler.handlers = append(defaultHandler.handlers, initHeartbeatHandler())
	defaultHandler.handlers = append(defaultHandler.handlers, initCommandHandler())
	defaultHandler.handlers = append(defaultHandler.handlers, InitMessageHandler())

	return defaultHandler
}
