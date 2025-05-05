package user

import "time"

type FriendRequest struct {
	AppID      string    `gorm:"primaryKey;column:app_id;size:50;comment:租户 ID"`
	RequestID  uint64    `gorm:"primaryKey;autoIncrement;column:request_id;comment:请求 ID"`
	FromUserID uint64    `gorm:"index:idx_app_from_user_id;column:from_user_id;not null;comment:发起者 ID"`
	ToUserID   uint64    `gorm:"index:idx_app_to_user_id;column:to_user_id;not null;comment:接收者 ID"`
	Status     string    `gorm:"column:status;size:20;default:pending;not null;comment:请求状态"`
	Message    string    `gorm:"column:message;size:200;comment:请求消息"`
	CreatedAt  time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;comment:创建时间"`
	UpdatedAt  time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;autoUpdateTime;comment:更新时间"`
}

func (FriendRequest) TableName() string {
	return "im_friend_request"
}
