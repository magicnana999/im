package broker

import (
	"context"
	"github.com/magicnana999/im/logger"
	"github.com/magicnana999/im/pb"
	"sync"
)

var defaultMessageHandler *messageHandler
var mhOnce sync.Once

type messageHandler struct {
	receiver *messageReceiver
	deliver  *deliver
}

func (m *messageHandler) handlePacket(ctx context.Context, p *pb.Packet) (*pb.Packet, error) {

	user, e := currentUserFromCtx(ctx)
	if e != nil {
		return nil, e
	}

	mb := p.GetMessageBody()
	if mb.IsRequest() {
		logger.DebugF("[%s#%s] receive request %s", user.ClientAddr, user.Label(), mb.MessageId)
		return m.receiver.receive(ctx, mb)
	}

	if mb.IsResponse() {
		logger.DebugF("[%s#%s] receive response %s", user.ClientAddr, user.Label(), mb.MessageId)
		m.deliver.ack(mb.MessageId)
	}

	return nil, nil
}

func (m *messageHandler) isSupport(ctx context.Context, packetType int32) bool {
	return packetType == pb.TypeMessage
}

func initMessageHandler() *messageHandler {

	mhOnce.Do(func() {
		defaultMessageHandler = &messageHandler{}
		defaultMessageHandler.receiver = initMessageReceiver()

		logger.DebugF("messageHandler init")
	})

	return defaultMessageHandler
}
