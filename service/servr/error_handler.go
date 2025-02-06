package servr

import (
	"context"
	"github.com/magicnana999/im/errors"
	"google.golang.org/grpc"
)

func ErrorHandlingInterceptor(
	ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {

	resp, err := handler(ctx, req)

	if err != nil {
		return nil, errors.FormatError(err)
	}

	return resp, err
}
