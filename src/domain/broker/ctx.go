package broker

import (
	"context"
	"github.com/magicnana999/im/dto/broker"
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

func UserFromCtx(ctx context.Context) (*broker.UserConnection, error) {
	if u, ok := ctx.Value(currentUserKey).(*broker.UserConnection); ok {
		return u, nil
	}

	return nil, errors.CurUserNotFound
}

func UserFromConn(c gnet.Conn) (*broker.UserConnection, error) {

	ctx, err := Context(c)
	if err != nil {
		return nil, err
	}
	return UserFromCtx(ctx)
}
