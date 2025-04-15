// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.29.3
// source: broker.proto

package api

import (
	context "context"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type DeliverRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MessageId  string   `protobuf:"bytes,1,opt,name=messageId,proto3" json:"messageId,omitempty"`
	UserLabels []string `protobuf:"bytes,2,rep,name=userLabels,proto3" json:"userLabels,omitempty"`
	Message    *Message `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *DeliverRequest) Reset() {
	*x = DeliverRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_broker_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeliverRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeliverRequest) ProtoMessage() {}

func (x *DeliverRequest) ProtoReflect() protoreflect.Message {
	mi := &file_broker_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeliverRequest.ProtoReflect.Descriptor instead.
func (*DeliverRequest) Descriptor() ([]byte, []int) {
	return file_broker_proto_rawDescGZIP(), []int{0}
}

func (x *DeliverRequest) GetMessageId() string {
	if x != nil {
		return x.MessageId
	}
	return ""
}

func (x *DeliverRequest) GetUserLabels() []string {
	if x != nil {
		return x.UserLabels
	}
	return nil
}

func (x *DeliverRequest) GetMessage() *Message {
	if x != nil {
		return x.Message
	}
	return nil
}

type DeliverReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MessageId string `protobuf:"bytes,1,opt,name=messageId,proto3" json:"messageId,omitempty"`
	Code      int32  `protobuf:"varint,2,opt,name=code,proto3" json:"code,omitempty"`
	Message   string `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *DeliverReply) Reset() {
	*x = DeliverReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_broker_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeliverReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeliverReply) ProtoMessage() {}

func (x *DeliverReply) ProtoReflect() protoreflect.Message {
	mi := &file_broker_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeliverReply.ProtoReflect.Descriptor instead.
func (*DeliverReply) Descriptor() ([]byte, []int) {
	return file_broker_proto_rawDescGZIP(), []int{1}
}

func (x *DeliverReply) GetMessageId() string {
	if x != nil {
		return x.MessageId
	}
	return ""
}

func (x *DeliverReply) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *DeliverReply) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_broker_proto protoreflect.FileDescriptor

var file_broker_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x62, 0x72, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03,
	0x61, 0x70, 0x69, 0x1a, 0x0c, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x76, 0x0a, 0x0e, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49,
	0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x75, 0x73, 0x65, 0x72, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x75, 0x73, 0x65, 0x72, 0x4c, 0x61, 0x62, 0x65, 0x6c,
	0x73, 0x12, 0x26, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x5a, 0x0a, 0x0c, 0x44, 0x65, 0x6c,
	0x69, 0x76, 0x65, 0x72, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x1c, 0x0a, 0x09, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x32, 0x44, 0x0a, 0x0d, 0x42, 0x72, 0x6f, 0x6b, 0x65, 0x72, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x33, 0x0a, 0x07, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65,
	0x72, 0x12, 0x13, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x44, 0x65, 0x6c,
	0x69, 0x76, 0x65, 0x72, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x42, 0x2e, 0x5a, 0x2c, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x61, 0x67, 0x69, 0x63, 0x6e,
	0x61, 0x6e, 0x61, 0x39, 0x39, 0x39, 0x2f, 0x69, 0x6d, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x6b, 0x69,
	0x74, 0x65, 0x78, 0x5f, 0x67, 0x65, 0x6e, 0x2f, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_broker_proto_rawDescOnce sync.Once
	file_broker_proto_rawDescData = file_broker_proto_rawDesc
)

func file_broker_proto_rawDescGZIP() []byte {
	file_broker_proto_rawDescOnce.Do(func() {
		file_broker_proto_rawDescData = protoimpl.X.CompressGZIP(file_broker_proto_rawDescData)
	})
	return file_broker_proto_rawDescData
}

var file_broker_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_broker_proto_goTypes = []interface{}{
	(*DeliverRequest)(nil), // 0: api.DeliverRequest
	(*DeliverReply)(nil),   // 1: api.DeliverReply
	(*Message)(nil),        // 2: api.Message
}
var file_broker_proto_depIdxs = []int32{
	2, // 0: api.DeliverRequest.message:type_name -> api.Message
	0, // 1: api.BrokerService.Deliver:input_type -> api.DeliverRequest
	1, // 2: api.BrokerService.Deliver:output_type -> api.DeliverReply
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_broker_proto_init() }
func file_broker_proto_init() {
	if File_broker_proto != nil {
		return
	}
	file_packet_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_broker_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeliverRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_broker_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeliverReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_broker_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_broker_proto_goTypes,
		DependencyIndexes: file_broker_proto_depIdxs,
		MessageInfos:      file_broker_proto_msgTypes,
	}.Build()
	File_broker_proto = out.File
	file_broker_proto_rawDesc = nil
	file_broker_proto_goTypes = nil
	file_broker_proto_depIdxs = nil
}

var _ context.Context

// Code generated by Kitex v0.12.3. DO NOT EDIT.

type BrokerService interface {
	Deliver(ctx context.Context, req *DeliverRequest) (res *DeliverReply, err error)
}
