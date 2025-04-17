package errors

import (
	"github.com/magicnana999/im/pkg/error"
)

var (
	DecodeError         = error.New(1101, "decode failed")
	EncodeError         = error.New(1102, "encode failed")
	NoHandlerSupport    = error.New(1103, "no handler support")
	HeartbeatError      = error.New(1104, "heartbeat failed")
	HeartbeatTimeout    = error.New(1105, "heartbeat timeout")
	MsgMQProduceError   = error.New(1106, "message produce failed")
	MsgDeliverTaskError = error.New(1107, "message deliver task failed")
	CurUserNotFound     = error.New(1108, "current user not found")

	CmdError        = error.New(1201, "cmd_service failed")
	CmdUnknownType  = error.New(1202, "unknown cmd_service type")
	CmdResponseNull = error.New(1203, "response is null")
	CmdReplyNull    = error.New(1204, "reply is null")
)
