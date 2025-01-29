package handler

import (
	"errors"
	"github.com/magicnana999/im/broker/protocol"
	"github.com/panjf2000/gnet/v2"
)

type PacketHandler interface {
	HandlePacket(c gnet.Conn, packet *protocol.Packet) error
	IsSupport(packetType int32) bool
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
		DefaultHeartbeatHandler)
	DefaultHandler.InitHandler()
}

func (p *PacketHandlerImpl) HandlePacket(c gnet.Conn, packet *protocol.Packet) error {
	for _, handler := range p.handlers {
		if handler.IsSupport(packet.Type) {
			return handler.HandlePacket(c, packet)
		}
	}
	return errors.New("not support")
}

func (p *PacketHandlerImpl) IsSupport(packetType int32) bool {
	for _, handler := range p.handlers {
		if handler.IsSupport(packetType) {
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
