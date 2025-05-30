package cmd_service

import (
	"context"
	"github.com/magicnana999/im/api/kitex_gen/api"
	"github.com/magicnana999/im/api/kitex_gen/api/businessservice"
	"github.com/magicnana999/im/broker/holder"
	"github.com/magicnana999/im/errors"
	"github.com/magicnana999/im/pkg/id"
	"go.uber.org/fx"
)

type UserService struct {
	userHolder  *holder.UserHolder
	businessCli businessservice.Client
}

func NewUserService(uh *holder.UserHolder, bc businessservice.Client, lf fx.Lifecycle) (*UserService, error) {
	bs := &UserService{userHolder: uh, businessCli: bc}
	return bs, nil
}

// Login 登录 RPC
func (s *UserService) Login(ctx context.Context, request *api.LoginRequest) (*api.LoginReply, error) {

	//rep, err := s.businessCli.Login(ctx, request)

	rep := &api.LoginReply{
		AppId:  "19860220",
		UserId: id.SnowflakeID(),
	}

	//if err != nil {
	//	return nil, errors.LoginErr.SetDetail(err.Error())
	//}

	if rep == nil {
		return nil, errors.LoginErr.SetDetail("reply is nil")
	}

	return rep, nil
}

// Logout 登出
func (s *UserService) Logout(ctx context.Context, request *api.LogoutRequest) (*api.LogoutReply, error) {
	return nil, nil
}
