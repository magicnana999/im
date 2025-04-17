package handler

import (
	"context"
	"github.com/magicnana999/im/api/kitex_gen/api"
	"github.com/magicnana999/im/api/kitex_gen/api/routerservice"
	"github.com/magicnana999/im/broker"
	"go.uber.org/fx"
)

type MessageHandler struct {
	mrs          *broker.MessageRetryServer
	routerClient routerservice.Client
}

func NewMessageHandler(mrs *broker.MessageRetryServer, rc routerservice.Client, lc fx.Lifecycle) (*MessageHandler, error) {
	h := &MessageHandler{
		mrs:          mrs,
		routerClient: rc,
	}

	return h, nil
}

func (m *MessageHandler) handlePacket(ctx context.Context, p *api.Packet) (*api.Packet, error) {

	mb := p.GetMessage()
	if mb.IsRequest() {
		_, err := m.routerClient.Route(ctx, mb)
		return mb.Response(nil, err).Wrap(), nil
	}

	if mb.IsResponse() {
		m.mrs.Ack(mb.MessageId)
	}

	return nil, nil
}
