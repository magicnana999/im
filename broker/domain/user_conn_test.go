package domain

import (
	"encoding/json"
	"fmt"
	"github.com/magicnana999/im/define"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUserConnToJSON(t *testing.T) {
	uc := &UserConn{}
	uc.Login("11", 100, define.Ios)
	uc.Refresh(time.Now().UnixMilli())
	bs, err := uc.ToJSON()
	assert.Nil(t, err)
	fmt.Println(string(bs))

	js := "{\"fd\":0,\"appId\":\"11\",\"userId\":100,\"os\":\"iOS\",\"clientAddr\":\"\",\"brokerAddr\":\"\",\"connectTime\":0,\"isLogin\":true,\"isClosed\":false,\"lastHeartbeat\":1745342982242}"

	var u UserConn
	err = json.Unmarshal([]byte(js), &u)
	assert.Nil(t, err)
	fmt.Println(u.AppId.Load(), u.UserId.Load())
}
