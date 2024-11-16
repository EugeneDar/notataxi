// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: toll_roads.proto

package toll_roads

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
	TollRoadsService_GetTollRoads_FullMethodName = "/toll_roads.TollRoadsService/GetTollRoads"
)

// TollRoadsServiceClient is the client API for TollRoadsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TollRoadsServiceClient interface {
	GetTollRoads(ctx context.Context, in *TollRoadsRequest, opts ...grpc.CallOption) (*TollRoadsResponse, error)
}

type tollRoadsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTollRoadsServiceClient(cc grpc.ClientConnInterface) TollRoadsServiceClient {
	return &tollRoadsServiceClient{cc}
}

func (c *tollRoadsServiceClient) GetTollRoads(ctx context.Context, in *TollRoadsRequest, opts ...grpc.CallOption) (*TollRoadsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TollRoadsResponse)
	err := c.cc.Invoke(ctx, TollRoadsService_GetTollRoads_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TollRoadsServiceServer is the server API for TollRoadsService service.
// All implementations must embed UnimplementedTollRoadsServiceServer
// for forward compatibility.
type TollRoadsServiceServer interface {
	GetTollRoads(context.Context, *TollRoadsRequest) (*TollRoadsResponse, error)
	mustEmbedUnimplementedTollRoadsServiceServer()
}

// UnimplementedTollRoadsServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedTollRoadsServiceServer struct{}

func (UnimplementedTollRoadsServiceServer) GetTollRoads(context.Context, *TollRoadsRequest) (*TollRoadsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTollRoads not implemented")
}
func (UnimplementedTollRoadsServiceServer) mustEmbedUnimplementedTollRoadsServiceServer() {}
func (UnimplementedTollRoadsServiceServer) testEmbeddedByValue()                          {}

// UnsafeTollRoadsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TollRoadsServiceServer will
// result in compilation errors.
type UnsafeTollRoadsServiceServer interface {
	mustEmbedUnimplementedTollRoadsServiceServer()
}

func RegisterTollRoadsServiceServer(s grpc.ServiceRegistrar, srv TollRoadsServiceServer) {
	// If the following call pancis, it indicates UnimplementedTollRoadsServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&TollRoadsService_ServiceDesc, srv)
}

func _TollRoadsService_GetTollRoads_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TollRoadsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TollRoadsServiceServer).GetTollRoads(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TollRoadsService_GetTollRoads_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TollRoadsServiceServer).GetTollRoads(ctx, req.(*TollRoadsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TollRoadsService_ServiceDesc is the grpc.ServiceDesc for TollRoadsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TollRoadsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "toll_roads.TollRoadsService",
	HandlerType: (*TollRoadsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetTollRoads",
			Handler:    _TollRoadsService_GetTollRoads_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "toll_roads.proto",
}
