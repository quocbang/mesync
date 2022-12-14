// Code generated by protoc-gen-go. DO NOT EDIT.
// source: kenda/mesync/limitary_hour.proto

package mesync

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

type LimitaryHour struct {
	ProductType          string                 `protobuf:"bytes,1,opt,name=product_type,json=productType,proto3" json:"product_type,omitempty"`
	LimitaryHour         *LimitaryHourParameter `protobuf:"bytes,2,opt,name=limitary_hour,json=limitaryHour,proto3" json:"limitary_hour,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *LimitaryHour) Reset()         { *m = LimitaryHour{} }
func (m *LimitaryHour) String() string { return proto.CompactTextString(m) }
func (*LimitaryHour) ProtoMessage()    {}
func (*LimitaryHour) Descriptor() ([]byte, []int) {
	return fileDescriptor_57021d57362868dc, []int{0}
}

func (m *LimitaryHour) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LimitaryHour.Unmarshal(m, b)
}
func (m *LimitaryHour) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LimitaryHour.Marshal(b, m, deterministic)
}
func (m *LimitaryHour) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LimitaryHour.Merge(m, src)
}
func (m *LimitaryHour) XXX_Size() int {
	return xxx_messageInfo_LimitaryHour.Size(m)
}
func (m *LimitaryHour) XXX_DiscardUnknown() {
	xxx_messageInfo_LimitaryHour.DiscardUnknown(m)
}

var xxx_messageInfo_LimitaryHour proto.InternalMessageInfo

func (m *LimitaryHour) GetProductType() string {
	if m != nil {
		return m.ProductType
	}
	return ""
}

func (m *LimitaryHour) GetLimitaryHour() *LimitaryHourParameter {
	if m != nil {
		return m.LimitaryHour
	}
	return nil
}

type LimitaryHourParameter struct {
	Min                  int32    `protobuf:"varint,1,opt,name=min,proto3" json:"min,omitempty"`
	Max                  int32    `protobuf:"varint,2,opt,name=max,proto3" json:"max,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LimitaryHourParameter) Reset()         { *m = LimitaryHourParameter{} }
func (m *LimitaryHourParameter) String() string { return proto.CompactTextString(m) }
func (*LimitaryHourParameter) ProtoMessage()    {}
func (*LimitaryHourParameter) Descriptor() ([]byte, []int) {
	return fileDescriptor_57021d57362868dc, []int{1}
}

func (m *LimitaryHourParameter) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LimitaryHourParameter.Unmarshal(m, b)
}
func (m *LimitaryHourParameter) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LimitaryHourParameter.Marshal(b, m, deterministic)
}
func (m *LimitaryHourParameter) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LimitaryHourParameter.Merge(m, src)
}
func (m *LimitaryHourParameter) XXX_Size() int {
	return xxx_messageInfo_LimitaryHourParameter.Size(m)
}
func (m *LimitaryHourParameter) XXX_DiscardUnknown() {
	xxx_messageInfo_LimitaryHourParameter.DiscardUnknown(m)
}

var xxx_messageInfo_LimitaryHourParameter proto.InternalMessageInfo

func (m *LimitaryHourParameter) GetMin() int32 {
	if m != nil {
		return m.Min
	}
	return 0
}

func (m *LimitaryHourParameter) GetMax() int32 {
	if m != nil {
		return m.Max
	}
	return 0
}

type CreateLimitaryHourRequest struct {
	LimitaryHour         []*LimitaryHour `protobuf:"bytes,1,rep,name=limitary_hour,json=limitaryHour,proto3" json:"limitary_hour,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *CreateLimitaryHourRequest) Reset()         { *m = CreateLimitaryHourRequest{} }
func (m *CreateLimitaryHourRequest) String() string { return proto.CompactTextString(m) }
func (*CreateLimitaryHourRequest) ProtoMessage()    {}
func (*CreateLimitaryHourRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_57021d57362868dc, []int{2}
}

func (m *CreateLimitaryHourRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateLimitaryHourRequest.Unmarshal(m, b)
}
func (m *CreateLimitaryHourRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateLimitaryHourRequest.Marshal(b, m, deterministic)
}
func (m *CreateLimitaryHourRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateLimitaryHourRequest.Merge(m, src)
}
func (m *CreateLimitaryHourRequest) XXX_Size() int {
	return xxx_messageInfo_CreateLimitaryHourRequest.Size(m)
}
func (m *CreateLimitaryHourRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateLimitaryHourRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateLimitaryHourRequest proto.InternalMessageInfo

func (m *CreateLimitaryHourRequest) GetLimitaryHour() []*LimitaryHour {
	if m != nil {
		return m.LimitaryHour
	}
	return nil
}

func init() {
	proto.RegisterType((*LimitaryHour)(nil), "kenda.mesync.LimitaryHour")
	proto.RegisterType((*LimitaryHourParameter)(nil), "kenda.mesync.LimitaryHourParameter")
	proto.RegisterType((*CreateLimitaryHourRequest)(nil), "kenda.mesync.CreateLimitaryHourRequest")
}

func init() { proto.RegisterFile("kenda/mesync/limitary_hour.proto", fileDescriptor_57021d57362868dc) }

var fileDescriptor_57021d57362868dc = []byte{
	// 243 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0xc8, 0x4e, 0xcd, 0x4b,
	0x49, 0xd4, 0xcf, 0x4d, 0x2d, 0xae, 0xcc, 0x4b, 0xd6, 0xcf, 0xc9, 0xcc, 0xcd, 0x2c, 0x49, 0x2c,
	0xaa, 0x8c, 0xcf, 0xc8, 0x2f, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x01, 0xab,
	0xd0, 0x83, 0xa8, 0x50, 0xaa, 0xe6, 0xe2, 0xf1, 0x81, 0x2a, 0xf2, 0xc8, 0x2f, 0x2d, 0x12, 0x52,
	0xe4, 0xe2, 0x29, 0x28, 0xca, 0x4f, 0x29, 0x4d, 0x2e, 0x89, 0x2f, 0xa9, 0x2c, 0x48, 0x95, 0x60,
	0x54, 0x60, 0xd4, 0xe0, 0x0c, 0xe2, 0x86, 0x8a, 0x85, 0x54, 0x16, 0xa4, 0x0a, 0x79, 0x70, 0xf1,
	0xa2, 0x98, 0x2b, 0xc1, 0xa4, 0xc0, 0xa8, 0xc1, 0x6d, 0xa4, 0xac, 0x87, 0x6c, 0xb0, 0x1e, 0xb2,
	0xa9, 0x01, 0x89, 0x45, 0x89, 0xb9, 0xa9, 0x25, 0xa9, 0x45, 0x41, 0x3c, 0x39, 0x48, 0xc2, 0x4a,
	0xd6, 0x5c, 0xa2, 0x58, 0x95, 0x09, 0x09, 0x70, 0x31, 0xe7, 0x66, 0xe6, 0x81, 0x2d, 0x67, 0x0d,
	0x02, 0x31, 0xc1, 0x22, 0x89, 0x15, 0x60, 0xab, 0x40, 0x22, 0x89, 0x15, 0x4a, 0x31, 0x5c, 0x92,
	0xce, 0x45, 0xa9, 0x89, 0x25, 0xa9, 0xc8, 0x46, 0x04, 0xa5, 0x16, 0x96, 0xa6, 0x16, 0x97, 0x08,
	0xd9, 0xa3, 0xbb, 0x91, 0x51, 0x81, 0x59, 0x83, 0xdb, 0x48, 0x0a, 0xb7, 0x1b, 0x51, 0x9d, 0xe6,
	0x64, 0x13, 0x65, 0x95, 0x9e, 0x59, 0x92, 0x93, 0x98, 0x04, 0xd5, 0x91, 0x9c, 0x9f, 0xab, 0x57,
	0x52, 0xae, 0x8f, 0x12, 0xba, 0x05, 0xd9, 0xe9, 0xfa, 0xe0, 0x30, 0x4d, 0x2a, 0x4d, 0x43, 0x91,
	0x49, 0x62, 0x03, 0x0b, 0x1b, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0xba, 0x3f, 0xe9, 0x63, 0x8e,
	0x01, 0x00, 0x00,
}
