package handler

import (
	"context"
	"fmt"
	"github.com/magicnana999/im/common/pb"
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

func (c2 *CommandHandler) HandlePacket(ctx context.Context, packet *pb.Packet) (*pb.Packet, error) {

	var body pb.CommandBody
	if err := packet.Body.UnmarshalTo(&body); err != nil {
		return nil, err
	}

	reply, err := c2.HandleCommand(ctx, body.CType, body.Request)

	return pb.NewCommandResponse(packet, body.CType, reply, err)

}

func (c2 *CommandHandler) IsSupport(ctx context.Context, packetType int32) bool {
	return pb.BTypeCommand == packetType
}

func (c2 *CommandHandler) InitHandler() {
	conn, err := grpc.NewClient("localhost:7540", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.FatalF("did not connect: %v", err)
	}
	c2.conn = conn
	c2.userApiClient = pb.NewUserApiClient(conn)
}

func (c2 *CommandHandler) HandleCommand(ctx context.Context, cType string, content *anypb.Any) (proto.Message, error) {

	switch cType {
	case pb.CTypeUserLogin:
		var src pb.LoginRequest
		if err := content.UnmarshalTo(&src); err != nil {
			break
		}
		return c2.userApiClient.Login(ctx, &src)
	default:
		return nil, fmt.Errorf("unsupported command type: %v", cType)
	}

	return nil, fmt.Errorf("unsupported command type: %v", cType)
}
