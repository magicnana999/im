package api

import (
	imerror "github.com/magicnana999/im/pkg/error"
	"github.com/magicnana999/im/pkg/id"
	"google.golang.org/protobuf/proto"
	"strings"
	"time"
)

const (
	FlowRequest int32 = iota + 1
	FlowResponse
)

// Type
const (
	TypeHeartbeat int32 = iota + 1
	TypeCommand
	TypeMessage
	TypeEvent
)

// NeedAck
const (
	NO int32 = iota
	YES
)

const (
	CTypeUserLogin      string = "USER_LOGIN"
	CTypeUserLogout            = "USER_LOGOUT"
	CTypeFriendAdd             = "FRIEND_ADD"
	CTypeFriendAddAgree        = "FRIEND_ADD_AGREE"
	CTypeFriendReject          = "FRIEND_ADD_REJECT"
)

const (
	CTypeText  string = "TEXT"
	CTypeImage string = "IMAGE"
	CTypeAudio string = "AUDIO"
	CTypeVideo string = "VIDEO"
)

const (
	TTypeSingle int32 = iota + 1
	TTypeGroup
)

func NewHeartbeat(v int32) *Packet {
	return &Packet{
		Type: TypeHeartbeat,
		Body: &Packet_Heartbeat{Heartbeat: &Heartbeat{Value: v}},
	}
}

func NewCommand(data proto.Message) *Packet {

	mb := &Command{
		Id: strings.ToLower(id.GenerateXId()),
	}

	mb.SetRequest(data)
	return mb.Wrap()
}

func NewMessage(
	userId, to, groupId, sequence int64,
	appId, cId string,
	c proto.Message) *Message {

	mb := &Message{
		MessageId: strings.ToLower(id.GenerateXId()),
		AppId:     appId,
		UserId:    userId,
		ConvId:    cId,
		To:        to,
		GroupId:   groupId,
		Sequence:  sequence,
		Flow:      FlowRequest,
		NeedAck:   YES,
		CTime:     time.Now().UnixMilli(),
	}
	mb.SetContent(c)
	return mb
}

func (p *Packet) IsHeartbeat() bool {
	return p.Type == TypeHeartbeat
}

func (p *Packet) IsCommand() bool {
	return p.Type == TypeCommand
}

func (p *Packet) IsMessage() bool {
	return p.Type == TypeMessage
}

func (p *Packet) Failure(e error) *Packet {
	switch p.Type {
	case TypeHeartbeat:
		return nil
	case TypeCommand:
		return p.GetCommand().Failure(e).Wrap()
	case TypeMessage:
		return p.GetMessage().Failure(e).Wrap()
	default:
		return nil
	}
}

func (p *Packet) Success(c proto.Message) *Packet {
	switch p.Type {
	case TypeHeartbeat:
		return nil
	case TypeCommand:
		return p.GetCommand().Success(c).Wrap()
	case TypeMessage:
		return p.GetMessage().Success(c).Wrap()
	default:
		return nil
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
		Id:    mb.Id,
		CType: mb.CType,
		Code:  0,
	}

	ack.SetReply(reply)

	return ack
}

func (mb *Command) Failure(e error) *Command {

	ack := mb.Success(nil)

	ee := imerror.Format(e)
	ack.Code = int32(ee.Code)
	ack.Message = ee.Message + "," + ee.Detail

	return ack
}

func (mb *Message) Success(content proto.Message) *Message {

	if mb.Flow == FlowResponse {
		return mb
	}

	ack := &Message{
		MessageId: mb.MessageId,
		Flow:      FlowResponse,
		NeedAck:   NO,
		Code:      0,
		//Content:  mb.Content,
	}

	return ack
}

func (mb *Message) Failure(e error) *Message {

	ack := mb.Success(nil)

	ee := imerror.Format(e)
	ack.Code = int32(ee.Code)
	ack.Message = ee.Message + "," + ee.Detail
	return ack
}

func (mb *Command) Wrap() *Packet {
	return &Packet{
		Type: TypeCommand,
		Body: &Packet_Command{
			Command: mb,
		},
	}
}

func (mb *Message) Wrap() *Packet {
	return &Packet{
		Type: TypeMessage,
		Body: &Packet_Message{
			Message: mb,
		},
	}
}

func (mb *Command) SetRequest(content proto.Message) {

	if content == nil {
		return
	}

	switch c := content.(type) {
	case *LoginRequest:
		mb.CType = CTypeUserLogin
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
		mb.CType = CTypeUserLogin
		mb.Reply = &Command_LoginReply{LoginReply: c}
	default:
	}
}

func (mb *Message) SetContent(content proto.Message) {

	if content == nil {
		return
	}

	switch content := content.(type) {
	case *Text:
		mb.CType = CTypeText
		mb.Content = &Message_Text{Text: content}
	case *Image:
		mb.CType = CTypeImage
		mb.Content = &Message_Image{Image: content}
	case *Audio:
		mb.CType = CTypeAudio
		mb.Content = &Message_Audio{Audio: content}
	case *Video:
		mb.CType = CTypeVideo
		mb.Content = &Message_Video{Video: content}
	default:
	}
}

func (mb *Message) IsToGroup() bool {
	return mb.GroupId > 0
}

func (mb *Message) IsRequest() bool {
	if mb.Flow == FlowRequest {
		return true
	}
	return false
}

func (mb *Message) IsResponse() bool {
	if mb.Flow == FlowResponse {
		return true
	}
	return false
}
