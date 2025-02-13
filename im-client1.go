package main

func main() {

	c1 := &Client{
		userId:            1200120,
		to:                1200121,
		userSig:           "cukpovu1a37hpofg6sj0",
		serverAddress:     "127.0.0.1:7539",
		heartbeatInterval: 10,
	}

	c1.start()

}
