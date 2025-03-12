package entity

import "time"

type MessageOffline struct {
	Message
	CreatedAt time.Time `gorm:"created_at",json:"createdAt"`
	UpdatedAt time.Time `gorm:"updated_at",json:"updatedAt"`
}

func (MessageOffline) TableName() string {
	return "im_message_offline"
}
