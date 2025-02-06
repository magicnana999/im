package handler

import (
	"context"
	"github.com/magicnana999/im/broker/state"
	"github.com/magicnana999/im/common/pb"
	"github.com/magicnana999/im/errors"
	"github.com/magicnana999/im/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

var DefaultCommandHandler = &CommandHandler{}

type CommandHandler struct {
	conn          *grpc.ClientConn
	userApiClient pb.UserApiClient
}

func (c *CommandHandler) HandlePacket(ctx context.Context, packet *pb.Packet) (*pb.Packet, error) {

	var body pb.CommandBody
	if err := packet.Body.UnmarshalTo(&body); err != nil {
		return nil, errors.HandleUnmarshalError.Fill(err.Error())
	}

	reply, err := c.HandleCommand(ctx, body.CType, body.Request)

	return pb.NewCommandResponse(packet, body.CType, reply, err)

}

func (c *CommandHandler) IsSupport(ctx context.Context, packetType int32) bool {
	return pb.BTypeCommand == packetType
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

func (c *CommandHandler) HandleCommand(ctx context.Context, cType string, content *anypb.Any) (proto.Message, error) {

	uc, e := state.CurrentUserFromCtx(ctx)
	if e != nil {
		return nil, e
	}

	switch cType {
	case pb.CTypeUserLogin:
		var src pb.LoginRequest
		if err := content.UnmarshalTo(&src); err != nil {
			return nil, errors.HandleUnmarshalError.Fill(err.Error())
		}
		//return c.userApiClient.Login(ctx, &src)

		rep, err := login(ctx, &src)
		if err != nil {
			return nil, err
		}

		if err = uc.Store(ctx, rep.AppId, rep.UserId); err != nil {
			return nil, err
		}

		return rep, nil

	default:
		return nil, errors.HandleInvalidCType
	}
}

func login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginReply, error) {

	switch request.UserSig {
	case "cuf5ofe1a37nfi3p4b5g":
		return &pb.LoginReply{
			AppId:  "19860220",
			UserId: 1200120,
		}, nil
	case "cuf5ofe1a37nfi3p4b6g":
		return &pb.LoginReply{
			AppId:  "19860220",
			UserId: 1200122,
		}, nil
	case "cuf5ofe1a37nfi3p4b60":
		return &pb.LoginReply{
			AppId:  "19860221",
			UserId: 1200122,
		}, nil
	default:
		return &pb.LoginReply{
			AppId:  "19860220",
			UserId: 0,
		}, nil
	}
}
