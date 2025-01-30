package pb

import (
	"fmt"
	"github.com/magicnana999/im/broker/protocol"
	"github.com/magicnana999/im/util"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"testing"
	"time"
)

func Test(t *testing.T) {
	text, _ := anypb.New(&TextContent{
		Text: "hello",
	})

	message, _ := anypb.New(&MessageBody{
		MType:    protocol.MText,
		CId:      "sjfjsdf",
		To:       "200",
		GroupId:  "",
		TType:    protocol.TSingle,
		Sequence: 1212,
		Content:  text,
	})

	packet := &Packet{
		Id:      util.GenerateXId(),
		AppId:   "100",
		UserId:  "10011",
		Flow:    protocol.FlowRequest,
		NeedAck: protocol.YES,
		Type:    protocol.TypeMessage,
		CTime:   time.Now().UnixMilli(),
		STime:   time.Now().UnixMilli(),
		Body:    message,
	}

	bs, _ := proto.Marshal(packet)
	fmt.Println(len(bs))

	proto.Unmarshal(bs, packet)

	var mb MessageBody
	packet.Body.UnmarshalTo(&mb)

	fmt.Println(packet.Id)
	fmt.Println(mb.Sequence)

}
