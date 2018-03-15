// Code generated by protoc-gen-go. DO NOT EDIT.
// source: shard_group_reader.proto

package logcache_v1

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "google.golang.org/genproto/googleapis/api/annotations"
import loggregator_v2 "code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type SetShardGroupRequest struct {
	Name     string            `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	SubGroup *GroupedSourceIds `protobuf:"bytes,2,opt,name=sub_group,json=subGroup" json:"sub_group,omitempty"`
	// local_only is used for internals only. A client should not set this.
	LocalOnly bool `protobuf:"varint,3,opt,name=local_only,json=localOnly" json:"local_only,omitempty"`
}

func (m *SetShardGroupRequest) Reset()                    { *m = SetShardGroupRequest{} }
func (m *SetShardGroupRequest) String() string            { return proto.CompactTextString(m) }
func (*SetShardGroupRequest) ProtoMessage()               {}
func (*SetShardGroupRequest) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{0} }

func (m *SetShardGroupRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *SetShardGroupRequest) GetSubGroup() *GroupedSourceIds {
	if m != nil {
		return m.SubGroup
	}
	return nil
}

func (m *SetShardGroupRequest) GetLocalOnly() bool {
	if m != nil {
		return m.LocalOnly
	}
	return false
}

type GroupedSourceIds struct {
	SourceIds []string `protobuf:"bytes,1,rep,name=source_ids,json=sourceIds" json:"source_ids,omitempty"`
}

func (m *GroupedSourceIds) Reset()                    { *m = GroupedSourceIds{} }
func (m *GroupedSourceIds) String() string            { return proto.CompactTextString(m) }
func (*GroupedSourceIds) ProtoMessage()               {}
func (*GroupedSourceIds) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{1} }

func (m *GroupedSourceIds) GetSourceIds() []string {
	if m != nil {
		return m.SourceIds
	}
	return nil
}

type SetShardGroupResponse struct {
}

func (m *SetShardGroupResponse) Reset()                    { *m = SetShardGroupResponse{} }
func (m *SetShardGroupResponse) String() string            { return proto.CompactTextString(m) }
func (*SetShardGroupResponse) ProtoMessage()               {}
func (*SetShardGroupResponse) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{2} }

type ShardGroupReadRequest struct {
	Name          string         `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	RequesterId   uint64         `protobuf:"varint,2,opt,name=requester_id,json=requesterId" json:"requester_id,omitempty"`
	StartTime     int64          `protobuf:"varint,3,opt,name=start_time,json=startTime" json:"start_time,omitempty"`
	EndTime       int64          `protobuf:"varint,4,opt,name=end_time,json=endTime" json:"end_time,omitempty"`
	Limit         int64          `protobuf:"varint,5,opt,name=limit" json:"limit,omitempty"`
	EnvelopeTypes []EnvelopeType `protobuf:"varint,6,rep,packed,name=envelope_types,json=envelopeTypes,enum=logcache.v1.EnvelopeType" json:"envelope_types,omitempty"`
	// local_only is used for internals only. A client should not set this.
	LocalOnly bool `protobuf:"varint,7,opt,name=local_only,json=localOnly" json:"local_only,omitempty"`
}

func (m *ShardGroupReadRequest) Reset()                    { *m = ShardGroupReadRequest{} }
func (m *ShardGroupReadRequest) String() string            { return proto.CompactTextString(m) }
func (*ShardGroupReadRequest) ProtoMessage()               {}
func (*ShardGroupReadRequest) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{3} }

func (m *ShardGroupReadRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ShardGroupReadRequest) GetRequesterId() uint64 {
	if m != nil {
		return m.RequesterId
	}
	return 0
}

func (m *ShardGroupReadRequest) GetStartTime() int64 {
	if m != nil {
		return m.StartTime
	}
	return 0
}

func (m *ShardGroupReadRequest) GetEndTime() int64 {
	if m != nil {
		return m.EndTime
	}
	return 0
}

func (m *ShardGroupReadRequest) GetLimit() int64 {
	if m != nil {
		return m.Limit
	}
	return 0
}

func (m *ShardGroupReadRequest) GetEnvelopeTypes() []EnvelopeType {
	if m != nil {
		return m.EnvelopeTypes
	}
	return nil
}

func (m *ShardGroupReadRequest) GetLocalOnly() bool {
	if m != nil {
		return m.LocalOnly
	}
	return false
}

type ShardGroupReadResponse struct {
	Envelopes *loggregator_v2.EnvelopeBatch `protobuf:"bytes,1,opt,name=envelopes" json:"envelopes,omitempty"`
}

func (m *ShardGroupReadResponse) Reset()                    { *m = ShardGroupReadResponse{} }
func (m *ShardGroupReadResponse) String() string            { return proto.CompactTextString(m) }
func (*ShardGroupReadResponse) ProtoMessage()               {}
func (*ShardGroupReadResponse) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{4} }

func (m *ShardGroupReadResponse) GetEnvelopes() *loggregator_v2.EnvelopeBatch {
	if m != nil {
		return m.Envelopes
	}
	return nil
}

type ShardGroupRequest struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	// local_only is used for internals only. A client should not set this.
	LocalOnly bool `protobuf:"varint,2,opt,name=local_only,json=localOnly" json:"local_only,omitempty"`
}

func (m *ShardGroupRequest) Reset()                    { *m = ShardGroupRequest{} }
func (m *ShardGroupRequest) String() string            { return proto.CompactTextString(m) }
func (*ShardGroupRequest) ProtoMessage()               {}
func (*ShardGroupRequest) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{5} }

func (m *ShardGroupRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ShardGroupRequest) GetLocalOnly() bool {
	if m != nil {
		return m.LocalOnly
	}
	return false
}

type ShardGroupResponse struct {
	SubGroups    []*GroupedSourceIds `protobuf:"bytes,1,rep,name=sub_groups,json=subGroups" json:"sub_groups,omitempty"`
	RequesterIds []uint64            `protobuf:"varint,2,rep,packed,name=requester_ids,json=requesterIds" json:"requester_ids,omitempty"`
}

func (m *ShardGroupResponse) Reset()                    { *m = ShardGroupResponse{} }
func (m *ShardGroupResponse) String() string            { return proto.CompactTextString(m) }
func (*ShardGroupResponse) ProtoMessage()               {}
func (*ShardGroupResponse) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{6} }

func (m *ShardGroupResponse) GetSubGroups() []*GroupedSourceIds {
	if m != nil {
		return m.SubGroups
	}
	return nil
}

func (m *ShardGroupResponse) GetRequesterIds() []uint64 {
	if m != nil {
		return m.RequesterIds
	}
	return nil
}

func init() {
	proto.RegisterType((*SetShardGroupRequest)(nil), "logcache.v1.SetShardGroupRequest")
	proto.RegisterType((*GroupedSourceIds)(nil), "logcache.v1.GroupedSourceIds")
	proto.RegisterType((*SetShardGroupResponse)(nil), "logcache.v1.SetShardGroupResponse")
	proto.RegisterType((*ShardGroupReadRequest)(nil), "logcache.v1.ShardGroupReadRequest")
	proto.RegisterType((*ShardGroupReadResponse)(nil), "logcache.v1.ShardGroupReadResponse")
	proto.RegisterType((*ShardGroupRequest)(nil), "logcache.v1.ShardGroupRequest")
	proto.RegisterType((*ShardGroupResponse)(nil), "logcache.v1.ShardGroupResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for ShardGroupReader service

type ShardGroupReaderClient interface {
	SetShardGroup(ctx context.Context, in *SetShardGroupRequest, opts ...grpc.CallOption) (*SetShardGroupResponse, error)
	Read(ctx context.Context, in *ShardGroupReadRequest, opts ...grpc.CallOption) (*ShardGroupReadResponse, error)
	ShardGroup(ctx context.Context, in *ShardGroupRequest, opts ...grpc.CallOption) (*ShardGroupResponse, error)
}

type shardGroupReaderClient struct {
	cc *grpc.ClientConn
}

func NewShardGroupReaderClient(cc *grpc.ClientConn) ShardGroupReaderClient {
	return &shardGroupReaderClient{cc}
}

func (c *shardGroupReaderClient) SetShardGroup(ctx context.Context, in *SetShardGroupRequest, opts ...grpc.CallOption) (*SetShardGroupResponse, error) {
	out := new(SetShardGroupResponse)
	err := grpc.Invoke(ctx, "/logcache.v1.ShardGroupReader/SetShardGroup", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shardGroupReaderClient) Read(ctx context.Context, in *ShardGroupReadRequest, opts ...grpc.CallOption) (*ShardGroupReadResponse, error) {
	out := new(ShardGroupReadResponse)
	err := grpc.Invoke(ctx, "/logcache.v1.ShardGroupReader/Read", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shardGroupReaderClient) ShardGroup(ctx context.Context, in *ShardGroupRequest, opts ...grpc.CallOption) (*ShardGroupResponse, error) {
	out := new(ShardGroupResponse)
	err := grpc.Invoke(ctx, "/logcache.v1.ShardGroupReader/ShardGroup", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ShardGroupReader service

type ShardGroupReaderServer interface {
	SetShardGroup(context.Context, *SetShardGroupRequest) (*SetShardGroupResponse, error)
	Read(context.Context, *ShardGroupReadRequest) (*ShardGroupReadResponse, error)
	ShardGroup(context.Context, *ShardGroupRequest) (*ShardGroupResponse, error)
}

func RegisterShardGroupReaderServer(s *grpc.Server, srv ShardGroupReaderServer) {
	s.RegisterService(&_ShardGroupReader_serviceDesc, srv)
}

func _ShardGroupReader_SetShardGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetShardGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShardGroupReaderServer).SetShardGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/logcache.v1.ShardGroupReader/SetShardGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShardGroupReaderServer).SetShardGroup(ctx, req.(*SetShardGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ShardGroupReader_Read_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShardGroupReadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShardGroupReaderServer).Read(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/logcache.v1.ShardGroupReader/Read",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShardGroupReaderServer).Read(ctx, req.(*ShardGroupReadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ShardGroupReader_ShardGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShardGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShardGroupReaderServer).ShardGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/logcache.v1.ShardGroupReader/ShardGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShardGroupReaderServer).ShardGroup(ctx, req.(*ShardGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ShardGroupReader_serviceDesc = grpc.ServiceDesc{
	ServiceName: "logcache.v1.ShardGroupReader",
	HandlerType: (*ShardGroupReaderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SetShardGroup",
			Handler:    _ShardGroupReader_SetShardGroup_Handler,
		},
		{
			MethodName: "Read",
			Handler:    _ShardGroupReader_Read_Handler,
		},
		{
			MethodName: "ShardGroup",
			Handler:    _ShardGroupReader_ShardGroup_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "shard_group_reader.proto",
}

func init() { proto.RegisterFile("shard_group_reader.proto", fileDescriptor3) }

var fileDescriptor3 = []byte{
	// 550 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x54, 0xd1, 0x8e, 0xd2, 0x40,
	0x14, 0x4d, 0x81, 0xdd, 0xa5, 0x17, 0xd8, 0xec, 0x4e, 0x76, 0xd7, 0x2e, 0x2b, 0x6b, 0x29, 0x2f,
	0x8d, 0x0f, 0x34, 0xe0, 0xdb, 0xea, 0x83, 0x31, 0x51, 0xb3, 0x4f, 0x26, 0x65, 0x7d, 0x6e, 0x06,
	0x7a, 0x53, 0x9a, 0x94, 0x99, 0x3a, 0x33, 0x60, 0x88, 0xf1, 0xc5, 0xc4, 0x2f, 0xf0, 0xd3, 0xfc,
	0x05, 0xbf, 0xc0, 0x1f, 0xd0, 0x30, 0x2d, 0x58, 0x2a, 0xa0, 0x6f, 0xcc, 0xb9, 0xb7, 0xf7, 0x9c,
	0x39, 0xf7, 0x0c, 0x60, 0xc9, 0x29, 0x15, 0x61, 0x10, 0x09, 0x3e, 0x4f, 0x03, 0x81, 0x34, 0x44,
	0xd1, 0x4f, 0x05, 0x57, 0x9c, 0x34, 0x12, 0x1e, 0x4d, 0xe8, 0x64, 0x8a, 0xfd, 0xc5, 0xa0, 0xfd,
	0x38, 0xe2, 0x3c, 0x4a, 0xd0, 0xa3, 0x69, 0xec, 0x51, 0xc6, 0xb8, 0xa2, 0x2a, 0xe6, 0x4c, 0x66,
	0xad, 0xed, 0xf3, 0xc5, 0xd0, 0x43, 0xb6, 0xc0, 0x84, 0xa7, 0x98, 0x43, 0x4d, 0x8c, 0x04, 0xca,
	0xbc, 0xc1, 0xf9, 0x6a, 0xc0, 0xc5, 0x08, 0xd5, 0x68, 0xc5, 0xf5, 0x76, 0x45, 0xe5, 0xe3, 0x87,
	0x39, 0x4a, 0x45, 0x08, 0xd4, 0x18, 0x9d, 0xa1, 0x65, 0xd8, 0x86, 0x6b, 0xfa, 0xfa, 0x37, 0xb9,
	0x03, 0x53, 0xce, 0xc7, 0x99, 0x24, 0xab, 0x62, 0x1b, 0x6e, 0x63, 0xd8, 0xe9, 0x17, 0xc4, 0xf4,
	0xf5, 0x04, 0x0c, 0x47, 0x7c, 0x2e, 0x26, 0x78, 0x1f, 0x4a, 0xbf, 0x2e, 0xe7, 0x63, 0x0d, 0x92,
	0x0e, 0x40, 0xc2, 0x27, 0x34, 0x09, 0x38, 0x4b, 0x96, 0x56, 0xd5, 0x36, 0xdc, 0xba, 0x6f, 0x6a,
	0xe4, 0x1d, 0x4b, 0x96, 0xce, 0x00, 0xce, 0xca, 0x1f, 0xaf, 0x3e, 0x91, 0xfa, 0x10, 0xc4, 0xa1,
	0xb4, 0x0c, 0xbb, 0xea, 0x9a, 0xbe, 0x29, 0xd7, 0x65, 0xe7, 0x11, 0x5c, 0x96, 0x94, 0xcb, 0x94,
	0x33, 0x89, 0xce, 0x2f, 0x03, 0x2e, 0x8b, 0x30, 0x0d, 0x0f, 0x5d, 0xaa, 0x0b, 0x4d, 0x91, 0x95,
	0x51, 0x04, 0x71, 0xa8, 0xef, 0x55, 0xf3, 0x1b, 0x1b, 0xec, 0x3e, 0xd4, 0x42, 0x14, 0x15, 0x2a,
	0x50, 0xf1, 0x0c, 0xb5, 0xf6, 0xaa, 0x6f, 0x6a, 0xe4, 0x21, 0x9e, 0x21, 0xb9, 0x86, 0x3a, 0xb2,
	0x30, 0x2b, 0xd6, 0x74, 0xf1, 0x04, 0x59, 0xa8, 0x4b, 0x17, 0x70, 0x94, 0xc4, 0xb3, 0x58, 0x59,
	0x47, 0x1a, 0xcf, 0x0e, 0xe4, 0x25, 0x9c, 0xae, 0x97, 0x12, 0xa8, 0x65, 0x8a, 0xd2, 0x3a, 0xb6,
	0xab, 0xee, 0xe9, 0xf0, 0x7a, 0xcb, 0xcc, 0xd7, 0x79, 0xcb, 0xc3, 0x32, 0x45, 0xbf, 0x85, 0x85,
	0x93, 0x2c, 0xb9, 0x79, 0x52, 0x76, 0xf3, 0x3d, 0x5c, 0x95, 0x0d, 0xc8, 0xbc, 0x21, 0xcf, 0xc1,
	0x5c, 0x4f, 0x92, 0xda, 0x86, 0x7c, 0x85, 0x91, 0xc0, 0x88, 0x2a, 0x2e, 0xfa, 0x8b, 0xe1, 0x86,
	0xf8, 0x15, 0x55, 0x93, 0xa9, 0xff, 0xa7, 0xdf, 0x79, 0x03, 0xe7, 0xff, 0x17, 0x94, 0x6d, 0x79,
	0x95, 0xb2, 0xbc, 0x8f, 0x40, 0xfe, 0x5e, 0x1b, 0x79, 0x01, 0xb0, 0x49, 0x57, 0xb6, 0xee, 0x7f,
	0xc6, 0xcb, 0x5c, 0xc7, 0x4b, 0x92, 0x1e, 0xb4, 0x8a, 0x6b, 0x94, 0x56, 0xc5, 0xae, 0xba, 0x35,
	0xbf, 0x59, 0xd8, 0xa3, 0x1c, 0xfe, 0xac, 0xc0, 0xd9, 0xb6, 0x31, 0x28, 0xc8, 0x12, 0x5a, 0x5b,
	0x39, 0x22, 0xdd, 0x2d, 0xd2, 0x5d, 0xaf, 0xa3, 0xed, 0x1c, 0x6a, 0xc9, 0x63, 0xd8, 0xfd, 0xf2,
	0xfd, 0xc7, 0xb7, 0xca, 0x4d, 0xfb, 0xca, 0x5b, 0x0c, 0xbc, 0xc2, 0x63, 0xf6, 0x3e, 0xad, 0x3c,
	0xfa, 0x7c, 0x67, 0x3c, 0x25, 0x1c, 0x6a, 0x2b, 0x11, 0xa4, 0x34, 0x6e, 0x57, 0x76, 0xdb, 0xbd,
	0x83, 0x3d, 0x39, 0xe7, 0xad, 0xe6, 0xb4, 0xc8, 0x1e, 0x4e, 0x22, 0x00, 0x0a, 0x17, 0xbd, 0xdd,
	0x3b, 0x32, 0xa3, 0x7c, 0xb2, 0xb7, 0x9e, 0xd3, 0xf5, 0x34, 0x5d, 0x87, 0xdc, 0xec, 0xa6, 0xf3,
	0x66, 0xa8, 0xe8, 0xf8, 0x58, 0xff, 0xd3, 0x3c, 0xfb, 0x1d, 0x00, 0x00, 0xff, 0xff, 0x61, 0x92,
	0xa2, 0x72, 0xd1, 0x04, 0x00, 0x00,
}