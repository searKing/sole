// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v4.25.2
// source: sole/types/v1/configuration/web.proto

package configuration

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Web struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BindAddr              *Web_Net             `protobuf:"bytes,1,opt,name=bind_addr,json=bindAddr,proto3" json:"bind_addr,omitempty"`                                          // for listen
	ShutdownDelayDuration *durationpb.Duration `protobuf:"bytes,2,opt,name=shutdown_delay_duration,json=shutdownDelayDuration,proto3" json:"shutdown_delay_duration,omitempty"` // ShutdownDelayDuration allows to block shutdown for graceful exit.
	Middlewares           *Web_Middlewares     `protobuf:"bytes,3,opt,name=middlewares,proto3" json:"middlewares,omitempty"`                                                    // for middlewares
	// for debug
	ForceDisableTls                bool `protobuf:"varint,20,opt,name=force_disable_tls,json=forceDisableTls,proto3" json:"force_disable_tls,omitempty"`                                                  // disable tls
	PreferRegisterHttpFromEndpoint bool `protobuf:"varint,30,opt,name=prefer_register_http_from_endpoint,json=preferRegisterHttpFromEndpoint,proto3" json:"prefer_register_http_from_endpoint,omitempty"` // prefer register http from endpoint instead of function call, grpc middleware takes effect.
}

func (x *Web) Reset() {
	*x = Web{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sole_types_v1_configuration_web_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Web) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Web) ProtoMessage() {}

func (x *Web) ProtoReflect() protoreflect.Message {
	mi := &file_sole_types_v1_configuration_web_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Web.ProtoReflect.Descriptor instead.
func (*Web) Descriptor() ([]byte, []int) {
	return file_sole_types_v1_configuration_web_proto_rawDescGZIP(), []int{0}
}

func (x *Web) GetBindAddr() *Web_Net {
	if x != nil {
		return x.BindAddr
	}
	return nil
}

func (x *Web) GetShutdownDelayDuration() *durationpb.Duration {
	if x != nil {
		return x.ShutdownDelayDuration
	}
	return nil
}

func (x *Web) GetMiddlewares() *Web_Middlewares {
	if x != nil {
		return x.Middlewares
	}
	return nil
}

func (x *Web) GetForceDisableTls() bool {
	if x != nil {
		return x.ForceDisableTls
	}
	return false
}

func (x *Web) GetPreferRegisterHttpFromEndpoint() bool {
	if x != nil {
		return x.PreferRegisterHttpFromEndpoint
	}
	return false
}

type Web_Net struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Host string `protobuf:"bytes,1,opt,name=host,proto3" json:"host,omitempty"`
	Port int32  `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
}

func (x *Web_Net) Reset() {
	*x = Web_Net{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sole_types_v1_configuration_web_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Web_Net) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Web_Net) ProtoMessage() {}

func (x *Web_Net) ProtoReflect() protoreflect.Message {
	mi := &file_sole_types_v1_configuration_web_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Web_Net.ProtoReflect.Descriptor instead.
func (*Web_Net) Descriptor() ([]byte, []int) {
	return file_sole_types_v1_configuration_web_proto_rawDescGZIP(), []int{0, 0}
}

func (x *Web_Net) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

func (x *Web_Net) GetPort() int32 {
	if x != nil {
		return x.Port
	}
	return 0
}

type Web_Middlewares struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MaxConcurrencyUnary          int64                `protobuf:"varint,1,opt,name=max_concurrency_unary,json=maxConcurrencyUnary,proto3" json:"max_concurrency_unary,omitempty"`                                  // for concurrent parallel requests of unary server, The default is 0 (no limit is given)
	MaxConcurrencyStream         int64                `protobuf:"varint,2,opt,name=max_concurrency_stream,json=maxConcurrencyStream,proto3" json:"max_concurrency_stream,omitempty"`                               // for concurrent parallel requests of stream server, The default is 0 (no limit is given)
	BurstLimitTimeoutUnary       *durationpb.Duration `protobuf:"bytes,3,opt,name=burst_limit_timeout_unary,json=burstLimitTimeoutUnary,proto3" json:"burst_limit_timeout_unary,omitempty"`                        // for concurrent parallel requests of unary server, The default is 0 (no limit is given)
	BurstLimitTimeoutStream      *durationpb.Duration `protobuf:"bytes,4,opt,name=burst_limit_timeout_stream,json=burstLimitTimeoutStream,proto3" json:"burst_limit_timeout_stream,omitempty"`                     // for concurrent parallel requests of stream server, The default is 0 (no limit is given)
	HandledTimeoutUnary          *durationpb.Duration `protobuf:"bytes,5,opt,name=handled_timeout_unary,json=handledTimeoutUnary,proto3" json:"handled_timeout_unary,omitempty"`                                   // for max handing time of unary server, The default is 0 (no limit is given)
	HandledTimeoutStream         *durationpb.Duration `protobuf:"bytes,6,opt,name=handled_timeout_stream,json=handledTimeoutStream,proto3" json:"handled_timeout_stream,omitempty"`                                // for max handing time of unary server, The default is 0 (no limit is given)
	MaxReceiveMessageSizeInBytes int64                `protobuf:"varint,7,opt,name=max_receive_message_size_in_bytes,json=maxReceiveMessageSizeInBytes,proto3" json:"max_receive_message_size_in_bytes,omitempty"` // sets the maximum message size in bytes the grpc server can receive, The default is 0 (no limit is given).
	MaxSendMessageSizeInBytes    int64                `protobuf:"varint,8,opt,name=max_send_message_size_in_bytes,json=maxSendMessageSizeInBytes,proto3" json:"max_send_message_size_in_bytes,omitempty"`          // sets the maximum message size in bytes the grpc server can send, The default is 0 (no limit is given).
	StatsHandling                bool                 `protobuf:"varint,9,opt,name=stats_handling,json=statsHandling,proto3" json:"stats_handling,omitempty"`                                                      // log for the related stats handling (e.g., RPCs, connections).
	FillRequestId                bool                 `protobuf:"varint,10,opt,name=fill_request_id,json=fillRequestId,proto3" json:"fill_request_id,omitempty"`                                                   // for the field "RequestId" filling in Request and Response.
	OtelHandling                 bool                 `protobuf:"varint,11,opt,name=otel_handling,json=otelHandling,proto3" json:"otel_handling,omitempty"`                                                        // captures traces and metrics and send them to an observability platform by OpenTelemetry.
}

func (x *Web_Middlewares) Reset() {
	*x = Web_Middlewares{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sole_types_v1_configuration_web_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Web_Middlewares) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Web_Middlewares) ProtoMessage() {}

func (x *Web_Middlewares) ProtoReflect() protoreflect.Message {
	mi := &file_sole_types_v1_configuration_web_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Web_Middlewares.ProtoReflect.Descriptor instead.
func (*Web_Middlewares) Descriptor() ([]byte, []int) {
	return file_sole_types_v1_configuration_web_proto_rawDescGZIP(), []int{0, 1}
}

func (x *Web_Middlewares) GetMaxConcurrencyUnary() int64 {
	if x != nil {
		return x.MaxConcurrencyUnary
	}
	return 0
}

func (x *Web_Middlewares) GetMaxConcurrencyStream() int64 {
	if x != nil {
		return x.MaxConcurrencyStream
	}
	return 0
}

func (x *Web_Middlewares) GetBurstLimitTimeoutUnary() *durationpb.Duration {
	if x != nil {
		return x.BurstLimitTimeoutUnary
	}
	return nil
}

func (x *Web_Middlewares) GetBurstLimitTimeoutStream() *durationpb.Duration {
	if x != nil {
		return x.BurstLimitTimeoutStream
	}
	return nil
}

func (x *Web_Middlewares) GetHandledTimeoutUnary() *durationpb.Duration {
	if x != nil {
		return x.HandledTimeoutUnary
	}
	return nil
}

func (x *Web_Middlewares) GetHandledTimeoutStream() *durationpb.Duration {
	if x != nil {
		return x.HandledTimeoutStream
	}
	return nil
}

func (x *Web_Middlewares) GetMaxReceiveMessageSizeInBytes() int64 {
	if x != nil {
		return x.MaxReceiveMessageSizeInBytes
	}
	return 0
}

func (x *Web_Middlewares) GetMaxSendMessageSizeInBytes() int64 {
	if x != nil {
		return x.MaxSendMessageSizeInBytes
	}
	return 0
}

func (x *Web_Middlewares) GetStatsHandling() bool {
	if x != nil {
		return x.StatsHandling
	}
	return false
}

func (x *Web_Middlewares) GetFillRequestId() bool {
	if x != nil {
		return x.FillRequestId
	}
	return false
}

func (x *Web_Middlewares) GetOtelHandling() bool {
	if x != nil {
		return x.OtelHandling
	}
	return false
}

var File_sole_types_v1_configuration_web_proto protoreflect.FileDescriptor

var file_sole_types_v1_configuration_web_proto_rawDesc = []byte{
	0x0a, 0x25, 0x73, 0x6f, 0x6c, 0x65, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f, 0x76, 0x31, 0x2f,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x77, 0x65,
	0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x2d, 0x73, 0x65, 0x61, 0x72, 0x6b, 0x69, 0x6e,
	0x67, 0x2e, 0x73, 0x6f, 0x6c, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x6f, 0x6c, 0x65, 0x2e,
	0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xfe, 0x08, 0x0a, 0x03, 0x57, 0x65, 0x62, 0x12, 0x53,
	0x0a, 0x09, 0x62, 0x69, 0x6e, 0x64, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x36, 0x2e, 0x73, 0x65, 0x61, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x73, 0x6f, 0x6c,
	0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x6f, 0x6c, 0x65, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73,
	0x2e, 0x76, 0x31, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x2e, 0x57, 0x65, 0x62, 0x2e, 0x4e, 0x65, 0x74, 0x52, 0x08, 0x62, 0x69, 0x6e, 0x64, 0x41,
	0x64, 0x64, 0x72, 0x12, 0x51, 0x0a, 0x17, 0x73, 0x68, 0x75, 0x74, 0x64, 0x6f, 0x77, 0x6e, 0x5f,
	0x64, 0x65, 0x6c, 0x61, 0x79, 0x5f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x15, 0x73, 0x68, 0x75, 0x74, 0x64, 0x6f, 0x77, 0x6e, 0x44, 0x65, 0x6c, 0x61, 0x79, 0x44, 0x75,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x60, 0x0a, 0x0b, 0x6d, 0x69, 0x64, 0x64, 0x6c, 0x65,
	0x77, 0x61, 0x72, 0x65, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x3e, 0x2e, 0x73, 0x65,
	0x61, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x73, 0x6f, 0x6c, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x73, 0x6f, 0x6c, 0x65, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x57, 0x65, 0x62, 0x2e,
	0x4d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x73, 0x52, 0x0b, 0x6d, 0x69, 0x64,
	0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x73, 0x12, 0x2a, 0x0a, 0x11, 0x66, 0x6f, 0x72, 0x63,
	0x65, 0x5f, 0x64, 0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x74, 0x6c, 0x73, 0x18, 0x14, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x0f, 0x66, 0x6f, 0x72, 0x63, 0x65, 0x44, 0x69, 0x73, 0x61, 0x62, 0x6c,
	0x65, 0x54, 0x6c, 0x73, 0x12, 0x4a, 0x0a, 0x22, 0x70, 0x72, 0x65, 0x66, 0x65, 0x72, 0x5f, 0x72,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x68, 0x74, 0x74, 0x70, 0x5f, 0x66, 0x72, 0x6f,
	0x6d, 0x5f, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x1e, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x1e, 0x70, 0x72, 0x65, 0x66, 0x65, 0x72, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72,
	0x48, 0x74, 0x74, 0x70, 0x46, 0x72, 0x6f, 0x6d, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74,
	0x1a, 0x2d, 0x0a, 0x03, 0x4e, 0x65, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70,
	0x6f, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x1a,
	0xc5, 0x05, 0x0a, 0x0b, 0x4d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x73, 0x12,
	0x32, 0x0a, 0x15, 0x6d, 0x61, 0x78, 0x5f, 0x63, 0x6f, 0x6e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e,
	0x63, 0x79, 0x5f, 0x75, 0x6e, 0x61, 0x72, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x13,
	0x6d, 0x61, 0x78, 0x43, 0x6f, 0x6e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x55, 0x6e,
	0x61, 0x72, 0x79, 0x12, 0x34, 0x0a, 0x16, 0x6d, 0x61, 0x78, 0x5f, 0x63, 0x6f, 0x6e, 0x63, 0x75,
	0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x5f, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x14, 0x6d, 0x61, 0x78, 0x43, 0x6f, 0x6e, 0x63, 0x75, 0x72, 0x72, 0x65,
	0x6e, 0x63, 0x79, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x54, 0x0a, 0x19, 0x62, 0x75, 0x72,
	0x73, 0x74, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74,
	0x5f, 0x75, 0x6e, 0x61, 0x72, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44,
	0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x16, 0x62, 0x75, 0x72, 0x73, 0x74, 0x4c, 0x69,
	0x6d, 0x69, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x55, 0x6e, 0x61, 0x72, 0x79, 0x12,
	0x56, 0x0a, 0x1a, 0x62, 0x75, 0x72, 0x73, 0x74, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x5f, 0x74,
	0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x5f, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x17,
	0x62, 0x75, 0x72, 0x73, 0x74, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75,
	0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x4d, 0x0a, 0x15, 0x68, 0x61, 0x6e, 0x64, 0x6c,
	0x65, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x5f, 0x75, 0x6e, 0x61, 0x72, 0x79,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x13, 0x68, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75,
	0x74, 0x55, 0x6e, 0x61, 0x72, 0x79, 0x12, 0x4f, 0x0a, 0x16, 0x68, 0x61, 0x6e, 0x64, 0x6c, 0x65,
	0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x5f, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x14, 0x68, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75,
	0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x47, 0x0a, 0x21, 0x6d, 0x61, 0x78, 0x5f, 0x72,
	0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x73,
	0x69, 0x7a, 0x65, 0x5f, 0x69, 0x6e, 0x5f, 0x62, 0x79, 0x74, 0x65, 0x73, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x1c, 0x6d, 0x61, 0x78, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x49, 0x6e, 0x42, 0x79, 0x74, 0x65, 0x73,
	0x12, 0x41, 0x0a, 0x1e, 0x6d, 0x61, 0x78, 0x5f, 0x73, 0x65, 0x6e, 0x64, 0x5f, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x5f, 0x69, 0x6e, 0x5f, 0x62, 0x79, 0x74,
	0x65, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x03, 0x52, 0x19, 0x6d, 0x61, 0x78, 0x53, 0x65, 0x6e,
	0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x49, 0x6e, 0x42, 0x79,
	0x74, 0x65, 0x73, 0x12, 0x25, 0x0a, 0x0e, 0x73, 0x74, 0x61, 0x74, 0x73, 0x5f, 0x68, 0x61, 0x6e,
	0x64, 0x6c, 0x69, 0x6e, 0x67, 0x18, 0x09, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0d, 0x73, 0x74, 0x61,
	0x74, 0x73, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x69, 0x6e, 0x67, 0x12, 0x26, 0x0a, 0x0f, 0x66, 0x69,
	0x6c, 0x6c, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x0a, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x0d, 0x66, 0x69, 0x6c, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x49, 0x64, 0x12, 0x23, 0x0a, 0x0d, 0x6f, 0x74, 0x65, 0x6c, 0x5f, 0x68, 0x61, 0x6e, 0x64, 0x6c,
	0x69, 0x6e, 0x67, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x6f, 0x74, 0x65, 0x6c, 0x48,
	0x61, 0x6e, 0x64, 0x6c, 0x69, 0x6e, 0x67, 0x42, 0x56, 0x5a, 0x54, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x65, 0x61, 0x72, 0x4b, 0x69, 0x6e, 0x67, 0x2f, 0x73,
	0x6f, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2d, 0x73, 0x70, 0x65, 0x63, 0x2f, 0x73, 0x6f, 0x6c, 0x65, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73,
	0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x3b, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_sole_types_v1_configuration_web_proto_rawDescOnce sync.Once
	file_sole_types_v1_configuration_web_proto_rawDescData = file_sole_types_v1_configuration_web_proto_rawDesc
)

func file_sole_types_v1_configuration_web_proto_rawDescGZIP() []byte {
	file_sole_types_v1_configuration_web_proto_rawDescOnce.Do(func() {
		file_sole_types_v1_configuration_web_proto_rawDescData = protoimpl.X.CompressGZIP(file_sole_types_v1_configuration_web_proto_rawDescData)
	})
	return file_sole_types_v1_configuration_web_proto_rawDescData
}

var file_sole_types_v1_configuration_web_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_sole_types_v1_configuration_web_proto_goTypes = []interface{}{
	(*Web)(nil),                 // 0: searking.sole.api.sole.types.v1.configuration.Web
	(*Web_Net)(nil),             // 1: searking.sole.api.sole.types.v1.configuration.Web.Net
	(*Web_Middlewares)(nil),     // 2: searking.sole.api.sole.types.v1.configuration.Web.Middlewares
	(*durationpb.Duration)(nil), // 3: google.protobuf.Duration
}
var file_sole_types_v1_configuration_web_proto_depIdxs = []int32{
	1, // 0: searking.sole.api.sole.types.v1.configuration.Web.bind_addr:type_name -> searking.sole.api.sole.types.v1.configuration.Web.Net
	3, // 1: searking.sole.api.sole.types.v1.configuration.Web.shutdown_delay_duration:type_name -> google.protobuf.Duration
	2, // 2: searking.sole.api.sole.types.v1.configuration.Web.middlewares:type_name -> searking.sole.api.sole.types.v1.configuration.Web.Middlewares
	3, // 3: searking.sole.api.sole.types.v1.configuration.Web.Middlewares.burst_limit_timeout_unary:type_name -> google.protobuf.Duration
	3, // 4: searking.sole.api.sole.types.v1.configuration.Web.Middlewares.burst_limit_timeout_stream:type_name -> google.protobuf.Duration
	3, // 5: searking.sole.api.sole.types.v1.configuration.Web.Middlewares.handled_timeout_unary:type_name -> google.protobuf.Duration
	3, // 6: searking.sole.api.sole.types.v1.configuration.Web.Middlewares.handled_timeout_stream:type_name -> google.protobuf.Duration
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_sole_types_v1_configuration_web_proto_init() }
func file_sole_types_v1_configuration_web_proto_init() {
	if File_sole_types_v1_configuration_web_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_sole_types_v1_configuration_web_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Web); i {
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
		file_sole_types_v1_configuration_web_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Web_Net); i {
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
		file_sole_types_v1_configuration_web_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Web_Middlewares); i {
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
			RawDescriptor: file_sole_types_v1_configuration_web_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_sole_types_v1_configuration_web_proto_goTypes,
		DependencyIndexes: file_sole_types_v1_configuration_web_proto_depIdxs,
		MessageInfos:      file_sole_types_v1_configuration_web_proto_msgTypes,
	}.Build()
	File_sole_types_v1_configuration_web_proto = out.File
	file_sole_types_v1_configuration_web_proto_rawDesc = nil
	file_sole_types_v1_configuration_web_proto_goTypes = nil
	file_sole_types_v1_configuration_web_proto_depIdxs = nil
}