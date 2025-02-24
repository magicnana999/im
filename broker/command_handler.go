package broker

import (
	"context"
	"github.com/magicnana999/im/conf"
	"github.com/magicnana999/im/errors"
	"github.com/magicnana999/im/logger"
	"github.com/magicnana999/im/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"sync"
)

var defaultCommandHandler *commandHandler
var cmdHandlerOnce sync.Once

type commandHandler struct {
	conn          *grpc.ClientConn
	userApiClient pb.UserApiClient
	userState     *userState
}

func initCommandHandler() *commandHandler {

	cmdHandlerOnce.Do(func() {

		defaultCommandHandler = &commandHandler{}
		userApiHost := conf.Global.Service.Addr

		conn, err := grpc.NewClient(
			userApiHost,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithUnaryInterceptor(logger.UnaryClientInterceptor()))

		if err != nil {
			logger.Fatalf("init command handler user api provider error: %v", err)

		}
		defaultCommandHandler.conn = conn
		defaultCommandHandler.userApiClient = pb.NewUserApiClient(conn)
		defaultCommandHandler.userState = initUserState()

	})

	return defaultCommandHandler
}

func (c *commandHandler) handlePacket(ctx context.Context, packet *pb.Packet) (*pb.Packet, error) {

	var reply proto.Message
	var err error

	mb := packet.GetCommandBody()

	switch mb.CType {
	case pb.CTypeUserLogin:
		reply, err = c.login(ctx, mb.GetLoginRequest())
	case pb.CTypeUserLogout:
		reply, err = c.logout(ctx, mb.GetLogoutRequest())
	case pb.CTypeFriendAdd:
	case pb.CTypeFriendAddAgree:
	case pb.CTypeFriendReject:
	default:
		err = errors.CmdUnknownType
	}

	return packet.GetCommandBody().Response(reply, err).Wrap(), nil

}

func (c *commandHandler) isSupport(ctx context.Context, packetType int32) bool {

	return pb.TypeCommand == packetType
}

func (c *commandHandler) login(ctx context.Context, req *pb.LoginRequest) (proto.Message, error) {

	uc, e := currentUserFromCtx(ctx)
	if e != nil {
		return nil, errors.CurUserNotFound.Detail(e)
	}

	rep, err := c.userApiClient.Login(ctx, req)

	if err != nil {
		return nil, errors.CmdError.Detail(err)
	}

	if rep == nil {
		return nil, errors.CmdResponseNull
	}

	if rep.Code != 0 {
		return nil, errors.New(int(rep.Code), rep.Message)
	}

	if rep.GetLoginReply() == nil {
		return nil, errors.CmdReplyNull
	}

	ret := rep.GetLoginReply()

	if err = c.userState.storeUser(ctx, uc, ret.AppId, ret.UserId, req.Os); err != nil {
		return nil, errors.CmdError.Detail(err)
	}

	return ret, nil

}

func (c *commandHandler) logout(ctx context.Context, req *pb.LogoutRequest) (proto.Message, error) {
	rep, err := c.userApiClient.Logout(ctx, req)

	if err != nil {
		return nil, errors.CmdError.Detail(err)
	}

	if rep == nil {
		return nil, errors.CmdResponseNull
	}

	if rep.Code != 0 {
		return nil, errors.New(int(rep.Code), rep.Message)
	}

	return nil, nil
}
