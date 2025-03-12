package infrastructure

import (
	"github.com/go-redis/redis/v8"
	"github.com/magicnana999/im/pkg/logger"
	"sync"
)

var (
	RDS     *redis.Client
	rdsOnce sync.Once
)

type RedisConfig struct {
	*redis.Options
}

func InitRedis(config *RedisConfig) *redis.Client {

	if config == nil {
		logger.Fatalf("redis configuration not found")
	}

	rdsOnce.Do(func() {
		RDS = redis.NewClient(config.Options)
	})
	return RDS
}
