// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v4.25.2
// source: sole/date/v1/date.service.proto

// Date Query API

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	DateService_Now_FullMethodName   = "/searking.sole.api.sole.date.v1.DateService/Now"
	DateService_Error_FullMethodName = "/searking.sole.api.sole.date.v1.DateService/Error"
)

// DateServiceClient is the client API for DateService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Date Service
type DateServiceClient interface {
	// Date Query
	Now(ctx context.Context, in *DateRequest, opts ...grpc.CallOption) (*DateResponse, error)
	// Date Query, only return error, for test only
	Error(ctx context.Context, in *DateRequest, opts ...grpc.CallOption) (*DateResponse, error)
}

type dateServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDateServiceClient(cc grpc.ClientConnInterface) DateServiceClient {
	return &dateServiceClient{cc}
}

func (c *dateServiceClient) Now(ctx context.Context, in *DateRequest, opts ...grpc.CallOption) (*DateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DateResponse)
	err := c.cc.Invoke(ctx, DateService_Now_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dateServiceClient) Error(ctx context.Context, in *DateRequest, opts ...grpc.CallOption) (*DateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DateResponse)
	err := c.cc.Invoke(ctx, DateService_Error_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DateServiceServer is the server API for DateService service.
// All implementations must embed UnimplementedDateServiceServer
// for forward compatibility.
//
// Date Service
type DateServiceServer interface {
	// Date Query
	Now(context.Context, *DateRequest) (*DateResponse, error)
	// Date Query, only return error, for test only
	Error(context.Context, *DateRequest) (*DateResponse, error)
	mustEmbedUnimplementedDateServiceServer()
}

// UnimplementedDateServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedDateServiceServer struct{}

func (UnimplementedDateServiceServer) Now(context.Context, *DateRequest) (*DateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Now not implemented")
}
func (UnimplementedDateServiceServer) Error(context.Context, *DateRequest) (*DateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Error not implemented")
}
func (UnimplementedDateServiceServer) mustEmbedUnimplementedDateServiceServer() {}
func (UnimplementedDateServiceServer) testEmbeddedByValue()                     {}

// UnsafeDateServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DateServiceServer will
// result in compilation errors.
type UnsafeDateServiceServer interface {
	mustEmbedUnimplementedDateServiceServer()
}

func RegisterDateServiceServer(s grpc.ServiceRegistrar, srv DateServiceServer) {
	// If the following call pancis, it indicates UnimplementedDateServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&DateService_ServiceDesc, srv)
}

func _DateService_Now_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DateServiceServer).Now(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DateService_Now_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DateServiceServer).Now(ctx, req.(*DateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DateService_Error_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DateServiceServer).Error(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DateService_Error_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DateServiceServer).Error(ctx, req.(*DateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DateService_ServiceDesc is the grpc.ServiceDesc for DateService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DateService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "searking.sole.api.sole.date.v1.DateService",
	HandlerType: (*DateServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Now",
			Handler:    _DateService_Now_Handler,
		},
		{
			MethodName: "Error",
			Handler:    _DateService_Error_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sole/date/v1/date.service.proto",
}
