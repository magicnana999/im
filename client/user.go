package main

import (
	"io"
	"sync/atomic"
)

type User struct {
	Writer   io.Writer
	Reader   io.Reader
	AppID    string
	UserID   int64
	IsClosed atomic.Bool
}
