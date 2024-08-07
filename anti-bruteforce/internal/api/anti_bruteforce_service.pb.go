// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.21.12
// source: anti_bruteforce_service.proto

package api

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ReqAllowAuthorization struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Login    string `protobuf:"bytes,1,opt,name=login,proto3" json:"login,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	Ip       string `protobuf:"bytes,3,opt,name=ip,proto3" json:"ip,omitempty"`
}

func (x *ReqAllowAuthorization) Reset() {
	*x = ReqAllowAuthorization{}
	if protoimpl.UnsafeEnabled {
		mi := &file_anti_bruteforce_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReqAllowAuthorization) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqAllowAuthorization) ProtoMessage() {}

func (x *ReqAllowAuthorization) ProtoReflect() protoreflect.Message {
	mi := &file_anti_bruteforce_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqAllowAuthorization.ProtoReflect.Descriptor instead.
func (*ReqAllowAuthorization) Descriptor() ([]byte, []int) {
	return file_anti_bruteforce_service_proto_rawDescGZIP(), []int{0}
}

func (x *ReqAllowAuthorization) GetLogin() string {
	if x != nil {
		return x.Login
	}
	return ""
}

func (x *ReqAllowAuthorization) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *ReqAllowAuthorization) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

type ResAllowAuthorization struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Allow bool `protobuf:"varint,1,opt,name=allow,proto3" json:"allow,omitempty"`
}

func (x *ResAllowAuthorization) Reset() {
	*x = ResAllowAuthorization{}
	if protoimpl.UnsafeEnabled {
		mi := &file_anti_bruteforce_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResAllowAuthorization) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResAllowAuthorization) ProtoMessage() {}

func (x *ResAllowAuthorization) ProtoReflect() protoreflect.Message {
	mi := &file_anti_bruteforce_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResAllowAuthorization.ProtoReflect.Descriptor instead.
func (*ResAllowAuthorization) Descriptor() ([]byte, []int) {
	return file_anti_bruteforce_service_proto_rawDescGZIP(), []int{1}
}

func (x *ResAllowAuthorization) GetAllow() bool {
	if x != nil {
		return x.Allow
	}
	return false
}

type ReqNetwork struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Network string `protobuf:"bytes,1,opt,name=network,proto3" json:"network,omitempty"`
}

func (x *ReqNetwork) Reset() {
	*x = ReqNetwork{}
	if protoimpl.UnsafeEnabled {
		mi := &file_anti_bruteforce_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReqNetwork) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqNetwork) ProtoMessage() {}

func (x *ReqNetwork) ProtoReflect() protoreflect.Message {
	mi := &file_anti_bruteforce_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqNetwork.ProtoReflect.Descriptor instead.
func (*ReqNetwork) Descriptor() ([]byte, []int) {
	return file_anti_bruteforce_service_proto_rawDescGZIP(), []int{2}
}

func (x *ReqNetwork) GetNetwork() string {
	if x != nil {
		return x.Network
	}
	return ""
}

type ResListNetworks struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Networks []string `protobuf:"bytes,1,rep,name=networks,proto3" json:"networks,omitempty"`
}

func (x *ResListNetworks) Reset() {
	*x = ResListNetworks{}
	if protoimpl.UnsafeEnabled {
		mi := &file_anti_bruteforce_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResListNetworks) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResListNetworks) ProtoMessage() {}

func (x *ResListNetworks) ProtoReflect() protoreflect.Message {
	mi := &file_anti_bruteforce_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResListNetworks.ProtoReflect.Descriptor instead.
func (*ResListNetworks) Descriptor() ([]byte, []int) {
	return file_anti_bruteforce_service_proto_rawDescGZIP(), []int{3}
}

func (x *ResListNetworks) GetNetworks() []string {
	if x != nil {
		return x.Networks
	}
	return nil
}

type ReqCleanBucket struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Login string `protobuf:"bytes,1,opt,name=login,proto3" json:"login,omitempty"`
	Ip    string `protobuf:"bytes,3,opt,name=ip,proto3" json:"ip,omitempty"`
}

func (x *ReqCleanBucket) Reset() {
	*x = ReqCleanBucket{}
	if protoimpl.UnsafeEnabled {
		mi := &file_anti_bruteforce_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReqCleanBucket) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqCleanBucket) ProtoMessage() {}

func (x *ReqCleanBucket) ProtoReflect() protoreflect.Message {
	mi := &file_anti_bruteforce_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqCleanBucket.ProtoReflect.Descriptor instead.
func (*ReqCleanBucket) Descriptor() ([]byte, []int) {
	return file_anti_bruteforce_service_proto_rawDescGZIP(), []int{4}
}

func (x *ReqCleanBucket) GetLogin() string {
	if x != nil {
		return x.Login
	}
	return ""
}

func (x *ReqCleanBucket) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

var File_anti_bruteforce_service_proto protoreflect.FileDescriptor

var file_anti_bruteforce_service_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x61, 0x6e, 0x74, 0x69, 0x5f, 0x62, 0x72, 0x75, 0x74, 0x65, 0x66, 0x6f, 0x72, 0x63,
	0x65, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0e, 0x61, 0x6e, 0x74, 0x69, 0x62, 0x72, 0x75, 0x74, 0x65, 0x66, 0x6f, 0x72, 0x63, 0x65, 0x1a,
	0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x59, 0x0a, 0x15,
	0x52, 0x65, 0x71, 0x41, 0x6c, 0x6c, 0x6f, 0x77, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x70,
	0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70,
	0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x70, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x70, 0x22, 0x2d, 0x0a, 0x15, 0x52, 0x65, 0x73, 0x41, 0x6c,
	0x6c, 0x6f, 0x77, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x14, 0x0a, 0x05, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x05, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x22, 0x26, 0x0a, 0x0a, 0x52, 0x65, 0x71, 0x4e, 0x65, 0x74,
	0x77, 0x6f, 0x72, 0x6b, 0x12, 0x18, 0x0a, 0x07, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x22, 0x2d,
	0x0a, 0x0f, 0x52, 0x65, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b,
	0x73, 0x12, 0x1a, 0x0a, 0x08, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x08, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x22, 0x36, 0x0a,
	0x0e, 0x52, 0x65, 0x71, 0x43, 0x6c, 0x65, 0x61, 0x6e, 0x42, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x12,
	0x14, 0x0a, 0x05, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x70, 0x32, 0x9c, 0x05, 0x0a, 0x15, 0x41, 0x6e, 0x74, 0x69, 0x42, 0x72,
	0x75, 0x74, 0x65, 0x46, 0x6f, 0x72, 0x63, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x64, 0x0a, 0x12, 0x41, 0x6c, 0x6c, 0x6f, 0x77, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x25, 0x2e, 0x61, 0x6e, 0x74, 0x69, 0x62, 0x72, 0x75, 0x74,
	0x65, 0x66, 0x6f, 0x72, 0x63, 0x65, 0x2e, 0x52, 0x65, 0x71, 0x41, 0x6c, 0x6c, 0x6f, 0x77, 0x41,
	0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x25, 0x2e, 0x61,
	0x6e, 0x74, 0x69, 0x62, 0x72, 0x75, 0x74, 0x65, 0x66, 0x6f, 0x72, 0x63, 0x65, 0x2e, 0x52, 0x65,
	0x73, 0x41, 0x6c, 0x6c, 0x6f, 0x77, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x22, 0x00, 0x12, 0x4d, 0x0a, 0x15, 0x41, 0x64, 0x64, 0x4e, 0x65, 0x74, 0x77,
	0x6f, 0x72, 0x6b, 0x54, 0x6f, 0x57, 0x68, 0x69, 0x74, 0x65, 0x6c, 0x69, 0x73, 0x74, 0x12, 0x1a,
	0x2e, 0x61, 0x6e, 0x74, 0x69, 0x62, 0x72, 0x75, 0x74, 0x65, 0x66, 0x6f, 0x72, 0x63, 0x65, 0x2e,
	0x52, 0x65, 0x71, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x22, 0x00, 0x12, 0x52, 0x0a, 0x1a, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4e, 0x65,
	0x74, 0x77, 0x6f, 0x72, 0x6b, 0x46, 0x72, 0x6f, 0x6d, 0x57, 0x68, 0x69, 0x74, 0x65, 0x6c, 0x69,
	0x73, 0x74, 0x12, 0x1a, 0x2e, 0x61, 0x6e, 0x74, 0x69, 0x62, 0x72, 0x75, 0x74, 0x65, 0x66, 0x6f,
	0x72, 0x63, 0x65, 0x2e, 0x52, 0x65, 0x71, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x1a, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x46, 0x0a, 0x09, 0x57, 0x68, 0x69, 0x74,
	0x65, 0x6c, 0x69, 0x73, 0x74, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x1f, 0x2e,
	0x61, 0x6e, 0x74, 0x69, 0x62, 0x72, 0x75, 0x74, 0x65, 0x66, 0x6f, 0x72, 0x63, 0x65, 0x2e, 0x52,
	0x65, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x22, 0x00,
	0x12, 0x4d, 0x0a, 0x15, 0x41, 0x64, 0x64, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x54, 0x6f,
	0x42, 0x6c, 0x61, 0x63, 0x6b, 0x6c, 0x69, 0x73, 0x74, 0x12, 0x1a, 0x2e, 0x61, 0x6e, 0x74, 0x69,
	0x62, 0x72, 0x75, 0x74, 0x65, 0x66, 0x6f, 0x72, 0x63, 0x65, 0x2e, 0x52, 0x65, 0x71, 0x4e, 0x65,
	0x74, 0x77, 0x6f, 0x72, 0x6b, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12,
	0x52, 0x0a, 0x1a, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b,
	0x46, 0x72, 0x6f, 0x6d, 0x42, 0x6c, 0x61, 0x63, 0x6b, 0x6c, 0x69, 0x73, 0x74, 0x12, 0x1a, 0x2e,
	0x61, 0x6e, 0x74, 0x69, 0x62, 0x72, 0x75, 0x74, 0x65, 0x66, 0x6f, 0x72, 0x63, 0x65, 0x2e, 0x52,
	0x65, 0x71, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x22, 0x00, 0x12, 0x46, 0x0a, 0x09, 0x42, 0x6c, 0x61, 0x63, 0x6b, 0x6c, 0x69, 0x73, 0x74,
	0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x1f, 0x2e, 0x61, 0x6e, 0x74, 0x69, 0x62,
	0x72, 0x75, 0x74, 0x65, 0x66, 0x6f, 0x72, 0x63, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x4c, 0x69, 0x73,
	0x74, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x22, 0x00, 0x12, 0x47, 0x0a, 0x0b, 0x43,
	0x6c, 0x65, 0x61, 0x6e, 0x42, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x1e, 0x2e, 0x61, 0x6e, 0x74,
	0x69, 0x62, 0x72, 0x75, 0x74, 0x65, 0x66, 0x6f, 0x72, 0x63, 0x65, 0x2e, 0x52, 0x65, 0x71, 0x43,
	0x6c, 0x65, 0x61, 0x6e, 0x42, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x22, 0x00, 0x42, 0x07, 0x5a, 0x05, 0x2e, 0x3b, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_anti_bruteforce_service_proto_rawDescOnce sync.Once
	file_anti_bruteforce_service_proto_rawDescData = file_anti_bruteforce_service_proto_rawDesc
)

func file_anti_bruteforce_service_proto_rawDescGZIP() []byte {
	file_anti_bruteforce_service_proto_rawDescOnce.Do(func() {
		file_anti_bruteforce_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_anti_bruteforce_service_proto_rawDescData)
	})
	return file_anti_bruteforce_service_proto_rawDescData
}

var file_anti_bruteforce_service_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_anti_bruteforce_service_proto_goTypes = []interface{}{
	(*ReqAllowAuthorization)(nil), // 0: antibruteforce.ReqAllowAuthorization
	(*ResAllowAuthorization)(nil), // 1: antibruteforce.ResAllowAuthorization
	(*ReqNetwork)(nil),            // 2: antibruteforce.ReqNetwork
	(*ResListNetworks)(nil),       // 3: antibruteforce.ResListNetworks
	(*ReqCleanBucket)(nil),        // 4: antibruteforce.ReqCleanBucket
	(*emptypb.Empty)(nil),         // 5: google.protobuf.Empty
}
var file_anti_bruteforce_service_proto_depIdxs = []int32{
	0, // 0: antibruteforce.AntiBruteForceService.AllowAuthorization:input_type -> antibruteforce.ReqAllowAuthorization
	2, // 1: antibruteforce.AntiBruteForceService.AddNetworkToWhitelist:input_type -> antibruteforce.ReqNetwork
	2, // 2: antibruteforce.AntiBruteForceService.DeleteNetworkFromWhitelist:input_type -> antibruteforce.ReqNetwork
	5, // 3: antibruteforce.AntiBruteForceService.Whitelist:input_type -> google.protobuf.Empty
	2, // 4: antibruteforce.AntiBruteForceService.AddNetworkToBlacklist:input_type -> antibruteforce.ReqNetwork
	2, // 5: antibruteforce.AntiBruteForceService.DeleteNetworkFromBlacklist:input_type -> antibruteforce.ReqNetwork
	5, // 6: antibruteforce.AntiBruteForceService.Blacklist:input_type -> google.protobuf.Empty
	4, // 7: antibruteforce.AntiBruteForceService.CleanBucket:input_type -> antibruteforce.ReqCleanBucket
	1, // 8: antibruteforce.AntiBruteForceService.AllowAuthorization:output_type -> antibruteforce.ResAllowAuthorization
	5, // 9: antibruteforce.AntiBruteForceService.AddNetworkToWhitelist:output_type -> google.protobuf.Empty
	5, // 10: antibruteforce.AntiBruteForceService.DeleteNetworkFromWhitelist:output_type -> google.protobuf.Empty
	3, // 11: antibruteforce.AntiBruteForceService.Whitelist:output_type -> antibruteforce.ResListNetworks
	5, // 12: antibruteforce.AntiBruteForceService.AddNetworkToBlacklist:output_type -> google.protobuf.Empty
	5, // 13: antibruteforce.AntiBruteForceService.DeleteNetworkFromBlacklist:output_type -> google.protobuf.Empty
	3, // 14: antibruteforce.AntiBruteForceService.Blacklist:output_type -> antibruteforce.ResListNetworks
	5, // 15: antibruteforce.AntiBruteForceService.CleanBucket:output_type -> google.protobuf.Empty
	8, // [8:16] is the sub-list for method output_type
	0, // [0:8] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_anti_bruteforce_service_proto_init() }
func file_anti_bruteforce_service_proto_init() {
	if File_anti_bruteforce_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_anti_bruteforce_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReqAllowAuthorization); i {
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
		file_anti_bruteforce_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResAllowAuthorization); i {
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
		file_anti_bruteforce_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReqNetwork); i {
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
		file_anti_bruteforce_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResListNetworks); i {
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
		file_anti_bruteforce_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReqCleanBucket); i {
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
			RawDescriptor: file_anti_bruteforce_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_anti_bruteforce_service_proto_goTypes,
		DependencyIndexes: file_anti_bruteforce_service_proto_depIdxs,
		MessageInfos:      file_anti_bruteforce_service_proto_msgTypes,
	}.Build()
	File_anti_bruteforce_service_proto = out.File
	file_anti_bruteforce_service_proto_rawDesc = nil
	file_anti_bruteforce_service_proto_goTypes = nil
	file_anti_bruteforce_service_proto_depIdxs = nil
}
