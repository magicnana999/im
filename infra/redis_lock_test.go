package infra

import (
	"context"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

var (
	rds  *redis.Client
	lock *RedisLock
)

func TestMain(m *testing.M) {

	miniRedis, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	rds = redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})

	lock = NewRedisLock(rds)

	ret := m.Run()
	miniRedis.Close()
	rds.Close()
	os.Exit(ret)
}
func TestSpin(t *testing.T) {

	var wg sync.WaitGroup
	key := "test:key"
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			ret := lock.SpinLock(context.Background(), key, "vv", time.Minute, 100)
			fmt.Println("SpinLock", i, ret)
			wg.Done()
		}(i)
	}

	wg.Wait()

	fmt.Println(rds.Get(context.Background(), key).Result())

	time.Sleep(time.Millisecond * 10)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			ret := lock.Release(context.Background(), key, "vv")
			fmt.Println("Release", i, ret)
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println("OK")
}

func TestReentrant(t *testing.T) {

	var wg sync.WaitGroup
	key := "test:key"
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			vv := "vv"
			if i%3 == 0 {
				vv = "vv"
			} else {
				vv = strconv.Itoa(i)
			}
			ret := lock.ReentrantLock(context.Background(), key, vv, time.Minute)
			fmt.Println("ReentrantLock", i, ret)
			wg.Done()
		}(i)
	}

	wg.Wait()

	fmt.Println(rds.Get(context.Background(), key).Result())

	time.Sleep(time.Millisecond * 10)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			ret := lock.Release(context.Background(), key, "vv")
			fmt.Println("Release", i, ret)
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println("OK")
}
