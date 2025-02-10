package entity

import "time"

type User struct {
	UserId       int64     `db:"user_id",json:"userId"`
	UserName     string    `db:"user_name",json:"userName"`
	AppId        string    `db:"app_id",json:"appId"`
	Deleted      int       `db:"deleted",json:"deleted"`
	Os           int       `db:"os",json:"os"`
	PushEnable   int       `db:"push_enable",json:"pushEnable"`
	PushDeviceId string    `db:"push_device_id",json:"pushDeviceId"`
	CreateTime   time.Time `db:"create_time",json:"createTime"`
	UpdateTime   time.Time `db:"update_time",json:"updateTime"`
}
