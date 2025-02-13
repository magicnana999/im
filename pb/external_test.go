package pb

import (
	"encoding/json"
	"errors"
	"fmt"
	imerror "github.com/magicnana999/im/errors"
	"testing"
)

func TestHeartbeat(t *testing.T) {
	p := NewHeartbeat(100)
	fmt.Println(&p, p)
	fmt.Println(&p, p.GetHeartbeatBody().Value)
}

func TestCommand(t *testing.T) {
	request := &LoginRequest{
		AppId:   "xxx",
		UserSig: "hello",
	}
	p := NewCommand(request)
	fmt.Println(&p, p)
	fmt.Println(&p, p.GetCommandBody().GetId(), p.GetCommandBody().GetLoginRequest().UserSig)

	e1 := imerror.UserSigNotFound
	e2 := errors.New("hello world error")
	reply := &LoginReply{
		AppId:  "appId",
		UserId: 100,
	}

	p1 := p.GetCommandBody().Response(reply, nil).Wrap()
	fmt.Println(&p1, p1)
	fmt.Println(&p1, p1.GetCommandBody().Id, p1.GetCommandBody().Code, p1.GetCommandBody().GetLoginReply().UserId)

	p2 := p.GetCommandBody().Response(nil, e1).Wrap()
	fmt.Println(&p2, p2)
	fmt.Println(&p2, p2.GetCommandBody().Id, p2.GetCommandBody().Code)

	p3 := p.GetCommandBody().Response(nil, e2).Wrap()
	fmt.Println(&p3, p3)
	fmt.Println(&p3, p3.GetCommandBody().Id, p3.GetCommandBody().Code)
}

func TestMessage(t *testing.T) {

	text := &TextContent{
		Text: "hello world content",
	}

	mb := NewMessage(1000, 2000, 0, 11111, "appId", "cId", text)
	js, _ := json.Marshal(mb)
	fmt.Println(mb)
	fmt.Println(string(js))

	var mb2 MessageBody
	json.Unmarshal(js, &mb2)
	fmt.Println(mb2)
	fmt.Println(mb2.Id, mb2.CId, mb2.GetTextContent().Text)

	p := mb.Wrap()
	fmt.Println(&p, p)
	fmt.Println(&p, p.GetMessageBody().Id, p.GetMessageBody().GetTextContent().Text)

	e1 := imerror.UserSigNotFound
	e2 := errors.New("hello world error")

	p2 := p.GetMessageBody().Failure(e1).Wrap()
	fmt.Println(&p2, p2)
	fmt.Println(&p2, p2.GetMessageBody().Id, p2.GetMessageBody().Code)

	p3 := p.GetMessageBody().Failure(e2).Wrap()
	fmt.Println(&p3, p3)
	fmt.Println(&p3, p3.GetMessageBody().Id, p3.GetMessageBody().Code)
}
