package api

import (
	imerror "github.com/magicnana999/im/pkg/error"
	"github.com/magicnana999/im/pkg/id"
	"google.golang.org/protobuf/proto"
	"strings"
	"time"
)

func NewMessage(
	userId, to, groupId, sequence int64,
	appId, convId string,
	c proto.Message) *Message {

	mb := &Message{
		MessageId: strings.ToLower(id.GenerateXId()),
		AppId:     appId,
		UserId:    userId,
		ConvId:    convId,
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

func (mb *Message) Wrap() *Packet {
	return &Packet{
		Type: TypeMessage,
		Body: &Packet_Message{
			Message: mb,
		},
	}
}

func (mb *Message) Response(content proto.Message, e error) *Message {
	if e == nil {
		return mb.Success(content)
	} else {
		return mb.Failure(e)
	}
}
func (mb *Message) Success(content proto.Message) *Message {

	if mb.Flow == FlowResponse {
		return mb
	}

	ack := &Message{
		MessageId:   mb.MessageId,
		MessageType: mb.MessageType,
		Flow:        FlowResponse,
		NeedAck:     NO,
		Code:        0,
	}

	ack.SetContent(content)

	return ack
}

func (mb *Message) Failure(e error) *Message {

	if mb.Flow == FlowResponse {
		return mb
	}

	ack := &Message{
		MessageId:   mb.MessageId,
		MessageType: mb.MessageType,
		Flow:        FlowResponse,
		NeedAck:     NO,
	}

	ee := imerror.Format(e)
	ack.Code = int32(ee.Code)
	ack.Message = ee.Message + "," + ee.Detail
	return ack
}

func (mb *Message) SetContent(content proto.Message) {

	if content == nil {
		return
	}

	switch content := content.(type) {
	case *Text:
		mb.MessageType = MessageTypeText
		mb.Content = &Message_Text{Text: content}
	case *Image:
		mb.MessageType = MessageTypeImage
		mb.Content = &Message_Image{Image: content}
	case *Audio:
		mb.MessageType = MessageTypeAudio
		mb.Content = &Message_Audio{Audio: content}
	case *Video:
		mb.MessageType = MessageTypeVideo
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
