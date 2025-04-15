package holder

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/magicnana999/im/broker/domain"
	"github.com/magicnana999/im/infra"
	"go.uber.org/fx"
	"time"
)

type BrokerHolder struct {
	rds *redis.Client
}

func NewBrokerHolder(rds *redis.Client, lf fx.Lifecycle) (*BrokerHolder, error) {
	return &BrokerHolder{rds: rds}, nil
}

func (s *BrokerHolder) StoreBroker(ctx context.Context, broker domain.BrokerInfo) (string, error) {
	json, err := json.Marshal(broker)
	if err != nil {
		return "", err
	}

	key := infra.KeyBroker(broker.Addr)

	ret := s.rds.Set(ctx, key, json, time.Minute)
	return ret.Val(), ret.Err()
}

func (s *BrokerHolder) RefreshBroker(ctx context.Context, broker domain.BrokerInfo) (bool, error) {

	key := infra.KeyBroker(broker.Addr)

	ret := s.rds.Expire(ctx, key, time.Minute)

	return ret.Val(), ret.Err()
}
