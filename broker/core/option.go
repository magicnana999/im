package core

import "time"

type Option struct {
	Addr         string
	TickDuration time.Duration
}

var DefaultOption = Option{
	Addr:         "127.0.0.1:7539",
	TickDuration: 30 * time.Second,
}
