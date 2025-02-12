package consumer

import (
	"context"
	"github.com/magicnana999/im/conf"
	"github.com/magicnana999/im/kafka"
	"github.com/magicnana999/im/pb"
	"runtime"
)

func InitConsumer(ctx context.Context) {
	cpu := runtime.NumCPU()
	brokers := []string{conf.Global.Kafka.String()}
	topic := kafka.Route
	handler := consumeRoute
	consumer, e := kafka.InitConsumer(brokers, cpu, topic, handler)
	if e != nil {
		panic(e)
	}

	consumer.Start(ctx)
}

func consumeRoute(message *pb.MessageBody) error {

}
