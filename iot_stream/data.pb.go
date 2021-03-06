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
	Id    int64 `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Value int64 `protobuf:"varint,2,opt,name=value" json:"value,omitempty"`
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

func (m *Measurement) GetValue() int64 {
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

// Client API for MeterCommunicator service

type MeterCommunicatorClient interface {
	SimpleRPC(ctx context.Context, opts ...grpc.CallOption) (MeterCommunicator_SimpleRPCClient, error)
}

type meterCommunicatorClient struct {
	cc *grpc.ClientConn
}

func NewMeterCommunicatorClient(cc *grpc.ClientConn) MeterCommunicatorClient {
	return &meterCommunicatorClient{cc}
}

func (c *meterCommunicatorClient) SimpleRPC(ctx context.Context, opts ...grpc.CallOption) (MeterCommunicator_SimpleRPCClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_MeterCommunicator_serviceDesc.Streams[0], c.cc, "/data.MeterCommunicator/SimpleRPC", opts...)
	if err != nil {
		return nil, err
	}
	x := &meterCommunicatorSimpleRPCClient{stream}
	return x, nil
}

type MeterCommunicator_SimpleRPCClient interface {
	Send(*Measurement) error
	Recv() (*Measurement, error)
	grpc.ClientStream
}

type meterCommunicatorSimpleRPCClient struct {
	grpc.ClientStream
}

func (x *meterCommunicatorSimpleRPCClient) Send(m *Measurement) error {
	return x.ClientStream.SendMsg(m)
}

func (x *meterCommunicatorSimpleRPCClient) Recv() (*Measurement, error) {
	m := new(Measurement)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for MeterCommunicator service

type MeterCommunicatorServer interface {
	SimpleRPC(MeterCommunicator_SimpleRPCServer) error
}

func RegisterMeterCommunicatorServer(s *grpc.Server, srv MeterCommunicatorServer) {
	s.RegisterService(&_MeterCommunicator_serviceDesc, srv)
}

func _MeterCommunicator_SimpleRPC_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MeterCommunicatorServer).SimpleRPC(&meterCommunicatorSimpleRPCServer{stream})
}

type MeterCommunicator_SimpleRPCServer interface {
	Send(*Measurement) error
	Recv() (*Measurement, error)
	grpc.ServerStream
}

type meterCommunicatorSimpleRPCServer struct {
	grpc.ServerStream
}

func (x *meterCommunicatorSimpleRPCServer) Send(m *Measurement) error {
	return x.ServerStream.SendMsg(m)
}

func (x *meterCommunicatorSimpleRPCServer) Recv() (*Measurement, error) {
	m := new(Measurement)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _MeterCommunicator_serviceDesc = grpc.ServiceDesc{
	ServiceName: "data.MeterCommunicator",
	HandlerType: (*MeterCommunicatorServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SimpleRPC",
			Handler:       _MeterCommunicator_SimpleRPC_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "data.proto",
}

func init() { proto.RegisterFile("data.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 141 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4a, 0x49, 0x2c, 0x49,
	0xd4, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x01, 0xb1, 0x95, 0x8c, 0xb9, 0xb8, 0x7d, 0x53,
	0x13, 0x8b, 0x4b, 0x8b, 0x52, 0x73, 0x53, 0xf3, 0x4a, 0x84, 0xf8, 0xb8, 0x98, 0x32, 0x53, 0x24,
	0x18, 0x15, 0x18, 0x35, 0x98, 0x83, 0x98, 0x32, 0x53, 0x84, 0x44, 0xb8, 0x58, 0xcb, 0x12, 0x73,
	0x4a, 0x53, 0x25, 0x98, 0xc0, 0x42, 0x10, 0x8e, 0x91, 0x0f, 0x97, 0xa0, 0x6f, 0x6a, 0x49, 0x6a,
	0x91, 0x73, 0x7e, 0x6e, 0x6e, 0x69, 0x5e, 0x66, 0x72, 0x62, 0x49, 0x7e, 0x91, 0x90, 0x39, 0x17,
	0x67, 0x70, 0x66, 0x6e, 0x41, 0x4e, 0x6a, 0x50, 0x80, 0xb3, 0x90, 0xa0, 0x1e, 0xd8, 0x26, 0x24,
	0xa3, 0xa5, 0x30, 0x85, 0x94, 0x18, 0x34, 0x18, 0x0d, 0x18, 0x93, 0xd8, 0xc0, 0xee, 0x31, 0x06,
	0x04, 0x00, 0x00, 0xff, 0xff, 0x79, 0xbf, 0xb7, 0x7b, 0x9d, 0x00, 0x00, 0x00,
}
