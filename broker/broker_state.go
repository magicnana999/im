package broker

import (
	"github.com/magicnana999/im/redis"
	"sync"
)

var defaultBrokerState *brokerState
var dbsOnce sync.Once

type brokerState struct {
	*redis.BrokerStorage
}

func initBrokerState() *brokerState {

	dbsOnce.Do(func() {
		defaultBrokerState = &brokerState{}
		defaultBrokerState.BrokerStorage = redis.InitBrokerStorage()
	})

	return defaultBrokerState
}
