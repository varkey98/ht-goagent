// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package internal

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// PersonRegistryClient is the client API for PersonRegistry service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PersonRegistryClient interface {
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterReply, error)
}

type personRegistryClient struct {
	cc grpc.ClientConnInterface
}

func NewPersonRegistryClient(cc grpc.ClientConnInterface) PersonRegistryClient {
	return &personRegistryClient{cc}
}

func (c *personRegistryClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterReply, error) {
	out := new(RegisterReply)
	err := c.cc.Invoke(ctx, "/helloworld.PersonRegistry/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PersonRegistryServer is the server API for PersonRegistry service.
// All implementations must embed UnimplementedPersonRegistryServer
// for forward compatibility
type PersonRegistryServer interface {
	Register(context.Context, *RegisterRequest) (*RegisterReply, error)
	mustEmbedUnimplementedPersonRegistryServer()
}

// UnimplementedPersonRegistryServer must be embedded to have forward compatible implementations.
type UnimplementedPersonRegistryServer struct {
}

func (*UnimplementedPersonRegistryServer) Register(context.Context, *RegisterRequest) (*RegisterReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (*UnimplementedPersonRegistryServer) mustEmbedUnimplementedPersonRegistryServer() {}

func RegisterPersonRegistryServer(s *grpc.Server, srv PersonRegistryServer) {
	s.RegisterService(&_PersonRegistry_serviceDesc, srv)
}

func _PersonRegistry_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PersonRegistryServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/helloworld.PersonRegistry/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PersonRegistryServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _PersonRegistry_serviceDesc = grpc.ServiceDesc{
	ServiceName: "helloworld.PersonRegistry",
	HandlerType: (*PersonRegistryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _PersonRegistry_Register_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "person.proto",
}