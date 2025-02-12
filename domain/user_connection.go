package domain

import (
	"github.com/magicnana999/im/enum"
	"github.com/panjf2000/gnet/v2"
	"strconv"
)

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

func (u *UserConnection) Label() string {
	if len(u.AppId) == 0 {
		return ""
	}
	if u.UserId == 0 {
		return ""
	}

	dt := u.OS.GetDeviceType()
	if !dt.Valid() {
		return ""
	}

	if len(dt.String()) == 0 {
		return ""
	}

	return u.AppId + "#" + strconv.FormatInt(u.UserId, 10) + "#" + dt.String()
}
