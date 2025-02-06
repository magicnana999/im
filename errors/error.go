package errors

import (
	"encoding/json"
	"github.com/magicnana999/im/util/str"
)

const (
	brokerSetup = iota + 1101
	brokerRefresh
)

var (
	BrokerSetupError   = New(brokerSetup, "broker setup error")
	BrokerRefreshError = New(brokerRefresh, "broker refresh error")
)

const (
	connectionHeartbeatInit = iota + 1201
	connectionDecode
	connectionEncode
)

var (
	ConnectionHeartbeatInitError = New(connectionHeartbeatInit, "connection heartbeat init error")
	ConnectionDecodeError        = New(connectionDecode, "connection decode error")
	ConnectionEncodeError        = New(connectionEncode, "connection encode error")
)

const (
	ucNotExists = iota + 1301
	ctxNotExists
)

var (
	UcNotExists  = New(ucNotExists, "no such user connection")
	CtxNotExists = New(ctxNotExists, "no such context")
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e Error) Error() string {

	js, err := json.Marshal(e)
	if err != nil {
		return "{}"
	}
	return string(js)
}

func (e Error) Fill(detail string) Error {
	if str.IsNotBlank(detail) {
		e.Details = detail
	}
	return e
}

func New(code int, message string) Error {
	return Error{
		Code:    code,
		Message: message,
	}
}

func CompleteError(code int, message string, detail string) Error {
	return Error{
		Code:    code,
		Message: message,
		Details: detail,
	}
}
