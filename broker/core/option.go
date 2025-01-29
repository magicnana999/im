package core

import "time"

const (
	DefaultPort = "7539"
)

type Option struct {
	Name              string        `json:"name"`
	TickInterval      time.Duration `json:"tickInterval"`
	HeartbeatInterval time.Duration `json:"heartbeatInterval"`
}
