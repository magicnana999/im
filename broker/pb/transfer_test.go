package pb

import (
	"encoding/json"
	"fmt"
	"github.com/magicnana999/im/broker/enum"
	"github.com/magicnana999/im/broker/protocol"
	"github.com/magicnana999/im/util"
	"testing"
)

func TestMessage(t *testing.T) {
	packet := &protocol.Packet{
		Id:    util.GenerateXId(),
		AppId: "appId",
		Type:  protocol.TypeMessage,
		Body: protocol.MessageBody{
			MType: protocol.MText,
			Content: protocol.TextContent{
				Text: "hello world",
			},
			At: []*protocol.At{{
				Name: "张三",
			}},
			Refer: []*protocol.Refer{{
				UserId: "231",
				MType:  protocol.MText,
				Content: protocol.TextContent{
					Text: "你好世界",
				},
			}},
		},
	}
	if b, e := json.Marshal(packet); e == nil {
		fmt.Println(string(b))
	}

	if p, err := ConvertPacket(packet); err == nil {
		var mb MessageBody
		p.Body.UnmarshalTo(&mb)

		var tt TextContent
		mb.Content.UnmarshalTo(&tt)

		fmt.Println(p)
		fmt.Println(mb)
		fmt.Println(tt)

		if pp, e := RevertPacket(p); e == nil {
			fmt.Println(pp)
		} else {
			fmt.Println(e)
		}
	}

}

func TestCommand(t *testing.T) {
	packet := &protocol.Packet{
		Id:    util.GenerateXId(),
		AppId: "appId",
		Type:  protocol.TypeCommand,
		Body: protocol.CommandBody{
			MType: protocol.MLogin,
			Token: util.GenerateXId(),
			Content: protocol.LoginContent{
				Version:      "hello world",
				OS:           enum.MacOS,
				DeviceId:     util.GenerateXId(),
				PushDeviceId: util.GenerateXId(),
			},
		},
	}
	if b, e := json.Marshal(packet); e == nil {
		fmt.Println(string(b))
	}

	if p, err := ConvertPacket(packet); err == nil {
		var mb CommandBody
		p.Body.UnmarshalTo(&mb)

		var tt LoginContent
		mb.Content.UnmarshalTo(&tt)

		fmt.Println(p)
		fmt.Println(mb)
		fmt.Println(tt)

		if pp, e := RevertPacket(p); e == nil {
			fmt.Println(pp)
		} else {
			fmt.Println(e)
		}
	}

}
