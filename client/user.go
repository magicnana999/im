package main

import (
	"io"
	"sync/atomic"
)

type User struct {
	FD       int
	Writer   io.Writer
	Reader   io.Reader
	AppID    string
	UserID   int64
	IsClosed atomic.Bool
	IsLogin  atomic.Bool
	LastHTS  atomic.Int64
}
