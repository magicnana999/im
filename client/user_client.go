package main

import "github.com/panjf2000/gnet/v2"

type UserClient struct {
	AppID  string
	UserID int64
	conn   gnet.Conn
}
