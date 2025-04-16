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

func Context(c gnet.Conn) (context.Context, error) {
	if ctx, o := c.Context().(context.Context); o {
		return ctx, nil
	}
	return nil, errors.CurUserNotFound
}

func UserFromCtx(ctx context.Context) (*domain.UserConn, error) {
	if u, ok := ctx.Value(currentUserKey).(*domain.UserConn); ok {
		return u, nil
	}

	return nil, errors.CurUserNotFound
}

func UserFromConn(c gnet.Conn) (*domain.UserConn, error) {

	ctx, err := Context(c)
	if err != nil {
		return nil, err
	}
	return UserFromCtx(ctx)
}
