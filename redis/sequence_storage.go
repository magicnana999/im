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

func (s *SequenceStorage) Increase(ctx context.Context, appId, seqId string) (int64, error) {
	key := KeySequence(appId, seqId)
	cmd := rds.Incr(ctx, key)
	return cmd.Val(), cmd.Err()
}

func (s *SequenceStorage) Store(ctx context.Context, appId, seqId string, batch int64) (bool, error) {
	key := KeySequence(appId, seqId)
	cmd := rds.SetNX(ctx, key, batch, time.Minute)
	return cmd.Val(), cmd.Err()
}

func (s *SequenceStorage) Lock(ctx context.Context, appId, seqId string) bool {
	key := KeySequence(appId, seqId)
	cmd := rds.SetNX(ctx, key, time.Now().UnixMilli(), time.Hour)
	if cmd.Err() == nil {
		return cmd.Val()
	} else {
		return false
	}
}

func (s *SequenceStorage) Unlock(ctx context.Context, appId, seqId string) int64 {
	key := KeySequence(appId, seqId)
	cmd := rds.Del(ctx, key)
	if cmd.Err() == nil {
		return cmd.Val()
	} else {
		return 0
	}
}
