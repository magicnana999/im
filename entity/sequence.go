package entity

import "time"

type Sequence struct {
	UserId     int64      `gorm:"primaryKey",json:"userId"`
	AppId      string     `gorm:"primaryKey",json:"appId"`
	Sequence   int64      `gorm:"sequence",json:"sequence"`
	Batch      int64      `gorm:"batch",json:"batch"`
	CreateTime *time.Time `gorm:"default:null,create_time",json:"createTime"`
	UpdateTime *time.Time `gorm:"default:null,update_time",json:"updateTime"`
}

func (Sequence) TableName() string {
	return "im_sequence"
}
