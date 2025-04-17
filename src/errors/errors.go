package errors

import (
	"github.com/magicnana999/im/pkg/errext"
)

var (
	DecodeError         = errext.New(1101, "decode failed")
	EncodeError         = errext.New(1102, "encode failed")
	NoHandlerSupport    = errext.New(1103, "no handler support")
	HeartbeatError      = errext.New(1104, "heartbeat failed")
	HeartbeatTimeout    = errext.New(1105, "heartbeat timeout")
	MsgMQProduceError   = errext.New(1106, "message produce failed")
	MsgDeliverTaskError = errext.New(1107, "message deliver task failed")
	CurUserNotFound     = errext.New(1108, "current user not found")

	CmdError        = errext.New(1201, "cmd_service failed")
	CmdUnknownType  = errext.New(1202, "unknown cmd_service type")
	CmdResponseNull = errext.New(1203, "response is null")
	CmdReplyNull    = errext.New(1204, "reply is null")
)
