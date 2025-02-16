package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/magicnana999/im/domain"
	"github.com/magicnana999/im/errors"
	"strconv"
	"time"
)

var DefaultUserStorage = &UserStorage{}

type UserStorage struct {
}

func InitUserStorage() *UserStorage {
	initRedis()
	return DefaultUserStorage
}

func (s *UserStorage) LoadUserConn(ctx context.Context, appId string, userId int64) (map[string]*domain.UserConnection, error) {
	key := KeyUserClients(appId, userId)
	cmd := rds.HGetAll(ctx, key)
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

	key := KeyUserConn(uc.AppId, uc.Label())

	js, err := json.Marshal(uc)
	if err != nil {
		return "", errors.UserStoreError.Detail(err)
	}

	ret := rds.Set(ctx, key, string(js), time.Minute)

	return ret.Val(), ret.Err()
}

func (s *UserStorage) RefreshUserConn(ctx context.Context, uc *domain.UserConnection) (bool, error) {
	key := KeyUserConn(uc.AppId, uc.Label())
	ret := rds.Expire(ctx, key, time.Minute)
	return ret.Val(), ret.Err()

}

func (s *UserStorage) StoreUserClients(ctx context.Context, uc *domain.UserConnection) (int64, error) {

	key := KeyUserClients(uc.AppId, uc.UserId)

	js, err := json.Marshal(uc)
	if err != nil {
		return 0, errors.UserStoreError.Detail(err)
	}

	ping := rds.Ping(ctx)
	fmt.Println(ping)

	ret := rds.HSet(ctx, key, uc.Label(), string(js))

	return ret.Val(), ret.Err()
}

func (s *UserStorage) Lock(ctx context.Context, appId, ucLabel string) (string, error) {
	key := KeyUserLock(appId, ucLabel)
	val := time.Now().UnixMilli()
	ret := rds.SetNX(ctx, key, strconv.FormatInt(val, 10), time.Minute)
	return strconv.FormatInt(val, 10), ret.Err()
}

func (s *UserStorage) UnLock(ctx context.Context, appId, ucLabel, lock string) (int64, error) {
	key := KeyUserLock(appId, ucLabel)
	ret := rds.Del(ctx, key)
	return ret.Val(), ret.Err()
}
