package handler

import (
	"context"
	"github.com/magicnana999/im/api/kitex_gen/api"
	"github.com/magicnana999/im/api/kitex_gen/api/routerservice"
	"github.com/magicnana999/im/broker"
	"github.com/magicnana999/im/infra"
	"github.com/magicnana999/im/pkg/singleton"
)

var messageHandleSingleton = singleton.NewSingleton[*MessageHandler]()

type MessageHandler struct {
	mss       *broker.MessageSendServer
	routerCli routerservice.Client
}

func NewMessageHandler(mss *broker.MessageSendServer) *MessageHandler {
	return messageHandleSingleton.Get(func() *MessageHandler {
		return &MessageHandler{
			mss:       mss,
			routerCli: infra.NewRouterClient(),
		}
	})
}

func (m *MessageHandler) handlePacket(ctx context.Context, p *api.Packet) (*api.Packet, error) {

	mb := p.GetMessage()
	if mb.IsRequest() {
		reply, err := m.routerCli.Route(ctx, mb)
		return nil, nil
	}

	//if mb.IsResponse() {
	//	m.deliver.ack(mb.MessageId)
	//}

	return nil, nil
}
