package impl

import (
	"context"
	"fmt"
	"github.com/magicnana999/im/common/pb"
	"testing"
)

func TestUserAPIImpl_Login(t *testing.T) {

	client := &UserAPIImpl{}
	p := pb.LoginContent{
		AppId:   "19860220",
		UserSig: "cuf5ofe1a37nfi3p4b5g1",
	}
	c, e := client.Login(context.Background(), &p)
	if e != nil {
		t.Error(e)
	}

	fmt.Println(c)
}
