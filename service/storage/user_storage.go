package storage

import (
	"context"
	"encoding/json"
	"github.com/magicnana999/im/common/entity"
	"github.com/magicnana999/im/redis"
	"github.com/magicnana999/im/util/id"
)

func GetUserByUserSig(ctx context.Context, appId, userSig string) (string, error) {
	cmd := redis.RDS.Get(ctx, redis.KeyUserSig(appId, userSig))
	return cmd.Val(), cmd.Err()
}

func SetUserByUserSig(ctx context.Context, appId string, user *entity.User) (string, error) {
	sig := id.GenerateXId()
	json1, _ := json.Marshal(user)

	cmd := redis.RDS.Set(ctx, redis.KeyUserSig(appId, sig), json1, -1)
	return cmd.Val(), cmd.Err()
}
