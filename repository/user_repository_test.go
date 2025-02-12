package repository

import (
	"fmt"
	"github.com/magicnana999/im/conf"
	"github.com/magicnana999/im/entity"
	"testing"
)

func _init() {
	conf.LoadConfig("/Users/jinsong/source/github/im/conf/im-router.yaml")
	initDB()
}

func TestSelectUserByUserId(t *testing.T) {

	_init()

	u, e := SelectUserByAppIdUserId("19860220", 1200120)
	if e != nil {
		t.Error(e)
	}

	fmt.Println(u)

	u2 := &entity.User{
		AppId:  "19860220",
		UserId: 100002,
	}
	e3 := InsertUser(u2)
	if e != nil {
		t.Error(e3)
	}
}
