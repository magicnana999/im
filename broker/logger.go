package broker

import (
	"context"
	"fmt"
	"github.com/magicnana999/im/conf"
	"github.com/magicnana999/im/logger"
	"github.com/panjf2000/gnet/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func initLogger() {
	writeSyncer := logger.Writer("log/im-broker.log")
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

	logger.InitTracer("im-broker")

}

func traffic(ctx context.Context, c gnet.Conn, template string, args ...any) (string, []any) {

	clientAddr := ""
	ucLabel := ""
	traceId := logger.TraceID(ctx)

	uc, _ := currentUserFromConn(c)

	if uc != nil {
		clientAddr = uc.ClientAddr
		ucLabel = uc.Label()
	}

	param := make([]any, 0)
	param = append(param, traceId)
	param = append(param, clientAddr)
	param = append(param, ucLabel)
	param = append(param, args...)

	return "%s [%s#%s] " + template, param
}
