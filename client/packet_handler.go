package main

import (
	"fmt"
	"github.com/magicnana999/im/api/kitex_gen/api"
	"github.com/magicnana999/im/broker"
	bb "github.com/panjf2000/gnet/v2/pkg/pool/bytebuffer"
	"google.golang.org/protobuf/proto"
	"io"
	"sync"
)

type PacketHandler struct {
	codec *broker.Codec
	mm    *sync.Map
}

func NewPacketHandler() *PacketHandler {
	return &PacketHandler{
		codec: broker.NewCodec(),
	}
}

func (s *PacketHandler) GetMessageSent(msgID string) *api.Message {
	if v, ok := s.mm.Load(msgID); v != nil && ok {
		if req, ok := v.(*api.Message); req != nil && ok {
			return req
		}
	}
	return nil
}

func (s *PacketHandler) ACK(messageId string) {
	s.mm.Delete(messageId)
	return
}

func (s *PacketHandler) Send(request *api.Packet, user *User) {
	if user != nil && !user.IsClosed.Load() {
		s.Write(request, user.Writer)
	}

	if request.IsMessage() {
		m := request.GetMessage()

		if m.IsRequest() {
			s.mm.Store(m.MessageId, s.ACK)
		}
	}
}

func (s *PacketHandler) Handle(p *api.Packet, user *User) *api.Packet {
	if p.IsMessage() {
		if p.GetMessage().IsRequest() {
			json, _ := proto.Marshal(p.GetMessage())
			fmt.Println("收到消息：", string(json))
			return p.GetMessage().Wrap()
		}

		if p.GetMessage().IsResponse() {
			s.ACK(p.GetMessage().GetMessageId())
			return nil
		}

		if p.IsCommand() && p.GetCommand().CommandType == api.CommandTypeUserLogin {
			cmd := p.GetCommand()
			if cmd.GetLoginReply() != nil {
				user.UserID = cmd.GetLoginReply().GetUserId()
				user.AppID = cmd.GetLoginReply().GetAppId()
				user.IsLogin.Store(true)
			}
		}
	}
	return nil
}

func (s *PacketHandler) Write(ret *api.Packet, writer io.Writer) error {
	buffer, err := s.codec.Encode(ret)
	defer bb.Put(buffer)
	if err != nil {
		return err
	}

	total := buffer.Len()
	sent := 0
	for sent < total {
		n, err := writer.Write(buffer.Bytes()[sent:])
		if err != nil {
			return err
		}
		sent += n
	}

	return nil
}
