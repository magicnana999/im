package storage

import (
	"context"
	"encoding/json"
	"github.com/magicnana999/im/broker/domain"
	entity "github.com/magicnana999/im/entities"
	inf "github.com/magicnana999/im/infra"
	"github.com/magicnana999/im/pkg/utils"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	DefaultUserStorage *UserStorage
	udsOnce            sync.Once
)

type UserStorage struct {
}

func InitUserStorage() *UserStorage {
	udsOnce.Do(func() {
		inf.NewRedisClient(nil)
		DefaultUserStorage = &UserStorage{}
	})

	return DefaultUserStorage
}

func (s *UserStorage) LoadUserConn(ctx context.Context, appId string, userId int64) (map[string]*domain.UserConnection, error) {
	key := inf.KeyUserClients(appId, userId)
	cmd := inf.RDS.HGetAll(ctx, key)
	if cmd.Err() == nil {
		m := cmd.Val()
		ret := make(map[string]*domain.UserConnection, len(m))
		for k, v := range m {

			var uc domain.UserConnection
			ee := json.Unmarshal([]byte(v), &uc)
			if ee != nil {
				return nil, ee
			}

			ret[k] = &uc
		}
		return ret, nil
	} else {
		return nil, cmd.Err()
	}
}

func (s *UserStorage) StoreUserConn(ctx context.Context, uc *domain.UserConnection) (string, error) {

	key := inf.KeyUserConn(uc.AppId, uc.Label())

	js, err := json.Marshal(uc)
	if err != nil {
		return "", err
	}

	ret := inf.RDS.Set(ctx, key, string(js), time.Minute)

	return ret.Val(), ret.Err()
}

func (s *UserStorage) RefreshUserConn(ctx context.Context, uc *domain.UserConnection) (bool, error) {
	key := inf.KeyUserConn(uc.AppId, uc.Label())
	ret := inf.RDS.Expire(ctx, key, time.Minute)
	return ret.Val(), ret.Err()

}

func (s *UserStorage) StoreUserClients(ctx context.Context, uc *domain.UserConnection) (int64, error) {

	key := inf.KeyUserClients(uc.AppId, uc.UserId)

	js, err := json.Marshal(uc)
	if err != nil {
		return 0, err
	}

	ret := inf.RDS.HSet(ctx, key, uc.Label(), string(js))

	return ret.Val(), ret.Err()
}

func (s *UserStorage) Lock(ctx context.Context, appId, ucLabel string) (string, error) {
	key := inf.KeyUserConnLock(appId, ucLabel)
	val := time.Now().UnixMilli()
	ret := inf.RDS.SetNX(ctx, key, strconv.FormatInt(val, 10), time.Minute)
	return strconv.FormatInt(val, 10), ret.Err()
}

func (s *UserStorage) UnLock(ctx context.Context, appId, ucLabel, lock string) (int64, error) {
	key := inf.KeyUserConnLock(appId, ucLabel)
	ret := inf.RDS.Del(ctx, key)
	return ret.Val(), ret.Err()
}

func (s *UserStorage) LoadUserSig(ctx context.Context, appId, userSig string) (*entity.User, error) {
	cmd := inf.RDS.Get(ctx, inf.KeyUserSig(appId, userSig))

	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	var user entity.User
	err := json.Unmarshal([]byte(cmd.Val()), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserStorage) StoreUserSig(ctx context.Context, appId string, user *entity.User) (string, error) {
	sig := strings.ToLower(utils.GenerateXId())
	json1, _ := json.Marshal(user)

	cmd := inf.RDS.Set(ctx, inf.KeyUserSig(appId, sig), json1, -1)
	return cmd.Val(), cmd.Err()
}

func (s *UserStorage) LoadByUserId(ctx context.Context, appId string, userId int64) (*entity.User, error) {
	cmd := inf.RDS.Get(ctx, inf.KeyUser(appId, userId))

	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	var user entity.User
	err := json.Unmarshal([]byte(cmd.Val()), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserStorage) StoreByUserId(ctx context.Context, user *entity.User) (string, error) {
	json1, _ := json.Marshal(user)
	cmd := inf.RDS.Set(ctx, inf.KeyUser(user.AppId, user.UserId), json1, -1)
	return cmd.Val(), cmd.Err()
}
