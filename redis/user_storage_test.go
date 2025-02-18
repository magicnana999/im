package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/magicnana999/im/conf"
	"github.com/magicnana999/im/constants"
	"github.com/magicnana999/im/domain"
	"github.com/magicnana999/im/logger"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	conf.LoadConfig("/Users/jinsong/source/github/im/conf/im-router.yaml")
	logger.InitLogger("debug")

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestUserStorage_StoreUserClients(t *testing.T) {
	ss := InitUserStorage()

	uc := &domain.UserConnection{
		AppId:  constants.AppId,
		UserId: 1001,
		OS:     constants.Ios,
	}

	key := KeyUserClients(uc.AppId, uc.UserId)
	label := uc.Label()
	js, _ := json.Marshal(uc)
	fmt.Println(key, label, string(js))

	ret, err := ss.StoreUserClients(context.Background(), uc)
	fmt.Println(ret, err)
}
