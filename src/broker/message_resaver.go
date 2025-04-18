package broker

import (
	"github.com/magicnana999/im/api/kitex_gen/api"
	"github.com/magicnana999/im/broker/domain"
)

type MessageResaver struct {
}

func NewMessageResaver() *MessageResaver {
	return &MessageResaver{}
}

func (mr *MessageResaver) Resave(m *api.Message, uc *domain.UserConn) {
}
