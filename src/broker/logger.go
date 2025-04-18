package broker

import (
	"github.com/magicnana999/im/pkg/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type EventZapField string

var (
	SrvLifecycle  EventZapField = "server lifecycle"
	ConnLifecycle EventZapField = "conn lifecycle"
	MsgTracking   EventZapField = "message tracking"
)

// Logger 用于Broker内部的日志
type Logger struct {
	*logger.Logger
	debugMode bool
}

func NewLogger(name string, debugMode bool) *Logger {
	return &Logger{
		Logger:    logger.Named(name),
		debugMode: debugMode,
	}
}

func (s *Logger) SrvInfo(msg string, ezf EventZapField, err error, field ...zap.Field) {
	s._log(zap.InfoLevel, msg, "", "", ezf, err, field...)
}

func (s *Logger) ConnDebug(msg, connDesc string, ezf EventZapField, err error, field ...zap.Field) {
	s._log(zap.DebugLevel, msg, connDesc, "", ezf, err, field...)
}

func (s *Logger) MsgDebug(msg, connDesc, messageID string, ezf EventZapField, err error, field ...zap.Field) {
	s._log(zap.DebugLevel, msg, connDesc, messageID, ezf, err, field...)
}

func (s *Logger) _log(level zapcore.Level, msg string, connDesc string, messageID string, ezf EventZapField, err error, fs ...zap.Field) {
	if s.Logger == nil {
		return
	}

	if level == zap.DebugLevel && (!s.debugMode || !s.Logger.IsDebugEnabled()) {
		return
	}

	fields := make([]zap.Field, 0)

	//if err == nil {
	//	fields = append(fields, zap.Bool(define.Status, true))
	//} else {
	//	fields = append(fields, zap.Bool(define.Status, false))
	//}

	if ezf != "" {
		fields = append(fields, zap.String("event", string(ezf)))
	}

	if connDesc != "" {
		fields = append(fields, zap.String("conn", connDesc))
	}

	if messageID != "" {
		fields = append(fields, zap.String("messageID", messageID))
	}

	if fs != nil && len(fs) > 0 {
		fields = append(fields, fs...)
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
		s.Logger.Error(msg, fields...)
	} else {
		s.Logger.Log(level, msg, fields...)
	}
}
