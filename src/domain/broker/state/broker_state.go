package state

import (
	"github.com/magicnana999/im/domain/broker/storage"
	"sync"
)

var DefaultBrokerState *BrokerState
var dbsOnce sync.Once

type BrokerState struct {
	*storage.BrokerStorage
}

func InitBrokerState() *BrokerState {

	dbsOnce.Do(func() {
		DefaultBrokerState = &BrokerState{storage.InitBrokerStorage()}
	})

	return DefaultBrokerState
}
