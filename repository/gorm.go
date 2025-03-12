package repository

import (
	"github.com/magicnana999/im/conf"
	"github.com/magicnana999/im/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

var DB *gorm.DB
var _dbOnce sync.Once

func InitGorm(config *gorm.Config) *gorm.DB {
	_dbOnce.Do(func() {

		dsn := conf.Global.Mysql.String()

		d, err := gorm.Open(mysql.Open(dsn), config)
		if err != nil {
			logger.Z.Fatalf("Failed to connect to MySQL:%v", err)
		}
		DB = d
	})
	return DB
}
