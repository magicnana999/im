package state

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/magicnana999/im/enum"
	"github.com/magicnana999/im/errors"
	"github.com/magicnana999/im/logger"
	"github.com/magicnana999/im/redis"
	"github.com/magicnana999/im/util/str"
	"github.com/panjf2000/gnet/v2"
	"strconv"
	"sync"
	"time"
)

const (
	KeyBrokerInfo = "im:broker:"

	keyLock = ":im:user:lock:"

	keyUserConn    = ":im:user:connect:"
	keyUserClients = ":im:user:clients:"

	KeyBrokerConnections               = "im:broker:connections:"
	KeyUserClients                     = "im:user:clients:"
	Expire               time.Duration = 60 * time.Second
)

const (
	CurrentUser string = `CurrentUser`
)

var (
	m  map[string]*UserConnection
	mu sync.RWMutex
)

type BrokerInfo struct {
	Addr    string `json:"addr"`
	StartAt int64  `json:"startAt"`
}

func NewBrokerInfo(addr string) BrokerInfo {
	return BrokerInfo{
		Addr:    addr,
		StartAt: time.Now().Unix(),
	}
}

func init() {
	m = make(map[string]*UserConnection)
}

type UserConnection struct {
	Fd          int         `json:"fd"`
	AppId       string      `json:"appId"`
	UserId      int64       `json:"userId"`
	ClientAddr  string      `json:"clientAddr"`
	BrokerAddr  string      `json:"brokerAddr"`
	OS          enum.OSType `json:"os"`
	ConnectTime int64       `json:"connectTime"`
	IsLogin     bool        `json:"isLogin"`
	C           gnet.Conn   `json:"-"`
}

func OpenUserConnection(c gnet.Conn) *UserConnection {
	return &UserConnection{
		Fd:          c.Fd(),
		AppId:       "",
		UserId:      0,
		ClientAddr:  c.RemoteAddr().String(),
		BrokerAddr:  c.LocalAddr().String(),
		OS:          enum.OSType(0),
		ConnectTime: time.Now().UnixMilli(),
		C:           c,
	}
}

func (u *UserConnection) Label() string {
	if u.UserId == 0 {
		return ""
	}

	dt := u.OS.GetDeviceType()
	if !dt.Valid() {
		return ""
	}

	if str.IsBlank(dt.String()) {
		return ""
	}

	return strconv.FormatInt(u.UserId, 10) + "#" + dt.String()
}

func (u *UserConnection) Store(ctx context.Context, appId string, userId int64, os enum.OSType) error {
	mu.Lock()
	defer mu.Unlock()
	u.AppId = appId
	u.UserId = userId
	u.OS = os
	u.IsLogin = true
	m[u.Label()] = u

	return StoreUserToRedis(ctx, u)
}

func (u *UserConnection) Refresh(ctx context.Context) error {
	return RefreshUserToRedis(ctx, u)
}

func Load(ucLabel string) (*UserConnection, error) {
	mu.RLock()
	defer mu.RUnlock()
	uc := m[ucLabel]
	if uc == nil {
		return nil, errors.UserNotLogin.Fill(ucLabel)
	}

	return uc, nil

}

func CurrentContextFromConn(c gnet.Conn) (context.Context, error) {
	if ctx, o := c.Context().(context.Context); o {
		return ctx, nil
	}

	return nil, errors.CtxNotExists
}

func CurrentUserFromConn(c gnet.Conn) (*UserConnection, error) {

	ctx, err := CurrentContextFromConn(c)
	if err != nil {
		return nil, err
	}
	return CurrentUserFromCtx(ctx)
}

func CurrentUserFromCtx(ctx context.Context) (*UserConnection, error) {
	if u, ok := ctx.Value(CurrentUser).(*UserConnection); ok {
		return u, nil
	}

	return nil, errors.UcNotExists
}

func SetupBroker(ctx context.Context, broker BrokerInfo) (string, error) {
	json, err := json.Marshal(broker)
	if err != nil {
		return "", errors.BrokerSetupError.Fill(err.Error())
	}

	key := fmt.Sprintf("%s%s", KeyBrokerInfo, broker.Addr)
	ret := redis.RDS.Set(ctx, key, json, Expire)

	if ret.Err() != nil {
		return "", errors.BrokerSetupError.Fill(ret.Err().Error())
	}

	logger.DebugF("BrokerInfo setup,key:%s,result:%s",
		key,
		ret.Val())

	return ret.Val(), ret.Err()
}

func RefreshBroker(ctx context.Context, broker BrokerInfo) (bool, error) {
	key := fmt.Sprintf("%s%s", KeyBrokerInfo, broker.Addr)
	ret := redis.RDS.Expire(ctx, key, Expire)

	if ret.Err() != nil {
		return false, errors.BrokerRefreshError.Fill(ret.Err().Error())
	}

	//logger.DebugF("BrokerInfo refresh,key:%s,result:%t",
	//	key,
	//	ret.Val())
	return ret.Val(), ret.Err()
}

func StoreUserToRedis(ctx context.Context, uc *UserConnection) error {
	lock, e := redisLock(ctx, uc.AppId, uc.Label())
	if e != nil {
		return e
	}

	defer redisUnlock(ctx, uc.AppId, uc.Label(), lock)

	_, e1 := setupUserConn(ctx, uc)
	if e1 != nil {
		return e1
	}

	_, e2 := setupUserClients(ctx, uc)
	if e2 != nil {
		return e2
	}

	return nil
}

func RefreshUserToRedis(ctx context.Context, uc *UserConnection) error {
	lock, e := redisLock(ctx, uc.AppId, uc.Label())
	if e != nil {
		return e
	}

	defer redisUnlock(ctx, uc.AppId, uc.Label(), lock)

	_, e1 := refreshUserConn(ctx, uc)
	if e1 != nil {
		return e1
	}

	return nil
}

func redisLock(ctx context.Context, appId, ucLabel string) (string, error) {
	key := fmt.Sprintf("%s%s%s", appId, keyLock, ucLabel)
	val := time.Now().UnixMilli()
	ret := redis.RDS.SetNX(ctx, key, strconv.FormatInt(val, 10), time.Minute)
	if ret.Err() != nil {
		return "", errors.UserStoreError.Fill(ret.Err().Error())
	}

	if !ret.Val() {
		return "", errors.UserStoreError.Fill("redis is not ok")
	}

	return strconv.FormatInt(val, 10), nil
}

func redisUnlock(ctx context.Context, appId, ucLabel, lock string) (string, error) {
	key := fmt.Sprintf("%s%s%s", appId, keyLock, ucLabel)
	ret := redis.RDS.Get(ctx, key)
	if ret.Err() != nil {
		return "", errors.UserStoreError.Fill(ret.Err().Error())
	}

	if ret.Val() == lock {
		redis.RDS.Del(ctx, key)
	}

	return ret.Val(), nil
}

func setupUserConn(ctx context.Context, uc *UserConnection) (string, error) {
	key := fmt.Sprintf("%s%s%s", uc.AppId, keyUserConn, uc.Label())

	js, err := json.Marshal(uc)
	if err != nil {
		return "", errors.UserStoreError.Fill(err.Error())
	}

	ret := redis.RDS.Set(ctx, key, string(js), time.Minute)
	if ret.Err() != nil {
		return "", errors.UserStoreError.Fill(ret.Err().Error())
	}

	return ret.Val(), nil
}

func refreshUserConn(ctx context.Context, uc *UserConnection) (bool, error) {
	key := fmt.Sprintf("%s%s%s", uc.AppId, keyUserConn, uc.Label())
	ret := redis.RDS.Expire(ctx, key, time.Minute)
	if ret.Err() != nil {
		return false, errors.UserRefreshError.Fill(ret.Err().Error())
	}
	return ret.Val(), nil

}

func setupUserClients(ctx context.Context, uc *UserConnection) (int64, error) {
	key := fmt.Sprintf("%s%s%d", uc.AppId, keyUserConn, uc.UserId)

	js, err := json.Marshal(uc)
	if err != nil {
		return 0, errors.UserStoreError.Fill(err.Error())
	}

	ret := redis.RDS.HSet(ctx, key, uc.Label(), string(js))
	if ret.Err() != nil {
		return 0, errors.UserStoreError.Fill(ret.Err().Error())
	}

	//TODO... 过期时间没设置

	return ret.Val(), nil
}
