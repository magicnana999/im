package repository

import (
	"github.com/magicnana999/im/conf"
	"github.com/magicnana999/im/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func initGorm() {

	dsn := conf.Global.Mysql.String()

	d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.FatalF("Failed to connect to MySQL:%v", err)
	}
	db = d
}
