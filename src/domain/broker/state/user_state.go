package state

import (
	"context"
	"github.com/magicnana999/im/broker/domain"
	"github.com/magicnana999/im/constants"
	"github.com/magicnana999/im/domain/broker/storage"
	"sync"
)

var DefaultUserState *UserState
var dusOnce sync.Once

type UserState struct {
	storage *storage.UserStorage
	m       sync.Map
}

func InitUserState() *UserState {

	dusOnce.Do(func() {
		DefaultUserState = &UserState{
			storage: storage.InitUserStorage(),
		}
	})
	return DefaultUserState
}

func (s *UserState) StoreUser(ctx context.Context, u *domain.UserConnection, appId string, userId int64, os constants.OSType) error {

	lock, e := s.storage.Lock(ctx, appId, u.Label())
	if e != nil {
		return e
	}
	defer s.storage.UnLock(ctx, appId, u.Label(), lock)

	u.AppId = appId
	u.UserId = userId
	u.OS = os
	u.IsLogin = true

	s.m.Store(u.Label(), u)

	_, e1 := s.storage.StoreUserConn(ctx, u)
	if e1 != nil {
		return e1

	}

	_, e2 := s.storage.StoreUserClients(ctx, u)
	if e2 != nil {
		return e1
	}

	return nil
}

func (s *UserState) RefreshUser(ctx context.Context, uc *domain.UserConnection) error {
	lock, e := s.storage.Lock(ctx, uc.AppId, uc.Label())
	if e != nil {
		return e
	}
	defer s.storage.UnLock(ctx, uc.AppId, uc.Label(), lock)

	_, e1 := s.storage.RefreshUserConn(ctx, uc)
	if e1 != nil {
		return e1
	}
	return nil
}

func (s *UserState) loadLocalUser(label string) *domain.UserConnection {
	if val, ok := s.m.Load(label); ok {
		return val.(*domain.UserConnection)
	}
	return nil
}
