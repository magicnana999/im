package pb

import (
	"fmt"
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

		if content, err := ConvertContent(src.CType, src.Content); err == nil {
			dest.Content = content
		} else {
			return nil, err
		}
		ret = dest
	case BTypeCommand:
		src := body.(protocol.CommandBody)

		dest := &CommandBody{
			CType:   src.CType,
			Code:    src.Code,
			Message: src.Message,
		}

		if content, err := ConvertContent(src.CType, src.Content); err == nil {
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

func ConvertContent(mType string, content any) (*anypb.Any, error) {
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
	case CTypeUserLogin:
		src := content.(protocol.LoginContent)

		ret = &LoginContent{
			AppId:        src.AppId,
			UserSig:      src.UserSig,
			Version:      src.Version,
			Os:           OSType(int32(src.OS)),
			PushDeviceId: src.PushDeviceId,
		}
	case CTypeUserLogout:
		src := content.(protocol.LogoutContent)

		ret = &LogoutContent{
			AppId:  src.AppId,
			UserId: src.UserId,
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

		if referContent, err := ConvertContent(refer.CType, refer.Content); err == nil {
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

		if content, err := RevertContent(src.CType, src.Content); err == nil {
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
			CType:   src.CType,
			Code:    src.Code,
			Message: src.Message,
		}

		if content, err := RevertContent(src.CType, src.Content); err == nil {
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

func RevertContent(mType string, content *anypb.Any) (any, error) {
	var ret any
	switch mType {
	case CTypeText:

		var src TextContent
		if err := content.UnmarshalTo(&src); err != nil {
			return nil, err
		}

		ret = &protocol.TextContent{
			Text: src.Text,
		}
	case CTypeImage:
		var src ImageContent
		if err := content.UnmarshalTo(&src); err != nil {
			return nil, err
		}

		ret = &protocol.ImageContent{
			Url:    src.Url,
			Width:  src.Width,
			Height: src.Height,
		}
	case CTypeAudio:
		var src AudioContent
		if err := content.UnmarshalTo(&src); err != nil {
			return nil, err
		}
		ret = &protocol.AudioContent{
			Url:    src.Url,
			Length: src.Length,
		}
	case CTypeVideo:
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
	case CTypeUserLogin:
		var src LoginContent
		if err := content.UnmarshalTo(&src); err != nil {
			return nil, err
		}

		ret = &protocol.LoginContent{
			AppId:        src.AppId,
			UserSig:      src.UserSig,
			Version:      src.Version,
			OS:           OSType(int32(src.Os)),
			PushDeviceId: src.PushDeviceId,
		}

	case CTypeUserLogout:

		var src LogoutContent
		if err := content.UnmarshalTo(&src); err != nil {
			return nil, err
		}

		ret = &protocol.LogoutContent{
			AppId:  src.AppId,
			UserId: src.UserId,
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

		if referContent, err := RevertContent(refer.CType, refer.Content); err == nil {
			referDest.CType = refer.CType
			referDest.Content = referContent
		} else {
			return nil, err
		}

		refers = append(refers, referDest)
	}
	return refers, nil
}
