package kafka

import (
	"bytes"
	"context"
	"encoding/binary"
	"github.com/magicnana999/im/logger"
	"github.com/magicnana999/im/pb"
	"github.com/panjf2000/ants/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	goPool "github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"time"
)

const (
	kafkaBroker = "localhost:9092"
	maxWorkers  = 24 // worker pool 中的最大并发数

)

var (
	MessageRoute = topicAndGroup{
		"im-message-route", "im-message-route-group", 1}

	executor *goPool.Pool
)

var (
	Producer *kafka.Writer
)

type handle func(message *pb.MessageBody) error

type topicAndGroup struct {
	topic     string
	group     string
	partition int
}

type consumer struct {
	topic     string
	group     string
	partition int
	executor  *goPool.Pool
	handle    handle
}

func (c *consumer) Start(ctx context.Context) error {
	for i := range c.partition {
		e := c.executor.Submit(func() {

			reader := kafka.NewReader(kafka.ReaderConfig{
				Brokers:   []string{kafkaBroker},
				GroupID:   c.group,
				Topic:     c.topic,
				Partition: i,
				MinBytes:  10e3, // 10KB
				MaxBytes:  10e6, // 10MB
			})

			defer reader.Close()

			for {
				select {
				case <-ctx.Done():
					return
				}

				message, er := reader.ReadMessage(ctx)
				if er != nil {
					continue
				}

				if err := handleMessageRoute(c.handle, &message); err != nil {
					continue
				}

				reader.CommitMessages(ctx, message)
			}
		})
		if e != nil {
			return e
		}
	}
	return nil

}

func init() {

	var (
		DefaultAntsPoolSize = maxWorkers
		ExpiryDuration      = 10 * time.Second
		Nonblocking         = true
	)

	options := ants.Options{
		ExpiryDuration: ExpiryDuration,
		Nonblocking:    Nonblocking,
		Logger:         logger.Logger,
		PanicHandler: func(a any) {
			logging.Errorf("goroutine pool panic: %v", a)
		},
	}
	e, _ := ants.NewPool(DefaultAntsPoolSize, ants.WithOptions(options))
	executor = e
}

func InitConsumer(tg topicAndGroup, handle handle) (*consumer, error) {

	c := &consumer{
		topic:     tg.topic,
		group:     tg.group,
		partition: tg.partition,
		executor:  executor,
		handle:    handle,
	}

	return c, nil
}

func InitProducer() {
	Producer = &kafka.Writer{
		Addr:                   kafka.TCP("localhost:9092"), //TCP函数参数为不定长参数，可以传多个地址组成集群
		Topic:                  MessageRoute.topic,
		Balancer:               &kafka.Hash{}, // 用于对key进行hash，决定消息发送到哪个分区
		MaxAttempts:            0,
		WriteBackoffMin:        0,
		WriteBackoffMax:        0,
		BatchSize:              0,
		BatchBytes:             0,
		BatchTimeout:           0,
		ReadTimeout:            0,
		WriteTimeout:           time.Second,      // kafka有时候可能负载很高，写不进去，那么超时后可以放弃写入，用于可以丢消息的场景
		RequiredAcks:           kafka.RequireOne, // 不需要任何节点确认就返回
		Async:                  false,
		Completion:             nil,
		Compression:            0,
		Logger:                 nil,
		ErrorLogger:            nil,
		Transport:              nil,
		AllowAutoTopicCreation: false, // 第一次发消息的时候，如果topic不存在，就自动创建topic，工作中禁止使用
	}
}

func sendMessageRoute(ctx context.Context, packet *pb.Packet) error {

	msg := packet.GetMessageBody()
	bs, e := proto.Marshal(msg)
	if e != nil {
		return e
	}

	buf := new(bytes.Buffer)
	k := binary.Write(buf, binary.BigEndian, msg.GetTo())
	if k != nil {
		return k
	}
	m := kafka.Message{
		Topic:      MessageRoute.topic,
		Value:      bs,
		Headers:    nil,
		WriterData: nil,
		Time:       time.Time{},
	}

	return Producer.WriteMessages(ctx, m)
}

func handleMessageRoute(h handle, m *kafka.Message) error {
	var body pb.MessageBody
	proto.Unmarshal(m.Value, &body)

	return h(&body)
}
