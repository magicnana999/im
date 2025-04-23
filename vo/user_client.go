package vo

import "go.uber.org/atomic"

type UserClient struct {
	Fd          int           `json:"fd"`
	AppId       atomic.String `json:"appId"`
	UserId      atomic.Int64  `json:"userId"`
	OS          atomic.String `json:"os"`
	ClientAddr  string        `json:"clientAddr"`
	BrokerAddr  string        `json:"brokerAddr"`
	ConnectTime int64         `json:"connectTime"` //首次连接时间 毫秒
	Label       string        `json:"label"`
}
