package redis

import (
	"context"
	"encoding/json"
	"github.com/magicnana999/im/domain"
	"github.com/magicnana999/im/errors"
	"time"
)

var DefaultBrokerStorage = &BrokerStorage{}

type BrokerStorage struct {
}

func InitBrokerStorage() *BrokerStorage {
	initRedis()
	return DefaultBrokerStorage
}

func (s *BrokerStorage) StoreBroker(ctx context.Context, broker domain.BrokerInfo) (string, error) {
	json, err := json.Marshal(broker)
	if err != nil {
		return "", errors.BrokerStoreError.Detail(err)
	}

	key := KeyBroker(broker.Addr)

	ret := rds.Set(ctx, key, json, time.Minute)
	return ret.Val(), ret.Err()
}

func (s *BrokerStorage) RefreshBroker(ctx context.Context, broker domain.BrokerInfo) (bool, error) {

	key := KeyBroker(broker.Addr)

	ret := rds.Expire(ctx, key, time.Minute)

	return ret.Val(), ret.Err()
}
