// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1-devel
// 	protoc        v3.17.3
// source: pkg/event/proto/event.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PublishedStatus int32

const (
	PublishedStatus_PublishedStatus_UNSPECIFIED PublishedStatus = 0
	PublishedStatus_PublishedStatus_Draft       PublishedStatus = 1
	PublishedStatus_PublishedStatus_Published   PublishedStatus = 2
)

// Enum value maps for PublishedStatus.
var (
	PublishedStatus_name = map[int32]string{
		0: "PublishedStatus_UNSPECIFIED",
		1: "PublishedStatus_Draft",
		2: "PublishedStatus_Published",
	}
	PublishedStatus_value = map[string]int32{
		"PublishedStatus_UNSPECIFIED": 0,
		"PublishedStatus_Draft":       1,
		"PublishedStatus_Published":   2,
	}
)

func (x PublishedStatus) Enum() *PublishedStatus {
	p := new(PublishedStatus)
	*p = x
	return p
}

func (x PublishedStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PublishedStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_pkg_event_proto_event_proto_enumTypes[0].Descriptor()
}

func (PublishedStatus) Type() protoreflect.EnumType {
	return &file_pkg_event_proto_event_proto_enumTypes[0]
}

func (x PublishedStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PublishedStatus.Descriptor instead.
func (PublishedStatus) EnumDescriptor() ([]byte, []int) {
	return file_pkg_event_proto_event_proto_rawDescGZIP(), []int{0}
}

type Event struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id              int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Title           string                 `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Description     string                 `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	PublishedStatus PublishedStatus        `protobuf:"varint,4,opt,name=published_status,json=publishedStatus,proto3,enum=proto.PublishedStatus" json:"published_status,omitempty"`
	CreatedAt       *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt       *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_event_proto_event_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_event_proto_event_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Event.ProtoReflect.Descriptor instead.
func (*Event) Descriptor() ([]byte, []int) {
	return file_pkg_event_proto_event_proto_rawDescGZIP(), []int{0}
}

func (x *Event) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Event) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Event) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Event) GetPublishedStatus() PublishedStatus {
	if x != nil {
		return x.PublishedStatus
	}
	return PublishedStatus_PublishedStatus_UNSPECIFIED
}

func (x *Event) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Event) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

type GetEventsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id              int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Title           string                 `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	PublishedStatus PublishedStatus        `protobuf:"varint,4,opt,name=published_status,json=publishedStatus,proto3,enum=proto.PublishedStatus" json:"published_status,omitempty"`
	CreatedAt       *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
}

func (x *GetEventsRequest) Reset() {
	*x = GetEventsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_event_proto_event_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetEventsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetEventsRequest) ProtoMessage() {}

func (x *GetEventsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_event_proto_event_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetEventsRequest.ProtoReflect.Descriptor instead.
func (*GetEventsRequest) Descriptor() ([]byte, []int) {
	return file_pkg_event_proto_event_proto_rawDescGZIP(), []int{1}
}

func (x *GetEventsRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *GetEventsRequest) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *GetEventsRequest) GetPublishedStatus() PublishedStatus {
	if x != nil {
		return x.PublishedStatus
	}
	return PublishedStatus_PublishedStatus_UNSPECIFIED
}

func (x *GetEventsRequest) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

type GetEventsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Events []*Event `protobuf:"bytes,1,rep,name=Events,proto3" json:"Events,omitempty"`
}

func (x *GetEventsResponse) Reset() {
	*x = GetEventsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_event_proto_event_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetEventsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetEventsResponse) ProtoMessage() {}

func (x *GetEventsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_event_proto_event_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetEventsResponse.ProtoReflect.Descriptor instead.
func (*GetEventsResponse) Descriptor() ([]byte, []int) {
	return file_pkg_event_proto_event_proto_rawDescGZIP(), []int{2}
}

func (x *GetEventsResponse) GetEvents() []*Event {
	if x != nil {
		return x.Events
	}
	return nil
}

type UpdatePublishStatusRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventId         int64           `protobuf:"varint,1,opt,name=event_id,json=eventId,proto3" json:"event_id,omitempty"`
	TransId         string          `protobuf:"bytes,2,opt,name=trans_id,json=transId,proto3" json:"trans_id,omitempty"`
	PublishedStatus PublishedStatus `protobuf:"varint,3,opt,name=published_status,json=publishedStatus,proto3,enum=proto.PublishedStatus" json:"published_status,omitempty"`
}

func (x *UpdatePublishStatusRequest) Reset() {
	*x = UpdatePublishStatusRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_event_proto_event_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdatePublishStatusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdatePublishStatusRequest) ProtoMessage() {}

func (x *UpdatePublishStatusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_event_proto_event_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdatePublishStatusRequest.ProtoReflect.Descriptor instead.
func (*UpdatePublishStatusRequest) Descriptor() ([]byte, []int) {
	return file_pkg_event_proto_event_proto_rawDescGZIP(), []int{3}
}

func (x *UpdatePublishStatusRequest) GetEventId() int64 {
	if x != nil {
		return x.EventId
	}
	return 0
}

func (x *UpdatePublishStatusRequest) GetTransId() string {
	if x != nil {
		return x.TransId
	}
	return ""
}

func (x *UpdatePublishStatusRequest) GetPublishedStatus() PublishedStatus {
	if x != nil {
		return x.PublishedStatus
	}
	return PublishedStatus_PublishedStatus_UNSPECIFIED
}

var File_pkg_event_proto_event_proto protoreflect.FileDescriptor

var file_pkg_event_proto_event_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x70, 0x6b, 0x67, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x88, 0x02, 0x0a, 0x05, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05,
	0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74,
	0x6c, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x41, 0x0a, 0x10, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x65,
	0x64, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x16,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x65, 0x64,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x0f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x65,
	0x64, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64,
	0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0xb6, 0x01,
	0x0a, 0x10, 0x47, 0x65, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x41, 0x0a, 0x10, 0x70, 0x75, 0x62, 0x6c,
	0x69, 0x73, 0x68, 0x65, 0x64, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x75, 0x62, 0x6c, 0x69,
	0x73, 0x68, 0x65, 0x64, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x0f, 0x70, 0x75, 0x62, 0x6c,
	0x69, 0x73, 0x68, 0x65, 0x64, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x39, 0x0a, 0x0a, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x39, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x24, 0x0a, 0x06, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x06, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x73, 0x22, 0x95, 0x01, 0x0a, 0x1a, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x75, 0x62, 0x6c,
	0x69, 0x73, 0x68, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x19, 0x0a, 0x08, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x07, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x74,
	0x72, 0x61, 0x6e, 0x73, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x74,
	0x72, 0x61, 0x6e, 0x73, 0x49, 0x64, 0x12, 0x41, 0x0a, 0x10, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73,
	0x68, 0x65, 0x64, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68,
	0x65, 0x64, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x0f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73,
	0x68, 0x65, 0x64, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2a, 0x6c, 0x0a, 0x0f, 0x50, 0x75, 0x62,
	0x6c, 0x69, 0x73, 0x68, 0x65, 0x64, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1f, 0x0a, 0x1b,
	0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x65, 0x64, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f,
	0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x19, 0x0a,
	0x15, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x65, 0x64, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x5f, 0x44, 0x72, 0x61, 0x66, 0x74, 0x10, 0x01, 0x12, 0x1d, 0x0a, 0x19, 0x50, 0x75, 0x62, 0x6c,
	0x69, 0x73, 0x68, 0x65, 0x64, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x50, 0x75, 0x62, 0x6c,
	0x69, 0x73, 0x68, 0x65, 0x64, 0x10, 0x02, 0x32, 0xa0, 0x01, 0x0a, 0x0c, 0x45, 0x76, 0x65, 0x6e,
	0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3e, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65,
	0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x50, 0x0a, 0x13, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12,
	0x21, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x75,
	0x62, 0x6c, 0x69, 0x73, 0x68, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x36, 0x5a, 0x34, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6e, 0x69, 0x74, 0x65, 0x2d, 0x63, 0x6f,
	0x64, 0x65, 0x72, 0x2f, 0x62, 0x6c, 0x61, 0x63, 0x6b, 0x62, 0x65, 0x61, 0x72, 0x2d, 0x64, 0x65,
	0x6d, 0x6f, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_event_proto_event_proto_rawDescOnce sync.Once
	file_pkg_event_proto_event_proto_rawDescData = file_pkg_event_proto_event_proto_rawDesc
)

func file_pkg_event_proto_event_proto_rawDescGZIP() []byte {
	file_pkg_event_proto_event_proto_rawDescOnce.Do(func() {
		file_pkg_event_proto_event_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_event_proto_event_proto_rawDescData)
	})
	return file_pkg_event_proto_event_proto_rawDescData
}

var file_pkg_event_proto_event_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_pkg_event_proto_event_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_pkg_event_proto_event_proto_goTypes = []interface{}{
	(PublishedStatus)(0),               // 0: proto.PublishedStatus
	(*Event)(nil),                      // 1: proto.Event
	(*GetEventsRequest)(nil),           // 2: proto.GetEventsRequest
	(*GetEventsResponse)(nil),          // 3: proto.GetEventsResponse
	(*UpdatePublishStatusRequest)(nil), // 4: proto.UpdatePublishStatusRequest
	(*timestamppb.Timestamp)(nil),      // 5: google.protobuf.Timestamp
	(*emptypb.Empty)(nil),              // 6: google.protobuf.Empty
}
var file_pkg_event_proto_event_proto_depIdxs = []int32{
	0, // 0: proto.Event.published_status:type_name -> proto.PublishedStatus
	5, // 1: proto.Event.created_at:type_name -> google.protobuf.Timestamp
	5, // 2: proto.Event.updated_at:type_name -> google.protobuf.Timestamp
	0, // 3: proto.GetEventsRequest.published_status:type_name -> proto.PublishedStatus
	5, // 4: proto.GetEventsRequest.created_at:type_name -> google.protobuf.Timestamp
	1, // 5: proto.GetEventsResponse.Events:type_name -> proto.Event
	0, // 6: proto.UpdatePublishStatusRequest.published_status:type_name -> proto.PublishedStatus
	2, // 7: proto.EventService.GetEvents:input_type -> proto.GetEventsRequest
	4, // 8: proto.EventService.UpdatePublishStatus:input_type -> proto.UpdatePublishStatusRequest
	3, // 9: proto.EventService.GetEvents:output_type -> proto.GetEventsResponse
	6, // 10: proto.EventService.UpdatePublishStatus:output_type -> google.protobuf.Empty
	9, // [9:11] is the sub-list for method output_type
	7, // [7:9] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_pkg_event_proto_event_proto_init() }
func file_pkg_event_proto_event_proto_init() {
	if File_pkg_event_proto_event_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_event_proto_event_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Event); i {
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
		file_pkg_event_proto_event_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetEventsRequest); i {
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
		file_pkg_event_proto_event_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetEventsResponse); i {
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
		file_pkg_event_proto_event_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdatePublishStatusRequest); i {
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
			RawDescriptor: file_pkg_event_proto_event_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_event_proto_event_proto_goTypes,
		DependencyIndexes: file_pkg_event_proto_event_proto_depIdxs,
		EnumInfos:         file_pkg_event_proto_event_proto_enumTypes,
		MessageInfos:      file_pkg_event_proto_event_proto_msgTypes,
	}.Build()
	File_pkg_event_proto_event_proto = out.File
	file_pkg_event_proto_event_proto_rawDesc = nil
	file_pkg_event_proto_event_proto_goTypes = nil
	file_pkg_event_proto_event_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// EventServiceClient is the client API for EventService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type EventServiceClient interface {
	GetEvents(ctx context.Context, in *GetEventsRequest, opts ...grpc.CallOption) (*GetEventsResponse, error)
	UpdatePublishStatus(ctx context.Context, in *UpdatePublishStatusRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type eventServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEventServiceClient(cc grpc.ClientConnInterface) EventServiceClient {
	return &eventServiceClient{cc}
}

func (c *eventServiceClient) GetEvents(ctx context.Context, in *GetEventsRequest, opts ...grpc.CallOption) (*GetEventsResponse, error) {
	out := new(GetEventsResponse)
	err := c.cc.Invoke(ctx, "/proto.EventService/GetEvents", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) UpdatePublishStatus(ctx context.Context, in *UpdatePublishStatusRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/proto.EventService/UpdatePublishStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EventServiceServer is the server API for EventService service.
type EventServiceServer interface {
	GetEvents(context.Context, *GetEventsRequest) (*GetEventsResponse, error)
	UpdatePublishStatus(context.Context, *UpdatePublishStatusRequest) (*emptypb.Empty, error)
}

// UnimplementedEventServiceServer can be embedded to have forward compatible implementations.
type UnimplementedEventServiceServer struct {
}

func (*UnimplementedEventServiceServer) GetEvents(context.Context, *GetEventsRequest) (*GetEventsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEvents not implemented")
}
func (*UnimplementedEventServiceServer) UpdatePublishStatus(context.Context, *UpdatePublishStatusRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePublishStatus not implemented")
}

func RegisterEventServiceServer(s *grpc.Server, srv EventServiceServer) {
	s.RegisterService(&_EventService_serviceDesc, srv)
}

func _EventService_GetEvents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEventsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).GetEvents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.EventService/GetEvents",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).GetEvents(ctx, req.(*GetEventsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_UpdatePublishStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePublishStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).UpdatePublishStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.EventService/UpdatePublishStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).UpdatePublishStatus(ctx, req.(*UpdatePublishStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _EventService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.EventService",
	HandlerType: (*EventServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetEvents",
			Handler:    _EventService_GetEvents_Handler,
		},
		{
			MethodName: "UpdatePublishStatus",
			Handler:    _EventService_UpdatePublishStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/event/proto/event.proto",
}
