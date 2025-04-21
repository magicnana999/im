package main

import "time"

func main() {
	tcp := NewTcpServer(NewPacketHandler(), NewHeartbeatServer(time.Second*30))
	tcp.Start()
	defer tcp.Stop()
}
