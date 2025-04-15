package infra

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/magicnana999/im/pkg/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"math/rand"
	"time"
)

type SpinLock struct {
	rds *redis.Client
}

func NewSpinLock(rds *redis.Client, lc fx.Lifecycle) *SpinLock {
	return &SpinLock{rds: rds}
}

func (s *SpinLock) Acquire(ctx context.Context, lockKey string, value string, lockExpire time.Duration, maxRetry int) bool {
	for i := 0; i < maxRetry; i++ {
		ok, err := s.rds.SetNX(ctx, lockKey, value, lockExpire).Result()
		if err != nil {
			logger.Error("Error acquiring", zap.Error(err))
			return false
		}
		if ok {
			// 锁成功获取
			return true
		}

		waitTime := time.Millisecond * time.Duration(rand.Intn(100*(i+1)))
		time.Sleep(waitTime)
	}

	return false
}

func (s *SpinLock) Release(ctx context.Context, lockKey string, value string) bool {
	script := `
		if cache.call("GET", KEYS[1]) == ARGV[1] then
			return cache.call("DEL", KEYS[1])
		else
			return 0
		end
	`
	result, err := s.rds.Eval(ctx, script, []string{lockKey}, value).Result()
	if err != nil {
		logger.Error("Error releasing", zap.Error(err))
		return false
	}

	return result.(int64) == 1
}
