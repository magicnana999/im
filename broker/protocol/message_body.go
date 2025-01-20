package protocol

const (
	MText  string = "text"
	MImage string = "image"
	MAudio string = "audio"
	MVideo string = "video"
)

const (
	TSingle int32 = iota + 1
	TGroup
)

type MessageBody struct {
	MType    string   `json:"mType"`
	CId      string   `json:"cId"`
	To       string   `json:"to"`
	GroupId  string   `json:"groupId"`
	TType    int32    `json:"tType"`
	Sequence int64    `json:"sequence"`
	Content  any      `json:"content"`
	At       []*At    `json:"at"`
	Refer    []*Refer `json:"refer"`
}

type At struct {
	UserId string `json:"userId"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type Refer struct {
	UserId  string `json:"userId"`
	Name    string `json:"name"`
	Avatar  string `json:"avatar"`
	MType   string `json:"mType"`
	Content any    `json:"body"`
}

type TextBody struct {
	Text string `json:"text"`
}

type ImageBody struct {
	Url    string `json:"url"`
	Width  int32  `json:"width"`
	Height int32  `json:"height"`
}

type AudioBody struct {
	Url    string `json:"url"`
	Length int32  `json:"length"`
}

type VideoBody struct {
	Url    string `json:"url"`
	Cover  string `json:"cover"`
	Length int32  `json:"length"`
	Width  int32  `json:"width"`
	Height int32  `json:"height"`
}
