package kafka

import (
	"context"
	"github.com/magicnana999/im/pb"
	"github.com/magicnana999/im/util/id"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestProducer(t *testing.T) {
	ctx, _ := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)

	producer := InitProducer([]string{"localhost:9092"})
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				producer.SendRoute(ctx, NewMessage().GetMessageBody(), 1)
				producer.SendStore(ctx, NewMessage().GetMessageBody(), 1)
				producer.SendPush(ctx, NewMessage().GetMessageBody(), 1)
				producer.SendOffline(ctx, NewMessage().GetMessageBody(), 1)
			}
		}
	}()

	wg.Wait()

}
func TestConsumer(t *testing.T) {

	ctx, _ := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)

	InitConsumer([]string{"localhost:9092"}, Route, &TestMessageHandler{}).Start(ctx)
	InitConsumer([]string{"localhost:9092"}, Store, &TestMessageHandler{}).Start(ctx)
	InitConsumer([]string{"localhost:9092"}, Push, &TestMessageHandler{}).Start(ctx)
	InitConsumer([]string{"localhost:9092"}, Offline, &TestMessageHandler{}).Start(ctx)

	wg.Wait()

}

type TestMessageHandler struct {
}

func (t TestMessageHandler) Consume(ctx context.Context, msg *pb.MQMessage) error {
	msg.GetMessage()
	return nil
}

func NewMessage() *pb.Packet {
	mb := &pb.MessageBody{
		Id:       strings.ToUpper(id.GenerateXId()),
		AppId:    "1212",
		UserId:   1212,
		CId:      "sdf",
		To:       1111,
		GroupId:  100,
		Sequence: 100,
		Flow:     pb.FlowRequest,
		NeedAck:  pb.YES,
		CTime:    time.Now().UnixMilli(),
	}

	mb.SetContent(&pb.TextContent{
		Text: "hello world",
	})

	return mb.Wrap()
}
