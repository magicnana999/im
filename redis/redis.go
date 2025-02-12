package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/magicnana999/im/conf"
	"sync"
)

var (
	rds *redis.Client

	lock sync.RWMutex
)

func initRedis() {
	lock.Lock()
	defer lock.Unlock()
	if rds == nil {
		rds = redis.NewClient(&redis.Options{
			Addr:     conf.Global.Redis.String(),
			Password: "",
			DB:       conf.Global.Redis.Db,
		})
	}
}
