// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: kv/kvserver/closedts/ctpb/service.proto

package ctpb

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import roachpb "github.com/cockroachdb/cockroach/pkg/roachpb"
import hlc "github.com/cockroachdb/cockroach/pkg/util/hlc"

import github_com_cockroachdb_cockroach_pkg_roachpb "github.com/cockroachdb/cockroach/pkg/roachpb"

import (
	context "context"
	grpc "google.golang.org/grpc"
)

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

// Update contains information about (the advancement of) closed timestamps for
// ranges with leases on the sender node. Updates are of two types: snapshots
// and incrementals. Snapshots are stand-alone messages, explicitly containing
// state about a bunch of ranges. Incrementals are deltas since the previous
// message (the previous message can be a snapshot or another incremental); they
// contain info about which new ranges are included in the info provided, which
// ranges are removed, and how the closed timestamps for different categories of
// ranges advanced. Ranges communicated by a previous message and not touched by
// an incremental are "implicitly" referenced by the incremental. In order to
// properly handle incrementals, the recipient maintains a "stream's state": the
// group of ranges that can be implicitly referenced by the next message.
type Update struct {
	// node_id identifies the sending node.
	NodeID github_com_cockroachdb_cockroach_pkg_roachpb.NodeID `protobuf:"varint,1,opt,name=node_id,json=nodeId,proto3,casttype=github.com/cockroachdb/cockroach/pkg/roachpb.NodeID" json:"node_id,omitempty"`
	// seq_num identifies this update across all updates produced by a node. The
	// sequence is reset when the node restarts, so a recipient can only count on
	// it increasing within a single PushUpdates stream.
	//
	// All messages have sequence numbers, including snapshots. A snapshot can be
	// applied on top of any state (i.e. it can be applied after having skipped
	// messages); its sequence number tells the recipient what incremental message
	// it should expect afterwards.
	SeqNum SeqNum `protobuf:"varint,2,opt,name=seq_num,json=seqNum,proto3,casttype=SeqNum" json:"seq_num,omitempty"`
	// snapshot indicates whether this message is standalone, or whether it's just
	// a delta since the messages with the previous seq_num. A snapshot
	// re-initializes all of the recipient's state. The first message on a stream
	// is always a snapshot. Afterwards, there could be others if the sender is
	// temporarily slowed down or if the stream experience network problems and
	// some incremental messages are dropped  (although generally we expect that
	// to result in a stream failing and a new one being established).
	Snapshot         bool                 `protobuf:"varint,3,opt,name=snapshot,proto3" json:"snapshot,omitempty"`
	ClosedTimestamps []Update_GroupUpdate `protobuf:"bytes,4,rep,name=closed_timestamps,json=closedTimestamps,proto3" json:"closed_timestamps"`
	// removed contains the set of ranges that are no longer registered on the
	// stream and who future updates are no longer applicable to.
	//
	// The field will be empty if snapshot is true, as a snapshot message implies
	// that all ranges not present in the snapshot's added_or_updated list are no
	// longer tracked.
	Removed        []github_com_cockroachdb_cockroach_pkg_roachpb.RangeID `protobuf:"varint,5,rep,packed,name=removed,proto3,casttype=github.com/cockroachdb/cockroach/pkg/roachpb.RangeID" json:"removed,omitempty"`
	AddedOrUpdated []Update_RangeUpdate                                   `protobuf:"bytes,6,rep,name=added_or_updated,json=addedOrUpdated,proto3" json:"added_or_updated"`
}

func (m *Update) Reset()         { *m = Update{} }
func (m *Update) String() string { return proto.CompactTextString(m) }
func (*Update) ProtoMessage()    {}
func (*Update) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_96a1a4bff833e11e, []int{0}
}
func (m *Update) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Update) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalTo(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (dst *Update) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Update.Merge(dst, src)
}
func (m *Update) XXX_Size() int {
	return m.Size()
}
func (m *Update) XXX_DiscardUnknown() {
	xxx_messageInfo_Update.DiscardUnknown(m)
}

var xxx_messageInfo_Update proto.InternalMessageInfo

// closed_timestamps represents the timestamps that are being closed for each
// group of ranges, with a group being represented by its policy.
//
// The recipient is supposed to forward the closed timestamps of the affected
// ranges to these values. Upon receiving one of these updates, the recipient
// should generally not assume that it hasn't been informed of a higher closed
// timestamp for any range in particular - races between this side-transport
// and the regular Raft transport are possible, as are races between two
// side-transport streams for an outgoing and incoming leaseholder.
type Update_GroupUpdate struct {
	Policy          roachpb.RangeClosedTimestampPolicy `protobuf:"varint,1,opt,name=policy,proto3,enum=cockroach.roachpb.RangeClosedTimestampPolicy" json:"policy,omitempty"`
	ClosedTimestamp hlc.Timestamp                      `protobuf:"bytes,2,opt,name=closed_timestamp,json=closedTimestamp,proto3" json:"closed_timestamp"`
}

func (m *Update_GroupUpdate) Reset()         { *m = Update_GroupUpdate{} }
func (m *Update_GroupUpdate) String() string { return proto.CompactTextString(m) }
func (*Update_GroupUpdate) ProtoMessage()    {}
func (*Update_GroupUpdate) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_96a1a4bff833e11e, []int{0, 0}
}
func (m *Update_GroupUpdate) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Update_GroupUpdate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalTo(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (dst *Update_GroupUpdate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Update_GroupUpdate.Merge(dst, src)
}
func (m *Update_GroupUpdate) XXX_Size() int {
	return m.Size()
}
func (m *Update_GroupUpdate) XXX_DiscardUnknown() {
	xxx_messageInfo_Update_GroupUpdate.DiscardUnknown(m)
}

var xxx_messageInfo_Update_GroupUpdate proto.InternalMessageInfo

// added_or_updated contains the set of ranges that are either being added to
// the tracked ranges set with a given (lai, policy) or updated within the
// tracked range set with a new (lai, policy). All future updates on the
// stream are applicable to these ranges until they are removed, either
// explicitly by being included in a future removed set or implicitly by not
// being included in the added_or_updated field of a future snapshot.
type Update_RangeUpdate struct {
	RangeID github_com_cockroachdb_cockroach_pkg_roachpb.RangeID `protobuf:"varint,1,opt,name=range_id,json=rangeId,proto3,casttype=github.com/cockroachdb/cockroach/pkg/roachpb.RangeID" json:"range_id,omitempty"`
	LAI     LAI                                                  `protobuf:"varint,2,opt,name=lai,proto3,casttype=LAI" json:"lai,omitempty"`
	Policy  roachpb.RangeClosedTimestampPolicy                   `protobuf:"varint,3,opt,name=policy,proto3,enum=cockroach.roachpb.RangeClosedTimestampPolicy" json:"policy,omitempty"`
}

func (m *Update_RangeUpdate) Reset()         { *m = Update_RangeUpdate{} }
func (m *Update_RangeUpdate) String() string { return proto.CompactTextString(m) }
func (*Update_RangeUpdate) ProtoMessage()    {}
func (*Update_RangeUpdate) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_96a1a4bff833e11e, []int{0, 1}
}
func (m *Update_RangeUpdate) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Update_RangeUpdate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalTo(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (dst *Update_RangeUpdate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Update_RangeUpdate.Merge(dst, src)
}
func (m *Update_RangeUpdate) XXX_Size() int {
	return m.Size()
}
func (m *Update_RangeUpdate) XXX_DiscardUnknown() {
	xxx_messageInfo_Update_RangeUpdate.DiscardUnknown(m)
}

var xxx_messageInfo_Update_RangeUpdate proto.InternalMessageInfo

type Response struct {
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_96a1a4bff833e11e, []int{1}
}
func (m *Response) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalTo(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (dst *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(dst, src)
}
func (m *Response) XXX_Size() int {
	return m.Size()
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Update)(nil), "cockroach.kv.kvserver.ctupdate.Update")
	proto.RegisterType((*Update_GroupUpdate)(nil), "cockroach.kv.kvserver.ctupdate.Update.GroupUpdate")
	proto.RegisterType((*Update_RangeUpdate)(nil), "cockroach.kv.kvserver.ctupdate.Update.RangeUpdate")
	proto.RegisterType((*Response)(nil), "cockroach.kv.kvserver.ctupdate.Response")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ClosedTimestampClient is the client API for ClosedTimestamp service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ClosedTimestampClient interface {
	Get(ctx context.Context, opts ...grpc.CallOption) (ClosedTimestamp_GetClient, error)
}

type closedTimestampClient struct {
	cc *grpc.ClientConn
}

func NewClosedTimestampClient(cc *grpc.ClientConn) ClosedTimestampClient {
	return &closedTimestampClient{cc}
}

func (c *closedTimestampClient) Get(ctx context.Context, opts ...grpc.CallOption) (ClosedTimestamp_GetClient, error) {
	stream, err := c.cc.NewStream(ctx, &_ClosedTimestamp_serviceDesc.Streams[0], "/cockroach.kv.kvserver.ctupdate.ClosedTimestamp/Get", opts...)
	if err != nil {
		return nil, err
	}
	x := &closedTimestampGetClient{stream}
	return x, nil
}

type ClosedTimestamp_GetClient interface {
	Send(*Reaction) error
	Recv() (*Entry, error)
	grpc.ClientStream
}

type closedTimestampGetClient struct {
	grpc.ClientStream
}

func (x *closedTimestampGetClient) Send(m *Reaction) error {
	return x.ClientStream.SendMsg(m)
}

func (x *closedTimestampGetClient) Recv() (*Entry, error) {
	m := new(Entry)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ClosedTimestampServer is the server API for ClosedTimestamp service.
type ClosedTimestampServer interface {
	Get(ClosedTimestamp_GetServer) error
}

func RegisterClosedTimestampServer(s *grpc.Server, srv ClosedTimestampServer) {
	s.RegisterService(&_ClosedTimestamp_serviceDesc, srv)
}

func _ClosedTimestamp_Get_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ClosedTimestampServer).Get(&closedTimestampGetServer{stream})
}

type ClosedTimestamp_GetServer interface {
	Send(*Entry) error
	Recv() (*Reaction, error)
	grpc.ServerStream
}

type closedTimestampGetServer struct {
	grpc.ServerStream
}

func (x *closedTimestampGetServer) Send(m *Entry) error {
	return x.ServerStream.SendMsg(m)
}

func (x *closedTimestampGetServer) Recv() (*Reaction, error) {
	m := new(Reaction)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _ClosedTimestamp_serviceDesc = grpc.ServiceDesc{
	ServiceName: "cockroach.kv.kvserver.ctupdate.ClosedTimestamp",
	HandlerType: (*ClosedTimestampServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Get",
			Handler:       _ClosedTimestamp_Get_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "kv/kvserver/closedts/ctpb/service.proto",
}

// SideTransportClient is the client API for SideTransport service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SideTransportClient interface {
	PushUpdates(ctx context.Context, opts ...grpc.CallOption) (SideTransport_PushUpdatesClient, error)
}

type sideTransportClient struct {
	cc *grpc.ClientConn
}

func NewSideTransportClient(cc *grpc.ClientConn) SideTransportClient {
	return &sideTransportClient{cc}
}

func (c *sideTransportClient) PushUpdates(ctx context.Context, opts ...grpc.CallOption) (SideTransport_PushUpdatesClient, error) {
	stream, err := c.cc.NewStream(ctx, &_SideTransport_serviceDesc.Streams[0], "/cockroach.kv.kvserver.ctupdate.SideTransport/PushUpdates", opts...)
	if err != nil {
		return nil, err
	}
	x := &sideTransportPushUpdatesClient{stream}
	return x, nil
}

type SideTransport_PushUpdatesClient interface {
	Send(*Update) error
	Recv() (*Response, error)
	grpc.ClientStream
}

type sideTransportPushUpdatesClient struct {
	grpc.ClientStream
}

func (x *sideTransportPushUpdatesClient) Send(m *Update) error {
	return x.ClientStream.SendMsg(m)
}

func (x *sideTransportPushUpdatesClient) Recv() (*Response, error) {
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SideTransportServer is the server API for SideTransport service.
type SideTransportServer interface {
	PushUpdates(SideTransport_PushUpdatesServer) error
}

func RegisterSideTransportServer(s *grpc.Server, srv SideTransportServer) {
	s.RegisterService(&_SideTransport_serviceDesc, srv)
}

func _SideTransport_PushUpdates_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SideTransportServer).PushUpdates(&sideTransportPushUpdatesServer{stream})
}

type SideTransport_PushUpdatesServer interface {
	Send(*Response) error
	Recv() (*Update, error)
	grpc.ServerStream
}

type sideTransportPushUpdatesServer struct {
	grpc.ServerStream
}

func (x *sideTransportPushUpdatesServer) Send(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *sideTransportPushUpdatesServer) Recv() (*Update, error) {
	m := new(Update)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _SideTransport_serviceDesc = grpc.ServiceDesc{
	ServiceName: "cockroach.kv.kvserver.ctupdate.SideTransport",
	HandlerType: (*SideTransportServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "PushUpdates",
			Handler:       _SideTransport_PushUpdates_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "kv/kvserver/closedts/ctpb/service.proto",
}

func (m *Update) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Update) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.NodeID != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintService(dAtA, i, uint64(m.NodeID))
	}
	if m.SeqNum != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintService(dAtA, i, uint64(m.SeqNum))
	}
	if m.Snapshot {
		dAtA[i] = 0x18
		i++
		if m.Snapshot {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if len(m.ClosedTimestamps) > 0 {
		for _, msg := range m.ClosedTimestamps {
			dAtA[i] = 0x22
			i++
			i = encodeVarintService(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if len(m.Removed) > 0 {
		dAtA2 := make([]byte, len(m.Removed)*10)
		var j1 int
		for _, num1 := range m.Removed {
			num := uint64(num1)
			for num >= 1<<7 {
				dAtA2[j1] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j1++
			}
			dAtA2[j1] = uint8(num)
			j1++
		}
		dAtA[i] = 0x2a
		i++
		i = encodeVarintService(dAtA, i, uint64(j1))
		i += copy(dAtA[i:], dAtA2[:j1])
	}
	if len(m.AddedOrUpdated) > 0 {
		for _, msg := range m.AddedOrUpdated {
			dAtA[i] = 0x32
			i++
			i = encodeVarintService(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func (m *Update_GroupUpdate) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Update_GroupUpdate) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Policy != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintService(dAtA, i, uint64(m.Policy))
	}
	dAtA[i] = 0x12
	i++
	i = encodeVarintService(dAtA, i, uint64(m.ClosedTimestamp.Size()))
	n3, err := m.ClosedTimestamp.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n3
	return i, nil
}

func (m *Update_RangeUpdate) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Update_RangeUpdate) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.RangeID != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintService(dAtA, i, uint64(m.RangeID))
	}
	if m.LAI != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintService(dAtA, i, uint64(m.LAI))
	}
	if m.Policy != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintService(dAtA, i, uint64(m.Policy))
	}
	return i, nil
}

func (m *Response) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Response) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	return i, nil
}

func encodeVarintService(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *Update) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.NodeID != 0 {
		n += 1 + sovService(uint64(m.NodeID))
	}
	if m.SeqNum != 0 {
		n += 1 + sovService(uint64(m.SeqNum))
	}
	if m.Snapshot {
		n += 2
	}
	if len(m.ClosedTimestamps) > 0 {
		for _, e := range m.ClosedTimestamps {
			l = e.Size()
			n += 1 + l + sovService(uint64(l))
		}
	}
	if len(m.Removed) > 0 {
		l = 0
		for _, e := range m.Removed {
			l += sovService(uint64(e))
		}
		n += 1 + sovService(uint64(l)) + l
	}
	if len(m.AddedOrUpdated) > 0 {
		for _, e := range m.AddedOrUpdated {
			l = e.Size()
			n += 1 + l + sovService(uint64(l))
		}
	}
	return n
}

func (m *Update_GroupUpdate) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Policy != 0 {
		n += 1 + sovService(uint64(m.Policy))
	}
	l = m.ClosedTimestamp.Size()
	n += 1 + l + sovService(uint64(l))
	return n
}

func (m *Update_RangeUpdate) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.RangeID != 0 {
		n += 1 + sovService(uint64(m.RangeID))
	}
	if m.LAI != 0 {
		n += 1 + sovService(uint64(m.LAI))
	}
	if m.Policy != 0 {
		n += 1 + sovService(uint64(m.Policy))
	}
	return n
}

func (m *Response) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func sovService(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozService(x uint64) (n int) {
	return sovService(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Update) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowService
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
			return fmt.Errorf("proto: Update: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Update: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NodeID", wireType)
			}
			m.NodeID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.NodeID |= (github_com_cockroachdb_cockroach_pkg_roachpb.NodeID(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SeqNum", wireType)
			}
			m.SeqNum = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SeqNum |= (SeqNum(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Snapshot", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Snapshot = bool(v != 0)
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClosedTimestamps", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowService
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
				return ErrInvalidLengthService
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ClosedTimestamps = append(m.ClosedTimestamps, Update_GroupUpdate{})
			if err := m.ClosedTimestamps[len(m.ClosedTimestamps)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType == 0 {
				var v github_com_cockroachdb_cockroach_pkg_roachpb.RangeID
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowService
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= (github_com_cockroachdb_cockroach_pkg_roachpb.RangeID(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.Removed = append(m.Removed, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowService
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= (int(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthService
				}
				postIndex := iNdEx + packedLen
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				var count int
				for _, integer := range dAtA {
					if integer < 128 {
						count++
					}
				}
				elementCount = count
				if elementCount != 0 && len(m.Removed) == 0 {
					m.Removed = make([]github_com_cockroachdb_cockroach_pkg_roachpb.RangeID, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v github_com_cockroachdb_cockroach_pkg_roachpb.RangeID
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowService
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= (github_com_cockroachdb_cockroach_pkg_roachpb.RangeID(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.Removed = append(m.Removed, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field Removed", wireType)
			}
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AddedOrUpdated", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowService
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
				return ErrInvalidLengthService
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AddedOrUpdated = append(m.AddedOrUpdated, Update_RangeUpdate{})
			if err := m.AddedOrUpdated[len(m.AddedOrUpdated)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipService(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthService
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
func (m *Update_GroupUpdate) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowService
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
			return fmt.Errorf("proto: GroupUpdate: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GroupUpdate: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Policy", wireType)
			}
			m.Policy = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Policy |= (roachpb.RangeClosedTimestampPolicy(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClosedTimestamp", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowService
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
				return ErrInvalidLengthService
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ClosedTimestamp.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipService(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthService
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
func (m *Update_RangeUpdate) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowService
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
			return fmt.Errorf("proto: RangeUpdate: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RangeUpdate: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RangeID", wireType)
			}
			m.RangeID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RangeID |= (github_com_cockroachdb_cockroach_pkg_roachpb.RangeID(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LAI", wireType)
			}
			m.LAI = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LAI |= (LAI(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Policy", wireType)
			}
			m.Policy = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Policy |= (roachpb.RangeClosedTimestampPolicy(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipService(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthService
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
func (m *Response) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowService
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
			return fmt.Errorf("proto: Response: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Response: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipService(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthService
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
func skipService(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowService
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
					return 0, ErrIntOverflowService
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
					return 0, ErrIntOverflowService
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
				return 0, ErrInvalidLengthService
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowService
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
				next, err := skipService(dAtA[start:])
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
	ErrInvalidLengthService = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowService   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("kv/kvserver/closedts/ctpb/service.proto", fileDescriptor_service_96a1a4bff833e11e)
}

var fileDescriptor_service_96a1a4bff833e11e = []byte{
	// 628 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x54, 0xc1, 0x4e, 0xdb, 0x30,
	0x18, 0x6e, 0x96, 0x92, 0x56, 0xae, 0x06, 0xcc, 0xda, 0x21, 0x8a, 0xb6, 0xa4, 0x62, 0x82, 0xe5,
	0xb2, 0x78, 0x2a, 0x3b, 0xec, 0x4a, 0x19, 0x42, 0x48, 0x13, 0x43, 0x81, 0x5d, 0xd0, 0xb4, 0xca,
	0x8d, 0xad, 0x34, 0x6a, 0x1a, 0x07, 0xdb, 0xa9, 0xc4, 0x5b, 0xec, 0x21, 0xf6, 0x28, 0x3b, 0x70,
	0xe4, 0x34, 0x71, 0x8a, 0xb6, 0xf0, 0x16, 0x3d, 0x4d, 0x89, 0x43, 0x57, 0x2a, 0x6d, 0x80, 0x38,
	0xf9, 0xf7, 0x6f, 0xff, 0x9f, 0xbf, 0xef, 0xf3, 0x6f, 0x83, 0xd7, 0xe3, 0x29, 0x1a, 0x4f, 0x05,
	0xe5, 0x53, 0xca, 0x51, 0x10, 0x33, 0x41, 0x89, 0x14, 0x28, 0x90, 0xe9, 0x10, 0x95, 0xc9, 0x28,
	0xa0, 0x5e, 0xca, 0x99, 0x64, 0xd0, 0x0e, 0x58, 0x30, 0xe6, 0x0c, 0x07, 0x23, 0x6f, 0x3c, 0xf5,
	0x6e, 0x4a, 0xbc, 0x40, 0x66, 0x29, 0xc1, 0x92, 0x5a, 0x9b, 0xff, 0x06, 0xa2, 0x89, 0xe4, 0xe7,
	0x0a, 0xc6, 0x82, 0x15, 0x44, 0x3a, 0x44, 0x04, 0x4b, 0x5c, 0xe7, 0xcc, 0x4c, 0x46, 0x31, 0x1a,
	0xc5, 0x01, 0x92, 0xd1, 0x84, 0x0a, 0x89, 0x27, 0x69, 0xbd, 0xf2, 0x3c, 0x64, 0x21, 0xab, 0x42,
	0x54, 0x46, 0x75, 0xf6, 0x45, 0xc8, 0x58, 0x18, 0x53, 0x84, 0xd3, 0x08, 0xe1, 0x24, 0x61, 0x12,
	0xcb, 0x88, 0x25, 0x42, 0xad, 0x6e, 0xfc, 0x30, 0x80, 0xf1, 0xb9, 0xe2, 0x04, 0x4f, 0x41, 0x2b,
	0x61, 0x84, 0x0e, 0x22, 0x62, 0x6a, 0x5d, 0xcd, 0x5d, 0xe9, 0xef, 0x14, 0xb9, 0x63, 0x1c, 0x32,
	0x42, 0x0f, 0x3e, 0xcc, 0x72, 0x67, 0x3b, 0x8c, 0xe4, 0x28, 0x1b, 0x7a, 0x01, 0x9b, 0xa0, 0xb9,
	0x3a, 0x32, 0xfc, 0x1b, 0xa3, 0x74, 0x1c, 0xa2, 0x9a, 0xb0, 0xa7, 0xca, 0x7c, 0xa3, 0x44, 0x3c,
	0x20, 0xf0, 0x15, 0x68, 0x09, 0x7a, 0x36, 0x48, 0xb2, 0x89, 0xf9, 0xa4, 0xab, 0xb9, 0x7a, 0x1f,
	0xcc, 0x72, 0xc7, 0x38, 0xa6, 0x67, 0x87, 0xd9, 0xc4, 0x37, 0x44, 0x35, 0x42, 0x0b, 0xb4, 0x45,
	0x82, 0x53, 0x31, 0x62, 0xd2, 0xd4, 0xbb, 0x9a, 0xdb, 0xf6, 0xe7, 0x73, 0x48, 0xc1, 0x33, 0x65,
	0xd3, 0x60, 0xae, 0x5a, 0x98, 0xcd, 0xae, 0xee, 0x76, 0x7a, 0x3d, 0xef, 0xff, 0x66, 0x7b, 0x4a,
	0x9f, 0xb7, 0xcf, 0x59, 0x96, 0xaa, 0xb8, 0xdf, 0xbc, 0xc8, 0x9d, 0x86, 0xbf, 0xae, 0x20, 0x4f,
	0xe6, 0x88, 0xd0, 0x07, 0x2d, 0x4e, 0x27, 0x6c, 0x4a, 0x89, 0xb9, 0xd2, 0xd5, 0xdd, 0x95, 0xfe,
	0xfb, 0x59, 0xee, 0xbc, 0x7b, 0x90, 0x72, 0x1f, 0x27, 0x61, 0x29, 0xfd, 0x06, 0x08, 0x0e, 0xc1,
	0x3a, 0x26, 0x84, 0x92, 0x01, 0xe3, 0x03, 0xc5, 0x88, 0x98, 0xc6, 0x83, 0x98, 0x57, 0x90, 0xb7,
	0x98, 0xaf, 0x56, 0x88, 0x9f, 0xb8, 0x4a, 0x12, 0xeb, 0xbb, 0x06, 0x3a, 0x0b, 0xfa, 0xe0, 0x1e,
	0x30, 0x52, 0x16, 0x47, 0xc1, 0x79, 0x75, 0x95, 0xab, 0xbd, 0x37, 0x0b, 0x27, 0xdd, 0x22, 0xba,
	0x7b, 0xdb, 0x81, 0xa3, 0xaa, 0xc8, 0xaf, 0x8b, 0xe1, 0x21, 0x58, 0x5f, 0x76, 0xbd, 0xba, 0xbf,
	0x4e, 0xef, 0xe5, 0x02, 0x60, 0xd9, 0x90, 0xde, 0x28, 0x0e, 0xbc, 0x39, 0x4c, 0xcd, 0x72, 0x6d,
	0xc9, 0x5f, 0xeb, 0xa7, 0x06, 0x3a, 0x0b, 0x62, 0xe0, 0x57, 0xd0, 0xe6, 0xe5, 0xf4, 0xa6, 0xe7,
	0x9a, 0xfd, 0xdd, 0x22, 0x77, 0x5a, 0xb5, 0x85, 0x8f, 0xb0, 0xbe, 0x0a, 0x08, 0xec, 0x02, 0x3d,
	0xc6, 0x51, 0xdd, 0x72, 0xab, 0x45, 0xee, 0xe8, 0x1f, 0x77, 0x0e, 0x66, 0x6a, 0xf0, 0xcb, 0xa5,
	0x05, 0xa3, 0xf4, 0x47, 0x18, 0xb5, 0x01, 0x40, 0xdb, 0xa7, 0x22, 0x65, 0x89, 0xa0, 0x3d, 0x06,
	0xd6, 0x96, 0x36, 0xc3, 0x2f, 0x40, 0xdf, 0xa7, 0x12, 0xba, 0x77, 0xdd, 0xb7, 0x4f, 0x71, 0x50,
	0xbe, 0x4e, 0x6b, 0xf3, 0xae, 0x9d, 0x7b, 0xe5, 0x2f, 0xb1, 0xd1, 0x70, 0xb5, 0xb7, 0x5a, 0x6f,
	0x0a, 0x9e, 0x1e, 0x47, 0x84, 0x9e, 0x70, 0x9c, 0x88, 0x94, 0xf1, 0xf2, 0xb1, 0x74, 0x8e, 0x32,
	0x31, 0x52, 0x26, 0x0b, 0xb8, 0x75, 0xbf, 0x36, 0xb3, 0xee, 0x41, 0x4f, 0x49, 0x54, 0xe7, 0xf6,
	0xb7, 0x2e, 0x7e, 0xdb, 0x8d, 0x8b, 0xc2, 0xd6, 0x2e, 0x0b, 0x5b, 0xbb, 0x2a, 0x6c, 0xed, 0x57,
	0x61, 0x6b, 0xdf, 0xae, 0xed, 0xc6, 0xe5, 0xb5, 0xdd, 0xb8, 0xba, 0xb6, 0x1b, 0xa7, 0xcd, 0xf2,
	0x47, 0x1b, 0x1a, 0xd5, 0x57, 0xb3, 0xfd, 0x27, 0x00, 0x00, 0xff, 0xff, 0x2d, 0x8d, 0x65, 0x3b,
	0x3e, 0x05, 0x00, 0x00,
}
