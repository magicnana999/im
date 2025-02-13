package kafka

import (
	"context"
	"github.com/magicnana999/im/logger"
	"github.com/magicnana999/im/pb"
	"github.com/panjf2000/ants/v2"
	goPool "github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
	"github.com/segmentio/kafka-go"
	"github.com/timandy/routine"
	"google.golang.org/protobuf/proto"
	"sync"
	"time"
)

var (
	executor *goPool.Pool
	lock     sync.RWMutex
)

type MQMessageHandler interface {
	Consume(ctx context.Context, msg *pb.MQMessage) error
}

type Consumer struct {
	brokers   []string
	topicInfo TopicInfo
	handle    MQMessageHandler
}

type Producer struct {
	writer *kafka.Writer
}

func (c *Consumer) Start(ctx context.Context) error {
	go func() {
		logger.InfoF("%d start consumer,Topic:%s", routine.Goid(), c.topicInfo)

		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers: c.brokers,
			GroupID: c.topicInfo.Group,
			Topic:   c.topicInfo.Topic,
		})

		defer reader.Close()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				message, er := reader.ReadMessage(ctx)
				if er != nil {
					continue
				}

				if err := handleMessageRoute(ctx, c.handle, &message); err != nil {
					logger.ErrorF("%d consume message,Topic:%s,error:%v", routine.Goid(), c.topicInfo.Topic, err)
					return
				}
				reader.CommitMessages(ctx, message)

			}
		}
	}()
	return nil
}

func handleMessageRoute(ctx context.Context, h MQMessageHandler, m *kafka.Message) error {
	var msg pb.MQMessage
	if err := proto.Unmarshal(m.Value, &msg); err != nil {
		return err
	}
	logger.InfoF("%d consume message,Topic:%s,id:%s", routine.Goid(), m.Topic, msg.Id)

	return h.Consume(ctx, &msg)
}

func initExecutor(maxWorkers int) *goPool.Pool {

	lock.Lock()
	defer lock.Unlock()

	if executor != nil {
		return executor
	}

	executor, _ = ants.NewPool(maxWorkers)

	logger.InfoF("%d init executor,max:%d", routine.Goid(), maxWorkers)

	return executor
}

func InitConsumer(brokers []string, topic TopicInfo, handle MQMessageHandler) *Consumer {

	c := &Consumer{
		topicInfo: topic,
		handle:    handle,
		brokers:   brokers,
	}

	logger.InfoF("%d init consumer", routine.Goid())

	return c
}

func InitProducer(brokers []string) *Producer {
	writer := &kafka.Writer{
		Addr: kafka.TCP(brokers...), //TCP函数参数为不定长参数，可以传多个地址组成集群
		//Topic:                  TopicRoute.Topic,
		Balancer:               &kafka.Hash{}, // 用于对key进行hash，决定消息发送到哪个分区
		MaxAttempts:            0,
		WriteBackoffMin:        0,
		WriteBackoffMax:        0,
		BatchSize:              0,
		BatchBytes:             0,
		BatchTimeout:           0,
		ReadTimeout:            0,
		WriteTimeout:           time.Second,      // kafka有时候可能负载很高，写不进去，那么超时后可以放弃写入，用于可以丢消息的场景
		RequiredAcks:           kafka.RequireAll, // 不需要任何节点确认就返回
		Async:                  true,
		Completion:             nil,
		Compression:            0,
		Logger:                 nil,
		ErrorLogger:            nil,
		Transport:              nil,
		AllowAutoTopicCreation: false, // 第一次发消息的时候，如果topic不存在，就自动创建topic，工作中禁止使用
	}

	logger.InfoF("%d init producer", routine.Goid())

	return &Producer{
		writer: writer,
	}
}

func (p *Producer) send(ctx context.Context, topic string, m *pb.MessageBody, count int32) error {
	if count == 0 {
		count = 1
	}
	mq := &pb.MQMessage{
		Id:      m.Id,
		Count:   count,
		Message: m,
	}

	bs, e := proto.Marshal(mq)
	if e != nil {
		return e
	}

	body := kafka.Message{
		Topic:      topic,
		Value:      bs,
		Headers:    nil,
		WriterData: nil,
		Time:       time.Time{},
	}

	logger.InfoF("%d produce message,Topic:%s,id:%s", routine.Goid(), body.Topic, mq.Id)

	return p.writer.WriteMessages(ctx, body)
}

func (p *Producer) SendRoute(ctx context.Context, m *pb.MessageBody, count int32) error {
	return p.send(ctx, Route.Topic, m, count)
}

func (p *Producer) SendStore(ctx context.Context, m *pb.MessageBody, count int32) error {
	return p.send(ctx, Store.Topic, m, count)
}

func (p *Producer) SendOffline(ctx context.Context, m *pb.MessageBody, count int32) error {
	return p.send(ctx, Offline.Topic, m, count)
}

func (p *Producer) SendPush(ctx context.Context, m *pb.MessageBody, count int32) error {
	return p.send(ctx, Push.Topic, m, count)
}
