package main

import (
	"context"
	"flag"
	"github.com/magicnana999/im/broker"
	"github.com/magicnana999/im/conf"
)

func main() {

	root, cancel := context.WithCancel(context.Background())

	var confFile string
	flag.StringVar(&confFile, "conf", "conf/im-broker.yaml", "config file path")
	flag.Parse()

	conf.LoadConfig(confFile)
	broker.Start(root, cancel)

}
