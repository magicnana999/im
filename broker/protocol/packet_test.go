package protocol

import (
	"encoding/json"
	"fmt"
	"github.com/magicnana999/im/util"
	"testing"
)

func Test(t *testing.T) {

	text := TextContent{
		Text: "Hello",
	}

	message := MessageBody{
		MType:    MText,
		CId:      "sdfsdf",
		To:       "sdfsdf",
		GroupId:  "",
		TType:    TSingle,
		Sequence: 100,
		Content:  text,
		At:       nil,
		Refer:    nil,
	}

	packet := &Packet{
		Id:      util.GenerateXId(),
		AppId:   "STARTSPACE",
		UserId:  "sdifejrjersdf",
		Flow:    FlowRequest,
		NeedAck: YES,
		Type:    TypeMessage,
		CTime:   12123123,
		STime:   12123123,
		Body:    message,
	}

	b, _ := json.Marshal(packet)
	fmt.Println(string(b))

	json.Unmarshal(b, packet)

	fmt.Println(*packet)

	if packet.Type == TypeMessage {

	}

}
