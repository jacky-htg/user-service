// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v5.26.1
// source: users/group_service.proto

package users

import (
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

type ListGroupRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pagination *Pagination `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
	CompanyId  string      `protobuf:"bytes,2,opt,name=company_id,json=companyId,proto3" json:"company_id,omitempty"`
}

func (x *ListGroupRequest) Reset() {
	*x = ListGroupRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_users_group_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListGroupRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListGroupRequest) ProtoMessage() {}

func (x *ListGroupRequest) ProtoReflect() protoreflect.Message {
	mi := &file_users_group_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListGroupRequest.ProtoReflect.Descriptor instead.
func (*ListGroupRequest) Descriptor() ([]byte, []int) {
	return file_users_group_service_proto_rawDescGZIP(), []int{0}
}

func (x *ListGroupRequest) GetPagination() *Pagination {
	if x != nil {
		return x.Pagination
	}
	return nil
}

func (x *ListGroupRequest) GetCompanyId() string {
	if x != nil {
		return x.CompanyId
	}
	return ""
}

type GroupPaginationResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pagination *Pagination `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
	Count      uint32      `protobuf:"varint,3,opt,name=count,proto3" json:"count,omitempty"`
}

func (x *GroupPaginationResponse) Reset() {
	*x = GroupPaginationResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_users_group_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GroupPaginationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GroupPaginationResponse) ProtoMessage() {}

func (x *GroupPaginationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_users_group_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GroupPaginationResponse.ProtoReflect.Descriptor instead.
func (*GroupPaginationResponse) Descriptor() ([]byte, []int) {
	return file_users_group_service_proto_rawDescGZIP(), []int{1}
}

func (x *GroupPaginationResponse) GetPagination() *Pagination {
	if x != nil {
		return x.Pagination
	}
	return nil
}

func (x *GroupPaginationResponse) GetCount() uint32 {
	if x != nil {
		return x.Count
	}
	return 0
}

type ListGroupResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pagination *GroupPaginationResponse `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
	Group      *Group                   `protobuf:"bytes,2,opt,name=group,proto3" json:"group,omitempty"`
}

func (x *ListGroupResponse) Reset() {
	*x = ListGroupResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_users_group_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListGroupResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListGroupResponse) ProtoMessage() {}

func (x *ListGroupResponse) ProtoReflect() protoreflect.Message {
	mi := &file_users_group_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListGroupResponse.ProtoReflect.Descriptor instead.
func (*ListGroupResponse) Descriptor() ([]byte, []int) {
	return file_users_group_service_proto_rawDescGZIP(), []int{2}
}

func (x *ListGroupResponse) GetPagination() *GroupPaginationResponse {
	if x != nil {
		return x.Pagination
	}
	return nil
}

func (x *ListGroupResponse) GetGroup() *Group {
	if x != nil {
		return x.Group
	}
	return nil
}

type GrantAccessRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GroupId  string `protobuf:"bytes,1,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`
	AccessId string `protobuf:"bytes,2,opt,name=access_id,json=accessId,proto3" json:"access_id,omitempty"`
}

func (x *GrantAccessRequest) Reset() {
	*x = GrantAccessRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_users_group_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GrantAccessRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GrantAccessRequest) ProtoMessage() {}

func (x *GrantAccessRequest) ProtoReflect() protoreflect.Message {
	mi := &file_users_group_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GrantAccessRequest.ProtoReflect.Descriptor instead.
func (*GrantAccessRequest) Descriptor() ([]byte, []int) {
	return file_users_group_service_proto_rawDescGZIP(), []int{3}
}

func (x *GrantAccessRequest) GetGroupId() string {
	if x != nil {
		return x.GroupId
	}
	return ""
}

func (x *GrantAccessRequest) GetAccessId() string {
	if x != nil {
		return x.AccessId
	}
	return ""
}

var File_users_group_service_proto protoreflect.FileDescriptor

var file_users_group_service_proto_rawDesc = []byte{
	0x0a, 0x19, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x77, 0x69, 0x72,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x1a, 0x19, 0x75, 0x73, 0x65,
	0x72, 0x73, 0x2f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2f, 0x67, 0x65,
	0x6e, 0x65, 0x72, 0x69, 0x63, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x6d, 0x0a, 0x10, 0x4c, 0x69, 0x73, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x3a, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x77, 0x69,
	0x72, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x50, 0x61, 0x67,
	0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x6e, 0x79, 0x5f, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x6e, 0x79,
	0x49, 0x64, 0x22, 0x6b, 0x0a, 0x17, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x50, 0x61, 0x67, 0x69, 0x6e,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3a, 0x0a,
	0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x77, 0x69, 0x72, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x75, 0x73, 0x65,
	0x72, 0x73, 0x2e, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x70,
	0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22,
	0x89, 0x01, 0x0a, 0x11, 0x4c, 0x69, 0x73, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x47, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x77, 0x69, 0x72, 0x61,
	0x64, 0x61, 0x74, 0x61, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70,
	0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x52, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x2b,
	0x0a, 0x05, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e,
	0x77, 0x69, 0x72, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x47,
	0x72, 0x6f, 0x75, 0x70, 0x52, 0x05, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x22, 0x4c, 0x0a, 0x12, 0x47,
	0x72, 0x61, 0x6e, 0x74, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x19, 0x0a, 0x08, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09,
	0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x49, 0x64, 0x32, 0xe0, 0x03, 0x0a, 0x0c, 0x47, 0x72,
	0x6f, 0x75, 0x70, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x38, 0x0a, 0x06, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x12, 0x15, 0x2e, 0x77, 0x69, 0x72, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e,
	0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x1a, 0x15, 0x2e, 0x77, 0x69,
	0x72, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x47, 0x72, 0x6f,
	0x75, 0x70, 0x22, 0x00, 0x12, 0x38, 0x0a, 0x06, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x15,
	0x2e, 0x77, 0x69, 0x72, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e,
	0x47, 0x72, 0x6f, 0x75, 0x70, 0x1a, 0x15, 0x2e, 0x77, 0x69, 0x72, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x22, 0x00, 0x12, 0x33,
	0x0a, 0x04, 0x56, 0x69, 0x65, 0x77, 0x12, 0x12, 0x2e, 0x77, 0x69, 0x72, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x49, 0x64, 0x1a, 0x15, 0x2e, 0x77, 0x69, 0x72,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x47, 0x72, 0x6f, 0x75,
	0x70, 0x22, 0x00, 0x12, 0x39, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x12, 0x2e,
	0x77, 0x69, 0x72, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x49,
	0x64, 0x1a, 0x19, 0x2e, 0x77, 0x69, 0x72, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x75, 0x73, 0x65,
	0x72, 0x73, 0x2e, 0x4d, 0x79, 0x42, 0x6f, 0x6f, 0x6c, 0x65, 0x61, 0x6e, 0x22, 0x00, 0x12, 0x4f,
	0x0a, 0x04, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x20, 0x2e, 0x77, 0x69, 0x72, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x47, 0x72, 0x6f, 0x75,
	0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x77, 0x69, 0x72, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x47, 0x72,
	0x6f, 0x75, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x30, 0x01, 0x12,
	0x4c, 0x0a, 0x0b, 0x47, 0x72, 0x61, 0x6e, 0x74, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x22,
	0x2e, 0x77, 0x69, 0x72, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e,
	0x47, 0x72, 0x61, 0x6e, 0x74, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x17, 0x2e, 0x77, 0x69, 0x72, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x75, 0x73,
	0x65, 0x72, 0x73, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x00, 0x12, 0x4d, 0x0a,
	0x0c, 0x52, 0x65, 0x76, 0x6f, 0x6b, 0x65, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x22, 0x2e,
	0x77, 0x69, 0x72, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x47,
	0x72, 0x61, 0x6e, 0x74, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x17, 0x2e, 0x77, 0x69, 0x72, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x75, 0x73, 0x65,
	0x72, 0x73, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x00, 0x42, 0x35, 0x0a, 0x21,
	0x63, 0x6f, 0x6d, 0x2e, 0x77, 0x69, 0x72, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x65, 0x72, 0x70,
	0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x75, 0x73, 0x65, 0x72,
	0x73, 0x50, 0x01, 0x5a, 0x0e, 0x70, 0x62, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x3b, 0x75, 0x73,
	0x65, 0x72, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_users_group_service_proto_rawDescOnce sync.Once
	file_users_group_service_proto_rawDescData = file_users_group_service_proto_rawDesc
)

func file_users_group_service_proto_rawDescGZIP() []byte {
	file_users_group_service_proto_rawDescOnce.Do(func() {
		file_users_group_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_users_group_service_proto_rawDescData)
	})
	return file_users_group_service_proto_rawDescData
}

var file_users_group_service_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_users_group_service_proto_goTypes = []interface{}{
	(*ListGroupRequest)(nil),        // 0: wiradata.users.ListGroupRequest
	(*GroupPaginationResponse)(nil), // 1: wiradata.users.GroupPaginationResponse
	(*ListGroupResponse)(nil),       // 2: wiradata.users.ListGroupResponse
	(*GrantAccessRequest)(nil),      // 3: wiradata.users.GrantAccessRequest
	(*Pagination)(nil),              // 4: wiradata.users.Pagination
	(*Group)(nil),                   // 5: wiradata.users.Group
	(*Id)(nil),                      // 6: wiradata.users.Id
	(*MyBoolean)(nil),               // 7: wiradata.users.MyBoolean
	(*Message)(nil),                 // 8: wiradata.users.Message
}
var file_users_group_service_proto_depIdxs = []int32{
	4,  // 0: wiradata.users.ListGroupRequest.pagination:type_name -> wiradata.users.Pagination
	4,  // 1: wiradata.users.GroupPaginationResponse.pagination:type_name -> wiradata.users.Pagination
	1,  // 2: wiradata.users.ListGroupResponse.pagination:type_name -> wiradata.users.GroupPaginationResponse
	5,  // 3: wiradata.users.ListGroupResponse.group:type_name -> wiradata.users.Group
	5,  // 4: wiradata.users.GroupService.Create:input_type -> wiradata.users.Group
	5,  // 5: wiradata.users.GroupService.Update:input_type -> wiradata.users.Group
	6,  // 6: wiradata.users.GroupService.View:input_type -> wiradata.users.Id
	6,  // 7: wiradata.users.GroupService.Delete:input_type -> wiradata.users.Id
	0,  // 8: wiradata.users.GroupService.List:input_type -> wiradata.users.ListGroupRequest
	3,  // 9: wiradata.users.GroupService.GrantAccess:input_type -> wiradata.users.GrantAccessRequest
	3,  // 10: wiradata.users.GroupService.RevokeAccess:input_type -> wiradata.users.GrantAccessRequest
	5,  // 11: wiradata.users.GroupService.Create:output_type -> wiradata.users.Group
	5,  // 12: wiradata.users.GroupService.Update:output_type -> wiradata.users.Group
	5,  // 13: wiradata.users.GroupService.View:output_type -> wiradata.users.Group
	7,  // 14: wiradata.users.GroupService.Delete:output_type -> wiradata.users.MyBoolean
	2,  // 15: wiradata.users.GroupService.List:output_type -> wiradata.users.ListGroupResponse
	8,  // 16: wiradata.users.GroupService.GrantAccess:output_type -> wiradata.users.Message
	8,  // 17: wiradata.users.GroupService.RevokeAccess:output_type -> wiradata.users.Message
	11, // [11:18] is the sub-list for method output_type
	4,  // [4:11] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_users_group_service_proto_init() }
func file_users_group_service_proto_init() {
	if File_users_group_service_proto != nil {
		return
	}
	file_users_group_message_proto_init()
	file_users_generic_message_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_users_group_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListGroupRequest); i {
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
		file_users_group_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GroupPaginationResponse); i {
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
		file_users_group_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListGroupResponse); i {
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
		file_users_group_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GrantAccessRequest); i {
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
			RawDescriptor: file_users_group_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_users_group_service_proto_goTypes,
		DependencyIndexes: file_users_group_service_proto_depIdxs,
		MessageInfos:      file_users_group_service_proto_msgTypes,
	}.Build()
	File_users_group_service_proto = out.File
	file_users_group_service_proto_rawDesc = nil
	file_users_group_service_proto_goTypes = nil
	file_users_group_service_proto_depIdxs = nil
}
