package handler

import (
	"context"
	"github.com/magicnana999/im/api/kitex_gen/api"
	"github.com/magicnana999/im/broker/cmd_service"
	"github.com/magicnana999/im/broker/holder"
	"github.com/magicnana999/im/errors"
	"google.golang.org/protobuf/proto"
)

type CommandHandler struct {
	userHolder  *holder.UserHolder
	userService *cmd_service.UserService
}

func NewCommandHandler(uh *holder.UserHolder, us *cmd_service.UserService) (*CommandHandler, error) {
	return &CommandHandler{
		userHolder:  uh,
		userService: us,
	}, nil

}

func (c *CommandHandler) HandlePacket(ctx context.Context, packet *api.Packet) (*api.Packet, error) {

	var reply proto.Message
	var err error

	mb := packet.GetCommand()

	switch mb.CommandType {
	case api.CommandTypeUserLogin:
		reply, err = c.userService.Login(ctx, mb.GetLoginRequest())
	case api.CommandTypeUserLogout:
		reply, err = c.userService.Logout(ctx, mb.GetLogoutRequest())
	default:
		err = errors.CmdUnknownType
	}

	return packet.GetCommand().Response(reply, err).Wrap(), nil

}
