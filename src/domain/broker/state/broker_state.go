package state

import (
	"github.com/magicnana999/im/domain/broker/storage"
	"sync"
)

var DefaultBrokerState *brokerState
var dbsOnce sync.Once

type brokerState struct {
	*storage.BrokerStorage
}

func InitBrokerState() *brokerState {

	dbsOnce.Do(func() {
		DefaultBrokerState = &brokerState{storage.InitBrokerStorage()}
	})

	return DefaultBrokerState
}
