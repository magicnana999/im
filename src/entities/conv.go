package entity

import "time"

type Conv struct {
	ConvId      string    `gorm:"primaryKey",json:"convId"`
	AppId       string    `gorm:"primaryKey",json:"appId"`
	ConvType    string    `gorm:"conv_type",json:"convType"`
	Sequence    int64     `gorm:"sequence",json:"sequence"`
	ReadSeq     int64     `gorm:"read_seq",json:"readSeq"`
	LastMsgId   string    `gorm:"last_msg_id",json:"lastMsgId"`
	LastMsgBody string    `gorm:"last_msg_body",json:"lastMsgBody"`
	IsHide      int       `gorm:"is_hide",json:"isHide"`
	IsTop       int       `gorm:"is_top",json:"isTop"`
	IsDisturb   int       `gorm:"is_disturb",json:"isDisturb"`
	CustomType  string    `gorm:"custom_type",json:"customType"`
	Custom1     string    `gorm:"custom_1",json:"custom1"`
	Custom2     string    `gorm:"custom_2",json:"custom2"`
	CreatedAt   time.Time `gorm:"created_at",json:"createdAt"`
	UpdatedAt   time.Time `gorm:"updated_at",json:"updatedAt"`
}

func (Conv) TableName() string {
	return "im_conv"
}
