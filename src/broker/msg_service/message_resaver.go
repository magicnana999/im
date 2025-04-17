package msg_service

import (
	"github.com/magicnana999/im/api/kitex_gen/api"
	"github.com/magicnana999/im/broker"
	"github.com/magicnana999/im/broker/domain"
	"github.com/magicnana999/im/define"
)

type MessageResaver struct {
	logger *broker.Logger
}

func NewMessageResaver(logger *broker.Logger) *MessageResaver {
	return &MessageResaver{
		logger: logger,
	}
}

func (mr *MessageResaver) Resave(m *api.Message, uc *domain.UserConn) {
	mr.logger.DebugOrError("resave message", uc.Desc(), define.OpResave, m.MessageId, nil)
}
