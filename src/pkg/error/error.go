package error

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/magicnana999/im/pkg/logger"
	"github.com/magicnana999/im/pkg/utils"
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
	return e.Message + " " + e.Detail
}

func (e Error) JsonString() string {
	js, err := json.Marshal(e)
	if err != nil {
		logger.Errorf("error json marshal err: %v", err)
		return "{}"
	}
	return string(js)
}

func (e Error) Error() string {
	return e.JsonString()
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

func (e Error) SetDetail(a any) Error {
	d, ee := utils.Any2String(a)
	if ee != nil {
		fmt.Printf("parse to string error:%v", ee)
	}
	e.Detail = d
	return e
}

func (e Error) FmtDetail(template string, args ...any) Error {
	s := fmt.Sprintf(template, args)
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
		return New(-1, "no error found")
	}

	var ime Error
	if ok := errors.As(e, &ime); ok {
		return ime
	}

	ime = New(-1, e.Error())
	return ime
}
