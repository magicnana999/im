package userstate

import "github.com/magicnana999/im/broker/enum"

type UserClient struct {
	AppId       string      `json:"appId"`
	UserId      uint64      `json:"userId"`
	ClientAddr  string      `json:"clientAddr"`
	BrokerAddr  string      `json:"brokerAddr"`
	OS          enum.OSType `json:"os"`
	ConnectTime uint64      `json:"connectTime"`
}
