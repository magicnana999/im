package service

import "github.com/magicnana999/im/broker/enum"

type UserLogin struct {
	AppId   string      `json:"appId"`
	UserId  string      `json:"userId"`
	Addr    string      `json:"addr"`
	Os      enum.OSType `json:"os"`
	Version string      `json:"version"`
}

func Login(userLogin UserLogin) error {
	return nil
}
