package logger

import (
	"context"
	"fmt"
	"github.com/magicnana999/im/tracer"
	"github.com/natefinch/lumberjack"
	//"github.com/timandy/routine"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var (
	Log *zap.SugaredLogger
)

func init() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	Log = logger.Sugar()
	defer func() {
		if err := Log.Sync(); err != nil {
			fmt.Printf("Error syncing log: %v\n", err)
		}
	}()
}

func demo() {
	ctx := context.Background()
	ctx = NewSpan(ctx, "root")
	Info(ctx, "haha1")
	Info(ctx, "haha2")

	ctx = NewSpan(ctx, "sub")
	Info(ctx, "haha3")

}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.LevelKey = "level"
	encoderConfig.NameKey = "logger"
	encoderConfig.MessageKey = "message"
	encoderConfig.StacktraceKey = "stack"
	encoderConfig.CallerKey = "" // 显示调用者（如果需要）
	encoderConfig.LineEnding = zapcore.DefaultLineEnding
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConfig.ConsoleSeparator = " "

	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "logs/im.log",
		MaxSize:    10, // 10M
		MaxBackups: 5,  // 5个
		MaxAge:     30, // 最多30天
		Compress:   false,
	}
	return zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
		zapcore.AddSync(lumberJackLogger))
}

func NewSpan(ctx context.Context, name string) context.Context {
	return tracer.NewSpan(ctx, name)
}

func EndSpan(ctx context.Context) {
	tracer.EndSpan(ctx)
}

func TraceID(ctx context.Context) string {
	return tracer.TraceID(ctx)
}

func SpanID(ctx context.Context) string {
	return tracer.SpanID(ctx)
}
func _log(ctx context.Context, lvl zapcore.Level, template string, args ...interface{}) {
	tid := TraceID(ctx)
	sid := SpanID(ctx)
	var param []interface{}
	param = append(param, tid, sid)
	param = append(param, args...)
	Log.Logf(lvl, "%s %s "+template, param...)
}

func InfoF(ctx context.Context, template string, args ...interface{}) {
	_log(ctx, zapcore.InfoLevel, template, args...)
}

func Info(ctx context.Context, args ...interface{}) {
	_log(ctx, zapcore.InfoLevel, "%s", args...)
}

func DebugF(ctx context.Context, template string, args ...interface{}) {
	_log(ctx, zapcore.DebugLevel, template, args...)
}

func Debug(ctx context.Context, args ...interface{}) {
	_log(ctx, zapcore.DebugLevel, "", args...)
}

func ErrorF(ctx context.Context, template string, args ...interface{}) {
	_log(ctx, zapcore.ErrorLevel, template, args...)
}

func Error(ctx context.Context, args ...interface{}) {
	_log(ctx, zapcore.ErrorLevel, "", args...)
}

func WarnF(ctx context.Context, template string, args ...interface{}) {
	_log(ctx, zapcore.ErrorLevel, template, args...)
}

func Warn(ctx context.Context, args ...interface{}) {
	_log(ctx, zapcore.ErrorLevel, "", args...)
}

func FatalF(ctx context.Context, template string, args ...interface{}) {
	_log(ctx, zapcore.FatalLevel, template, args...)
}

func Fatal(ctx context.Context, args ...interface{}) {
	_log(ctx, zapcore.FatalLevel, "", args...)
}
