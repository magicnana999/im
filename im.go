package main

import (
	"context"
	"fmt"
	broker "github.com/magicnana999/im/broker/brokerstate"
	"github.com/magicnana999/im/logger"
	"github.com/petermattis/goid"
	"sync"
)

func main() {

	fmt.Println(goid.Get())
	ctx := context.Background()
	logger.Info(ctx, "Start im broker ...")

	var wait sync.WaitGroup
	wait.Add(1)
	go startBroker(ctx, &wait)

	wait.Wait()
	logger.Info(ctx, "Im broker is stop")
}

func startBroker(ctx context.Context, wg *sync.WaitGroup) {
	broker := broker.NewBroker(ctx, wg, "192.168.1.1", 0)
	broker.Start()
	logger.Info(ctx, "Start im broker success")
}
