package logger

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var (
	Z *zap.SugaredLogger
)

func InitLogger(logfile, tracerName string, level int8) {
	writeSyncer := Writer(logfile)
	encoder := Encoder()

	lvl := zapcore.Level(level)
	core := zapcore.NewCore(encoder, writeSyncer, lvl)
	zp := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	l := zp.Sugar()
	defer func() {
		if err := l.Sync(); err != nil {
			fmt.Printf("Error syncing log: %v\n", err)
		}
	}()

	Z = l

	if tracerName != "" {
		InitTracer(tracerName)
	}
}

func Encoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.LevelKey = "level"
	encoderConfig.NameKey = "logger"
	encoderConfig.MessageKey = "message"
	encoderConfig.StacktraceKey = "stack"
	encoderConfig.CallerKey = "caller" // 显示调用者（如果需要）
	encoderConfig.LineEnding = zapcore.DefaultLineEnding
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConfig.ConsoleSeparator = " "

	return zapcore.NewConsoleEncoder(encoderConfig)
}

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

func Debug(args ...interface{}) {
	Z.Debug(args...)
}

func Info(args ...interface{}) {
	Z.Info(args...)
}

func Warn(args ...interface{}) {
	Z.Warn(args)
}

func Error(args ...interface{}) {
	Z.Error(args...)
}

func Fatal(args ...interface{}) {
	Z.Fatal(args...)
}

func Debugf(template string, args ...interface{}) {
	Z.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	Z.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	Z.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	Z.Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	Z.Fatalf(template, args...)
}

func IsDebugEnable() bool {
	return Z.Level() == zapcore.DebugLevel
}
