package protocol

import (
	"math/rand"
	"time"
)

// Flow
const (
	FlowRequest int32 = iota + 1
	FlowResponse
)

// Type
const (
	TypeHeartbeat int32 = iota + 1
	TypeCommand
	TypeMessage
	TypeNotice
	TypeTips
)

// NeedAck
const (
	NO int32 = iota
	YES
)

var (
	r               *rand.Rand
	HeartbeatPacket *Packet
)

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	HeartbeatPacket = &Packet{
		Type: TypeHeartbeat,
		Body: r.Uint32(),
	}
}

type Packet struct {
	Id      string `json:"id"`
	AppId   string `json:"appId"`
	UserId  string `json:"userId"`
	Flow    int32  `json:"flow"`
	NeedAck int32  `json:"needAck"`
	Type    int32  `json:"type"`
	CTime   int64  `json:"cTime"`
	STime   int64  `json:"sTime"`
	Body    any    `json:"body"`
}

func (p *Packet) IsHeartbeat() bool {
	return p.Type == TypeHeartbeat
}
