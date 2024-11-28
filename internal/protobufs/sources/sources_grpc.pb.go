// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.26.1
// source: sources.proto

package sources

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

// SourcesServiceClient is the client API for SourcesService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SourcesServiceClient interface {
	GetOrderInfo(ctx context.Context, in *SourcesRequest, opts ...grpc.CallOption) (*SourcesResponse, error)
}

type sourcesServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSourcesServiceClient(cc grpc.ClientConnInterface) SourcesServiceClient {
	return &sourcesServiceClient{cc}
}

func (c *sourcesServiceClient) GetOrderInfo(ctx context.Context, in *SourcesRequest, opts ...grpc.CallOption) (*SourcesResponse, error) {
	out := new(SourcesResponse)
	err := c.cc.Invoke(ctx, "/origins.SourcesService/GetOrderInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SourcesServiceServer is the server API for SourcesService service.
// All implementations must embed UnimplementedSourcesServiceServer
// for forward compatibility
type SourcesServiceServer interface {
	GetOrderInfo(context.Context, *SourcesRequest) (*SourcesResponse, error)
	mustEmbedUnimplementedSourcesServiceServer()
}

// UnimplementedSourcesServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSourcesServiceServer struct {
}

func (UnimplementedSourcesServiceServer) GetOrderInfo(context.Context, *SourcesRequest) (*SourcesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOrderInfo not implemented")
}
func (UnimplementedSourcesServiceServer) mustEmbedUnimplementedSourcesServiceServer() {}

// UnsafeSourcesServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SourcesServiceServer will
// result in compilation errors.
type UnsafeSourcesServiceServer interface {
	mustEmbedUnimplementedSourcesServiceServer()
}

func RegisterSourcesServiceServer(s grpc.ServiceRegistrar, srv SourcesServiceServer) {
	s.RegisterService(&SourcesService_ServiceDesc, srv)
}

func _SourcesService_GetOrderInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SourcesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SourcesServiceServer).GetOrderInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/origins.SourcesService/GetOrderInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SourcesServiceServer).GetOrderInfo(ctx, req.(*SourcesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SourcesService_ServiceDesc is the grpc.ServiceDesc for SourcesService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SourcesService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "origins.SourcesService",
	HandlerType: (*SourcesServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetOrderInfo",
			Handler:    _SourcesService_GetOrderInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sources.proto",
}
