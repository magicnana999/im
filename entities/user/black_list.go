package user

import "time"

type Blacklist struct {
	AppID     string    `gorm:"primaryKey;column:app_id;size:50;comment:租户 ID"`
	UserID    uint64    `gorm:"primaryKey;column:user_id;comment:用户 ID"`
	BlockedID uint64    `gorm:"primaryKey;column:blocked_id;comment:被拉黑用户 ID"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;comment:拉黑时间"`
}

func (Blacklist) TableName() string {
	return "im_blacklist"
}
