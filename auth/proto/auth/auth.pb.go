// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/auth/auth.proto

package fullbottle_srv_auth

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type GenerateJwtTokenRequest struct {
	UserId               int64    `protobuf:"varint,1,opt,name=userId,proto3" json:"userId,omitempty"`
	Expire               int64    `protobuf:"varint,2,opt,name=expire,proto3" json:"expire,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GenerateJwtTokenRequest) Reset()         { *m = GenerateJwtTokenRequest{} }
func (m *GenerateJwtTokenRequest) String() string { return proto.CompactTextString(m) }
func (*GenerateJwtTokenRequest) ProtoMessage()    {}
func (*GenerateJwtTokenRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_82b5829f48cfb8e5, []int{0}
}

func (m *GenerateJwtTokenRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GenerateJwtTokenRequest.Unmarshal(m, b)
}
func (m *GenerateJwtTokenRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GenerateJwtTokenRequest.Marshal(b, m, deterministic)
}
func (m *GenerateJwtTokenRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenerateJwtTokenRequest.Merge(m, src)
}
func (m *GenerateJwtTokenRequest) XXX_Size() int {
	return xxx_messageInfo_GenerateJwtTokenRequest.Size(m)
}
func (m *GenerateJwtTokenRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GenerateJwtTokenRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GenerateJwtTokenRequest proto.InternalMessageInfo

func (m *GenerateJwtTokenRequest) GetUserId() int64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *GenerateJwtTokenRequest) GetExpire() int64 {
	if m != nil {
		return m.Expire
	}
	return 0
}

type GenerateJwtTokenResponse struct {
	Code                 int64    `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg                  string   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Token                string   `protobuf:"bytes,3,opt,name=token,proto3" json:"token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GenerateJwtTokenResponse) Reset()         { *m = GenerateJwtTokenResponse{} }
func (m *GenerateJwtTokenResponse) String() string { return proto.CompactTextString(m) }
func (*GenerateJwtTokenResponse) ProtoMessage()    {}
func (*GenerateJwtTokenResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_82b5829f48cfb8e5, []int{1}
}

func (m *GenerateJwtTokenResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GenerateJwtTokenResponse.Unmarshal(m, b)
}
func (m *GenerateJwtTokenResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GenerateJwtTokenResponse.Marshal(b, m, deterministic)
}
func (m *GenerateJwtTokenResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenerateJwtTokenResponse.Merge(m, src)
}
func (m *GenerateJwtTokenResponse) XXX_Size() int {
	return xxx_messageInfo_GenerateJwtTokenResponse.Size(m)
}
func (m *GenerateJwtTokenResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GenerateJwtTokenResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GenerateJwtTokenResponse proto.InternalMessageInfo

func (m *GenerateJwtTokenResponse) GetCode() int64 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *GenerateJwtTokenResponse) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func (m *GenerateJwtTokenResponse) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type ParseJwtTokenRequest struct {
	Token                string   `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ParseJwtTokenRequest) Reset()         { *m = ParseJwtTokenRequest{} }
func (m *ParseJwtTokenRequest) String() string { return proto.CompactTextString(m) }
func (*ParseJwtTokenRequest) ProtoMessage()    {}
func (*ParseJwtTokenRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_82b5829f48cfb8e5, []int{2}
}

func (m *ParseJwtTokenRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ParseJwtTokenRequest.Unmarshal(m, b)
}
func (m *ParseJwtTokenRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ParseJwtTokenRequest.Marshal(b, m, deterministic)
}
func (m *ParseJwtTokenRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ParseJwtTokenRequest.Merge(m, src)
}
func (m *ParseJwtTokenRequest) XXX_Size() int {
	return xxx_messageInfo_ParseJwtTokenRequest.Size(m)
}
func (m *ParseJwtTokenRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ParseJwtTokenRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ParseJwtTokenRequest proto.InternalMessageInfo

func (m *ParseJwtTokenRequest) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type ParseJwtTokenResponse struct {
	Code                 int64    `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg                  string   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	UserId               int64    `protobuf:"varint,3,opt,name=userId,proto3" json:"userId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ParseJwtTokenResponse) Reset()         { *m = ParseJwtTokenResponse{} }
func (m *ParseJwtTokenResponse) String() string { return proto.CompactTextString(m) }
func (*ParseJwtTokenResponse) ProtoMessage()    {}
func (*ParseJwtTokenResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_82b5829f48cfb8e5, []int{3}
}

func (m *ParseJwtTokenResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ParseJwtTokenResponse.Unmarshal(m, b)
}
func (m *ParseJwtTokenResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ParseJwtTokenResponse.Marshal(b, m, deterministic)
}
func (m *ParseJwtTokenResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ParseJwtTokenResponse.Merge(m, src)
}
func (m *ParseJwtTokenResponse) XXX_Size() int {
	return xxx_messageInfo_ParseJwtTokenResponse.Size(m)
}
func (m *ParseJwtTokenResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ParseJwtTokenResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ParseJwtTokenResponse proto.InternalMessageInfo

func (m *ParseJwtTokenResponse) GetCode() int64 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *ParseJwtTokenResponse) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func (m *ParseJwtTokenResponse) GetUserId() int64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

func init() {
	proto.RegisterType((*GenerateJwtTokenRequest)(nil), "fullbottle.srv.auth.GenerateJwtTokenRequest")
	proto.RegisterType((*GenerateJwtTokenResponse)(nil), "fullbottle.srv.auth.GenerateJwtTokenResponse")
	proto.RegisterType((*ParseJwtTokenRequest)(nil), "fullbottle.srv.auth.ParseJwtTokenRequest")
	proto.RegisterType((*ParseJwtTokenResponse)(nil), "fullbottle.srv.auth.ParseJwtTokenResponse")
}

func init() { proto.RegisterFile("proto/auth/auth.proto", fileDescriptor_82b5829f48cfb8e5) }

var fileDescriptor_82b5829f48cfb8e5 = []byte{
	// 263 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2d, 0x28, 0xca, 0x2f,
	0xc9, 0xd7, 0x4f, 0x2c, 0x2d, 0xc9, 0x00, 0x13, 0x7a, 0x60, 0xbe, 0x90, 0x70, 0x5a, 0x69, 0x4e,
	0x4e, 0x52, 0x7e, 0x49, 0x49, 0x4e, 0xaa, 0x5e, 0x71, 0x51, 0x99, 0x1e, 0x48, 0x4a, 0xc9, 0x93,
	0x4b, 0xdc, 0x3d, 0x35, 0x2f, 0xb5, 0x28, 0xb1, 0x24, 0xd5, 0xab, 0xbc, 0x24, 0x24, 0x3f, 0x3b,
	0x35, 0x2f, 0x28, 0xb5, 0xb0, 0x34, 0xb5, 0xb8, 0x44, 0x48, 0x8c, 0x8b, 0xad, 0xb4, 0x38, 0xb5,
	0xc8, 0x33, 0x45, 0x82, 0x51, 0x81, 0x51, 0x83, 0x39, 0x08, 0xca, 0x03, 0x89, 0xa7, 0x56, 0x14,
	0x64, 0x16, 0xa5, 0x4a, 0x30, 0x41, 0xc4, 0x21, 0x3c, 0xa5, 0x30, 0x2e, 0x09, 0x4c, 0xa3, 0x8a,
	0x0b, 0xf2, 0xf3, 0x8a, 0x53, 0x85, 0x84, 0xb8, 0x58, 0x92, 0xf3, 0x53, 0x52, 0xa1, 0x26, 0x81,
	0xd9, 0x42, 0x02, 0x5c, 0xcc, 0xb9, 0xc5, 0xe9, 0x60, 0x43, 0x38, 0x83, 0x40, 0x4c, 0x21, 0x11,
	0x2e, 0xd6, 0x12, 0x90, 0x36, 0x09, 0x66, 0xb0, 0x18, 0x84, 0xa3, 0xa4, 0xc3, 0x25, 0x12, 0x90,
	0x58, 0x54, 0x8c, 0xe1, 0x3e, 0xb8, 0x6a, 0x46, 0x64, 0xd5, 0xa1, 0x5c, 0xa2, 0x68, 0xaa, 0x49,
	0x72, 0x02, 0xc2, 0xd3, 0xcc, 0xc8, 0x9e, 0x36, 0x7a, 0xc5, 0xc8, 0xc5, 0xed, 0x58, 0x5a, 0x92,
	0x11, 0x9c, 0x5a, 0x54, 0x96, 0x99, 0x9c, 0x2a, 0x54, 0xc8, 0x25, 0x80, 0xee, 0x59, 0x21, 0x1d,
	0x3d, 0x2c, 0x21, 0xac, 0x87, 0x23, 0x78, 0xa5, 0x74, 0x89, 0x54, 0x0d, 0x71, 0xbe, 0x12, 0x83,
	0x50, 0x06, 0x17, 0x2f, 0x8a, 0xcf, 0x84, 0x34, 0xb1, 0x9a, 0x80, 0x2d, 0xac, 0xa4, 0xb4, 0x88,
	0x51, 0x0a, 0xb3, 0x29, 0x89, 0x0d, 0x9c, 0x60, 0x8c, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0xb6,
	0x4e, 0x5c, 0xd0, 0x49, 0x02, 0x00, 0x00,
}