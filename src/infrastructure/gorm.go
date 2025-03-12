package infrastructure

import (
	"fmt"
	"github.com/magicnana999/im/pkg/logger"
	"github.com/magicnana999/im/pkg/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"sync"
)

var DB *gorm.DB
var dbOnce sync.Once

type GormConfig struct {
	gorm.Config
	User      string `yaml:"user"`
	Password  string `yaml:"password"`
	Host      string `yaml:"host"`
	Port      string `yaml:"port"`
	Schema    string `yaml:"schema"`
	Charset   string `yaml:"charset"`
	ParseTime string `yaml:"parseTime"`
	Location  string `yaml:"location"`
}

func (m *GormConfig) String() string {
	format := "%s:%s@tcp(%s:%s)/%s"
	first := fmt.Sprintf(format, m.User, m.Password, m.Host, m.Port, m.Schema)

	s := make([]string, 0)
	if utils.IsNotBlank(m.Charset) {
		s = append(s, fmt.Sprintf("?charset=%s", m.Charset))
	}
	if utils.IsNotBlank(m.ParseTime) {
		s = append(s, fmt.Sprintf("parseTime=%s", m.ParseTime))
	}
	if utils.IsNotBlank(m.Location) {
		s = append(s, fmt.Sprintf("loc=%s", m.Location))
	}
	second := strings.Join(s, "&")
	return first + second
}

func InitGorm(config *GormConfig) *gorm.DB {

	if config == nil {
		logger.Fatalf("gorm configuration not found")
	}

	dbOnce.Do(func() {

		dsn := config.String()

		d, err := gorm.Open(mysql.Open(dsn), config)
		if err != nil {
			logger.Fatalf("Failed to connect to MySQL:%v", err)
		}
		DB = d
	})

	return DB

}
