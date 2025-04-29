package handler

import (
	"context"
	"github.com/magicnana999/im/api/kitex_gen/api"
	"github.com/magicnana999/im/api/kitex_gen/api/routerservice"
	"go.uber.org/fx"
)

type MessageHandler struct {
	routerClient routerservice.Client
}

func NewMessageHandler(rc routerservice.Client, lc fx.Lifecycle) (*MessageHandler, error) {
	h := &MessageHandler{
		routerClient: rc,
	}

	return h, nil
}

// HandlePacket 处理Message类型的Packet
func (m *MessageHandler) HandlePacket(ctx context.Context, p *api.Packet) (*api.Packet, error) {

	mb := p.GetMessage()
	if mb.IsRequest() {
		_, err := m.routerClient.Route(ctx, mb)
		return mb.Response(nil, err).Wrap(), err
	}

	return nil, nil
}
