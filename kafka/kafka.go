package kafka

import (
	"bytes"
	"context"
	"encoding/binary"
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

type handle func(message *pb.MessageBody) error

type Consumer struct {
	brokers   []string
	topicInfo TopicInfo
	executor  *goPool.Pool
	handle    handle
}

type Producer struct {
	writer *kafka.Writer
}

func (c *Consumer) Start(ctx context.Context) error {
	go func() {
		logger.InfoF("%d start consumer,topic:%s", routine.Goid(), c.topicInfo)

		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers: c.brokers,
			GroupID: c.topicInfo.group,
			Topic:   c.topicInfo.topic,
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

				c.executor.Submit(func() {
					if err := handleMessageRoute(c.handle, &message); err != nil {
						logger.ErrorF("%d consume message,topic:%s,error:%v", routine.Goid(), c.topicInfo.topic, err)
						return
					}
					reader.CommitMessages(ctx, message)
				})

			}
		}
	}()
	return nil
}

func handleMessageRoute(h handle, m *kafka.Message) error {
	var msg pb.MessageBody
	if err := proto.Unmarshal(m.Value, &msg); err != nil {
		return err
	}
	logger.InfoF("%d consume message,topic:%s,id:%s", routine.Goid(), m.Topic, msg.Id)

	return h(&msg)
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

func InitConsumer(brokers []string, maxWorkers int, topic TopicInfo, handle handle) (*Consumer, error) {

	e := initExecutor(maxWorkers)

	c := &Consumer{
		topicInfo: topic,
		executor:  e,
		handle:    handle,
		brokers:   brokers,
	}

	logger.InfoF("%d init consumer,executor:%p", routine.Goid(), e)

	return c, nil
}

func InitProducer(brokers []string) *Producer {
	writer := &kafka.Writer{
		Addr: kafka.TCP(brokers...), //TCP函数参数为不定长参数，可以传多个地址组成集群
		//Topic:                  TopicRoute.topic,
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

func (p *Producer) sendMessageRoute(ctx context.Context, packet *pb.Packet) error {

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
		Topic:      Route.topic,
		Value:      bs,
		Headers:    nil,
		WriterData: nil,
		Time:       time.Time{},
	}

	logger.InfoF("%d produce message,topic:%s,id:%s", routine.Goid(), m.Topic, msg.Id)

	return p.writer.WriteMessages(ctx, m)
}
