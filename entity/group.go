package entity

import "time"

type Group struct {
	GroupId      int64      `gorm:"primaryKey",json:"groupId"`
	AppId        string     `gorm:"primaryKey",json:"appId"`
	UserId       int64      `gorm:"user_id",json:"userId"`
	GroupName    string     `gorm:"group_name",json:"groupName"`
	GroupType    string     `gorm:"group_type",json:"groupType"`
	GroupAvatar  string     `gorm:"group_avatar",json:"groupAvatar"`
	CustomType   string     `gorm:"custom_type",json:"customType"`
	Introduction string     `gorm:"introduction",json:"introduction"`
	Notification string     `gorm:"notification",json:"notification"`
	Custom1      string     `gorm:"custom1",json:"custom1"`
	Custom2      string     `gorm:"custom2",json:"custom2"`
	CreateTime   *time.Time `gorm:"default:null,create_time",json:"createTime"`
	UpdateTime   *time.Time `gorm:"default:null,update_time",json:"updateTime"`
}

func (Group) TableName() string {
	return "im_group"
}
