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
	TCP   *TCPConfig   `yaml:"tcp,omitempty"`
	Gorm  *GormConfig  `yaml:"gorm"`
	Redis *RedisConfig `yaml:"redis"`
	Kafka *KafkaConfig `yaml:"kafka"`
	Etcd  *EtcdConfig  `yaml:"etcd"`
	MRS   *MRSConfig   `yaml:"mrs"`
	MSS   *MSSConfig   `yaml:"mss"`
	RBS   *RBSConfig   `yaml:"rbs"`
	RRS   *RRSConfig   `yaml:"rrs,omitempty"`
}

type TCPConfig struct {
	Addr      string              `yaml:"addr",json:"addr"`
	Interval  time.Duration       `yaml:"interval"`
	Heartbeat *TcpHeartbeatConfig `yaml:"heartbeat"`
	Worker    *TcpWorkerConfig    `yaml:"worker"`
}

type TcpHeartbeatConfig struct {
	Timeout           time.Duration `yaml:"timeout"`
	SlotTick          time.Duration `yaml:"slotTick",json:"slotTick"`
	SlotCount         int           `yaml:"slotCount",json:"slotCount"`
	SlotMaxLength     int64         `yaml:"slotMaxLength",json:"slotMaxLength"`
	WorkerCount       int           `yaml:"workerCount",json:"workerCount"`
	WorkerNonBlocking bool          `yaml:"workerNonBlocking",json:"workerNonBlocking"`
	WorkerExpire      time.Duration `yaml:"workerExpire",json:"workerExpire"`
	WorkerPreAlloc    bool          `yaml:"workerPreAlloc",json:"workerPreAlloc"`
}

type TcpWorkerConfig struct {
	Size             int           `yaml:"size"`
	ExpireDuration   time.Duration `yaml:"expireDuration"`
	MaxBlockingTasks int           `yaml:"maxBlockingTasks"`
	Nonblocking      bool          `yaml:"nonblocking",json:"nonblocking"`
	PreAlloc         bool          `yaml:"preAlloc",json:"preAlloc"`
	DisablePurge     bool          `yaml:"disablePurge",json:"disablePurge"`
}

type RBSConfig struct {
	Network   string `yaml:"network"`
	Addr      string `yaml:"addr"`
	DebugMode bool   `yaml:"debugMode"`
}

type RRSConfig struct {
	Network   string `yaml:"network"`
	Addr      string `yaml:"addr"`
	DebugMode bool   `yaml:"debugMode"`
}

type MSSConfig struct {
	MaxRemaining int  `yaml:"maxRemaining"`
	DebugMode    bool `yaml:"debugMode"`
}

type MRSConfig struct {
	Timeout   time.Duration `yaml:"timeout"` //重发超时时间
	DebugMode bool          `yaml:"debugMode"`
}

// GormConfig infra.NewGorm()时使用
type GormConfig struct {
	gorm.Config                   //继承
	Dsn             string        `yaml:"dsn"` //链接字符串
	MaxOpenConns    int           `yaml:"maxOpenConns"` //最大连接数
	MaxIdleConns    int           `yaml:"maxIdleConns"` //最大空闲数
	ConnMaxLifetime time.Duration `yaml:"connMaxLifetime"` //最大存活时间
	ConnMaxIdleTime time.Duration `yaml:"connMaxIdleTime"` //最大空闲时间
	SlowThreshold   time.Duration `yaml:"slowThreshold"` // 慢查询阈值
	ConnTimeout     time.Duration `yaml:"connTimeout"` // 连接超时
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
	Addr         string        `yaml:"addr"`     // Redis 地址，如 "127.0.0.1:6379"
	Password     string        `yaml:"password"` // Redis 密码
	DB           int           `yaml:"db"`       // Redis 数据库编号
	Timeout      time.Duration `yaml:"timeout"`  // 连接超时，可选
	PoolSize     int           `yaml:"poolSize"`
	MinIdleConns int           `yaml:"minIdleConns"`
}

func (c *RedisConfig) ToOptions() *redis.Options {
	return &redis.Options{
		Addr:         c.Addr,
		Password:     c.Password,
		DB:           c.DB,
		DialTimeout:  c.Timeout,
		PoolSize:     c.PoolSize,
		MinIdleConns: c.MinIdleConns,
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
