package broker

import (
	"context"
	"github.com/magicnana999/im/errors"
	"strconv"
	"sync"
)

type packetHandler interface {
	handlePacket(ctx context.Context, packet *pb.Packet) (*pb.Packet, error)
	isSupport(ctx context.Context, packetType int32) bool
}

var (
	DefaultHandler *packetHandlerImpl
	hOnce          sync.Once
)

type packetHandlerImpl struct {
	handlers []packetHandler
}

func (p *packetHandlerImpl) handlePacket(ctx context.Context, packet *pb.Packet) (*pb.Packet, error) {
	for _, handler := range p.handlers {
		if handler.isSupport(ctx, packet.Type) {
			return handler.handlePacket(ctx, packet)
		}
	}
	return nil, errors.NoHandlerSupport.SetDetail("type:" + strconv.Itoa(int(packet.Type)))
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

func InitHandler() *packetHandlerImpl {
	hOnce.Do(func() {
		DefaultHandler = &packetHandlerImpl{
			handlers: []packetHandler{initHeartbeatHandler(), initCommandHandler(), initMessageHandler()}}
	})
	return DefaultHandler
}
