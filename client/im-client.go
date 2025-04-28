package main

import (
	"github.com/magicnana999/im/pkg/logger"
	"time"
)

func main() {
	logger.Init(nil)
	defer logger.Close()

	tcp := NewTcpServer(NewPacketHandler(), NewHeartbeatServer(time.Second*30))
	tcp.Start()
	defer tcp.Stop()
}
