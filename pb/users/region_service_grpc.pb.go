// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: users/region_service.proto

package users

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

const (
	RegionService_Create_FullMethodName = "/wiradata.users.RegionService/Create"
	RegionService_Update_FullMethodName = "/wiradata.users.RegionService/Update"
	RegionService_View_FullMethodName   = "/wiradata.users.RegionService/View"
	RegionService_Delete_FullMethodName = "/wiradata.users.RegionService/Delete"
	RegionService_List_FullMethodName   = "/wiradata.users.RegionService/List"
)

// RegionServiceClient is the client API for RegionService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RegionServiceClient interface {
	Create(ctx context.Context, in *Region, opts ...grpc.CallOption) (*Region, error)
	Update(ctx context.Context, in *Region, opts ...grpc.CallOption) (*Region, error)
	View(ctx context.Context, in *Id, opts ...grpc.CallOption) (*Region, error)
	Delete(ctx context.Context, in *Id, opts ...grpc.CallOption) (*MyBoolean, error)
	List(ctx context.Context, in *ListRegionRequest, opts ...grpc.CallOption) (RegionService_ListClient, error)
}

type regionServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRegionServiceClient(cc grpc.ClientConnInterface) RegionServiceClient {
	return &regionServiceClient{cc}
}

func (c *regionServiceClient) Create(ctx context.Context, in *Region, opts ...grpc.CallOption) (*Region, error) {
	out := new(Region)
	err := c.cc.Invoke(ctx, RegionService_Create_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *regionServiceClient) Update(ctx context.Context, in *Region, opts ...grpc.CallOption) (*Region, error) {
	out := new(Region)
	err := c.cc.Invoke(ctx, RegionService_Update_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *regionServiceClient) View(ctx context.Context, in *Id, opts ...grpc.CallOption) (*Region, error) {
	out := new(Region)
	err := c.cc.Invoke(ctx, RegionService_View_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *regionServiceClient) Delete(ctx context.Context, in *Id, opts ...grpc.CallOption) (*MyBoolean, error) {
	out := new(MyBoolean)
	err := c.cc.Invoke(ctx, RegionService_Delete_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *regionServiceClient) List(ctx context.Context, in *ListRegionRequest, opts ...grpc.CallOption) (RegionService_ListClient, error) {
	stream, err := c.cc.NewStream(ctx, &RegionService_ServiceDesc.Streams[0], RegionService_List_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &regionServiceListClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type RegionService_ListClient interface {
	Recv() (*ListRegionResponse, error)
	grpc.ClientStream
}

type regionServiceListClient struct {
	grpc.ClientStream
}

func (x *regionServiceListClient) Recv() (*ListRegionResponse, error) {
	m := new(ListRegionResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// RegionServiceServer is the server API for RegionService service.
// All implementations must embed UnimplementedRegionServiceServer
// for forward compatibility
type RegionServiceServer interface {
	Create(context.Context, *Region) (*Region, error)
	Update(context.Context, *Region) (*Region, error)
	View(context.Context, *Id) (*Region, error)
	Delete(context.Context, *Id) (*MyBoolean, error)
	List(*ListRegionRequest, RegionService_ListServer) error
	mustEmbedUnimplementedRegionServiceServer()
}

// UnimplementedRegionServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRegionServiceServer struct {
}

func (UnimplementedRegionServiceServer) Create(context.Context, *Region) (*Region, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedRegionServiceServer) Update(context.Context, *Region) (*Region, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedRegionServiceServer) View(context.Context, *Id) (*Region, error) {
	return nil, status.Errorf(codes.Unimplemented, "method View not implemented")
}
func (UnimplementedRegionServiceServer) Delete(context.Context, *Id) (*MyBoolean, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedRegionServiceServer) List(*ListRegionRequest, RegionService_ListServer) error {
	return status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedRegionServiceServer) mustEmbedUnimplementedRegionServiceServer() {}

// UnsafeRegionServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RegionServiceServer will
// result in compilation errors.
type UnsafeRegionServiceServer interface {
	mustEmbedUnimplementedRegionServiceServer()
}

func RegisterRegionServiceServer(s grpc.ServiceRegistrar, srv RegionServiceServer) {
	s.RegisterService(&RegionService_ServiceDesc, srv)
}

func _RegionService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Region)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegionServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RegionService_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegionServiceServer).Create(ctx, req.(*Region))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegionService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Region)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegionServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RegionService_Update_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegionServiceServer).Update(ctx, req.(*Region))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegionService_View_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Id)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegionServiceServer).View(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RegionService_View_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegionServiceServer).View(ctx, req.(*Id))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegionService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Id)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegionServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RegionService_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegionServiceServer).Delete(ctx, req.(*Id))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegionService_List_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListRegionRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(RegionServiceServer).List(m, &regionServiceListServer{stream})
}

type RegionService_ListServer interface {
	Send(*ListRegionResponse) error
	grpc.ServerStream
}

type regionServiceListServer struct {
	grpc.ServerStream
}

func (x *regionServiceListServer) Send(m *ListRegionResponse) error {
	return x.ServerStream.SendMsg(m)
}

// RegionService_ServiceDesc is the grpc.ServiceDesc for RegionService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RegionService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "wiradata.users.RegionService",
	HandlerType: (*RegionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _RegionService_Create_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _RegionService_Update_Handler,
		},
		{
			MethodName: "View",
			Handler:    _RegionService_View_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _RegionService_Delete_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "List",
			Handler:       _RegionService_List_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "users/region_service.proto",
}
