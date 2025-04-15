package logger

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	RotationTime = 24 * time.Hour      // 每天轮转
	MaxAge       = 30 * 24 * time.Hour // 保留30天
	RotationSize = 500 * 1024 * 1024
	//RotationSize = 1024 // 10MB 分割
)

const (
	t       = "t"
	level   = "lvl"
	name    = "log"
	message = "msg"
	stack   = "stk"
	caller  = "clr"
)

type EncodeType string

const (
	ConsoleEncode EncodeType = "console"
	JSONEncode    EncodeType = "json"
)

var (
	instance *Logger
	once     sync.Once
)

type Config struct {
	Dir        string        `json:"dir"` // 日志目录
	TracerName string        `json:"tracerName"`
	Level      zapcore.Level `json:"level"`
	Encode     EncodeType    `json:"encode"`
	TimeFormat string        `json:"timeFormat"`
}

var defaultConfig = Config{
	Dir:        "./logs",
	TracerName: "",
	Level:      0,
	Encode:     JSONEncode,
	TimeFormat: time.DateTime,
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
		c.TimeFormat = time.DateTime
	}

	return c
}

func Init(c *Config) (*Logger, error) {

	var err error
	once.Do(func() {

		c = getDefaultConfig(c)

		// 确保日志目录存在
		if err = os.MkdirAll(c.Dir, 0755); err != nil {
			return
		}

		writers := make([]*rotatelogs.RotateLogs, 0)
		// 创建不同级别的 writer
		infoWriter, e := rotatelogs.New(
			filepath.Join(c.Dir, "info.%Y-%m-%d.log"),
			rotatelogs.WithLinkName(filepath.Join(c.Dir, "info.log")),
			rotatelogs.WithRotationTime(RotationTime),
			rotatelogs.WithMaxAge(MaxAge),
			rotatelogs.WithRotationSize(RotationSize),
		)
		if e != nil {
			err = e
			return
		}

		errorWriter, e := rotatelogs.New(
			filepath.Join(c.Dir, "error.%Y-%m-%d.log"),
			rotatelogs.WithLinkName(filepath.Join(c.Dir, "error.log")),
			rotatelogs.WithRotationTime(RotationTime),
			rotatelogs.WithMaxAge(MaxAge),
			rotatelogs.WithRotationSize(RotationSize),
		)
		if e != nil {
			err = e
			return
		}

		debugWriter, e := rotatelogs.New(
			filepath.Join(c.Dir, "debug.%Y-%m-%d.log"),
			rotatelogs.WithLinkName(filepath.Join(c.Dir, "debug.log")),
			rotatelogs.WithRotationTime(RotationTime),
			rotatelogs.WithMaxAge(MaxAge),
			rotatelogs.WithRotationSize(RotationSize),
		)
		if e != nil {
			err = e
			return
		}

		writers = append(writers, infoWriter, debugWriter, errorWriter)

		// 创建 encoder
		encoding := encoder(c.Encode, c.TimeFormat)

		// 设置默认日志级别
		if c.Level == 0 {
			c.Level = zapcore.DebugLevel
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
			zapcore.NewCore(encoding, zapcore.AddSync(infoWriter), infoLevel),
			zapcore.NewCore(encoding, zapcore.AddSync(errorWriter), errorLevel),
			zapcore.NewCore(encoding, zapcore.AddSync(debugWriter), debugLevel),
			zapcore.NewCore(encoding, zapcore.AddSync(os.Stdout), debugLevel),
		}

		core := zapcore.NewTee(cores...)
		//z = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

		// 创建实例 Logger（用于返回，跳过0层）
		instanceLogger := zap.New(core, zap.AddCaller()) // 无 AddCallerSkip

		if c.TracerName != "" {
			InitTracer(c.TracerName)
		}

		instance = &Logger{Logger: instanceLogger, writers: writers}
	})

	return instance, err
}

func encoder(et EncodeType, format string) zapcore.Encoder {
	if format == "" {
		format = time.DateTime
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = t
	encoderConfig.LevelKey = level
	encoderConfig.StacktraceKey = stack
	encoderConfig.CallerKey = caller
	encoderConfig.NameKey = name
	encoderConfig.MessageKey = message
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

func Named(name string) *Logger {
	return instance.Named(name)
}

func Close() error {
	return instance.Close()
}
