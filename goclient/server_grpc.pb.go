// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package main

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// SocketGuideClient is the client API for SocketGuide service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SocketGuideClient interface {
	GetSocketFamilyList(ctx context.Context, in *SocketTree, opts ...grpc.CallOption) (SocketGuide_GetSocketFamilyListClient, error)
	GetSocketTypeList(ctx context.Context, in *SocketFamily, opts ...grpc.CallOption) (SocketGuide_GetSocketTypeListClient, error)
	GetSocketProtocolList(ctx context.Context, in *SocketTypeAndFamily, opts ...grpc.CallOption) (SocketGuide_GetSocketProtocolListClient, error)
}

type socketGuideClient struct {
	cc grpc.ClientConnInterface
}

func NewSocketGuideClient(cc grpc.ClientConnInterface) SocketGuideClient {
	return &socketGuideClient{cc}
}

func (c *socketGuideClient) GetSocketFamilyList(ctx context.Context, in *SocketTree, opts ...grpc.CallOption) (SocketGuide_GetSocketFamilyListClient, error) {
	stream, err := c.cc.NewStream(ctx, &SocketGuide_ServiceDesc.Streams[0], "/main.SocketGuide/GetSocketFamilyList", opts...)
	if err != nil {
		return nil, err
	}
	x := &socketGuideGetSocketFamilyListClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type SocketGuide_GetSocketFamilyListClient interface {
	Recv() (*SocketFamily, error)
	grpc.ClientStream
}

type socketGuideGetSocketFamilyListClient struct {
	grpc.ClientStream
}

func (x *socketGuideGetSocketFamilyListClient) Recv() (*SocketFamily, error) {
	m := new(SocketFamily)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *socketGuideClient) GetSocketTypeList(ctx context.Context, in *SocketFamily, opts ...grpc.CallOption) (SocketGuide_GetSocketTypeListClient, error) {
	stream, err := c.cc.NewStream(ctx, &SocketGuide_ServiceDesc.Streams[1], "/main.SocketGuide/GetSocketTypeList", opts...)
	if err != nil {
		return nil, err
	}
	x := &socketGuideGetSocketTypeListClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type SocketGuide_GetSocketTypeListClient interface {
	Recv() (*SocketType, error)
	grpc.ClientStream
}

type socketGuideGetSocketTypeListClient struct {
	grpc.ClientStream
}

func (x *socketGuideGetSocketTypeListClient) Recv() (*SocketType, error) {
	m := new(SocketType)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *socketGuideClient) GetSocketProtocolList(ctx context.Context, in *SocketTypeAndFamily, opts ...grpc.CallOption) (SocketGuide_GetSocketProtocolListClient, error) {
	stream, err := c.cc.NewStream(ctx, &SocketGuide_ServiceDesc.Streams[2], "/main.SocketGuide/GetSocketProtocolList", opts...)
	if err != nil {
		return nil, err
	}
	x := &socketGuideGetSocketProtocolListClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type SocketGuide_GetSocketProtocolListClient interface {
	Recv() (*SocketProtocol, error)
	grpc.ClientStream
}

type socketGuideGetSocketProtocolListClient struct {
	grpc.ClientStream
}

func (x *socketGuideGetSocketProtocolListClient) Recv() (*SocketProtocol, error) {
	m := new(SocketProtocol)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SocketGuideServer is the server API for SocketGuide service.
// All implementations must embed UnimplementedSocketGuideServer
// for forward compatibility
type SocketGuideServer interface {
	GetSocketFamilyList(*SocketTree, SocketGuide_GetSocketFamilyListServer) error
	GetSocketTypeList(*SocketFamily, SocketGuide_GetSocketTypeListServer) error
	GetSocketProtocolList(*SocketTypeAndFamily, SocketGuide_GetSocketProtocolListServer) error
	mustEmbedUnimplementedSocketGuideServer()
}

// UnimplementedSocketGuideServer must be embedded to have forward compatible implementations.
type UnimplementedSocketGuideServer struct {
}

func (UnimplementedSocketGuideServer) GetSocketFamilyList(*SocketTree, SocketGuide_GetSocketFamilyListServer) error {
	return status.Errorf(codes.Unimplemented, "method GetSocketFamilyList not implemented")
}
func (UnimplementedSocketGuideServer) GetSocketTypeList(*SocketFamily, SocketGuide_GetSocketTypeListServer) error {
	return status.Errorf(codes.Unimplemented, "method GetSocketTypeList not implemented")
}
func (UnimplementedSocketGuideServer) GetSocketProtocolList(*SocketTypeAndFamily, SocketGuide_GetSocketProtocolListServer) error {
	return status.Errorf(codes.Unimplemented, "method GetSocketProtocolList not implemented")
}
func (UnimplementedSocketGuideServer) mustEmbedUnimplementedSocketGuideServer() {}

// UnsafeSocketGuideServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SocketGuideServer will
// result in compilation errors.
type UnsafeSocketGuideServer interface {
	mustEmbedUnimplementedSocketGuideServer()
}

func RegisterSocketGuideServer(s grpc.ServiceRegistrar, srv SocketGuideServer) {
	s.RegisterService(&SocketGuide_ServiceDesc, srv)
}

func _SocketGuide_GetSocketFamilyList_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SocketTree)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SocketGuideServer).GetSocketFamilyList(m, &socketGuideGetSocketFamilyListServer{stream})
}

type SocketGuide_GetSocketFamilyListServer interface {
	Send(*SocketFamily) error
	grpc.ServerStream
}

type socketGuideGetSocketFamilyListServer struct {
	grpc.ServerStream
}

func (x *socketGuideGetSocketFamilyListServer) Send(m *SocketFamily) error {
	return x.ServerStream.SendMsg(m)
}

func _SocketGuide_GetSocketTypeList_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SocketFamily)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SocketGuideServer).GetSocketTypeList(m, &socketGuideGetSocketTypeListServer{stream})
}

type SocketGuide_GetSocketTypeListServer interface {
	Send(*SocketType) error
	grpc.ServerStream
}

type socketGuideGetSocketTypeListServer struct {
	grpc.ServerStream
}

func (x *socketGuideGetSocketTypeListServer) Send(m *SocketType) error {
	return x.ServerStream.SendMsg(m)
}

func _SocketGuide_GetSocketProtocolList_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SocketTypeAndFamily)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SocketGuideServer).GetSocketProtocolList(m, &socketGuideGetSocketProtocolListServer{stream})
}

type SocketGuide_GetSocketProtocolListServer interface {
	Send(*SocketProtocol) error
	grpc.ServerStream
}

type socketGuideGetSocketProtocolListServer struct {
	grpc.ServerStream
}

func (x *socketGuideGetSocketProtocolListServer) Send(m *SocketProtocol) error {
	return x.ServerStream.SendMsg(m)
}

// SocketGuide_ServiceDesc is the grpc.ServiceDesc for SocketGuide service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SocketGuide_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "main.SocketGuide",
	HandlerType: (*SocketGuideServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetSocketFamilyList",
			Handler:       _SocketGuide_GetSocketFamilyList_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "GetSocketTypeList",
			Handler:       _SocketGuide_GetSocketTypeList_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "GetSocketProtocolList",
			Handler:       _SocketGuide_GetSocketProtocolList_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "server.proto",
}
