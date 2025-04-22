package main

import (
	"io"
	"sync/atomic"
)

type User struct {
	FD       int          `json:"fd"`
	Writer   io.Writer    `json:"-"`
	Reader   io.Reader    `json:"-"`
	AppID    string       `json:"appId"`
	UserID   int64        `json:"userId"`
	IsClosed atomic.Bool  `json:"-"`
	IsLogin  atomic.Bool  `json:"-"`
	LastHTS  atomic.Int64 `json:"-"`
}
