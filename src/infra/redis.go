package infra

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/magicnana999/im/define"
	"github.com/magicnana999/im/global"
	"github.com/magicnana999/im/pkg/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"time"
)

func getOrDefaultRedisConfig(g *global.Config) *redis.Options {
	c := &global.RedisConfig{}
	if g != nil && g.Redis != nil {
		*c = *g.Redis
	}

	if c.Addr == "" {
		c.Addr = "127.0.0.1:6379"
	}

	if c.Timeout == 0 {
		c.Timeout = 1 * time.Second
	}

	return c.ToOptions()

}

func NewRedisClient(g *global.Config, lc fx.Lifecycle) *redis.Client {

	log := logger.Named("redis")

	c := getOrDefaultRedisConfig(g)

	rds := redis.NewClient(c)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("redis established",
				zap.String(define.OP, define.OpInit))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if e := rds.Close(); e != nil {
				log.Error("reds could not close",
					zap.String(define.OP, define.OpClose),
					zap.Error(e))
				return e
			} else {
				log.Info("redis closed",
					zap.String(define.OP, define.OpClose))
				return nil
			}
		},
	})

	return rds
}
