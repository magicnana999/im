package state

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/magicnana999/im/common/enum"
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
	KeyBrokerInfo        string        = "im:broker:"
	KeyBrokerConnections               = "im:broker:connections:"
	KeyUserClients       string        = "im:user:clients:"
	Expire               time.Duration = 60 * time.Second
)

const (
	CurrentUser string = `CurrentUser`
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

var (
	m  map[string]*UserConnection
	mu sync.RWMutex
)

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

func (u *UserConnection) Store(appId string, userId int64) error {
	mu.Lock()
	defer mu.Unlock()
	u.AppId = appId
	u.UserId = userId
	m[u.Label()] = u
	return nil
}

func Load(ucLabel string) (*UserConnection, error) {
	mu.RLock()
	defer mu.RUnlock()
	uc := m[ucLabel]
	if uc == nil {
		return nil, errors.UserConnectionNotHeld
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

	logger.DebugF("BrokerInfo refresh,key:%s,result:%t",
		key,
		ret.Val())
	return ret.Val(), ret.Err()
}
