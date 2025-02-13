package broker

import (
	"context"
	"github.com/magicnana999/im/logger"
	"github.com/magicnana999/im/pb"
	"sync"
)

var defaultMessageHandler = &messageHandler{}
var messageHandlerMu sync.RWMutex

type messageHandler struct {
	receiver *messageReceiver
	deliver  *deliver
}

func (m *messageHandler) handlePacket(ctx context.Context, p *pb.Packet) (*pb.Packet, error) {

	user, e := currentUserFromCtx(ctx)
	if e != nil {
		return nil, e
	}

	if p.IsRequest() {
		return m.receiver.receive(ctx, p.GetMessageBody())
	}

	if p.IsResponse() {
		m.deliver.ack(&delivery{packet: p, uc: user})
	}

	return nil, nil
}

func (m *messageHandler) isSupport(ctx context.Context, packetType int32) bool {
	return packetType == pb.TypeMessage
}

func initMessageHandler() *messageHandler {

	messageHandlerMu.Lock()
	defer messageHandlerMu.Unlock()

	if defaultMessageHandler.receiver != nil {
		return defaultMessageHandler
	}

	defaultMessageHandler.deliver = defaultDeliver
	defaultMessageHandler.receiver = initMessageReceiver()

	logger.DebugF("messageHandler init")

	return defaultMessageHandler
}
