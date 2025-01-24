package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/magicnana999/im/broker/core"
	"github.com/magicnana999/im/logger"
	"github.com/magicnana999/im/util/ip"
	"github.com/magicnana999/im/util/json"
	"github.com/magicnana999/im/util/str"
	"github.com/timandy/routine"
	"strconv"
	time "time"
)

func main() {

	goid := routine.Goid()

	root := context.Background()
	option := parseFlag()
	logger.InfoF("[%d] Start broker instance: %s", goid, json.IgnoreErrorMarshal(option))

	core.Start(root, option)

}

func parseFlag() *core.Option {
	var name string
	var interval string
	flag.StringVar(&name, "name", "", "The name of the broker instance")
	flag.StringVar(&interval, "interval", "30", "Ticking interval of broker instance")
	flag.Parse()

	if str.IsBlank(&name) {
		ipAddress, e := ip.GetLocalIP()
		if str.IsBlank(&ipAddress) || e != nil {
			logger.FatalF("Failed to get local IP address: %v", e)
		}
		name = fmt.Sprintf("%s:%s", ipAddress, core.DefaultPort)
	}

	if str.IsBlank(&name) {
		logger.Fatal("Could not found the name of broker instance")
	}

	if str.IsBlank(&interval) {
		logger.Fatal("Could not found the interval of broker instance")
	}

	i, e := strconv.Atoi(interval)
	if e != nil {
		logger.FatalF("Failed to parse interval of broker instance: %v", e)
	}

	option := &core.Option{
		Name:     name,
		Interval: time.Duration(i) * time.Second,
	}

	return option
}
