package broker

import (
	"errors"
	"github.com/magicnana999/im/api/kitex_gen/api"
	"github.com/magicnana999/im/broker/domain"
	"github.com/panjf2000/gnet/v2"
	bb "github.com/panjf2000/gnet/v2/pkg/pool/bytebuffer"
	"sync/atomic"
)

var (
	errUcIsClosed  = errors.New("uc is closed")
	errGetShutting = errors.New("tcp server is shutting down")
)

// PacketWriter 消息写入服务
type PacketWriter struct {
	codec     *Codec
	logger    *Logger
	IsClose   atomic.Bool
	IsWriting atomic.Bool
}

func NewPacketWriter(codec *Codec, logger *Logger) *PacketWriter {
	return &PacketWriter{
		codec:  codec,
		logger: logger,
	}
}

func (s *PacketWriter) Write(packet *api.Packet, uc *domain.UserConn) error {

	if s.IsClose.Load() {
		s.logger.PktDebug("tcp is shutting down", uc.Desc(), packet.GetPacketId(), nil, PacketTracking, errGetShutting)
		return errGetShutting
	}

	s.IsWriting.Store(true)

	if uc.IsClosed.Load() {
		s.logger.PktDebug("failed to write message,client closed", uc.Desc(), packet.GetPacketId(), nil, PacketTracking, errUcIsClosed)
		return errUcIsClosed
	}

	buffer, err := s.codec.Encode(packet)
	defer bb.Put(buffer)
	if err != nil {
		s.logger.PktDebug("failed to decode message", uc.Desc(), packet.GetPacketId(), nil, PacketTracking, err)
		return err
	}

	err = uc.Conn.AsyncWrite(buffer.Bytes(), func(c gnet.Conn, err error) error {
		s.logger.PktDebug("write completed", uc.Desc(), packet.GetPacketId(), nil, PacketTracking, err)
		s.IsWriting.Store(false)
		return nil
	})
	return err

	//关闭要区分主动关闭还是被动关闭，主动的关闭，需要等写完之后才能关
}
