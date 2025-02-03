package repository

import "github.com/magicnana999/im/common/entity"

func SelectUserByUserId(appId string, userId int64) (*entity.User, error) {
	var user entity.User
	err := db.Get(&user, "select * from im_user where app_id = ? and user_id = ?", appId, userId)
	return &user, err
}
