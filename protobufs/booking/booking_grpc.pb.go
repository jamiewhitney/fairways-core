// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: protobufs/booking/booking.proto

package booking

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
	BookingService_GetBooking_FullMethodName           = "/booking.BookingService/GetBooking"
	BookingService_GetBookings_FullMethodName          = "/booking.BookingService/GetBookings"
	BookingService_CreateBooking_FullMethodName        = "/booking.BookingService/CreateBooking"
	BookingService_GetConfirmedBookings_FullMethodName = "/booking.BookingService/GetConfirmedBookings"
)

// BookingServiceClient is the client API for BookingService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BookingServiceClient interface {
	GetBooking(ctx context.Context, in *GetBookingRequest, opts ...grpc.CallOption) (*GetBookingResponse, error)
	GetBookings(ctx context.Context, in *GetBookingsRequest, opts ...grpc.CallOption) (*GetBookingsResponse, error)
	CreateBooking(ctx context.Context, in *CreateBookingRequest, opts ...grpc.CallOption) (*CreateBookingResponse, error)
	GetConfirmedBookings(ctx context.Context, in *GetConfirmedBookingsRequest, opts ...grpc.CallOption) (*GetConfirmedBookingResponse, error)
}

type bookingServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBookingServiceClient(cc grpc.ClientConnInterface) BookingServiceClient {
	return &bookingServiceClient{cc}
}

func (c *bookingServiceClient) GetBooking(ctx context.Context, in *GetBookingRequest, opts ...grpc.CallOption) (*GetBookingResponse, error) {
	out := new(GetBookingResponse)
	err := c.cc.Invoke(ctx, BookingService_GetBooking_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bookingServiceClient) GetBookings(ctx context.Context, in *GetBookingsRequest, opts ...grpc.CallOption) (*GetBookingsResponse, error) {
	out := new(GetBookingsResponse)
	err := c.cc.Invoke(ctx, BookingService_GetBookings_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bookingServiceClient) CreateBooking(ctx context.Context, in *CreateBookingRequest, opts ...grpc.CallOption) (*CreateBookingResponse, error) {
	out := new(CreateBookingResponse)
	err := c.cc.Invoke(ctx, BookingService_CreateBooking_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bookingServiceClient) GetConfirmedBookings(ctx context.Context, in *GetConfirmedBookingsRequest, opts ...grpc.CallOption) (*GetConfirmedBookingResponse, error) {
	out := new(GetConfirmedBookingResponse)
	err := c.cc.Invoke(ctx, BookingService_GetConfirmedBookings_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BookingServiceServer is the server API for BookingService service.
// All implementations must embed UnimplementedBookingServiceServer
// for forward compatibility
type BookingServiceServer interface {
	GetBooking(context.Context, *GetBookingRequest) (*GetBookingResponse, error)
	GetBookings(context.Context, *GetBookingsRequest) (*GetBookingsResponse, error)
	CreateBooking(context.Context, *CreateBookingRequest) (*CreateBookingResponse, error)
	GetConfirmedBookings(context.Context, *GetConfirmedBookingsRequest) (*GetConfirmedBookingResponse, error)
	mustEmbedUnimplementedBookingServiceServer()
}

// UnimplementedBookingServiceServer must be embedded to have forward compatible implementations.
type UnimplementedBookingServiceServer struct {
}

func (UnimplementedBookingServiceServer) GetBooking(context.Context, *GetBookingRequest) (*GetBookingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBooking not implemented")
}
func (UnimplementedBookingServiceServer) GetBookings(context.Context, *GetBookingsRequest) (*GetBookingsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBookings not implemented")
}
func (UnimplementedBookingServiceServer) CreateBooking(context.Context, *CreateBookingRequest) (*CreateBookingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateBooking not implemented")
}
func (UnimplementedBookingServiceServer) GetConfirmedBookings(context.Context, *GetConfirmedBookingsRequest) (*GetConfirmedBookingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConfirmedBookings not implemented")
}
func (UnimplementedBookingServiceServer) mustEmbedUnimplementedBookingServiceServer() {}

// UnsafeBookingServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BookingServiceServer will
// result in compilation errors.
type UnsafeBookingServiceServer interface {
	mustEmbedUnimplementedBookingServiceServer()
}

func RegisterBookingServiceServer(s grpc.ServiceRegistrar, srv BookingServiceServer) {
	s.RegisterService(&BookingService_ServiceDesc, srv)
}

func _BookingService_GetBooking_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBookingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookingServiceServer).GetBooking(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BookingService_GetBooking_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookingServiceServer).GetBooking(ctx, req.(*GetBookingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BookingService_GetBookings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBookingsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookingServiceServer).GetBookings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BookingService_GetBookings_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookingServiceServer).GetBookings(ctx, req.(*GetBookingsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BookingService_CreateBooking_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateBookingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookingServiceServer).CreateBooking(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BookingService_CreateBooking_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookingServiceServer).CreateBooking(ctx, req.(*CreateBookingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BookingService_GetConfirmedBookings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetConfirmedBookingsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookingServiceServer).GetConfirmedBookings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BookingService_GetConfirmedBookings_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookingServiceServer).GetConfirmedBookings(ctx, req.(*GetConfirmedBookingsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// BookingService_ServiceDesc is the grpc.ServiceDesc for BookingService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BookingService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "booking.BookingService",
	HandlerType: (*BookingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetBooking",
			Handler:    _BookingService_GetBooking_Handler,
		},
		{
			MethodName: "GetBookings",
			Handler:    _BookingService_GetBookings_Handler,
		},
		{
			MethodName: "CreateBooking",
			Handler:    _BookingService_CreateBooking_Handler,
		},
		{
			MethodName: "GetConfirmedBookings",
			Handler:    _BookingService_GetConfirmedBookings_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protobufs/booking/booking.proto",
}
