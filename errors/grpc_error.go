package errors

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	UserSigNotFound = NewGrpcError(codes.NotFound, "user sig not found %s")
	UnmarshalError  = NewGrpcError(codes.Aborted, "unmarshal error")
)

func NewGrpcError(code codes.Code, message string) error {
	return status.Error(code, message)
}

func NewGrpcErrorF(code codes.Code, format string, v ...any) error {
	return status.Errorf(code, format, v)
}

func FormatError(err error) error {
	if err == nil {
		return status.Errorf(codes.Internal, "unknown error")
	}

	if s, b := status.FromError(err); b {
		return s.Err()
	}

	var target Error
	if b := errors.As(err, &target); b {
		return status.Errorf(codes.Internal, target.Message+" "+target.Details)
	}

	return status.Errorf(codes.Internal, err.Error())
}

func Format2ImError(e error) *Error {
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
