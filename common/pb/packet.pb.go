//
//
//╰$ cd broker/pb
//╰$ protoc --go_out=. packet.proto
//╰$ protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative packet.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.2
// 	protoc        v4.25.3
// source: packet.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type OSType int32

const (
	OSType_UNKNOWN   OSType = 0
	OSType_OSWindows OSType = 1
	OSType_MacOS     OSType = 2
	OSType_LinuxOS   OSType = 3
	OSType_OSIos     OSType = 4
	OSType_Xiaomi    OSType = 5
	OSType_Huawei    OSType = 6
	OSType_Samsung   OSType = 7
	OSType_Honor     OSType = 8
	OSType_Oppo      OSType = 9
	OSType_Vivo      OSType = 10
)

// Enum value maps for OSType.
var (
	OSType_name = map[int32]string{
		0:  "UNKNOWN",
		1:  "OSWindows",
		2:  "MacOS",
		3:  "LinuxOS",
		4:  "OSIos",
		5:  "Xiaomi",
		6:  "Huawei",
		7:  "Samsung",
		8:  "Honor",
		9:  "Oppo",
		10: "Vivo",
	}
	OSType_value = map[string]int32{
		"UNKNOWN":   0,
		"OSWindows": 1,
		"MacOS":     2,
		"LinuxOS":   3,
		"OSIos":     4,
		"Xiaomi":    5,
		"Huawei":    6,
		"Samsung":   7,
		"Honor":     8,
		"Oppo":      9,
		"Vivo":      10,
	}
)

func (x OSType) Enum() *OSType {
	p := new(OSType)
	*p = x
	return p
}

func (x OSType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (OSType) Descriptor() protoreflect.EnumDescriptor {
	return file_packet_proto_enumTypes[0].Descriptor()
}

func (OSType) Type() protoreflect.EnumType {
	return &file_packet_proto_enumTypes[0]
}

func (x OSType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use OSType.Descriptor instead.
func (OSType) EnumDescriptor() ([]byte, []int) {
	return file_packet_proto_rawDescGZIP(), []int{0}
}

type Packet struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	AppId         string                 `protobuf:"bytes,2,opt,name=appId,proto3" json:"appId,omitempty"`
	UserId        int64                  `protobuf:"varint,3,opt,name=userId,proto3" json:"userId,omitempty"`
	Flow          int32                  `protobuf:"varint,4,opt,name=flow,proto3" json:"flow,omitempty"`
	NeedAck       int32                  `protobuf:"varint,5,opt,name=needAck,proto3" json:"needAck,omitempty"`
	CTime         int64                  `protobuf:"varint,6,opt,name=cTime,proto3" json:"cTime,omitempty"`
	STime         int64                  `protobuf:"varint,7,opt,name=sTime,proto3" json:"sTime,omitempty"`
	BType         int32                  `protobuf:"varint,8,opt,name=bType,proto3" json:"bType,omitempty"`
	Body          *anypb.Any             `protobuf:"bytes,9,opt,name=body,proto3" json:"body,omitempty"` // for the "any" type field
	Status        *Status                `protobuf:"bytes,10,opt,name=status,proto3" json:"status,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Packet) Reset() {
	*x = Packet{}
	mi := &file_packet_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Packet) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Packet) ProtoMessage() {}

func (x *Packet) ProtoReflect() protoreflect.Message {
	mi := &file_packet_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Packet.ProtoReflect.Descriptor instead.
func (*Packet) Descriptor() ([]byte, []int) {
	return file_packet_proto_rawDescGZIP(), []int{0}
}

func (x *Packet) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Packet) GetAppId() string {
	if x != nil {
		return x.AppId
	}
	return ""
}

func (x *Packet) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *Packet) GetFlow() int32 {
	if x != nil {
		return x.Flow
	}
	return 0
}

func (x *Packet) GetNeedAck() int32 {
	if x != nil {
		return x.NeedAck
	}
	return 0
}

func (x *Packet) GetCTime() int64 {
	if x != nil {
		return x.CTime
	}
	return 0
}

func (x *Packet) GetSTime() int64 {
	if x != nil {
		return x.STime
	}
	return 0
}

func (x *Packet) GetBType() int32 {
	if x != nil {
		return x.BType
	}
	return 0
}

func (x *Packet) GetBody() *anypb.Any {
	if x != nil {
		return x.Body
	}
	return nil
}

func (x *Packet) GetStatus() *Status {
	if x != nil {
		return x.Status
	}
	return nil
}

type MessageBody struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CType         string                 `protobuf:"bytes,1,opt,name=cType,proto3" json:"cType,omitempty"`
	CId           string                 `protobuf:"bytes,2,opt,name=cId,proto3" json:"cId,omitempty"`
	To            string                 `protobuf:"bytes,3,opt,name=to,proto3" json:"to,omitempty"`
	GroupId       string                 `protobuf:"bytes,4,opt,name=groupId,proto3" json:"groupId,omitempty"`
	TargetType    int32                  `protobuf:"varint,5,opt,name=targetType,proto3" json:"targetType,omitempty"`
	Sequence      int64                  `protobuf:"varint,6,opt,name=sequence,proto3" json:"sequence,omitempty"`
	Content       *anypb.Any             `protobuf:"bytes,7,opt,name=content,proto3" json:"content,omitempty"`
	At            []*At                  `protobuf:"bytes,8,rep,name=at,proto3" json:"at,omitempty"`
	Refer         []*Refer               `protobuf:"bytes,9,rep,name=refer,proto3" json:"refer,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MessageBody) Reset() {
	*x = MessageBody{}
	mi := &file_packet_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MessageBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageBody) ProtoMessage() {}

func (x *MessageBody) ProtoReflect() protoreflect.Message {
	mi := &file_packet_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageBody.ProtoReflect.Descriptor instead.
func (*MessageBody) Descriptor() ([]byte, []int) {
	return file_packet_proto_rawDescGZIP(), []int{1}
}

func (x *MessageBody) GetCType() string {
	if x != nil {
		return x.CType
	}
	return ""
}

func (x *MessageBody) GetCId() string {
	if x != nil {
		return x.CId
	}
	return ""
}

func (x *MessageBody) GetTo() string {
	if x != nil {
		return x.To
	}
	return ""
}

func (x *MessageBody) GetGroupId() string {
	if x != nil {
		return x.GroupId
	}
	return ""
}

func (x *MessageBody) GetTargetType() int32 {
	if x != nil {
		return x.TargetType
	}
	return 0
}

func (x *MessageBody) GetSequence() int64 {
	if x != nil {
		return x.Sequence
	}
	return 0
}

func (x *MessageBody) GetContent() *anypb.Any {
	if x != nil {
		return x.Content
	}
	return nil
}

func (x *MessageBody) GetAt() []*At {
	if x != nil {
		return x.At
	}
	return nil
}

func (x *MessageBody) GetRefer() []*Refer {
	if x != nil {
		return x.Refer
	}
	return nil
}

type At struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        int64                  `protobuf:"varint,1,opt,name=userId,proto3" json:"userId,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Avatar        string                 `protobuf:"bytes,3,opt,name=avatar,proto3" json:"avatar,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *At) Reset() {
	*x = At{}
	mi := &file_packet_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *At) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*At) ProtoMessage() {}

func (x *At) ProtoReflect() protoreflect.Message {
	mi := &file_packet_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use At.ProtoReflect.Descriptor instead.
func (*At) Descriptor() ([]byte, []int) {
	return file_packet_proto_rawDescGZIP(), []int{2}
}

func (x *At) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *At) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *At) GetAvatar() string {
	if x != nil {
		return x.Avatar
	}
	return ""
}

type Refer struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        int64                  `protobuf:"varint,1,opt,name=userId,proto3" json:"userId,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Avatar        string                 `protobuf:"bytes,3,opt,name=avatar,proto3" json:"avatar,omitempty"`
	CType         string                 `protobuf:"bytes,4,opt,name=cType,proto3" json:"cType,omitempty"`
	Content       *anypb.Any             `protobuf:"bytes,5,opt,name=content,proto3" json:"content,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Refer) Reset() {
	*x = Refer{}
	mi := &file_packet_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Refer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Refer) ProtoMessage() {}

func (x *Refer) ProtoReflect() protoreflect.Message {
	mi := &file_packet_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Refer.ProtoReflect.Descriptor instead.
func (*Refer) Descriptor() ([]byte, []int) {
	return file_packet_proto_rawDescGZIP(), []int{3}
}

func (x *Refer) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *Refer) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Refer) GetAvatar() string {
	if x != nil {
		return x.Avatar
	}
	return ""
}

func (x *Refer) GetCType() string {
	if x != nil {
		return x.CType
	}
	return ""
}

func (x *Refer) GetContent() *anypb.Any {
	if x != nil {
		return x.Content
	}
	return nil
}

type TextContent struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Text          string                 `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TextContent) Reset() {
	*x = TextContent{}
	mi := &file_packet_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TextContent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TextContent) ProtoMessage() {}

func (x *TextContent) ProtoReflect() protoreflect.Message {
	mi := &file_packet_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TextContent.ProtoReflect.Descriptor instead.
func (*TextContent) Descriptor() ([]byte, []int) {
	return file_packet_proto_rawDescGZIP(), []int{4}
}

func (x *TextContent) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

type ImageContent struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Url           string                 `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	Width         int32                  `protobuf:"varint,2,opt,name=width,proto3" json:"width,omitempty"`
	Height        int32                  `protobuf:"varint,3,opt,name=height,proto3" json:"height,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ImageContent) Reset() {
	*x = ImageContent{}
	mi := &file_packet_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ImageContent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ImageContent) ProtoMessage() {}

func (x *ImageContent) ProtoReflect() protoreflect.Message {
	mi := &file_packet_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ImageContent.ProtoReflect.Descriptor instead.
func (*ImageContent) Descriptor() ([]byte, []int) {
	return file_packet_proto_rawDescGZIP(), []int{5}
}

func (x *ImageContent) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *ImageContent) GetWidth() int32 {
	if x != nil {
		return x.Width
	}
	return 0
}

func (x *ImageContent) GetHeight() int32 {
	if x != nil {
		return x.Height
	}
	return 0
}

type AudioContent struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Url           string                 `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	Length        int32                  `protobuf:"varint,2,opt,name=length,proto3" json:"length,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AudioContent) Reset() {
	*x = AudioContent{}
	mi := &file_packet_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AudioContent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AudioContent) ProtoMessage() {}

func (x *AudioContent) ProtoReflect() protoreflect.Message {
	mi := &file_packet_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AudioContent.ProtoReflect.Descriptor instead.
func (*AudioContent) Descriptor() ([]byte, []int) {
	return file_packet_proto_rawDescGZIP(), []int{6}
}

func (x *AudioContent) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *AudioContent) GetLength() int32 {
	if x != nil {
		return x.Length
	}
	return 0
}

type VideoContent struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Url           string                 `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	Cover         string                 `protobuf:"bytes,2,opt,name=cover,proto3" json:"cover,omitempty"`
	Length        int32                  `protobuf:"varint,3,opt,name=length,proto3" json:"length,omitempty"`
	Width         int32                  `protobuf:"varint,4,opt,name=width,proto3" json:"width,omitempty"`
	Height        int32                  `protobuf:"varint,5,opt,name=height,proto3" json:"height,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *VideoContent) Reset() {
	*x = VideoContent{}
	mi := &file_packet_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *VideoContent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VideoContent) ProtoMessage() {}

func (x *VideoContent) ProtoReflect() protoreflect.Message {
	mi := &file_packet_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VideoContent.ProtoReflect.Descriptor instead.
func (*VideoContent) Descriptor() ([]byte, []int) {
	return file_packet_proto_rawDescGZIP(), []int{7}
}

func (x *VideoContent) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *VideoContent) GetCover() string {
	if x != nil {
		return x.Cover
	}
	return ""
}

func (x *VideoContent) GetLength() int32 {
	if x != nil {
		return x.Length
	}
	return 0
}

func (x *VideoContent) GetWidth() int32 {
	if x != nil {
		return x.Width
	}
	return 0
}

func (x *VideoContent) GetHeight() int32 {
	if x != nil {
		return x.Height
	}
	return 0
}

type CommandBody struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CType         string                 `protobuf:"bytes,1,opt,name=cType,proto3" json:"cType,omitempty"`
	Request       *anypb.Any             `protobuf:"bytes,2,opt,name=Request,proto3" json:"Request,omitempty"`
	Reply         *anypb.Any             `protobuf:"bytes,3,opt,name=Reply,proto3" json:"Reply,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CommandBody) Reset() {
	*x = CommandBody{}
	mi := &file_packet_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CommandBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommandBody) ProtoMessage() {}

func (x *CommandBody) ProtoReflect() protoreflect.Message {
	mi := &file_packet_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommandBody.ProtoReflect.Descriptor instead.
func (*CommandBody) Descriptor() ([]byte, []int) {
	return file_packet_proto_rawDescGZIP(), []int{8}
}

func (x *CommandBody) GetCType() string {
	if x != nil {
		return x.CType
	}
	return ""
}

func (x *CommandBody) GetRequest() *anypb.Any {
	if x != nil {
		return x.Request
	}
	return nil
}

func (x *CommandBody) GetReply() *anypb.Any {
	if x != nil {
		return x.Reply
	}
	return nil
}

type LoginRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	AppId         string                 `protobuf:"bytes,1,opt,name=appId,proto3" json:"appId,omitempty"`
	UserSig       string                 `protobuf:"bytes,2,opt,name=userSig,proto3" json:"userSig,omitempty"`
	Version       string                 `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`
	Os            OSType                 `protobuf:"varint,4,opt,name=os,proto3,enum=pb.OSType" json:"os,omitempty"`
	PushDeviceId  string                 `protobuf:"bytes,5,opt,name=pushDeviceId,proto3" json:"pushDeviceId,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LoginRequest) Reset() {
	*x = LoginRequest{}
	mi := &file_packet_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LoginRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginRequest) ProtoMessage() {}

func (x *LoginRequest) ProtoReflect() protoreflect.Message {
	mi := &file_packet_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginRequest.ProtoReflect.Descriptor instead.
func (*LoginRequest) Descriptor() ([]byte, []int) {
	return file_packet_proto_rawDescGZIP(), []int{9}
}

func (x *LoginRequest) GetAppId() string {
	if x != nil {
		return x.AppId
	}
	return ""
}

func (x *LoginRequest) GetUserSig() string {
	if x != nil {
		return x.UserSig
	}
	return ""
}

func (x *LoginRequest) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *LoginRequest) GetOs() OSType {
	if x != nil {
		return x.Os
	}
	return OSType_UNKNOWN
}

func (x *LoginRequest) GetPushDeviceId() string {
	if x != nil {
		return x.PushDeviceId
	}
	return ""
}

type LoginReply struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	AppId         string                 `protobuf:"bytes,1,opt,name=appId,proto3" json:"appId,omitempty"`
	UserId        int64                  `protobuf:"varint,2,opt,name=userId,proto3" json:"userId,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LoginReply) Reset() {
	*x = LoginReply{}
	mi := &file_packet_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LoginReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginReply) ProtoMessage() {}

func (x *LoginReply) ProtoReflect() protoreflect.Message {
	mi := &file_packet_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginReply.ProtoReflect.Descriptor instead.
func (*LoginReply) Descriptor() ([]byte, []int) {
	return file_packet_proto_rawDescGZIP(), []int{10}
}

func (x *LoginReply) GetAppId() string {
	if x != nil {
		return x.AppId
	}
	return ""
}

func (x *LoginReply) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type Status struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Code          int32                  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Status) Reset() {
	*x = Status{}
	mi := &file_packet_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Status) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Status) ProtoMessage() {}

func (x *Status) ProtoReflect() protoreflect.Message {
	mi := &file_packet_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Status.ProtoReflect.Descriptor instead.
func (*Status) Descriptor() ([]byte, []int) {
	return file_packet_proto_rawDescGZIP(), []int{11}
}

func (x *Status) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *Status) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_packet_proto protoreflect.FileDescriptor

var file_packet_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02,
	0x70, 0x62, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x84, 0x02,
	0x0a, 0x06, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x70, 0x70, 0x49,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x61, 0x70, 0x70, 0x49, 0x64, 0x12, 0x16,
	0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06,
	0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x6c, 0x6f, 0x77, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x66, 0x6c, 0x6f, 0x77, 0x12, 0x18, 0x0a, 0x07, 0x6e, 0x65,
	0x65, 0x64, 0x41, 0x63, 0x6b, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x6e, 0x65, 0x65,
	0x64, 0x41, 0x63, 0x6b, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x05, 0x63, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x54,
	0x69, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x73, 0x54, 0x69, 0x6d, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x62, 0x54, 0x79, 0x70, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x05, 0x62, 0x54, 0x79, 0x70, 0x65, 0x12, 0x28, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x09,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79,
	0x12, 0x22, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0a, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x22, 0x84, 0x02, 0x0a, 0x0b, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x42, 0x6f, 0x64, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x63, 0x54, 0x79, 0x70, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x63, 0x49,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x63, 0x49, 0x64, 0x12, 0x0e, 0x0a, 0x02,
	0x74, 0x6f, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x74, 0x6f, 0x12, 0x18, 0x0a, 0x07,
	0x67, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x67,
	0x72, 0x6f, 0x75, 0x70, 0x49, 0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74,
	0x54, 0x79, 0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x74, 0x61, 0x72, 0x67,
	0x65, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e,
	0x63, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e,
	0x63, 0x65, 0x12, 0x2e, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x12, 0x16, 0x0a, 0x02, 0x61, 0x74, 0x18, 0x08, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x06,
	0x2e, 0x70, 0x62, 0x2e, 0x41, 0x74, 0x52, 0x02, 0x61, 0x74, 0x12, 0x1f, 0x0a, 0x05, 0x72, 0x65,
	0x66, 0x65, 0x72, 0x18, 0x09, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x70, 0x62, 0x2e, 0x52,
	0x65, 0x66, 0x65, 0x72, 0x52, 0x05, 0x72, 0x65, 0x66, 0x65, 0x72, 0x22, 0x48, 0x0a, 0x02, 0x41,
	0x74, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a,
	0x06, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61,
	0x76, 0x61, 0x74, 0x61, 0x72, 0x22, 0x91, 0x01, 0x0a, 0x05, 0x52, 0x65, 0x66, 0x65, 0x72, 0x12,
	0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x61,
	0x76, 0x61, 0x74, 0x61, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x76, 0x61,
	0x74, 0x61, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x54, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x63, 0x54, 0x79, 0x70, 0x65, 0x12, 0x2e, 0x0a, 0x07, 0x63, 0x6f, 0x6e,
	0x74, 0x65, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79,
	0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x21, 0x0a, 0x0b, 0x54, 0x65, 0x78,
	0x74, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x22, 0x4e, 0x0a, 0x0c,
	0x49, 0x6d, 0x61, 0x67, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x10, 0x0a, 0x03,
	0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x14,
	0x0a, 0x05, 0x77, 0x69, 0x64, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x77,
	0x69, 0x64, 0x74, 0x68, 0x12, 0x16, 0x0a, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x22, 0x38, 0x0a, 0x0c,
	0x41, 0x75, 0x64, 0x69, 0x6f, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x10, 0x0a, 0x03,
	0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x16,
	0x0a, 0x06, 0x6c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06,
	0x6c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x22, 0x7c, 0x0a, 0x0c, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x43,
	0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x76, 0x65,
	0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x12, 0x16,
	0x0a, 0x06, 0x6c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06,
	0x6c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x12, 0x14, 0x0a, 0x05, 0x77, 0x69, 0x64, 0x74, 0x68, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x77, 0x69, 0x64, 0x74, 0x68, 0x12, 0x16, 0x0a, 0x06,
	0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x68, 0x65,
	0x69, 0x67, 0x68, 0x74, 0x22, 0x7f, 0x0a, 0x0b, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x42,
	0x6f, 0x64, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x63, 0x54, 0x79, 0x70, 0x65, 0x12, 0x2e, 0x0a, 0x07, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79,
	0x52, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2a, 0x0a, 0x05, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x05,
	0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x98, 0x01, 0x0a, 0x0c, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x70, 0x70, 0x49, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x61, 0x70, 0x70, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07,
	0x75, 0x73, 0x65, 0x72, 0x53, 0x69, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x75,
	0x73, 0x65, 0x72, 0x53, 0x69, 0x67, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x12, 0x1a, 0x0a, 0x02, 0x6f, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0a, 0x2e, 0x70,
	0x62, 0x2e, 0x4f, 0x53, 0x54, 0x79, 0x70, 0x65, 0x52, 0x02, 0x6f, 0x73, 0x12, 0x22, 0x0a, 0x0c,
	0x70, 0x75, 0x73, 0x68, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0c, 0x70, 0x75, 0x73, 0x68, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64,
	0x22, 0x3a, 0x0a, 0x0a, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x14,
	0x0a, 0x05, 0x61, 0x70, 0x70, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x61,
	0x70, 0x70, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x36, 0x0a, 0x06,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x2a, 0x8b, 0x01, 0x0a, 0x06, 0x4f, 0x53, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x0d, 0x0a, 0x09,
	0x4f, 0x53, 0x57, 0x69, 0x6e, 0x64, 0x6f, 0x77, 0x73, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x4d,
	0x61, 0x63, 0x4f, 0x53, 0x10, 0x02, 0x12, 0x0b, 0x0a, 0x07, 0x4c, 0x69, 0x6e, 0x75, 0x78, 0x4f,
	0x53, 0x10, 0x03, 0x12, 0x09, 0x0a, 0x05, 0x4f, 0x53, 0x49, 0x6f, 0x73, 0x10, 0x04, 0x12, 0x0a,
	0x0a, 0x06, 0x58, 0x69, 0x61, 0x6f, 0x6d, 0x69, 0x10, 0x05, 0x12, 0x0a, 0x0a, 0x06, 0x48, 0x75,
	0x61, 0x77, 0x65, 0x69, 0x10, 0x06, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x61, 0x6d, 0x73, 0x75, 0x6e,
	0x67, 0x10, 0x07, 0x12, 0x09, 0x0a, 0x05, 0x48, 0x6f, 0x6e, 0x6f, 0x72, 0x10, 0x08, 0x12, 0x08,
	0x0a, 0x04, 0x4f, 0x70, 0x70, 0x6f, 0x10, 0x09, 0x12, 0x08, 0x0a, 0x04, 0x56, 0x69, 0x76, 0x6f,
	0x10, 0x0a, 0x32, 0x36, 0x0a, 0x07, 0x55, 0x73, 0x65, 0x72, 0x41, 0x70, 0x69, 0x12, 0x2b, 0x0a,
	0x05, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x10, 0x2e, 0x70, 0x62, 0x2e, 0x4c, 0x6f, 0x67, 0x69,
	0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x70, 0x62, 0x2e, 0x4c, 0x6f,
	0x67, 0x69, 0x6e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x42, 0x07, 0x5a, 0x05, 0x2e, 0x2e,
	0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_packet_proto_rawDescOnce sync.Once
	file_packet_proto_rawDescData = file_packet_proto_rawDesc
)

func file_packet_proto_rawDescGZIP() []byte {
	file_packet_proto_rawDescOnce.Do(func() {
		file_packet_proto_rawDescData = protoimpl.X.CompressGZIP(file_packet_proto_rawDescData)
	})
	return file_packet_proto_rawDescData
}

var file_packet_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_packet_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_packet_proto_goTypes = []any{
	(OSType)(0),          // 0: pb.OSType
	(*Packet)(nil),       // 1: pb.Packet
	(*MessageBody)(nil),  // 2: pb.MessageBody
	(*At)(nil),           // 3: pb.At
	(*Refer)(nil),        // 4: pb.Refer
	(*TextContent)(nil),  // 5: pb.TextContent
	(*ImageContent)(nil), // 6: pb.ImageContent
	(*AudioContent)(nil), // 7: pb.AudioContent
	(*VideoContent)(nil), // 8: pb.VideoContent
	(*CommandBody)(nil),  // 9: pb.CommandBody
	(*LoginRequest)(nil), // 10: pb.LoginRequest
	(*LoginReply)(nil),   // 11: pb.LoginReply
	(*Status)(nil),       // 12: pb.Status
	(*anypb.Any)(nil),    // 13: google.protobuf.Any
}
var file_packet_proto_depIdxs = []int32{
	13, // 0: pb.Packet.body:type_name -> google.protobuf.Any
	12, // 1: pb.Packet.status:type_name -> pb.Status
	13, // 2: pb.MessageBody.content:type_name -> google.protobuf.Any
	3,  // 3: pb.MessageBody.at:type_name -> pb.At
	4,  // 4: pb.MessageBody.refer:type_name -> pb.Refer
	13, // 5: pb.Refer.content:type_name -> google.protobuf.Any
	13, // 6: pb.CommandBody.Request:type_name -> google.protobuf.Any
	13, // 7: pb.CommandBody.Reply:type_name -> google.protobuf.Any
	0,  // 8: pb.LoginRequest.os:type_name -> pb.OSType
	10, // 9: pb.UserApi.Login:input_type -> pb.LoginRequest
	11, // 10: pb.UserApi.Login:output_type -> pb.LoginReply
	10, // [10:11] is the sub-list for method output_type
	9,  // [9:10] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_packet_proto_init() }
func file_packet_proto_init() {
	if File_packet_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_packet_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_packet_proto_goTypes,
		DependencyIndexes: file_packet_proto_depIdxs,
		EnumInfos:         file_packet_proto_enumTypes,
		MessageInfos:      file_packet_proto_msgTypes,
	}.Build()
	File_packet_proto = out.File
	file_packet_proto_rawDesc = nil
	file_packet_proto_goTypes = nil
	file_packet_proto_depIdxs = nil
}
