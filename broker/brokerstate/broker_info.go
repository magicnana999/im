package brokerstate

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/magicnana999/im/logger"
	"github.com/magicnana999/im/redis"
	"github.com/timandy/routine"
	"time"
)

type BrokerInfo struct {
	Addr    string `json:"addr"`
	StartAt int64  `json:"startAt"`
}

const (
	KeyBrokerInfo    string        = "im:broker:"
	ExpireBrokerInfo time.Duration = 60 * time.Second
)

func SetBroker(ctx context.Context, broker *BrokerInfo) (string, error) {
	json, err := json.Marshal(broker)
	if err != nil {
		return "", err
	}

	key := fmt.Sprintf("%s%s", KeyBrokerInfo, broker.Addr)
	ret := redis.RDS.Set(ctx, key, json, ExpireBrokerInfo)

	logger.InfoF("[%d] Set broker info ok %v", routine.Goid(), ret)

	return ret.Val(), ret.Err()
}

func RefreshBroker(ctx context.Context, broker *BrokerInfo) (bool, error) {
	key := fmt.Sprintf("%s%s", KeyBrokerInfo, broker.Addr)
	ret := redis.RDS.Expire(ctx, key, ExpireBrokerInfo)
	return ret.Val(), ret.Err()
}
