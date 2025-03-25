package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
	"testing"
)

func TestConsoleEncode(t *testing.T) {

	log := Init(&Config{File: "console.log", Encode: ConsoleEncode})
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

func TestJsonEncode(t *testing.T) {
	log := Init(&Config{File: "json.json"})

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
	config := &Config{
		File:       "console_logger.log",
		Encode:     ConsoleEncode,
		Level:      int8(zapcore.InfoLevel),
		TimeFormat: YYYYMMDDHHMMSS,
	}

	logger := Init(config)

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
	config := &Config{
		File:       "json_logger.log",
		Encode:     JSONEncode,
		Level:      int8(zapcore.InfoLevel),
		TimeFormat: YYYYMMDDHHMMSS,
	}

	logger := Init(config)

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
	// 初始化 logger 配置为 Console 编码
	config := &Config{
		File:       "/dev/null", // 输出到 /dev/null，避免实际写文件影响性能
		Level:      int8(zapcore.InfoLevel),
		Encode:     ConsoleEncode,
		TimeFormat: YYYYMMDDHHMMSS,
	}
	logger := Init(config)
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
	// 初始化 logger 配置为 JSON 编码
	config := &Config{
		File:       "/dev/null", // 输出到 /dev/null，避免实际写文件影响性能
		Level:      int8(zapcore.InfoLevel),
		Encode:     JSONEncode,
		TimeFormat: YYYYMMDDHHMMSS,
	}
	logger := Init(config)
	defer logger.Sync() // 确保日志同步

	// 重置计时器
	b.ResetTimer()

	// 基准测试循环
	for i := 0; i < b.N; i++ {
		logger.Info("message")
	}
}
