package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/magicnana999/im/broker"
	"github.com/magicnana999/im/broker/cmd_service"
	"github.com/magicnana999/im/broker/handler"
	"github.com/magicnana999/im/broker/holder"
	"github.com/magicnana999/im/global"
	"github.com/magicnana999/im/infra"
	"github.com/magicnana999/im/pkg/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger.Init(nil)
	defer logger.Close()

	var confFile string
	flag.StringVar(&confFile, "conf", "conf/im-broker.yaml", "config file path")
	flag.Parse()

	f := func() (*global.Config, error) {
		return global.Load(confFile)
	}

	log := logger.Named("main")
	app := fx.New(
		fx.NopLogger,
		fx.Provide(
			f,
			infra.NewRedisClient,
			infra.NewGorm,
			infra.NewKafkaProducer,
			infra.NewEtcdRegistry,
			infra.NewEtcdResolver,
			infra.NewBusinessClient,
			infra.NewRouterClient,
			infra.NewBrokerClientResolver,
			holder.NewBrokerHolder,
			holder.NewUserHolder,
			broker.NewHeartbeatServer,
			broker.NewMessageRetryServer,
			broker.NewMessageSendServer,
			cmd_service.NewUserService,
			handler.NewCommandHandler,
			handler.NewMessageHandler,
			broker.NewRpcBrokerServer,
			broker.NewTcpServer,
		),
		fx.Invoke(func(tcp *broker.TcpServer, rpc *broker.RpcBrokerServer, delivery *broker.MessageSendServer) {
			go func() {
				fmt.Println("sdfsdfsdfsdfsfdsf ")
			}()
		}),
	)

	// 捕获信号
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// 启动 Fx
	startCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.Start(startCtx); err != nil {
		log.Fatal("Failed to start app", zap.Error(err))
	}

	<-sigs
	log.Info("shutdown...")

	// 停止 Fx
	stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.Stop(stopCtx); err != nil {
		log.Error("Failed to stop app", zap.Error(err))
	}

	log.Info("Shutdown complete")

}
