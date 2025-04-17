package handler

import (
	"context"
	"github.com/magicnana999/im/api/kitex_gen/api"
	"github.com/magicnana999/im/broker/cmd_service"
	"github.com/magicnana999/im/broker/holder"
	"github.com/magicnana999/im/errors"
	"github.com/magicnana999/im/pb"
	"google.golang.org/protobuf/proto"
)

type CommandHandler struct {
	userHolder  *holder.UserHolder
	userService *cmd_service.BusinessService
}

func NewCommandHandler(uh *holder.UserHolder) *CommandHandler {

	return commandHandlerSingleton.Get(func() *CommandHandler {
		return &CommandHandler{
			userHolder:  holder.NewUserHolder(),
			userService: cmd_service.NewUserService(),
		}
	})
}

func (c *CommandHandler) handlePacket(ctx context.Context, packet *api.Packet) (*api.Packet, error) {

	var reply proto.Message
	var err error

	mb := packet.GetCommand()

	switch mb.CommandType {
	case pb.CommandTypeUserLogin:
		reply, err = c.userService.Login(ctx, mb.GetLoginRequest())
	case pb.CommandTypeUserLogout:
		reply, err = c.userService.Logout(ctx, mb.GetLogoutRequest())
	default:
		err = errors.CmdUnknownType
	}

	return packet.GetCommand().Response(reply, err).Wrap(), nil

}
