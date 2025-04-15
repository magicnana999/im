package infra

import (
	"context"
	"github.com/magicnana999/im/global"
	"github.com/magicnana999/im/pkg/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

const (
	DefaultDSN             = "root:root1234@tcp(127.0.0.1:3306)/im?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai"
	DefaultMaxOpenConns    = 200                    // 高并发，增加连接数
	DefaultMaxIdleConns    = 50                     // 适量空闲，减少创建开销
	DefaultConnMaxLifetime = 30 * time.Minute       // 避免连接老化
	DefaultConnMaxIdleTime = 10 * time.Minute       // 空闲连接回收
	DefaultSlowThreshold   = 200 * time.Millisecond // 慢查询阈值
	DefaultConnTimeout     = 5 * time.Second        // 连接超时
)

type GormConfig struct {
	gorm.Config
	Dsn             string        `yaml:"dsn"`
	MaxOpenConns    int           `yaml:"maxOpenConns"`
	MaxIdleConns    int           `yaml:"maxIdleConns"`
	ConnMaxLifetime time.Duration `yaml:"connMaxLifetime"`
	ConnMaxIdleTime time.Duration `yaml:"connMaxIdleTime"`
	SlowThreshold   time.Duration `yaml:"slowThreshold"` // 慢查询阈值
	ConnTimeout     time.Duration `yaml:"connTimeout"`   // 连接超时
}

func getOrDefaultConfig(global *global.Config) *GormConfig {

	var c *GormConfig
	if global == nil || global.Gorm == nil {
		c = &GormConfig{}
	} else {
		c = global.Gorm
	}

	if c.Dsn == "" {
		c.Dsn = DefaultDSN
	}

	if c.MaxOpenConns <= 0 {
		c.MaxOpenConns = DefaultMaxOpenConns
	}
	if c.MaxIdleConns <= 0 {
		c.MaxIdleConns = DefaultMaxIdleConns
	}
	if c.MaxIdleConns > c.MaxOpenConns {
		c.MaxIdleConns = c.MaxOpenConns
	}
	if c.ConnMaxLifetime == 0 {
		c.ConnMaxLifetime = DefaultConnMaxLifetime
	}
	if c.ConnMaxIdleTime == 0 {
		c.ConnMaxIdleTime = DefaultConnMaxIdleTime
	}
	if c.SlowThreshold == 0 {
		c.SlowThreshold = DefaultSlowThreshold
	}
	if c.ConnTimeout == 0 {
		c.ConnTimeout = DefaultConnTimeout
	}

	c.Config = gorm.Config{
		SkipDefaultTransaction:   true, // 禁用默认事务，高并发
		PrepareStmt:              true, // 缓存 SQL
		DisableNestedTransaction: true, // 禁用嵌套事务，简化逻辑
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 表名单数，如 message
		},
		// GORM 日志
		Logger: logger.New(
			zap.NewStdLog(logger.Named("gorm")),
			logger.Config{
				SlowThreshold:             c.SlowThreshold,
				LogLevel:                  logger.Warn,
				IgnoreRecordNotFoundError: true,
				ParameterizedQueries:      true,
			},
		),
	}
	return c
}

func NewGorm(g *global.Config, lc fx.Lifecycle) *gorm.DB {

	log := logger.Named("gorm")

	c := getOrDefaultConfig(g)

	db, err := gorm.Open(mysql.Open(c.Dsn), c)
	if err != nil {
		log.Fatal("failed to open gorm",
			zap.String(logger.Op, logger.OpInit),
			zap.Error(err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("failed to get sql db",
			zap.String(logger.Op, logger.OpInit),
			zap.Error(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		log.Fatal("failed to ping sql db",
			zap.String(logger.Op, logger.OpInit),
			zap.Error(err))
	}

	sqlDB.SetMaxOpenConns(c.MaxOpenConns)
	sqlDB.SetMaxIdleConns(c.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(c.ConnMaxLifetime)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("gorm connected",
				zap.String(logger.Op, logger.OpInit))
			return nil
		},
		OnStop: func(ctx context.Context) error {

			if err := sqlDB.Close(); err != nil {
				log.Error("failed to close gorm",
					zap.String(logger.Op, logger.OpInit),
					zap.Error(err))
				return err
			}

			log.Info("gorm closed", zap.String(logger.Op, logger.OpInit))
			return nil
		},
	})
	return db
}
