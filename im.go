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
	"strconv"
	"time"
)

func main() {

	root := context.Background()
	option := parseFlag()
	logger.InfoF("BrokerOption:%s", json.IgnoreErrorMarshal(option))

	core.Start(root, option)

}

func parseFlag() *core.Option {
	var name string
	var interval string
	var heartbeatInterval string
	flag.StringVar(&name, "name", "", "The name of the broker instance")
	flag.StringVar(&interval, "interval", "30", "Ticking interval of broker instance")
	flag.StringVar(&heartbeatInterval, "heartbeatInterval", "30", "Heartbeat interval")
	flag.Parse()

	if str.IsBlank(name) {
		ipAddress, e := ip.GetLocalIP()
		if str.IsBlank(ipAddress) || e != nil {
			logger.FatalF("Failed to get local IP address: %v", e)
		}
		name = fmt.Sprintf("%s:%s", ipAddress, core.DefaultPort)
	}

	if str.IsBlank(name) {
		logger.Fatal("Could not found the name of broker instance")
	}

	if str.IsBlank(interval) {
		logger.Fatal("Could not found the interval of broker instance")
	}

	i, e := strconv.Atoi(interval)
	if e != nil {
		logger.FatalF("Failed to parse interval of broker instance: %v", e)
	}

	if str.IsBlank(heartbeatInterval) {
		logger.Fatal("Could not found the heartbeat interval")
	}

	ih, eh := strconv.Atoi(heartbeatInterval)
	if eh != nil {
		logger.FatalF("Failed to parse heartbeat interval: %v", e)
	}

	option := &core.Option{
		Name:              name,
		TickInterval:      time.Duration(i) * time.Second,
		HeartbeatInterval: time.Duration(ih) * time.Second,
	}

	return option
}
