package pb

import (
	imerror "github.com/magicnana999/im/pkg/error"
	"github.com/magicnana999/im/pkg/id"
	"google.golang.org/protobuf/proto"
	"strings"
)

func NewCommand(data proto.Message) *Packet {

	mb := &Command{
		CommandId: strings.ToLower(id.GenerateXId()),
	}

	mb.SetRequest(data)
	return mb.Wrap()
}

func (mb *Command) Wrap() *Packet {
	return &Packet{
		Type: TypeCommand,
		Body: &Packet_Command{
			Command: mb,
		},
	}
}

func (mb *Command) Response(reply proto.Message, e error) *Command {
	if e == nil {
		return mb.Success(reply)
	} else {
		return mb.Failure(e)
	}
}

func (mb *Command) Success(reply proto.Message) *Command {

	ack := &Command{
		CommandId:   mb.CommandId,
		CommandType: mb.CommandType,
		Code:        0,
	}

	ack.SetReply(reply)

	return ack
}

func (mb *Command) Failure(e error) *Command {

	ack := &Command{
		CommandId:   mb.CommandId,
		CommandType: mb.CommandType,
	}

	ee := imerror.Format(e)
	ack.Code = int32(ee.Code)
	ack.Message = ee.Message + "," + ee.Detail

	return ack
}

func (mb *Command) SetRequest(content proto.Message) {

	if content == nil {
		return
	}

	switch c := content.(type) {
	case *LoginRequest:
		mb.CommandType = CommandTypeUserLogin
		mb.Request = &Command_LoginRequest{LoginRequest: c}
	default:
	}
}

func (mb *Command) SetReply(content proto.Message) {

	if content == nil {
		return
	}

	switch c := content.(type) {
	case *LoginReply:
		mb.CommandType = CommandTypeUserLogin
		mb.Reply = &Command_LoginReply{LoginReply: c}
	default:
	}
}
