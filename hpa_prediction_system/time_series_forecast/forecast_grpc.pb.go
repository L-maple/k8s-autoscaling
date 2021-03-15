// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package time_series_forecast

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// ForecastServiceClient is the client API for ForecastService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ForecastServiceClient interface {
	GetForeCastValue(ctx context.Context, in *ForecastRequest, opts ...grpc.CallOption) (*ForecastResponse, error)
}

type forecastServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewForecastServiceClient(cc grpc.ClientConnInterface) ForecastServiceClient {
	return &forecastServiceClient{cc}
}

func (c *forecastServiceClient) GetForeCastValue(ctx context.Context, in *ForecastRequest, opts ...grpc.CallOption) (*ForecastResponse, error) {
	out := new(ForecastResponse)
	err := c.cc.Invoke(ctx, "/time_series_forecast.ForecastService/GetForeCastValue", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ForecastServiceServer is the server API for ForecastService service.
// All implementations must embed UnimplementedForecastServiceServer
// for forward compatibility
type ForecastServiceServer interface {
	GetForeCastValue(context.Context, *ForecastRequest) (*ForecastResponse, error)
	mustEmbedUnimplementedForecastServiceServer()
}

// UnimplementedForecastServiceServer must be embedded to have forward compatible implementations.
type UnimplementedForecastServiceServer struct {
}

func (UnimplementedForecastServiceServer) GetForeCastValue(context.Context, *ForecastRequest) (*ForecastResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetForeCastValue not implemented")
}
func (UnimplementedForecastServiceServer) mustEmbedUnimplementedForecastServiceServer() {}

// UnsafeForecastServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ForecastServiceServer will
// result in compilation errors.
type UnsafeForecastServiceServer interface {
	mustEmbedUnimplementedForecastServiceServer()
}

func RegisterForecastServiceServer(s grpc.ServiceRegistrar, srv ForecastServiceServer) {
	s.RegisterService(&_ForecastService_serviceDesc, srv)
}

func _ForecastService_GetForeCastValue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ForecastRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ForecastServiceServer).GetForeCastValue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/time_series_forecast.ForecastService/GetForeCastValue",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ForecastServiceServer).GetForeCastValue(ctx, req.(*ForecastRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ForecastService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "time_series_forecast.ForecastService",
	HandlerType: (*ForecastServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetForeCastValue",
			Handler:    _ForecastService_GetForeCastValue_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "time_series_forecast/forecast.proto",
}
