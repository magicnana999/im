package errext

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	NilError     = New(-1, "nil error")
	DefaultError = New(-2, "default error")
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

func (e Error) ShortString() string {
	return e.Message
}

func (e Error) LongString() string {
	if e.Detail != "" {
		return e.Message + "," + e.Detail
	} else {
		return e.Message
	}
}

func (e Error) JsonBytes() ([]byte, error) {
	return json.Marshal(e)
}

func (e Error) JsonString() (string, error) {
	if bs, err := e.JsonBytes(); err != nil {
		return "", err
	} else {
		return string(bs), nil
	}
}

func (e Error) Error() string {
	return e.LongString()
}

func (e Error) GetMessage() string {
	return e.Message
}

func (e Error) GetDetail() string {
	return e.Detail
}

func (e Error) GetCode() int {
	return e.Code
}

func (e Error) SetMessage(msg string) Error {
	e.Message = msg
	return e
}

func (e Error) SetDetail(detail string) Error {
	e.Detail = detail
	return e
}

func (e Error) FmtDetail(template string, args ...any) Error {
	s := fmt.Sprintf(template, args...)
	e.Detail = s
	return e
}

func (e Error) SetCode(code int) Error {
	e.Code = code
	return e
}

func New(code int, message string) Error {
	return Error{
		Code:    code,
		Message: message,
	}
}

func Format(e error) Error {
	if e == nil {
		return NilError
	}

	var ime Error
	if ok := errors.As(e, &ime); ok {
		return ime
	}

	ime = New(DefaultError.Code, e.Error())
	return ime
}
