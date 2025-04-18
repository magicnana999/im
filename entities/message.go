package entity

type Message struct {
	MessageId string `gorm:"primaryKey",json:"id"`
	AppId     string `gorm:"appId",json:"appId"`
	UserId    int64  `gorm:"user_id",json:"appId"`
	To        int64  `gorm:"to",json:"to"`
	GroupId   int64  `gorm:"group_id",json:"groupId"`
	ConvId    string `gorm:"conv_id",json:"convId"`
	Sequence  int64  `gorm:"sequence",json:"sequence"`
	cTime     int64  `gorm:"c_time",json:"cTime"`
	sTime     int64  `gorm:"s_time",json:"sTime"`
	cType     string `gorm:"c_type",json:"cType"`
	At        string `gorm:"at",json:"at"`
	Refer     string `gorm:"refer",json:"refer"`
	Content   string `gorm:"content",json:"content"`
}
