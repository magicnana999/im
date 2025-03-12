package storage

import (
	"context"
	"encoding/json"
	"github.com/magicnana999/im/dto/broker"
	inf "github.com/magicnana999/im/infrastructure"
	"sync"
	"time"
)

var (
	DefaultBrokerStorage *BrokerStorage
	dbsOnce              sync.Once
)

type BrokerStorage struct {
}

func InitBrokerStorage() *BrokerStorage {

	dbsOnce.Do(func() {
		inf.InitRedis(nil)
		DefaultBrokerStorage = &BrokerStorage{}
	})
	return DefaultBrokerStorage
}

func (s *BrokerStorage) StoreBroker(ctx context.Context, broker broker.BrokerInfo) (string, error) {
	json, err := json.Marshal(broker)
	if err != nil {
		return "", err
	}

	key := inf.KeyBroker(broker.Addr)

	ret := inf.RDS.Set(ctx, key, json, time.Minute)
	return ret.Val(), ret.Err()
}

func (s *BrokerStorage) RefreshBroker(ctx context.Context, broker broker.BrokerInfo) (bool, error) {

	key := inf.KeyBroker(broker.Addr)

	ret := inf.RDS.Expire(ctx, key, time.Minute)

	return ret.Val(), ret.Err()
}
