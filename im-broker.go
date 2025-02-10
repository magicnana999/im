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
	var serverInterval int
	var heartbeatInterval int
	flag.StringVar(&name, "name", "", "The name of the broker instance")
	flag.IntVar(&serverInterval, "interval", 30, "Ticking interval of broker instance")
	flag.IntVar(&heartbeatInterval, "heartbeatInterval", 30, "Heartbeat interval")
	flag.Parse()

	if str.IsBlank(name) {
		ipAddress, e := ip.GetLocalIP()
		if str.IsBlank(ipAddress) || e != nil {
			logger.FatalF("Failed to get local IP address: %v", e)
		}
		name = fmt.Sprintf("%s:%s", ipAddress, core.DefaultPort)
	}

	if serverInterval <= 0 {
		logger.FatalF("Invalid server interval value: %d (must be > 0)", serverInterval)
	}
	if heartbeatInterval <= 0 {
		logger.FatalF("Invalid heartbeat interval value: %d (must be > 0)", heartbeatInterval)
	}

	return &core.Option{
		Name:              name,
		ServerInterval:    time.Duration(serverInterval) * time.Second,
		HeartbeatInterval: time.Duration(heartbeatInterval) * time.Second,
	}
}
