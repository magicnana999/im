package kafka

import (
	"fmt"
	"github.com/magicnana999/im/pb"
	"github.com/segmentio/kafka-go"
	"time"
)

const (
	topicMessageRoute      = "im-message-route"
	topicMessageRouteGroup = "im-message-route-group"
)

var (
	Producer *kafka.Writer
)

func init() {
	Producer = &kafka.Writer{
		Addr:                   kafka.TCP("localhost:9092"), //TCP函数参数为不定长参数，可以传多个地址组成集群
		Topic:                  topicMessageRoute,
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

func SendToRoute(packet *pb.Packet) error {

	var body *pb.MessageBody
	packet.Body.UnmarshalTo(body)

	msg := kafka.Message{
		Topic:      topicMessageRoute,
		Key:        []byte(fmt.Sprintf("%d", user.Id)),
		Value:      msgContent,
		Headers:    nil,
		WriterData: nil,
		Time:       time.Time{},
	}

	err = Producer.WriteMessages(ctx, msg)
	if err != nil {
		fmt.Println(fmt.Sprintf("写入kafka失败，user:%v,err:%v", user, err))
	}
}
