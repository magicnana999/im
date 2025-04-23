package api

import "errors"

var (
	InvalidMessage     = errors.New("message is nil")
	InvalidMessageId   = errors.New("message ID is nil")
	InvalidMessageType = errors.New("message type is nil")
	InvalidAppId       = errors.New("appId is empty")
	InvalidFlow        = errors.New("flow is zero")
	InvalidUserId      = errors.New("userId is zero")
	InvalidConvId      = errors.New("convId is empty")
	InvalidSequence    = errors.New("sequence is zero")
	InvalidCTime       = errors.New("cTime is zero")
	InvalidToGroupId   = errors.New("both to and groupId are zero")
)

func (mb *Message) Validate() error {
	if mb == nil {
		return InvalidMessage
	}

	if mb.MessageId == "" {
		return InvalidMessageId
	}

	if mb.MessageType == "" {
		return InvalidMessageType
	}

	if mb.AppId == "" {
		return InvalidAppId
	}

	if mb.Flow == 0 {
		return InvalidFlow
	}

	if mb.UserId == 0 {
		return InvalidUserId
	}

	if mb.ConvId == "" {
		return InvalidConvId
	}

	if mb.Sequence == 0 {
		return InvalidSequence
	}

	if mb.CTime == 0 {
		return InvalidCTime
	}

	if mb.To == 0 && mb.GroupId == 0 {
		return InvalidToGroupId
	}

	return nil
}
