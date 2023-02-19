// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.0--rc2
// source: proto/challenge.proto

package proto

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

// RemoteExecuteAPIClient is the client API for RemoteExecuteAPI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RemoteExecuteAPIClient interface {
	Execute(ctx context.Context, in *RequestExecute, opts ...grpc.CallOption) (*ResponseExecute, error)
}

type remoteExecuteAPIClient struct {
	cc grpc.ClientConnInterface
}

func NewRemoteExecuteAPIClient(cc grpc.ClientConnInterface) RemoteExecuteAPIClient {
	return &remoteExecuteAPIClient{cc}
}

func (c *remoteExecuteAPIClient) Execute(ctx context.Context, in *RequestExecute, opts ...grpc.CallOption) (*ResponseExecute, error) {
	out := new(ResponseExecute)
	err := c.cc.Invoke(ctx, "/com.github.anight.mysticchallenge.RemoteExecuteAPI/Execute", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RemoteExecuteAPIServer is the server API for RemoteExecuteAPI service.
// All implementations must embed UnimplementedRemoteExecuteAPIServer
// for forward compatibility
type RemoteExecuteAPIServer interface {
	Execute(context.Context, *RequestExecute) (*ResponseExecute, error)
	mustEmbedUnimplementedRemoteExecuteAPIServer()
}

// UnimplementedRemoteExecuteAPIServer must be embedded to have forward compatible implementations.
type UnimplementedRemoteExecuteAPIServer struct {
}

func (UnimplementedRemoteExecuteAPIServer) Execute(context.Context, *RequestExecute) (*ResponseExecute, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Execute not implemented")
}
func (UnimplementedRemoteExecuteAPIServer) mustEmbedUnimplementedRemoteExecuteAPIServer() {}

// UnsafeRemoteExecuteAPIServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RemoteExecuteAPIServer will
// result in compilation errors.
type UnsafeRemoteExecuteAPIServer interface {
	mustEmbedUnimplementedRemoteExecuteAPIServer()
}

func RegisterRemoteExecuteAPIServer(s grpc.ServiceRegistrar, srv RemoteExecuteAPIServer) {
	s.RegisterService(&RemoteExecuteAPI_ServiceDesc, srv)
}

func _RemoteExecuteAPI_Execute_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestExecute)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RemoteExecuteAPIServer).Execute(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.github.anight.mysticchallenge.RemoteExecuteAPI/Execute",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RemoteExecuteAPIServer).Execute(ctx, req.(*RequestExecute))
	}
	return interceptor(ctx, in, info, handler)
}

// RemoteExecuteAPI_ServiceDesc is the grpc.ServiceDesc for RemoteExecuteAPI service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RemoteExecuteAPI_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "com.github.anight.mysticchallenge.RemoteExecuteAPI",
	HandlerType: (*RemoteExecuteAPIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Execute",
			Handler:    _RemoteExecuteAPI_Execute_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/challenge.proto",
}
