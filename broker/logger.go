package broker

import (
	"fmt"
	"github.com/magicnana999/im/api/kitex_gen/api"
	"github.com/magicnana999/im/pkg/jsonext"
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
}

func NewLogger(name string) *Logger {
	return &Logger{
		Logger: logger.NameWithOptions(name, zap.AddCallerSkip(2)),
	}
}

func (s *Logger) Debugf(format string, args ...any) {
	s.Debug(fmt.Sprintf(format, args...))
}

func (s *Logger) Infof(format string, args ...any) {
	s.Info(fmt.Sprintf(format, args...))
}

func (s *Logger) Warnf(format string, args ...any) {
	s.Warn(fmt.Sprintf(format, args...))
}

func (s *Logger) Errorf(format string, args ...any) {
	s.Error(fmt.Sprintf(format, args...))
}

func (s *Logger) Fatalf(format string, args ...any) {
	s.Fatal(fmt.Sprintf(format, args...))
}

func (s *Logger) Printf(format string, args ...any) {
	if !s.Logger.IsDebugEnabled() {
		return
	}
	s.Debug(fmt.Sprintf(format, args))
}

func (s *Logger) SrvInfo(msg string, ezf EventZapField, err error, field ...zap.Field) {
	s._log(zap.InfoLevel, msg, "", "", nil, ezf, err, field...)
}

func (s *Logger) ConnDebug(msg, connDesc string, ezf EventZapField, err error, field ...zap.Field) {
	s._log(zap.DebugLevel, msg, connDesc, "", nil, ezf, err, field...)
}

func (s *Logger) PktDebug(msg, connDesc, messageID string, packet *api.Packet, ezf EventZapField, err error, field ...zap.Field) {
	s._log(zap.DebugLevel, msg, connDesc, messageID, packet, ezf, err, field...)
}

// _log 记录结构化日志，支持不同级别和事件类型。
// level: 日志级别（Debug、Info 等）。
// msg: 日志消息。
// connDesc: 连接描述。
// packetID: 数据包 ID。
// packet: 数据包对象（序列化为字符串）。
// ezf: 事件类型。
// err: 错误信息。
// fs: 附加字段。
func (s *Logger) _log(level zapcore.Level, msg string, connDesc string, packetID string, packet *api.Packet, ezf EventZapField, err error, fs ...zap.Field) {
	if s.Logger == nil {
		return
	}

	if level == zap.DebugLevel && !s.Logger.IsDebugEnabled() {
		return
	}

	fields := make([]zap.Field, 0, 7)

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

	if packet != nil {
		fields = append(fields, zap.String("packet", string(jsonext.PbMarshalNoErr(packet))))
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
