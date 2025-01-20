package pb

import (
	"fmt"
	"github.com/magicnana999/im/broker/protocol"
	"github.com/magicnana999/im/util"
	"testing"
)

func TestProtoSerializer_Serialize(t *testing.T) {

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

	b, e := ProtoSerializerInstance.Serialize(packet)
	if e != nil {
		t.Error(e)
	}

	fmt.Println(len(b))

	p, ee := ProtoSerializerInstance.Deserialize(b)
	if ee != nil {
		t.Error(ee)
	}

	fmt.Println(p)

}
