package state

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/magicnana999/im/broker/enum"
	"github.com/magicnana999/im/logger"
	"github.com/magicnana999/im/redis"
	"github.com/magicnana999/im/util/str"
	"github.com/panjf2000/gnet/v2"
	"time"
)

const (
	KeyBrokerInfo        string        = "im:broker:"
	KeyBrokerConnections               = "im:broker:connections:"
	KeyUserClients       string        = "im:user:clients:"
	Expire               time.Duration = 60 * time.Second
)

type BrokerInfo struct {
	Addr    string `json:"addr"`
	StartAt int64  `json:"startAt"`
}

type UserConnection struct {
	Fd          int         `json:"fd"`
	AppId       string      `json:"appId"`
	UserId      string      `json:"userId"`
	ClientAddr  string      `json:"clientAddr"`
	BrokerAddr  string      `json:"brokerAddr"`
	OS          enum.OSType `json:"os"`
	ConnectTime int64       `json:"connectTime"`
}

func EmptyUserConnection(c gnet.Conn) UserConnection {
	return UserConnection{
		Fd:          c.Fd(),
		AppId:       "",
		UserId:      "",
		ClientAddr:  c.RemoteAddr().String(),
		BrokerAddr:  c.LocalAddr().String(),
		OS:          enum.OSType(0),
		ConnectTime: time.Now().UnixMilli(),
	}
}

func (u UserConnection) Label() string {
	if str.IsBlank(u.UserId) {
		return ""
	}

	dt := u.OS.GetDeviceType()
	if !dt.Valid() {
		return ""
	}

	if str.IsBlank(dt.Name()) {
		return ""
	}

	return u.UserId + "#" + dt.Name()
}

func SetBroker(ctx context.Context, broker BrokerInfo) (string, error) {
	json, err := json.Marshal(broker)
	if err != nil {
		return "", err
	}

	key := fmt.Sprintf("%s%s", KeyBrokerInfo, broker.Addr)
	ret := redis.RDS.Set(ctx, key, json, Expire)

	logger.InfoF("BrokerInfo set,key:%s,result:%s,error:%v",
		key,
		ret.Val(),
		ret.Err())

	return ret.Val(), ret.Err()
}

func RefreshBroker(ctx context.Context, broker BrokerInfo) (bool, error) {
	key := fmt.Sprintf("%s%s", KeyBrokerInfo, broker.Addr)
	ret := redis.RDS.Expire(ctx, key, Expire)

	//logger.DebugF("BrokerInfo refresh,key:%s,result:%t,error:%v",
	//	key,
	//	ret.Val(),
	//	ret.Err())
	return ret.Val(), ret.Err()
}

//// SetConnection TODO.. 更换为Lua
//func SetConnection(ctx context.Context, uc UserConnection) (any, error) {
//	jsonStr, err := json.Marshal(uc)
//	if err != nil {
//		return nil, err
//	}
//
//	if str.IsBlank(uc.Label()) {
//		return nil, errors.New("empty label")
//	}
//
//	logger.DebugF("[%d] Get UserConnection label:%s", routine.Goid(), label)
//
//	{
//		key := fmt.Sprintf("%s%s", KeyUserClients, uc.UserId)
//		ret := redis.RDS.HSetNX(ctx, key, label, jsonStr)
//
//		logger.InfoF("[%d] Set UserClient [%s:%s:%v]", routine.Goid(), key, label, ret)
//
//		if ret.Err() != nil {
//			return nil, ret.Err()
//		}
//
//		if !ret.Val() {
//			r := redis.RDS.HGet(ctx, key, label)
//
//			logger.InfoF("[%d] Set existed UserConnection [%s:%s:%v]", routine.Goid(), key, label, r)
//
//			if r.Err() != nil {
//				return nil, fmt.Errorf("[%d] Failed to set UserConnection because it already exists, but the existing one could not be found [%s:%s:%v]", routine.Goid(), key, label, r)
//			}
//
//			if str.IsBlank(r.Val()) {
//				return nil, fmt.Errorf("[%d] Failed to set UserConnection because it already exists, but an empty string was found [%s:%s:%v]", routine.Goid(), key, label, r)
//			}
//
//			existed := UserConnection{}
//
//			if e := json.Unmarshal([]byte(r.String()), existed); e != nil {
//				return nil, e
//			}
//
//			return existed, nil
//		}
//	}
//
//	{
//		key := fmt.Sprintf("%s%s", KeyBrokerConnections, uc.BrokerAddr)
//		ret := redis.RDS.HSet(ctx, key, label, jsonStr)
//
//		logger.InfoF("[%d] Set BrokerConnections [%s:%s:%v]", routine.Goid(), key, label, ret)
//
//		if ret.Err() != nil {
//			return nil, ret.Err()
//		}
//	}
//
//	return nil, nil
//}
