package timewheel

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/magicnana999/im/global"
	"github.com/magicnana999/im/pkg/id"
	"github.com/magicnana999/im/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

var rds *redis.Client

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

func TestMain(m *testing.M) {
	redisConfig := global.RedisConfig{
		Addr:    "127.0.0.1:6379",
		DB:      0,
		Timeout: time.Second,
	}
	o := redisConfig.ToOptions()
	rds = redis.NewClient(o)
	defer rds.Close()

	m.Run()
}

func deleteAllKeys() {

	rds.FlushDB(context.Background())
}

func rangeAllKeys(fun func(ctx context.Context, keys ...string) error) {
	ctx := context.Background()

	var (
		cursor    uint64
		batchSize = 10000 // 每次 SCAN 返回的 key 数
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

		err = fun(ctx, keys...)
		if err != nil {
			break
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
}

func TestLogic1(t *testing.T) {

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
		WorkerCount:       8 * 100,
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

func TestLogic2(t *testing.T) {

	var wg sync.WaitGroup

	fun := func(now time.Time) TaskResult {
		return Retry
	}

	logger.Init(nil)
	defer logger.Close()

	log := logger.NameWithOptions("timewheel", zap.IncreaseLevel(zapcore.DebugLevel))
	tw, _ := NewTimewheel(&Config{
		SlotCount:         1,
		SlotTick:          time.Second,
		WorkerCount:       8 * 100,
		WorkerNonBlocking: false,
	}, log, nil)

	for i := 0; i < 1; i++ {
		for j := 0; j < 1000000; j++ {
			tw.Put(task(fun), i)
		}

	}

	tw.Start(context.Background())

	wg.Add(1)
	wg.Wait()
}

type RedisTask struct {
	slot int
	id   string
}

func (rt *RedisTask) Execute(now time.Time) TaskResult {
	rds.Set(context.Background(), fmt.Sprintf("%d#%s#%d", rt.slot, rt.id, now.Unix()), 1, time.Hour)
	return Break
}

func TestDeleteAllKeys(t *testing.T) {
	deleteAllKeys()
}
func TestLogic1WithRedisTask(t *testing.T) {

	var wg sync.WaitGroup

	logger.Init(nil)
	defer logger.Close()

	log := logger.NameWithOptions("timewheel", zap.IncreaseLevel(zapcore.DebugLevel))
	tw, _ := NewTimewheel(&Config{
		SlotCount:         60,
		SlotTick:          time.Second,
		WorkerCount:       8 * 100,
		WorkerNonBlocking: false,
	}, log, nil)

	for i := 0; i < tw.cfg.SlotCount; i++ {
		for j := 0; j < 18000; j++ {
			tw.Put(&RedisTask{i, id.GenerateXId()}, i)
		}

	}
	fmt.Println("starting:", time.Now().Unix())
	tw.Start(context.Background())

	wg.Add(1)
	wg.Wait()
}

func TestLogic1WithRedisTaskResult(t *testing.T) {

	fmt.Println(rds.DBSize(context.Background()))

	keys, err := rds.Keys(context.Background(), "*").Result()
	if err != nil {
		fmt.Println(err)
		return
	}

	slotCount := int64(60)
	slotStart := int64(58)
	starting := int64(1745696277)
	var total *int = new(int)
	var success *int = new(int)
	var failed *int = new(int)

	for _, key := range keys {
		*total++

		s := strings.Split(key, "#")
		if len(s) != 3 {
			fmt.Println("error key:", key)
			*failed++
			continue
		}

		slot, err := strconv.ParseInt(s[0], 10, 64)
		if err != nil {
			*failed++
			continue
		}

		ticking, err := strconv.ParseInt(s[2], 10, 64)
		if err != nil {
			*failed++
			continue
		}

		if slot < slotStart {
			if ticking-starting == slot+(slotCount-slotStart)+1 {
				*success++
				continue
			} else {
				*failed++
				fmt.Println("error key:", key, "slot:", slot, "starting:", starting, "ticking:", ticking, "gap:", ticking-starting)
			}
		} else if slot == slotStart {
			if ticking-starting == 1 {
				*success++
				continue
			} else {
				*failed++
				fmt.Println("error key:", key, "slot:", slot, "starting:", starting, "ticking:", ticking, "gap:", ticking-starting)
			}
		} else {
			if ticking-starting == slot-slotStart+1 {
				*success++
				continue
			} else {
				*failed++
				fmt.Println("error key:", key, "slot:", slot, "starting:", starting, "ticking:", ticking, "gap:", ticking-starting)
			}
		}
	}

	fmt.Println("total:", *total, "success:", *success, "failed:", *failed)
}

func TestSubmitAVG(t *testing.T) {

}
