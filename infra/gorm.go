package infra

import (
	"context"
	"fmt"
	"github.com/magicnana999/im/global"
	"github.com/magicnana999/im/pkg/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"time"
)

const (
	DefaultDSN             = "root:root@tcp(127.0.0.1:3306)/im?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai"
	DefaultMaxOpenConns    = 200                    // 高并发，增加连接数
	DefaultMaxIdleConns    = 50                     // 适量空闲，减少创建开销
	DefaultConnMaxLifetime = 30 * time.Minute       // 避免连接老化
	DefaultConnMaxIdleTime = 10 * time.Minute       // 空闲连接回收
	DefaultSlowThreshold   = 200 * time.Millisecond // 慢查询阈值
	DefaultConnTimeout     = 1 * time.Second        // 连接超时
)

// getOrDefaultConfig 返回 GORM 配置，优先使用全局配置，缺失时应用默认值。
// 不会修改输入的 global.Config。
func getOrDefaultGormConfig(g *global.Config) *global.GormConfig {

	c := &global.GormConfig{}
	if g != nil && g.Gorm != nil {
		*c = *g.Gorm
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

	c.SkipDefaultTransaction = true
	c.PrepareStmt = true
	c.DisableNestedTransaction = true

	c.Logger = newGormLogger(c.SlowThreshold)

	return c
}

// NewGorm 初始化 GORM 数据库连接。
// 使用 global.Config 提供配置，通过 fx.Lifecycle 管理生命周期。
// 返回已配置的 gorm.DB 实例和错误（如果有）。
func NewGorm(g *global.Config, lc fx.Lifecycle) (*gorm.DB, error) {

	log := logger.Named("gorm")

	c := getOrDefaultGormConfig(g)

	db, err := gorm.Open(mysql.Open(c.Dsn), c)
	if err != nil {
		log.Error("failed to open gorm",
			zap.Error(err))
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Error("failed to get sql db",
			zap.Error(err))
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.ConnTimeout)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		log.Error("failed to ping sql db",
			zap.Error(err))
		return nil, err
	}

	sqlDB.SetMaxOpenConns(c.MaxOpenConns)
	sqlDB.SetMaxIdleConns(c.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(c.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(c.ConnMaxIdleTime)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("gorm connected")
			return nil
		},
		OnStop: func(ctx context.Context) error {

			if err := sqlDB.Close(); err != nil {
				log.Error("failed to close gorm",
					zap.Error(err))
				return err
			}

			log.Info("gorm closed")
			return nil
		},
	})
	return db, nil
}

type GormLogger struct {
	logger      *logger.Logger
	slowSQL     time.Duration
	sourceField string
}

func newGormLogger(slowThreshold time.Duration) *GormLogger {
	return &GormLogger{
		logger:      logger.Named("gorm"),
		slowSQL:     slowThreshold,
		sourceField: "gorm",
	}
}

func (l *GormLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	return l
}

func (l *GormLogger) Info(_ context.Context, msg string, data ...interface{}) {
	l.logger.Info(fmt.Sprintf(msg, data...))
}

func (l *GormLogger) Warn(_ context.Context, msg string, data ...interface{}) {
	l.logger.Warn(fmt.Sprintf(msg, data...))
}

func (l *GormLogger) Error(_ context.Context, msg string, data ...interface{}) {
	l.logger.Error(fmt.Sprintf(msg, data...))
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	fields := []zap.Field{
		zap.String("sql", sql),
		zap.Int64("rows", rows),
		zap.Duration("elapsed", elapsed),
		zap.String("source", l.sourceField),
	}
	if err != nil {
		fields = append(fields, zap.Error(err))
		l.logger.Error("SQL 执行错误", fields...)
		return
	}
	if elapsed > l.slowSQL {
		l.logger.Warn(fmt.Sprintf("慢 SQL >= %v", l.slowSQL), fields...)
	}
}
