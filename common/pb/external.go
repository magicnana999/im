package pb

import (
	"errors"
	imerror "github.com/magicnana999/im/errors"
	"github.com/magicnana999/im/util/id"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"strings"
	"time"
)

func IsHeartbeat(p *Packet) bool {
	return p.BType == BTypeHeartbeat
}

func IsCommand(p *Packet) bool {
	return p.BType == BTypeCommand
}

func IsMessage(packet *Packet) bool {
	return packet.BType == BTypeMessage
}

func IsResponse(p *Packet) bool {
	return p.Flow == FlowResponse
}

func NewHeartbeatRequest(v int32) (*Packet, error) {
	body, _ := anypb.New(wrapperspb.UInt32(uint32(v)))
	packet := &Packet{
		BType: BTypeHeartbeat,
		Body:  body,
	}
	return packet, nil
}

func NewHeartbeatResponse() (*Packet, error) {
	body, _ := anypb.New(wrapperspb.UInt32(0))
	packet := &Packet{
		BType: BTypeHeartbeat,
		Body:  body,
	}
	return packet, nil
}

func NewCommandRequest(userId int64, ct string, content proto.Message) (*Packet, error) {

	c, e := anypb.New(content)
	if e != nil {
		return nil, imerror.HandleWrapRequestError.Fill(e.Error())
	}

	b, e := anypb.New(&CommandBody{
		CType:   ct,
		Request: c,
	})
	if e != nil {
		return nil, imerror.HandleWrapRequestError.Fill(e.Error())
	}

	packet := Packet{
		Id:      id.GenerateXId(),
		AppId:   "19860220",
		UserId:  userId,
		Flow:    FlowRequest,
		NeedAck: YES,
		CTime:   time.Now().UnixMilli(),
		BType:   BTypeCommand,
		Body:    b,
	}

	return &packet, nil
}

func NewCommandResponse(packet *Packet, ct string, content proto.Message, err error) (*Packet, error) {

	normal := &Packet{
		Id:      packet.Id,
		AppId:   packet.AppId,
		UserId:  packet.UserId,
		Flow:    FlowResponse,
		NeedAck: NO,
		CTime:   packet.CTime,
		STime:   time.Now().UnixMilli(),
		BType:   BTypeCommand,
	}

	if err != nil {
		c, m := formatCommandError(err)

		status := &Status{
			Code:    int32(c),
			Message: m,
		}

		normal.Status = status
		return normal, nil
	}

	var body *anypb.Any

	c, e := anypb.New(content)
	if e != nil {
		return nil, imerror.HandleWrapReplyError.Fill(e.Error())
	}

	body, e = anypb.New(&CommandBody{
		CType: ct,
		Reply: c,
	})
	if e != nil {
		return nil, imerror.HandleWrapReplyError.Fill(e.Error())
	}
	status := &Status{
		Code:    int32(0),
		Message: "",
	}
	normal.Status = status
	normal.Body = body
	return normal, nil
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
