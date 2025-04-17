package cmd_service

import (
	"context"
	"github.com/magicnana999/im/api/kitex_gen/api"
	"github.com/magicnana999/im/api/kitex_gen/api/businessservice"
	"github.com/magicnana999/im/broker"
	"github.com/magicnana999/im/broker/holder"
	"github.com/magicnana999/im/errors"
	"go.uber.org/fx"
)

type BusinessService struct {
	userHolder  *holder.UserHolder
	businessCli businessservice.Client
}

func NewBusinessService(uh *holder.UserHolder, bc businessservice.Client, lf fx.Lifecycle) *BusinessService {
	return &BusinessService{userHolder: uh, businessCli: bc}
}

func (s *BusinessService) Login(ctx context.Context, request *api.LoginRequest) (*api.LoginReply, error) {
	uc, err := broker.CurUserFromCtx(ctx)
	if err != nil {
		return nil, errors.CurUserNotFound.SetDetail(err)
	}

	rep, err := s.businessCli.Login(ctx, request)

	if err != nil {
		return nil, errors.CmdError.SetDetail(err)
	}

	if rep == nil {
		return nil, errors.CmdResponseNull
	}

	if err = s.userHolder.StoreUser(ctx, uc, rep.AppId, rep.UserId, request.Os); err != nil {
		return nil, errors.CmdError.SetDetail(err)
	}

	return rep, nil
}

func (s *BusinessService) Logout(ctx context.Context, request *api.LogoutRequest) (*api.LogoutReply, error) {
	return nil, nil
}
