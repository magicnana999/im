package cmd_service

import (
	"context"
	"github.com/magicnana999/im/api/kitex_gen/api"
	"github.com/magicnana999/im/api/kitex_gen/api/businessservice"
	brokerctx "github.com/magicnana999/im/broker/ctx"
	"github.com/magicnana999/im/broker/holder"
	"github.com/magicnana999/im/errors"
	"github.com/magicnana999/im/pkg/id"
	"go.uber.org/fx"
	"time"
)

type UserService struct {
	userHolder  *holder.UserHolder
	businessCli businessservice.Client
}

func NewUserService(uh *holder.UserHolder, bc businessservice.Client, lf fx.Lifecycle) (*UserService, error) {
	bs := &UserService{userHolder: uh, businessCli: bc}
	return bs, nil
}

func (s *UserService) Login(ctx context.Context, request *api.LoginRequest) (*api.LoginReply, error) {

	uc, err := brokerctx.GetCurUserConn(ctx)
	if err != nil {
		return nil, errors.LoginErr.SetDetail(err.Error())
	}

	//rep, err := s.businessCli.Login(ctx, request)

	rep := &api.LoginReply{
		AppId:  "19860220",
		UserId: id.SnowflakeID(),
	}

	if err != nil {
		return nil, errors.LoginErr.SetDetail(err.Error())
	}

	if rep == nil {
		return nil, errors.LoginErr.SetDetail("reply is nil")
	}

	uc.AppId.Store(rep.GetAppId())
	uc.UserId.Store(rep.GetUserId())
	uc.OS.Store(request.Os)
	uc.IsLogin.Store(true)
	uc.Refresh(time.Now())

	//TODO ... 这里要做互踢
	s.userHolder.HoldUserConn(uc)

	return rep, nil
}

func (s *UserService) Logout(ctx context.Context, request *api.LogoutRequest) (*api.LogoutReply, error) {
	return nil, nil
}
