package redis

import (
	"context"
	red "github.com/go-redis/redis/v8"
	"github.com/magicnana999/im/logger"
	"time"
)

var DefaultGroupMemberStorage = &GroupMemberStorage{}

type GroupMemberStorage struct {
}

func InitGroupMemberStorage() *GroupMemberStorage {
	initRedis()
	return DefaultGroupMemberStorage
}

func (s *GroupMemberStorage) LoadMembers(ctx context.Context, appId string, groupId int64) ([]string, error) {
	key := KeyGroupMembers(appId, groupId)
	cmd := rds.ZRange(ctx, key, 0, -1)
	return cmd.Val(), cmd.Err()
}

func (s *GroupMemberStorage) StoreMembers(ctx context.Context, appId string, groupId int64, m ...*red.Z) (int64, error) {
	key := KeyGroupMembers(appId, groupId)
	ret := rds.ZAdd(ctx, key, m...)
	return ret.Val(), ret.Err()
}

func (s *GroupMemberStorage) Lock(ctx context.Context, appId string, groupId int64) bool {
	key := KeyGroupMembersLock(appId, groupId)
	now := time.Now().Second()
	ret := rds.Set(ctx, key, now, time.Hour)
	if ret.Err() != nil {
		logger.Error(ret.Err().Error())
		return false
	}
	return true
}

func (s *GroupMemberStorage) UnLock(ctx context.Context, appId string, groupId int64) (int64, error) {
	key := KeyGroupMembersLock(appId, groupId)
	ret := rds.Del(ctx, key)
	return ret.Val(), ret.Err()
}
