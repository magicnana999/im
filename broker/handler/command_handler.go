package handler

import (
	"context"
	"github.com/magicnana999/im/broker/state"
	"github.com/magicnana999/im/enum"
	"github.com/magicnana999/im/errors"
	"github.com/magicnana999/im/logger"
	"github.com/magicnana999/im/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

var DefaultCommandHandler = &CommandHandler{}

type CommandHandler struct {
	conn          *grpc.ClientConn
	userApiClient pb.UserApiClient
}

func (c *CommandHandler) HandlePacket(ctx context.Context, packet *pb.Packet) (*pb.Packet, error) {
	reply, err := c.HandleCommand(ctx, packet.GetCommandBody())

	return pb.NewCommandResponse(reply, err)

}

func (c *CommandHandler) IsSupport(ctx context.Context, packetType int32) bool {
	return pb.TypeCommand == packetType
}

func (c *CommandHandler) InitHandler() error {
	conn, err := grpc.NewClient("localhost:7540", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.FatalF("did not connect: %v", err)
	}
	c.conn = conn
	c.userApiClient = pb.NewUserApiClient(conn)

	return nil
}

func (c *CommandHandler) HandleCommand(ctx context.Context, cmd *pb.CommandBody) (proto.Message, error) {

	uc, e := state.CurrentUserFromCtx(ctx)
	if e != nil {
		return nil, e
	}

	switch cmd.CType {
	case pb.CTypeUserLogin:
		src := cmd.GetLoginRequest()

		rep, err := login(ctx, src)
		if err != nil {
			return nil, err
		}

		if err = uc.Store(ctx, rep.AppId, rep.UserId, enum.OSType(src.Os)); err != nil {
			return nil, err
		}

		return rep, nil

	default:
		return nil, errors.HandleInvalidCType
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
