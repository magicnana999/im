package main

import (
	"context"
	"flag"
	"github.com/magicnana999/im/conf"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	_init()
}

func _init() {
	var file string
	flag.StringVar(&file, "conf", "conf/im-router.yaml", "config file path")
	flag.Parse()
	conf.LoadConfig(file)
}
