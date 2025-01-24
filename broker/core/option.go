package core

import "time"

const (
	DefaultPort = "7539"
)

type Option struct {
	Name     string
	Interval time.Duration
}
