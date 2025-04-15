package global

import "C"
import (
	"fmt"
	"github.com/magicnana999/im/broker"
	"github.com/magicnana999/im/infra"
	"github.com/magicnana999/im/pkg/ip"
	"github.com/magicnana999/im/pkg/logger"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"os"
)

const (
	Scope = "config"
	OP    = "load"
)

var c = &Config{}

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
	Gorm         *infra.GormConfig  `yaml:"gorm"`
	Redis        *infra.RedisConfig `yaml:"redis"`
	Kafka        *infra.KafkaConfig `yaml:"kafka"`
	Etcd         *infra.EtcdConfig  `yaml:"etcd"`
	HTS          *broker.HTSConfig  `yaml:"hts"`
	MSS          *broker.MSSConfig  `yaml:"mss"`
}

func Load(path string) (*Config, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		logger.Fatal("load config failed",
			zap.String(logger.SCOPE, Scope),
			zap.String(logger.OP, OP),
			zap.Error(err))
	}

	err = yaml.Unmarshal(data, &c)
	if err != nil {
		logger.Fatal("unmarshal config failed",
			zap.String(logger.SCOPE, Scope),
			zap.String(logger.OP, OP),
			zap.Error(err))
	}

	if c.Broker.Addr == "" {
		i, e := ip.GetLocalIP()
		if e != nil {
			logger.Fatal("get localIP failed",
				zap.String(logger.SCOPE, Scope),
				zap.String(logger.OP, OP),
				zap.Error(err))

		}

		c.Broker.Addr = fmt.Sprintf("%s:7539", i)
	}

	return c, nil
}

func GetMicroService() MicroService {
	return c.MicroService
}

func GetGorm() *infra.GormConfig {
	return c.Gorm
}

func GetRedis() *infra.RedisConfig {
	return c.Redis
}

func GetKafka() *infra.KafkaConfig {
	return c.Kafka
}

func GetBroker() Broker {
	return c.Broker
}

func GetEtcd() *infra.EtcdConfig {
	return c.Etcd
}

func GetHTS() *broker.HTSConfig {
	return c.HTS
}

func GetMSS() *broker.MSSConfig {
	return c.MSS
}
