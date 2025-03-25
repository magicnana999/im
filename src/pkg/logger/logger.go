package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

const (
	YYYYMMDDHHMMSS = "2006-01-02 15:04:05"
)

type EncodeType string

const (
	ConsoleEncode EncodeType = "console"
	JSONEncode    EncodeType = "json"
)

var (
	Z       *zap.Logger // 改为 *zap.Logger
	tracing = false
)

type Config struct {
	File       string     `json:"file"`
	TracerName string     `json:"tracerName"`
	Level      int8       `json:"level"`
	Encode     EncodeType `json:"encode"`
	TimeFormat string     `json:"timeFormat"`
}

var defaultConfig = Config{
	File:       "./log/logger.log",
	TracerName: "",
	Level:      0,
	Encode:     JSONEncode,
	TimeFormat: YYYYMMDDHHMMSS,
}

func getDefaultConfig(c *Config) *Config {
	if c == nil {
		return &defaultConfig
	}

	if c.File == "" {
		c.File = "logger.log"
	}

	if c.Encode == "" {
		c.Encode = JSONEncode
	}

	if c.TimeFormat == "" {
		c.TimeFormat = YYYYMMDDHHMMSS
	}

	return c
}

// Init 返回 *zap.Logger
func Init(c *Config) *zap.Logger {
	c = getDefaultConfig(c)

	writeSyncer := Writer(c.File)
	encoder := Encoder(c.Encode, c.TimeFormat)

	if c.Level == 0 {
		c.Level = int8(zapcore.InfoLevel)
	}

	core := zapcore.NewCore(encoder, writeSyncer, zapcore.Level(c.Level))
	zp := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	Z = zp // 直接赋值给全局 Z

	if c.TracerName != "" {
		tracing = true
		InitTracer(c.TracerName)
	}

	return Z
}

// Encoder 保持不变
func Encoder(et EncodeType, format string) zapcore.Encoder {
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

// Writer 保持不变
func Writer(file string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   file,
		MaxSize:    10, // 10M
		MaxBackups: 5,  // 5个
		MaxAge:     30, // 最多30天
		Compress:   false,
	}
	return zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
		zapcore.AddSync(lumberJackLogger))
}

// Level 返回日志级别
func Level() zapcore.Level {
	return Z.Level()
}

// With 返回带上下文字段的新 Logger
func With(fields ...zap.Field) *zap.Logger {
	return Z.With(fields...)
}

// 基本日志方法，使用 zap.Field
func Debug(msg string, fields ...zap.Field) {
	Z.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	Z.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Z.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Z.Error(msg, fields...)
}

func DPanic(msg string, fields ...zap.Field) {
	Z.DPanic(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	Z.Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	Z.Fatal(msg, fields...)
}

// IsDebugEnable 检查是否启用 Debug 级别
func IsDebugEnable() bool {
	return Z.Level() == zapcore.DebugLevel
}
