package pb

import (
	"errors"
	"fmt"
	"github.com/magicnana999/im/broker/protocol"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type ProtoSerializer struct {
}

var ProtoSerializerInstance = &ProtoSerializer{}

func (p ProtoSerializer) Serialize(packet *Packet) ([]byte, error) {
	return proto.Marshal(packet)
}

func (p ProtoSerializer) Deserialization(data []byte) (*Packet, error) {
	packet := new(Packet)
	if err := proto.Unmarshal(data, packet); err != nil {
		return nil, err
	}

	switch packet.Type {
	case protocol.TypeMessage:
		var body MessageBody
		if err := packet.Body.UnmarshalTo(&body); err != nil {
			return nil, err
		}

		bodyAny, err := anypb.New(&body)
		if err != nil {
			return nil, err
		}
		packet.Body = bodyAny

	}
	return nil, nil

}

func unmarshalPacketBody(packet *Packet) error {
	switch packet.Type {
	case protocol.TypeMessage:
		var body MessageBody
		if err := packet.Body.UnmarshalTo(&body); err != nil {
			return err
		}
		unmarshalMessageBody(&body)
		bodyAny, err := anypb.New(&body)
		if err != nil {
			return err
		}
		packet.Body = bodyAny
	default:
		return errors.New(fmt.Sprintf("unknown packet type %d", packet.Type))
	}
	return nil
}

func unmarshalMessageBody(mb *MessageBody) error {
	switch mb.MType {
	case protocol.MText:
		var body TextBody
		bodyAny, err := unmarshalContent(&body, mb)
		if err != nil {
			return err
		}
		mb.Content = bodyAny
	case protocol.MImage:
		var body ImageBody
		bodyAny, err := unmarshalContent(&body, mb)
		if err != nil {
			return err
		}
		mb.Content = bodyAny
	case protocol.MAudio:
		var body AudioBody
		bodyAny, err := unmarshalContent(&body, mb)
		if err != nil {
			return err
		}
		mb.Content = bodyAny
	case protocol.MVideo:
		var body VideoBody
		bodyAny, err := unmarshalContent(&body, mb)
		if err != nil {
			return err
		}
		mb.Content = bodyAny
	default:
		return errors.New(fmt.Sprintf("unknown message body mtype %s", mb.MType))
	}
	return nil
}

func unmarshalContent(m proto.Message, mb *MessageBody) (*anypb.Any, error) {
	if err := mb.Content.UnmarshalTo(m); err != nil {
		return nil, err
	}

	bodyAny, err := anypb.New(m)
	if err != nil {
		return nil, err
	}

	return bodyAny, nil
}
