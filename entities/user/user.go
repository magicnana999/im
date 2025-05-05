package user

import (
	"time"
)

// User 用户表
type User struct {
	AppID        string    `gorm:"primaryKey;column:app_id;size:50;comment:租户 ID"`
	UserID       uint64    `gorm:"primaryKey;column:user_id;comment:用户 ID，唯一标识"`
	Username     string    `gorm:"uniqueIndex:idx_app_username;column:username;size:50;not null;comment:用户名"`
	Nickname     string    `gorm:"column:nickname;size:50;comment:昵称"`
	PhoneNumber  string    `gorm:"uniqueIndex:idx_app_phone_number;column:phone_number;size:20;comment:手机号码"`
	PasswordHash string    `gorm:"column:password_hash;size:256;not null;comment:密码哈希"`
	Status       string    `gorm:"column:status;size:20;default:active;not null;comment:用户状态"`
	CreatedAt    time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;comment:创建时间"`
	UpdatedAt    time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;autoUpdateTime;comment:更新时间"`
}

func (User) TableName() string {
	return "im_user"
}
