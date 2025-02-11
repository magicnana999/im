package kafka

import (
	"context"
	"github.com/magicnana999/im/pb"
	"github.com/magicnana999/im/util/id"
	"runtime"
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
				producer.sendMessageRoute(ctx, NewMessage())
			}
		}
	}()

	wg.Wait()

}
func TestConsumer(t *testing.T) {

	ctx, _ := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)

	c1, _ := InitConsumer([]string{"localhost:9092"}, runtime.NumCPU()*2, Route, process)
	c1.Start(ctx)

	wg.Wait()

}

func process(m *pb.MessageBody) error {
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

	mb.Set(&pb.TextContent{
		Text: "hello world",
	})

	return mb.Wrap()
}
