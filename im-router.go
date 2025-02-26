package main

import (
	"context"
	"flag"
	"github.com/magicnana999/im/conf"
	"github.com/magicnana999/im/router"
	"sync"
)

func main() {

	ctx, _ := context.WithCancel(context.Background())

	var confFile string
	flag.StringVar(&confFile, "conf", "conf/im-router.yaml", "config file path")
	flag.Parse()

	conf.LoadConfig(confFile)
	router.Start(ctx)

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
