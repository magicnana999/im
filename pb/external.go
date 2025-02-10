package pb

import (
	"errors"
	imerror "github.com/magicnana999/im/errors"
	"google.golang.org/grpc/status"
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

func NewCommandRequest(data proto.Message) (*Packet, error) {
	var request isCommandBody_Request
	var cType string
	switch req := data.(type) {
	case *LoginRequest:
		cType = CTypeUserLogin
		request = &CommandBody_LoginRequest{
			LoginRequest: req,
		}
	default:
		return nil, imerror.HandleWrapRequestError.Fill("invalid request type")
	}

	return &Packet{
		Type: TypeCommand,
		Body: &Packet_CommandBody{CommandBody: &CommandBody{CType: cType, Request: request}},
	}, nil

}

func NewCommandResponse(data proto.Message, err error) (*Packet, error) {

	body := &Packet_CommandBody{
		CommandBody: &CommandBody{},
	}

	packet := &Packet{
		Type: TypeCommand,
		Body: body,
	}

	if err != nil {
		c, m := formatCommandError(err)
		body.CommandBody.Code = int32(c)
		body.CommandBody.Message = m
		return packet, nil
	}

	var reply isCommandBody_Reply
	var cType string
	switch rep := data.(type) {
	case *LoginReply:
		cType = CTypeUserLogin
		reply = &CommandBody_LoginReply{
			LoginReply: rep,
		}
	default:
		return nil, imerror.HandleWrapRequestError.Fill("invalid request type")
	}

	body.CommandBody.Code = 0
	body.CommandBody.Message = ""
	body.CommandBody.Reply = reply
	body.CommandBody.CType = cType
	return packet, nil
}

func formatCommandError(err error) (int, string) {

	if err == nil {
		return 0, ""
	}

	if _, ok := status.FromError(err); ok {
		return imerror.HandleGrpcError.Code, err.Error()
	}

	var e imerror.Error
	if b := errors.Is(err, &e); b {
		return imerror.HandleGrpcError.Code, strings.TrimRight(e.Message+" "+e.Details, " ")
	}

	return imerror.HandleInternalError.Code, err.Error()
}

func (p *Packet) IsRequest() bool {
	if p.Type == TypeMessage {
		return p.GetMessageBody().Flow == FlowRequest
	}
	return false
}

func (p *Packet) IsResponse() bool {
	if p.Type == TypeMessage {
		return p.GetMessageBody().Flow == FlowResponse
	}
	return false
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

func (mb *MessageBody) Reply() *MessageBody {

	if mb.Flow == FlowResponse {
		return mb
	}
	return &MessageBody{
		Id:       mb.Id,
		AppId:    mb.AppId,
		UserId:   mb.UserId,
		CId:      mb.CId,
		To:       mb.To,
		GroupId:  mb.GroupId,
		Sequence: mb.Sequence,
		Flow:     FlowResponse,
		NeedAck:  NO,
		CTime:    mb.CTime,
		STime:    time.Now().UnixMilli(),
		CType:    mb.CType,
		Content:  mb.Content,
	}
}

func (mb *MessageBody) Set(content proto.Message) {
	switch content := content.(type) {
	case *TextContent:
		mb.CType = CTypeText
		mb.Content = &MessageBody_TextContent{
			TextContent: content,
		}
	case *ImageContent:
		mb.CType = CTypeImage
		mb.Content = &MessageBody_ImageContent{
			ImageContent: content,
		}
	case *AudioContent:
		mb.CType = CTypeAudio
		mb.Content = &MessageBody_AudioContent{
			AudioContent: content,
		}
	case *VideoContent:
		mb.CType = CTypeVideo
		mb.Content = &MessageBody_VideoContent{
			VideoContent: content,
		}
	default:
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
