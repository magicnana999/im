package logger

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Interface interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)

	Close() error
	IsDebugEnabled() bool
}

type Logger struct {
	*zap.Logger
	writers []*rotatelogs.RotateLogs
}

func (s *Logger) IsDebugEnabled() bool {
	return s.Level() == zapcore.DebugLevel
}

func (s *Logger) Named(name string) *Logger {
	return &Logger{Logger: s.Logger.Named(name)}
}

func (s *Logger) Close() error {

	for _, writer := range s.writers {
		if err := writer.Close(); err != nil {
			return err
		}
	}

	if err := s.Sync(); err != nil {
		return err
	}

	if tracer != nil {
		if err := ShutdownTracer(); err != nil {
			return err
		}
	}
	return nil
}
