package pb

import (
	"errors"
	"github.com/magicnana999/im/util/id"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"time"
)

func IsHeartbeat(p *Packet) bool {
	return p.BType == BTypeHeartbeat
}

func IsCommand(p *Packet) bool {
	return p.BType == BTypeCommand
}

func isMessage(packet *Packet) bool {
	return packet.BType == BTypeMessage
}

func NewCommandRequest(userId int64, ct string, content proto.Message) (*Packet, error) {

	c, e := anypb.New(content)
	if e != nil {
		return nil, e
	}

	b, e := anypb.New(&CommandBody{
		CType:   ct,
		Content: c,
	})
	if e != nil {
		return nil, e
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

	var body *anypb.Any

	if err != nil {
		var pbError Error
		if errors.As(err, &pbError) {
			b, e := anypb.New(&CommandBody{
				CType:   ct,
				Code:    int32(pbError.Code),
				Message: pbError.Message,
			})
			if e != nil {
				return nil, e
			}
			body = b
		} else {
			b, e := anypb.New(&CommandBody{
				CType:   ct,
				Code:    int32(codes.Unknown),
				Message: err.Error(),
			})
			if e != nil {
				return nil, e
			}
			body = b
		}
	} else {
		c, e := anypb.New(content)
		if e != nil {
			return nil, e
		}

		body, e = anypb.New(&CommandBody{
			CType:   ct,
			Content: c,
		})
		if e != nil {
			return nil, e
		}
	}

	response := Packet{
		Id:      packet.Id,
		AppId:   packet.AppId,
		UserId:  packet.UserId,
		Flow:    FlowResponse,
		NeedAck: NO,
		CTime:   packet.CTime,
		STime:   time.Now().UnixMilli(),
		BType:   BTypeCommand,
		Body:    body,
	}

	return &response, nil
}
