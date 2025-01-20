package main

import (
	"context"
	"github.com/magicnana999/im/broker/core"
	"time"
)

func main() {
	ctx := context.Background()
	option := core.Option{
		Addr:         "localhost:8080",
		TickDuration: time.Second * 30,
	}
	core.Start(ctx, option)

}
