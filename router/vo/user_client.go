package vo

type UserClient struct {
	Fd          int    `json:"fd"`
	AppId       string `json:"appId"`
	UserId      int64  `json:"userId"`
	OS          string `json:"os"`
	ClientAddr  string `json:"clientAddr"`
	BrokerAddr  string `json:"brokerAddr"`
	ConnectTime int64  `json:"connectTime"` //首次连接时间 毫秒
	Label       string `json:"label"`
}
