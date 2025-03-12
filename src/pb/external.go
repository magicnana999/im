package pb

import (
	"google.golang.org/protobuf/proto"
	"strings"
	"time"
)

func NewHeartbeat(v int32) *Packet {
	return &Packet{
		Type: TypeHeartbeat,
		Body: &Packet_HeartbeatBody{HeartbeatBody: &HeartbeatBody{Value: v}},
	}
}

func NewCommand(data proto.Message) *Packet {

	mb := &CommandBody{
		Id: strings.ToLower(id.GenerateXId()),
	}

	mb.SetRequest(data)
	return mb.Wrap()
}

func NewMessage(
	userId, to, groupId, sequence int64,
	appId, cId string,
	c proto.Message) *MessageBody {

	mb := &MessageBody{
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
		return p.GetCommandBody().Failure(e).Wrap()
	case TypeMessage:
		return p.GetMessageBody().Failure(e).Wrap()
	default:
		return nil
	}
}

func (p *Packet) Success(c proto.Message) *Packet {
	switch p.Type {
	case TypeHeartbeat:
		return nil
	case TypeCommand:
		return p.GetCommandBody().Success(c).Wrap()
	case TypeMessage:
		return p.GetMessageBody().Success(c).Wrap()
	default:
		return nil
	}
}

func (mb *CommandBody) Response(reply proto.Message, e error) *CommandBody {
	if e == nil {
		return mb.Success(reply)
	} else {
		return mb.Failure(e)
	}
}

func (mb *CommandBody) Success(reply proto.Message) *CommandBody {

	ack := &CommandBody{
		Id:    mb.Id,
		CType: mb.CType,
		Code:  0,
	}

	ack.SetReply(reply)

	return ack
}

func (mb *CommandBody) Failure(e error) *CommandBody {

	ack := mb.Success(nil)

	ee := imerror.Format(e)
	if ee != nil {
		ack.Code = int32(ee.Code)
		ack.Message = ee.Message + " " + ee.Details
	}

	return ack
}

func (mb *MessageBody) Success(content proto.Message) *MessageBody {

	if mb.Flow == FlowResponse {
		return mb
	}

	ack := &MessageBody{
		MessageId: mb.MessageId,
		Flow:      FlowResponse,
		NeedAck:   NO,
		Code:      0,
		//Content:  mb.Content,
	}

	return ack
}

func (mb *MessageBody) Failure(e error) *MessageBody {

	ack := mb.Success(nil)

	ee := imerror.Format(e)
	if ee != nil {
		ack.Code = int32(ee.Code)
		ack.Message = ee.Message
	}

	return ack
}

func (mb *CommandBody) Wrap() *Packet {
	return &Packet{
		Type: TypeCommand,
		Body: &Packet_CommandBody{
			CommandBody: mb,
		},
	}
}

func (mb *MessageBody) Wrap() *Packet {
	return &Packet{
		Type: TypeMessage,
		Body: &Packet_MessageBody{
			MessageBody: mb,
		},
	}
}

func (mb *CommandBody) SetRequest(content proto.Message) {

	if content == nil {
		return
	}

	switch c := content.(type) {
	case *LoginRequest:
		mb.CType = CTypeUserLogin
		mb.Request = &CommandBody_LoginRequest{LoginRequest: c}
	default:
	}
}

func (mb *CommandBody) SetReply(content proto.Message) {

	if content == nil {
		return
	}

	switch c := content.(type) {
	case *LoginReply:
		mb.CType = CTypeUserLogin
		mb.Reply = &CommandBody_LoginReply{LoginReply: c}
	default:
	}
}

func (mb *MessageBody) SetContent(content proto.Message) {

	if content == nil {
		return
	}

	switch content := content.(type) {
	case *TextContent:
		mb.CType = CTypeText
		mb.Content = &MessageBody_TextContent{TextContent: content}
	case *ImageContent:
		mb.CType = CTypeImage
		mb.Content = &MessageBody_ImageContent{ImageContent: content}
	case *AudioContent:
		mb.CType = CTypeAudio
		mb.Content = &MessageBody_AudioContent{AudioContent: content}
	case *VideoContent:
		mb.CType = CTypeVideo
		mb.Content = &MessageBody_VideoContent{VideoContent: content}
	default:
	}
}

func (mb *MessageBody) IsToGroup() bool {
	return mb.GroupId > 0
}

func (mb *MessageBody) IsRequest() bool {
	if mb.Flow == FlowRequest {
		return true
	}
	return false
}

func (mb *MessageBody) IsResponse() bool {
	if mb.Flow == FlowResponse {
		return true
	}
	return false
}
