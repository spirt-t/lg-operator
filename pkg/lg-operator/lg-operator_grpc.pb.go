// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.23.4
// source: lg-operator/lg-operator.proto

package lg_operator

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// LoadGeneratorOperatorServiceClient is the client API for LoadGeneratorOperatorService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LoadGeneratorOperatorServiceClient interface {
	// Debug entrypoint.
	Hello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloResponse, error)
	// Create pod, service and ingress of load-generator according to the passed parameters.
	CreateGenerators(ctx context.Context, in *CreateGeneratorsRequest, opts ...grpc.CallOption) (*CreateGeneratorsResponse, error)
	// Delete pod, service and ingress by load-generator name.
	DeleteGenerators(ctx context.Context, in *DeleteGeneratorsRequest, opts ...grpc.CallOption) (*DeleteGeneratorsResponse, error)
	// Get list of all load-generators in cluster.
	GeneratorsList(ctx context.Context, in *GeneratorsListRequest, opts ...grpc.CallOption) (*GeneratorsListResponse, error)
	// Delete all pods, services and ingresses of generators. Use carefully!
	ClearAll(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type loadGeneratorOperatorServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLoadGeneratorOperatorServiceClient(cc grpc.ClientConnInterface) LoadGeneratorOperatorServiceClient {
	return &loadGeneratorOperatorServiceClient{cc}
}

func (c *loadGeneratorOperatorServiceClient) Hello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloResponse, error) {
	out := new(HelloResponse)
	err := c.cc.Invoke(ctx, "/lg_operator.LoadGeneratorOperatorService/Hello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *loadGeneratorOperatorServiceClient) CreateGenerators(ctx context.Context, in *CreateGeneratorsRequest, opts ...grpc.CallOption) (*CreateGeneratorsResponse, error) {
	out := new(CreateGeneratorsResponse)
	err := c.cc.Invoke(ctx, "/lg_operator.LoadGeneratorOperatorService/CreateGenerators", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *loadGeneratorOperatorServiceClient) DeleteGenerators(ctx context.Context, in *DeleteGeneratorsRequest, opts ...grpc.CallOption) (*DeleteGeneratorsResponse, error) {
	out := new(DeleteGeneratorsResponse)
	err := c.cc.Invoke(ctx, "/lg_operator.LoadGeneratorOperatorService/DeleteGenerators", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *loadGeneratorOperatorServiceClient) GeneratorsList(ctx context.Context, in *GeneratorsListRequest, opts ...grpc.CallOption) (*GeneratorsListResponse, error) {
	out := new(GeneratorsListResponse)
	err := c.cc.Invoke(ctx, "/lg_operator.LoadGeneratorOperatorService/GeneratorsList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *loadGeneratorOperatorServiceClient) ClearAll(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/lg_operator.LoadGeneratorOperatorService/ClearAll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LoadGeneratorOperatorServiceServer is the server API for LoadGeneratorOperatorService service.
// All implementations must embed UnimplementedLoadGeneratorOperatorServiceServer
// for forward compatibility
type LoadGeneratorOperatorServiceServer interface {
	// Debug entrypoint.
	Hello(context.Context, *HelloRequest) (*HelloResponse, error)
	// Create pod, service and ingress of load-generator according to the passed parameters.
	CreateGenerators(context.Context, *CreateGeneratorsRequest) (*CreateGeneratorsResponse, error)
	// Delete pod, service and ingress by load-generator name.
	DeleteGenerators(context.Context, *DeleteGeneratorsRequest) (*DeleteGeneratorsResponse, error)
	// Get list of all load-generators in cluster.
	GeneratorsList(context.Context, *GeneratorsListRequest) (*GeneratorsListResponse, error)
	// Delete all pods, services and ingresses of generators. Use carefully!
	ClearAll(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	mustEmbedUnimplementedLoadGeneratorOperatorServiceServer()
}

// UnimplementedLoadGeneratorOperatorServiceServer must be embedded to have forward compatible implementations.
type UnimplementedLoadGeneratorOperatorServiceServer struct {
}

func (UnimplementedLoadGeneratorOperatorServiceServer) Hello(context.Context, *HelloRequest) (*HelloResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Hello not implemented")
}
func (UnimplementedLoadGeneratorOperatorServiceServer) CreateGenerators(context.Context, *CreateGeneratorsRequest) (*CreateGeneratorsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateGenerators not implemented")
}
func (UnimplementedLoadGeneratorOperatorServiceServer) DeleteGenerators(context.Context, *DeleteGeneratorsRequest) (*DeleteGeneratorsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteGenerators not implemented")
}
func (UnimplementedLoadGeneratorOperatorServiceServer) GeneratorsList(context.Context, *GeneratorsListRequest) (*GeneratorsListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GeneratorsList not implemented")
}
func (UnimplementedLoadGeneratorOperatorServiceServer) ClearAll(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClearAll not implemented")
}
func (UnimplementedLoadGeneratorOperatorServiceServer) mustEmbedUnimplementedLoadGeneratorOperatorServiceServer() {
}

// UnsafeLoadGeneratorOperatorServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LoadGeneratorOperatorServiceServer will
// result in compilation errors.
type UnsafeLoadGeneratorOperatorServiceServer interface {
	mustEmbedUnimplementedLoadGeneratorOperatorServiceServer()
}

func RegisterLoadGeneratorOperatorServiceServer(s grpc.ServiceRegistrar, srv LoadGeneratorOperatorServiceServer) {
	s.RegisterService(&LoadGeneratorOperatorService_ServiceDesc, srv)
}

func _LoadGeneratorOperatorService_Hello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoadGeneratorOperatorServiceServer).Hello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lg_operator.LoadGeneratorOperatorService/Hello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoadGeneratorOperatorServiceServer).Hello(ctx, req.(*HelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LoadGeneratorOperatorService_CreateGenerators_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateGeneratorsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoadGeneratorOperatorServiceServer).CreateGenerators(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lg_operator.LoadGeneratorOperatorService/CreateGenerators",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoadGeneratorOperatorServiceServer).CreateGenerators(ctx, req.(*CreateGeneratorsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LoadGeneratorOperatorService_DeleteGenerators_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteGeneratorsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoadGeneratorOperatorServiceServer).DeleteGenerators(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lg_operator.LoadGeneratorOperatorService/DeleteGenerators",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoadGeneratorOperatorServiceServer).DeleteGenerators(ctx, req.(*DeleteGeneratorsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LoadGeneratorOperatorService_GeneratorsList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GeneratorsListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoadGeneratorOperatorServiceServer).GeneratorsList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lg_operator.LoadGeneratorOperatorService/GeneratorsList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoadGeneratorOperatorServiceServer).GeneratorsList(ctx, req.(*GeneratorsListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LoadGeneratorOperatorService_ClearAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoadGeneratorOperatorServiceServer).ClearAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lg_operator.LoadGeneratorOperatorService/ClearAll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoadGeneratorOperatorServiceServer).ClearAll(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// LoadGeneratorOperatorService_ServiceDesc is the grpc.ServiceDesc for LoadGeneratorOperatorService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LoadGeneratorOperatorService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "lg_operator.LoadGeneratorOperatorService",
	HandlerType: (*LoadGeneratorOperatorServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Hello",
			Handler:    _LoadGeneratorOperatorService_Hello_Handler,
		},
		{
			MethodName: "CreateGenerators",
			Handler:    _LoadGeneratorOperatorService_CreateGenerators_Handler,
		},
		{
			MethodName: "DeleteGenerators",
			Handler:    _LoadGeneratorOperatorService_DeleteGenerators_Handler,
		},
		{
			MethodName: "GeneratorsList",
			Handler:    _LoadGeneratorOperatorService_GeneratorsList_Handler,
		},
		{
			MethodName: "ClearAll",
			Handler:    _LoadGeneratorOperatorService_ClearAll_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "lg-operator/lg-operator.proto",
}
