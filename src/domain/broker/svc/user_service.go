package svc

import (
	"context"
	"github.com/magicnana999/im/api/kitex_gen/api"
	"github.com/magicnana999/im/api/kitex_gen/api/serverservice"
	"github.com/magicnana999/im/constants"
	"github.com/magicnana999/im/domain/broker"
	"github.com/magicnana999/im/domain/broker/state"
	"github.com/magicnana999/im/errors"
	"github.com/magicnana999/im/infra"
	"google.golang.org/protobuf/proto"
	"sync"
)

var (
	usvcOnce       sync.Once
	DefaultUserSvc *UserSvc
)

type UserSvc struct {
	userState *state.UserState
	serverCli serverservice.Client
}

func InitUserSvc() *UserSvc {
	usvcOnce.Do(func() {
		DefaultUserSvc = &UserSvc{
			userState: state.InitUserState(),
			serverCli: infra.InitServerCli(),
		}
	})

	return DefaultUserSvc
}

func (s *UserSvc) Login(ctx context.Context, request *api.LoginRequest) (proto.Message, error) {
	uc, err := broker.UserFromCtx(ctx)
	if err != nil {
		return nil, errors.CurUserNotFound.SetDetail(err)
	}

	rep, err := s.serverCli.Login(ctx, request)

	if err != nil {
		return nil, errors.CmdError.SetDetail(err)
	}

	if rep == nil {
		return nil, errors.CmdResponseNull
	}

	if rep == nil {
		return nil, errors.CmdReplyNull
	}

	if err = s.userState.StoreUser(ctx, uc, rep.AppId, rep.UserId, constants.OSType(request.Os)); err != nil {
		return nil, errors.CmdError.SetDetail(err)
	}

	return rep, nil
}
