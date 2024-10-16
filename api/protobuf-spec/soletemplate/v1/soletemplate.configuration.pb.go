// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v4.25.2
// source: soletemplate/v1/soletemplate.configuration.proto

package v1

import (
	_ "github.com/searKing/sole/api/protobuf-spec/sole/types/v1"
	configuration "github.com/searKing/sole/api/protobuf-spec/sole/types/v1/configuration"
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

type Configuration struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Log      *configuration.Log      `protobuf:"bytes,1,opt,name=log,proto3" json:"log,omitempty"`
	Web      *configuration.Web      `protobuf:"bytes,2,opt,name=web,proto3" json:"web,omitempty"`
	Otel     *configuration.Otel     `protobuf:"bytes,3,opt,name=otel,proto3" json:"otel,omitempty"`
	Category *Configuration_Category `protobuf:"bytes,4,opt,name=category,proto3" json:"category,omitempty"`
}

func (x *Configuration) Reset() {
	*x = Configuration{}
	if protoimpl.UnsafeEnabled {
		mi := &file_soletemplate_v1_soletemplate_configuration_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Configuration) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Configuration) ProtoMessage() {}

func (x *Configuration) ProtoReflect() protoreflect.Message {
	mi := &file_soletemplate_v1_soletemplate_configuration_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Configuration.ProtoReflect.Descriptor instead.
func (*Configuration) Descriptor() ([]byte, []int) {
	return file_soletemplate_v1_soletemplate_configuration_proto_rawDescGZIP(), []int{0}
}

func (x *Configuration) GetLog() *configuration.Log {
	if x != nil {
		return x.Log
	}
	return nil
}

func (x *Configuration) GetWeb() *configuration.Web {
	if x != nil {
		return x.Web
	}
	return nil
}

func (x *Configuration) GetOtel() *configuration.Otel {
	if x != nil {
		return x.Otel
	}
	return nil
}

func (x *Configuration) GetCategory() *Configuration_Category {
	if x != nil {
		return x.Category
	}
	return nil
}

type Configuration_Category struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DynamicEnvironments map[string]string            `protobuf:"bytes,1,rep,name=dynamic_environments,json=dynamicEnvironments,proto3" json:"dynamic_environments,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	System              *configuration.System        `protobuf:"bytes,2,opt,name=system,proto3" json:"system,omitempty"`
	FileCleaners        []*configuration.FileCleaner `protobuf:"bytes,3,rep,name=file_cleaners,json=fileCleaners,proto3" json:"file_cleaners,omitempty"`
}

func (x *Configuration_Category) Reset() {
	*x = Configuration_Category{}
	if protoimpl.UnsafeEnabled {
		mi := &file_soletemplate_v1_soletemplate_configuration_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Configuration_Category) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Configuration_Category) ProtoMessage() {}

func (x *Configuration_Category) ProtoReflect() protoreflect.Message {
	mi := &file_soletemplate_v1_soletemplate_configuration_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Configuration_Category.ProtoReflect.Descriptor instead.
func (*Configuration_Category) Descriptor() ([]byte, []int) {
	return file_soletemplate_v1_soletemplate_configuration_proto_rawDescGZIP(), []int{0, 0}
}

func (x *Configuration_Category) GetDynamicEnvironments() map[string]string {
	if x != nil {
		return x.DynamicEnvironments
	}
	return nil
}

func (x *Configuration_Category) GetSystem() *configuration.System {
	if x != nil {
		return x.System
	}
	return nil
}

func (x *Configuration_Category) GetFileCleaners() []*configuration.FileCleaner {
	if x != nil {
		return x.FileCleaners
	}
	return nil
}

var File_soletemplate_v1_soletemplate_configuration_proto protoreflect.FileDescriptor

var file_soletemplate_v1_soletemplate_configuration_proto_rawDesc = []byte{
	0x0a, 0x30, 0x73, 0x6f, 0x6c, 0x65, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x2f, 0x76,
	0x31, 0x2f, 0x73, 0x6f, 0x6c, 0x65, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x2e, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x21, 0x73, 0x65, 0x61, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x73, 0x6f, 0x6c,
	0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x6f, 0x6c, 0x65, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61,
	0x74, 0x65, 0x2e, 0x76, 0x31, 0x1a, 0x19, 0x73, 0x6f, 0x6c, 0x65, 0x2f, 0x74, 0x79, 0x70, 0x65,
	0x73, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0xc8, 0x05, 0x0a, 0x0d, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x44, 0x0a, 0x03, 0x6c, 0x6f, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x32, 0x2e, 0x73, 0x65, 0x61, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x73, 0x6f, 0x6c, 0x65, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x73, 0x6f, 0x6c, 0x65, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x76,
	0x31, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e,
	0x4c, 0x6f, 0x67, 0x52, 0x03, 0x6c, 0x6f, 0x67, 0x12, 0x44, 0x0a, 0x03, 0x77, 0x65, 0x62, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x32, 0x2e, 0x73, 0x65, 0x61, 0x72, 0x6b, 0x69, 0x6e, 0x67,
	0x2e, 0x73, 0x6f, 0x6c, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x6f, 0x6c, 0x65, 0x2e, 0x74,
	0x79, 0x70, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x57, 0x65, 0x62, 0x52, 0x03, 0x77, 0x65, 0x62, 0x12, 0x47,
	0x0a, 0x04, 0x6f, 0x74, 0x65, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x33, 0x2e, 0x73,
	0x65, 0x61, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x73, 0x6f, 0x6c, 0x65, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x73, 0x6f, 0x6c, 0x65, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x4f, 0x74, 0x65,
	0x6c, 0x52, 0x04, 0x6f, 0x74, 0x65, 0x6c, 0x12, 0x55, 0x0a, 0x08, 0x63, 0x61, 0x74, 0x65, 0x67,
	0x6f, 0x72, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x39, 0x2e, 0x73, 0x65, 0x61, 0x72,
	0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x73, 0x6f, 0x6c, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x6f,
	0x6c, 0x65, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x43, 0x61, 0x74, 0x65,
	0x67, 0x6f, 0x72, 0x79, 0x52, 0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x1a, 0x8a,
	0x03, 0x0a, 0x08, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x12, 0x85, 0x01, 0x0a, 0x14,
	0x64, 0x79, 0x6e, 0x61, 0x6d, 0x69, 0x63, 0x5f, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d,
	0x65, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x52, 0x2e, 0x73, 0x65, 0x61,
	0x72, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x73, 0x6f, 0x6c, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73,
	0x6f, 0x6c, 0x65, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x43, 0x61, 0x74,
	0x65, 0x67, 0x6f, 0x72, 0x79, 0x2e, 0x44, 0x79, 0x6e, 0x61, 0x6d, 0x69, 0x63, 0x45, 0x6e, 0x76,
	0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x13,
	0x64, 0x79, 0x6e, 0x61, 0x6d, 0x69, 0x63, 0x45, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65,
	0x6e, 0x74, 0x73, 0x12, 0x4d, 0x0a, 0x06, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x35, 0x2e, 0x73, 0x65, 0x61, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x73,
	0x6f, 0x6c, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x6f, 0x6c, 0x65, 0x2e, 0x74, 0x79, 0x70,
	0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x52, 0x06, 0x73, 0x79, 0x73, 0x74,
	0x65, 0x6d, 0x12, 0x5f, 0x0a, 0x0d, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x63, 0x6c, 0x65, 0x61, 0x6e,
	0x65, 0x72, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x3a, 0x2e, 0x73, 0x65, 0x61, 0x72,
	0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x73, 0x6f, 0x6c, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x6f,
	0x6c, 0x65, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x43, 0x6c,
	0x65, 0x61, 0x6e, 0x65, 0x72, 0x52, 0x0c, 0x66, 0x69, 0x6c, 0x65, 0x43, 0x6c, 0x65, 0x61, 0x6e,
	0x65, 0x72, 0x73, 0x1a, 0x46, 0x0a, 0x18, 0x44, 0x79, 0x6e, 0x61, 0x6d, 0x69, 0x63, 0x45, 0x6e,
	0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0x3f, 0x5a, 0x3d, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x65, 0x61, 0x72, 0x4b, 0x69,
	0x6e, 0x67, 0x2f, 0x73, 0x6f, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2d, 0x73, 0x70, 0x65, 0x63, 0x2f, 0x73, 0x6f, 0x6c, 0x65, 0x74, 0x65,
	0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x31, 0x3b, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_soletemplate_v1_soletemplate_configuration_proto_rawDescOnce sync.Once
	file_soletemplate_v1_soletemplate_configuration_proto_rawDescData = file_soletemplate_v1_soletemplate_configuration_proto_rawDesc
)

func file_soletemplate_v1_soletemplate_configuration_proto_rawDescGZIP() []byte {
	file_soletemplate_v1_soletemplate_configuration_proto_rawDescOnce.Do(func() {
		file_soletemplate_v1_soletemplate_configuration_proto_rawDescData = protoimpl.X.CompressGZIP(file_soletemplate_v1_soletemplate_configuration_proto_rawDescData)
	})
	return file_soletemplate_v1_soletemplate_configuration_proto_rawDescData
}

var file_soletemplate_v1_soletemplate_configuration_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_soletemplate_v1_soletemplate_configuration_proto_goTypes = []interface{}{
	(*Configuration)(nil),             // 0: searking.sole.api.soletemplate.v1.Configuration
	(*Configuration_Category)(nil),    // 1: searking.sole.api.soletemplate.v1.Configuration.Category
	nil,                               // 2: searking.sole.api.soletemplate.v1.Configuration.Category.DynamicEnvironmentsEntry
	(*configuration.Log)(nil),         // 3: searking.sole.api.sole.types.v1.configuration.Log
	(*configuration.Web)(nil),         // 4: searking.sole.api.sole.types.v1.configuration.Web
	(*configuration.Otel)(nil),        // 5: searking.sole.api.sole.types.v1.configuration.Otel
	(*configuration.System)(nil),      // 6: searking.sole.api.sole.types.v1.configuration.System
	(*configuration.FileCleaner)(nil), // 7: searking.sole.api.sole.types.v1.configuration.FileCleaner
}
var file_soletemplate_v1_soletemplate_configuration_proto_depIdxs = []int32{
	3, // 0: searking.sole.api.soletemplate.v1.Configuration.log:type_name -> searking.sole.api.sole.types.v1.configuration.Log
	4, // 1: searking.sole.api.soletemplate.v1.Configuration.web:type_name -> searking.sole.api.sole.types.v1.configuration.Web
	5, // 2: searking.sole.api.soletemplate.v1.Configuration.otel:type_name -> searking.sole.api.sole.types.v1.configuration.Otel
	1, // 3: searking.sole.api.soletemplate.v1.Configuration.category:type_name -> searking.sole.api.soletemplate.v1.Configuration.Category
	2, // 4: searking.sole.api.soletemplate.v1.Configuration.Category.dynamic_environments:type_name -> searking.sole.api.soletemplate.v1.Configuration.Category.DynamicEnvironmentsEntry
	6, // 5: searking.sole.api.soletemplate.v1.Configuration.Category.system:type_name -> searking.sole.api.sole.types.v1.configuration.System
	7, // 6: searking.sole.api.soletemplate.v1.Configuration.Category.file_cleaners:type_name -> searking.sole.api.sole.types.v1.configuration.FileCleaner
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_soletemplate_v1_soletemplate_configuration_proto_init() }
func file_soletemplate_v1_soletemplate_configuration_proto_init() {
	if File_soletemplate_v1_soletemplate_configuration_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_soletemplate_v1_soletemplate_configuration_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Configuration); i {
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
		file_soletemplate_v1_soletemplate_configuration_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Configuration_Category); i {
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
			RawDescriptor: file_soletemplate_v1_soletemplate_configuration_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_soletemplate_v1_soletemplate_configuration_proto_goTypes,
		DependencyIndexes: file_soletemplate_v1_soletemplate_configuration_proto_depIdxs,
		MessageInfos:      file_soletemplate_v1_soletemplate_configuration_proto_msgTypes,
	}.Build()
	File_soletemplate_v1_soletemplate_configuration_proto = out.File
	file_soletemplate_v1_soletemplate_configuration_proto_rawDesc = nil
	file_soletemplate_v1_soletemplate_configuration_proto_goTypes = nil
	file_soletemplate_v1_soletemplate_configuration_proto_depIdxs = nil
}
