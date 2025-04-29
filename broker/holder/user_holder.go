package holder

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/magicnana999/im/broker/domain"
	"github.com/magicnana999/im/infra"
	"go.uber.org/fx"
	"sync"
	"time"
)

type UserHolder struct {
	rds *redis.Client
	m   sync.Map
}

func NewUserHolder(rds *redis.Client, lf fx.Lifecycle) (*UserHolder, error) {
	return &UserHolder{rds: rds, m: sync.Map{}}, nil
}

//func (s *UserHolder) StoreUser(ctx context.Context, u *domain.UserConn, appId string, userId int64, os string) error {
//
//	lock, e := s.Lock(ctx, appId, u.Label())
//	if e != nil {
//		return e
//	}
//	defer s.UnLock(ctx, appId, u.Label(), lock)
//
//	u.AppId = appId
//	u.UserId = userId
//	u.OS = define.OSType(os)
//	u.IsLogin.Store(true)
//
//	s.m.Store(u.Label(), u)
//
//	_, e1 := s.StoreUserConn(ctx, u)
//	if e1 != nil {
//		return e1
//
//	}
//
//	_, e2 := s.StoreUserClients(ctx, u)
//	if e2 != nil {
//		return e1
//	}
//
//	return nil
//}

//func (s *UserHolder) RefreshUser(ctx context.Context, uc *domain.UserConn) error {
//	lock, e := s.Lock(ctx, uc.AppId, uc.Label())
//	if e != nil {
//		return e
//	}
//	defer s.UnLock(ctx, uc.AppId, uc.Label(), lock)
//
//	_, e1 := s.RefreshUserConn(ctx, uc)
//	if e1 != nil {
//		return e1
//	}
//	return nil
//}

func (s *UserHolder) DeleteUserConn(ctx context.Context, uc *domain.UserConn) (int64, error) {
	key := infra.KeyUserConn(uc.AppId.Load(), uc.Label())
	ret := s.rds.Del(ctx, key)
	return ret.Val(), ret.Err()
}

func (s *UserHolder) StoreUserConn(ctx context.Context, uc *domain.UserConn) (string, error) {
	key := infra.KeyUserConn(uc.AppId.Load(), uc.Label())
	js, err := json.Marshal(uc)
	if err != nil {
		return "", err
	}

	ret := s.rds.Set(ctx, key, string(js), time.Minute)
	return ret.Val(), ret.Err()
}

func (s *UserHolder) RefreshUserConn(ctx context.Context, uc *domain.UserConn) (bool, error) {
	key := infra.KeyUserConn(uc.AppId.Load(), uc.Label())
	ret := s.rds.Expire(ctx, key, time.Minute)
	return ret.Val(), ret.Err()

}

func (s *UserHolder) DeleteUserClient(ctx context.Context, uc *domain.UserConn) (int64, error) {
	key := infra.KeyUserClients(uc.AppId.Load(), uc.UserId.Load())
	ret := s.rds.HDel(ctx, key, uc.Label())
	return ret.Val(), ret.Err()
}

func (s *UserHolder) StoreUserClients(ctx context.Context, uc *domain.UserConn) (int64, error) {
	key := infra.KeyUserClients(uc.AppId.Load(), uc.UserId.Load())
	js, err := json.Marshal(uc)
	if err != nil {
		return 0, err
	}
	ret := s.rds.HSet(ctx, key, uc.Label(), string(js))
	return ret.Val(), ret.Err()
}

//	func (s *UserHolder) Lock(ctx context.Context, appId, ucLabel string) (string, error) {
//		key := infra.KeyUserConnLock(appId, ucLabel)
//		val := time.Now().UnixMilli()
//		ret := s.rds.SetNX(ctx, key, strconv.FormatInt(val, 10), time.Minute)
//		return strconv.FormatInt(val, 10), ret.Err()
//	}
//
//	func (s *UserHolder) UnLock(ctx context.Context, appId, ucLabel, lock string) (int64, error) {
//		key := infra.KeyUserConnLock(appId, ucLabel)
//		ret := s.rds.Del(ctx, key)
//		return ret.Val(), ret.Err()
//	}
//
//	func (s *UserHolder) LoadUserSig(ctx context.Context, appId, userSig string) (*entity.User, error) {
//		cmd := s.rds.Get(ctx, infra.KeyUserSig(appId, userSig))
//
//		if cmd.Err() != nil {
//			return nil, cmd.Err()
//		}
//
//		var user entity.User
//		err := json.Unmarshal([]byte(cmd.Val()), &user)
//		if err != nil {
//			return nil, err
//		}
//		return &user, nil
//	}
//
//	func (s *UserHolder) StoreUserSig(ctx context.Context, appId string, user *entity.User) (string, error) {
//		sig := strings.ToLower(id.GenerateXId())
//		json1, _ := json.Marshal(user)
//
//		cmd := s.rds.Set(ctx, infra.KeyUserSig(appId, sig), json1, -1)
//		return cmd.Val(), cmd.Err()
//	}
//
//	func (s *UserHolder) LoadByUserId(ctx context.Context, appId string, userId int64) (*entity.User, error) {
//		cmd := s.rds.Get(ctx, infra.KeyUser(appId, userId))
//
//		if cmd.Err() != nil {
//			return nil, cmd.Err()
//		}
//
//		var user entity.User
//		err := json.Unmarshal([]byte(cmd.Val()), &user)
//		if err != nil {
//			return nil, err
//		}
//		return &user, nil
//	}
//func (s *UserHolder) StoreByUserId(ctx context.Context, user *entity.User) (string, error) {
//	json1, _ := json.Marshal(user)
//	cmd := s.rds.Set(ctx, infra.KeyUser(user.AppId, user.UserId), json1, -1)
//	return cmd.Val(), cmd.Err()
//}

// GetUserConn 从本地获取，内含gnet.Conn
func (s *UserHolder) GetUserConn(label string) *domain.UserConn {
	if val, ok := s.m.Load(label); ok {
		return val.(*domain.UserConn)
	}
	return nil
}

// HoldUserConn 保存到本地，内含gnet.Conn
func (s *UserHolder) HoldUserConn(uc *domain.UserConn) bool {
	_, ok := s.m.LoadOrStore(uc.Label(), uc)
	return !ok
}

// RemoveUserConn 删除本地UC
func (s *UserHolder) RemoveUserConn(uc *domain.UserConn) bool {
	_, ok := s.m.LoadAndDelete(uc.Label())
	return ok
}

// RangeAllUserConn 遍历本地UC
func (s *UserHolder) RangeAllUserConn(fun func(*domain.UserConn) bool) {
	s.m.Range(func(key, value any) bool {
		return fun(value.(*domain.UserConn))
	})
}
