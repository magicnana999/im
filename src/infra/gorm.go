package infra

import (
	"context"
	"github.com/magicnana999/im/global"
	"github.com/magicnana999/im/pkg/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

const (
	Gorm  = "gorm"
	Init  = "init"
	Close = "close"
)

type GormConfig struct {
	gorm.Config
	Dsn             string        `yaml:"dsn"`
	MaxOpenConns    int           `yaml:"maxOpenConns"`
	MaxIdleConns    int           `yaml:"maxIdleConns"`
	ConnMaxLifetime time.Duration `yaml:"connMaxLifetime"`
}

func newGormConfig(c *global.Config) *GormConfig {

	var gc *GormConfig
	if c == nil {
		gc
	}
}

func NewGorm(lc fx.Lifecycle) *gorm.DB {

	c := global.GetGorm()

	if c == nil {
		logger.Fatal("gorm configuration not found",
			zap.String(logger.SCOPE, Gorm),
			zap.String(logger.OP, Init))
	}

	db, err := gorm.Open(mysql.Open(c.Dsn), c)
	if err != nil {
		logger.Fatal("gorm could not be open",
			zap.String(logger.SCOPE, Gorm),
			zap.String(logger.OP, Init),
			zap.Error(err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal("failed to get sql db",
			zap.String(logger.SCOPE, Gorm),
			zap.String(logger.OP, Init),
			zap.Error(err))
	}

	if err := sqlDB.Ping(); err != nil {
		logger.Fatal("failed to ping sql db",
			zap.String(logger.SCOPE, Gorm),
			zap.String(logger.OP, Init),
			zap.Error(err))
	}

	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("gorm connection established",
				zap.String(logger.SCOPE, Gorm),
				zap.String(logger.OP, Init))
			return nil
		},
		OnStop: func(ctx context.Context) error {

			done := make(chan error)
			go func() {
				done <- sqlDB.Close()
			}()

			select {
			case err := <-done:
				if err != nil {
					logger.Error("gorm could not close",
						zap.String(logger.SCOPE, Gorm),
						zap.String(logger.OP, Close),
						zap.Error(err))
					return err
				}
				logger.Info("gorm closed",
					zap.String(logger.SCOPE, Gorm),
					zap.String(logger.OP, Close))
				return nil
			case <-ctx.Done():
				logger.Error("gorm close timeout",
					zap.String(logger.SCOPE, Gorm),
					zap.String(logger.OP, Close),
					zap.Error(ctx.Err()))
				return ctx.Err()
			}
		},
	})
	return db
}
