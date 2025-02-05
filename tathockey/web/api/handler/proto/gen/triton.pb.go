// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.12.4
// source: gen/triton.proto

package tat_ProfoundPick

import (
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

type VideoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	VideoName string `protobuf:"bytes,1,opt,name=video_name,json=videoName,proto3" json:"video_name,omitempty"`
}

func (x *VideoRequest) Reset() {
	*x = VideoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gen_triton_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VideoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VideoRequest) ProtoMessage() {}

func (x *VideoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gen_triton_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VideoRequest.ProtoReflect.Descriptor instead.
func (*VideoRequest) Descriptor() ([]byte, []int) {
	return file_gen_triton_proto_rawDescGZIP(), []int{0}
}

func (x *VideoRequest) GetVideoName() string {
	if x != nil {
		return x.VideoName
	}
	return ""
}

type VideoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Seconds   []uint64 `protobuf:"varint,1,rep,packed,name=seconds,proto3" json:"seconds,omitempty"`
	EventType []int32  `protobuf:"varint,2,rep,packed,name=event_type,json=eventType,proto3" json:"event_type,omitempty"`
}

func (x *VideoResponse) Reset() {
	*x = VideoResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gen_triton_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VideoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VideoResponse) ProtoMessage() {}

func (x *VideoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_gen_triton_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VideoResponse.ProtoReflect.Descriptor instead.
func (*VideoResponse) Descriptor() ([]byte, []int) {
	return file_gen_triton_proto_rawDescGZIP(), []int{1}
}

func (x *VideoResponse) GetSeconds() []uint64 {
	if x != nil {
		return x.Seconds
	}
	return nil
}

func (x *VideoResponse) GetEventType() []int32 {
	if x != nil {
		return x.EventType
	}
	return nil
}

var File_gen_triton_proto protoreflect.FileDescriptor

var file_gen_triton_proto_rawDesc = []byte{
	0x0a, 0x10, 0x67, 0x65, 0x6e, 0x2f, 0x74, 0x72, 0x69, 0x74, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x10, 0x74, 0x61, 0x74, 0x5f, 0x50, 0x72, 0x6f, 0x66, 0x6f, 0x75, 0x6e, 0x64,
	0x50, 0x69, 0x63, 0x6b, 0x22, 0x2d, 0x0a, 0x0c, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x4e,
	0x61, 0x6d, 0x65, 0x22, 0x48, 0x0a, 0x0d, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x04, 0x52, 0x07, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x12, 0x1d,
	0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x05, 0x52, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x32, 0x5d, 0x0a,
	0x0a, 0x54, 0x72, 0x69, 0x74, 0x6f, 0x6e, 0x50, 0x69, 0x63, 0x6b, 0x12, 0x4f, 0x0a, 0x0c, 0x50,
	0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x12, 0x1e, 0x2e, 0x74, 0x61,
	0x74, 0x5f, 0x50, 0x72, 0x6f, 0x66, 0x6f, 0x75, 0x6e, 0x64, 0x50, 0x69, 0x63, 0x6b, 0x2e, 0x56,
	0x69, 0x64, 0x65, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x74, 0x61,
	0x74, 0x5f, 0x50, 0x72, 0x6f, 0x66, 0x6f, 0x75, 0x6e, 0x64, 0x50, 0x69, 0x63, 0x6b, 0x2e, 0x56,
	0x69, 0x64, 0x65, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x1b, 0x5a, 0x19,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x65, 0x6e, 0x3b, 0x74, 0x61, 0x74, 0x5f, 0x50, 0x72, 0x6f,
	0x66, 0x6f, 0x75, 0x6e, 0x64, 0x50, 0x69, 0x63, 0x6b, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_gen_triton_proto_rawDescOnce sync.Once
	file_gen_triton_proto_rawDescData = file_gen_triton_proto_rawDesc
)

func file_gen_triton_proto_rawDescGZIP() []byte {
	file_gen_triton_proto_rawDescOnce.Do(func() {
		file_gen_triton_proto_rawDescData = protoimpl.X.CompressGZIP(file_gen_triton_proto_rawDescData)
	})
	return file_gen_triton_proto_rawDescData
}

var file_gen_triton_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_gen_triton_proto_goTypes = []any{
	(*VideoRequest)(nil),  // 0: tat_ProfoundPick.VideoRequest
	(*VideoResponse)(nil), // 1: tat_ProfoundPick.VideoResponse
}
var file_gen_triton_proto_depIdxs = []int32{
	0, // 0: tat_ProfoundPick.TritonPick.ProcessVideo:input_type -> tat_ProfoundPick.VideoRequest
	1, // 1: tat_ProfoundPick.TritonPick.ProcessVideo:output_type -> tat_ProfoundPick.VideoResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_gen_triton_proto_init() }
func file_gen_triton_proto_init() {
	if File_gen_triton_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_gen_triton_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*VideoRequest); i {
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
		file_gen_triton_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*VideoResponse); i {
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
			RawDescriptor: file_gen_triton_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_gen_triton_proto_goTypes,
		DependencyIndexes: file_gen_triton_proto_depIdxs,
		MessageInfos:      file_gen_triton_proto_msgTypes,
	}.Build()
	File_gen_triton_proto = out.File
	file_gen_triton_proto_rawDesc = nil
	file_gen_triton_proto_goTypes = nil
	file_gen_triton_proto_depIdxs = nil
}
