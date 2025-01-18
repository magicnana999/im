package broker

import (
	"encoding/json"
	"fmt"
	"github.com/magicnana999/im/logger"
	"github.com/magicnana999/im/redis"
	"time"
)

const (
	KeyBrokerInfo    string        = "im:broker:"
	ExpireBrokerInfo time.Duration = 30 * time.Second
)

type BrokerInfo struct {
	Addr    string `json:"addr"`
	StartAt int64  `json:"startAt"`
}

func NewBrokerInfo(broker *Broker) *BrokerInfo {
	return &BrokerInfo{
		Addr:    broker.GetBrokerAddr(),
		StartAt: broker.StartAt.UnixMilli(),
	}
}

func SetNewBroker(broker *Broker) {
	brokerInfo := NewBrokerInfo(broker)

	json, err := json.Marshal(brokerInfo)
	if err != nil {
		logger.Error(broker.Ctx, "Could not marshal BrokerInfo")
		panic(err)
	}

	key := fmt.Sprintf("%s%s", KeyBrokerInfo, broker.GetBrokerAddr())
	redis.RDS.Set(broker.Ctx, key, json, ExpireBrokerInfo)
}
