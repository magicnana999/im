package redis

import (
	"context"
	"time"
)

var DefaultSequenceStorage = &SequenceStorage{}

type SequenceStorage struct {
}

func InitSequenceStorage() *SequenceStorage {
	initRedis()
	return DefaultSequenceStorage
}

func (s *SequenceStorage) Increase(ctx context.Context, appId, cId string) (int64, error) {
	key := KeySequence(appId, cId)
	cmd := rds.Incr(ctx, key)
	return cmd.Val(), cmd.Err()
}

func (s *SequenceStorage) Lock(ctx context.Context, appId, cId string) (bool, error) {
	key := KeySequence(appId, cId)
	cmd := rds.SetNX(ctx, key, time.Now().UnixMilli(), time.Hour)
	return cmd.Val(), cmd.Err()
}

func (s *SequenceStorage) Unlock(ctx context.Context, appId, cId string) (int64, error) {
	key := KeySequence(appId, cId)
	cmd := rds.Del(ctx, key)
	return cmd.Val(), cmd.Err()
}
