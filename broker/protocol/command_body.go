package protocol

import "github.com/magicnana999/im/broker/enum"

const (
	MLogin      string = "login"
	MLoginReply string = "loginReply"
	MLogout     string = "logout"
)

type CommandBody struct {
	MType   string `json:"cType"`
	Token   string `json:"token"`
	Code    int32  `json:"code"`
	Message string `json:"message"`
	Content any    `json:"content"`
}

type LoginContent struct {
	Version      string      `json:"version"`
	OS           enum.OSType `json:"os"`
	DeviceId     string      `json:"deviceId"`
	PushDeviceId string      `json:"pushDeviceId"`
}

type LoginReply struct {
	UserId int64 `json:"userId"`
}

type LogoutContent struct {
}
