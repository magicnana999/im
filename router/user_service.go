package router

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/magicnana999/im/infra"
	"github.com/magicnana999/im/vo"
	"go.uber.org/fx"
)

type UserService struct {
	rds *redis.Client
}

func NewUserService(rds *redis.Client, lc fx.Lifecycle) *UserService {
	return &UserService{
		rds: rds,
	}
}

func (s *UserService) GetUserClients(ctx context.Context, appId string, userId int64) ([]vo.UserClient, error) {

	key := infra.KeyUserClients(appId, userId)
	m, e := s.rds.HGetAll(ctx, key).Result()
	if e != nil {
		return nil, e
	}

	var clients []vo.UserClient
	for k, v := range m {
		var client vo.UserClient
		err := json.Unmarshal([]byte(v), &client)
		if err != nil {
			return nil, err
		}
		client.Label = k
		clients = append(clients, client)
	}

	if len(clients) == 0 {
		return nil, errors.New("no user client found")
	}

	return clients, nil
}
