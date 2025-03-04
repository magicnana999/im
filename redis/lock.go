package redis

import (
	"context"
	"github.com/magicnana999/im/logger"
	"math/rand"
	"time"
)

func AcquireLock(ctx context.Context, lockKey string, value string, lockExpire time.Duration, maxRetry int) bool {
	for i := 0; i < maxRetry; i++ {
		ok, err := rds.SetNX(ctx, lockKey, value, lockExpire).Result()
		if err != nil {
			logger.Error("Error acquiring lock: %v", err)
			return false
		}
		if ok {
			// 锁成功获取
			return true
		}

		// 自旋等待，指数退避
		waitTime := time.Millisecond * time.Duration(rand.Intn(100*(i+1)))
		time.Sleep(waitTime)
	}

	return false
}

func ReleaseLock(ctx context.Context, lockKey string, value string) bool {
	script := `
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		else
			return 0
		end
	`
	result, err := rds.Eval(ctx, script, []string{lockKey}, value).Result()
	if err != nil {
		logger.Errorf("Error releasing lock: %v", err)
		return false
	}

	return result.(int64) == 1
}
