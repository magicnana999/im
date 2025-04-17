package broker

import (
	"context"
	"github.com/magicnana999/im/broker/domain"
)

const (
	currentUserKey string = `CurrentUser`
)

// GetCurUserConn 从当前协程（和Connection绑定）中获取UserConn
func GetCurUserConn(ctx context.Context) (*domain.UserConn, error) {
	if u, ok := ctx.Value(currentUserKey).(*domain.UserConn); ok {
		return u, nil
	}
	return nil, UserConnNotExist
}
