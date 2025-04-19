package main

import (
	"errors"
	"github.com/magicnana999/im/api/kitex_gen/api"
	"github.com/magicnana999/im/broker"
	bb "github.com/panjf2000/gnet/v2/pkg/pool/bytebuffer"
)

type MessageWriter struct {
	codec *broker.Codec
}

func NewMessageWriter() *MessageWriter {
	return &MessageWriter{
		codec: broker.NewCodec(),
	}
}

func (s *MessageWriter) Write(m *api.Message, user *User) error {

	if user.IsClosed.Load() {
		return errors.New("user is closed")
	}

	buffer, err := s.codec.Encode(m.Wrap())
	defer bb.Put(buffer)
	if err != nil {
		return err
	}

	total := buffer.Len()
	sent := 0
	for sent < total {
		n, err := user.Writer.Write(buffer.Bytes()[sent:])
		if err != nil {
			return err
		}
		sent += n
	}

	return nil
}
