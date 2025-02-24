package apiimpl

import (
	"context"
	"github.com/magicnana999/im/errors"
	"github.com/magicnana999/im/pb"
	"github.com/magicnana999/im/svc"
	"sync"
)

var DefaultUserApi *UserAPIImpl
var userApiOnce sync.Once

type UserAPIImpl struct {
	pb.UnimplementedUserApiServer
	userSvc *svc.UserSvc
}

func InitUserApi() *UserAPIImpl {
	userApiOnce.Do(func() {
		DefaultUserApi = &UserAPIImpl{
			userSvc: svc.InitUserSvc(),
		}
	})

	return DefaultUserApi
}

func (s *UserAPIImpl) Login(ctx context.Context, in *pb.LoginRequest) (*pb.ApiResult, error) {

	var ret *pb.ApiResult

	reply, err := s.userSvc.Login(ctx, in)
	if err != nil {

		e := errors.Format(err)

		ret = &pb.ApiResult{
			Code:    int32(e.Code),
			Message: e.String(),
		}
	} else {
		ret = &pb.ApiResult{
			Code: 0,
			Data: &pb.ApiResult_LoginReply{LoginReply: reply},
		}
	}

	return ret, nil
}

func (s *UserAPIImpl) Logout(ctx context.Context, in *pb.LogoutRequest) (*pb.ApiResult, error) {
	err := s.userSvc.Logout(ctx, in)
	if err != nil {
		e := errors.Format(err)
		return &pb.ApiResult{
			Code:    int32(e.Code),
			Message: e.String(),
		}, nil
	} else {

	}

	return &pb.ApiResult{
		Code: 0,
	}, nil

}
