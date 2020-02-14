// Code generated by protoc-gen-go. DO NOT EDIT.
// source: search.proto

package proto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/options"
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

type SearchRequest struct {
	// The search query.
	Query                string   `protobuf:"bytes,1,opt,name=query,proto3" json:"query,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SearchRequest) Reset()         { *m = SearchRequest{} }
func (m *SearchRequest) String() string { return proto.CompactTextString(m) }
func (*SearchRequest) ProtoMessage()    {}
func (*SearchRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_453745cff914010e, []int{0}
}

func (m *SearchRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SearchRequest.Unmarshal(m, b)
}
func (m *SearchRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SearchRequest.Marshal(b, m, deterministic)
}
func (m *SearchRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SearchRequest.Merge(m, src)
}
func (m *SearchRequest) XXX_Size() int {
	return xxx_messageInfo_SearchRequest.Size(m)
}
func (m *SearchRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SearchRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SearchRequest proto.InternalMessageInfo

func (m *SearchRequest) GetQuery() string {
	if m != nil {
		return m.Query
	}
	return ""
}

type SearchResponse struct {
	// The search hits.
	Hits                 []*DocumentResult `protobuf:"bytes,1,rep,name=hits,proto3" json:"hits,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *SearchResponse) Reset()         { *m = SearchResponse{} }
func (m *SearchResponse) String() string { return proto.CompactTextString(m) }
func (*SearchResponse) ProtoMessage()    {}
func (*SearchResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_453745cff914010e, []int{1}
}

func (m *SearchResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SearchResponse.Unmarshal(m, b)
}
func (m *SearchResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SearchResponse.Marshal(b, m, deterministic)
}
func (m *SearchResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SearchResponse.Merge(m, src)
}
func (m *SearchResponse) XXX_Size() int {
	return xxx_messageInfo_SearchResponse.Size(m)
}
func (m *SearchResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SearchResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SearchResponse proto.InternalMessageInfo

func (m *SearchResponse) GetHits() []*DocumentResult {
	if m != nil {
		return m.Hits
	}
	return nil
}

type DocumentResult struct {
	// The document ID that contains the query terms.
	DocID string `protobuf:"bytes,1,opt,name=docID,proto3" json:"docID,omitempty"`
	// The score of the search result.
	Score                float64  `protobuf:"fixed64,2,opt,name=score,proto3" json:"score,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DocumentResult) Reset()         { *m = DocumentResult{} }
func (m *DocumentResult) String() string { return proto.CompactTextString(m) }
func (*DocumentResult) ProtoMessage()    {}
func (*DocumentResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_453745cff914010e, []int{2}
}

func (m *DocumentResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DocumentResult.Unmarshal(m, b)
}
func (m *DocumentResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DocumentResult.Marshal(b, m, deterministic)
}
func (m *DocumentResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DocumentResult.Merge(m, src)
}
func (m *DocumentResult) XXX_Size() int {
	return xxx_messageInfo_DocumentResult.Size(m)
}
func (m *DocumentResult) XXX_DiscardUnknown() {
	xxx_messageInfo_DocumentResult.DiscardUnknown(m)
}

var xxx_messageInfo_DocumentResult proto.InternalMessageInfo

func (m *DocumentResult) GetDocID() string {
	if m != nil {
		return m.DocID
	}
	return ""
}

func (m *DocumentResult) GetScore() float64 {
	if m != nil {
		return m.Score
	}
	return 0
}

type GetIndexSizeRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetIndexSizeRequest) Reset()         { *m = GetIndexSizeRequest{} }
func (m *GetIndexSizeRequest) String() string { return proto.CompactTextString(m) }
func (*GetIndexSizeRequest) ProtoMessage()    {}
func (*GetIndexSizeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_453745cff914010e, []int{3}
}

func (m *GetIndexSizeRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetIndexSizeRequest.Unmarshal(m, b)
}
func (m *GetIndexSizeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetIndexSizeRequest.Marshal(b, m, deterministic)
}
func (m *GetIndexSizeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetIndexSizeRequest.Merge(m, src)
}
func (m *GetIndexSizeRequest) XXX_Size() int {
	return xxx_messageInfo_GetIndexSizeRequest.Size(m)
}
func (m *GetIndexSizeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetIndexSizeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetIndexSizeRequest proto.InternalMessageInfo

type GetIndexSizeResponse struct {
	// The size of the index.
	Size                 int64    `protobuf:"varint,1,opt,name=size,proto3" json:"size,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetIndexSizeResponse) Reset()         { *m = GetIndexSizeResponse{} }
func (m *GetIndexSizeResponse) String() string { return proto.CompactTextString(m) }
func (*GetIndexSizeResponse) ProtoMessage()    {}
func (*GetIndexSizeResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_453745cff914010e, []int{4}
}

func (m *GetIndexSizeResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetIndexSizeResponse.Unmarshal(m, b)
}
func (m *GetIndexSizeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetIndexSizeResponse.Marshal(b, m, deterministic)
}
func (m *GetIndexSizeResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetIndexSizeResponse.Merge(m, src)
}
func (m *GetIndexSizeResponse) XXX_Size() int {
	return xxx_messageInfo_GetIndexSizeResponse.Size(m)
}
func (m *GetIndexSizeResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetIndexSizeResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetIndexSizeResponse proto.InternalMessageInfo

func (m *GetIndexSizeResponse) GetSize() int64 {
	if m != nil {
		return m.Size
	}
	return 0
}

func init() {
	proto.RegisterType((*SearchRequest)(nil), "proto.SearchRequest")
	proto.RegisterType((*SearchResponse)(nil), "proto.SearchResponse")
	proto.RegisterType((*DocumentResult)(nil), "proto.DocumentResult")
	proto.RegisterType((*GetIndexSizeRequest)(nil), "proto.GetIndexSizeRequest")
	proto.RegisterType((*GetIndexSizeResponse)(nil), "proto.GetIndexSizeResponse")
}

func init() { proto.RegisterFile("search.proto", fileDescriptor_453745cff914010e) }

var fileDescriptor_453745cff914010e = []byte{
	// 377 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x50, 0x4d, 0x8b, 0xdb, 0x30,
	0x10, 0x45, 0xdd, 0x0f, 0xa8, 0xda, 0x6e, 0x17, 0x75, 0xb7, 0x0d, 0x7b, 0x58, 0xc4, 0x9c, 0x52,
	0x37, 0x9b, 0x78, 0xbd, 0xd0, 0x83, 0xd8, 0x8b, 0x43, 0xa0, 0xe4, 0x56, 0x9c, 0x63, 0x08, 0xc5,
	0xb1, 0x07, 0x47, 0x34, 0x91, 0x12, 0x4b, 0xee, 0x47, 0x8c, 0xef, 0x3d, 0x16, 0xff, 0xe2, 0x62,
	0xd9, 0x39, 0xa4, 0xe4, 0xa4, 0x99, 0xd1, 0x7b, 0x6f, 0xe6, 0x3d, 0xfa, 0xda, 0x60, 0x9c, 0x27,
	0xab, 0xe1, 0x36, 0xd7, 0x56, 0xb3, 0x0b, 0xf7, 0xdc, 0x0d, 0xdc, 0x93, 0x3c, 0x64, 0xa8, 0x1e,
	0xcc, 0xcf, 0x38, 0xcb, 0x30, 0x1f, 0xe9, 0xad, 0x95, 0x5a, 0x99, 0x51, 0xac, 0x94, 0xb6, 0xb1,
	0xab, 0x5b, 0x12, 0x2c, 0xe8, 0x9b, 0x99, 0x13, 0x89, 0x70, 0x57, 0xa0, 0xb1, 0xec, 0x86, 0x5e,
	0xec, 0x0a, 0xcc, 0x7f, 0xf7, 0x08, 0x27, 0xfd, 0x97, 0x51, 0xdb, 0x88, 0xcf, 0x75, 0xf8, 0x44,
	0xdf, 0x7a, 0xc7, 0xd8, 0x80, 0xb3, 0xfb, 0x12, 0x1c, 0x00, 0x04, 0x87, 0xe2, 0x7b, 0x1e, 0x4b,
	0x85, 0x3c, 0xd5, 0x49, 0xb1, 0x41, 0x65, 0x0d, 0x54, 0xf0, 0x97, 0xd0, 0xab, 0x03, 0xc7, 0x6c,
	0xb5, 0x32, 0xc8, 0x3e, 0xd2, 0xf3, 0x95, 0xb4, 0xa6, 0x47, 0xf8, 0x59, 0xff, 0x55, 0x70, 0xdb,
	0xde, 0x31, 0x9c, 0x74, 0xa4, 0x08, 0x4d, 0xb1, 0xb6, 0x91, 0x83, 0x88, 0x59, 0x1d, 0x7e, 0xa5,
	0xd7, 0xde, 0x7f, 0x0a, 0xc1, 0x33, 0x13, 0x25, 0x34, 0x00, 0x10, 0x7c, 0x5e, 0x42, 0xaa, 0x93,
	0xe9, 0xa4, 0x39, 0x60, 0x29, 0x53, 0x54, 0xdf, 0xa4, 0xfa, 0x81, 0xc6, 0xca, 0xcc, 0xd9, 0x84,
	0x01, 0x07, 0x93, 0xe8, 0x1c, 0x41, 0xf0, 0x47, 0xdf, 0xaf, 0x16, 0x15, 0x3c, 0xd3, 0xab, 0xe3,
	0x65, 0x8d, 0x65, 0xa7, 0x71, 0xb0, 0xec, 0x9a, 0x66, 0xea, 0x98, 0xbd, 0x17, 0x9c, 0xf4, 0x49,
	0xd4, 0x36, 0xf0, 0x48, 0xdf, 0x7d, 0x41, 0x3b, 0x55, 0x29, 0xfe, 0x9a, 0xc9, 0x3d, 0x76, 0x49,
	0x88, 0xbb, 0x3a, 0xfc, 0x40, 0x6f, 0xbd, 0x53, 0x7f, 0x30, 0xa7, 0x37, 0xc7, 0xe3, 0x2e, 0x08,
	0x46, 0xcf, 0x8d, 0xdc, 0xa3, 0xdb, 0x7a, 0x16, 0xb9, 0x5a, 0x8c, 0xea, 0x70, 0x40, 0xdf, 0x7b,
	0x27, 0x09, 0x01, 0x63, 0xd7, 0x25, 0x34, 0xb0, 0xd6, 0x8c, 0xef, 0xfb, 0xd5, 0xf8, 0x13, 0xbd,
	0x4f, 0xf4, 0x66, 0x68, 0x6c, 0xae, 0x55, 0x66, 0xe2, 0xb5, 0xed, 0xca, 0x54, 0x27, 0x6d, 0xb2,
	0xe3, 0xcb, 0x36, 0xbd, 0x3f, 0x84, 0x2c, 0x2f, 0xdd, 0xe4, 0xe9, 0x5f, 0x00, 0x00, 0x00, 0xff,
	0xff, 0xc4, 0x74, 0x88, 0x54, 0x38, 0x02, 0x00, 0x00,
}
