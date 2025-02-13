package repository

import (
	"github.com/magicnana999/im/conf"
	"github.com/magicnana999/im/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

var db *gorm.DB
var lock sync.RWMutex

func initGorm() {

	lock.Lock()
	defer lock.Unlock()

	if db == nil {
		dsn := conf.Global.Mysql.String()

		d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			logger.FatalF("Failed to connect to MySQL:%v", err)
		}
		db = d
	}

}
