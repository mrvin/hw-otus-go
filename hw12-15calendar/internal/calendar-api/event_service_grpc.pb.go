// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: event_service.proto

package calendarapi

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

const (
	EventService_CreateEvent_FullMethodName       = "/event.EventService/CreateEvent"
	EventService_GetEventByID_FullMethodName      = "/event.EventService/GetEventByID"
	EventService_UpdateEvent_FullMethodName       = "/event.EventService/UpdateEvent"
	EventService_DeleteEvent_FullMethodName       = "/event.EventService/DeleteEvent"
	EventService_ListEventsForUser_FullMethodName = "/event.EventService/ListEventsForUser"
)

// EventServiceClient is the client API for EventService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EventServiceClient interface {
	CreateEvent(ctx context.Context, in *CreateEventRequest, opts ...grpc.CallOption) (*CreateEventResponse, error)
	GetEventByID(ctx context.Context, in *GetEventByIDRequest, opts ...grpc.CallOption) (*EventResponse, error)
	UpdateEvent(ctx context.Context, in *UpdateEventRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	DeleteEvent(ctx context.Context, in *DeleteEventRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	ListEventsForUser(ctx context.Context, in *ListEventsForUserRequest, opts ...grpc.CallOption) (*ListEventsResponse, error)
}

type eventServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEventServiceClient(cc grpc.ClientConnInterface) EventServiceClient {
	return &eventServiceClient{cc}
}

func (c *eventServiceClient) CreateEvent(ctx context.Context, in *CreateEventRequest, opts ...grpc.CallOption) (*CreateEventResponse, error) {
	out := new(CreateEventResponse)
	err := c.cc.Invoke(ctx, EventService_CreateEvent_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) GetEventByID(ctx context.Context, in *GetEventByIDRequest, opts ...grpc.CallOption) (*EventResponse, error) {
	out := new(EventResponse)
	err := c.cc.Invoke(ctx, EventService_GetEventByID_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) UpdateEvent(ctx context.Context, in *UpdateEventRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, EventService_UpdateEvent_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) DeleteEvent(ctx context.Context, in *DeleteEventRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, EventService_DeleteEvent_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) ListEventsForUser(ctx context.Context, in *ListEventsForUserRequest, opts ...grpc.CallOption) (*ListEventsResponse, error) {
	out := new(ListEventsResponse)
	err := c.cc.Invoke(ctx, EventService_ListEventsForUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EventServiceServer is the server API for EventService service.
// All implementations should embed UnimplementedEventServiceServer
// for forward compatibility
type EventServiceServer interface {
	CreateEvent(context.Context, *CreateEventRequest) (*CreateEventResponse, error)
	GetEventByID(context.Context, *GetEventByIDRequest) (*EventResponse, error)
	UpdateEvent(context.Context, *UpdateEventRequest) (*emptypb.Empty, error)
	DeleteEvent(context.Context, *DeleteEventRequest) (*emptypb.Empty, error)
	ListEventsForUser(context.Context, *ListEventsForUserRequest) (*ListEventsResponse, error)
}

// UnimplementedEventServiceServer should be embedded to have forward compatible implementations.
type UnimplementedEventServiceServer struct {
}

func (UnimplementedEventServiceServer) CreateEvent(context.Context, *CreateEventRequest) (*CreateEventResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateEvent not implemented")
}
func (UnimplementedEventServiceServer) GetEventByID(context.Context, *GetEventByIDRequest) (*EventResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEventByID not implemented")
}
func (UnimplementedEventServiceServer) UpdateEvent(context.Context, *UpdateEventRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateEvent not implemented")
}
func (UnimplementedEventServiceServer) DeleteEvent(context.Context, *DeleteEventRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteEvent not implemented")
}
func (UnimplementedEventServiceServer) ListEventsForUser(context.Context, *ListEventsForUserRequest) (*ListEventsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEventsForUser not implemented")
}

// UnsafeEventServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EventServiceServer will
// result in compilation errors.
type UnsafeEventServiceServer interface {
	mustEmbedUnimplementedEventServiceServer()
}

func RegisterEventServiceServer(s grpc.ServiceRegistrar, srv EventServiceServer) {
	s.RegisterService(&EventService_ServiceDesc, srv)
}

func _EventService_CreateEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).CreateEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventService_CreateEvent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).CreateEvent(ctx, req.(*CreateEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_GetEventByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEventByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).GetEventByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventService_GetEventByID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).GetEventByID(ctx, req.(*GetEventByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_UpdateEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).UpdateEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventService_UpdateEvent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).UpdateEvent(ctx, req.(*UpdateEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_DeleteEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).DeleteEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventService_DeleteEvent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).DeleteEvent(ctx, req.(*DeleteEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_ListEventsForUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListEventsForUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).ListEventsForUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventService_ListEventsForUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).ListEventsForUser(ctx, req.(*ListEventsForUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// EventService_ServiceDesc is the grpc.ServiceDesc for EventService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EventService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "event.EventService",
	HandlerType: (*EventServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateEvent",
			Handler:    _EventService_CreateEvent_Handler,
		},
		{
			MethodName: "GetEventByID",
			Handler:    _EventService_GetEventByID_Handler,
		},
		{
			MethodName: "UpdateEvent",
			Handler:    _EventService_UpdateEvent_Handler,
		},
		{
			MethodName: "DeleteEvent",
			Handler:    _EventService_DeleteEvent_Handler,
		},
		{
			MethodName: "ListEventsForUser",
			Handler:    _EventService_ListEventsForUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "event_service.proto",
}
