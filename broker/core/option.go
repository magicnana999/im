package core

import "time"

const (
	DefaultPort = "7539"
)

type Option struct {
	Name              string        `json:"name"`
	ServerInterval    time.Duration `json:"serverInterval"`
	HeartbeatInterval time.Duration `json:"heartbeatInterval"`
}
