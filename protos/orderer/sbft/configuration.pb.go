// Code generated by protoc-gen-go. DO NOT EDIT.
// source: orderer/sbft/configuration.proto

package sbft // import "github.com/hyperledger/fabric/protos/orderer/sbft"

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

// ConfigMetadata is serialized and set as the value of ConsensusType.Metadata in
// a channel configuration when the ConsensusType.Type is set "etcdraft".
type ConfigMetadata struct {
	Consenters           []*Consenter `protobuf:"bytes,1,rep,name=consenters,proto3" json:"consenters,omitempty"`
	Options              *Options     `protobuf:"bytes,2,opt,name=options,proto3" json:"options,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *ConfigMetadata) Reset()         { *m = ConfigMetadata{} }
func (m *ConfigMetadata) String() string { return proto.CompactTextString(m) }
func (*ConfigMetadata) ProtoMessage()    {}
func (*ConfigMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_configuration_4241b6fb51698640, []int{0}
}
func (m *ConfigMetadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfigMetadata.Unmarshal(m, b)
}
func (m *ConfigMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfigMetadata.Marshal(b, m, deterministic)
}
func (dst *ConfigMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigMetadata.Merge(dst, src)
}
func (m *ConfigMetadata) XXX_Size() int {
	return xxx_messageInfo_ConfigMetadata.Size(m)
}
func (m *ConfigMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfigMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_ConfigMetadata proto.InternalMessageInfo

func (m *ConfigMetadata) GetConsenters() []*Consenter {
	if m != nil {
		return m.Consenters
	}
	return nil
}

func (m *ConfigMetadata) GetOptions() *Options {
	if m != nil {
		return m.Options
	}
	return nil
}

// Consenter represents a consenting node (i.e. replica).
type Consenter struct {
	Host                 string   `protobuf:"bytes,1,opt,name=host,proto3" json:"host,omitempty"`
	Port                 uint32   `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
	ClientSignCert       []byte   `protobuf:"bytes,3,opt,name=client_sign_cert,json=clientSignCert,proto3" json:"client_sign_cert,omitempty"`
	ServerTlsCert        []byte   `protobuf:"bytes,4,opt,name=server_tls_cert,json=serverTlsCert,proto3" json:"server_tls_cert,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Consenter) Reset()         { *m = Consenter{} }
func (m *Consenter) String() string { return proto.CompactTextString(m) }
func (*Consenter) ProtoMessage()    {}
func (*Consenter) Descriptor() ([]byte, []int) {
	return fileDescriptor_configuration_4241b6fb51698640, []int{1}
}
func (m *Consenter) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Consenter.Unmarshal(m, b)
}
func (m *Consenter) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Consenter.Marshal(b, m, deterministic)
}
func (dst *Consenter) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Consenter.Merge(dst, src)
}
func (m *Consenter) XXX_Size() int {
	return xxx_messageInfo_Consenter.Size(m)
}
func (m *Consenter) XXX_DiscardUnknown() {
	xxx_messageInfo_Consenter.DiscardUnknown(m)
}

var xxx_messageInfo_Consenter proto.InternalMessageInfo

func (m *Consenter) GetHost() string {
	if m != nil {
		return m.Host
	}
	return ""
}

func (m *Consenter) GetPort() uint32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *Consenter) GetClientSignCert() []byte {
	if m != nil {
		return m.ClientSignCert
	}
	return nil
}

func (m *Consenter) GetServerTlsCert() []byte {
	if m != nil {
		return m.ServerTlsCert
	}
	return nil
}

// Defines the SBFT parameters when 'sbft' is specified as the 'OrdererType'
type Options struct {
	N                    uint64   `protobuf:"varint,1,opt,name=n,proto3" json:"n,omitempty"`
	F                    uint64   `protobuf:"varint,2,opt,name=f,proto3" json:"f,omitempty"`
	RequestTimeoutNsec   uint64   `protobuf:"varint,3,opt,name=request_timeout_nsec,json=requestTimeoutNsec,proto3" json:"request_timeout_nsec,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Options) Reset()         { *m = Options{} }
func (m *Options) String() string { return proto.CompactTextString(m) }
func (*Options) ProtoMessage()    {}
func (*Options) Descriptor() ([]byte, []int) {
	return fileDescriptor_configuration_4241b6fb51698640, []int{2}
}
func (m *Options) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Options.Unmarshal(m, b)
}
func (m *Options) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Options.Marshal(b, m, deterministic)
}
func (dst *Options) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Options.Merge(dst, src)
}
func (m *Options) XXX_Size() int {
	return xxx_messageInfo_Options.Size(m)
}
func (m *Options) XXX_DiscardUnknown() {
	xxx_messageInfo_Options.DiscardUnknown(m)
}

var xxx_messageInfo_Options proto.InternalMessageInfo

func (m *Options) GetN() uint64 {
	if m != nil {
		return m.N
	}
	return 0
}

func (m *Options) GetF() uint64 {
	if m != nil {
		return m.F
	}
	return 0
}

func (m *Options) GetRequestTimeoutNsec() uint64 {
	if m != nil {
		return m.RequestTimeoutNsec
	}
	return 0
}

type ConsensusConfig struct {
	Consensus            *Options              `protobuf:"bytes,1,opt,name=consensus,proto3" json:"consensus,omitempty"`
	Peers                map[string]*Consenter `protobuf:"bytes,2,rep,name=peers,proto3" json:"peers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *ConsensusConfig) Reset()         { *m = ConsensusConfig{} }
func (m *ConsensusConfig) String() string { return proto.CompactTextString(m) }
func (*ConsensusConfig) ProtoMessage()    {}
func (*ConsensusConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_configuration_4241b6fb51698640, []int{3}
}
func (m *ConsensusConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConsensusConfig.Unmarshal(m, b)
}
func (m *ConsensusConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConsensusConfig.Marshal(b, m, deterministic)
}
func (dst *ConsensusConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConsensusConfig.Merge(dst, src)
}
func (m *ConsensusConfig) XXX_Size() int {
	return xxx_messageInfo_ConsensusConfig.Size(m)
}
func (m *ConsensusConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_ConsensusConfig.DiscardUnknown(m)
}

var xxx_messageInfo_ConsensusConfig proto.InternalMessageInfo

func (m *ConsensusConfig) GetConsensus() *Options {
	if m != nil {
		return m.Consensus
	}
	return nil
}

func (m *ConsensusConfig) GetPeers() map[string]*Consenter {
	if m != nil {
		return m.Peers
	}
	return nil
}

type JsonConfig struct {
	Consensus            *Options `protobuf:"bytes,1,opt,name=consensus,proto3" json:"consensus,omitempty"`
	Peers                []*Peer  `protobuf:"bytes,2,rep,name=peers,proto3" json:"peers,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *JsonConfig) Reset()         { *m = JsonConfig{} }
func (m *JsonConfig) String() string { return proto.CompactTextString(m) }
func (*JsonConfig) ProtoMessage()    {}
func (*JsonConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_configuration_4241b6fb51698640, []int{4}
}
func (m *JsonConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JsonConfig.Unmarshal(m, b)
}
func (m *JsonConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JsonConfig.Marshal(b, m, deterministic)
}
func (dst *JsonConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JsonConfig.Merge(dst, src)
}
func (m *JsonConfig) XXX_Size() int {
	return xxx_messageInfo_JsonConfig.Size(m)
}
func (m *JsonConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_JsonConfig.DiscardUnknown(m)
}

var xxx_messageInfo_JsonConfig proto.InternalMessageInfo

func (m *JsonConfig) GetConsensus() *Options {
	if m != nil {
		return m.Consensus
	}
	return nil
}

func (m *JsonConfig) GetPeers() []*Peer {
	if m != nil {
		return m.Peers
	}
	return nil
}

type Peer struct {
	Address              string   `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Cert                 string   `protobuf:"bytes,2,opt,name=cert,proto3" json:"cert,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Peer) Reset()         { *m = Peer{} }
func (m *Peer) String() string { return proto.CompactTextString(m) }
func (*Peer) ProtoMessage()    {}
func (*Peer) Descriptor() ([]byte, []int) {
	return fileDescriptor_configuration_4241b6fb51698640, []int{5}
}
func (m *Peer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Peer.Unmarshal(m, b)
}
func (m *Peer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Peer.Marshal(b, m, deterministic)
}
func (dst *Peer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Peer.Merge(dst, src)
}
func (m *Peer) XXX_Size() int {
	return xxx_messageInfo_Peer.Size(m)
}
func (m *Peer) XXX_DiscardUnknown() {
	xxx_messageInfo_Peer.DiscardUnknown(m)
}

var xxx_messageInfo_Peer proto.InternalMessageInfo

func (m *Peer) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *Peer) GetCert() string {
	if m != nil {
		return m.Cert
	}
	return ""
}

func init() {
	proto.RegisterType((*ConfigMetadata)(nil), "sbft.ConfigMetadata")
	proto.RegisterType((*Consenter)(nil), "sbft.Consenter")
	proto.RegisterType((*Options)(nil), "sbft.Options")
	proto.RegisterType((*ConsensusConfig)(nil), "sbft.consensus_config")
	proto.RegisterMapType((map[string]*Consenter)(nil), "sbft.consensus_config.PeersEntry")
	proto.RegisterType((*JsonConfig)(nil), "sbft.json_config")
	proto.RegisterType((*Peer)(nil), "sbft.peer")
}

func init() {
	proto.RegisterFile("orderer/sbft/configuration.proto", fileDescriptor_configuration_4241b6fb51698640)
}

var fileDescriptor_configuration_4241b6fb51698640 = []byte{
	// 458 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x53, 0x4d, 0x8b, 0xd4, 0x40,
	0x10, 0xa5, 0x67, 0xb2, 0x0e, 0x53, 0xb3, 0xb3, 0x3b, 0x34, 0x1e, 0x82, 0xa7, 0x18, 0x50, 0x83,
	0x42, 0xa2, 0xab, 0xa0, 0x78, 0x74, 0xf0, 0xe0, 0xc1, 0x0f, 0xda, 0x05, 0x41, 0x84, 0x90, 0x49,
	0x2a, 0x99, 0xac, 0xd9, 0xee, 0xd8, 0xd5, 0x59, 0x98, 0x1f, 0xe0, 0x2f, 0xf3, 0x8f, 0x49, 0x77,
	0x67, 0x76, 0x47, 0xf7, 0xb6, 0xb7, 0xca, 0x7b, 0xaf, 0xab, 0x1e, 0xaf, 0x2a, 0x10, 0x29, 0x5d,
	0xa1, 0x46, 0x9d, 0xd1, 0xa6, 0x36, 0x59, 0xa9, 0x64, 0xdd, 0x36, 0x83, 0x2e, 0x4c, 0xab, 0x64,
	0xda, 0x6b, 0x65, 0x14, 0x0f, 0x2c, 0x13, 0x5f, 0xc0, 0xc9, 0xda, 0x91, 0x1f, 0xd1, 0x14, 0x55,
	0x61, 0x0a, 0x9e, 0x01, 0x94, 0x4a, 0x12, 0x4a, 0x83, 0x9a, 0x42, 0x16, 0x4d, 0x93, 0xc5, 0xd9,
	0x69, 0x6a, 0xc5, 0xe9, 0x7a, 0x8f, 0x8b, 0x03, 0x09, 0x7f, 0x02, 0x33, 0xd5, 0xdb, 0xc6, 0x14,
	0x4e, 0x22, 0x96, 0x2c, 0xce, 0x96, 0x5e, 0xfd, 0xd9, 0x83, 0x62, 0xcf, 0xc6, 0xbf, 0x19, 0xcc,
	0xaf, 0x5b, 0x70, 0x0e, 0xc1, 0x56, 0x91, 0x09, 0x59, 0xc4, 0x92, 0xb9, 0x70, 0xb5, 0xc5, 0x7a,
	0xa5, 0x8d, 0xeb, 0xb3, 0x14, 0xae, 0xe6, 0x09, 0xac, 0xca, 0xae, 0x45, 0x69, 0x72, 0x6a, 0x1b,
	0x99, 0x97, 0xa8, 0x4d, 0x38, 0x8d, 0x58, 0x72, 0x2c, 0x4e, 0x3c, 0xfe, 0xb5, 0x6d, 0xe4, 0x1a,
	0xb5, 0xe1, 0x8f, 0xe1, 0x94, 0x50, 0x5f, 0xa1, 0xce, 0x4d, 0x47, 0x5e, 0x18, 0x38, 0xe1, 0xd2,
	0xc3, 0xe7, 0x1d, 0x59, 0x5d, 0xfc, 0x0d, 0x66, 0xa3, 0x37, 0x7e, 0x0c, 0x4c, 0x3a, 0x07, 0x81,
	0x60, 0xd2, 0x7e, 0xd5, 0x6e, 0x76, 0x20, 0x58, 0xcd, 0x9f, 0xc3, 0x7d, 0x8d, 0xbf, 0x06, 0x24,
	0x93, 0x9b, 0xf6, 0x12, 0xd5, 0x60, 0x72, 0x49, 0x58, 0xba, 0xe1, 0x81, 0xe0, 0x23, 0x77, 0xee,
	0xa9, 0x4f, 0x84, 0x65, 0xfc, 0x87, 0xc1, 0xca, 0x07, 0x43, 0x03, 0xe5, 0x3e, 0x74, 0xfe, 0x0c,
	0xe6, 0xd7, 0x98, 0x1b, 0x75, 0x2b, 0xa0, 0x1b, 0x9e, 0xbf, 0x86, 0xa3, 0x1e, 0x6d, 0xee, 0x13,
	0x97, 0xfb, 0x43, 0x2f, 0xfc, 0xbf, 0x67, 0xfa, 0xc5, 0x6a, 0xde, 0x4b, 0xa3, 0x77, 0xc2, 0xeb,
	0x1f, 0x7c, 0x00, 0xb8, 0x01, 0xf9, 0x0a, 0xa6, 0x3f, 0x71, 0x37, 0x46, 0x6b, 0x4b, 0xfe, 0x08,
	0x8e, 0xae, 0x8a, 0x6e, 0xc0, 0x71, 0x45, 0xb7, 0x16, 0xea, 0xd9, 0xb7, 0x93, 0x37, 0x2c, 0xfe,
	0x01, 0x8b, 0x0b, 0x52, 0xf2, 0x4e, 0xfe, 0xa3, 0x7f, 0xfd, 0x83, 0x17, 0x5a, 0x68, 0x34, 0x1a,
	0xbf, 0x82, 0xc0, 0x16, 0x3c, 0x84, 0x59, 0x51, 0x55, 0x1a, 0x89, 0x46, 0x9b, 0xfb, 0x4f, 0x7b,
	0x04, 0x6e, 0x77, 0x13, 0x7f, 0x18, 0xb6, 0x7e, 0x97, 0xc3, 0x53, 0xa5, 0x9b, 0x74, 0xbb, 0xeb,
	0x51, 0x77, 0x58, 0x35, 0xa8, 0xd3, 0xba, 0xd8, 0xe8, 0xb6, 0xf4, 0xc7, 0x4c, 0xe9, 0x78, 0xee,
	0x6e, 0xde, 0xf7, 0x17, 0x4d, 0x6b, 0xb6, 0xc3, 0x26, 0x2d, 0xd5, 0x65, 0x76, 0xf0, 0x24, 0xf3,
	0x4f, 0x32, 0xff, 0x24, 0x3b, 0xfc, 0x43, 0x36, 0xf7, 0x1c, 0xf8, 0xf2, 0x6f, 0x00, 0x00, 0x00,
	0xff, 0xff, 0x8b, 0xcb, 0x41, 0x1c, 0x38, 0x03, 0x00, 0x00,
}
