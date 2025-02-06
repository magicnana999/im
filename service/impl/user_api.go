package impl

import (
	"context"
	"encoding/json"
	"github.com/magicnana999/im/common/entity"
	"github.com/magicnana999/im/common/pb"
	"github.com/magicnana999/im/logger"
	"github.com/magicnana999/im/service/storage"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserAPIImpl struct {
	pb.UnimplementedUserApiServer
}

func (u *UserAPIImpl) Login(ctx context.Context, p *pb.LoginRequest) (*pb.LoginReply, error) {
	logger.InfoF("Command %s,request:%s", pb.CTypeUserLogin, p)

	js, e := storage.GetUserByUserSig(ctx, p.AppId, p.UserSig)
	if e != nil {
		return nil, pb.UserSigNotFound.Format(p.UserSig)
	}

	var user entity.User
	if e := json.Unmarshal([]byte(js), &user); e != nil {
		return nil, pb.UnmarshalError
	}

	response := &pb.LoginReply{
		AppId:  user.AppId,
		UserId: user.UserId,
	}

	logger.InfoF("Command %s,response:%s", pb.CTypeUserLogin, response)
	return response, nil

}

func (u *UserAPIImpl) Logout(ctx context.Context, p *pb.LogoutRequest) (*emptypb.Empty, error) {
	return nil, nil
}

//func (u UserAPIImpl) Login(ctx context.Context, p *pb.LoginContent) (*pb.LoginReply, error) {
//	jsonstr, e := storage.GetUserByUserSig(ctx, p.AppId, p.UserSign)
//	if e != nil {
//		return nil, e
//	}
//
//	var user entity.User
//	if e := json.Unmarshal([]byte(jsonstr), &user); e != nil {
//		return nil, e
//	}
//
//	return &pb.LoginReply{
//		AppId:  user.AppId,
//		UserId: user.UserId,
//	}, nil
//}
//
//func (u UserAPIImpl) Logout(ctx context.Context, p *protocol.LogoutContent) error {
//	//TODO implement me
//	panic("implement me")
//}
