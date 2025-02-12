package entity

import "time"

type User struct {
	UserId       int64      `gorm:"primaryKey",json:"userId"`
	UserName     string     `gorm:"user_name",json:"userName"`
	AppId        string     `gorm:"primaryKey",json:"appId"`
	Deleted      int        `gorm:"deleted",json:"deleted"`
	Os           int        `gorm:"os",json:"os"`
	PushEnable   int        `gorm:"push_enable",json:"pushEnable"`
	PushDeviceId string     `gorm:"push_device_id",json:"pushDeviceId"`
	CreateTime   *time.Time `gorm:"default:null,create_time",json:"createTime"`
	UpdateTime   *time.Time `gorm:"default:null,update_time",json:"updateTime"`
}

func (User) TableName() string {
	return "im_user"
}
