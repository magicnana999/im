package pb

import (
	"errors"
	"fmt"
	"github.com/magicnana999/im/util/id"
	"strings"
	"testing"
)

func TestNewCommandRequest(t *testing.T) {
	c := LoginContent{
		UserSig:      "xx",
		AppId:        "19860229",
		Version:      "1.0.0",
		Os:           OSType_OSIos,
		PushDeviceId: strings.ToUpper(id.GenerateXId()),
	}
	r, _ := NewCommandRequest(0, CTypeUserLogin, &c)
	fmt.Println(r)

	rep := LoginReply{
		AppId:  "19860220",
		UserId: 1200120,
	}
	ret, _ := NewCommandResponse(r, CTypeUserLogin, &rep, nil)
	fmt.Println(ret)

	ret1, _ := NewCommandResponse(r, CTypeUserLogin, &rep, UserSigNotFound)
	fmt.Println(ret1)

	ret2, _ := NewCommandResponse(r, CTypeUserLogin, &rep, errors.New("test"))
	fmt.Println(ret2)
}
