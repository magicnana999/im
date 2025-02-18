package errors

import (
	"encoding/json"
	"errors"
)

const (
	decode = iota + 1101
	encode
	noHandler
	heartbeat
	heartbeatTask
	command
	cmdUnknownType
	msgMQProduce
	msgDeliverTask
)

var (
	DecodeError         = New(decode, "decode failed")
	EncodeError         = New(encode, "encode failed")
	NoHandlerSupport    = New(noHandler, "no handler support")
	HeartbeatError      = New(heartbeat, "heartbeat failed")
	HeartbeatTaskError  = New(heartbeatTask, "heartbeatTask failed")
	CommandHandleError  = New(command, "command failed")
	CmdUnknownTypeError = New(cmdUnknownType, "unknown command type")
	MsgMQProduceError   = New(msgMQProduce, "message produce failed")
	MsgDeliverTaskError = New(msgDeliverTask, "message deliver task failed")
)

const (
	stateStoreBroker = iota + 1151
	stateStoreUser
	stateGetCtx
	stateGetUser
)

var (
	BrokerStoreError = New(stateStoreBroker, "store broker failed")
	UserStoreError   = New(stateStoreUser, "store user failed")
	GetCtxError      = New(stateGetCtx, "get context failed")
	GetUserError     = New(stateGetUser, "get user failed")
)

// GRPC API
const (
	userSigNotFound = iota + 2001
)

var (
	UserSigNotFound = New(userSigNotFound, "user sig not found")
)

///////////////////////////////

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e Error) String() string {
	return e.Message + " " + e.Details
}

func (e Error) Error() string {

	js, err := json.Marshal(e)
	if err != nil {
		return "{}"
	}
	return string(js)
}

func (e Error) DetailString(str string) Error {
	e.Details = str
	return e
}

func (e Error) DetailJson(m map[string]any) Error {
	js, err := json.Marshal(m)
	if err == nil {
		e.Details = string(js)
	}
	return e
}

func (e Error) Detail(err error) Error {
	var ime Error
	if errors.As(err, &ime) {
		e.Details = ime.Details
	} else {
		e.Details = err.Error()
	}
	return e
}

func New(code int, message string) Error {
	return Error{
		Code:    code,
		Message: message,
	}
}

func Format(e error) *Error {
	if e == nil {
		return nil
	}

	var ime Error
	if ok := errors.As(e, &ime); ok {
		return &ime
	}

	ime = New(-1, e.Error())
	return &ime
}
