package infrastructure

//
//import (
//	"context"
//	"github.com/magicnana999/im/logger"
//	"github.com/magicnana999/im/api"
//	"github.com/segmentio/kafka-go"
//	"go.opentelemetry.io/otel"
//	"go.opentelemetry.io/otel/propagation"
//	"google.golang.org/protobuf/encoding/protojson"
//	"google.golang.org/protobuf/api"
//	"strings"
//	"sync"
//	"time"
//)
//
//var (
//	once sync.Once
//)
//
//type MQMessageHandler interface {
//	Consume(ctx context.Context, msg *api.MQMessage) error
//}
//
//type Consumer struct {
//	brokers   []string
//	topicInfo TopicInfo
//	handle    MQMessageHandler
//}
//
//func InitConsumer(brokers []string, topic TopicInfo, handle MQMessageHandler) *Consumer {
//
//	c := &Consumer{
//		topicInfo: topic,
//		handle:    handle,
//		brokers:   brokers,
//	}
//
//	return c
//}
//
//func (c *Consumer) Start(ctx context.Context) error {
//	go func() {
//		logger.Infof("consumer start,topic:%s", c.topicInfo)
//
//		reader := kafka.NewReader(kafka.ReaderConfig{
//			Brokers: c.brokers,
//			GroupID: c.topicInfo.Group,
//			Topic:   c.topicInfo.Topic,
//		})
//
//		defer reader.Close()
//
//		for {
//			select {
//			case <-ctx.Done():
//				return
//			default:
//				message, er := reader.ReadMessage(ctx)
//				if er != nil {
//					continue
//				}
//
//				carrier := propagation.MapCarrier{}
//				var traceId string
//				for _, header := range message.Headers {
//					if header.Key == "X-B3-TraceId" {
//						traceId = string(header.Value)
//						carrier.Set("X-B3-TraceId", traceId)
//					}
//				}
//
//				subCtx, _ := context.WithCancel(ctx)
//				subCtx = otel.GetTextMapPropagator().Extract(subCtx, carrier)
//
//				sub, span := logger.Tracer.Start(ctx, c.topicInfo.Topic+"-consumer")
//
//				var msg api.MQMessage
//				if err := api.Unmarshal(message.Value, &msg); err != nil {
//					continue
//				}
//
//				if logger.IsDebugEnable() {
//					js, _ := protojson.Marshal(&msg)
//					logger.Debugf("%s consume message %s input:%s", traceId, c.topicInfo.Topic, string(js))
//				}
//
//				if err := c.handle.Consume(sub, &msg); err != nil && logger.IsDebugEnable() {
//					logger.Errorf("%s consume message %s error:%s", traceId, c.topicInfo.Topic, err.Error())
//				}
//
//				reader.CommitMessages(ctx, message)
//
//				span.End()
//
//			}
//		}
//	}()
//	return nil
//}
//
//func handleMessageRoute(ctx context.Context, topic string, h MQMessageHandler, m *kafka.Message) error {
//	var msg api.MQMessage
//	if err := api.Unmarshal(m.Value, &msg); err != nil {
//		return err
//	}
//
//	logger.Debugf("consume message,topic:%s,id:%s", topic, msg.Id)
//
//	return h.Consume(ctx, &msg)
//}
//
//type Producer struct {
//	writer *kafka.Writer
//}
//
//var defaultProducer *Producer
//
//func InitProducer(brokers []string) *Producer {
//
//	once.Do(func() {
//
//		writer := &kafka.Writer{
//			Addr: kafka.TCP(brokers...), //TCP函数参数为不定长参数，可以传多个地址组成集群
//			//Topic:                  TopicRoute.Topic,
//			Balancer:               &kafka.Hash{}, // 用于对key进行hash，决定消息发送到哪个分区
//			MaxAttempts:            0,
//			WriteBackoffMin:        0,
//			WriteBackoffMax:        0,
//			BatchSize:              0,
//			BatchBytes:             0,
//			BatchTimeout:           0,
//			ReadTimeout:            0,
//			WriteTimeout:           time.Second,      // kafka有时候可能负载很高，写不进去，那么超时后可以放弃写入，用于可以丢消息的场景
//			RequiredAcks:           kafka.RequireAll, // 不需要任何节点确认就返回
//			Async:                  true,
//			Completion:             nil,
//			Compression:            0,
//			Logger:                 nil,
//			ErrorLogger:            nil,
//			Transport:              nil,
//			AllowAutoTopicCreation: false, // 第一次发消息的时候，如果topic不存在，就自动创建topic，工作中禁止使用
//		}
//
//		defaultProducer = &Producer{
//			writer: writer,
//		}
//
//	})
//	return defaultProducer
//}
//
//func (p *Producer) send(ctx context.Context, topic string, m *api.MessageBody, userIds []int64, userLabels []string, count int32) error {
//	if count == 0 {
//		count = 1
//	}
//	mq := &api.MQMessage{
//		Id:         m.MessageId,
//		Len:      count,
//		UserIds:    userIds,
//		UserLabels: userLabels,
//		Message:    m,
//	}
//
//	bs, e := api.Marshal(mq)
//	if e != nil {
//		return e
//	}
//
//	ctx, span := logger.Tracer.Start(ctx, topic)
//	defer span.End()
//
//	traceID := span.SpanContext().TraceID().String()
//
//	if logger.IsDebugEnable() {
//		js, _ := protojson.Marshal(mq)
//		logger.Debugf("%s produce message %s input:%s", traceID, topic, string(js))
//	}
//
//	body := kafka.Message{
//		Topic:      topic,
//		Value:      bs,
//		WriterData: nil,
//		Time:       time.Time{},
//		Headers: []kafka.Header{
//			{Key: "X-B3-TraceId", Value: []byte(traceID)},
//		},
//	}
//
//	err := p.writer.WriteMessages(ctx, body)
//	if err != nil {
//		if logger.IsDebugEnable() && err != nil {
//			logger.Errorf("%s produce message %s error:%s", traceID, topic, err.Error())
//		}
//	}
//
//	return err
//}
//
//func (p *Producer) SendRoute(ctx context.Context, m *api.MessageBody, count int32) error {
//	return p.send(ctx, Route.Topic, m, nil, nil, count)
//}
//
//func (p *Producer) SendRouteDLQ(ctx context.Context, m *api.MessageBody) error {
//	return p.send(ctx, RouteDLQ.Topic, m, nil, nil, 0)
//}
//
//func (p *Producer) SendStore(ctx context.Context, m *api.MessageBody) error {
//	return p.send(ctx, Store.Topic, m, nil, nil, 0)
//}
//
//func (p *Producer) SendOffline(ctx context.Context, m *api.MessageBody, userIds []int64) error {
//	return p.send(ctx, Offline.Topic, m, userIds, nil, 0)
//}
//
//func (p *Producer) SendPush(ctx context.Context, m *api.MessageBody) error {
//	return p.send(ctx, Push.Topic, m, nil, nil, 0)
//}
//
//func (p *Producer) SendDeliver(ctx context.Context, topic string, m *api.MessageBody, labels []string) error {
//	topic = strings.Replace(topic, ":", "-", -1)
//	return p.send(ctx, topic, m, nil, labels, 0)
//}
