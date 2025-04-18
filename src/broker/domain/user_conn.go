package domain

import (
	"fmt"
	"github.com/magicnana999/im/define"
	"io"
	"strconv"
	"sync/atomic"
)

type UserConn struct {
	Fd            int           `json:"fd"`
	AppId         string        `json:"appId"`
	UserId        int64         `json:"userId"`
	ClientAddr    string        `json:"clientAddr"`
	BrokerAddr    string        `json:"brokerAddr"`
	OS            define.OSType `json:"os"`
	ConnectTime   int64         `json:"connectTime"`
	IsLogin       atomic.Bool   `json:"isLogin"`
	IsClosed      atomic.Bool   `json:"isClosed"`
	LastHeartbeat atomic.Int64  `json:"lastHeartbeat"`
	Reader        io.Reader     `json:"-"`
	Writer        io.Writer     `json:"-"`
	ConnLabel     string        `json:"connLabel"`
	ConnDesc      string        `json:"connDesc"`
}

func (u *UserConn) Close() {
	u.IsClosed.Store(true)
}

func (u *UserConn) Login(appId string, userId int64, os define.OSType) (bool, error) {
	u.IsLogin.Store(true)
	u.AppId = appId
	u.UserId = userId
	u.OS = os
	u.ConnLabel = u.parseLabel()
	u.ConnDesc = u.parseDesc()
	return true, nil
}

func (u *UserConn) parseDesc() string {
	return fmt.Sprintf("%s#%s", u.ClientAddr, u.parseLabel())
}

func (u *UserConn) parseLabel() string {
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

	return fmt.Sprintf("%s#%s#%s", u.AppId, strconv.FormatInt(u.UserId, 10), dt.String())
}

func (u *UserConn) Label() string {
	if u.IsLogin.Load() {
		return u.ConnLabel
	} else {
		return u.parseLabel()
	}
}

func (u *UserConn) Desc() string {
	if u.IsLogin.Load() {
		return u.ConnDesc
	} else {
		return u.parseDesc()
	}
}

func (u *UserConn) RefreshHeartbeat(milli int64) {
	u.LastHeartbeat.Store(milli)
}
