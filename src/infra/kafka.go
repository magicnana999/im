package infra

import (
	"context"
	"fmt"
	"github.com/magicnana999/im/define"
	"github.com/magicnana999/im/global"
	log "github.com/magicnana999/im/pkg/logger"
	"github.com/segmentio/kafka-go"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

var (
	Route    = TopicInfo{"msg-route", "msg-route-group"}
	RouteDLQ = TopicInfo{"msg-route-dlq", "msg-route-dlq-group"}
	Store    = TopicInfo{"msg-store", "msg-store-group"}
	Offline  = TopicInfo{"msg-offline", "msg-offline-group"}
	Push     = TopicInfo{"msg-push", "msg-push-group"}
)

type TopicInfo struct {
	Topic string
	Group string
}

// getOrDefaultKafkaConfig 返回 Kafka 配置，优先使用全局配置，缺失时应用默认值。
// 不会修改输入的 global.Config。
func getOrDefaultKafkaConfig(g *global.Config) *global.KafkaConfig {

	c := &global.KafkaConfig{}
	if g != nil && g.Kafka != nil {
		*c = *g.Kafka
	}

	if c.Brokers == nil || len(c.Brokers) == 0 {
		log.Named("kafka").Warn("no Kafka brokers configured, using default (testing only)")
		c.Brokers = []string{"127.0.0.1:9092"}
	}
	return c
}

// NewKafkaProducer 初始化 Kafka 生产者。
// 使用 global.Config 提供配置，通过 fx.Lifecycle 管理生命周期。
// 返回已配置的 kafka.Writer 实例和错误（如果有）。
func NewKafkaProducer(g *global.Config, lc fx.Lifecycle) (*kafka.Writer, error) {

	logger := log.Named("kafka")

	c := getOrDefaultKafkaConfig(g)

	kw := &kafka.Writer{
		Addr: kafka.TCP(c.Brokers...), //TCP函数参数为不定长参数，可以传多个地址组成集群
		//Topic:                  TopicRoute.Topic,
		Balancer:               &kafka.Hash{},                      // 用于对key进行hash，决定消息发送到哪个分区
		MaxAttempts:            3,                                  // 重试 3 次
		WriteBackoffMin:        100 * time.Millisecond,             // 最小退避时间
		WriteBackoffMax:        1 * time.Second,                    // 最大退避时间
		BatchSize:              100,                                // 批量大小
		BatchBytes:             1 << 20,                            // 1MB 批量字节
		BatchTimeout:           100 * time.Millisecond,             // 批量超时
		ReadTimeout:            10 * time.Second,                   // 读取超时
		WriteTimeout:           1 * time.Second,                    // 写入超时
		RequiredAcks:           kafka.RequireOne,                   // 不需要任何节点确认就返回
		Async:                  true,                               //异步
		AllowAutoTopicCreation: false,                              //不允许自动创建topic
		Completion:             nil,                                //write成功后的回调
		Compression:            kafka.Snappy,                       // 使用 Snappy 压缩
		Logger:                 newKafkaLogger(zapcore.DebugLevel), //日志
		ErrorLogger:            newKafkaLogger(zapcore.ErrorLevel), //日志
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("kafka writer established",
				zap.String(define.OP, define.OpInit))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if e := kw.Close(); e != nil {
				logger.Error("kafka writer could not close",
					zap.String(define.OP, define.OpClose),
					zap.Error(e))
				return e
			} else {
				logger.Info("kafka closed",
					zap.String(define.OP, define.OpClose))
				return nil
			}
		},
	})

	return kw, nil
}

type KafkaLogger struct {
	*log.Logger
	level zapcore.Level
}

func newKafkaLogger(level zapcore.Level) *KafkaLogger {
	return &KafkaLogger{Logger: log.Named("kafka"), level: level}
}
func (k KafkaLogger) Printf(s string, i ...interface{}) {
	msg := fmt.Sprintf(s, i...)
	switch k.level {
	case zapcore.DebugLevel:
		k.Debug(msg)
	case zapcore.InfoLevel:
		k.Info(msg)
	case zapcore.WarnLevel:
		k.Warn(msg)
	case zapcore.ErrorLevel:
		k.Error(msg)
	default:
		k.Debug(msg) // 默认 Info
	}
}
