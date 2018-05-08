// Code generated by protoc-gen-go. DO NOT EDIT.
// source: ledger/crosstranset/cross_transet.proto

/*
Package crosstranset is a generated protocol buffer package.

It is generated from these files:
	ledger/crosstranset/cross_transet.proto

It has these top-level messages:
	Version
	CrossTranSet
*/
package crosstranset

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Version struct {
	BlockNum uint64 `protobuf:"varint,1,opt,name=block_num,json=blockNum" json:"block_num,omitempty"`
	TxNum    uint64 `protobuf:"varint,2,opt,name=tx_num,json=txNum" json:"tx_num,omitempty"`
}

func (m *Version) Reset()                    { *m = Version{} }
func (m *Version) String() string            { return proto.CompactTextString(m) }
func (*Version) ProtoMessage()               {}
func (*Version) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Version) GetBlockNum() uint64 {
	if m != nil {
		return m.BlockNum
	}
	return 0
}

func (m *Version) GetTxNum() uint64 {
	if m != nil {
		return m.TxNum
	}
	return 0
}

type CrossTranSet struct {
	From      string `protobuf:"bytes,1,opt,name=from" json:"from,omitempty"`
	Ctranset  []byte `protobuf:"bytes,2,opt,name=ctranset,proto3" json:"ctranset,omitempty"`
	TokenAddr string `protobuf:"bytes,3,opt,name=token_addr,json=tokenAddr" json:"token_addr,omitempty"`
	TokenType string `protobuf:"bytes,4,opt,name=token_type,json=tokenType" json:"token_type,omitempty"`
}

func (m *CrossTranSet) Reset()                    { *m = CrossTranSet{} }
func (m *CrossTranSet) String() string            { return proto.CompactTextString(m) }
func (*CrossTranSet) ProtoMessage()               {}
func (*CrossTranSet) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *CrossTranSet) GetFrom() string {
	if m != nil {
		return m.From
	}
	return ""
}

func (m *CrossTranSet) GetCtranset() []byte {
	if m != nil {
		return m.Ctranset
	}
	return nil
}

func (m *CrossTranSet) GetTokenAddr() string {
	if m != nil {
		return m.TokenAddr
	}
	return ""
}

func (m *CrossTranSet) GetTokenType() string {
	if m != nil {
		return m.TokenType
	}
	return ""
}

func init() {
	proto.RegisterType((*Version)(nil), "crosstranset.Version")
	proto.RegisterType((*CrossTranSet)(nil), "crosstranset.CrossTranSet")
}

func init() { proto.RegisterFile("ledger/crosstranset/cross_transet.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 258 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0xcf, 0x4b, 0xfb, 0x40,
	0x10, 0xc5, 0xe9, 0xf7, 0x1b, 0x6b, 0x33, 0xe4, 0xb4, 0x20, 0x04, 0x45, 0x90, 0x5e, 0xf4, 0x94,
	0x1c, 0x3c, 0x79, 0x10, 0xfc, 0x71, 0xef, 0x21, 0x16, 0x0f, 0x5e, 0xc2, 0x26, 0xbb, 0x4d, 0x43,
	0x92, 0x9d, 0x30, 0x3b, 0x0b, 0x2d, 0xe8, 0xff, 0x2e, 0xdd, 0xad, 0x1a, 0xd0, 0xdb, 0xbe, 0xf7,
	0xd9, 0xf7, 0x86, 0x19, 0xb8, 0xee, 0xb5, 0x6a, 0x34, 0xe5, 0x35, 0xa1, 0xb5, 0x4c, 0xd2, 0x58,
	0xcd, 0x41, 0x94, 0x47, 0x95, 0x8d, 0x84, 0x8c, 0x22, 0x99, 0xfe, 0x58, 0xde, 0xc3, 0xe9, 0xab,
	0x26, 0xdb, 0xa2, 0x11, 0x17, 0x10, 0x57, 0x3d, 0xd6, 0x5d, 0x69, 0xdc, 0x90, 0xce, 0xae, 0x66,
	0x37, 0x51, 0xb1, 0xf0, 0xc6, 0xca, 0x0d, 0xe2, 0x0c, 0xe6, 0xbc, 0xf3, 0xe4, 0x9f, 0x27, 0x27,
	0xbc, 0x5b, 0xb9, 0x61, 0xf9, 0x0e, 0xc9, 0xf3, 0xa1, 0x6e, 0x4d, 0xd2, 0xbc, 0x68, 0x16, 0x02,
	0xa2, 0x0d, 0x61, 0x88, 0xc7, 0x85, 0x7f, 0x8b, 0x73, 0x58, 0xd4, 0xc7, 0x71, 0x3e, 0x9c, 0x14,
	0xdf, 0x5a, 0x5c, 0x02, 0x30, 0x76, 0xda, 0x94, 0x52, 0x29, 0x4a, 0xff, 0xfb, 0x54, 0xec, 0x9d,
	0x47, 0xa5, 0xe8, 0x07, 0xf3, 0x7e, 0xd4, 0x69, 0x34, 0xc1, 0xeb, 0xfd, 0xa8, 0x9f, 0x3e, 0xe0,
	0x0e, 0xa9, 0xc9, 0x5a, 0xd3, 0xf5, 0xb2, 0xb2, 0x1b, 0x74, 0x46, 0x49, 0x6e, 0xd1, 0x1c, 0x9c,
	0x7a, 0x2b, 0x5b, 0x13, 0x56, 0xb6, 0x59, 0xb8, 0x4d, 0x36, 0xdd, 0xfc, 0xed, 0xa1, 0x69, 0x79,
	0xeb, 0xaa, 0xac, 0xc6, 0x21, 0xff, 0xd5, 0x90, 0x7f, 0x35, 0xe4, 0xa1, 0x21, 0xff, 0xe3, 0xba,
	0xd5, 0xdc, 0xb3, 0xdb, 0xcf, 0x00, 0x00, 0x00, 0xff, 0xff, 0xf8, 0x27, 0x61, 0xaa, 0x7b, 0x01,
	0x00, 0x00,
}
