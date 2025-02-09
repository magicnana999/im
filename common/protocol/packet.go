package protocol

import (
	"github.com/magicnana999/im/common/enum"
)

type Packet struct {
	Id      string `json:"id"`
	AppId   string `json:"appId"`
	UserId  int64  `json:"userId"`
	Flow    int32  `json:"flow"`
	NeedAck int32  `json:"needAck"`
	CTime   int64  `json:"cTime"`
	STime   int64  `json:"sTime"`
	BType   int32  `json:"bType"`
	Body    any    `json:"body"`
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

type MessageBody struct {
	CType      string   `json:"mType"`
	CId        string   `json:"cId"`
	To         string   `json:"to"`
	GroupId    string   `json:"groupId"`
	TargetType int32    `json:"targetType"`
	Sequence   int64    `json:"sequence"`
	Content    any      `json:"content"`
	At         []*At    `json:"at"`
	Refer      []*Refer `json:"refer"`
}

type At struct {
	UserId int64  `json:"userId"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type Refer struct {
	UserId  int64  `json:"userId"`
	Name    string `json:"name"`
	Avatar  string `json:"avatar"`
	CType   string `json:"mType"`
	Content any    `json:"body"`
}

type TextContent struct {
	Text string `json:"text"`
}

type ImageContent struct {
	Url    string `json:"url"`
	Width  int32  `json:"width"`
	Height int32  `json:"height"`
}

type AudioContent struct {
	Url    string `json:"url"`
	Length int32  `json:"length"`
}

type VideoContent struct {
	Url    string `json:"url"`
	Cover  string `json:"cover"`
	Length int32  `json:"length"`
	Width  int32  `json:"width"`
	Height int32  `json:"height"`
}

type CommandBody struct {
	CType   string `json:"cType"`
	Request any    `json:"request"`
	Reply   any    `json:"reply"`
}

type LoginRequest struct {
	AppId        string      `json:"appId"`
	UserSig      string      `json:"userSig"`
	Version      string      `json:"version"`
	OS           enum.OSType `json:"os"`
	PushDeviceId string      `json:"pushDeviceId"`
}

type LoginReply struct {
	AppId  string `json:"appId"`
	UserId int64  `json:"userId"`
}
