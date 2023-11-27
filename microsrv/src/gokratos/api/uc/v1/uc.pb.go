// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.22.2
// source: api/uc/v1/uc.proto

package v1

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type LoginRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AccountType string `protobuf:"bytes,1,opt,name=accountType,proto3" json:"accountType,omitempty"`
	Account     string `protobuf:"bytes,2,opt,name=account,proto3" json:"account,omitempty"`
	Pwd         string `protobuf:"bytes,3,opt,name=pwd,proto3" json:"pwd,omitempty"`
}

func (x *LoginRequest) Reset() {
	*x = LoginRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_uc_v1_uc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoginRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginRequest) ProtoMessage() {}

func (x *LoginRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_uc_v1_uc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginRequest.ProtoReflect.Descriptor instead.
func (*LoginRequest) Descriptor() ([]byte, []int) {
	return file_api_uc_v1_uc_proto_rawDescGZIP(), []int{0}
}

func (x *LoginRequest) GetAccountType() string {
	if x != nil {
		return x.AccountType
	}
	return ""
}

func (x *LoginRequest) GetAccount() string {
	if x != nil {
		return x.Account
	}
	return ""
}

func (x *LoginRequest) GetPwd() string {
	if x != nil {
		return x.Pwd
	}
	return ""
}

type LoginReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AccessToken  string `protobuf:"bytes,1,opt,name=accessToken,proto3" json:"accessToken,omitempty"`
	RefreshToken string `protobuf:"bytes,2,opt,name=refreshToken,proto3" json:"refreshToken,omitempty"`
	TokenType    string `protobuf:"bytes,3,opt,name=tokenType,proto3" json:"tokenType,omitempty"`
	Expires      int32  `protobuf:"varint,4,opt,name=expires,proto3" json:"expires,omitempty"`
	Scrope       string `protobuf:"bytes,5,opt,name=scrope,proto3" json:"scrope,omitempty"`
	Uid          int64  `protobuf:"varint,6,opt,name=uid,proto3" json:"uid,omitempty"`
}

func (x *LoginReply) Reset() {
	*x = LoginReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_uc_v1_uc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoginReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginReply) ProtoMessage() {}

func (x *LoginReply) ProtoReflect() protoreflect.Message {
	mi := &file_api_uc_v1_uc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginReply.ProtoReflect.Descriptor instead.
func (*LoginReply) Descriptor() ([]byte, []int) {
	return file_api_uc_v1_uc_proto_rawDescGZIP(), []int{1}
}

func (x *LoginReply) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

func (x *LoginReply) GetRefreshToken() string {
	if x != nil {
		return x.RefreshToken
	}
	return ""
}

func (x *LoginReply) GetTokenType() string {
	if x != nil {
		return x.TokenType
	}
	return ""
}

func (x *LoginReply) GetExpires() int32 {
	if x != nil {
		return x.Expires
	}
	return 0
}

func (x *LoginReply) GetScrope() string {
	if x != nil {
		return x.Scrope
	}
	return ""
}

func (x *LoginReply) GetUid() int64 {
	if x != nil {
		return x.Uid
	}
	return 0
}

var File_api_uc_v1_uc_proto protoreflect.FileDescriptor

var file_api_uc_v1_uc_proto_rawDesc = []byte{
	0x0a, 0x12, 0x61, 0x70, 0x69, 0x2f, 0x75, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x75, 0x63, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x75, 0x63, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x5c, 0x0a, 0x0c, 0x4c, 0x6f, 0x67,
	0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x61, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b,
	0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x61,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x70, 0x77, 0x64, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x70, 0x77, 0x64, 0x22, 0xb4, 0x01, 0x0a, 0x0a, 0x4c, 0x6f, 0x67, 0x69,
	0x6e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x20, 0x0a, 0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x22, 0x0a, 0x0c, 0x72, 0x65, 0x66, 0x72,
	0x65, 0x73, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c,
	0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x1c, 0x0a, 0x09,
	0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x65, 0x78,
	0x70, 0x69, 0x72, 0x65, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x65, 0x78, 0x70,
	0x69, 0x72, 0x65, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x63, 0x72, 0x6f, 0x70, 0x65, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x63, 0x72, 0x6f, 0x70, 0x65, 0x12, 0x10, 0x0a, 0x03,
	0x75, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x75, 0x69, 0x64, 0x32, 0x5a,
	0x0a, 0x0a, 0x55, 0x73, 0x65, 0x72, 0x43, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x12, 0x4c, 0x0a, 0x05,
	0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x13, 0x2e, 0x75, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x6f,
	0x67, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x75, 0x63, 0x2e,
	0x76, 0x31, 0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x1b, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x15, 0x3a, 0x01, 0x2a, 0x22, 0x10, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x75,
	0x63, 0x2f, 0x76, 0x31, 0x2f, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x42, 0x17, 0x5a, 0x15, 0x67, 0x6f,
	0x6b, 0x72, 0x61, 0x74, 0x6f, 0x73, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x75, 0x63, 0x2f, 0x76, 0x31,
	0x3b, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_uc_v1_uc_proto_rawDescOnce sync.Once
	file_api_uc_v1_uc_proto_rawDescData = file_api_uc_v1_uc_proto_rawDesc
)

func file_api_uc_v1_uc_proto_rawDescGZIP() []byte {
	file_api_uc_v1_uc_proto_rawDescOnce.Do(func() {
		file_api_uc_v1_uc_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_uc_v1_uc_proto_rawDescData)
	})
	return file_api_uc_v1_uc_proto_rawDescData
}

var file_api_uc_v1_uc_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_api_uc_v1_uc_proto_goTypes = []interface{}{
	(*LoginRequest)(nil), // 0: uc.v1.LoginRequest
	(*LoginReply)(nil),   // 1: uc.v1.LoginReply
}
var file_api_uc_v1_uc_proto_depIdxs = []int32{
	0, // 0: uc.v1.UserCenter.Login:input_type -> uc.v1.LoginRequest
	1, // 1: uc.v1.UserCenter.Login:output_type -> uc.v1.LoginReply
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_uc_v1_uc_proto_init() }
func file_api_uc_v1_uc_proto_init() {
	if File_api_uc_v1_uc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_uc_v1_uc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoginRequest); i {
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
		file_api_uc_v1_uc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoginReply); i {
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
			RawDescriptor: file_api_uc_v1_uc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_uc_v1_uc_proto_goTypes,
		DependencyIndexes: file_api_uc_v1_uc_proto_depIdxs,
		MessageInfos:      file_api_uc_v1_uc_proto_msgTypes,
	}.Build()
	File_api_uc_v1_uc_proto = out.File
	file_api_uc_v1_uc_proto_rawDesc = nil
	file_api_uc_v1_uc_proto_goTypes = nil
	file_api_uc_v1_uc_proto_depIdxs = nil
}
