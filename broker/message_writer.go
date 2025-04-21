package broker

import (
	"errors"
	"github.com/magicnana999/im/api/kitex_gen/api"
	"github.com/magicnana999/im/broker/domain"
	bb "github.com/panjf2000/gnet/v2/pkg/pool/bytebuffer"
)

var (
	errUcIsClosed = errors.New("uc is closed")
)

// MessageWriter 消息写入服务
type MessageWriter struct {
	codec  *Codec
	logger *Logger
}

func NewMessageWriter(codec *Codec, logger *Logger) *MessageWriter {
	return &MessageWriter{
		codec:  codec,
		logger: logger,
	}
}

func (s *MessageWriter) Write(m *api.Message, uc *domain.UserConn) error {

	if uc.IsClosed.Load() {
		s.logger.MsgDebug("failed to write message", uc.Desc(), m.MessageId, MsgTracking, errUcIsClosed)
		return errUcIsClosed
	}

	buffer, err := s.codec.Encode(m.Wrap())
	defer bb.Put(buffer)
	if err != nil {
		s.logger.MsgDebug("failed to write message", uc.Desc(), m.MessageId, MsgTracking, err)
		return err
	}

	total := buffer.Len()
	sent := 0
	for sent < total {
		n, err := uc.Writer.Write(buffer.Bytes()[sent:])
		if err != nil {
			s.logger.MsgDebug("failed to write message", uc.Desc(), m.MessageId, MsgTracking, err)
			return err
		}
		sent += n
	}

	return nil
}
