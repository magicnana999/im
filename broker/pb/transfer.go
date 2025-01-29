package pb

import (
	"fmt"
	"github.com/magicnana999/im/broker/protocol"
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
		Type:    src.Type,
		CTime:   src.CTime,
		STime:   src.STime,
	}

	if body, err := ConvertPacketBody(src.Type, src.Body); err == nil {
		dest.Body = body
		return dest, nil
	} else {
		return nil, err
	}
}

func ConvertPacketBody(pType int32, body any) (*anypb.Any, error) {
	var ret any
	switch pType {
	case protocol.TypeMessage:
		src := body.(protocol.MessageBody)

		dest := &MessageBody{
			MType:    src.MType,
			CId:      src.CId,
			To:       src.To,
			GroupId:  src.GroupId,
			TType:    src.TType,
			Sequence: src.Sequence,
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

		if content, err := ConvertMessageContent(src.MType, src.Content); err == nil {
			dest.Content = content
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
	case protocol.MText:
		src := content.(protocol.TextContent)
		ret = &TextContent{
			Text: src.Text,
		}
	case protocol.MImage:
		src := content.(protocol.ImageContent)
		ret = &ImageContent{
			Url:    src.Url,
			Width:  src.Width,
			Height: src.Height,
		}
	case protocol.MAudio:
		src := content.(protocol.AudioContent)
		ret = &AudioContent{
			Url:    src.Url,
			Length: src.Length,
		}
	case protocol.MVideo:
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

		if referContent, err := ConvertMessageContent(refer.MType, refer.Content); err == nil {
			referDest.MType = refer.MType
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
		Type:    src.Type,
		CTime:   src.CTime,
		STime:   src.STime,
	}

	if body, err := RevertPacketBody(src.Type, src.Body); err == nil {
		dest.Body = body
		return dest, nil
	} else {
		return nil, err
	}
}

func RevertPacketBody(pType int32, body *anypb.Any) (any, error) {
	var ret any
	switch pType {
	case protocol.TypeMessage:

		var src MessageBody
		if err := body.UnmarshalTo(&src); err != nil {
			return nil, err
		}

		dest := &protocol.MessageBody{
			MType:    src.MType,
			CId:      src.CId,
			To:       src.To,
			GroupId:  src.GroupId,
			TType:    src.TType,
			Sequence: src.Sequence,
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

		if content, err := RevertMessageContent(src.MType, src.Content); err == nil {
			dest.Content = content
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
	var ret any
	switch mType {
	case protocol.MText:

		var src TextContent
		if err := content.UnmarshalTo(&src); err != nil {
			return nil, err
		}

		ret = &protocol.TextContent{
			Text: src.Text,
		}
	case protocol.MImage:
		var src ImageContent
		if err := content.UnmarshalTo(&src); err != nil {
			return nil, err
		}

		ret = &protocol.ImageContent{
			Url:    src.Url,
			Width:  src.Width,
			Height: src.Height,
		}
	case protocol.MAudio:
		var src AudioContent
		if err := content.UnmarshalTo(&src); err != nil {
			return nil, err
		}
		ret = &protocol.AudioContent{
			Url:    src.Url,
			Length: src.Length,
		}
	case protocol.MVideo:
		var src VideoContent
		if err := content.UnmarshalTo(&src); err != nil {
			return nil, err
		}
		ret = &protocol.VideoContent{
			Url:    src.Url,
			Cover:  src.Cover,
			Length: src.Length,
			Width:  src.Width,
			Height: src.Height,
		}
	default:
		return nil, fmt.Errorf("unsupported message type: %v", mType)
	}
	return ret, nil
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

		if referContent, err := RevertMessageContent(refer.MType, refer.Content); err == nil {
			referDest.MType = refer.MType
			referDest.Content = referContent
		} else {
			return nil, err
		}

		refers = append(refers, referDest)
	}
	return refers, nil
}
