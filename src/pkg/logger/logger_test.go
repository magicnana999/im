package logger

import (
	"go.uber.org/zap"
	"sync"
	"testing"
)

func TestConsoleEncode(t *testing.T) {

	log := Init(nil)
	defer log.Sync()

	for i := 0; i < 3; i++ {

		log.Debug("debug1")
		log.Info("info1")
		log.Warn("warn1")
		log.Error("error1")

		Debug("debug2")
		Info("info2")
		Warn("warn2")
		Error("error2")

	}

}

func TestJsonEncode(t *testing.T) {
	log := Init(nil)

	log.Debug("debug")
	log.Info("info")
	log.Warn("warn")
	log.Error("error")

	Debug("debug")
	Info("info")
	Warn("warn")
	Error("error")

	With(zap.String("key", "value"), zap.Int("haha", 100)).Info("debug333")
}

// go test -v -race -run TestConcurrentConsoleLoggingRace
func TestConcurrentConsoleLoggingRace(t *testing.T) {
	logger := Init(nil)

	var wg sync.WaitGroup
	numWorkers := 100

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				logger.Info("message")
			}
		}(i)
	}

	wg.Wait()
}

// go test -v -race -run TestConcurrentJSONLoggingRace
func TestConcurrentJSONLoggingRace(t *testing.T) {

	logger := Init(nil)

	var wg sync.WaitGroup
	numWorkers := 100

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				logger.Info("message")
			}
		}(i)
	}

	wg.Wait()
}

// go test -bench=BenchmarkConsoleEncoding -benchmem -benchtime=3s
//
//	132058             27536 ns/op             435 B/op         10 allocs/op
//	132058			运行次数
//	27536 ns/op		每次操作平均时间（纳秒）		总运行时间 ≈ 132058 × 27536 ns ≈ 3.64 秒
//	435 B/op		每次操作分配内存大小（字节）		每次分配435字节
//	10 allocs/op	每次操作内存分配次数（次数）		每次涉及10次内存分配
func BenchmarkConsoleEncoding(b *testing.B) {
	logger := Init(nil)
	defer logger.Sync() // 确保日志同步

	// 重置计时器，避免初始化影响结果
	b.ResetTimer()

	// 基准测试循环
	for i := 0; i < b.N; i++ {
		logger.Info("message")
	}
}

// go test -bench=BenchmarkJSONEncoding -benchmem -benchtime=3s
//
// 116582             29181 ns/op             265 B/op          2 allocs/op
func BenchmarkJSONEncoding(b *testing.B) {
	logger := Init(nil)
	defer logger.Sync() // 确保日志同步

	// 重置计时器
	b.ResetTimer()

	// 基准测试循环
	for i := 0; i < b.N; i++ {
		logger.Info("message")
	}
}
