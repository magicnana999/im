package broker

import (
	"context"
	"errors"
	"github.com/magicnana999/im/broker/domain"
)

const (
	currentUserKey string = `CurrentUser`
)

var (
	curUserNil      = errors.New("current user is nil")
	curUserNotExist = errors.New("current user not exist")
)

// GetCurUserConn 从当前协程（和Connection绑定）中获取UserConn
func GetCurUserConn(ctx context.Context) (*domain.UserConn, error) {
	if u, ok := ctx.Value(currentUserKey).(*domain.UserConn); ok {
		if u == nil {
			return nil, curUserNil
		} else {
			return u, nil
		}
	}
	return nil, curUserNotExist
}
