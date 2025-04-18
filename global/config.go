package global

import (
	"github.com/go-redis/redis/v8"
	"github.com/magicnana999/im/pkg/logger"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
	"os"
	"time"
)

type Config struct {
	Gorm  *GormConfig  `yaml:"gorm"`
	Redis *RedisConfig `yaml:"redis"`
	Kafka *KafkaConfig `yaml:"kafka"`
	Etcd  *EtcdConfig  `yaml:"etcd"`
	HTS   *HTSConfig   `yaml:"hts"`
	MRS   *MRSConfig   `yaml:"mrs"`
	MSS   *MSSConfig   `yaml:"mss"`
	RBS   *RBSConfig   `yaml:"rbs"`
}

type RBSConfig struct {
	Network   string `yaml:"network"`
	Addr      string `yaml:"addr"`
	DebugMode bool   `yaml:"debugMode"`
}
type MSSConfig struct {
	MaxRemaining int  `yaml:"maxRemaining"`
	DebugMode    bool `yaml:"debugMode"`
}

type MRSConfig struct {
	Interval  time.Duration `yaml:"interval"` //重发间隔
	Timeout   time.Duration `yaml:"timeout"`  //重发超时时间
	DebugMode bool          `yaml:"debugMode"`
}

type HTSConfig struct {
	Interval  time.Duration `yaml:"interval"` //心跳间隔
	Timeout   time.Duration `yaml:"timeout"`  //心跳超时时间
	DebugMode bool          `yaml:"debugMode"`
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

type EtcdConfig struct {
	Endpoints   []string      `yaml:"endpoints"`
	DialTimeout time.Duration `yaml:"dial-timeout"`
}

type RedisConfig struct {
	Addr     string        `yaml:"addr"`     // Redis 地址，如 "127.0.0.1:6379"
	Password string        `yaml:"password"` // Redis 密码
	DB       int           `yaml:"db"`       // Redis 数据库编号
	Timeout  time.Duration `yaml:"timeout"`  // 连接超时，可选
}

func (c *RedisConfig) ToOptions() *redis.Options {
	return &redis.Options{
		Addr:        c.Addr,
		Password:    c.Password,
		DB:          c.DB,
		DialTimeout: c.Timeout,
	}
}

func Load(path string) (*Config, error) {

	log := logger.Named("global")

	c := &Config{}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Error("load config failed", zap.Error(err))
		return nil, err
	}

	err = yaml.Unmarshal(data, &c)
	if err != nil {
		log.Error("unmarshal config failed", zap.Error(err))
		return nil, err
	}

	return c, nil
}
