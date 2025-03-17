package broker

import (
	"context"
	"github.com/magicnana999/im/api/kitex_gen/api"
	"github.com/magicnana999/im/domain/broker/svc"
	"github.com/magicnana999/im/errors"
	"github.com/magicnana999/im/pb"
	"google.golang.org/protobuf/proto"
	"sync"
)

var defaultCommandHandler *CommandHandler
var cmdHandlerOnce sync.Once

type CommandHandler struct {
	userSvc *svc.UserSvc
}

func initCommandHandler() *CommandHandler {
	cmdHandlerOnce.Do(func() {
		defaultCommandHandler = &CommandHandler{
			userSvc: svc.InitUserSvc(),
		}
	})

	return defaultCommandHandler
}

func (c *CommandHandler) handlePacket(ctx context.Context, packet *api.Packet) (*api.Packet, error) {

	var reply proto.Message
	var err error

	mb := packet.GetCommand()

	switch mb.CType {
	case pb.CTypeUserLogin:
		reply, err = c.userSvc.Login(ctx, mb.GetLoginRequest())
	case pb.CTypeUserLogout:
		//reply, err = c.logout(ctx, mb.GetLogoutRequest())
	case pb.CTypeFriendAdd:
	case pb.CTypeFriendAddAgree:
	case pb.CTypeFriendReject:
	default:
		err = errors.CmdUnknownType
	}

	return packet.GetCommand().Response(reply, err).Wrap(), nil

}

func (c *CommandHandler) isSupport(ctx context.Context, packetType int32) bool {

	return api.TypeCommand == packetType
}
