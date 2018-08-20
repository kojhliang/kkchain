// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: p2p/handshake/pb/handshake.proto

package handshake

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"
import _ "github.com/gogo/protobuf/types"

import bytes "bytes"

import strings "strings"
import reflect "reflect"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type Message_Type int32

const (
	Message_HELLO       Message_Type = 0
	Message_HELLO_OK    Message_Type = 1
	Message_HELLO_ERROR Message_Type = 2
)

var Message_Type_name = map[int32]string{
	0: "HELLO",
	1: "HELLO_OK",
	2: "HELLO_ERROR",
}
var Message_Type_value = map[string]int32{
	"HELLO":       0,
	"HELLO_OK":    1,
	"HELLO_ERROR": 2,
}

func (x Message_Type) String() string {
	return proto.EnumName(Message_Type_name, int32(x))
}
func (Message_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_handshake_735fac6794902d4f, []int{1, 0}
}

type Error struct {
	// error code
	Code int32 `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	// description
	Desc                 string   `protobuf:"bytes,2,opt,name=desc,proto3" json:"desc,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Error) Reset()      { *m = Error{} }
func (*Error) ProtoMessage() {}
func (*Error) Descriptor() ([]byte, []int) {
	return fileDescriptor_handshake_735fac6794902d4f, []int{0}
}
func (m *Error) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Error) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Error.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *Error) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Error.Merge(dst, src)
}
func (m *Error) XXX_Size() int {
	return m.Size()
}
func (m *Error) XXX_DiscardUnknown() {
	xxx_messageInfo_Error.DiscardUnknown(m)
}

var xxx_messageInfo_Error proto.InternalMessageInfo

func (m *Error) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *Error) GetDesc() string {
	if m != nil {
		return m.Desc
	}
	return ""
}

type Message struct {
	// defines what type of message it is.
	Type Message_Type `protobuf:"varint,1,opt,name=type,proto3,enum=pb.Message_Type" json:"type,omitempty"`
	// protodcol version
	ProtocolVersion string `protobuf:"bytes,2,opt,name=protocol_version,json=protocolVersion,proto3" json:"protocol_version,omitempty"`
	// chain ID
	ChainId []byte `protobuf:"bytes,3,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	// genesis hash
	GenesisHash []byte `protobuf:"bytes,4,opt,name=genesis_hash,json=genesisHash,proto3" json:"genesis_hash,omitempty"`
	// head hash
	HeadHash []byte `protobuf:"bytes,5,opt,name=head_hash,json=headHash,proto3" json:"head_hash,omitempty"`
	// head num
	HeadNum uint64 `protobuf:"varint,6,opt,name=head_num,json=headNum,proto3" json:"head_num,omitempty"`
	// error info for HELLO_ERROR
	Error                *Error   `protobuf:"bytes,7,opt,name=error" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Message) Reset()      { *m = Message{} }
func (*Message) ProtoMessage() {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_handshake_735fac6794902d4f, []int{1}
}
func (m *Message) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Message.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(dst, src)
}
func (m *Message) XXX_Size() int {
	return m.Size()
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetType() Message_Type {
	if m != nil {
		return m.Type
	}
	return Message_HELLO
}

func (m *Message) GetProtocolVersion() string {
	if m != nil {
		return m.ProtocolVersion
	}
	return ""
}

func (m *Message) GetChainId() []byte {
	if m != nil {
		return m.ChainId
	}
	return nil
}

func (m *Message) GetGenesisHash() []byte {
	if m != nil {
		return m.GenesisHash
	}
	return nil
}

func (m *Message) GetHeadHash() []byte {
	if m != nil {
		return m.HeadHash
	}
	return nil
}

func (m *Message) GetHeadNum() uint64 {
	if m != nil {
		return m.HeadNum
	}
	return 0
}

func (m *Message) GetError() *Error {
	if m != nil {
		return m.Error
	}
	return nil
}

func init() {
	proto.RegisterType((*Error)(nil), "handshake.Error")
	proto.RegisterType((*Message)(nil), "handshake.Message")
	proto.RegisterEnum("handshake.Message_Type", Message_Type_name, Message_Type_value)
}
func (this *Error) VerboseEqual(that interface{}) error {
	if that == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that == nil && this != nil")
	}

	that1, ok := that.(*Error)
	if !ok {
		that2, ok := that.(Error)
		if ok {
			that1 = &that2
		} else {
			return fmt.Errorf("that is not of type *Error")
		}
	}
	if that1 == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that is type *Error but is nil && this != nil")
	} else if this == nil {
		return fmt.Errorf("that is type *Error but is not nil && this == nil")
	}
	if this.Code != that1.Code {
		return fmt.Errorf("Code this(%v) Not Equal that(%v)", this.Code, that1.Code)
	}
	if this.Desc != that1.Desc {
		return fmt.Errorf("Desc this(%v) Not Equal that(%v)", this.Desc, that1.Desc)
	}
	return nil
}
func (this *Error) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Error)
	if !ok {
		that2, ok := that.(Error)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Code != that1.Code {
		return false
	}
	if this.Desc != that1.Desc {
		return false
	}
	return true
}
func (this *Message) VerboseEqual(that interface{}) error {
	if that == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that == nil && this != nil")
	}

	that1, ok := that.(*Message)
	if !ok {
		that2, ok := that.(Message)
		if ok {
			that1 = &that2
		} else {
			return fmt.Errorf("that is not of type *Message")
		}
	}
	if that1 == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that is type *Message but is nil && this != nil")
	} else if this == nil {
		return fmt.Errorf("that is type *Message but is not nil && this == nil")
	}
	if this.Type != that1.Type {
		return fmt.Errorf("Type this(%v) Not Equal that(%v)", this.Type, that1.Type)
	}
	if this.ProtocolVersion != that1.ProtocolVersion {
		return fmt.Errorf("ProtocolVersion this(%v) Not Equal that(%v)", this.ProtocolVersion, that1.ProtocolVersion)
	}
	if !bytes.Equal(this.ChainId, that1.ChainId) {
		return fmt.Errorf("ChainId this(%v) Not Equal that(%v)", this.ChainId, that1.ChainId)
	}
	if !bytes.Equal(this.GenesisHash, that1.GenesisHash) {
		return fmt.Errorf("GenesisHash this(%v) Not Equal that(%v)", this.GenesisHash, that1.GenesisHash)
	}
	if !bytes.Equal(this.HeadHash, that1.HeadHash) {
		return fmt.Errorf("HeadHash this(%v) Not Equal that(%v)", this.HeadHash, that1.HeadHash)
	}
	if this.HeadNum != that1.HeadNum {
		return fmt.Errorf("HeadNum this(%v) Not Equal that(%v)", this.HeadNum, that1.HeadNum)
	}
	if !this.Error.Equal(that1.Error) {
		return fmt.Errorf("Error this(%v) Not Equal that(%v)", this.Error, that1.Error)
	}
	return nil
}
func (this *Message) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Message)
	if !ok {
		that2, ok := that.(Message)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Type != that1.Type {
		return false
	}
	if this.ProtocolVersion != that1.ProtocolVersion {
		return false
	}
	if !bytes.Equal(this.ChainId, that1.ChainId) {
		return false
	}
	if !bytes.Equal(this.GenesisHash, that1.GenesisHash) {
		return false
	}
	if !bytes.Equal(this.HeadHash, that1.HeadHash) {
		return false
	}
	if this.HeadNum != that1.HeadNum {
		return false
	}
	if !this.Error.Equal(that1.Error) {
		return false
	}
	return true
}
func (this *Error) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 6)
	s = append(s, "&pb.Error{")
	s = append(s, "Code: "+fmt.Sprintf("%#v", this.Code)+",\n")
	s = append(s, "Desc: "+fmt.Sprintf("%#v", this.Desc)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *Message) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 11)
	s = append(s, "&pb.Message{")
	s = append(s, "Type: "+fmt.Sprintf("%#v", this.Type)+",\n")
	s = append(s, "ProtocolVersion: "+fmt.Sprintf("%#v", this.ProtocolVersion)+",\n")
	s = append(s, "ChainId: "+fmt.Sprintf("%#v", this.ChainId)+",\n")
	s = append(s, "GenesisHash: "+fmt.Sprintf("%#v", this.GenesisHash)+",\n")
	s = append(s, "HeadHash: "+fmt.Sprintf("%#v", this.HeadHash)+",\n")
	s = append(s, "HeadNum: "+fmt.Sprintf("%#v", this.HeadNum)+",\n")
	if this.Error != nil {
		s = append(s, "Error: "+fmt.Sprintf("%#v", this.Error)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringHandshake(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func (m *Error) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Error) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Code != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintHandshake(dAtA, i, uint64(m.Code))
	}
	if len(m.Desc) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintHandshake(dAtA, i, uint64(len(m.Desc)))
		i += copy(dAtA[i:], m.Desc)
	}
	return i, nil
}

func (m *Message) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Message) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Type != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintHandshake(dAtA, i, uint64(m.Type))
	}
	if len(m.ProtocolVersion) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintHandshake(dAtA, i, uint64(len(m.ProtocolVersion)))
		i += copy(dAtA[i:], m.ProtocolVersion)
	}
	if len(m.ChainId) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintHandshake(dAtA, i, uint64(len(m.ChainId)))
		i += copy(dAtA[i:], m.ChainId)
	}
	if len(m.GenesisHash) > 0 {
		dAtA[i] = 0x22
		i++
		i = encodeVarintHandshake(dAtA, i, uint64(len(m.GenesisHash)))
		i += copy(dAtA[i:], m.GenesisHash)
	}
	if len(m.HeadHash) > 0 {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintHandshake(dAtA, i, uint64(len(m.HeadHash)))
		i += copy(dAtA[i:], m.HeadHash)
	}
	if m.HeadNum != 0 {
		dAtA[i] = 0x30
		i++
		i = encodeVarintHandshake(dAtA, i, uint64(m.HeadNum))
	}
	if m.Error != nil {
		dAtA[i] = 0x3a
		i++
		i = encodeVarintHandshake(dAtA, i, uint64(m.Error.Size()))
		n1, err := m.Error.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	return i, nil
}

func encodeVarintHandshake(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *Error) Size() (n int) {
	var l int
	_ = l
	if m.Code != 0 {
		n += 1 + sovHandshake(uint64(m.Code))
	}
	l = len(m.Desc)
	if l > 0 {
		n += 1 + l + sovHandshake(uint64(l))
	}
	return n
}

func (m *Message) Size() (n int) {
	var l int
	_ = l
	if m.Type != 0 {
		n += 1 + sovHandshake(uint64(m.Type))
	}
	l = len(m.ProtocolVersion)
	if l > 0 {
		n += 1 + l + sovHandshake(uint64(l))
	}
	l = len(m.ChainId)
	if l > 0 {
		n += 1 + l + sovHandshake(uint64(l))
	}
	l = len(m.GenesisHash)
	if l > 0 {
		n += 1 + l + sovHandshake(uint64(l))
	}
	l = len(m.HeadHash)
	if l > 0 {
		n += 1 + l + sovHandshake(uint64(l))
	}
	if m.HeadNum != 0 {
		n += 1 + sovHandshake(uint64(m.HeadNum))
	}
	if m.Error != nil {
		l = m.Error.Size()
		n += 1 + l + sovHandshake(uint64(l))
	}
	return n
}

func sovHandshake(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozHandshake(x uint64) (n int) {
	return sovHandshake(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *Error) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Error{`,
		`Code:` + fmt.Sprintf("%v", this.Code) + `,`,
		`Desc:` + fmt.Sprintf("%v", this.Desc) + `,`,
		`}`,
	}, "")
	return s
}
func (this *Message) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Message{`,
		`Type:` + fmt.Sprintf("%v", this.Type) + `,`,
		`ProtocolVersion:` + fmt.Sprintf("%v", this.ProtocolVersion) + `,`,
		`ChainId:` + fmt.Sprintf("%v", this.ChainId) + `,`,
		`GenesisHash:` + fmt.Sprintf("%v", this.GenesisHash) + `,`,
		`HeadHash:` + fmt.Sprintf("%v", this.HeadHash) + `,`,
		`HeadNum:` + fmt.Sprintf("%v", this.HeadNum) + `,`,
		`Error:` + strings.Replace(fmt.Sprintf("%v", this.Error), "Error", "Error", 1) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringHandshake(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *Error) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowHandshake
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Error: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Error: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Code", wireType)
			}
			m.Code = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHandshake
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Code |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Desc", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHandshake
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthHandshake
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Desc = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipHandshake(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthHandshake
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Message) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowHandshake
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Message: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Message: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			m.Type = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHandshake
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Type |= (Message_Type(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProtocolVersion", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHandshake
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthHandshake
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ProtocolVersion = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChainId", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHandshake
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthHandshake
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ChainId = append(m.ChainId[:0], dAtA[iNdEx:postIndex]...)
			if m.ChainId == nil {
				m.ChainId = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GenesisHash", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHandshake
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthHandshake
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.GenesisHash = append(m.GenesisHash[:0], dAtA[iNdEx:postIndex]...)
			if m.GenesisHash == nil {
				m.GenesisHash = []byte{}
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field HeadHash", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHandshake
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthHandshake
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.HeadHash = append(m.HeadHash[:0], dAtA[iNdEx:postIndex]...)
			if m.HeadHash == nil {
				m.HeadHash = []byte{}
			}
			iNdEx = postIndex
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field HeadNum", wireType)
			}
			m.HeadNum = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHandshake
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.HeadNum |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Error", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHandshake
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthHandshake
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Error == nil {
				m.Error = &Error{}
			}
			if err := m.Error.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipHandshake(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthHandshake
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipHandshake(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowHandshake
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowHandshake
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowHandshake
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthHandshake
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowHandshake
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipHandshake(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthHandshake = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowHandshake   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("p2p/handshake/pb/handshake.proto", fileDescriptor_handshake_735fac6794902d4f)
}

var fileDescriptor_handshake_735fac6794902d4f = []byte{
	// 401 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x44, 0x91, 0x4f, 0x6e, 0xd3, 0x40,
	0x14, 0xc6, 0x33, 0xc1, 0x6e, 0x92, 0x97, 0x88, 0x5a, 0xb3, 0x72, 0x8a, 0x34, 0xb8, 0x11, 0x0b,
	0xb3, 0xb1, 0x51, 0xb8, 0x41, 0x51, 0xa4, 0xa2, 0x14, 0x52, 0x8d, 0x10, 0x2c, 0x2d, 0xff, 0x19,
	0x6c, 0x2b, 0xed, 0xcc, 0xc8, 0xd3, 0x54, 0xca, 0x8e, 0x23, 0x70, 0x07, 0x36, 0x1c, 0x85, 0x25,
	0x4b, 0x96, 0x8d, 0xb9, 0x00, 0x47, 0x40, 0xf3, 0x26, 0xd0, 0xdd, 0xf7, 0xfd, 0x7e, 0x63, 0x8f,
	0xde, 0x3c, 0x88, 0xf4, 0x52, 0xa7, 0x4d, 0x2e, 0x2b, 0xd3, 0xe4, 0x5b, 0x91, 0xea, 0xe2, 0xb1,
	0x24, 0xba, 0x53, 0x77, 0x8a, 0x0e, 0x75, 0x71, 0x36, 0xaf, 0x95, 0xaa, 0x6f, 0x44, 0x8a, 0xa4,
	0xd8, 0x7d, 0x4e, 0x73, 0xb9, 0x77, 0xfa, 0x6c, 0x51, 0xab, 0x5a, 0x3d, 0x0a, 0xdb, 0xb0, 0x60,
	0x72, 0x67, 0x16, 0x29, 0xf8, 0xab, 0xae, 0x53, 0x1d, 0xa5, 0xe0, 0x95, 0xaa, 0x12, 0x21, 0x89,
	0x48, 0xec, 0x73, 0xcc, 0x96, 0x55, 0xc2, 0x94, 0xe1, 0x30, 0x22, 0xf1, 0x84, 0x63, 0x5e, 0x7c,
	0x1b, 0xc2, 0xe8, 0x9d, 0x30, 0x26, 0xaf, 0x05, 0x7d, 0x01, 0xde, 0xdd, 0x5e, 0xbb, 0x6f, 0x9e,
	0x2e, 0x83, 0x44, 0x17, 0xc9, 0x51, 0x25, 0x1f, 0xf6, 0x5a, 0x70, 0xb4, 0xf4, 0x25, 0x04, 0x78,
	0x57, 0xa9, 0x6e, 0xb2, 0x7b, 0xd1, 0x99, 0x56, 0xc9, 0xe3, 0x1f, 0x4f, 0xff, 0xf1, 0x8f, 0x0e,
	0xd3, 0x39, 0x8c, 0xcb, 0x26, 0x6f, 0x65, 0xd6, 0x56, 0xe1, 0x93, 0x88, 0xc4, 0x33, 0x3e, 0xc2,
	0xfe, 0xb6, 0xa2, 0xe7, 0x30, 0xab, 0x85, 0x14, 0xa6, 0x35, 0x59, 0x93, 0x9b, 0x26, 0xf4, 0x50,
	0x4f, 0x8f, 0xec, 0x32, 0x37, 0x0d, 0x7d, 0x06, 0x93, 0x46, 0xe4, 0x95, 0xf3, 0x3e, 0xfa, 0xb1,
	0x05, 0x28, 0xe7, 0x80, 0x39, 0x93, 0xbb, 0xdb, 0xf0, 0x24, 0x22, 0xb1, 0xc7, 0x47, 0xb6, 0xbf,
	0xdf, 0xdd, 0xd2, 0xe7, 0xe0, 0x0b, 0xfb, 0x06, 0xe1, 0x28, 0x22, 0xf1, 0x74, 0x39, 0xb1, 0x73,
	0xe0, 0xa3, 0x70, 0xc7, 0x17, 0xaf, 0xc0, 0xb3, 0xf3, 0xd0, 0x09, 0xf8, 0x97, 0xab, 0xab, 0xab,
	0x4d, 0x30, 0xa0, 0x33, 0x18, 0x63, 0xcc, 0x36, 0xeb, 0x80, 0xd0, 0x53, 0x98, 0xba, 0xb6, 0xe2,
	0x7c, 0xc3, 0x83, 0xe1, 0xc5, 0xa7, 0x5f, 0x07, 0x36, 0x78, 0x38, 0x30, 0xf2, 0xe7, 0xc0, 0xc8,
	0x97, 0x9e, 0x91, 0xef, 0x3d, 0x23, 0x3f, 0x7a, 0x46, 0x7e, 0xf6, 0x8c, 0x3c, 0xf4, 0x8c, 0x7c,
	0xfd, 0xcd, 0x06, 0x70, 0xae, 0xba, 0x3a, 0x69, 0xe5, 0x7d, 0x2b, 0x93, 0xed, 0x16, 0x07, 0x75,
	0x5b, 0x49, 0xfe, 0x2f, 0xfa, 0x62, 0xb6, 0x5e, 0xbf, 0xb1, 0xe2, 0xda, 0xf2, 0x6b, 0x52, 0x9c,
	0xe0, 0x81, 0xd7, 0x7f, 0x03, 0x00, 0x00, 0xff, 0xff, 0x85, 0xb8, 0x1f, 0x2b, 0x1d, 0x02, 0x00,
	0x00,
}
