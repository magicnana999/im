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
	TCP   *TCPConfig   `yaml:"tcp,omitempty" json:"tcp,omitempty"`
	Gorm  *GormConfig  `yaml:"gorm" json:"gorm"`
	Redis *RedisConfig `yaml:"redis" json:"redis"`
	Kafka *KafkaConfig `yaml:"kafka" json:"kafka"`
	Etcd  *EtcdConfig  `yaml:"etcd" json:"etcd"`
	MRS   *MRSConfig   `yaml:"mrs" json:"mrs"`
	MSS   *MSSConfig   `yaml:"mss" json:"mss"`
	RBS   *RBSConfig   `yaml:"rbs" json:"rbs"`
	RRS   *RRSConfig   `yaml:"rrs,omitempty" json:"rrs,omitempty"`
}

type TCPConfig struct {
	Addr      string              `yaml:"addr" json:"addr"`
	Interval  time.Duration       `yaml:"interval" json:"interval"`
	Heartbeat *TcpHeartbeatConfig `yaml:"heartbeat" json:"heartbeat"`
	Worker    *TcpWorkerConfig    `yaml:"worker" json:"worker"`
}

type TcpHeartbeatConfig struct {
	Timeout           time.Duration `yaml:"timeout" json:"timeout"`
	SlotTick          time.Duration `yaml:"slotTick" json:"slotTick"`
	SlotCount         int           `yaml:"slotCount" json:"slotCount"`
	SlotMaxLength     int64         `yaml:"slotMaxLength" json:"slotMaxLength"`
	WorkerCount       int           `yaml:"workerCount" json:"workerCount"`
	WorkerNonBlocking bool          `yaml:"workerNonBlocking" json:"workerNonBlocking"`
	WorkerExpire      time.Duration `yaml:"workerExpire" json:"workerExpire"`
	WorkerPreAlloc    bool          `yaml:"workerPreAlloc" json:"workerPreAlloc"`
}

type TcpWorkerConfig struct {
	Size             int           `yaml:"size" json:"size"`
	ExpireDuration   time.Duration `yaml:"expireDuration" json:"expireDuration"`
	MaxBlockingTasks int           `yaml:"maxBlockingTasks" json:"maxBlockingTasks"`
	Nonblocking      bool          `yaml:"nonblocking" json:"nonblocking"`
	PreAlloc         bool          `yaml:"preAlloc" json:"preAlloc"`
	DisablePurge     bool          `yaml:"disablePurge" json:"disablePurge"`
}

type RBSConfig struct {
	Network   string `yaml:"network" json:"network"`
	Addr      string `yaml:"addr" json:"addr"`
	DebugMode bool   `yaml:"debugMode" json:"debugMode"`
}

type RRSConfig struct {
	Network   string `yaml:"network" json:"network"`
	Addr      string `yaml:"addr" json:"addr"`
	DebugMode bool   `yaml:"debugMode" json:"debugMode"`
}

type MSSConfig struct {
	MaxRemaining int  `yaml:"maxRemaining" json:"maxRemaining"`
	DebugMode    bool `yaml:"debugMode" json:"debugMode"`
}

type MRSConfig struct {
	Timeout   time.Duration `yaml:"timeout" json:"timeout"`
	DebugMode bool          `yaml:"debugMode" json:"debugMode"`
}

type GormConfig struct {
	gorm.Config     `json:"-"`
	Dsn             string        `yaml:"dsn" json:"dsn"`
	MaxOpenConns    int           `yaml:"maxOpenConns" json:"maxOpenConns"`
	MaxIdleConns    int           `yaml:"maxIdleConns" json:"maxIdleConns"`
	ConnMaxLifetime time.Duration `yaml:"connMaxLifetime" json:"connMaxLifetime"`
	ConnMaxIdleTime time.Duration `yaml:"connMaxIdleTime" json:"connMaxIdleTime"`
	SlowThreshold   time.Duration `yaml:"slowThreshold" json:"slowThreshold"`
	ConnTimeout     time.Duration `yaml:"connTimeout" json:"connTimeout"`
}

type KafkaConfig struct {
	Brokers []string `yaml:"brokers" json:"brokers"`
}

type EtcdConfig struct {
	Endpoints   []string      `yaml:"endpoints" json:"endpoints"`
	DialTimeout time.Duration `yaml:"dial-timeout" json:"dial-timeout"`
}

type RedisConfig struct {
	Addr         string        `yaml:"addr" json:"addr"`
	Password     string        `yaml:"password" json:"password"`
	DB           int           `yaml:"db" json:"db"`
	Timeout      time.Duration `yaml:"timeout" json:"timeout"`
	PoolSize     int           `yaml:"poolSize" json:"poolSize"`
	MinIdleConns int           `yaml:"minIdleConns" json:"minIdleConns"`
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
