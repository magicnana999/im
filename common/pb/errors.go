package pb

import (
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
)

var (
	UserSigNotFound = newError(codes.NotFound, "user sig not found %s")
	UnmarshalError  = newError(codes.Aborted, "unmarshal error")
)

type Error struct {
	Code    codes.Code
	Message string
}

func (e Error) Format(v ...string) Error {
	e.Message = fmt.Sprintf(e.Message, v)
	return e
}

func newError(code codes.Code, message string) Error {
	return Error{
		Code:    code,
		Message: message,
	}
}

func (e Error) Error() string {
	return fmt.Sprintf("Code: %v, Message: %s", e.Code, e.Message)
}

func FromError(err error) *Error {
	var e *Error
	if errors.As(err, &e) {
		return e
	}
	return nil
}
