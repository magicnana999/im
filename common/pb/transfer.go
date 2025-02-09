package pb

import (
	"fmt"
	"github.com/magicnana999/im/common/enum"
	"github.com/magicnana999/im/common/protocol"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func ConvertPacket(src *protocol.Packet) (*Packet, error) {

	dest := &Packet{
		Id:      src.Id,
		AppId:   src.AppId,
		UserId:  src.UserId,
		Flow:    src.Flow,
		NeedAck: src.NeedAck,
		BType:   src.BType,
		CTime:   src.CTime,
		STime:   src.STime,
		Status: &Status{
			Code:    src.Code,
			Message: src.Message,
		},
	}

	if body, err := ConvertPacketBody(src.BType, src.Body); err == nil {
		dest.Body = body
		return dest, nil
	} else {
		return nil, err
	}
}

func ConvertPacketBody(pType int32, body any) (*anypb.Any, error) {
	var ret any
	switch pType {
	case BTypeMessage:
		src := body.(protocol.MessageBody)

		dest := &MessageBody{
			CType:      src.CType,
			CId:        src.CId,
			To:         src.To,
			GroupId:    src.GroupId,
			TargetType: src.TargetType,
			Sequence:   src.Sequence,
		}

		if ats, err := ConvertMessageAt(src.At); err == nil {
			dest.At = ats
		} else {
			return nil, err
		}

		if refers, err := ConvertMessageRefer(src.Refer); err == nil {
			dest.Refer = refers
		} else {
			return nil, err
		}

		if content, err := ConvertMessageContent(src.CType, src.Content); err == nil {
			dest.Content = content
		} else {
			return nil, err
		}
		ret = dest
	case BTypeCommand:
		src := body.(protocol.CommandBody)

		dest := &CommandBody{
			CType: src.CType,
		}

		if req, res, err := ConvertCommandContent(src.CType, src.Request, src.Reply); err == nil {
			dest.Request = req
			dest.Reply = res
		} else {
			return nil, err
		}

		ret = dest
	default:
		return nil, fmt.Errorf("unsupported packet type: %v", pType)
	}

	if c, e := anypb.New(ret.(proto.Message)); e == nil {
		return c, nil
	} else {
		return nil, e
	}
}

func ConvertMessageContent(mType string, content any) (*anypb.Any, error) {
	var ret any
	switch mType {
	case CTypeText:
		src := content.(protocol.TextContent)
		ret = &TextContent{
			Text: src.Text,
		}
	case CTypeImage:
		src := content.(protocol.ImageContent)
		ret = &ImageContent{
			Url:    src.Url,
			Width:  src.Width,
			Height: src.Height,
		}
	case CTypeAudio:
		src := content.(protocol.AudioContent)
		ret = &AudioContent{
			Url:    src.Url,
			Length: src.Length,
		}
	case CTypeVideo:
		src := content.(protocol.VideoContent)
		ret = &VideoContent{
			Url:    src.Url,
			Cover:  src.Cover,
			Length: src.Length,
			Width:  src.Width,
			Height: src.Height,
		}
	default:
		return nil, fmt.Errorf("unsupported message type: %v", mType)
	}

	if c, e := anypb.New(ret.(proto.Message)); e == nil {
		return c, nil
	} else {
		return nil, e
	}
}

func ConvertCommandContent(mType string, request any, reply any) (*anypb.Any, *anypb.Any, error) {
	switch mType {

	case CTypeUserLogin:

		var req *anypb.Any
		var res *anypb.Any

		if request != nil {

			input := request.(protocol.LoginRequest)
			ret, err := anypb.New(&LoginRequest{
				AppId:        input.AppId,
				UserSig:      input.UserSig,
				Version:      input.Version,
				Os:           OSType(int32(input.OS)),
				PushDeviceId: input.PushDeviceId,
			})

			if err != nil {
				return nil, nil, err
			}

			req = ret
		}

		if reply != nil {

			output := reply.(protocol.LoginReply)
			ret, err := anypb.New(&LoginReply{
				AppId:  output.AppId,
				UserId: output.UserId,
			})

			if err != nil {
				return nil, nil, err
			}

			res = ret
		}

		return req, res, nil

	default:
		return nil, nil, fmt.Errorf("unsupported message type: %v", mType)
	}
}

func ConvertMessageAt(src []*protocol.At) ([]*At, error) {

	if src == nil || len(src) == 0 {
		return nil, nil
	}

	ats := make([]*At, 0, len(src))
	for _, at := range src {
		ats = append(ats, &At{
			UserId: at.UserId,
			Name:   at.Name,
			Avatar: at.Avatar,
		})
	}
	return ats, nil
}

func ConvertMessageRefer(src []*protocol.Refer) ([]*Refer, error) {

	if src == nil || len(src) == 0 {
		return nil, nil
	}

	refers := make([]*Refer, 0, len(src))
	for _, refer := range src {
		referDest := &Refer{
			UserId: refer.UserId,
			Name:   refer.Name,
			Avatar: refer.Avatar,
		}

		if referContent, err := ConvertMessageContent(refer.CType, refer.Content); err == nil {
			referDest.CType = refer.CType
			referDest.Content = referContent
		} else {
			return nil, err
		}

		refers = append(refers, referDest)
	}
	return refers, nil
}

func RevertPacket(src *Packet) (*protocol.Packet, error) {
	dest := &protocol.Packet{
		Id:      src.Id,
		AppId:   src.AppId,
		UserId:  src.UserId,
		Flow:    src.Flow,
		NeedAck: src.NeedAck,
		BType:   src.BType,
		CTime:   src.CTime,
		STime:   src.STime,
		Code:    src.Status.Code,
		Message: src.Status.Message,
	}

	if body, err := RevertPacketBody(src.BType, src.Body); err == nil {
		dest.Body = body
		return dest, nil
	} else {
		return nil, err
	}
}

func RevertPacketBody(pType int32, body *anypb.Any) (any, error) {
	var ret any
	switch pType {
	case BTypeMessage:

		var src MessageBody
		if err := body.UnmarshalTo(&src); err != nil {
			return nil, err
		}

		dest := &protocol.MessageBody{
			CType:      src.CType,
			CId:        src.CId,
			To:         src.To,
			GroupId:    src.GroupId,
			TargetType: src.TargetType,
			Sequence:   src.Sequence,
		}

		if ats, err := RevertMessageAt(src.At); err == nil {
			dest.At = ats
		} else {
			return nil, err
		}

		if refers, err := RevertMessageRefer(src.Refer); err == nil {
			dest.Refer = refers
		} else {
			return nil, err
		}

		if content, err := RevertMessageContent(src.CType, src.Content); err == nil {
			dest.Content = content
		} else {
			return nil, err
		}
		ret = dest
	case BTypeCommand:
		var src CommandBody
		if err := body.UnmarshalTo(&src); err != nil {
			return nil, err
		}

		dest := &protocol.CommandBody{
			CType: src.CType,
		}

		if req, res, err := RevertCommandContent(src.CType, src.Request, src.Reply); err == nil {
			dest.Request = req
			dest.Reply = res
		} else {
			return nil, err
		}
		ret = dest
	default:
		return nil, fmt.Errorf("unsupported packet type: %v", pType)
	}

	return ret, nil

}

func RevertMessageContent(mType string, content *anypb.Any) (any, error) {
	switch mType {
	case CTypeText:

		var src TextContent
		if err := content.UnmarshalTo(&src); err != nil {
			return nil, err
		}

		ret := &protocol.TextContent{
			Text: src.Text,
		}

		return ret, nil
	case CTypeImage:
		var src ImageContent
		if err := content.UnmarshalTo(&src); err != nil {
			return nil, err
		}

		ret := &protocol.ImageContent{
			Url:    src.Url,
			Width:  src.Width,
			Height: src.Height,
		}

		return ret, nil

	case CTypeAudio:
		var src AudioContent
		if err := content.UnmarshalTo(&src); err != nil {
			return nil, err
		}
		ret := &protocol.AudioContent{
			Url:    src.Url,
			Length: src.Length,
		}

		return ret, nil
	case CTypeVideo:
		var src VideoContent
		if err := content.UnmarshalTo(&src); err != nil {
			return nil, err
		}
		ret := &protocol.VideoContent{
			Url:    src.Url,
			Cover:  src.Cover,
			Length: src.Length,
			Width:  src.Width,
			Height: src.Height,
		}
		return ret, nil

	default:
		return nil, fmt.Errorf("unsupported message type: %v", mType)
	}
}

func RevertCommandContent(mType string, request *anypb.Any, reply *anypb.Any) (any, any, error) {
	switch mType {

	case CTypeUserLogin:

		var req *protocol.LoginRequest
		var res *protocol.LoginReply

		if request != nil {
			var input LoginRequest

			if err := request.UnmarshalTo(&input); err != nil {
				return nil, nil, err
			}

			req = &protocol.LoginRequest{
				AppId:        input.AppId,
				UserSig:      input.UserSig,
				Version:      input.Version,
				OS:           enum.OSType(int32(input.Os)),
				PushDeviceId: input.PushDeviceId,
			}
		}

		if reply != nil {
			var input LoginReply

			if err := reply.UnmarshalTo(&input); err != nil {
				return nil, nil, err
			}

			res = &protocol.LoginReply{
				AppId:  input.AppId,
				UserId: input.UserId,
			}
		}

		return req, res, nil

	default:
		return nil, nil, fmt.Errorf("unsupported message type: %v", mType)
	}
}

func RevertMessageAt(src []*At) ([]*protocol.At, error) {
	if src == nil || len(src) == 0 {
		return nil, nil
	}

	ats := make([]*protocol.At, 0, len(src))
	for _, at := range src {
		ats = append(ats, &protocol.At{
			UserId: at.UserId,
			Name:   at.Name,
			Avatar: at.Avatar,
		})
	}
	return ats, nil
}

func RevertMessageRefer(src []*Refer) ([]*protocol.Refer, error) {

	if src == nil || len(src) == 0 {
		return nil, nil
	}

	refers := make([]*protocol.Refer, 0, len(src))
	for _, refer := range src {
		referDest := &protocol.Refer{
			UserId: refer.UserId,
			Name:   refer.Name,
			Avatar: refer.Avatar,
		}

		if referContent, err := RevertMessageContent(refer.CType, refer.Content); err == nil {
			referDest.CType = refer.CType
			referDest.Content = referContent
		} else {
			return nil, err
		}

		refers = append(refers, referDest)
	}
	return refers, nil
}
