package service

import (
	"context"
	"github.com/magicnana999/im/entity"
	"github.com/magicnana999/im/errors"
	"github.com/magicnana999/im/pb"
	"github.com/magicnana999/im/redis"
	"github.com/magicnana999/im/repository"
	"gorm.io/gorm"
	"sync"
	"time"
)

var DefaultUserSvc *UserSvc
var userOnce sync.Once

type UserSvc struct {
	pb.UnimplementedUserApiServer
	db      *gorm.DB
	storage *redis.UserStorage
}

func InitUserSvc() *UserSvc {

	userOnce.Do(func() {
		DefaultUserSvc = &UserSvc{
			db:      repository.InitGorm(),
			storage: redis.InitUserStorage(),
		}
	})

	return DefaultUserSvc
}

func (s *UserSvc) LoadByUserId(ctx context.Context, appId string, userId int64) (*entity.User, error) {

}

func (s *UserSvc) Login(ctx context.Context, p *pb.LoginRequest) (*pb.LoginReply, error) {

	user, e := s.storage.LoadUserSig(ctx, p.AppId, p.UserSig)
	if e != nil {
		return nil, errors.UserSigNotFound.Detail(e)
	}

	if user == nil {
		return nil, &errors.UserSigNotFound
	}

	s.db.Begin()
	u := &entity.User{
		Os:        p.Os,
		IsLogin:   int(pb.YES),
		LastLogin: time.Now(),
	}
	err := s.db.Where("app_id = ? and user_id = ?", user.AppId, user.UserId).UpdateColumns(u).Error

	if err != nil {
		s.db.Rollback()
		return nil, errors.CmdError.Detail(err)
	}

	response := &pb.LoginReply{
		AppId:  user.AppId,
		UserId: user.UserId,
	}

	return response, nil

}

func (s *UserSvc) Logout(ctx context.Context, p *pb.LogoutRequest) error {

	s.db.Begin()
	u := &entity.User{
		IsLogin: int(pb.NO),
	}
	err := s.db.Where("app_id = ? and user_id = ?", p.AppId, p.UserId).UpdateColumns(u).Error

	if err != nil {
		s.db.Rollback()
		return errors.CmdError.Detail(err)
	}

	return nil

}
