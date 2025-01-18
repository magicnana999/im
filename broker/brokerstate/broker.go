package broker

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const (
	DefaultBrokerPort int = 7539
)

type Broker struct {
	IP        string
	Port      int
	StartAt   time.Time
	Heartbeat *time.Ticker
	Ctx       context.Context
	Wg        *sync.WaitGroup
}

func NewBroker(
	ctx context.Context,
	wg *sync.WaitGroup,
	ip string,
	port int) *Broker {

	if port <= 0 {
		port = DefaultBrokerPort
	}

	broker := &Broker{
		IP:        ip,
		Port:      port,
		StartAt:   time.Now(),
		Heartbeat: time.NewTicker(30 * time.Second),
		Ctx:       ctx,
		Wg:        wg,
	}

	defer broker.Heartbeat.Stop()
	return broker
}

func (b *Broker) Start() {
	SetNewBroker(b)
}

func (b *Broker) Stop() {
	b.Wg.Done()
}

func (b *Broker) GetBrokerAddr() string {
	return fmt.Sprintf("%s:%d", b.IP, b.Port)
}
