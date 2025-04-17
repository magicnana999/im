package handler

import (
	"context"
	"github.com/magicnana999/im/api/kitex_gen/api"
	"github.com/magicnana999/im/broker"
	"github.com/magicnana999/im/broker/domain"
	"github.com/magicnana999/im/broker/holder"
	"github.com/magicnana999/im/errors"
	"go.uber.org/fx"
	"time"
)

type HeartbeatHandler struct {
	userHolder *holder.UserHolder
	htServer   *broker.HeartbeatServer
}

func NewHeartbeatHandler(uh *holder.UserHolder, htServer *broker.HeartbeatServer, lc fx.Lifecycle) (*HeartbeatHandler, error) {
	return &HeartbeatHandler{
		userHolder: uh,
		htServer:   htServer,
	}, nil
}

func (h *HeartbeatHandler) HandlePacket(ctx context.Context, uc *domain.UserConn, packet *api.Packet) *api.Packet {
	uc, err := broker.CurUserFromCtx(ctx)
	if err != nil {
		return nil, errors.HeartbeatError.SetDetail(err)
	}

	uc.LastHeartbeat.Store(time.Now().Unix())
	return api.NewHeartbeat(int32(1)).Wrap(), nil
}

func (h *HeartbeatHandler) StartHeartbeat(ctx context.Context, uc *domain.UserConn) error {
	return h.htServer.Submit(ctx, uc)
}

func (h *HeartbeatHandler) StopHeartbeat(ctx context.Context, uc *domain.UserConn) error {
	uc.IsClosed.Store(true)
	return nil
}
