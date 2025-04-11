package logger

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	YYYYMMDDHHMMSS = "2006-01-02 15:04:05"
	RotationTime   = 24 * time.Hour      // 每天轮转
	MaxAge         = 30 * 24 * time.Hour // 保留30天
	RotationSize   = 100 * 1024 * 1024   // 100MB 分割
	//RotationSize = 1024 // 10MB 分割
)

type EncodeType string

const (
	ConsoleEncode EncodeType = "console"
	JSONEncode    EncodeType = "json"
)

var (
	z        *zap.Logger // 全局 Logger
	instance *Logger
	tracing  = false
	once     sync.Once
)

type Logger struct {
	*zap.Logger
}

type Config struct {
	Dir        string     `json:"dir"` // 日志目录
	TracerName string     `json:"tracerName"`
	Level      int8       `json:"level"`
	Encode     EncodeType `json:"encode"`
	TimeFormat string     `json:"timeFormat"`
}

var defaultConfig = Config{
	Dir:        "./logs",
	TracerName: "",
	Level:      0,
	Encode:     JSONEncode,
	TimeFormat: YYYYMMDDHHMMSS,
}

func getDefaultConfig(c *Config) *Config {
	if c == nil {
		return &defaultConfig
	}

	if c.Dir == "" {
		c.Dir = "./logs"
	}

	if c.Encode == "" {
		c.Encode = JSONEncode
	}

	if c.TimeFormat == "" {
		c.TimeFormat = YYYYMMDDHHMMSS
	}

	return c
}

func Init(c *Config) *Logger {
	once.Do(func() {
		instance = _init(c)
	})

	return instance
}

func _init(c *Config) *Logger {
	c = getDefaultConfig(c)

	// 确保日志目录存在
	if err := os.MkdirAll(c.Dir, 0755); err != nil {
		panic(fmt.Sprintf("create log dir failed: %v", err))
	}

	// 创建不同级别的 writer
	infoWriter, err := rotatelogs.New(
		filepath.Join(c.Dir, "info.%Y-%m-%d.log"),
		rotatelogs.WithLinkName(filepath.Join(c.Dir, "info.log")),
		rotatelogs.WithRotationTime(RotationTime),
		rotatelogs.WithMaxAge(MaxAge),
		rotatelogs.WithRotationSize(RotationSize),
	)
	if err != nil {
		panic(fmt.Sprintf("init info writer failed: %v", err))
	}

	errorWriter, err := rotatelogs.New(
		filepath.Join(c.Dir, "error.%Y-%m-%d.log"),
		rotatelogs.WithLinkName(filepath.Join(c.Dir, "error.log")),
		rotatelogs.WithRotationTime(RotationTime),
		rotatelogs.WithMaxAge(MaxAge),
		rotatelogs.WithRotationSize(RotationSize),
	)
	if err != nil {
		panic(fmt.Sprintf("init error writer failed: %v", err))
	}

	debugWriter, err := rotatelogs.New(
		filepath.Join(c.Dir, "debug.%Y-%m-%d.log"),
		rotatelogs.WithLinkName(filepath.Join(c.Dir, "debug.log")),
		rotatelogs.WithRotationTime(RotationTime),
		rotatelogs.WithMaxAge(MaxAge),
		rotatelogs.WithRotationSize(RotationSize),
	)
	if err != nil {
		panic(fmt.Sprintf("init debug writer failed: %v", err))
	}

	// 创建 encoder
	encoder := encoder(c.Encode, c.TimeFormat)

	// 设置默认日志级别
	if c.Level == 0 {
		c.Level = int8(zapcore.DebugLevel)
	}

	// 创建不同级别的 core
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel && lvl >= zapcore.Level(c.Level)
	})

	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel && lvl >= zapcore.Level(c.Level)
	})

	debugLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel && lvl >= zapcore.Level(c.Level)
	})

	// 配置 core，同时输出到 stdout
	cores := []zapcore.Core{
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(errorWriter), errorLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(debugWriter), debugLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), debugLevel),
	}

	core := zapcore.NewTee(cores...)
	z = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	// 创建实例 Logger（用于返回，跳过0层）
	instanceLogger := zap.New(core, zap.AddCaller()) // 无 AddCallerSkip

	if c.TracerName != "" {
		tracing = true
		InitTracer(c.TracerName)
	}

	return &Logger{instanceLogger}
}

func encoder(et EncodeType, format string) zapcore.Encoder {
	if format == "" {
		format = YYYYMMDDHHMMSS
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "t"
	encoderConfig.LevelKey = "lvl"
	encoderConfig.NameKey = "log"
	encoderConfig.MessageKey = "m"
	encoderConfig.StacktraceKey = "s"
	encoderConfig.CallerKey = "c"
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(format)
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConfig.ConsoleSeparator = " "

	switch et {
	case JSONEncode:
		return zapcore.NewJSONEncoder(encoderConfig)
	case ConsoleEncode:
		return zapcore.NewConsoleEncoder(encoderConfig)
	default:
		return zapcore.NewJSONEncoder(encoderConfig)
	}
}

// Level 返回日志级别
func Level() zapcore.Level {
	return z.Level()
}

// With 返回带上下文字段的新 Logger
func With(fields ...zap.Field) *zap.Logger {
	return z.With(fields...)
}

func Debug(msg string, fields ...zap.Field) {
	z.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	z.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	z.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	z.Error(msg, fields...)
}

func DPanic(msg string, fields ...zap.Field) {
	z.DPanic(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	z.Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	z.Fatal(msg, fields...)
}

func Sync() error {
	if z != nil {
		return z.Sync()
	}
	return nil
}

func IsDebugEnable() bool {
	return z.Level() == zapcore.DebugLevel
}
