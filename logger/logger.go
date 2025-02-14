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
	"sync"
)

type Log interface {
	Info(args ...any)
	Debug(args ...any)
	Error(args ...any)
	Warn(args ...any)
	Fatal(args ...any)
	InfoF(template string, args ...any)
	DebugF(template string, args ...any)
	ErrorF(template string, args ...any)
	WarnF(template string, args ...any)
	FatalF(template string, args ...any)
	Infof(template string, args ...any)
	Debugf(template string, args ...any)
	Errorf(template string, args ...any)
	Warnf(template string, args ...any)
	Fatalf(template string, args ...any)
	Printf(template string, args ...any)
	Level() string
}

type Tracer interface {
	Log
	NewSpan(ctx context.Context, name string) context.Context
	EndSpan(ctx context.Context)
	TraceID(ctx context.Context) string
	SpanID(ctx context.Context) string
}

type ZapLogger struct {
	zap         *zap.SugaredLogger
	threadLocal routine.ThreadLocal[context.Context]
	traceEnable bool
	level       string
}

var (
	DefaultLogger *ZapLogger
	lock          sync.Mutex
)

func InitLogger(logLevel string) *ZapLogger {
	lock.Lock()
	defer lock.Unlock()

	if DefaultLogger != nil {
		return DefaultLogger
	}

	var level zapcore.Level

	switch logLevel {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	default:
		level = zapcore.DebugLevel
	}

	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, level)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(3))
	//logger := zap.New(core, zap.AddCaller())
	log := logger.Sugar()
	defer func() {
		if err := log.Sync(); err != nil {
			fmt.Printf("Error syncing log: %v\n", err)
		}
	}()

	DefaultLogger = &ZapLogger{
		zap:         log,
		traceEnable: false,
		level:       logLevel,
	}

	return DefaultLogger

}

func getEncoder() zapcore.Encoder {
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

func (z *ZapLogger) NewSpan(ctx context.Context, name string) context.Context {
	ctx = tracer.NewSpan(ctx, name)
	z.threadLocal.Set(ctx)
	return ctx
}

func (z *ZapLogger) EndSpan(ctx context.Context) {
	tracer.EndSpan(ctx)
	z.threadLocal.Remove()
}

func (z *ZapLogger) TraceID(ctx context.Context) string {
	return tracer.TraceID(ctx)
}

func (z *ZapLogger) SpanID(ctx context.Context) string {
	return tracer.SpanID(ctx)
}

func (z *ZapLogger) log(level zapcore.Level, template string, args ...any) {

	if z.traceEnable && z.threadLocal.Get() != nil {

		ctx := z.threadLocal.Get().(context.Context)
		tid := z.TraceID(ctx)
		sid := z.SpanID(ctx)
		var param []any
		param = append(param, tid, sid)
		param = append(param, args...)
		z.zap.Logf(level, "%s %s "+template, tid, sid, args)
		return
	}

	z.zap.Logf(level, template, args...)

}

func (z *ZapLogger) InfoF(template string, args ...any) {
	z.log(zapcore.InfoLevel, template, args...)
}

func (z *ZapLogger) Info(args ...any) {
	z.log(zapcore.InfoLevel, "%s", args...)
}

func (z *ZapLogger) DebugF(template string, args ...any) {
	z.log(zapcore.DebugLevel, template, args...)
}

func (z *ZapLogger) Debug(args ...any) {
	z.log(zapcore.DebugLevel, "", args...)
}

func (z *ZapLogger) ErrorF(template string, args ...any) {
	z.log(zapcore.ErrorLevel, template, args...)
}

func (z *ZapLogger) Error(args ...any) {
	z.log(zapcore.ErrorLevel, "", args...)
}

func (z *ZapLogger) WarnF(template string, args ...any) {
	z.log(zapcore.ErrorLevel, template, args...)
}

func (z *ZapLogger) Warn(args ...any) {
	z.log(zapcore.ErrorLevel, "", args...)
}

func (z *ZapLogger) FatalF(template string, args ...any) {
	z.log(zapcore.FatalLevel, template, args...)
}

func (z *ZapLogger) Fatal(args ...any) {
	z.log(zapcore.FatalLevel, "", args...)
}

func (z *ZapLogger) Printf(format string, args ...any) {
	z.log(zapcore.InfoLevel, format, args)
}

func (z *ZapLogger) Debugf(format string, args ...any) {
	z.log(zapcore.DebugLevel, format, args...)
}

func (z *ZapLogger) Infof(format string, args ...any) {
	z.log(zapcore.InfoLevel, format, args...)
}

func (z *ZapLogger) Warnf(format string, args ...any) {
	z.log(zapcore.WarnLevel, format, args...)
}

func (z *ZapLogger) Errorf(format string, args ...any) {
	z.log(zapcore.ErrorLevel, format, args...)
}

func (z *ZapLogger) Fatalf(format string, args ...any) {
	z.log(zapcore.FatalLevel, format, args...)
}

func (z *ZapLogger) IsDebug() bool {
	return z.level == "debug"
}

func IsDebug() bool {
	return DefaultLogger.IsDebug()
}

func InfoF(template string, args ...any) {
	DefaultLogger.Infof(template, args...)
}

func Info(args ...any) {
	DefaultLogger.Info(args...)
}

func DebugF(template string, args ...any) {
	DefaultLogger.DebugF(template, args...)
}

func Debug(args ...any) {
	DefaultLogger.Debug(args...)
}

func ErrorF(template string, args ...any) {
	DefaultLogger.ErrorF(template, args...)
}

func Error(args ...any) {
	DefaultLogger.Error(args...)
}

func WarnF(template string, args ...any) {
	DefaultLogger.WarnF(template, args...)
}

func Warn(args ...any) {
	DefaultLogger.Warn(args...)
}

func FatalF(template string, args ...any) {
	DefaultLogger.Fatalf(template, args...)
}

func Fatal(args ...any) {
	DefaultLogger.Fatal(args...)
}
