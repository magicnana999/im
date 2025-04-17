package msg_service

import (
	"errors"
	"github.com/magicnana999/im/api/kitex_gen/api"
	"github.com/magicnana999/im/broker"
	"github.com/magicnana999/im/broker/domain"
	"github.com/magicnana999/im/define"
	bb "github.com/panjf2000/gnet/v2/pkg/pool/bytebuffer"
)

var (
	errUcIsNil      = errors.New("uc is nil")
	errUcIsClosed   = errors.New("uc is closed")
	errMessageIsNil = errors.New("message is nil")
)

// MessageWriter 消息写入服务
type MessageWriter struct {
	codec  *broker.Codec
	logger *broker.Logger
}

func NewMessageWriter(codec *broker.Codec, logger *broker.Logger) *MessageWriter {
	return &MessageWriter{
		codec:  codec,
		logger: logger,
	}
}

func (s *MessageWriter) Write(m *api.Message, uc *domain.UserConn) error {
	if m == nil {
		s.logger.DebugOrError("failed to Write message", "", define.OpWrite, "", errMessageIsNil)
		return errMessageIsNil
	}

	if uc == nil {
		s.logger.DebugOrError("failed to Write message", "", define.OpWrite, m.MessageId, errUcIsNil)
		return errUcIsNil
	}

	if uc.IsClosed.Load() {
		s.logger.DebugOrError("failed to Write message", uc.Desc(), define.OpWrite, m.MessageId, errUcIsClosed)
		return errUcIsClosed
	}

	buffer, err := s.codec.Encode(m.Wrap())
	defer bb.Put(buffer)
	if err != nil {
		s.logger.DebugOrError("failed to Write message", uc.Desc(), define.OpWrite, m.MessageId, err)
		return err
	}

	total := buffer.Len()
	sent := 0
	for sent < total {
		n, err := uc.Writer.Write(buffer.Bytes()[sent:])
		if err != nil {
			s.logger.DebugOrError("failed to Write message", uc.Desc(), define.OpWrite, m.MessageId, err)
			return err
		}
		sent += n
	}

	s.logger.DebugOrError("Write message ok", uc.Desc(), define.OpWrite, m.MessageId, nil)
	return nil
}
