package state

import (
	"context"
	"github.com/magicnana999/im/domain"
	"github.com/magicnana999/im/enum"
	"github.com/magicnana999/im/errors"
	"github.com/magicnana999/im/redis"
	"strconv"
	"sync"
)

var DefaultUserState = &UserState{}

type UserState struct {
	storage *redis.UserStorage
	m       sync.Map
}

func InitUserState() *UserState {
	DefaultUserState.storage = redis.InitUserStorage()
	return DefaultUserState
}

func (s *UserState) StoreUser(ctx context.Context, u *domain.UserConnection, appId string, userId int64, os enum.OSType) error {

	lock, e := s.storage.Lock(ctx, appId, u.Label())
	if e != nil {
		return errors.UserStoreError.Detail(e)
	}
	defer s.storage.UnLock(ctx, appId, u.Label(), lock)

	u.AppId = appId
	u.UserId = userId
	u.OS = os
	u.IsLogin = true

	s.m.Store(u.Label(), u)

	_, e1 := s.storage.StoreUserConn(ctx, u)
	if e1 != nil {
		return errors.UserStoreError.Detail(e1)

	}

	_, e2 := s.storage.StoreUserClients(ctx, u)
	if e2 != nil {
		return errors.UserStoreError.Detail(e2)
	}

	return nil
}

func (s *UserState) RefreshUser(ctx context.Context, uc *domain.UserConnection) error {
	lock, e := s.storage.Lock(ctx, uc.AppId, uc.Label())
	if e != nil {
		return errors.UserStoreError.Detail(e)
	}
	defer s.storage.UnLock(ctx, uc.AppId, uc.Label(), lock)

	_, e1 := s.storage.RefreshUserConn(ctx, uc)
	if e1 != nil {
		return errors.UserStoreError.Detail(e1)
	}
	return nil
}

func (s *UserState) LoadLocalUser(appId string, userId int64) []*domain.UserConnection {
	labels := label(appId, userId)
	var ret []*domain.UserConnection
	for _, v := range labels {
		if val, ok := s.m.Load(v); ok && val.(*domain.UserConnection).UserId == userId {
			ret = append(ret, val.(*domain.UserConnection))
		}
	}
	if len(ret) == 0 {
		return nil
	}
	return ret
}

func label(appId string, userId int64) []string {
	s1 := appId + "#" + strconv.FormatInt(userId, 10) + "#" + enum.Desktop.String()
	s2 := appId + "#" + strconv.FormatInt(userId, 10) + "#" + enum.Mobile.String()
	return []string{s1, s2}
}
