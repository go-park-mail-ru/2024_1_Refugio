// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: email.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type EmailIdAndLogin struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Login string `protobuf:"bytes,2,opt,name=login,proto3" json:"login,omitempty"`
}

func (x *EmailIdAndLogin) Reset() {
	*x = EmailIdAndLogin{}
	if protoimpl.UnsafeEnabled {
		mi := &file_email_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmailIdAndLogin) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmailIdAndLogin) ProtoMessage() {}

func (x *EmailIdAndLogin) ProtoReflect() protoreflect.Message {
	mi := &file_email_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmailIdAndLogin.ProtoReflect.Descriptor instead.
func (*EmailIdAndLogin) Descriptor() ([]byte, []int) {
	return file_email_proto_rawDescGZIP(), []int{0}
}

func (x *EmailIdAndLogin) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *EmailIdAndLogin) GetLogin() string {
	if x != nil {
		return x.Login
	}
	return ""
}

type LoginOffsetLimit struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Login  string `protobuf:"bytes,1,opt,name=login,proto3" json:"login,omitempty"`
	Offset int64  `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
	Limit  int64  `protobuf:"varint,3,opt,name=limit,proto3" json:"limit,omitempty"`
}

func (x *LoginOffsetLimit) Reset() {
	*x = LoginOffsetLimit{}
	if protoimpl.UnsafeEnabled {
		mi := &file_email_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoginOffsetLimit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginOffsetLimit) ProtoMessage() {}

func (x *LoginOffsetLimit) ProtoReflect() protoreflect.Message {
	mi := &file_email_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginOffsetLimit.ProtoReflect.Descriptor instead.
func (*LoginOffsetLimit) Descriptor() ([]byte, []int) {
	return file_email_proto_rawDescGZIP(), []int{1}
}

func (x *LoginOffsetLimit) GetLogin() string {
	if x != nil {
		return x.Login
	}
	return ""
}

func (x *LoginOffsetLimit) GetOffset() int64 {
	if x != nil {
		return x.Offset
	}
	return 0
}

func (x *LoginOffsetLimit) GetLimit() int64 {
	if x != nil {
		return x.Limit
	}
	return 0
}

type Emails struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Emails []*Email `protobuf:"bytes,1,rep,name=emails,proto3" json:"emails,omitempty"`
}

func (x *Emails) Reset() {
	*x = Emails{}
	if protoimpl.UnsafeEnabled {
		mi := &file_email_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Emails) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Emails) ProtoMessage() {}

func (x *Emails) ProtoReflect() protoreflect.Message {
	mi := &file_email_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Emails.ProtoReflect.Descriptor instead.
func (*Emails) Descriptor() ([]byte, []int) {
	return file_email_proto_rawDescGZIP(), []int{2}
}

func (x *Emails) GetEmails() []*Email {
	if x != nil {
		return x.Emails
	}
	return nil
}

type Email struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             uint64               `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Topic          string               `protobuf:"bytes,2,opt,name=topic,proto3" json:"topic,omitempty"`
	Text           string               `protobuf:"bytes,3,opt,name=text,proto3" json:"text,omitempty"`
	PhotoID        string               `protobuf:"bytes,4,opt,name=photoID,proto3" json:"photoID,omitempty"`
	ReadStatus     bool                 `protobuf:"varint,5,opt,name=readStatus,proto3" json:"readStatus,omitempty"`
	Flag           bool                 `protobuf:"varint,6,opt,name=flag,proto3" json:"flag,omitempty"`
	Deleted        bool                 `protobuf:"varint,7,opt,name=deleted,proto3" json:"deleted,omitempty"`
	DateOfDispatch *timestamp.Timestamp `protobuf:"bytes,8,opt,name=dateOfDispatch,proto3" json:"dateOfDispatch,omitempty"`
	ReplyToEmailID uint64               `protobuf:"varint,9,opt,name=replyToEmailID,proto3" json:"replyToEmailID,omitempty"`
	DraftStatus    bool                 `protobuf:"varint,10,opt,name=draftStatus,proto3" json:"draftStatus,omitempty"`
	SenderEmail    string               `protobuf:"bytes,11,opt,name=senderEmail,proto3" json:"senderEmail,omitempty"`
	RecipientEmail string               `protobuf:"bytes,12,opt,name=recipientEmail,proto3" json:"recipientEmail,omitempty"`
}

func (x *Email) Reset() {
	*x = Email{}
	if protoimpl.UnsafeEnabled {
		mi := &file_email_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Email) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Email) ProtoMessage() {}

func (x *Email) ProtoReflect() protoreflect.Message {
	mi := &file_email_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Email.ProtoReflect.Descriptor instead.
func (*Email) Descriptor() ([]byte, []int) {
	return file_email_proto_rawDescGZIP(), []int{3}
}

func (x *Email) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Email) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *Email) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *Email) GetPhotoID() string {
	if x != nil {
		return x.PhotoID
	}
	return ""
}

func (x *Email) GetReadStatus() bool {
	if x != nil {
		return x.ReadStatus
	}
	return false
}

func (x *Email) GetFlag() bool {
	if x != nil {
		return x.Flag
	}
	return false
}

func (x *Email) GetDeleted() bool {
	if x != nil {
		return x.Deleted
	}
	return false
}

func (x *Email) GetDateOfDispatch() *timestamp.Timestamp {
	if x != nil {
		return x.DateOfDispatch
	}
	return nil
}

func (x *Email) GetReplyToEmailID() uint64 {
	if x != nil {
		return x.ReplyToEmailID
	}
	return 0
}

func (x *Email) GetDraftStatus() bool {
	if x != nil {
		return x.DraftStatus
	}
	return false
}

func (x *Email) GetSenderEmail() string {
	if x != nil {
		return x.SenderEmail
	}
	return ""
}

func (x *Email) GetRecipientEmail() string {
	if x != nil {
		return x.RecipientEmail
	}
	return ""
}

type EmailWithID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Email *Email `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Id    uint64 `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *EmailWithID) Reset() {
	*x = EmailWithID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_email_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmailWithID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmailWithID) ProtoMessage() {}

func (x *EmailWithID) ProtoReflect() protoreflect.Message {
	mi := &file_email_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmailWithID.ProtoReflect.Descriptor instead.
func (*EmailWithID) Descriptor() ([]byte, []int) {
	return file_email_proto_rawDescGZIP(), []int{4}
}

func (x *EmailWithID) GetEmail() *Email {
	if x != nil {
		return x.Email
	}
	return nil
}

func (x *EmailWithID) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type LoginWithID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Login string `protobuf:"bytes,1,opt,name=login,proto3" json:"login,omitempty"`
	Id    uint64 `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *LoginWithID) Reset() {
	*x = LoginWithID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_email_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoginWithID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginWithID) ProtoMessage() {}

func (x *LoginWithID) ProtoReflect() protoreflect.Message {
	mi := &file_email_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginWithID.ProtoReflect.Descriptor instead.
func (*LoginWithID) Descriptor() ([]byte, []int) {
	return file_email_proto_rawDescGZIP(), []int{5}
}

func (x *LoginWithID) GetLogin() string {
	if x != nil {
		return x.Login
	}
	return ""
}

func (x *LoginWithID) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type IdSenderRecipient struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Sender    string `protobuf:"bytes,2,opt,name=sender,proto3" json:"sender,omitempty"`
	Recipient string `protobuf:"bytes,3,opt,name=recipient,proto3" json:"recipient,omitempty"`
}

func (x *IdSenderRecipient) Reset() {
	*x = IdSenderRecipient{}
	if protoimpl.UnsafeEnabled {
		mi := &file_email_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IdSenderRecipient) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IdSenderRecipient) ProtoMessage() {}

func (x *IdSenderRecipient) ProtoReflect() protoreflect.Message {
	mi := &file_email_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IdSenderRecipient.ProtoReflect.Descriptor instead.
func (*IdSenderRecipient) Descriptor() ([]byte, []int) {
	return file_email_proto_rawDescGZIP(), []int{6}
}

func (x *IdSenderRecipient) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *IdSenderRecipient) GetSender() string {
	if x != nil {
		return x.Sender
	}
	return ""
}

func (x *IdSenderRecipient) GetRecipient() string {
	if x != nil {
		return x.Recipient
	}
	return ""
}

type Recipient struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Recipient string `protobuf:"bytes,1,opt,name=recipient,proto3" json:"recipient,omitempty"`
}

func (x *Recipient) Reset() {
	*x = Recipient{}
	if protoimpl.UnsafeEnabled {
		mi := &file_email_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Recipient) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Recipient) ProtoMessage() {}

func (x *Recipient) ProtoReflect() protoreflect.Message {
	mi := &file_email_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Recipient.ProtoReflect.Descriptor instead.
func (*Recipient) Descriptor() ([]byte, []int) {
	return file_email_proto_rawDescGZIP(), []int{7}
}

func (x *Recipient) GetRecipient() string {
	if x != nil {
		return x.Recipient
	}
	return ""
}

type StatusEmail struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status bool `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *StatusEmail) Reset() {
	*x = StatusEmail{}
	if protoimpl.UnsafeEnabled {
		mi := &file_email_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StatusEmail) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StatusEmail) ProtoMessage() {}

func (x *StatusEmail) ProtoReflect() protoreflect.Message {
	mi := &file_email_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StatusEmail.ProtoReflect.Descriptor instead.
func (*StatusEmail) Descriptor() ([]byte, []int) {
	return file_email_proto_rawDescGZIP(), []int{8}
}

func (x *StatusEmail) GetStatus() bool {
	if x != nil {
		return x.Status
	}
	return false
}

type EmptyEmail struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *EmptyEmail) Reset() {
	*x = EmptyEmail{}
	if protoimpl.UnsafeEnabled {
		mi := &file_email_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmptyEmail) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmptyEmail) ProtoMessage() {}

func (x *EmptyEmail) ProtoReflect() protoreflect.Message {
	mi := &file_email_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmptyEmail.ProtoReflect.Descriptor instead.
func (*EmptyEmail) Descriptor() ([]byte, []int) {
	return file_email_proto_rawDescGZIP(), []int{9}
}

var File_email_proto protoreflect.FileDescriptor

var file_email_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x37, 0x0a, 0x0f, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x49, 0x64,
	0x41, 0x6e, 0x64, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x6f, 0x67, 0x69,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x22, 0x56,
	0x0a, 0x10, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x4f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x4c, 0x69, 0x6d,
	0x69, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x66, 0x66, 0x73,
	0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74,
	0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x22, 0x2e, 0x0a, 0x06, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x73,
	0x12, 0x24, 0x0a, 0x06, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x52, 0x06,
	0x65, 0x6d, 0x61, 0x69, 0x6c, 0x73, 0x22, 0x81, 0x03, 0x0a, 0x05, 0x45, 0x6d, 0x61, 0x69, 0x6c,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x68,
	0x6f, 0x74, 0x6f, 0x49, 0x44, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x68, 0x6f,
	0x74, 0x6f, 0x49, 0x44, 0x12, 0x1e, 0x0a, 0x0a, 0x72, 0x65, 0x61, 0x64, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x72, 0x65, 0x61, 0x64, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x6c, 0x61, 0x67, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x04, 0x66, 0x6c, 0x61, 0x67, 0x12, 0x18, 0x0a, 0x07, 0x64, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x64, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x64, 0x12, 0x42, 0x0a, 0x0e, 0x64, 0x61, 0x74, 0x65, 0x4f, 0x66, 0x44, 0x69, 0x73, 0x70,
	0x61, 0x74, 0x63, 0x68, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0e, 0x64, 0x61, 0x74, 0x65, 0x4f, 0x66, 0x44, 0x69,
	0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x12, 0x26, 0x0a, 0x0e, 0x72, 0x65, 0x70, 0x6c, 0x79, 0x54,
	0x6f, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x49, 0x44, 0x18, 0x09, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0e,
	0x72, 0x65, 0x70, 0x6c, 0x79, 0x54, 0x6f, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x49, 0x44, 0x12, 0x20,
	0x0a, 0x0b, 0x64, 0x72, 0x61, 0x66, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x0a, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x0b, 0x64, 0x72, 0x61, 0x66, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x20, 0x0a, 0x0b, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x18,
	0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x45, 0x6d, 0x61,
	0x69, 0x6c, 0x12, 0x26, 0x0a, 0x0e, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x45,
	0x6d, 0x61, 0x69, 0x6c, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x72, 0x65, 0x63, 0x69,
	0x70, 0x69, 0x65, 0x6e, 0x74, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x22, 0x41, 0x0a, 0x0b, 0x45, 0x6d,
	0x61, 0x69, 0x6c, 0x57, 0x69, 0x74, 0x68, 0x49, 0x44, 0x12, 0x22, 0x0a, 0x05, 0x65, 0x6d, 0x61,
	0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x22, 0x33, 0x0a,
	0x0b, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x57, 0x69, 0x74, 0x68, 0x49, 0x44, 0x12, 0x14, 0x0a, 0x05,
	0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6c, 0x6f, 0x67,
	0x69, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02,
	0x69, 0x64, 0x22, 0x59, 0x0a, 0x11, 0x49, 0x64, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x52, 0x65,
	0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x65,
	0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x12,
	0x1c, 0x0a, 0x09, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x22, 0x29, 0x0a,
	0x09, 0x52, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x72, 0x65,
	0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x72,
	0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x22, 0x25, 0x0a, 0x0b, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22,
	0x0c, 0x0a, 0x0a, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x32, 0xdc, 0x03,
	0x0a, 0x0c, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3a,
	0x0a, 0x0e, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x49, 0x6e, 0x63, 0x6f, 0x6d, 0x69, 0x6e, 0x67,
	0x12, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x4f, 0x66,
	0x66, 0x73, 0x65, 0x74, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x1a, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x73, 0x22, 0x00, 0x12, 0x36, 0x0a, 0x0a, 0x47, 0x65,
	0x74, 0x41, 0x6c, 0x6c, 0x53, 0x65, 0x6e, 0x74, 0x12, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x4f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x4c, 0x69, 0x6d, 0x69,
	0x74, 0x1a, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x73,
	0x22, 0x00, 0x12, 0x36, 0x0a, 0x0c, 0x47, 0x65, 0x74, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x42, 0x79,
	0x49, 0x44, 0x12, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x45, 0x6d, 0x61, 0x69, 0x6c,
	0x49, 0x64, 0x41, 0x6e, 0x64, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x1a, 0x0c, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x22, 0x00, 0x12, 0x31, 0x0a, 0x0b, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x1a, 0x12, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x45, 0x6d, 0x61, 0x69, 0x6c, 0x57, 0x69, 0x74, 0x68, 0x49, 0x44, 0x22, 0x00, 0x12, 0x43, 0x0a,
	0x12, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x45, 0x6d,
	0x61, 0x69, 0x6c, 0x12, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x49, 0x64, 0x53, 0x65,
	0x6e, 0x64, 0x65, 0x72, 0x52, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x1a, 0x11, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x45, 0x6d, 0x61, 0x69, 0x6c,
	0x22, 0x00, 0x12, 0x3c, 0x0a, 0x13, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x63, 0x69, 0x70,
	0x69, 0x65, 0x6e, 0x74, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x52, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x1a, 0x11, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x22, 0x00,
	0x12, 0x31, 0x0a, 0x0b, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12,
	0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x1a, 0x12, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x45, 0x6d, 0x61, 0x69,
	0x6c, 0x22, 0x00, 0x12, 0x37, 0x0a, 0x0b, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x45, 0x6d, 0x61,
	0x69, 0x6c, 0x12, 0x12, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e,
	0x57, 0x69, 0x74, 0x68, 0x49, 0x44, 0x1a, 0x12, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x22, 0x00, 0x42, 0x0a, 0x5a, 0x08,
	0x2e, 0x2f, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_email_proto_rawDescOnce sync.Once
	file_email_proto_rawDescData = file_email_proto_rawDesc
)

func file_email_proto_rawDescGZIP() []byte {
	file_email_proto_rawDescOnce.Do(func() {
		file_email_proto_rawDescData = protoimpl.X.CompressGZIP(file_email_proto_rawDescData)
	})
	return file_email_proto_rawDescData
}

var file_email_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_email_proto_goTypes = []interface{}{
	(*EmailIdAndLogin)(nil),     // 0: proto.EmailIdAndLogin
	(*LoginOffsetLimit)(nil),    // 1: proto.LoginOffsetLimit
	(*Emails)(nil),              // 2: proto.Emails
	(*Email)(nil),               // 3: proto.Email
	(*EmailWithID)(nil),         // 4: proto.EmailWithID
	(*LoginWithID)(nil),         // 5: proto.LoginWithID
	(*IdSenderRecipient)(nil),   // 6: proto.IdSenderRecipient
	(*Recipient)(nil),           // 7: proto.Recipient
	(*StatusEmail)(nil),         // 8: proto.StatusEmail
	(*EmptyEmail)(nil),          // 9: proto.EmptyEmail
	(*timestamp.Timestamp)(nil), // 10: google.protobuf.Timestamp
}
var file_email_proto_depIdxs = []int32{
	3,  // 0: proto.Emails.emails:type_name -> proto.Email
	10, // 1: proto.Email.dateOfDispatch:type_name -> google.protobuf.Timestamp
	3,  // 2: proto.EmailWithID.email:type_name -> proto.Email
	1,  // 3: proto.EmailService.GetAllIncoming:input_type -> proto.LoginOffsetLimit
	1,  // 4: proto.EmailService.GetAllSent:input_type -> proto.LoginOffsetLimit
	0,  // 5: proto.EmailService.GetEmailByID:input_type -> proto.EmailIdAndLogin
	3,  // 6: proto.EmailService.CreateEmail:input_type -> proto.Email
	6,  // 7: proto.EmailService.CreateProfileEmail:input_type -> proto.IdSenderRecipient
	7,  // 8: proto.EmailService.CheckRecipientEmail:input_type -> proto.Recipient
	3,  // 9: proto.EmailService.UpdateEmail:input_type -> proto.Email
	5,  // 10: proto.EmailService.DeleteEmail:input_type -> proto.LoginWithID
	2,  // 11: proto.EmailService.GetAllIncoming:output_type -> proto.Emails
	2,  // 12: proto.EmailService.GetAllSent:output_type -> proto.Emails
	3,  // 13: proto.EmailService.GetEmailByID:output_type -> proto.Email
	4,  // 14: proto.EmailService.CreateEmail:output_type -> proto.EmailWithID
	9,  // 15: proto.EmailService.CreateProfileEmail:output_type -> proto.EmptyEmail
	9,  // 16: proto.EmailService.CheckRecipientEmail:output_type -> proto.EmptyEmail
	8,  // 17: proto.EmailService.UpdateEmail:output_type -> proto.StatusEmail
	8,  // 18: proto.EmailService.DeleteEmail:output_type -> proto.StatusEmail
	11, // [11:19] is the sub-list for method output_type
	3,  // [3:11] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file_email_proto_init() }
func file_email_proto_init() {
	if File_email_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_email_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EmailIdAndLogin); i {
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
		file_email_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoginOffsetLimit); i {
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
		file_email_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Emails); i {
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
		file_email_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Email); i {
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
		file_email_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EmailWithID); i {
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
		file_email_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoginWithID); i {
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
		file_email_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IdSenderRecipient); i {
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
		file_email_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Recipient); i {
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
		file_email_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StatusEmail); i {
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
		file_email_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EmptyEmail); i {
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
			RawDescriptor: file_email_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_email_proto_goTypes,
		DependencyIndexes: file_email_proto_depIdxs,
		MessageInfos:      file_email_proto_msgTypes,
	}.Build()
	File_email_proto = out.File
	file_email_proto_rawDesc = nil
	file_email_proto_goTypes = nil
	file_email_proto_depIdxs = nil
}
