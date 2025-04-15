package infra

import (
	"context"
	"github.com/magicnana999/im/global"
	"github.com/magicnana999/im/pkg/logger"
	"github.com/segmentio/kafka-go"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"time"
)

const (
	Kafka = "kafka"
)

type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
}

func NewProducer(lc fx.Lifecycle) *kafka.Writer {

	c := global.GetKafka()

	if c == nil {
		logger.Fatal("kafka configuration not found",
			zap.String(logger.SCOPE, Kafka),
			zap.String(logger.OP, Init))
	}

	kw := &kafka.Writer{
		Addr: kafka.TCP(c.Brokers...), //TCP函数参数为不定长参数，可以传多个地址组成集群
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

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("kafka writer established",
				zap.String(logger.SCOPE, Kafka),
				zap.String(logger.OP, Init))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if e := kw.Close(); e != nil {
				logger.Error("kafka writer could not close",
					zap.String(logger.SCOPE, Gorm),
					zap.String(logger.OP, Close),
					zap.Error(e))
				return e
			} else {
				logger.Info("gorm closed",
					zap.String(logger.SCOPE, Gorm),
					zap.String(logger.OP, Close))
				return nil
			}
		},
	})

	return kw
}
