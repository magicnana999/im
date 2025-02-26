package svc

import (
	"context"
	"github.com/magicnana999/im/conf"
	"github.com/magicnana999/im/entity"
	"github.com/magicnana999/im/redis"
	"github.com/magicnana999/im/repository"
	"testing"
)

func TestUserSvc_UserSig(t *testing.T) {

	conf.LoadConfig("/Users/jinsong/source/github/im/conf/im-broker.yaml")

	db := repository.InitGorm()
	rds := redis.InitUserStorage()

	var user1 entity.User
	db.Where("app_id = ? and user_id = ?", "19860220", 1200120).First(&user1)

	rds.StoreUserSig(context.Background(), "19860220", &user1)

	var user2 entity.User
	db.Where("app_id = ? and user_id = ?", "19860220", 1200121).First(&user2)

	rds.StoreUserSig(context.Background(), "19860220", &user2)

}
