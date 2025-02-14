package entity

import "time"

type User struct {
	UserId       int64     `gorm:"primaryKey",json:"userId"`
	AppId        string    `gorm:"primaryKey",json:"appId"`
	UserName     string    `gorm:"user_name",json:"userName"`
	Deleted      int       `gorm:"deleted",json:"deleted"`
	Os           int       `gorm:"os",json:"os"`
	PushEnable   int       `gorm:"push_enable",json:"pushEnable"`
	PushDeviceId string    `gorm:"push_device_id",json:"pushDeviceId"`
	CreatedAt    time.Time `gorm:"created_at",json:"createdAt"`
	UpdatedAt    time.Time `gorm:"updated_at",json:"updatedAt"`
}

func (User) TableName() string {
	return "im_user"
}
