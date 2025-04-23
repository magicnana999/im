package infra

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisLock struct {
	rds *redis.Client
}

func NewRedisLock(rds *redis.Client) *RedisLock {
	return &RedisLock{rds: rds}
}

func (s *RedisLock) SpinLock(ctx context.Context, lockKey string, value string, lockExpire time.Duration, maxRetry int) bool {
	for i := 0; i < maxRetry; i++ {
		ok, err := s.rds.SetNX(ctx, lockKey, value, lockExpire).Result()
		if err != nil {
			return false
		}
		if ok {
			return true
		}

		waitTime := time.Millisecond * 10
		time.Sleep(waitTime)
	}

	return false
}

func (s *RedisLock) ReentrantLock(ctx context.Context, lockKey string, value string, lockExpire time.Duration) bool {
	script := `
		if redis.call("EXISTS", KEYS[1]) == 0 then
			redis.call("SET", KEYS[1], ARGV[1], "PX", ARGV[2])
			return 1
		end
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			redis.call("PEXPIRE", KEYS[1], ARGV[2])
			return 1
		end
		return 0
	`
	result, err := s.rds.Eval(ctx, script, []string{lockKey}, value, int64(lockExpire/time.Millisecond)).Result()
	if err != nil {
		return false
	}
	return result.(int64) == 1
}

func (s *RedisLock) Release(ctx context.Context, lockKey string, value string) bool {
	script := `
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		else
			return 0
		end
	`
	result, err := s.rds.Eval(ctx, script, []string{lockKey}, value).Result()
	if err != nil {
		return false
	}

	return result.(int64) == 1
}
