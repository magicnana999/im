package logger

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// RotationTime defines the duration after which log files are rotated (24 hours).
const RotationTime = 24 * time.Hour

// MaxAge defines the maximum duration to retain log files (30 days).
const MaxAge = 30 * 24 * time.Hour

// RotationSize defines the maximum file size before rotation (500 MB).
const RotationSize = 500 * 1024 * 1024

// TimeKey is the key for the timestamp field in logs.
const TimeKey = "time"

// LevelKey is the key for the log level field in logs.
const LevelKey = "level"

// NameKey is the key for the logger name field in logs.
const NameKey = "name"

// MessageKey is the key for the message field in logs.
const MessageKey = "message"

// StacktraceKey is the key for the stacktrace field in logs.
const StacktraceKey = "stacktrace"

// CallerKey is the key for the caller field in logs.
const CallerKey = "caller"

// EncodeType defines the encoding format for logs.
type EncodeType string

// Constants for log encoding types.
const (
	ConsoleEncode EncodeType = "console" // ConsoleEncode formats logs in a human-readable console format.
	JSONEncode    EncodeType = "json"    // JSONEncode formats logs in JSON format.
)

// Logger wraps a zap logger with file rotation writers.
type Logger struct {
	*zap.Logger                          // Embedded zap logger for logging operations.
	writers     []*rotatelogs.RotateLogs // writers handle log file rotation.
}

// Config defines the configuration for initializing a Logger.
type Config struct {
	Dir        string        // Dir is the directory for log files.
	TracerName string        // TracerName is the name of the tracer (optional).
	Level      zapcore.Level // Level is the minimum log level to record.
	Encode     EncodeType    // Encode specifies the log encoding format (JSON or console).
	TimeFormat string        // TimeFormat is the format for timestamp fields.
}

// defaultConfig provides default values for Logger configuration.
var defaultConfig = &Config{
	Dir:        "./logs",
	TracerName: "",
	Level:      zapcore.DebugLevel,
	Encode:     JSONEncode,
	TimeFormat: time.DateTime,
}

// instance holds the singleton Logger instance.
var instance *Logger

// once ensures the Logger is initialized only once.
var once sync.Once

// Init initializes the singleton Logger with the specified configuration.
// It creates log files with rotation (by time and size) and supports JSON or console encoding.
// It returns the initialized Logger and an error if initialization fails (e.g., directory creation or writer setup fails).
// Subsequent calls return the same instance without reinitializing.
func Init(c *Config) (*Logger, error) {
	var err error
	once.Do(func() {
		c = getDefaultConfig(c)

		// Ensure log directory exists.
		if err = os.MkdirAll(c.Dir, 0755); err != nil {
			return
		}

		// Create writers for different log levels.
		writers := make([]*rotatelogs.RotateLogs, 0)
		for _, logType := range []string{"info", "error", "debug"} {
			writer, e := rotatelogs.New(
				filepath.Join(c.Dir, fmt.Sprintf("%s.%%Y-%%m-%%d.log", logType)),
				rotatelogs.WithLinkName(filepath.Join(c.Dir, fmt.Sprintf("%s.log", logType))),
				rotatelogs.WithRotationTime(RotationTime),
				rotatelogs.WithMaxAge(MaxAge),
				rotatelogs.WithRotationSize(RotationSize),
			)
			if e != nil {
				err = fmt.Errorf("failed to create %s writer: %w", logType, e)
				return
			}
			writers = append(writers, writer)
		}

		// Create encoder.
		encoding := encoder(c.Encode, c.TimeFormat)

		// Configure log levels.
		infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl == zapcore.InfoLevel && lvl >= c.Level
		})
		errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel && lvl >= c.Level
		})
		debugLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.DebugLevel && lvl >= c.Level
		})

		// Create cores for different outputs.
		cores := []zapcore.Core{
			zapcore.NewCore(encoding, zapcore.AddSync(writers[0]), infoLevel),  // info
			zapcore.NewCore(encoding, zapcore.AddSync(writers[1]), errorLevel), // error
			zapcore.NewCore(encoding, zapcore.AddSync(writers[2]), debugLevel), // debug
			zapcore.NewCore(encoding, zapcore.AddSync(os.Stdout), debugLevel),  // stdout
		}

		// Combine cores.
		core := zapcore.NewTee(cores...)

		// Create zap logger.
		zapLogger := zap.New(core, zap.AddCaller())

		// Initialize tracer if specified.
		if c.TracerName != "" {
			InitTracer(c.TracerName)
		}

		instance = &Logger{
			Logger:  zapLogger,
			writers: writers,
		}
	})

	if err != nil {
		return nil, err
	}
	if instance == nil {
		return nil, errors.New("logger not initialized")
	}
	return instance, nil
}

// getDefaultConfig returns a Config with default values applied to any unset fields.
// It ensures Dir, Encode, and TimeFormat have valid values.
func getDefaultConfig(config *Config) *Config {

	c := &Config{
		Dir:        "./logs",
		TracerName: "",
		Level:      zapcore.DebugLevel,
		Encode:     JSONEncode,
		TimeFormat: time.DateTime,
	}

	if config != nil {
		*c = *config
	}

	if c.Dir == "" {
		c.Dir = defaultConfig.Dir
	}
	if c.Encode == "" {
		c.Encode = defaultConfig.Encode
	}
	if c.TimeFormat == "" {
		c.TimeFormat = defaultConfig.TimeFormat
	}
	if c.Level == 0 {
		c.Level = defaultConfig.Level
	}
	return c
}

// encoder creates a zap encoder based on the specified encoding type and time format.
// It returns a JSON encoder if the encoding type is invalid or unspecified.
func encoder(et EncodeType, format string) zapcore.Encoder {
	if format == "" {
		format = time.DateTime
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = TimeKey
	encoderConfig.LevelKey = LevelKey
	encoderConfig.StacktraceKey = StacktraceKey
	encoderConfig.CallerKey = CallerKey
	encoderConfig.NameKey = NameKey
	encoderConfig.MessageKey = MessageKey
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

// Named returns a new Logger with the specified name.
// It returns nil if the Logger instance is not initialized.
func Named(name string) *Logger {
	if instance == nil {
		return nil
	}
	return &Logger{Logger: instance.Logger.Named(name)}
}

func NamedAndAddSkip(name string, skip int) *Logger {
	if instance == nil {
		return nil
	}
	return &Logger{Logger: instance.Logger.WithOptions(zap.AddCallerSkip(skip)).Named(name)}
}

func Close() {
	if instance != nil {
		instance.Close()
	}
}
func (l *Logger) IsDebugEnabled() bool {
	return l.Level() == zapcore.DebugLevel
}

func (l *Logger) Named(name string) *Logger {
	return &Logger{Logger: l.Logger.Named(name)}
}

// Close shuts down the Logger and closes all associated log writers.
// It returns an error if any writer fails to close.
func (l *Logger) Close() error {
	if l == nil || l.Logger == nil {
		return nil
	}
	l.Logger.Sync()

	var errs []error
	for _, w := range l.writers {
		if err := w.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("failed to close writers: %v", errs)
	}
	return nil
}

type Interface interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)

	Close() error
	IsDebugEnabled() bool
}
