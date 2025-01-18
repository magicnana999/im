package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var (
	RDS *redis.Client
)

func init() {
	RDS = redis.NewClient(&redis.Options{
		Addr:     "testa.heguang.club:6379",
		Password: "Heguang@789...",
		DB:       15,
	})
}

func Get(ctx context.Context) {
	val, err := RDS.Get(ctx, "heguang:im:sequence:dm:19860220:100000000191_100000000193").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)
}
