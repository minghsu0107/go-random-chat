// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.9
// source: proto/forwarder/forwarder.proto

package forwarder

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RegisterChannelSessionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChannelId  uint64 `protobuf:"varint,1,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
	UserId     uint64 `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Subscriber string `protobuf:"bytes,3,opt,name=subscriber,proto3" json:"subscriber,omitempty"`
}

func (x *RegisterChannelSessionRequest) Reset() {
	*x = RegisterChannelSessionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_forwarder_forwarder_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterChannelSessionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterChannelSessionRequest) ProtoMessage() {}

func (x *RegisterChannelSessionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_forwarder_forwarder_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterChannelSessionRequest.ProtoReflect.Descriptor instead.
func (*RegisterChannelSessionRequest) Descriptor() ([]byte, []int) {
	return file_proto_forwarder_forwarder_proto_rawDescGZIP(), []int{0}
}

func (x *RegisterChannelSessionRequest) GetChannelId() uint64 {
	if x != nil {
		return x.ChannelId
	}
	return 0
}

func (x *RegisterChannelSessionRequest) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *RegisterChannelSessionRequest) GetSubscriber() string {
	if x != nil {
		return x.Subscriber
	}
	return ""
}

type RegisterChannelSessionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RegisterChannelSessionResponse) Reset() {
	*x = RegisterChannelSessionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_forwarder_forwarder_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterChannelSessionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterChannelSessionResponse) ProtoMessage() {}

func (x *RegisterChannelSessionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_forwarder_forwarder_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterChannelSessionResponse.ProtoReflect.Descriptor instead.
func (*RegisterChannelSessionResponse) Descriptor() ([]byte, []int) {
	return file_proto_forwarder_forwarder_proto_rawDescGZIP(), []int{1}
}

type RemoveChannelSessionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChannelId uint64 `protobuf:"varint,1,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
	UserId    uint64 `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *RemoveChannelSessionRequest) Reset() {
	*x = RemoveChannelSessionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_forwarder_forwarder_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemoveChannelSessionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveChannelSessionRequest) ProtoMessage() {}

func (x *RemoveChannelSessionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_forwarder_forwarder_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveChannelSessionRequest.ProtoReflect.Descriptor instead.
func (*RemoveChannelSessionRequest) Descriptor() ([]byte, []int) {
	return file_proto_forwarder_forwarder_proto_rawDescGZIP(), []int{2}
}

func (x *RemoveChannelSessionRequest) GetChannelId() uint64 {
	if x != nil {
		return x.ChannelId
	}
	return 0
}

func (x *RemoveChannelSessionRequest) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type RemoveChannelSessionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RemoveChannelSessionResponse) Reset() {
	*x = RemoveChannelSessionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_forwarder_forwarder_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemoveChannelSessionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveChannelSessionResponse) ProtoMessage() {}

func (x *RemoveChannelSessionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_forwarder_forwarder_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveChannelSessionResponse.ProtoReflect.Descriptor instead.
func (*RemoveChannelSessionResponse) Descriptor() ([]byte, []int) {
	return file_proto_forwarder_forwarder_proto_rawDescGZIP(), []int{3}
}

var File_proto_forwarder_forwarder_proto protoreflect.FileDescriptor

var file_proto_forwarder_forwarder_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x65,
	0x72, 0x2f, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x09, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x65, 0x72, 0x22, 0x77, 0x0a, 0x1d,
	0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x53,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a,
	0x0a, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x09, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07,
	0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x75,
	0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69,
	0x62, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x75, 0x62, 0x73, 0x63,
	0x72, 0x69, 0x62, 0x65, 0x72, 0x22, 0x20, 0x0a, 0x1e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65,
	0x72, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x55, 0x0a, 0x1b, 0x52, 0x65, 0x6d, 0x6f, 0x76,
	0x65, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65,
	0x6c, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x63, 0x68, 0x61, 0x6e,
	0x6e, 0x65, 0x6c, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x1e,
	0x0a, 0x1c, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x53,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0xec,
	0x01, 0x0a, 0x0e, 0x46, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x6f, 0x0a, 0x16, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x43, 0x68, 0x61,
	0x6e, 0x6e, 0x65, 0x6c, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x28, 0x2e, 0x66, 0x6f,
	0x72, 0x77, 0x61, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72,
	0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x29, 0x2e, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x65,
	0x72, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65,
	0x6c, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x69, 0x0a, 0x14, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x43, 0x68, 0x61, 0x6e,
	0x6e, 0x65, 0x6c, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x26, 0x2e, 0x66, 0x6f, 0x72,
	0x77, 0x61, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x43, 0x68, 0x61,
	0x6e, 0x6e, 0x65, 0x6c, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x27, 0x2e, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x52,
	0x65, 0x6d, 0x6f, 0x76, 0x65, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x53, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x1b, 0x5a,
	0x19, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x65, 0x72,
	0x3b, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_proto_forwarder_forwarder_proto_rawDescOnce sync.Once
	file_proto_forwarder_forwarder_proto_rawDescData = file_proto_forwarder_forwarder_proto_rawDesc
)

func file_proto_forwarder_forwarder_proto_rawDescGZIP() []byte {
	file_proto_forwarder_forwarder_proto_rawDescOnce.Do(func() {
		file_proto_forwarder_forwarder_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_forwarder_forwarder_proto_rawDescData)
	})
	return file_proto_forwarder_forwarder_proto_rawDescData
}

var file_proto_forwarder_forwarder_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_forwarder_forwarder_proto_goTypes = []interface{}{
	(*RegisterChannelSessionRequest)(nil),  // 0: forwarder.RegisterChannelSessionRequest
	(*RegisterChannelSessionResponse)(nil), // 1: forwarder.RegisterChannelSessionResponse
	(*RemoveChannelSessionRequest)(nil),    // 2: forwarder.RemoveChannelSessionRequest
	(*RemoveChannelSessionResponse)(nil),   // 3: forwarder.RemoveChannelSessionResponse
}
var file_proto_forwarder_forwarder_proto_depIdxs = []int32{
	0, // 0: forwarder.ForwardService.RegisterChannelSession:input_type -> forwarder.RegisterChannelSessionRequest
	2, // 1: forwarder.ForwardService.RemoveChannelSession:input_type -> forwarder.RemoveChannelSessionRequest
	1, // 2: forwarder.ForwardService.RegisterChannelSession:output_type -> forwarder.RegisterChannelSessionResponse
	3, // 3: forwarder.ForwardService.RemoveChannelSession:output_type -> forwarder.RemoveChannelSessionResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_forwarder_forwarder_proto_init() }
func file_proto_forwarder_forwarder_proto_init() {
	if File_proto_forwarder_forwarder_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_forwarder_forwarder_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterChannelSessionRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_forwarder_forwarder_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterChannelSessionResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_forwarder_forwarder_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemoveChannelSessionRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_forwarder_forwarder_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemoveChannelSessionResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_forwarder_forwarder_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_forwarder_forwarder_proto_goTypes,
		DependencyIndexes: file_proto_forwarder_forwarder_proto_depIdxs,
		MessageInfos:      file_proto_forwarder_forwarder_proto_msgTypes,
	}.Build()
	File_proto_forwarder_forwarder_proto = out.File
	file_proto_forwarder_forwarder_proto_rawDesc = nil
	file_proto_forwarder_forwarder_proto_goTypes = nil
	file_proto_forwarder_forwarder_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ForwardServiceClient is the client API for ForwardService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ForwardServiceClient interface {
	RegisterChannelSession(ctx context.Context, in *RegisterChannelSessionRequest, opts ...grpc.CallOption) (*RegisterChannelSessionResponse, error)
	RemoveChannelSession(ctx context.Context, in *RemoveChannelSessionRequest, opts ...grpc.CallOption) (*RemoveChannelSessionResponse, error)
}

type forwardServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewForwardServiceClient(cc grpc.ClientConnInterface) ForwardServiceClient {
	return &forwardServiceClient{cc}
}

func (c *forwardServiceClient) RegisterChannelSession(ctx context.Context, in *RegisterChannelSessionRequest, opts ...grpc.CallOption) (*RegisterChannelSessionResponse, error) {
	out := new(RegisterChannelSessionResponse)
	err := c.cc.Invoke(ctx, "/forwarder.ForwardService/RegisterChannelSession", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *forwardServiceClient) RemoveChannelSession(ctx context.Context, in *RemoveChannelSessionRequest, opts ...grpc.CallOption) (*RemoveChannelSessionResponse, error) {
	out := new(RemoveChannelSessionResponse)
	err := c.cc.Invoke(ctx, "/forwarder.ForwardService/RemoveChannelSession", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ForwardServiceServer is the server API for ForwardService service.
type ForwardServiceServer interface {
	RegisterChannelSession(context.Context, *RegisterChannelSessionRequest) (*RegisterChannelSessionResponse, error)
	RemoveChannelSession(context.Context, *RemoveChannelSessionRequest) (*RemoveChannelSessionResponse, error)
}

// UnimplementedForwardServiceServer can be embedded to have forward compatible implementations.
type UnimplementedForwardServiceServer struct {
}

func (*UnimplementedForwardServiceServer) RegisterChannelSession(context.Context, *RegisterChannelSessionRequest) (*RegisterChannelSessionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterChannelSession not implemented")
}
func (*UnimplementedForwardServiceServer) RemoveChannelSession(context.Context, *RemoveChannelSessionRequest) (*RemoveChannelSessionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveChannelSession not implemented")
}

func RegisterForwardServiceServer(s *grpc.Server, srv ForwardServiceServer) {
	s.RegisterService(&_ForwardService_serviceDesc, srv)
}

func _ForwardService_RegisterChannelSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterChannelSessionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ForwardServiceServer).RegisterChannelSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/forwarder.ForwardService/RegisterChannelSession",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ForwardServiceServer).RegisterChannelSession(ctx, req.(*RegisterChannelSessionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ForwardService_RemoveChannelSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveChannelSessionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ForwardServiceServer).RemoveChannelSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/forwarder.ForwardService/RemoveChannelSession",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ForwardServiceServer).RemoveChannelSession(ctx, req.(*RemoveChannelSessionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ForwardService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "forwarder.ForwardService",
	HandlerType: (*ForwardServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterChannelSession",
			Handler:    _ForwardService_RegisterChannelSession_Handler,
		},
		{
			MethodName: "RemoveChannelSession",
			Handler:    _ForwardService_RemoveChannelSession_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/forwarder/forwarder.proto",
}