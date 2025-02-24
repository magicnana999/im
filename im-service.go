package main

import (
	"flag"
	"github.com/magicnana999/im/conf"
	"github.com/magicnana999/im/service"
)

func main() {

	//root, cancel := context.WithCancel(context.Background())

	var confFile string
	flag.StringVar(&confFile, "conf", "conf/im-service.yaml", "config file path")
	flag.Parse()

	conf.LoadConfig(confFile)

	service.Start()
}
