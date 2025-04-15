package domain

type BrokerInfo struct {
	Addr    string `json:"addr"`
	StartAt int64  `json:"startAt"`
}
