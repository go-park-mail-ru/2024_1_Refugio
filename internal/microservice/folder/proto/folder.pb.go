// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.25.1
// source: folder.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
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

type Folder struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        uint32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	ProfileId uint32 `protobuf:"varint,2,opt,name=profileId,proto3" json:"profileId,omitempty"`
	Name      string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *Folder) Reset() {
	*x = Folder{}
	if protoimpl.UnsafeEnabled {
		mi := &file_folder_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Folder) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Folder) ProtoMessage() {}

func (x *Folder) ProtoReflect() protoreflect.Message {
	mi := &file_folder_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Folder.ProtoReflect.Descriptor instead.
func (*Folder) Descriptor() ([]byte, []int) {
	return file_folder_proto_rawDescGZIP(), []int{0}
}

func (x *Folder) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Folder) GetProfileId() uint32 {
	if x != nil {
		return x.ProfileId
	}
	return 0
}

func (x *Folder) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type Folders struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Folders []*Folder `protobuf:"bytes,1,rep,name=folders,proto3" json:"folders,omitempty"`
}

func (x *Folders) Reset() {
	*x = Folders{}
	if protoimpl.UnsafeEnabled {
		mi := &file_folder_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Folders) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Folders) ProtoMessage() {}

func (x *Folders) ProtoReflect() protoreflect.Message {
	mi := &file_folder_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Folders.ProtoReflect.Descriptor instead.
func (*Folders) Descriptor() ([]byte, []int) {
	return file_folder_proto_rawDescGZIP(), []int{1}
}

func (x *Folders) GetFolders() []*Folder {
	if x != nil {
		return x.Folders
	}
	return nil
}

type FolderWithID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Folder *Folder `protobuf:"bytes,1,opt,name=folder,proto3" json:"folder,omitempty"`
	Id     uint32  `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *FolderWithID) Reset() {
	*x = FolderWithID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_folder_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FolderWithID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FolderWithID) ProtoMessage() {}

func (x *FolderWithID) ProtoReflect() protoreflect.Message {
	mi := &file_folder_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FolderWithID.ProtoReflect.Descriptor instead.
func (*FolderWithID) Descriptor() ([]byte, []int) {
	return file_folder_proto_rawDescGZIP(), []int{2}
}

func (x *FolderWithID) GetFolder() *Folder {
	if x != nil {
		return x.Folder
	}
	return nil
}

func (x *FolderWithID) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

type FolderStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status bool `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *FolderStatus) Reset() {
	*x = FolderStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_folder_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FolderStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FolderStatus) ProtoMessage() {}

func (x *FolderStatus) ProtoReflect() protoreflect.Message {
	mi := &file_folder_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FolderStatus.ProtoReflect.Descriptor instead.
func (*FolderStatus) Descriptor() ([]byte, []int) {
	return file_folder_proto_rawDescGZIP(), []int{3}
}

func (x *FolderStatus) GetStatus() bool {
	if x != nil {
		return x.Status
	}
	return false
}

type GetAllFoldersData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     uint32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Limit  int64  `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
	Offset int64  `protobuf:"varint,3,opt,name=offset,proto3" json:"offset,omitempty"`
}

func (x *GetAllFoldersData) Reset() {
	*x = GetAllFoldersData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_folder_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAllFoldersData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAllFoldersData) ProtoMessage() {}

func (x *GetAllFoldersData) ProtoReflect() protoreflect.Message {
	mi := &file_folder_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAllFoldersData.ProtoReflect.Descriptor instead.
func (*GetAllFoldersData) Descriptor() ([]byte, []int) {
	return file_folder_proto_rawDescGZIP(), []int{4}
}

func (x *GetAllFoldersData) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *GetAllFoldersData) GetLimit() int64 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *GetAllFoldersData) GetOffset() int64 {
	if x != nil {
		return x.Offset
	}
	return 0
}

type DeleteFolderData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FolderID  uint32 `protobuf:"varint,1,opt,name=folderID,proto3" json:"folderID,omitempty"`
	ProfileID uint32 `protobuf:"varint,2,opt,name=profileID,proto3" json:"profileID,omitempty"`
}

func (x *DeleteFolderData) Reset() {
	*x = DeleteFolderData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_folder_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteFolderData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteFolderData) ProtoMessage() {}

func (x *DeleteFolderData) ProtoReflect() protoreflect.Message {
	mi := &file_folder_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteFolderData.ProtoReflect.Descriptor instead.
func (*DeleteFolderData) Descriptor() ([]byte, []int) {
	return file_folder_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteFolderData) GetFolderID() uint32 {
	if x != nil {
		return x.FolderID
	}
	return 0
}

func (x *DeleteFolderData) GetProfileID() uint32 {
	if x != nil {
		return x.ProfileID
	}
	return 0
}

type FolderEmail struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FolderID uint32 `protobuf:"varint,1,opt,name=folderID,proto3" json:"folderID,omitempty"`
	EmailID  uint32 `protobuf:"varint,2,opt,name=emailID,proto3" json:"emailID,omitempty"`
}

func (x *FolderEmail) Reset() {
	*x = FolderEmail{}
	if protoimpl.UnsafeEnabled {
		mi := &file_folder_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FolderEmail) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FolderEmail) ProtoMessage() {}

func (x *FolderEmail) ProtoReflect() protoreflect.Message {
	mi := &file_folder_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FolderEmail.ProtoReflect.Descriptor instead.
func (*FolderEmail) Descriptor() ([]byte, []int) {
	return file_folder_proto_rawDescGZIP(), []int{6}
}

func (x *FolderEmail) GetFolderID() uint32 {
	if x != nil {
		return x.FolderID
	}
	return 0
}

func (x *FolderEmail) GetEmailID() uint32 {
	if x != nil {
		return x.EmailID
	}
	return 0
}

type FolderEmailStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status bool `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *FolderEmailStatus) Reset() {
	*x = FolderEmailStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_folder_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FolderEmailStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FolderEmailStatus) ProtoMessage() {}

func (x *FolderEmailStatus) ProtoReflect() protoreflect.Message {
	mi := &file_folder_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FolderEmailStatus.ProtoReflect.Descriptor instead.
func (*FolderEmailStatus) Descriptor() ([]byte, []int) {
	return file_folder_proto_rawDescGZIP(), []int{7}
}

func (x *FolderEmailStatus) GetStatus() bool {
	if x != nil {
		return x.Status
	}
	return false
}

type FolderProfile struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FolderID  uint32 `protobuf:"varint,1,opt,name=folderID,proto3" json:"folderID,omitempty"`
	ProfileID uint32 `protobuf:"varint,2,opt,name=profileID,proto3" json:"profileID,omitempty"`
}

func (x *FolderProfile) Reset() {
	*x = FolderProfile{}
	if protoimpl.UnsafeEnabled {
		mi := &file_folder_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FolderProfile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FolderProfile) ProtoMessage() {}

func (x *FolderProfile) ProtoReflect() protoreflect.Message {
	mi := &file_folder_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FolderProfile.ProtoReflect.Descriptor instead.
func (*FolderProfile) Descriptor() ([]byte, []int) {
	return file_folder_proto_rawDescGZIP(), []int{8}
}

func (x *FolderProfile) GetFolderID() uint32 {
	if x != nil {
		return x.FolderID
	}
	return 0
}

func (x *FolderProfile) GetProfileID() uint32 {
	if x != nil {
		return x.ProfileID
	}
	return 0
}

type GetAllEmailsInFolderData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FolderID  uint32 `protobuf:"varint,1,opt,name=folderID,proto3" json:"folderID,omitempty"`
	ProfileID uint32 `protobuf:"varint,2,opt,name=profileID,proto3" json:"profileID,omitempty"`
	Limit     uint32 `protobuf:"varint,3,opt,name=limit,proto3" json:"limit,omitempty"`
	Offset    uint32 `protobuf:"varint,4,opt,name=offset,proto3" json:"offset,omitempty"`
}

func (x *GetAllEmailsInFolderData) Reset() {
	*x = GetAllEmailsInFolderData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_folder_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAllEmailsInFolderData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAllEmailsInFolderData) ProtoMessage() {}

func (x *GetAllEmailsInFolderData) ProtoReflect() protoreflect.Message {
	mi := &file_folder_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAllEmailsInFolderData.ProtoReflect.Descriptor instead.
func (*GetAllEmailsInFolderData) Descriptor() ([]byte, []int) {
	return file_folder_proto_rawDescGZIP(), []int{9}
}

func (x *GetAllEmailsInFolderData) GetFolderID() uint32 {
	if x != nil {
		return x.FolderID
	}
	return 0
}

func (x *GetAllEmailsInFolderData) GetProfileID() uint32 {
	if x != nil {
		return x.ProfileID
	}
	return 0
}

func (x *GetAllEmailsInFolderData) GetLimit() uint32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *GetAllEmailsInFolderData) GetOffset() uint32 {
	if x != nil {
		return x.Offset
	}
	return 0
}

type EmailProfile struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EmailID   uint32 `protobuf:"varint,1,opt,name=emailID,proto3" json:"emailID,omitempty"`
	ProfileID uint32 `protobuf:"varint,2,opt,name=profileID,proto3" json:"profileID,omitempty"`
}

func (x *EmailProfile) Reset() {
	*x = EmailProfile{}
	if protoimpl.UnsafeEnabled {
		mi := &file_folder_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmailProfile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmailProfile) ProtoMessage() {}

func (x *EmailProfile) ProtoReflect() protoreflect.Message {
	mi := &file_folder_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmailProfile.ProtoReflect.Descriptor instead.
func (*EmailProfile) Descriptor() ([]byte, []int) {
	return file_folder_proto_rawDescGZIP(), []int{10}
}

func (x *EmailProfile) GetEmailID() uint32 {
	if x != nil {
		return x.EmailID
	}
	return 0
}

func (x *EmailProfile) GetProfileID() uint32 {
	if x != nil {
		return x.ProfileID
	}
	return 0
}

type ObjectsEmail struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Emails []*ObjectEmail `protobuf:"bytes,1,rep,name=emails,proto3" json:"emails,omitempty"`
}

func (x *ObjectsEmail) Reset() {
	*x = ObjectsEmail{}
	if protoimpl.UnsafeEnabled {
		mi := &file_folder_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ObjectsEmail) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ObjectsEmail) ProtoMessage() {}

func (x *ObjectsEmail) ProtoReflect() protoreflect.Message {
	mi := &file_folder_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ObjectsEmail.ProtoReflect.Descriptor instead.
func (*ObjectsEmail) Descriptor() ([]byte, []int) {
	return file_folder_proto_rawDescGZIP(), []int{11}
}

func (x *ObjectsEmail) GetEmails() []*ObjectEmail {
	if x != nil {
		return x.Emails
	}
	return nil
}

type ObjectEmail struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Topic          string                 `protobuf:"bytes,2,opt,name=topic,proto3" json:"topic,omitempty"`
	Text           string                 `protobuf:"bytes,3,opt,name=text,proto3" json:"text,omitempty"`
	PhotoID        string                 `protobuf:"bytes,4,opt,name=photoID,proto3" json:"photoID,omitempty"`
	ReadStatus     bool                   `protobuf:"varint,5,opt,name=readStatus,proto3" json:"readStatus,omitempty"`
	Flag           bool                   `protobuf:"varint,6,opt,name=flag,proto3" json:"flag,omitempty"`
	Deleted        bool                   `protobuf:"varint,7,opt,name=deleted,proto3" json:"deleted,omitempty"`
	DateOfDispatch *timestamppb.Timestamp `protobuf:"bytes,8,opt,name=dateOfDispatch,proto3" json:"dateOfDispatch,omitempty"`
	ReplyToEmailID uint64                 `protobuf:"varint,9,opt,name=replyToEmailID,proto3" json:"replyToEmailID,omitempty"`
	DraftStatus    bool                   `protobuf:"varint,10,opt,name=draftStatus,proto3" json:"draftStatus,omitempty"`
	SpamStatus     bool                   `protobuf:"varint,11,opt,name=spamStatus,proto3" json:"spamStatus,omitempty"`
	SenderEmail    string                 `protobuf:"bytes,12,opt,name=senderEmail,proto3" json:"senderEmail,omitempty"`
	RecipientEmail string                 `protobuf:"bytes,13,opt,name=recipientEmail,proto3" json:"recipientEmail,omitempty"`
}

func (x *ObjectEmail) Reset() {
	*x = ObjectEmail{}
	if protoimpl.UnsafeEnabled {
		mi := &file_folder_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ObjectEmail) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ObjectEmail) ProtoMessage() {}

func (x *ObjectEmail) ProtoReflect() protoreflect.Message {
	mi := &file_folder_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ObjectEmail.ProtoReflect.Descriptor instead.
func (*ObjectEmail) Descriptor() ([]byte, []int) {
	return file_folder_proto_rawDescGZIP(), []int{12}
}

func (x *ObjectEmail) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *ObjectEmail) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *ObjectEmail) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *ObjectEmail) GetPhotoID() string {
	if x != nil {
		return x.PhotoID
	}
	return ""
}

func (x *ObjectEmail) GetReadStatus() bool {
	if x != nil {
		return x.ReadStatus
	}
	return false
}

func (x *ObjectEmail) GetFlag() bool {
	if x != nil {
		return x.Flag
	}
	return false
}

func (x *ObjectEmail) GetDeleted() bool {
	if x != nil {
		return x.Deleted
	}
	return false
}

func (x *ObjectEmail) GetDateOfDispatch() *timestamppb.Timestamp {
	if x != nil {
		return x.DateOfDispatch
	}
	return nil
}

func (x *ObjectEmail) GetReplyToEmailID() uint64 {
	if x != nil {
		return x.ReplyToEmailID
	}
	return 0
}

func (x *ObjectEmail) GetDraftStatus() bool {
	if x != nil {
		return x.DraftStatus
	}
	return false
}

func (x *ObjectEmail) GetSpamStatus() bool {
	if x != nil {
		return x.SpamStatus
	}
	return false
}

func (x *ObjectEmail) GetSenderEmail() string {
	if x != nil {
		return x.SenderEmail
	}
	return ""
}

func (x *ObjectEmail) GetRecipientEmail() string {
	if x != nil {
		return x.RecipientEmail
	}
	return ""
}

var File_folder_proto protoreflect.FileDescriptor

var file_folder_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x66, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4a, 0x0a, 0x06, 0x46, 0x6f, 0x6c, 0x64, 0x65, 0x72,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x1c, 0x0a, 0x09, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x64, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x22, 0x32, 0x0a, 0x07, 0x46, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x73, 0x12, 0x27, 0x0a,
	0x07, 0x66, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0d,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x46, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x52, 0x07, 0x66,
	0x6f, 0x6c, 0x64, 0x65, 0x72, 0x73, 0x22, 0x45, 0x0a, 0x0c, 0x46, 0x6f, 0x6c, 0x64, 0x65, 0x72,
	0x57, 0x69, 0x74, 0x68, 0x49, 0x44, 0x12, 0x25, 0x0a, 0x06, 0x66, 0x6f, 0x6c, 0x64, 0x65, 0x72,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x46,
	0x6f, 0x6c, 0x64, 0x65, 0x72, 0x52, 0x06, 0x66, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x69, 0x64, 0x22, 0x26, 0x0a,
	0x0c, 0x46, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x16, 0x0a,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x51, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x46,
	0x6f, 0x6c, 0x64, 0x65, 0x72, 0x73, 0x44, 0x61, 0x74, 0x61, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69,
	0x6d, 0x69, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74,
	0x12, 0x16, 0x0a, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x22, 0x4c, 0x0a, 0x10, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x46, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1a, 0x0a, 0x08,
	0x66, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08,
	0x66, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x49, 0x44, 0x12, 0x1c, 0x0a, 0x09, 0x70, 0x72, 0x6f, 0x66,
	0x69, 0x6c, 0x65, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x70, 0x72, 0x6f,
	0x66, 0x69, 0x6c, 0x65, 0x49, 0x44, 0x22, 0x43, 0x0a, 0x0b, 0x46, 0x6f, 0x6c, 0x64, 0x65, 0x72,
	0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x49,
	0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x66, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x49,
	0x44, 0x12, 0x18, 0x0a, 0x07, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x07, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x49, 0x44, 0x22, 0x2b, 0x0a, 0x11, 0x46,
	0x6f, 0x6c, 0x64, 0x65, 0x72, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x49, 0x0a, 0x0d, 0x46, 0x6f, 0x6c, 0x64,
	0x65, 0x72, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x6f, 0x6c,
	0x64, 0x65, 0x72, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x66, 0x6f, 0x6c,
	0x64, 0x65, 0x72, 0x49, 0x44, 0x12, 0x1c, 0x0a, 0x09, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65,
	0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c,
	0x65, 0x49, 0x44, 0x22, 0x82, 0x01, 0x0a, 0x18, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x45, 0x6d,
	0x61, 0x69, 0x6c, 0x73, 0x49, 0x6e, 0x46, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61,
	0x12, 0x1a, 0x0a, 0x08, 0x66, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x08, 0x66, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x49, 0x44, 0x12, 0x1c, 0x0a, 0x09,
	0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x09, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x44, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69,
	0x6d, 0x69, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74,
	0x12, 0x16, 0x0a, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x22, 0x46, 0x0a, 0x0c, 0x45, 0x6d, 0x61, 0x69,
	0x6c, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x65, 0x6d, 0x61, 0x69,
	0x6c, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x65, 0x6d, 0x61, 0x69, 0x6c,
	0x49, 0x44, 0x12, 0x1c, 0x0a, 0x09, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x44, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x44,
	0x22, 0x3a, 0x0a, 0x0c, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x45, 0x6d, 0x61, 0x69, 0x6c,
	0x12, 0x2a, 0x0a, 0x06, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x12, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x45,
	0x6d, 0x61, 0x69, 0x6c, 0x52, 0x06, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x73, 0x22, 0xa7, 0x03, 0x0a,
	0x0b, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05,
	0x74, 0x6f, 0x70, 0x69, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x70,
	0x69, 0x63, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x68, 0x6f, 0x74, 0x6f, 0x49,
	0x44, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x68, 0x6f, 0x74, 0x6f, 0x49, 0x44,
	0x12, 0x1e, 0x0a, 0x0a, 0x72, 0x65, 0x61, 0x64, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x72, 0x65, 0x61, 0x64, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x12, 0x0a, 0x04, 0x66, 0x6c, 0x61, 0x67, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04,
	0x66, 0x6c, 0x61, 0x67, 0x12, 0x18, 0x0a, 0x07, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x12, 0x42,
	0x0a, 0x0e, 0x64, 0x61, 0x74, 0x65, 0x4f, 0x66, 0x44, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x0e, 0x64, 0x61, 0x74, 0x65, 0x4f, 0x66, 0x44, 0x69, 0x73, 0x70, 0x61, 0x74,
	0x63, 0x68, 0x12, 0x26, 0x0a, 0x0e, 0x72, 0x65, 0x70, 0x6c, 0x79, 0x54, 0x6f, 0x45, 0x6d, 0x61,
	0x69, 0x6c, 0x49, 0x44, 0x18, 0x09, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0e, 0x72, 0x65, 0x70, 0x6c,
	0x79, 0x54, 0x6f, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x49, 0x44, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x72,
	0x61, 0x66, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x0b, 0x64, 0x72, 0x61, 0x66, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1e, 0x0a, 0x0a,
	0x73, 0x70, 0x61, 0x6d, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x0a, 0x73, 0x70, 0x61, 0x6d, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x20, 0x0a, 0x0b,
	0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x0c, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x26,
	0x0a, 0x0e, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x45, 0x6d, 0x61, 0x69, 0x6c,
	0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e,
	0x74, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x32, 0xe1, 0x04, 0x0a, 0x0d, 0x46, 0x6f, 0x6c, 0x64, 0x65,
	0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x34, 0x0a, 0x0c, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x46, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x12, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x46, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x1a, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x46, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x57, 0x69, 0x74, 0x68, 0x49, 0x44, 0x22, 0x00, 0x12, 0x3b,
	0x0a, 0x0d, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x46, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x73, 0x12,
	0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x46, 0x6f,
	0x6c, 0x64, 0x65, 0x72, 0x73, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x46, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x73, 0x22, 0x00, 0x12, 0x34, 0x0a, 0x0c, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x46, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x12, 0x0d, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x46, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x1a, 0x13, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x46, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22,
	0x00, 0x12, 0x3e, 0x0a, 0x0c, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x46, 0x6f, 0x6c, 0x64, 0x65,
	0x72, 0x12, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x46, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x13, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x46, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22,
	0x00, 0x12, 0x42, 0x0a, 0x10, 0x41, 0x64, 0x64, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x49, 0x6e, 0x46,
	0x6f, 0x6c, 0x64, 0x65, 0x72, 0x12, 0x12, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x46, 0x6f,
	0x6c, 0x64, 0x65, 0x72, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x1a, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x46, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x22, 0x00, 0x12, 0x45, 0x0a, 0x13, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x45,
	0x6d, 0x61, 0x69, 0x6c, 0x49, 0x6e, 0x46, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x12, 0x12, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x46, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x45, 0x6d, 0x61, 0x69, 0x6c,
	0x1a, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x46, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x45,
	0x6d, 0x61, 0x69, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x00, 0x12, 0x4e, 0x0a, 0x14,
	0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x73, 0x49, 0x6e, 0x46, 0x6f,
	0x6c, 0x64, 0x65, 0x72, 0x12, 0x1f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74,
	0x41, 0x6c, 0x6c, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x73, 0x49, 0x6e, 0x46, 0x6f, 0x6c, 0x64, 0x65,
	0x72, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4f, 0x62,
	0x6a, 0x65, 0x63, 0x74, 0x73, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x22, 0x00, 0x12, 0x46, 0x0a, 0x12,
	0x43, 0x68, 0x65, 0x63, 0x6b, 0x46, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x66, 0x69,
	0x6c, 0x65, 0x12, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x46, 0x6f, 0x6c, 0x64, 0x65,
	0x72, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x1a, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x46, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x22, 0x00, 0x12, 0x44, 0x0a, 0x11, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x45, 0x6d, 0x61,
	0x69, 0x6c, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x12, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x1a, 0x18,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x46, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x45, 0x6d, 0x61,
	0x69, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x00, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2f,
	0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_folder_proto_rawDescOnce sync.Once
	file_folder_proto_rawDescData = file_folder_proto_rawDesc
)

func file_folder_proto_rawDescGZIP() []byte {
	file_folder_proto_rawDescOnce.Do(func() {
		file_folder_proto_rawDescData = protoimpl.X.CompressGZIP(file_folder_proto_rawDescData)
	})
	return file_folder_proto_rawDescData
}

var file_folder_proto_msgTypes = make([]protoimpl.MessageInfo, 13)
var file_folder_proto_goTypes = []interface{}{
	(*Folder)(nil),                   // 0: proto.Folder
	(*Folders)(nil),                  // 1: proto.Folders
	(*FolderWithID)(nil),             // 2: proto.FolderWithID
	(*FolderStatus)(nil),             // 3: proto.FolderStatus
	(*GetAllFoldersData)(nil),        // 4: proto.GetAllFoldersData
	(*DeleteFolderData)(nil),         // 5: proto.DeleteFolderData
	(*FolderEmail)(nil),              // 6: proto.FolderEmail
	(*FolderEmailStatus)(nil),        // 7: proto.FolderEmailStatus
	(*FolderProfile)(nil),            // 8: proto.FolderProfile
	(*GetAllEmailsInFolderData)(nil), // 9: proto.GetAllEmailsInFolderData
	(*EmailProfile)(nil),             // 10: proto.EmailProfile
	(*ObjectsEmail)(nil),             // 11: proto.ObjectsEmail
	(*ObjectEmail)(nil),              // 12: proto.ObjectEmail
	(*timestamppb.Timestamp)(nil),    // 13: google.protobuf.Timestamp
}
var file_folder_proto_depIdxs = []int32{
	0,  // 0: proto.Folders.folders:type_name -> proto.Folder
	0,  // 1: proto.FolderWithID.folder:type_name -> proto.Folder
	12, // 2: proto.ObjectsEmail.emails:type_name -> proto.ObjectEmail
	13, // 3: proto.ObjectEmail.dateOfDispatch:type_name -> google.protobuf.Timestamp
	0,  // 4: proto.FolderService.CreateFolder:input_type -> proto.Folder
	4,  // 5: proto.FolderService.GetAllFolders:input_type -> proto.GetAllFoldersData
	0,  // 6: proto.FolderService.UpdateFolder:input_type -> proto.Folder
	5,  // 7: proto.FolderService.DeleteFolder:input_type -> proto.DeleteFolderData
	6,  // 8: proto.FolderService.AddEmailInFolder:input_type -> proto.FolderEmail
	6,  // 9: proto.FolderService.DeleteEmailInFolder:input_type -> proto.FolderEmail
	9,  // 10: proto.FolderService.GetAllEmailsInFolder:input_type -> proto.GetAllEmailsInFolderData
	8,  // 11: proto.FolderService.CheckFolderProfile:input_type -> proto.FolderProfile
	10, // 12: proto.FolderService.CheckEmailProfile:input_type -> proto.EmailProfile
	2,  // 13: proto.FolderService.CreateFolder:output_type -> proto.FolderWithID
	1,  // 14: proto.FolderService.GetAllFolders:output_type -> proto.Folders
	3,  // 15: proto.FolderService.UpdateFolder:output_type -> proto.FolderStatus
	3,  // 16: proto.FolderService.DeleteFolder:output_type -> proto.FolderStatus
	7,  // 17: proto.FolderService.AddEmailInFolder:output_type -> proto.FolderEmailStatus
	7,  // 18: proto.FolderService.DeleteEmailInFolder:output_type -> proto.FolderEmailStatus
	11, // 19: proto.FolderService.GetAllEmailsInFolder:output_type -> proto.ObjectsEmail
	7,  // 20: proto.FolderService.CheckFolderProfile:output_type -> proto.FolderEmailStatus
	7,  // 21: proto.FolderService.CheckEmailProfile:output_type -> proto.FolderEmailStatus
	13, // [13:22] is the sub-list for method output_type
	4,  // [4:13] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_folder_proto_init() }
func file_folder_proto_init() {
	if File_folder_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_folder_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Folder); i {
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
		file_folder_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Folders); i {
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
		file_folder_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FolderWithID); i {
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
		file_folder_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FolderStatus); i {
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
		file_folder_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAllFoldersData); i {
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
		file_folder_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteFolderData); i {
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
		file_folder_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FolderEmail); i {
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
		file_folder_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FolderEmailStatus); i {
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
		file_folder_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FolderProfile); i {
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
		file_folder_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAllEmailsInFolderData); i {
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
		file_folder_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EmailProfile); i {
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
		file_folder_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ObjectsEmail); i {
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
		file_folder_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ObjectEmail); i {
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
			RawDescriptor: file_folder_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   13,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_folder_proto_goTypes,
		DependencyIndexes: file_folder_proto_depIdxs,
		MessageInfos:      file_folder_proto_msgTypes,
	}.Build()
	File_folder_proto = out.File
	file_folder_proto_rawDesc = nil
	file_folder_proto_goTypes = nil
	file_folder_proto_depIdxs = nil
}
