package infrastructure

import (
	"github.com/segmentio/kafka-go"
	"sync"
	"time"
)

var (
	once     sync.Once
	Producer *KafkaProducer
)

type KafkaProducer struct {
	writer  *kafka.Writer
	tracing bool
}

func InitProducer(brokers []string, tracing bool) *KafkaProducer {

	once.Do(func() {

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

		Producer = &KafkaProducer{
			writer:  writer,
			tracing: tracing,
		}

	})
	return Producer
}
