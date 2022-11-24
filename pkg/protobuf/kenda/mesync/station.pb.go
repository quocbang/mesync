// Code generated by protoc-gen-go. DO NOT EDIT.
// source: kenda/mesync/station.proto

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

type StationGroupInformation struct {
	Stations             []string `protobuf:"bytes,1,rep,name=stations,proto3" json:"stations,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StationGroupInformation) Reset()         { *m = StationGroupInformation{} }
func (m *StationGroupInformation) String() string { return proto.CompactTextString(m) }
func (*StationGroupInformation) ProtoMessage()    {}
func (*StationGroupInformation) Descriptor() ([]byte, []int) {
	return fileDescriptor_00aae5f2060fc928, []int{0}
}

func (m *StationGroupInformation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StationGroupInformation.Unmarshal(m, b)
}
func (m *StationGroupInformation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StationGroupInformation.Marshal(b, m, deterministic)
}
func (m *StationGroupInformation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StationGroupInformation.Merge(m, src)
}
func (m *StationGroupInformation) XXX_Size() int {
	return xxx_messageInfo_StationGroupInformation.Size(m)
}
func (m *StationGroupInformation) XXX_DiscardUnknown() {
	xxx_messageInfo_StationGroupInformation.DiscardUnknown(m)
}

var xxx_messageInfo_StationGroupInformation proto.InternalMessageInfo

func (m *StationGroupInformation) GetStations() []string {
	if m != nil {
		return m.Stations
	}
	return nil
}

type CreateStationGroupRequest struct {
	Id                   string                   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Info                 *StationGroupInformation `protobuf:"bytes,2,opt,name=info,proto3" json:"info,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *CreateStationGroupRequest) Reset()         { *m = CreateStationGroupRequest{} }
func (m *CreateStationGroupRequest) String() string { return proto.CompactTextString(m) }
func (*CreateStationGroupRequest) ProtoMessage()    {}
func (*CreateStationGroupRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_00aae5f2060fc928, []int{1}
}

func (m *CreateStationGroupRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateStationGroupRequest.Unmarshal(m, b)
}
func (m *CreateStationGroupRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateStationGroupRequest.Marshal(b, m, deterministic)
}
func (m *CreateStationGroupRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateStationGroupRequest.Merge(m, src)
}
func (m *CreateStationGroupRequest) XXX_Size() int {
	return xxx_messageInfo_CreateStationGroupRequest.Size(m)
}
func (m *CreateStationGroupRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateStationGroupRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateStationGroupRequest proto.InternalMessageInfo

func (m *CreateStationGroupRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *CreateStationGroupRequest) GetInfo() *StationGroupInformation {
	if m != nil {
		return m.Info
	}
	return nil
}

type UpdateStationGroupRequest struct {
	Id                   string                   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Info                 *StationGroupInformation `protobuf:"bytes,2,opt,name=info,proto3" json:"info,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *UpdateStationGroupRequest) Reset()         { *m = UpdateStationGroupRequest{} }
func (m *UpdateStationGroupRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateStationGroupRequest) ProtoMessage()    {}
func (*UpdateStationGroupRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_00aae5f2060fc928, []int{2}
}

func (m *UpdateStationGroupRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateStationGroupRequest.Unmarshal(m, b)
}
func (m *UpdateStationGroupRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateStationGroupRequest.Marshal(b, m, deterministic)
}
func (m *UpdateStationGroupRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateStationGroupRequest.Merge(m, src)
}
func (m *UpdateStationGroupRequest) XXX_Size() int {
	return xxx_messageInfo_UpdateStationGroupRequest.Size(m)
}
func (m *UpdateStationGroupRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateStationGroupRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateStationGroupRequest proto.InternalMessageInfo

func (m *UpdateStationGroupRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *UpdateStationGroupRequest) GetInfo() *StationGroupInformation {
	if m != nil {
		return m.Info
	}
	return nil
}

type DeleteStationGroupRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteStationGroupRequest) Reset()         { *m = DeleteStationGroupRequest{} }
func (m *DeleteStationGroupRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteStationGroupRequest) ProtoMessage()    {}
func (*DeleteStationGroupRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_00aae5f2060fc928, []int{3}
}

func (m *DeleteStationGroupRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteStationGroupRequest.Unmarshal(m, b)
}
func (m *DeleteStationGroupRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteStationGroupRequest.Marshal(b, m, deterministic)
}
func (m *DeleteStationGroupRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteStationGroupRequest.Merge(m, src)
}
func (m *DeleteStationGroupRequest) XXX_Size() int {
	return xxx_messageInfo_DeleteStationGroupRequest.Size(m)
}
func (m *DeleteStationGroupRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteStationGroupRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteStationGroupRequest proto.InternalMessageInfo

func (m *DeleteStationGroupRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func init() {
	proto.RegisterType((*StationGroupInformation)(nil), "kenda.mesync.StationGroupInformation")
	proto.RegisterType((*CreateStationGroupRequest)(nil), "kenda.mesync.CreateStationGroupRequest")
	proto.RegisterType((*UpdateStationGroupRequest)(nil), "kenda.mesync.UpdateStationGroupRequest")
	proto.RegisterType((*DeleteStationGroupRequest)(nil), "kenda.mesync.DeleteStationGroupRequest")
}

func init() { proto.RegisterFile("kenda/mesync/station.proto", fileDescriptor_00aae5f2060fc928) }

var fileDescriptor_00aae5f2060fc928 = []byte{
	// 210 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0xca, 0x4e, 0xcd, 0x4b,
	0x49, 0xd4, 0xcf, 0x4d, 0x2d, 0xae, 0xcc, 0x4b, 0xd6, 0x2f, 0x2e, 0x49, 0x2c, 0xc9, 0xcc, 0xcf,
	0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x01, 0xcb, 0xe9, 0x41, 0xe4, 0x94, 0x4c, 0xb9,
	0xc4, 0x83, 0x21, 0xd2, 0xee, 0x45, 0xf9, 0xa5, 0x05, 0x9e, 0x79, 0x69, 0xf9, 0x45, 0xb9, 0x60,
	0xbe, 0x90, 0x14, 0x17, 0x07, 0x54, 0x67, 0xb1, 0x04, 0xa3, 0x02, 0xb3, 0x06, 0x67, 0x10, 0x9c,
	0xaf, 0x94, 0xc6, 0x25, 0xe9, 0x5c, 0x94, 0x9a, 0x58, 0x92, 0x8a, 0xac, 0x39, 0x28, 0xb5, 0xb0,
	0x34, 0xb5, 0xb8, 0x44, 0x88, 0x8f, 0x8b, 0x29, 0x33, 0x45, 0x82, 0x51, 0x81, 0x51, 0x83, 0x33,
	0x88, 0x29, 0x33, 0x45, 0xc8, 0x92, 0x8b, 0x25, 0x33, 0x2f, 0x2d, 0x5f, 0x82, 0x49, 0x81, 0x51,
	0x83, 0xdb, 0x48, 0x55, 0x0f, 0xd9, 0x01, 0x7a, 0x38, 0x6c, 0x0f, 0x02, 0x6b, 0x01, 0xd9, 0x13,
	0x5a, 0x90, 0x42, 0x7b, 0x7b, 0xb4, 0xb9, 0x24, 0x5d, 0x52, 0x73, 0x52, 0x89, 0xb2, 0xc7, 0xc9,
	0x26, 0xca, 0x2a, 0x3d, 0xb3, 0x24, 0x27, 0x31, 0x09, 0x6a, 0x43, 0x72, 0x7e, 0xae, 0x5e, 0x49,
	0xb9, 0x3e, 0x4a, 0x98, 0x17, 0x64, 0xa7, 0xeb, 0x83, 0xc3, 0x3b, 0xa9, 0x34, 0x0d, 0x45, 0x26,
	0x89, 0x0d, 0x2c, 0x6c, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x5b, 0xe4, 0x59, 0xf9, 0xa4, 0x01,
	0x00, 0x00,
}