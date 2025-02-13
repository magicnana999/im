package broker

import (
	"context"
	"github.com/magicnana999/im/conf"
	"github.com/magicnana999/im/enum"
	"github.com/magicnana999/im/errors"
	"github.com/magicnana999/im/logger"
	"github.com/magicnana999/im/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"sync"
)

var defaultCommandHandler = &commandHandler{}
var commandHandlerMu sync.RWMutex

type commandHandler struct {
	conn          *grpc.ClientConn
	userApiClient pb.UserApiClient
	userState     *userState
}

func initCommandHandler() *commandHandler {

	commandHandlerMu.Lock()
	defer commandHandlerMu.Unlock()

	if defaultCommandHandler.conn != nil {
		return defaultCommandHandler
	}

	//localhost:7540
	userApiHost := conf.Global.Grpc.UserApiHost

	conn, err := grpc.NewClient(userApiHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.FatalF("init command handler user api provider error: %v", err)

	}
	defaultCommandHandler.conn = conn
	defaultCommandHandler.userApiClient = pb.NewUserApiClient(conn)
	defaultCommandHandler.userState = initUserState()

	logger.DebugF("commandHandler init")

	return defaultCommandHandler
}

func (c *commandHandler) handlePacket(ctx context.Context, packet *pb.Packet) (*pb.Packet, error) {
	reply, err := c.handleCommand(ctx, packet.GetCommandBody())

	return packet.GetCommandBody().Response(reply, err).Wrap(), nil

}

func (c *commandHandler) isSupport(ctx context.Context, packetType int32) bool {
	return pb.TypeCommand == packetType
}

func (c *commandHandler) handleCommand(ctx context.Context, cmd *pb.CommandBody) (proto.Message, error) {

	uc, e := currentUserFromCtx(ctx)
	if e != nil {
		return nil, errors.CommandHandleError.Detail(e)
	}

	switch cmd.CType {
	case pb.CTypeUserLogin:
		src := cmd.GetLoginRequest()

		rep, err := login(ctx, src)
		if err != nil {
			return nil, err
		}

		if err = c.userState.storeUser(ctx, uc, rep.AppId, rep.UserId, enum.OSType(src.Os)); err != nil {
			return nil, err
		}

		return rep, nil

	default:
		return nil, errors.CmdUnknownTypeError.DetailString("unknown type:" + cmd.CType)
	}
}

func login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginReply, error) {

	switch request.UserSig {
	case "cukpovu1a37hpofg6sj0":
		return &pb.LoginReply{
			AppId:  "19860220",
			UserId: 1200120,
		}, nil
	case "cuf5ofe1a37nfi3p4b6g":
		return &pb.LoginReply{
			AppId:  "19860220",
			UserId: 1200122,
		}, nil
	default:
		return &pb.LoginReply{
			AppId:  "19860220",
			UserId: 0,
		}, nil
	}
}
