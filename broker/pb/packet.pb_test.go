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

/*
*

	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`

	AppId         string                 `protobuf:"bytes,2,opt,name=appId,proto3" json:"appId,omitempty"`
	UserId        string                 `protobuf:"bytes,3,opt,name=userId,proto3" json:"userId,omitempty"`
	Flow          uint32                 `protobuf:"varint,4,opt,name=flow,proto3" json:"flow,omitempty"`
	NeedAck       uint32                 `protobuf:"varint,5,opt,name=needAck,proto3" json:"needAck,omitempty"`
	Type          uint32                 `protobuf:"varint,6,opt,name=type,proto3" json:"type,omitempty"`
	CTime         uint64                 `protobuf:"varint,7,opt,name=cTime,proto3" json:"cTime,omitempty"`
	STime         uint64                 `protobuf:"varint,8,opt,name=sTime,proto3" json:"sTime,omitempty"`
*/

/*
*
MType         string                 `protobuf:"bytes,1,opt,name=mType,proto3" json:"mType,omitempty"`
CId           string                 `protobuf:"bytes,2,opt,name=cId,proto3" json:"cId,omitempty"`
To            string                 `protobuf:"bytes,3,opt,name=to,proto3" json:"to,omitempty"`
GroupId       string                 `protobuf:"bytes,4,opt,name=groupId,proto3" json:"groupId,omitempty"`
TType         int32                  `protobuf:"varint,5,opt,name=tType,proto3" json:"tType,omitempty"` // int8 in Go maps to int32 in Proto3
Sequence      int64                  `protobuf:"varint,6,opt,name=sequence,proto3" json:"sequence,omitempty"`
Content       *anypb.Any             `protobuf:"bytes,7,opt,name=content,proto3" json:"content,omitempty"` // For any type content
At            []*At                  `protobuf:"bytes,8,rep,name=at,proto3" json:"at,omitempty"`           // At list
Refer         []*Refer               `protobuf:"bytes,9,rep,name=refer,proto3" json:"refer,
*/
func Test(t *testing.T) {
	text, _ := anypb.New(&TextBody{
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
