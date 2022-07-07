// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: proxy.proto

package btrace_proxy

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

// ExporterClient is the client API for Exporter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ExporterClient interface {
	// used by both collector and agent
	Send(ctx context.Context, in *Span, opts ...grpc.CallOption) (*Response, error)
	// for the collector endpoint
	Stream(ctx context.Context, opts ...grpc.CallOption) (Exporter_StreamClient, error)
}

type exporterClient struct {
	cc grpc.ClientConnInterface
}

func NewExporterClient(cc grpc.ClientConnInterface) ExporterClient {
	return &exporterClient{cc}
}

func (c *exporterClient) Send(ctx context.Context, in *Span, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/Exporter/Send", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exporterClient) Stream(ctx context.Context, opts ...grpc.CallOption) (Exporter_StreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &Exporter_ServiceDesc.Streams[0], "/Exporter/Stream", opts...)
	if err != nil {
		return nil, err
	}
	x := &exporterStreamClient{stream}
	return x, nil
}

type Exporter_StreamClient interface {
	Send(*Span) error
	CloseAndRecv() (*Response, error)
	grpc.ClientStream
}

type exporterStreamClient struct {
	grpc.ClientStream
}

func (x *exporterStreamClient) Send(m *Span) error {
	return x.ClientStream.SendMsg(m)
}

func (x *exporterStreamClient) CloseAndRecv() (*Response, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ExporterServer is the server API for Exporter service.
// All implementations must embed UnimplementedExporterServer
// for forward compatibility
type ExporterServer interface {
	// used by both collector and agent
	Send(context.Context, *Span) (*Response, error)
	// for the collector endpoint
	Stream(Exporter_StreamServer) error
	mustEmbedUnimplementedExporterServer()
}

// UnimplementedExporterServer must be embedded to have forward compatible implementations.
type UnimplementedExporterServer struct {
}

func (UnimplementedExporterServer) Send(context.Context, *Span) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Send not implemented")
}
func (UnimplementedExporterServer) Stream(Exporter_StreamServer) error {
	return status.Errorf(codes.Unimplemented, "method Stream not implemented")
}
func (UnimplementedExporterServer) mustEmbedUnimplementedExporterServer() {}

// UnsafeExporterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ExporterServer will
// result in compilation errors.
type UnsafeExporterServer interface {
	mustEmbedUnimplementedExporterServer()
}

func RegisterExporterServer(s grpc.ServiceRegistrar, srv ExporterServer) {
	s.RegisterService(&Exporter_ServiceDesc, srv)
}

func _Exporter_Send_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Span)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExporterServer).Send(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Exporter/Send",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExporterServer).Send(ctx, req.(*Span))
	}
	return interceptor(ctx, in, info, handler)
}

func _Exporter_Stream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ExporterServer).Stream(&exporterStreamServer{stream})
}

type Exporter_StreamServer interface {
	SendAndClose(*Response) error
	Recv() (*Span, error)
	grpc.ServerStream
}

type exporterStreamServer struct {
	grpc.ServerStream
}

func (x *exporterStreamServer) SendAndClose(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *exporterStreamServer) Recv() (*Span, error) {
	m := new(Span)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Exporter_ServiceDesc is the grpc.ServiceDesc for Exporter service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Exporter_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Exporter",
	HandlerType: (*ExporterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Send",
			Handler:    _Exporter_Send_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Stream",
			Handler:       _Exporter_Stream_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "proxy.proto",
}
