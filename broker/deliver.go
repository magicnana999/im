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
	"sync"
	"time"
)

var defaultDeliver = &deliver{}
var lock sync.RWMutex

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
	deliverFailed    func(delivery *delivery)
	mqProducer       *kafka.Producer
}

func initDeliver(ctx context.Context, codec codec) *deliver {

	lock.Lock()
	defer lock.Unlock()

	if defaultDeliver.executor != nil {
		return defaultDeliver
	}

	pool, err := ants.NewPool(1024)
	if err != nil {
		logger.FatalF("init deliver err: %v", err)
	}

	defaultDeliver.ctx = ctx
	defaultDeliver.delivery = make(chan *delivery, 65536)
	defaultDeliver.heartbeatHandler = initHeartbeatHandler()
	defaultDeliver.codec = codec
	defaultDeliver.executor = pool
	defaultDeliver.mqProducer = kafka.InitProducer([]string{conf.Global.Kafka.String()})
	defaultDeliver.deliverFailed = func(delivery *delivery) {
		eee := defaultDeliver.mqProducer.SendOffline(ctx, delivery.packet.GetMessageBody(), 1)
		if eee != nil {
			logger.ErrorF("[%s#%s] deliver task creation error:%v", delivery.uc.ClientAddr, delivery.uc.Label(), err)

			defaultInstance.eng.Stop(defaultInstance.ctx)

		}
	}

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
				if d.packet.IsResponse() {
					s.stopPacketRetry(d.packet.GetMessageBody().GetId())
				} else {
					s.sendMessage(d)
				}
			}
		}
	}
}

func (s *deliver) sendMessage(delivery *delivery) {

	uc := delivery.uc
	packet := delivery.packet

	_, exist := s.m.Load(packet.GetMessageBody().GetId())

	if exist {
		return
	}

	subCtx, cancel := context.WithCancel(s.ctx)

	task := &deliverTask{id: packet.GetMessageBody().GetId()}
	_, loaded := s.m.LoadOrStore(task.id, task)
	if loaded {
		task = nil
		return
	}

	s.write(uc.C, packet, uc)
	logger.InfoF("[%s#%s] deliver,id:%s",
		uc.ClientAddr, uc.Label(), packet.GetMessageBody().GetId())

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
					s.stopPacketRetry(packet.GetMessageBody().GetId())
				}

				next := exponentialBackoff(task.interval)
				if next >= 8 {
					s.deliverFailed(task.delivery)
					s.stopPacketRetry(packet.GetMessageBody().GetId())
					return

				}
				s.write(task.conn, packet, uc)
				logger.InfoF("[%s#%s] deliver retry,id:%s",
					uc.ClientAddr, uc.Label(), packet.GetMessageBody().GetId())

				task.interval = next
				task.ticker.Reset(time.Duration(next) * time.Second)
			}
		}
	})

	if err != nil {
		logger.ErrorF("[%s#%s] deliver task creation error:%v", uc.ClientAddr, uc.Label(), err)
		defaultInstance.eng.Stop(defaultInstance.ctx)
	}
}

func (s *deliver) write(conn gnet.Conn, packet *pb.Packet, uc *domain.UserConnection) {

	buffer, err := s.codec.encode(conn, packet)
	defer bb.Put(buffer)

	if err != nil {

		logger.ErrorF("[%s#%s] deliver encode error:%v", uc.ClientAddr, uc.Label())
		s.stopPacketRetry(packet.GetMessageBody().Id)
		s.heartbeatHandler.stopTicker(conn)
	}

	total := buffer.Len()
	sent := 0
	for sent < total {
		n, err := conn.Write(buffer.Bytes()[sent:])
		if err != nil {
			logger.ErrorF("[%s#%s] deliver write error:%v", uc.ClientAddr, uc.Label())
			s.stopPacketRetry(packet.GetMessageBody().Id)
			s.heartbeatHandler.stopTicker(conn)
		}
		sent += n
	}
}

func (s *deliver) ack(d *delivery) error {
	s.send(d)
	return nil
}

func exponentialBackoff(retryCount int) int {
	interval := 1 << uint(retryCount)
	return interval
}

type delivery struct {
	packet *pb.Packet
	uc     *domain.UserConnection
}
