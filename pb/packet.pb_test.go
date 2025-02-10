package pb

import (
	"testing"
)

func Test(t *testing.T) {
	//text, _ := anypb.New(&TextContent{
	//	Text: "hello",
	//})
	//
	//val := wrapperspb.UInt32(uint32(100))
	//anypb.New(val)
	//
	//message, _ := anypb.New(&MessageBody{
	//	MType:    protocol2.MText,
	//	CId:      "sjfjsdf",
	//	To:       "200",
	//	GroupId:  "",
	//	TType:    protocol2.TSingle,
	//	Sequence: 1212,
	//	Content:  text,
	//})
	//
	//packet := &Packet{
	//	Id:      id.GenerateXId(),
	//	AppId:   "100",
	//	UserId:  10011,
	//	Flow:    protocol2.FlowRequest,
	//	NeedAck: protocol2.YES,
	//	Type:    protocol2.TypeMessage,
	//	CTime:   time.Now().UnixMilli(),
	//	STime:   time.Now().UnixMilli(),
	//	Body:    message,
	//}
	//
	//bs, _ := proto.Marshal(packet)
	//fmt.Println(len(bs))
	//
	//proto.Unmarshal(bs, packet)
	//
	//var mb MessageBody
	//packet.Body.UnmarshalTo(&mb)
	//
	//fmt.Println(packet.Id)
	//fmt.Println(mb.Sequence)

}
