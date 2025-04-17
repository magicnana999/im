package broker

import (
	"github.com/magicnana999/im/pkg/errext"
)

var (
	ErrLogin         = errext.New(1001, "login failed")
	UserConnNotExist = errext.New(1002, "UserConn not exist")
	UserConnIsNil    = errext.New(1003, "UserConn is nil")
	DecodeError      = errext.New(1004, "decode failed")
)
