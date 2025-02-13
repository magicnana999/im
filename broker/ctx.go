package broker

import (
	"context"
	"github.com/magicnana999/im/domain"
	"github.com/magicnana999/im/errors"
	"github.com/panjf2000/gnet/v2"
)

const (
	currentUserKey string = `CurrentUser`
)

func currentContext(c gnet.Conn) (context.Context, error) {
	if ctx, o := c.Context().(context.Context); o {
		return ctx, nil
	}
	return nil, errors.GetCtxError
}

func currentUserFromCtx(ctx context.Context) (*domain.UserConnection, error) {
	if u, ok := ctx.Value(currentUserKey).(*domain.UserConnection); ok {
		return u, nil
	}

	return nil, errors.GetUserError
}

func currentUserFromConn(c gnet.Conn) (*domain.UserConnection, error) {

	ctx, err := currentContext(c)
	if err != nil {
		return nil, err
	}
	return currentUserFromCtx(ctx)
}
