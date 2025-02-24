package broker

import (
	"context"
	"github.com/magicnana999/im/conf"
	"github.com/magicnana999/im/domain"
	"github.com/magicnana999/im/kafka"
	"github.com/magicnana999/im/logger"
	"github.com/magicnana999/im/pb"
	"github.com/panjf2000/ants/v2"
	"github.com/panjf2000/gnet/v2"
	bb "github.com/panjf2000/gnet/v2/pkg/pool/bytebuffer"
	goPool "github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
	"strconv"
	"strings"
	"sync"
	"time"
)

var defaultDeliver *deliver
var dOnce sync.Once

type deliverTask struct {
	ctx      context.Context
	cancel   context.CancelFunc
	id       string
	interval int
	delivery *delivery
	ticker   *time.Ticker
	codec    codec
	conn     gnet.Conn
}

type deliver struct {
	ctx              context.Context
	delivery         chan *delivery
	executor         *goPool.Pool
	m                sync.Map
	codec            codec
	heartbeatHandler *heartbeatHandler
	deliverFailed    func(delivery *delivery) error
	userState        *userState
	mqProducer       *kafka.Producer
	mqConsumer       *kafka.Consumer
}

func initDeliver(ctx context.Context, codec codec) *deliver {

	dOnce.Do(func() {
		d := &deliver{}
		pool, err := ants.NewPool(1024)
		if err != nil {
			logger.Fatalf("init deliver err: %v", err)
		}

		broker := []string{conf.Global.Kafka.String()}
		topic := getTopic()

		d.ctx = ctx
		d.delivery = make(chan *delivery, 4096)
		d.heartbeatHandler = initHeartbeatHandler()
		d.codec = codec
		d.executor = pool
		d.userState = initUserState()
		d.mqProducer = kafka.InitProducer(broker)
		d.mqConsumer = kafka.InitConsumer(broker, topic, d)
		d.deliverFailed = func(delivery *delivery) error {
			eee := d.mqProducer.SendOffline(ctx, delivery.packet.GetMessageBody(), []int64{delivery.uc.UserId})
			if eee != nil {
				defaultInstance.eng.Stop(defaultInstance.ctx)
				return eee
			}

			return nil
		}

		d.mqConsumer.Start(ctx)

		defaultDeliver = d
	})

	return defaultDeliver
}

func (s *deliver) stopAll() {
	s.m.Range(func(key, value interface{}) bool {
		s.stopPacketRetry(value.(*deliverTask).id)
		return true
	})
}

func (s *deliver) send(delivery *delivery) {
	s.delivery <- delivery
}

func (s *deliver) stopPacketRetry(id string) {

	if task, ok := s.m.Load(id); ok && task != nil {
		t := task.(*deliverTask)
		t.cancel()
		s.m.Delete(id)
	}
}

func (s *deliver) start() {

	logger.Infof("deliver start")

	for {

		select {
		case <-s.ctx.Done():
			s.stopAll()
			return
		case d, ok := <-s.delivery:

			if !ok {
				continue
			}

			switch d.packet.Type {
			case pb.TypeMessage:
				s.sendMessage(d)
			}
		}
	}
}

func (s *deliver) sendMessage(delivery *delivery) {

	uc := delivery.uc
	packet := delivery.packet

	_, exist := s.m.Load(packet.GetMessageBody().MessageId)

	if exist {
		return
	}

	subCtx, cancel := context.WithCancel(delivery.ctx)

	task := &deliverTask{id: packet.GetMessageBody().MessageId}
	_, loaded := s.m.LoadOrStore(task.id, task)
	if loaded {
		task = nil
		return
	}

	s.write(uc.C, packet, uc)

	task.ctx = subCtx
	task.cancel = cancel
	task.interval = 2
	task.delivery = delivery
	task.ticker = time.NewTicker(time.Duration(task.interval) * time.Second)
	task.codec = s.codec
	task.conn = delivery.uc.C

	err := s.executor.Submit(func() {

		for {
			select {
			case <-task.ctx.Done():
				return

			case <-task.ticker.C:

				if !s.heartbeatHandler.isRunning(task.conn.Fd()) {
					s.deliverFailed(task.delivery)
					s.stopPacketRetry(packet.GetMessageBody().MessageId)
				}

				next := exponentialBackoff(task.interval)
				if next >= 8 {
					s.deliverFailed(task.delivery)
					s.stopPacketRetry(packet.GetMessageBody().MessageId)
					return

				}
				s.write(task.conn, packet, uc)

				task.interval = next
				task.ticker.Reset(time.Duration(next) * time.Second)
			}
		}
	})

	if err != nil {
		defaultInstance.eng.Stop(defaultInstance.ctx)
	}
}

func (s *deliver) write(conn gnet.Conn, packet *pb.Packet, uc *domain.UserConnection) {

	buffer, err := s.codec.encode(conn, packet)
	defer bb.Put(buffer)

	mb := packet.GetMessageBody()

	if err != nil {

		s.stopPacketRetry(mb.MessageId)
		s.heartbeatHandler.stopTicker(conn)
	}

	total := buffer.Len()
	sent := 0
	for sent < total {
		n, err := conn.Write(buffer.Bytes()[sent:])
		if err != nil {
			s.stopPacketRetry(mb.MessageId)
			s.heartbeatHandler.stopTicker(conn)
		}
		sent += n
	}
}

func (s *deliver) ack(id string) error {
	s.stopPacketRetry(id)
	return nil
}

func (s *deliver) Consume(ctx context.Context, msg *pb.MQMessage) error {
	if msg == nil || msg.UserLabels == nil || len(msg.UserLabels) == 0 {
		return nil
	}

	for _, label := range msg.UserLabels {
		uc := s.userState.loadLocalUser(label)
		if uc == nil {
			userId := splitUserIdFromLabel(label)
			e := s.mqProducer.SendOffline(ctx, msg.Message, []int64{userId})
			if e != nil {
				return e
			}
		} else {
			s.send(&delivery{msg.Message.Wrap(), uc, ctx})
		}
	}
	return nil
}

func getTopic() kafka.TopicInfo {

	t := kafka.TopicInfo{Topic: conf.Global.Broker.Addr, Group: conf.Global.Broker.Addr + "-group"}

	t.Topic = strings.Replace(t.Topic, ":", "-", -1)
	t.Group = strings.Replace(t.Group, ":", "-", -1)
	return t

}

func splitUserIdFromLabel(input string) int64 {
	parts := strings.Split(input, "#")
	if len(parts) > 1 {
		num, _ := strconv.ParseInt(parts[1], 10, 64)
		return num
	}
	return 0
}
func exponentialBackoff(retryCount int) int {
	interval := 1 << uint(retryCount)
	return interval
}

type delivery struct {
	packet *pb.Packet
	uc     *domain.UserConnection
	ctx    context.Context
}
