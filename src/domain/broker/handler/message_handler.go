package broker

import (
	"context"
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

	mb := p.GetMessageBody()
	if mb.IsRequest() {
		return m.receiver.receive(ctx, mb)
	}

	if mb.IsResponse() {
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

	})

	return defaultMessageHandler
}
