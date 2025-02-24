package entity

import "time"

type User struct {
	UserId     int64     `gorm:"primaryKey",json:"userId"`
	AppId      string    `gorm:"primaryKey",json:"appId"`
	UserName   string    `gorm:"user_name",json:"userName"`
	Deleted    int       `gorm:"deleted",json:"deleted"`
	Os         string    `gorm:"os",json:"os"`
	PushEnable int       `gorm:"push_enable",json:"pushEnable"`
	DeviceId   string    `gorm:"device_id",json:"deviceId"`
	IsLogin    int       `gorm:"is_login",json:"isLogin"`
	LastLogin  time.Time `gorm:"last_login",json:"lastLogin"`
	CreatedAt  time.Time `gorm:"created_at",json:"createdAt"`
	UpdatedAt  time.Time `gorm:"updated_at",json:"updatedAt"`
}

func (User) TableName() string {
	return "im_user"
}
