package svc

import (
	"context"
	"github.com/magicnana999/im/errors"
	"github.com/magicnana999/im/pb"
	"github.com/magicnana999/im/redis"
	"github.com/magicnana999/im/repository"
	"gorm.io/gorm"
	"sync"
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

func (s *UserSvc) Login(ctx context.Context, p *pb.LoginRequest) (*pb.LoginReply, error) {

	user, e := s.storage.LoadUserSig(ctx, p.AppId, p.UserSig)
	if e != nil {
		return nil, errors.UserSigNotFound.Detail(e)
	}

	if user == nil {
		return nil, &errors.UserSigNotFound
	}

	response := &pb.LoginReply{
		AppId:  user.AppId,
		UserId: user.UserId,
	}

	return response, nil

}
