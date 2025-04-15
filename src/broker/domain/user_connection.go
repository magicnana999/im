package domain

import (
	"github.com/magicnana999/im/constants"
	"github.com/panjf2000/gnet/v2"
	"strconv"
	"sync/atomic"
)

type UserConnection struct {
	Fd            int              `json:"fd"`
	AppId         string           `json:"appId"`
	UserId        int64            `json:"userId"`
	ClientAddr    string           `json:"clientAddr"`
	BrokerAddr    string           `json:"brokerAddr"`
	OS            constants.OSType `json:"os"`
	ConnectTime   int64            `json:"connectTime"`
	IsLogin       atomic.Bool      `json:"isLogin"`
	IsClose       atomic.Bool      `json:"isClose"`
	C             gnet.Conn        `json:"-"`
	LastHeartbeat atomic.Int64     `json:"lastHeartbeat"`
}

func (u *UserConnection) Label() string {
	if len(u.AppId) == 0 {
		return ""
	}
	if u.UserId == 0 {
		return ""
	}

	dt := u.OS.GetDeviceType()

	if len(dt.String()) == 0 {
		return ""
	}

	return u.AppId + "#" + strconv.FormatInt(u.UserId, 10) + "#" + dt.String()
}
