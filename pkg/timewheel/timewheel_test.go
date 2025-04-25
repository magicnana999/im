package timewheel

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/magicnana999/im/global"
	"github.com/magicnana999/im/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"runtime"
	"sync"
	"testing"
	"time"
)

type task func(now time.Time) TaskResult

func (t task) Execute(now time.Time) TaskResult {
	return t(now)
}

func Test_getOrDefaultConfig(t *testing.T) {

	{
		c := &Config{}
		c = getOrDefaultConfig(c)
		assert.Equal(t, c.SlotTick, DefSlotTick)
		assert.Equal(t, c.SlotCount, DefSlotCount)
		assert.Equal(t, c.SlotMaxLength, DefMaxLength)
		assert.Equal(t, c.WorkerCount, runtime.NumCPU()*100)
		assert.Equal(t, c.WorkerNonBlocking, false)
		assert.Equal(t, c.WorkerPreAlloc, false)
		assert.Equal(t, c.WorkerExpire, DefExpire)
		assert.Equal(t, c.TaskInterval, c.SlotTick*time.Duration(c.SlotCount))
	}

	{
		origin := &Config{SlotTick: time.Hour, SlotCount: 10}
		c := getOrDefaultConfig(origin)
		assert.Equal(t, c.SlotTick, origin.SlotTick)
		assert.Equal(t, c.SlotCount, origin.SlotCount)
		assert.Equal(t, c.SlotMaxLength, DefMaxLength)
		assert.Equal(t, c.WorkerCount, runtime.NumCPU()*100)
		assert.Equal(t, c.WorkerNonBlocking, false)
		assert.Equal(t, c.WorkerPreAlloc, false)
		assert.Equal(t, c.WorkerExpire, DefExpire)
		assert.Equal(t, c.TaskInterval, c.SlotTick*time.Duration(c.SlotCount))
	}
}

func TestNewTimewheel(t *testing.T) {
	tw, err := NewTimewheel(&Config{}, nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, tw)
	assert.NotNil(t, tw.slots)
	assert.NotNil(t, tw.logger)
	assert.NotNil(t, tw.cfg)
	assert.NotNil(t, tw.pool)

	assert.Nil(t, tw.cancel)
	assert.False(t, tw.IsRunning.Load())

	ctx, _ := context.WithCancel(context.Background())
	tw.Start(ctx)
	assert.True(t, tw.IsRunning.Load())
	assert.NotNil(t, tw.cancel)

	tw.Stop()
	assert.False(t, tw.IsRunning.Load())
}

type RedisTask struct {
	id int64
}

var rds *redis.Client

func (s *RedisTask) Execute(now time.Time) TaskResult {
	ret := rds.Incr(context.Background(), fmt.Sprintf("%d", s.id))
	if ret.Val() >= 10 {
		return Break
	} else {
		return Retry
	}
}

func TestMain(m *testing.M) {
	redisConfig := global.RedisConfig{
		Addr:    "127.0.0.1:6379",
		DB:      0,
		Timeout: time.Second,
	}
	o := redisConfig.ToOptions()
	rds = redis.NewClient(o)
	defer rds.Close()

	deleteAllKeys()
	m.Run()
}

func deleteAllKeys() {
	ctx := context.Background()

	var (
		deletedCount int64
		cursor       uint64
		batchSize    = 1000 // 每次 SCAN 返回的 key 数
		maxRetries   = 10   // 最大重试次数
	)

	for {
		keys, nextCursor, err := rds.Scan(ctx, cursor, "*", int64(batchSize)).Result()
		if err != nil {
			return
		}

		if keys == nil || len(keys) == 0 || nextCursor == 0 {
			fmt.Println("empty")
			break
		}

		for i := 0; i < maxRetries; i++ {
			ret, err := rds.Del(ctx, keys...).Result()
			if err != nil {
				fmt.Println("delete fail，retry")
			} else {
				fmt.Println("delete", ret)
				deletedCount = deletedCount + ret
				break
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	fmt.Println("deletedCount: ", deletedCount)
}

func TestSubmit2(t *testing.T) {

	var wg sync.WaitGroup

	fun := func(now time.Time) TaskResult {
		return Retry
	}

	logger.Init(nil)
	defer logger.Close()

	log := logger.NameWithOptions("timewheel", zap.IncreaseLevel(zapcore.DebugLevel))
	tw, _ := NewTimewheel(&Config{
		SlotCount:         60,
		SlotTick:          time.Second,
		WorkerCount:       8 * 400,
		WorkerNonBlocking: false,
	}, log, nil)

	for i := 0; i < 60; i++ {
		for j := 0; j < 10000; j++ {
			tw.Put(task(fun), i)
		}

	}

	tw.Start(context.Background())

	wg.Add(1)
	wg.Wait()
}
