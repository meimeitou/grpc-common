// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: database/query.proto

package database

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Query struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ChQuery       *QueryDefine           `protobuf:"bytes,1,opt,name=ch_query,json=chQuery,proto3" json:"ch_query,omitempty"`
	Query         *QueryDefine           `protobuf:"bytes,2,opt,name=query,proto3" json:"query,omitempty"` // mysql / postgresql
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Query) Reset() {
	*x = Query{}
	mi := &file_database_query_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Query) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Query) ProtoMessage() {}

func (x *Query) ProtoReflect() protoreflect.Message {
	mi := &file_database_query_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Query.ProtoReflect.Descriptor instead.
func (*Query) Descriptor() ([]byte, []int) {
	return file_database_query_proto_rawDescGZIP(), []int{0}
}

func (x *Query) GetChQuery() *QueryDefine {
	if x != nil {
		return x.ChQuery
	}
	return nil
}

func (x *Query) GetQuery() *QueryDefine {
	if x != nil {
		return x.Query
	}
	return nil
}

type QueryDefine struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	SqlKeyMap     map[string]string      `protobuf:"bytes,1,rep,name=sql_key_map,json=sqlKeyMap,proto3" json:"sql_key_map,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"` // sql map key from sql file
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *QueryDefine) Reset() {
	*x = QueryDefine{}
	mi := &file_database_query_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QueryDefine) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryDefine) ProtoMessage() {}

func (x *QueryDefine) ProtoReflect() protoreflect.Message {
	mi := &file_database_query_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryDefine.ProtoReflect.Descriptor instead.
func (*QueryDefine) Descriptor() ([]byte, []int) {
	return file_database_query_proto_rawDescGZIP(), []int{1}
}

func (x *QueryDefine) GetSqlKeyMap() map[string]string {
	if x != nil {
		return x.SqlKeyMap
	}
	return nil
}

var file_database_query_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.MethodOptions)(nil),
		ExtensionType: (*Query)(nil),
		Field:         50002,
		Name:          "grpc_common.database.option.query",
		Tag:           "bytes,50002,opt,name=query",
		Filename:      "database/query.proto",
	},
}

// Extension fields to descriptorpb.MethodOptions.
var (
	// optional grpc_common.database.option.Query query = 50002;
	E_Query = &file_database_query_proto_extTypes[0]
)

var File_database_query_proto protoreflect.FileDescriptor

var file_database_query_proto_rawDesc = string([]byte{
	0x0a, 0x14, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x2f, 0x71, 0x75, 0x65, 0x72, 0x79,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1b, 0x67, 0x72, 0x70, 0x63, 0x5f, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x6f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x8c, 0x01, 0x0a, 0x05, 0x51, 0x75, 0x65, 0x72, 0x79, 0x12,
	0x43, 0x0a, 0x08, 0x63, 0x68, 0x5f, 0x71, 0x75, 0x65, 0x72, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x28, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e,
	0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x2e,
	0x51, 0x75, 0x65, 0x72, 0x79, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x52, 0x07, 0x63, 0x68, 0x51,
	0x75, 0x65, 0x72, 0x79, 0x12, 0x3e, 0x0a, 0x05, 0x71, 0x75, 0x65, 0x72, 0x79, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x6f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x52, 0x05, 0x71,
	0x75, 0x65, 0x72, 0x79, 0x22, 0xa4, 0x01, 0x0a, 0x0b, 0x51, 0x75, 0x65, 0x72, 0x79, 0x44, 0x65,
	0x66, 0x69, 0x6e, 0x65, 0x12, 0x57, 0x0a, 0x0b, 0x73, 0x71, 0x6c, 0x5f, 0x6b, 0x65, 0x79, 0x5f,
	0x6d, 0x61, 0x70, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x37, 0x2e, 0x67, 0x72, 0x70, 0x63,
	0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65,
	0x2e, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x44, 0x65, 0x66,
	0x69, 0x6e, 0x65, 0x2e, 0x53, 0x71, 0x6c, 0x4b, 0x65, 0x79, 0x4d, 0x61, 0x70, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x52, 0x09, 0x73, 0x71, 0x6c, 0x4b, 0x65, 0x79, 0x4d, 0x61, 0x70, 0x1a, 0x3c, 0x0a,
	0x0e, 0x53, 0x71, 0x6c, 0x4b, 0x65, 0x79, 0x4d, 0x61, 0x70, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x3a, 0x5a, 0x0a, 0x05, 0x71,
	0x75, 0x65, 0x72, 0x79, 0x12, 0x1e, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x4f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x18, 0xd2, 0x86, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x67,
	0x72, 0x70, 0x63, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x62,
	0x61, 0x73, 0x65, 0x2e, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79,
	0x52, 0x05, 0x71, 0x75, 0x65, 0x72, 0x79, 0x42, 0xf2, 0x01, 0x0a, 0x1f, 0x63, 0x6f, 0x6d, 0x2e,
	0x67, 0x72, 0x70, 0x63, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x64, 0x61, 0x74, 0x61,
	0x62, 0x61, 0x73, 0x65, 0x2e, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x0a, 0x51, 0x75, 0x65,
	0x72, 0x79, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x39, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x65, 0x69, 0x6d, 0x65, 0x69, 0x74, 0x6f, 0x75, 0x2f,
	0x67, 0x72, 0x70, 0x63, 0x2d, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x3b, 0x64, 0x61, 0x74, 0x61,
	0x62, 0x61, 0x73, 0x65, 0xa2, 0x02, 0x03, 0x47, 0x44, 0x4f, 0xaa, 0x02, 0x1a, 0x47, 0x72, 0x70,
	0x63, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65,
	0x2e, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0xca, 0x02, 0x1a, 0x47, 0x72, 0x70, 0x63, 0x43, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x5c, 0x44, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x5c, 0x4f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0xe2, 0x02, 0x26, 0x47, 0x72, 0x70, 0x63, 0x43, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x5c, 0x44, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x5c, 0x4f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x1c,
	0x47, 0x72, 0x70, 0x63, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x3a, 0x3a, 0x44, 0x61, 0x74, 0x61,
	0x62, 0x61, 0x73, 0x65, 0x3a, 0x3a, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_database_query_proto_rawDescOnce sync.Once
	file_database_query_proto_rawDescData []byte
)

func file_database_query_proto_rawDescGZIP() []byte {
	file_database_query_proto_rawDescOnce.Do(func() {
		file_database_query_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_database_query_proto_rawDesc), len(file_database_query_proto_rawDesc)))
	})
	return file_database_query_proto_rawDescData
}

var file_database_query_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_database_query_proto_goTypes = []any{
	(*Query)(nil),                      // 0: grpc_common.database.option.Query
	(*QueryDefine)(nil),                // 1: grpc_common.database.option.QueryDefine
	nil,                                // 2: grpc_common.database.option.QueryDefine.SqlKeyMapEntry
	(*descriptorpb.MethodOptions)(nil), // 3: google.protobuf.MethodOptions
}
var file_database_query_proto_depIdxs = []int32{
	1, // 0: grpc_common.database.option.Query.ch_query:type_name -> grpc_common.database.option.QueryDefine
	1, // 1: grpc_common.database.option.Query.query:type_name -> grpc_common.database.option.QueryDefine
	2, // 2: grpc_common.database.option.QueryDefine.sql_key_map:type_name -> grpc_common.database.option.QueryDefine.SqlKeyMapEntry
	3, // 3: grpc_common.database.option.query:extendee -> google.protobuf.MethodOptions
	0, // 4: grpc_common.database.option.query:type_name -> grpc_common.database.option.Query
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	4, // [4:5] is the sub-list for extension type_name
	3, // [3:4] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_database_query_proto_init() }
func file_database_query_proto_init() {
	if File_database_query_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_database_query_proto_rawDesc), len(file_database_query_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 1,
			NumServices:   0,
		},
		GoTypes:           file_database_query_proto_goTypes,
		DependencyIndexes: file_database_query_proto_depIdxs,
		MessageInfos:      file_database_query_proto_msgTypes,
		ExtensionInfos:    file_database_query_proto_extTypes,
	}.Build()
	File_database_query_proto = out.File
	file_database_query_proto_goTypes = nil
	file_database_query_proto_depIdxs = nil
}
