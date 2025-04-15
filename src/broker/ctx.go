package broker

import (
	"context"
	"github.com/magicnana999/im/broker/domain"
	"github.com/magicnana999/im/errors"
	"github.com/panjf2000/gnet/v2"
)

const (
	currentUserKey string = `CurrentUser`
)

func CurContext(c gnet.Conn) (context.Context, error) {
	if ctx, o := c.Context().(context.Context); o {
		return ctx, nil
	}
	return nil, errors.CurUserNotFound
}

func CurUserFromCtx(ctx context.Context) (*domain.UserConnection, error) {
	if u, ok := ctx.Value(currentUserKey).(*domain.UserConnection); ok {
		return u, nil
	}

	return nil, errors.CurUserNotFound
}

func CurUserFromConn(c gnet.Conn) (*domain.UserConnection, error) {

	ctx, err := CurContext(c)
	if err != nil {
		return nil, err
	}
	return CurUserFromCtx(ctx)
}
