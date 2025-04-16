package handler

import (
	"context"
	"github.com/magicnana999/im/api/kitex_gen/api"
	"github.com/magicnana999/im/broker"
	"github.com/magicnana999/im/broker/domain"
	"github.com/magicnana999/im/broker/holder"
	"github.com/magicnana999/im/errors"
	"github.com/magicnana999/im/pb"
	"github.com/magicnana999/im/pkg/singleton"
	"github.com/panjf2000/gnet/v2"
	"sync/atomic"
	"time"
)

var heartbeatHandlerSingleton = singleton.NewSingleton[*HeartbeatHandler]()

type HeartbeatHandler struct {
	userHolder *holder.UserHolder
	htServer   *broker.HeartbeatServer
}

func NewHeartbeatHandler(htServer *broker.HeartbeatServer) *HeartbeatHandler {
	return heartbeatHandlerSingleton.Get(func() *HeartbeatHandler {
		return &HeartbeatHandler{
			userHolder: holder.NewUserHolder(),
			htServer:   htServer,
		}
	})
}

func (h *HeartbeatHandler) handlePacket(ctx context.Context, packet *api.Packet) (*api.Packet, error) {
	uc, err := broker.CurUserFromCtx(ctx)
	if err != nil {
		return nil, errors.HeartbeatError.SetDetail(err)
	}

	atomic.StoreInt64(uc.LastHeartbeat, time.Now().Unix())
	return pb.NewHeartbeat(int32(1)), nil
}

func (h *HeartbeatHandler) StartHeartbeat(ctx context.Context, c gnet.Conn, uc *domain.UserConn) error {
	return h.htServer.Submit(ctx, c, uc)
}
