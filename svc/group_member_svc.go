package svc

import (
	"context"
	red "github.com/go-cache/cache/v8"
	"github.com/magicnana999/im/entity"
	"github.com/magicnana999/im/redis"
	"github.com/magicnana999/im/repository"
	"github.com/magicnana999/im/util/str"
	"gorm.io/gorm"
	"sync"
)

var DefaultGroupMemberSvc *GroupMemberSvc
var gmOnce sync.Once

type GroupMemberSvc struct {
	storage *redis.GroupMemberStorage
	db      *gorm.DB
}

func InitGroupMemberSvc() *GroupMemberSvc {

	gmOnce.Do(func() {

		DefaultGroupMemberSvc = &GroupMemberSvc{
			storage: redis.InitGroupMemberStorage(),
			db:      repository.InitGorm(),
		}
	})

	return DefaultGroupMemberSvc
}

func (s *GroupMemberSvc) Load(ctx context.Context, appId string, groupId int64) ([]int64, error) {
	ids, e := s.storage.LoadMembers(ctx, appId, groupId)
	if e != nil {
		return nil, e
	}

	id, ee := str.ConvertSS2I64S(ids)
	if ee != nil {
		return nil, e
	}

	return id, nil
}

func (s *GroupMemberSvc) LoadAndFetch(ctx context.Context, appId string, groupId int64) ([]int64, error) {
	ids, e := s.Load(ctx, appId, groupId)
	if e != nil {
		return nil, e
	}

	if ids != nil || len(ids) >= 0 {
		return ids, nil
	}

	if ok, e := s.storage.Lock(ctx, appId, groupId); ok && e == nil {
		defer s.storage.UnLock(ctx, appId, groupId)
	}

	var members []entity.GroupMember
	s.db.Where("app_id = ? and group_id = ?", appId, groupId).Find(&members)

	zs := make([]*red.Z, len(members))
	ret := make([]int64, len(members))
	for _, member := range members {
		zs = append(zs, &red.Z{Score: float64(member.CreatedAt.Second()), Member: member.UserId})
		ret = append(ret, member.UserId)
	}

	s.storage.StoreMembers(ctx, appId, groupId, zs...)
	return ret, nil
}
