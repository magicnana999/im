package domain

import (
	"encoding/json"
	"fmt"
	"github.com/panjf2000/gnet/v2"
	"go.uber.org/atomic"
	"io"
	"time"
)

type UserConn struct {
	Fd            int           `json:"fd"`
	AppId         atomic.String `json:"appId"`
	UserId        atomic.Int64  `json:"userId"`
	OS            atomic.String `json:"os"`
	ClientAddr    string        `json:"clientAddr"`
	BrokerAddr    string        `json:"brokerAddr"`
	ConnectTime   int64         `json:"connectTime"` //首次连接时间 毫秒
	IsLogin       atomic.Bool   `json:"-"`
	IsClosed      atomic.Bool   `json:"-"`
	LastHeartbeat atomic.Time   `json:"-"` //上次心跳 毫秒
	Reader        io.Reader     `json:"-"`
	Conn          gnet.Conn     `json:"-"`
}

func NewUserConn(c gnet.Conn) *UserConn {
	uc := &UserConn{
		Fd:          c.Fd(),
		ClientAddr:  c.RemoteAddr().String(),
		BrokerAddr:  c.LocalAddr().String(),
		ConnectTime: time.Now().UnixMilli(),
		Reader:      c,
		Conn:        c,
	}

	uc.Refresh(time.Now())

	return uc
}
func (u *UserConn) Close() bool {
	return u.IsClosed.CompareAndSwap(false, true)
}

func (u *UserConn) Login(appId string, userId int64, os string) bool {
	if u.IsLogin.CompareAndSwap(false, true) {
		u.AppId.Store(appId)
		u.UserId.Store(userId)
		u.OS.Store(os)
		return true
	}
	return false
}

func (u *UserConn) Desc() string {
	return fmt.Sprintf("%s#%s", u.ClientAddr, u.Label())
}

func (u *UserConn) Label() string {
	if u.AppId.Load() == "" {
		return ""
	}
	if u.UserId.Load() == 0 {
		return ""
	}

	if u.OS.Load() == "" {
		return ""
	}

	return fmt.Sprintf("%s#%s#%s", u.AppId.Load(), u.UserId.String(), u.OS.Load())
}

func (u *UserConn) Refresh(t time.Time) {
	u.LastHeartbeat.Store(t)
}

func (u *UserConn) ToJSON() ([]byte, error) {
	return json.Marshal(u)
}
