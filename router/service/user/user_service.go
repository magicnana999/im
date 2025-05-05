package user

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/magicnana999/im/global"
	"github.com/magicnana999/im/infra"
	"github.com/magicnana999/im/pkg/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
	"sync"
	"time"
)

type UserService struct {
	db       *gorm.DB
	redis    *redis.Client
	logger   *logger.Logger
	group    singleflight.Group // 防缓存击穿
	lockPool sync.Pool          // Redis 分布式锁池
}

// NewUserService 初始化用户服务
func NewUserService(g *global.Config, lc fx.Lifecycle) (*UserService, error) {
	db, err := infra.NewGorm(g, lc)
	if err != nil {
		return nil, err
	}

	redis := infra.NewRedisClient(g, lc)

	svc := &UserService{
		db:     db,
		redis:  redis,
		logger: logger.Named("user_service"),
		lockPool: sync.Pool{
			New: func() interface{} {
				return redislock.New(redis)
			},
		},
	}

	// 预加载热点用户到缓存
	svc.preloadHotUsers(context.Background())

	return svc, nil
}

// Close 关闭服务
func (s *UserService) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// SignUp 用户注册
func (s *UserService) SignUp(ctx context.Context, appID, username, password, nickname, phoneNumber string) (uint64, error) {
	if appID == "" || username == "" || password == "" {
		return 0, errors.New("app_id, username, and password are required")
	}

	userID, err := s.generateUserID(ctx, appID)
	if err != nil {
		return 0, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Failed to hash password", zap.Error(err), zap.String("app_id", appID), zap.String("username", username))
		return 0, err
	}

	user := entities.User{
		AppID:        appID,
		UserID:       userID,
		Username:     username,
		Nickname:     nickname,
		PhoneNumber:  phoneNumber,
		PasswordHash: string(hash),
		Status:       "active",
	}

	if err := s.db.WithContext(ctx).Create(&user).Error; err != nil {
		s.logger.Error("Failed to create user", zap.Error(err), zap.String("app_id", appID), zap.String("username", username))
		return 0, err
	}

	// 设置空缓存，防止穿透
	cacheKey := s.userCacheKey(appID, userID)
	if err := s.redis.Set(ctx, cacheKey, "exists", time.Hour*24).Err(); err != nil {
		s.logger.Warn("Failed to set user cache", zap.Error(err), zap.String("app_id", appID), zap.Uint64("user_id", userID))
	}

	return userID, nil
}

// SignIn 用户登录
func (s *UserService) SignIn(ctx context.Context, appID, username, password string) (uint64, error) {
	if appID == "" || username == "" || password == "" {
		return 0, errors.New("app_id, username, and password are required")
	}

	cacheKey := s.userCacheKeyByUsername(appID, username)
	result, err, _ := s.group.Do(cacheKey, func() (interface{}, error) {
		// 检查缓存（防穿透）
		if exists, err := s.redis.Get(ctx, cacheKey).Result(); err == nil && exists == "nil" {
			return uint64(0), fmt.Errorf("user %s not found in app %s", username, appID)
		}

		var user entities.User
		if err := s.db.WithContext(ctx).Where("app_id = ? AND username = ?", appID, username).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 缓存空结果
				if err := s.redis.Set(ctx, cacheKey, "nil", time.Minute*5).Err(); err != nil {
					s.logger.Warn("Failed to set nil cache", zap.Error(err), zap.String("app_id", appID), zap.String("username", username))
				}
				return uint64(0), fmt.Errorf("user %s not found in app %s", username, appID)
			}
			s.logger.Error("Failed to query user", zap.Error(err), zap.String("app_id", appID), zap.String("username", username))
			return uint64(0), err
		}

		if user.Status != "active" {
			return uint64(0), fmt.Errorf("user %s in app %s is not active (status: %s)", username, appID, user.Status)
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
			return uint64(0), errors.New("invalid password")
		}

		return user.UserID, nil
	})

	if err != nil {
		return 0, err
	}
	return result.(uint64), nil
}

// MarkUserInactive 标记用户为 inactive（逻辑删除）
func (s *UserService) MarkUserInactive(ctx context.Context, appID string, userID uint64) error {
	if err := s.CheckUserExistsAndActive(ctx, appID, userID); err != nil {
		return err
	}

	result := s.db.WithContext(ctx).Model(&entities.User{}).Where("app_id = ? AND user_id = ?", appID, userID).Update("status", "inactive")
	if result.Error != nil {
		s.logger.Error("Failed to mark user inactive", zap.Error(result.Error), zap.String("app_id", appID), zap.Uint64("user_id", userID))
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	// 失效用户缓存
	cacheKey := s.userCacheKey(appID, userID)
	if err := s.redis.Del(ctx, cacheKey).Err(); err != nil {
		s.logger.Warn("Failed to invalidate user cache", zap.Error(err), zap.String("app_id", appID), zap.Uint64("user_id", userID))
	}

	return nil
}

// CheckUserExistsAndActive 验证用户是否存在且状态为 active
func (s *UserService) CheckUserExistsAndActive(ctx context.Context, appID string, userIDs ...uint64) error {
	for _, userID := range userIDs {
		cacheKey := s.userCacheKey(appID, userID)
		result, err, _ := s.group.Do(cacheKey, func() (interface{}, error) {
			// 检查缓存（防穿透）
			if exists, err := s.redis.Get(ctx, cacheKey).Result(); err == nil && exists == "nil" {
				return nil, fmt.Errorf("user %d not found in app %s", userID, appID)
			}

			var user entities.User
			if err := s.db.WithContext(ctx).Where("app_id = ? AND user_id = ?", appID, userID).First(&user).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// 缓存空结果
					if err := s.redis.Set(ctx, cacheKey, "nil", time.Minute*5).Err(); err != nil {
						s.logger.Warn("Failed to set nil cache", zap.Error(err), zap.String("app_id", appID), zap.Uint64("user_id", userID))
					}
					return nil, fmt.Errorf("user %d not found in app %s", userID, appID)
				}
				s.logger.Error("Failed to check user existence", zap.Error(err), zap.String("app_id", appID), zap.Uint64("user_id", userID))
				return nil, err
			}

			if user.Status != "active" {
				return nil, fmt.Errorf("user %d in app %s is not active (status: %s)", userID, appID, user.Status)
			}

			// 缓存存在结果
			if err := s.redis.Set(ctx, cacheKey, "exists", time.Hour*24).Err(); err != nil {
				s.logger.Warn("Failed to set user cache", zap.Error(err), zap.String("app_id", appID), zap.Uint64("user_id", userID))
			}

			return user, nil
		})

		if err != nil {
			return err
		}
		if result == nil {
			return fmt.Errorf("user %d not found or not active in app %s", userID, appID)
		}
	}
	return nil
}

// preloadHotUsers 预加载热点用户到缓存
func (s *UserService) preloadHotUsers(ctx context.Context) {
	// 示例：加载最近活跃用户（实际可根据业务定义热点用户）
	var users []entities.User
	if err := s.db.WithContext(ctx).Where("status = ? AND updated_at > ?", "active", time.Now().Add(-24*time.Hour)).Limit(1000).Find(&users).Error; err != nil {
		s.logger.Warn("Failed to preload hot users", zap.Error(err))
		return
	}

	for _, user := range users {
		cacheKey := s.userCacheKey(user.AppID, user.UserID)
		if err := s.redis.Set(ctx, cacheKey, "exists", time.Hour*24).Err(); err != nil {
			s.logger.Warn("Failed to set user cache", zap.Error(err), zap.String("app_id", user.AppID), zap.Uint64("user_id", user.UserID))
		}
	}
}

// generateUserID 生成用户 ID（简单实现，实际可用雪花算法）
func (s *UserService) generateUserID(ctx context.Context, appID string) (uint64, error) {
	var maxID uint64
	if err := s.db.WithContext(ctx).Model(&entities.User{}).Where("app_id = ?", appID).Select("COALESCE(MAX(user_id), 0)").Scan(&maxID).Error; err != nil {
		return 0, err
	}
	return maxID + 1, nil
}

// userCacheKey 用户缓存键
func (s *UserService) userCacheKey(appID string, userID uint64) string {
	return fmt.Sprintf("user:%s:%d", appID, userID)
}

// userCacheKeyByUsername 用户名缓存键
func (s *UserService) userCacheKeyByUsername(appID, username string) string {
	return fmt.Sprintf("user:%s:username:%s", appID, username)
}
