package infra

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/magicnana999/im/global"
	"github.com/magicnana999/im/pkg/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	Redis = "redis"
)

type RedisConfig struct {
	*redis.Options
}

func NewRedisClient(lc fx.Lifecycle) *redis.Client {
	c := global.GetRedis()

	if c == nil {
		logger.Fatal("redis configuration not found",
			zap.String(logger.SCOPE, Redis),
			zap.String(logger.OP, Init))
	}

	rds := redis.NewClient(c.Options)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("redis established",
				zap.String(logger.SCOPE, Redis),
				zap.String(logger.OP, Init))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if e := rds.Close(); e != nil {
				logger.Error("reds could not close",
					zap.String(logger.SCOPE, Redis),
					zap.String(logger.OP, Close),
					zap.Error(e))
				return e
			} else {
				logger.Info("redis closed",
					zap.String(logger.SCOPE, Redis),
					zap.String(logger.OP, Close))
				return nil
			}
		},
	})

	return rds
}
