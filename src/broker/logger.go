package broker

import (
	"github.com/magicnana999/im/define"
	"github.com/magicnana999/im/pkg/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps a zap logger with additional debug mode control.
// It provides methods to log messages with consistent fields for connection, operation, and message ID.
type Logger struct {
	*logger.Logger      // Embedded zap logger for logging operations.
	debugMode      bool // debugMode enables debug-level logging when true.
}

// NewLogger creates a new Logger instance with the specified name and debug mode.
// The name is used to identify the logger in logs, and debugMode controls whether debug-level logs are enabled.
// It returns a configured Logger instance.
func NewLogger(name string, debugMode bool) *Logger {
	return &Logger{
		Logger:    logger.Named(name),
		debugMode: debugMode,
	}
}

// Print print info-level message with operation
// if err is non-nil, it logs an error;otherwise it logs a info
func (s *Logger) Print(msg string, operation string, err error) {
	s.log(zap.InfoLevel, msg, "", operation, "", err)
}

// Debug logs a debug-level message with connection, operation, and message ID fields.
// It logs only if debugMode is enabled and the logger's debug level is active.
// If err is non-nil, it logs an error with status false; otherwise, it logs a debug message with status true.
func (s *Logger) Debug(msg string, connDesc string, operation string, messageID string, err error) {
	s.log(zap.DebugLevel, msg, connDesc, operation, messageID, err)
}

// Info logs an info-level message with connection, operation, and message ID fields.
// If err is non-nil, it logs an error with status false; otherwise, it logs an info message with status true.
func (s *Logger) Info(msg string, connDesc string, operation string, messageID string, err error) {
	s.log(zap.InfoLevel, msg, connDesc, operation, messageID, err)
}

// log is a helper method to log messages at the specified level with consistent fields.
// It handles both success (status true) and error (status false) cases.
func (s *Logger) log(level zapcore.Level, msg string, connDesc string, operation string, messageID string, err error) {
	if s.Logger == nil {
		return // Safety check, though NewLogger ensures Logger is non-nil.
	}

	// Skip debug logs if debug mode is disabled or level is not enabled.
	if level == zap.DebugLevel && (!s.debugMode || !s.Logger.IsDebugEnabled()) {
		return
	}

	fields := make([]zap.Field, 0)

	if connDesc != "" {
		fields = append(fields, zap.String("conn", connDesc))
	}

	if operation != "" {
		fields = append(fields, zap.String("operation", operation))
	}
	if messageID != "" {
		fields = append(fields, zap.String(define.MessageId, messageID))
	}

	if err != nil {
		fields = append(fields, zap.Bool(define.Status, false), zap.Error(err))
		s.Logger.Error(msg, fields...)
	} else {
		fields = append(fields, zap.Bool(define.Status, true))
		s.Logger.Log(level, msg, fields...)
	}
}
