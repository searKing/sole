// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package swagger

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	httpbody "google.golang.org/genproto/googleapis/api/httpbody"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// SwaggerServiceClient is the client API for SwaggerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SwaggerServiceClient interface {
	// 静态Swagger JSON
	Json(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*httpbody.HttpBody, error)
	// 静态Swagger YAML
	Yaml(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*httpbody.HttpBody, error)
	// 静态Swagger UI
	UI(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*httpbody.HttpBody, error)
}

type swaggerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSwaggerServiceClient(cc grpc.ClientConnInterface) SwaggerServiceClient {
	return &swaggerServiceClient{cc}
}

func (c *swaggerServiceClient) Json(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*httpbody.HttpBody, error) {
	out := new(httpbody.HttpBody)
	err := c.cc.Invoke(ctx, "/sole.api.v1.doc.swagger.SwaggerService/Json", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *swaggerServiceClient) Yaml(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*httpbody.HttpBody, error) {
	out := new(httpbody.HttpBody)
	err := c.cc.Invoke(ctx, "/sole.api.v1.doc.swagger.SwaggerService/Yaml", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *swaggerServiceClient) UI(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*httpbody.HttpBody, error) {
	out := new(httpbody.HttpBody)
	err := c.cc.Invoke(ctx, "/sole.api.v1.doc.swagger.SwaggerService/UI", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SwaggerServiceServer is the server API for SwaggerService service.
// All implementations must embed UnimplementedSwaggerServiceServer
// for forward compatibility
type SwaggerServiceServer interface {
	// 静态Swagger JSON
	Json(context.Context, *empty.Empty) (*httpbody.HttpBody, error)
	// 静态Swagger YAML
	Yaml(context.Context, *empty.Empty) (*httpbody.HttpBody, error)
	// 静态Swagger UI
	UI(context.Context, *empty.Empty) (*httpbody.HttpBody, error)
	mustEmbedUnimplementedSwaggerServiceServer()
}

// UnimplementedSwaggerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSwaggerServiceServer struct {
}

func (UnimplementedSwaggerServiceServer) Json(context.Context, *empty.Empty) (*httpbody.HttpBody, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Json not implemented")
}
func (UnimplementedSwaggerServiceServer) Yaml(context.Context, *empty.Empty) (*httpbody.HttpBody, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Yaml not implemented")
}
func (UnimplementedSwaggerServiceServer) UI(context.Context, *empty.Empty) (*httpbody.HttpBody, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UI not implemented")
}
func (UnimplementedSwaggerServiceServer) mustEmbedUnimplementedSwaggerServiceServer() {}

// UnsafeSwaggerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SwaggerServiceServer will
// result in compilation errors.
type UnsafeSwaggerServiceServer interface {
	mustEmbedUnimplementedSwaggerServiceServer()
}

func RegisterSwaggerServiceServer(s grpc.ServiceRegistrar, srv SwaggerServiceServer) {
	s.RegisterService(&_SwaggerService_serviceDesc, srv)
}

func _SwaggerService_Json_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SwaggerServiceServer).Json(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sole.api.v1.doc.swagger.SwaggerService/Json",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SwaggerServiceServer).Json(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _SwaggerService_Yaml_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SwaggerServiceServer).Yaml(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sole.api.v1.doc.swagger.SwaggerService/Yaml",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SwaggerServiceServer).Yaml(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _SwaggerService_UI_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SwaggerServiceServer).UI(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sole.api.v1.doc.swagger.SwaggerService/UI",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SwaggerServiceServer).UI(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _SwaggerService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "sole.api.v1.doc.swagger.SwaggerService",
	HandlerType: (*SwaggerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Json",
			Handler:    _SwaggerService_Json_Handler,
		},
		{
			MethodName: "Yaml",
			Handler:    _SwaggerService_Yaml_Handler,
		},
		{
			MethodName: "UI",
			Handler:    _SwaggerService_UI_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "swagger.proto",
}