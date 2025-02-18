package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/magicnana999/im/conf"
	"sync"
)

var (
	rds *redis.Client

	once sync.Once
)

func initRedis() *redis.Client {

	once.Do(func() {

		rds = redis.NewClient(&redis.Options{
			Addr:     conf.Global.Redis.String(),
			Password: "",
			DB:       conf.Global.Redis.Db,
		})
	})
	return rds
}
