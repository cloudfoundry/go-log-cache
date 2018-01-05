// Code generated by protoc-gen-go. DO NOT EDIT.
// source: egress.proto

/*
Package logcache is a generated protocol buffer package.

It is generated from these files:
	egress.proto
	group_reader.proto
	ingress.proto

It has these top-level messages:
	ReadRequest
	ReadResponse
	AddToGroupRequest
	AddToGroupResponse
	RemoveFromGroupRequest
	RemoveFromGroupResponse
	GroupReadRequest
	GroupReadResponse
	GroupRequest
	GroupResponse
	SendRequest
	SendResponse
*/
package logcache

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import loggregator_v2 "code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
import _ "google.golang.org/genproto/googleapis/api/annotations"

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

type EnvelopeTypes int32

const (
	EnvelopeTypes_ANY     EnvelopeTypes = 0
	EnvelopeTypes_LOG     EnvelopeTypes = 1
	EnvelopeTypes_COUNTER EnvelopeTypes = 2
	EnvelopeTypes_GAUGE   EnvelopeTypes = 3
	EnvelopeTypes_TIMER   EnvelopeTypes = 4
	EnvelopeTypes_EVENT   EnvelopeTypes = 5
)

var EnvelopeTypes_name = map[int32]string{
	0: "ANY",
	1: "LOG",
	2: "COUNTER",
	3: "GAUGE",
	4: "TIMER",
	5: "EVENT",
}
var EnvelopeTypes_value = map[string]int32{
	"ANY":     0,
	"LOG":     1,
	"COUNTER": 2,
	"GAUGE":   3,
	"TIMER":   4,
	"EVENT":   5,
}

func (x EnvelopeTypes) String() string {
	return proto.EnumName(EnvelopeTypes_name, int32(x))
}
func (EnvelopeTypes) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type ReadRequest struct {
	SourceId     string        `protobuf:"bytes,1,opt,name=source_id,json=sourceId" json:"source_id,omitempty"`
	StartTime    int64         `protobuf:"varint,2,opt,name=start_time,json=startTime" json:"start_time,omitempty"`
	EndTime      int64         `protobuf:"varint,3,opt,name=end_time,json=endTime" json:"end_time,omitempty"`
	Limit        int64         `protobuf:"varint,4,opt,name=limit" json:"limit,omitempty"`
	EnvelopeType EnvelopeTypes `protobuf:"varint,5,opt,name=envelope_type,json=envelopeType,enum=logcache.EnvelopeTypes" json:"envelope_type,omitempty"`
}

func (m *ReadRequest) Reset()                    { *m = ReadRequest{} }
func (m *ReadRequest) String() string            { return proto.CompactTextString(m) }
func (*ReadRequest) ProtoMessage()               {}
func (*ReadRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ReadRequest) GetSourceId() string {
	if m != nil {
		return m.SourceId
	}
	return ""
}

func (m *ReadRequest) GetStartTime() int64 {
	if m != nil {
		return m.StartTime
	}
	return 0
}

func (m *ReadRequest) GetEndTime() int64 {
	if m != nil {
		return m.EndTime
	}
	return 0
}

func (m *ReadRequest) GetLimit() int64 {
	if m != nil {
		return m.Limit
	}
	return 0
}

func (m *ReadRequest) GetEnvelopeType() EnvelopeTypes {
	if m != nil {
		return m.EnvelopeType
	}
	return EnvelopeTypes_ANY
}

type ReadResponse struct {
	Envelopes *loggregator_v2.EnvelopeBatch `protobuf:"bytes,1,opt,name=envelopes" json:"envelopes,omitempty"`
}

func (m *ReadResponse) Reset()                    { *m = ReadResponse{} }
func (m *ReadResponse) String() string            { return proto.CompactTextString(m) }
func (*ReadResponse) ProtoMessage()               {}
func (*ReadResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ReadResponse) GetEnvelopes() *loggregator_v2.EnvelopeBatch {
	if m != nil {
		return m.Envelopes
	}
	return nil
}

func init() {
	proto.RegisterType((*ReadRequest)(nil), "logcache.ReadRequest")
	proto.RegisterType((*ReadResponse)(nil), "logcache.ReadResponse")
	proto.RegisterEnum("logcache.EnvelopeTypes", EnvelopeTypes_name, EnvelopeTypes_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Egress service

type EgressClient interface {
	Read(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*ReadResponse, error)
}

type egressClient struct {
	cc *grpc.ClientConn
}

func NewEgressClient(cc *grpc.ClientConn) EgressClient {
	return &egressClient{cc}
}

func (c *egressClient) Read(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*ReadResponse, error) {
	out := new(ReadResponse)
	err := grpc.Invoke(ctx, "/logcache.Egress/Read", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Egress service

type EgressServer interface {
	Read(context.Context, *ReadRequest) (*ReadResponse, error)
}

func RegisterEgressServer(s *grpc.Server, srv EgressServer) {
	s.RegisterService(&_Egress_serviceDesc, srv)
}

func _Egress_Read_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EgressServer).Read(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/logcache.Egress/Read",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EgressServer).Read(ctx, req.(*ReadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Egress_serviceDesc = grpc.ServiceDesc{
	ServiceName: "logcache.Egress",
	HandlerType: (*EgressServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Read",
			Handler:    _Egress_Read_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "egress.proto",
}

func init() { proto.RegisterFile("egress.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 380 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x91, 0xcf, 0xae, 0xd2, 0x40,
	0x14, 0xc6, 0x2d, 0xa5, 0x40, 0x0f, 0x60, 0xea, 0x04, 0xb1, 0x22, 0x24, 0x84, 0x15, 0x71, 0xd1,
	0xc6, 0xba, 0xd4, 0x0d, 0x9a, 0x86, 0x10, 0x15, 0x92, 0xb1, 0x98, 0xb8, 0x30, 0x64, 0x6c, 0x4f,
	0x4a, 0x93, 0xd2, 0xa9, 0x9d, 0xa1, 0x09, 0x31, 0x6e, 0x7c, 0x85, 0xfb, 0x32, 0xf7, 0x3d, 0xee,
	0x2b, 0xdc, 0x07, 0xb9, 0x69, 0xcb, 0xbf, 0x7b, 0x77, 0xe7, 0x7c, 0xbf, 0x99, 0x33, 0xdf, 0x99,
	0x0f, 0x3a, 0x18, 0x66, 0x28, 0x84, 0x95, 0x66, 0x5c, 0x72, 0xd2, 0x8a, 0x79, 0xe8, 0x33, 0x7f,
	0x8b, 0x83, 0x17, 0xb9, 0x63, 0x63, 0x92, 0x63, 0xcc, 0x53, 0xac, 0xe0, 0x60, 0x18, 0x72, 0x1e,
	0xc6, 0x68, 0xb3, 0x34, 0xb2, 0x59, 0x92, 0x70, 0xc9, 0x64, 0xc4, 0x93, 0xe3, 0xd5, 0xc9, 0xad,
	0x02, 0x6d, 0x8a, 0x2c, 0xa0, 0xf8, 0x67, 0x8f, 0x42, 0x92, 0x37, 0xa0, 0x0b, 0xbe, 0xcf, 0x7c,
	0xdc, 0x44, 0x81, 0xa9, 0x8c, 0x95, 0xa9, 0x4e, 0x5b, 0x95, 0xb0, 0x08, 0xc8, 0x08, 0x40, 0x48,
	0x96, 0xc9, 0x8d, 0x8c, 0x76, 0x68, 0xd6, 0xc6, 0xca, 0x54, 0xa5, 0x7a, 0xa9, 0x78, 0xd1, 0x0e,
	0xc9, 0x6b, 0x68, 0x61, 0x12, 0x54, 0x50, 0x2d, 0x61, 0x13, 0x93, 0xa0, 0x44, 0x3d, 0xd0, 0xe2,
	0x68, 0x17, 0x49, 0xb3, 0x5e, 0xea, 0x55, 0x43, 0x3e, 0x42, 0xf7, 0x64, 0x76, 0x23, 0x0f, 0x29,
	0x9a, 0xda, 0x58, 0x99, 0x3e, 0x77, 0x5e, 0x59, 0xa7, 0x7d, 0x2c, 0xf7, 0x88, 0xbd, 0x43, 0x8a,
	0x82, 0x76, 0xf0, 0xaa, 0x9d, 0x7c, 0x81, 0x4e, 0xe5, 0x5c, 0xa4, 0x3c, 0x11, 0x48, 0x3e, 0x80,
	0x7e, 0xe2, 0xa2, 0xb4, 0xde, 0x76, 0x46, 0xc5, 0xa4, 0x30, 0xc3, 0x90, 0x49, 0x9e, 0x59, 0xb9,
	0x73, 0x9e, 0xf7, 0x89, 0x49, 0x7f, 0x4b, 0x2f, 0xe7, 0xdf, 0xae, 0xa0, 0xfb, 0xe8, 0x2d, 0xd2,
	0x04, 0x75, 0xb6, 0xfc, 0x69, 0x3c, 0x2b, 0x8a, 0xaf, 0xab, 0xb9, 0xa1, 0x90, 0x36, 0x34, 0x3f,
	0xaf, 0xd6, 0x4b, 0xcf, 0xa5, 0x46, 0x8d, 0xe8, 0xa0, 0xcd, 0x67, 0xeb, 0xb9, 0x6b, 0xa8, 0x45,
	0xe9, 0x2d, 0xbe, 0xb9, 0xd4, 0xa8, 0x17, 0xa5, 0xfb, 0xc3, 0x5d, 0x7a, 0x86, 0xe6, 0xfc, 0x82,
	0x86, 0x5b, 0x66, 0x44, 0xbe, 0x43, 0xbd, 0xf0, 0x49, 0x5e, 0x5e, 0xd6, 0xba, 0xfa, 0xf1, 0x41,
	0xff, 0xa9, 0x5c, 0xad, 0x33, 0x19, 0xfe, 0xbf, 0xbb, 0xbf, 0xa9, 0xf5, 0x49, 0xcf, 0xce, 0xdf,
	0xd9, 0x19, 0xb2, 0xc0, 0xfe, 0x7b, 0x4e, 0xe6, 0xdf, 0xef, 0x46, 0x19, 0xdf, 0xfb, 0x87, 0x00,
	0x00, 0x00, 0xff, 0xff, 0x59, 0x8f, 0x1b, 0x1e, 0x09, 0x02, 0x00, 0x00,
}
