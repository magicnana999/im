package repository

import (
	"github.com/magicnana999/im/conf"
	"github.com/magicnana999/im/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

var DB *gorm.DB
var lock sync.RWMutex

func InitGorm() *gorm.DB {

	lock.Lock()
	defer lock.Unlock()

	if DB == nil {
		dsn := conf.Global.Mysql.String()

		d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			logger.FatalF("Failed to connect to MySQL:%v", err)
		}
		DB = d
	}

	return DB

}
