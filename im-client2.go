package main

import "github.com/magicnana999/im/client"

func main() {

	c1 := &client.Client{
		UserId:            1200121,
		To:                1200120,
		UserSig:           "cuv0ele1a37rf7ccvang",
		ServerAddress:     "127.0.0.1:7539",
		HeartbeatInterval: 10,
	}

	c1.Start()

}
