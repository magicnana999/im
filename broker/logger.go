package broker

import (
	"github.com/magicnana999/im/pkg/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type EventZapField string

var (
	SrvLifecycle   EventZapField = "server lifecycle"
	ConnLifecycle  EventZapField = "conn lifecycle"
	PacketTracking EventZapField = "packet tracking"
)

// Logger 用于Broker内部的日志
type Logger struct {
	*logger.Logger
	debugMode bool
}

func NewLogger(name string, debugMode bool) *Logger {
	return &Logger{
		Logger:    logger.NameWithOptions(name, zap.AddCallerSkip(2)),
		debugMode: debugMode,
	}
}

func (s *Logger) SrvInfo(msg string, ezf EventZapField, err error, field ...zap.Field) {
	s._log(zap.InfoLevel, msg, "", "", "", ezf, err, field...)
}

func (s *Logger) ConnDebug(msg, connDesc string, ezf EventZapField, err error, field ...zap.Field) {
	s._log(zap.DebugLevel, msg, connDesc, "", "", ezf, err, field...)
}

func (s *Logger) PktDebug(msg, connDesc, messageID, packet string, ezf EventZapField, err error, field ...zap.Field) {
	s._log(zap.DebugLevel, msg, connDesc, messageID, packet, ezf, err, field...)
}

func (s *Logger) _log(level zapcore.Level, msg string, connDesc string, packetID string, packet string, ezf EventZapField, err error, fs ...zap.Field) {
	if s.Logger == nil {
		return
	}

	if level == zap.DebugLevel && (!s.debugMode || !s.Logger.IsDebugEnabled()) {
		return
	}

	fields := make([]zap.Field, 0)

	if ezf == ConnLifecycle || ezf == PacketTracking {
		if err == nil {
			fields = append(fields, zap.Bool("status", true))
		} else {
			fields = append(fields, zap.Bool("status", false))
		}
	}

	if ezf != "" {
		fields = append(fields, zap.String("event", string(ezf)))
	}

	if connDesc != "" {
		fields = append(fields, zap.String("conn", connDesc))
	}

	if packetID != "" {
		fields = append(fields, zap.String("packetID", packetID))
	}

	if packet != "" {
		fields = append(fields, zap.String("packet", packet))
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
