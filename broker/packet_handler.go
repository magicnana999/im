package broker

import (
	"context"
	"github.com/magicnana999/im/errors"
	"github.com/magicnana999/im/pb"
	"strconv"
	"sync"
)

var hOnce sync.Once

type packetHandler interface {
	handlePacket(ctx context.Context, packet *pb.Packet) (*pb.Packet, error)
	isSupport(ctx context.Context, packetType int32) bool
}

type packetHandlerImpl struct {
	handlers []packetHandler
}

var defaultHandler *packetHandlerImpl

func (p *packetHandlerImpl) handlePacket(ctx context.Context, packet *pb.Packet) (*pb.Packet, error) {
	for _, handler := range p.handlers {
		if handler.isSupport(ctx, packet.Type) {
			return handler.handlePacket(ctx, packet)
		}
	}
	return nil, errors.NoHandlerSupport.DetailString("type:" + strconv.Itoa(int(packet.Type)))
}

func (p *packetHandlerImpl) isSupport(ctx context.Context, packetType int32) bool {
	for _, handler := range p.handlers {
		if handler.isSupport(ctx, packetType) {
			return true
		}
	}

	return false
}

func (p *packetHandlerImpl) HeartbeatHandler() *heartbeatHandler {
	return defaultHeartbeatHandler
}

func initHandler() *packetHandlerImpl {

	hOnce.Do(func() {
		defaultHandler = &packetHandlerImpl{handlers: make([]packetHandler, 0)}
		defaultHandler.handlers = append(defaultHandler.handlers, initHeartbeatHandler())
		defaultHandler.handlers = append(defaultHandler.handlers, initCommandHandler())
		defaultHandler.handlers = append(defaultHandler.handlers, initMessageHandler())
	})
	return defaultHandler
}
