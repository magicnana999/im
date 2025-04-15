package global

import (
	"fmt"
	"github.com/magicnana999/im/broker"
	"github.com/magicnana999/im/infra"
	"github.com/magicnana999/im/pkg/ip"
	"github.com/magicnana999/im/pkg/logger"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
	"os"
	"time"
)

type MicroService struct {
	EtcdAddr   []string `yaml:"etcdAddr"`
	BrokerName string   `yaml:"brokerName"`
	BrokerAddr string   `yaml:"brokerAddr"`
	RouterName string   `yaml:"routerName"`
	RouterAddr string   `yaml:"routerAddr"`
	ServerName string   `yaml:"serverName"`
	ServerAddr string   `yaml:"serverAddr"`
}

type Broker struct {
	Name              string `yaml:"name"`
	Addr              string `yaml:"addr"`
	TickerInterval    string `yaml:"tickerInterval"`
	HeartbeatInterval int64  `yaml:"heartbeatInterval"`
}

type Config struct {
	Broker       Broker             `yaml:"broker"`
	MicroService MicroService       `yaml:"microService"`
	Gorm         *GormConfig        `yaml:"gorm"`
	Redis        *infra.RedisConfig `yaml:"redis"`
	Kafka        *KafkaConfig       `yaml:"kafka"`
	Etcd         *infra.EtcdConfig  `yaml:"etcd"`
	HTS          *broker.HTSConfig  `yaml:"hts"`
	MSS          *broker.MSSConfig  `yaml:"mss"`
}

// GormConfig infra.NewGorm()时使用
type GormConfig struct {
	gorm.Config                   //继承
	Dsn             string        `yaml:"dsn"`             //链接字符串
	MaxOpenConns    int           `yaml:"maxOpenConns"`    //最大连接数
	MaxIdleConns    int           `yaml:"maxIdleConns"`    //最大空闲数
	ConnMaxLifetime time.Duration `yaml:"connMaxLifetime"` //最大存活时间
	ConnMaxIdleTime time.Duration `yaml:"connMaxIdleTime"` //最大空闲时间
	SlowThreshold   time.Duration `yaml:"slowThreshold"`   // 慢查询阈值
	ConnTimeout     time.Duration `yaml:"connTimeout"`     // 连接超时
}

// KafkaConfig infra.NewKafkaProducer() 时使用
type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
}

func Load(path string) (*Config, error) {

	log := logger.Named("global")

	c := &Config{}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Error("load config failed",
			zap.String(logger.Op, logger.OpInit),
			zap.Error(err))
		return nil, err
	}

	err = yaml.Unmarshal(data, &c)
	if err != nil {
		log.Error("unmarshal config failed",
			zap.String(logger.Op, logger.OpInit),
			zap.Error(err))
		return nil, err
	}

	if c.Broker.Addr == "" {
		i, err := ip.GetLocalIP()
		if err != nil {
			log.Error("get localIP failed",
				zap.String(logger.Op, logger.OpInit),
				zap.Error(err))
			return nil, err

		}

		c.Broker.Addr = fmt.Sprintf("%s:7539", i)
	}

	return c, nil
}
