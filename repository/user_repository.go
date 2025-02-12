package repository

import "github.com/magicnana999/im/entity"

func SelectUserByAppIdUserId(appId string, userId int64) (*entity.User, error) {
	var user entity.User
	err := db.Where("app_id = ? AND user_id = ?", appId, userId).First(&user).Error
	return &user, err
}

func InsertUser(user *entity.User) error {
	return db.Create(user).Error
}

func UpdateUser(user *entity.User) error {
	return db.Save(user).Error
}

func DeleteUser(user *entity.User) error {
	return db.Delete(user).Error
}
