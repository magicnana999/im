package state

import (
	"context"
	"github.com/magicnana999/im/domain"
	"github.com/magicnana999/im/errors"
	"github.com/panjf2000/gnet/v2"
)

const (
	CurrentUser string = `CurrentUser`
)

func CurrentContextFromConn(c gnet.Conn) (context.Context, error) {
	if ctx, o := c.Context().(context.Context); o {
		return ctx, nil
	}
	return nil, errors.GetCtxError
}

func CurrentUserFromCtx(ctx context.Context) (*domain.UserConnection, error) {
	if u, ok := ctx.Value(CurrentUser).(*domain.UserConnection); ok {
		return u, nil
	}

	return nil, errors.GetUserError
}

func CurrentUserFromConn(c gnet.Conn) (*domain.UserConnection, error) {

	ctx, err := CurrentContextFromConn(c)
	if err != nil {
		return nil, err
	}
	return CurrentUserFromCtx(ctx)
}
