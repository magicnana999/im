package broker

import (
	"context"
	"github.com/magicnana999/im/conf"
	"github.com/magicnana999/im/constants"
	"github.com/magicnana999/im/pb"
	"github.com/magicnana999/im/util/id"
	"strings"
	"testing"
)

func Test_commandHandler_login(t *testing.T) {

	conf.LoadConfig("/Users/jinsong/source/github/im/conf/im-broker.yaml")
	InitLogger()

	s := initCommandHandler()
	req := &pb.LoginRequest{
		AppId:    constants.AppId,
		UserSig:  "cuv0ele1a37rf7ccvan0",
		Version:  "1.0.0",
		Os:       constants.Ios,
		DeviceId: strings.ToLower(id.GenerateXId()),
	}
	s.login(context.Background(), req)
}
