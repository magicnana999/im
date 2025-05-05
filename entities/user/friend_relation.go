package user

import "time"

type FriendRelation struct {
	AppID     string    `gorm:"primaryKey;column:app_id;size:50;comment:租户 ID"`
	UserID    uint64    `gorm:"primaryKey;column:user_id;comment:用户 ID"`
	FriendID  uint64    `gorm:"primaryKey;column:friend_id;comment:好友 ID"`
	GroupID   uint64    `gorm:"column:group_id;default:0;comment:分组 ID，0 表示默认分组"`
	Remark    string    `gorm:"column:remark;size:50;comment:好友备注"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;comment:关系建立时间"`
}

func (FriendRelation) TableName() string {
	return "im_friend_relation"
}
