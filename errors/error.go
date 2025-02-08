package errors

import (
	"encoding/json"
	"github.com/magicnana999/im/util/str"
)

const (
	brokerSetup = iota + 1101
	brokerRefresh

	connectionHeartbeatInit
	connectionDecode
	connectionEncode
)

var (
	BrokerSetupError   = New(brokerSetup, "broker setup error")
	BrokerRefreshError = New(brokerRefresh, "broker refresh error")

	ConnectionHeartbeatInitError = New(connectionHeartbeatInit, "connection heartbeat init error")
	ConnectionDecodeError        = New(connectionDecode, "connection decode error")
	ConnectionEncodeError        = New(connectionEncode, "connection encode error")
)

const (
	internal = iota + 1201
	handlerNoSupport
	UnmarshalPacket
	invalidCType
	grpcError
	wrapRequest
	wrapReply
)

var (
	HandleInternalError    = New(internal, "internal error")
	HandlerNoSupportError  = New(handlerNoSupport, "no handler is support")
	HandleUnmarshalError   = New(UnmarshalPacket, "unmarshal error")
	HandleInvalidCType     = New(invalidCType, "invalid cType")
	HandleGrpcError        = New(grpcError, "grpc error")
	HandleWrapRequestError = New(wrapRequest, "wrap request error")
	HandleWrapReplyError   = New(wrapReply, "wrap reply error")
)

const (
	ucNotExists = iota + 1501
	ctxNotExists
	userNotLogin
	userStore
	userRefresh
)

var (
	UcNotExists      = New(ucNotExists, "no such user connection")
	CtxNotExists     = New(ctxNotExists, "no such user context")
	UserNotLogin     = New(userNotLogin, "user is not login")
	UserStoreError   = New(userStore, "user store error")
	UserRefreshError = New(userRefresh, "user refresh error")
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
