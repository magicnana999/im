package protocol

// Flow
const (
	FlowRequest int32 = iota + 1
	FlowResponse
)

// Type
const (
	TypeCommand int32 = iota + 1
	TypeMessage
	TypeNotice
	TypeTips
)

// NeedAck
const (
	NO int32 = iota
	YES
)

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
