package broker

import (
	"github.com/magicnana999/im/redis"
)

var defaultBrokerState = &brokerState{}

type brokerState struct {
	*redis.BrokerStorage
}

func initBrokerState() *brokerState {
	defaultBrokerState.BrokerStorage = redis.InitBrokerStorage()
	return defaultBrokerState
}
