package pb

import (
	"encoding/json"
	"fmt"
	"github.com/magicnana999/im/errors"
	"github.com/magicnana999/im/util/id"
	"strings"
	"testing"
	"time"
)

func TestNewCommandRequest(t *testing.T) {
	packet, err := NewCommandRequest(&LoginRequest{
		AppId:   "ssss",
		UserSig: "3234234",
	})

	if err != nil {
		t.Errorf("NewCommandRequest() error = %v", err)
	}

	fmt.Println(packet.IsCommand())
	fmt.Println(packet.GetCommandBody())

}

func TestNewCommandResponse(t *testing.T) {
	packet, err := NewCommandResponse(&LoginReply{
		AppId:  "ssss",
		UserId: 1212,
	}, errors.HandlerNoSupportError)

	if err != nil {
		t.Errorf("NewCommandResponse() error = %v", err)
	}

	fmt.Println(packet.IsCommand())
	fmt.Println(packet.GetCommandBody())
}

func TestNewHeartbeat(t *testing.T) {
	packet := NewHeartbeat(12)
	fmt.Println(packet.IsHeartbeat())
	fmt.Println(packet.GetHeartbeatBody())
}

/*
*
Id       string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`

	AppId    string                 `protobuf:"bytes,2,opt,name=appId,proto3" json:"appId,omitempty"`
	UserId   int64                  `protobuf:"varint,3,opt,name=userId,proto3" json:"userId,omitempty"`
	CId      string                 `protobuf:"bytes,4,opt,name=cId,proto3" json:"cId,omitempty"`
	To       string                 `protobuf:"bytes,5,opt,name=to,proto3" json:"to,omitempty"`
	GroupId  string                 `protobuf:"bytes,6,opt,name=groupId,proto3" json:"groupId,omitempty"`
	Sequence int64                  `protobuf:"varint,7,opt,name=sequence,proto3" json:"sequence,omitempty"`
	Flow     int32                  `protobuf:"varint,8,opt,name=flow,proto3" json:"flow,omitempty"`
	NeedAck  int32                  `protobuf:"varint,9,opt,name=needAck,proto3" json:"needAck,omitempty"`
	CTime    int64                  `protobuf:"varint,10,opt,name=cTime,proto3" json:"cTime,omitempty"`
	STime    int64                  `protobuf:"varint,11,opt,name=sTime,proto3" json:"sTime,omitempty"`
	CType    string                 `protobuf:"bytes,12,opt,name=cType,proto3" json:"cType,omitempty"`
	At       []*At                  `protobuf:"bytes,13,rep,name=at,proto3" json:"at,omitempty"`
	Refer    []*Refer               `protobuf:"bytes,14,rep,name=refer,proto3" json:"refer,omitempty"`
	Code     int32                  `protobuf:"varint,15,opt,name=code,proto3" json:"code,omitempty"`
	Message  string                 `protobuf:"bytes,16,opt,name=message,proto3" json:"message,omitempty"`
	// Types that are valid to be assigned to Content:
	//
	//	*MessageBody_TextContent
	//	*MessageBody_ImageContent
	//	*MessageBody_AudioContent
	//	*MessageBody_VideoContent
	Content       isMessageBody_Content `protobuf_oneof:"content"`
*/
func TestMessage(t *testing.T) {
	mb := &MessageBody{
		Id:       strings.ToUpper(id.GenerateXId()),
		AppId:    "1212",
		UserId:   1212,
		CId:      "sdf",
		To:       1111,
		GroupId:  100,
		Sequence: 100,
		Flow:     FlowRequest,
		NeedAck:  YES,
		CTime:    time.Now().UnixMilli(),
	}

	mb.Set(&TextContent{
		Text: "hello world",
	})

	js, _ := json.Marshal(mb)
	fmt.Println(string(js))

	j1, _ := json.Marshal(mb.Reply())
	fmt.Println(string(j1))

	j2, _ := json.Marshal(mb.Wrap())
	fmt.Println(string(j2))
}
