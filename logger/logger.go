package logger

import (
	"context"
	"fmt"
	"github.com/magicnana999/im/tracer"
	"github.com/natefinch/lumberjack"
	"github.com/timandy/routine"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

const (
	TraceEnable bool = true
)

var (
	Log    *zap.SugaredLogger
	m      = routine.NewThreadLocal[context.Context]()
	Logger LoggerAdapter
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
	Logger = LoggerAdapter{}
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
	ctx = tracer.NewSpan(ctx, name)
	m.Set(ctx)
	return ctx
}

func EndSpan(ctx context.Context) {
	tracer.EndSpan(ctx)
	m.Remove()
}

func TraceID(ctx context.Context) string {
	return tracer.TraceID(ctx)
}

func SpanID(ctx context.Context) string {
	return tracer.SpanID(ctx)
}
func _log(lvl zapcore.Level, template string, args ...interface{}) {
	trace := false

	if m.Get() == nil {
		trace = false
	} else if TraceEnable && m.Get() != nil {
		trace = true
	}

	if trace {
		ctx := m.Get().(context.Context)
		tid := TraceID(ctx)
		sid := SpanID(ctx)
		var param []interface{}
		param = append(param, tid, sid)
		param = append(param, args...)
		Log.Logf(lvl, "%s %s "+template, param...)
	} else {
		Log.Logf(lvl, template, args...)
	}

}

func InfoF(template string, args ...interface{}) {
	_log(zapcore.InfoLevel, template, args...)
}

func Info(args ...interface{}) {
	_log(zapcore.InfoLevel, "%s", args...)
}

func DebugF(template string, args ...interface{}) {
	_log(zapcore.DebugLevel, template, args...)
}

func Debug(args ...interface{}) {
	_log(zapcore.DebugLevel, "", args...)
}

func ErrorF(template string, args ...interface{}) {
	_log(zapcore.ErrorLevel, template, args...)
}

func Error(args ...interface{}) {
	_log(zapcore.ErrorLevel, "", args...)
}

func WarnF(template string, args ...interface{}) {
	_log(zapcore.ErrorLevel, template, args...)
}

func Warn(args ...interface{}) {
	_log(zapcore.ErrorLevel, "", args...)
}

func FatalF(template string, args ...interface{}) {
	_log(zapcore.FatalLevel, template, args...)
}

func Fatal(args ...interface{}) {
	_log(zapcore.FatalLevel, "", args...)
}

type LoggerAdapter struct {
}

func (l LoggerAdapter) Printf(format string, args ...any) {
	InfoF(format, args)
}

func (l LoggerAdapter) Debugf(format string, args ...any) {
	DebugF(format, args...)
}

func (l LoggerAdapter) Infof(format string, args ...any) {
	InfoF(format, args...)
}

func (l LoggerAdapter) Warnf(format string, args ...any) {
	WarnF(format, args...)
}

func (l LoggerAdapter) Errorf(format string, args ...any) {
	ErrorF(format, args...)
}

func (l LoggerAdapter) Fatalf(format string, args ...any) {
	FatalF(format, args...)
}
