//
//
//╰$ cd broker/pb
//╰$ protoc --go_out=. packet.proto
//╰$ protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative packet.proto

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v4.25.3
// source: packet.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	UserApi_Login_FullMethodName = "/pb.UserApi/Login"
)

// UserApiClient is the client API for UserApi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserApiClient interface {
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginReply, error)
}

type userApiClient struct {
	cc grpc.ClientConnInterface
}

func NewUserApiClient(cc grpc.ClientConnInterface) UserApiClient {
	return &userApiClient{cc}
}

func (c *userApiClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LoginReply)
	err := c.cc.Invoke(ctx, UserApi_Login_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserApiServer is the server API for UserApi service.
// All implementations must embed UnimplementedUserApiServer
// for forward compatibility.
type UserApiServer interface {
	Login(context.Context, *LoginRequest) (*LoginReply, error)
	mustEmbedUnimplementedUserApiServer()
}

// UnimplementedUserApiServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedUserApiServer struct{}

func (UnimplementedUserApiServer) Login(context.Context, *LoginRequest) (*LoginReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedUserApiServer) mustEmbedUnimplementedUserApiServer() {}
func (UnimplementedUserApiServer) testEmbeddedByValue()                 {}

// UnsafeUserApiServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserApiServer will
// result in compilation errors.
type UnsafeUserApiServer interface {
	mustEmbedUnimplementedUserApiServer()
}

func RegisterUserApiServer(s grpc.ServiceRegistrar, srv UserApiServer) {
	// If the following call pancis, it indicates UnimplementedUserApiServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&UserApi_ServiceDesc, srv)
}

func _UserApi_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserApiServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserApi_Login_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserApiServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserApi_ServiceDesc is the grpc.ServiceDesc for UserApi service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserApi_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.UserApi",
	HandlerType: (*UserApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Login",
			Handler:    _UserApi_Login_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "packet.proto",
}

const (
	MessageDeliverApi_Deliver_FullMethodName = "/pb.MessageDeliverApi/deliver"
)

// MessageDeliverApiClient is the client API for MessageDeliverApi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MessageDeliverApiClient interface {
	Deliver(ctx context.Context, in *CommandBody, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type messageDeliverApiClient struct {
	cc grpc.ClientConnInterface
}

func NewMessageDeliverApiClient(cc grpc.ClientConnInterface) MessageDeliverApiClient {
	return &messageDeliverApiClient{cc}
}

func (c *messageDeliverApiClient) Deliver(ctx context.Context, in *CommandBody, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, MessageDeliverApi_Deliver_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MessageDeliverApiServer is the server API for MessageDeliverApi service.
// All implementations must embed UnimplementedMessageDeliverApiServer
// for forward compatibility.
type MessageDeliverApiServer interface {
	Deliver(context.Context, *CommandBody) (*emptypb.Empty, error)
	mustEmbedUnimplementedMessageDeliverApiServer()
}

// UnimplementedMessageDeliverApiServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMessageDeliverApiServer struct{}

func (UnimplementedMessageDeliverApiServer) Deliver(context.Context, *CommandBody) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Deliver not implemented")
}
func (UnimplementedMessageDeliverApiServer) mustEmbedUnimplementedMessageDeliverApiServer() {}
func (UnimplementedMessageDeliverApiServer) testEmbeddedByValue()                           {}

// UnsafeMessageDeliverApiServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MessageDeliverApiServer will
// result in compilation errors.
type UnsafeMessageDeliverApiServer interface {
	mustEmbedUnimplementedMessageDeliverApiServer()
}

func RegisterMessageDeliverApiServer(s grpc.ServiceRegistrar, srv MessageDeliverApiServer) {
	// If the following call pancis, it indicates UnimplementedMessageDeliverApiServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&MessageDeliverApi_ServiceDesc, srv)
}

func _MessageDeliverApi_Deliver_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommandBody)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageDeliverApiServer).Deliver(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MessageDeliverApi_Deliver_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageDeliverApiServer).Deliver(ctx, req.(*CommandBody))
	}
	return interceptor(ctx, in, info, handler)
}

// MessageDeliverApi_ServiceDesc is the grpc.ServiceDesc for MessageDeliverApi service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MessageDeliverApi_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.MessageDeliverApi",
	HandlerType: (*MessageDeliverApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "deliver",
			Handler:    _MessageDeliverApi_Deliver_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "packet.proto",
}
