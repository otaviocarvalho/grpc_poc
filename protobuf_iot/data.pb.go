// Code generated by protoc-gen-go.
// source: data.proto
// DO NOT EDIT!

/*
Package data is a generated protocol buffer package.

It is generated from these files:
	data.proto

It has these top-level messages:
	Measurement
*/
package data

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Measurement struct {
	Id    int64   `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Value float64 `protobuf:"fixed64,2,opt,name=value" json:"value,omitempty"`
}

func (m *Measurement) Reset()                    { *m = Measurement{} }
func (m *Measurement) String() string            { return proto.CompactTextString(m) }
func (*Measurement) ProtoMessage()               {}
func (*Measurement) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Measurement) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Measurement) GetValue() float64 {
	if m != nil {
		return m.Value
	}
	return 0
}

func init() {
	proto.RegisterType((*Measurement)(nil), "data.Measurement")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Data service

type DataClient interface {
	SendMeasurement(ctx context.Context, in *Measurement, opts ...grpc.CallOption) (*Measurement, error)
}

type dataClient struct {
	cc *grpc.ClientConn
}

func NewDataClient(cc *grpc.ClientConn) DataClient {
	return &dataClient{cc}
}

func (c *dataClient) SendMeasurement(ctx context.Context, in *Measurement, opts ...grpc.CallOption) (*Measurement, error) {
	out := new(Measurement)
	err := grpc.Invoke(ctx, "/data.Data/SendMeasurement", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Data service

type DataServer interface {
	SendMeasurement(context.Context, *Measurement) (*Measurement, error)
}

func RegisterDataServer(s *grpc.Server, srv DataServer) {
	s.RegisterService(&_Data_serviceDesc, srv)
}

func _Data_SendMeasurement_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Measurement)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataServer).SendMeasurement(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/data.Data/SendMeasurement",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataServer).SendMeasurement(ctx, req.(*Measurement))
	}
	return interceptor(ctx, in, info, handler)
}

var _Data_serviceDesc = grpc.ServiceDesc{
	ServiceName: "data.Data",
	HandlerType: (*DataServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendMeasurement",
			Handler:    _Data_SendMeasurement_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "data.proto",
}

func init() { proto.RegisterFile("data.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 123 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4a, 0x49, 0x2c, 0x49,
	0xd4, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x01, 0xb1, 0x95, 0x8c, 0xb9, 0xb8, 0x7d, 0x53,
	0x13, 0x8b, 0x4b, 0x8b, 0x52, 0x73, 0x53, 0xf3, 0x4a, 0x84, 0xf8, 0xb8, 0x98, 0x32, 0x53, 0x24,
	0x18, 0x15, 0x18, 0x35, 0x98, 0x83, 0x98, 0x32, 0x53, 0x84, 0x44, 0xb8, 0x58, 0xcb, 0x12, 0x73,
	0x4a, 0x53, 0x25, 0x98, 0x14, 0x18, 0x35, 0x18, 0x83, 0x20, 0x1c, 0x23, 0x47, 0x2e, 0x16, 0x97,
	0xc4, 0x92, 0x44, 0x21, 0x4b, 0x2e, 0xfe, 0xe0, 0xd4, 0xbc, 0x14, 0x64, 0x03, 0x04, 0xf5, 0xc0,
	0x56, 0x20, 0x09, 0x49, 0x61, 0x0a, 0x29, 0x31, 0x24, 0xb1, 0x81, 0x1d, 0x61, 0x0c, 0x08, 0x00,
	0x00, 0xff, 0xff, 0x25, 0xaa, 0xdc, 0x3d, 0x92, 0x00, 0x00, 0x00,
}
