package ctx

import (
	"context"
	"errors"
	"github.com/magicnana999/im/broker/domain"
)

const (
	CurrentUserKey string = `CurrentUser`
)

var (
	curUserNil      = errors.New("current user is nil")
	curUserNotExist = errors.New("current user not exist")
	ctxNotExist     = errors.New("context nil")
)

// GetCurUserConn 从当前协程（和Connection绑定）中获取UserConn
func GetCurUserConn(ctx context.Context) (*domain.UserConn, error) {

	if ctx == nil {
		return nil, ctxNotExist
	}

	if u, ok := ctx.Value(CurrentUserKey).(*domain.UserConn); ok {
		if u == nil {
			return nil, curUserNil
		} else {
			return u, nil
		}
	}
	return nil, curUserNotExist
}
