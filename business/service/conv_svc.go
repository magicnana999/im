package service

import (
	"context"
	"github.com/magicnana999/im/entity"
	"github.com/magicnana999/im/repository"
	"gorm.io/gorm"
	"sync"
)

var DefaultConvSvc *ConvSvc
var convOnce sync.Once

type ConvSvc struct {
	db *gorm.DB
}

func InitConvSvc() *ConvSvc {
	convOnce.Do(func() {

		DefaultConvSvc = &ConvSvc{
			db: repository.InitGorm(),
		}
	})

	return DefaultConvSvc
}

func (s *ConvSvc) Save(ctx context.Context, conv *entity.Conv) error {
	return s.db.Save(conv).Error
}

func (s *ConvSvc) QueryOne(ctx context.Context, appId, convId string) (*entity.Conv, error) {
	var conv entity.Conv
	s.db.Where("app_id = ? and conv_id = ?", appId, convId).First(&conv)
	return &conv, nil
}

func (s *ConvSvc) QueryLatest(ctx context.Context, appId string) ([]entity.Conv, error) {
	var convs []entity.Conv
	s.db.Where("app_id = ? and user_id = ? order by updated_at desc limit 0,20", appId).Find(&convs)
	return convs, nil
}
