package state

import (
	"github.com/magicnana999/im/redis"
)

var DefaultBrokerState = &BrokerState{}

type BrokerState struct {
	*redis.BrokerStorage
}

func InitBrokerState() *BrokerState {
	DefaultBrokerState.BrokerStorage = redis.InitBrokerStorage()
	return DefaultBrokerState
}
