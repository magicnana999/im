package logger

import (
	"context"
	"fmt"
	"github.com/natefinch/lumberjack"
	"github.com/timandy/routine"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var (
	ThreadLocal routine.ThreadLocal[context.Context] = routine.NewThreadLocal[context.Context]()
	Log         *zap.SugaredLogger
)

var (
	PD                   *trace.TracerProvider
	loggerSpanContextKey string = "trace_current_span"
)

func init() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	Log = logger.Sugar()
	defer Log.Sync()
	InitTracer()
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
		Filename:   "logs/upc.log",
		MaxSize:    10, // 10M
		MaxBackups: 5,  // 5个
		MaxAge:     30, // 最多30天
		Compress:   false,
	}
	return zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
		zapcore.AddSync(lumberJackLogger))
}

func InitTracer() error {

	otel.SetTextMapPropagator(b3.New())

	f, _ := os.Create("trace.txt")

	exp, _ := stdouttrace.New(
		stdouttrace.WithWriter(f),
		stdouttrace.WithPrettyPrint(),
		stdouttrace.WithoutTimestamps(),
	)

	PD = trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithSampler(trace.NeverSample()),
	)

	otel.SetTracerProvider(PD)
	return nil
}

func CurrentSpan(ctx context.Context, spanName string) context.Context {
	started := ThreadLocal.Get()

	if started == nil {
		started = StartSpan(ctx, spanName)
		if started != nil {
			ThreadLocal.Set(started)
		}
	}
	return started
}

func StartSpan(ctx context.Context, spanName string) context.Context {
	spanCtx, span := PD.Tracer("").Start(ctx, spanName)
	startCtx := context.WithValue(spanCtx, loggerSpanContextKey, span)
	defer span.End()
	return startCtx
}

func TraceID(ctx context.Context) string {
	span, ok := ctx.Value(loggerSpanContextKey).(oteltrace.Span)
	if !ok {
		return ""
	}
	return span.SpanContext().TraceID().String()
}

func SpanID(ctx context.Context) string {
	span, ok := ctx.Value(loggerSpanContextKey).(oteltrace.Span)
	if !ok {
		return ""
	}
	return span.SpanContext().SpanID().String()
}

func _log(lvl zapcore.Level, template string, args ...interface{}) {
	ctx := CurrentSpan(context.Background(), "hello world")
	tid := TraceID(ctx)
	sid := SpanID(ctx)
	fmt.Println(tid)
	fmt.Println(sid)
	Log.Log(zapcore.InfoLevel, "%s %s "+template, append(args, tid, sid))
}

func Infof(template string, args ...interface{}) {
	_log(zapcore.InfoLevel, template, args...)
}

func Info(args ...interface{}) {
	_log(zapcore.InfoLevel, "", args)
}

func DebugF(template string, args ...interface{}) {
	_log(zapcore.DebugLevel, template, args...)
}

func Debug(args ...interface{}) {
	_log(zapcore.DebugLevel, "", args)
}

func ErrorF(template string, args ...interface{}) {
	_log(zapcore.ErrorLevel, template, args...)
}

func Error(args ...interface{}) {
	_log(zapcore.ErrorLevel, "", args)
}

func WarnF(template string, args ...interface{}) {
	_log(zapcore.ErrorLevel, template, args...)
}

func Warn(args ...interface{}) {
	_log(zapcore.ErrorLevel, "", args)
}

//func Debug(args ...interface{}) {
//	log.Debug(args...)
//}
//func Debugf(template string, args ...interface{}) {
//	log.Debugf(template, args...)
//}
//
//func Info(args ...interface{}) {
//	traceId, uri, err := _traceId()
//	if err != nil {
//		log.Info(args...)
//	}
//	log.Infof("%s %s %s", traceId, uri, args)
//}
//func Infof(template string, args ...interface{}) {
//
//	traceId, uri, err := _traceId()
//	if err != nil {
//		log.Infof(template, args)
//	}
//	log.Infof("%s %s "+template, traceId, uri, args)
//
//}
//
//func Warn(args ...interface{}) {
//	log.Warn(args...)
//}
//func Warnf(template string, args ...interface{}) {
//	log.Warnf(template, args...)
//}
//func Error(args ...interface{}) {
//	traceId, uri, err := _traceId()
//	if err != nil {
//		log.Error(args)
//	}
//
//	log.Errorf("%s %s %s", traceId, uri, args)
//}
//func Errorf(template string, args ...interface{}) {
//
//	traceId, uri, err := _traceId()
//	if err != nil {
//		log.Errorf("%s %s %s", traceId, uri, args)
//	}
//
//	log.Errorf("%s %s "+template, traceId, uri, args)
//}
//func DPanic(args ...interface{}) {
//	log.DPanic(args...)
//}
//func DPanicf(template string, args ...interface{}) {
//	log.DPanicf(template, args...)
//}
//func Panic(args ...interface{}) {
//	log.Panic(args...)
//}
//func Panicf(template string, args ...interface{}) {
//	log.Panicf(template, args...)
//}
//func Fatal(args ...interface{}) {
//	log.Fatal(args...)
//}
//func Fatalf(template string, args ...interface{}) {
//	log.Fatalf(template, args...)
//}
