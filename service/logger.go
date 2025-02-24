package service

import (
	"fmt"
	"github.com/magicnana999/im/conf"
	"github.com/magicnana999/im/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func initLogger() {
	writeSyncer := logger.Writer("log/im-service.log")
	encoder := logger.Encoder()

	level := zapcore.Level(conf.Global.Logger.Level)
	core := zapcore.NewCore(encoder, writeSyncer, level)
	zp := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	l := zp.Sugar()
	defer func() {
		if err := l.Sync(); err != nil {
			fmt.Printf("Error syncing log: %v\n", err)
		}
	}()

	logger.Z = l

	logger.InitTracer("im-service")

}
