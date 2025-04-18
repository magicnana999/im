package entity

import "time"

type GroupMember struct {
	GroupId    int64     `gorm:"primaryKey",json:"groupId"`
	AppId      string    `gorm:"primaryKey",json:"appId"`
	UserId     int64     `gorm:"primaryKey",json:"userId"`
	MemberType string    `gorm:"member_type",json:"memberType"`
	Sort       string    `gorm:"sort",json:"sort"`
	Alias      string    `gorm:"alias",json:"alias"`
	Role       string    `gorm:"role",json:"role"`
	Muted      int       `gorm:"muted",json:"muted"`
	MuteUntil  time.Time `gorm:"mute_until",json:"muteUntil"`
	CreatedAt  time.Time `gorm:"created_at",json:"createdAt"`
	UpdatedAt  time.Time `gorm:"updated_at",json:"updatedAt"`
}

func (GroupMember) TableName() string {
	return "im_group_member"
}
