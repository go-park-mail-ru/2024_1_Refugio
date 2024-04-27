// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.25.1
// source: question-answer.proto

package proto

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

type Question struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      uint32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Text    string `protobuf:"bytes,2,opt,name=text,proto3" json:"text,omitempty"`
	MinText string `protobuf:"bytes,3,opt,name=min_text,json=minText,proto3" json:"min_text,omitempty"`
	MaxText string `protobuf:"bytes,4,opt,name=max_text,json=maxText,proto3" json:"max_text,omitempty"`
}

func (x *Question) Reset() {
	*x = Question{}
	if protoimpl.UnsafeEnabled {
		mi := &file_question_answer_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Question) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Question) ProtoMessage() {}

func (x *Question) ProtoReflect() protoreflect.Message {
	mi := &file_question_answer_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Question.ProtoReflect.Descriptor instead.
func (*Question) Descriptor() ([]byte, []int) {
	return file_question_answer_proto_rawDescGZIP(), []int{0}
}

func (x *Question) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Question) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *Question) GetMinText() string {
	if x != nil {
		return x.MinText
	}
	return ""
}

func (x *Question) GetMaxText() string {
	if x != nil {
		return x.MaxText
	}
	return ""
}

type Answer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         uint32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	QuestionId uint32 `protobuf:"varint,2,opt,name=question_id,json=questionId,proto3" json:"question_id,omitempty"`
	Login      string `protobuf:"bytes,3,opt,name=login,proto3" json:"login,omitempty"`
	Mark       uint32 `protobuf:"varint,4,opt,name=mark,proto3" json:"mark,omitempty"`
}

func (x *Answer) Reset() {
	*x = Answer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_question_answer_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Answer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Answer) ProtoMessage() {}

func (x *Answer) ProtoReflect() protoreflect.Message {
	mi := &file_question_answer_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Answer.ProtoReflect.Descriptor instead.
func (*Answer) Descriptor() ([]byte, []int) {
	return file_question_answer_proto_rawDescGZIP(), []int{1}
}

func (x *Answer) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Answer) GetQuestionId() uint32 {
	if x != nil {
		return x.QuestionId
	}
	return 0
}

func (x *Answer) GetLogin() string {
	if x != nil {
		return x.Login
	}
	return ""
}

func (x *Answer) GetMark() uint32 {
	if x != nil {
		return x.Mark
	}
	return 0
}

type Statistic struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Text    string `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
	Average uint32 `protobuf:"varint,2,opt,name=average,proto3" json:"average,omitempty"`
}

func (x *Statistic) Reset() {
	*x = Statistic{}
	if protoimpl.UnsafeEnabled {
		mi := &file_question_answer_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Statistic) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Statistic) ProtoMessage() {}

func (x *Statistic) ProtoReflect() protoreflect.Message {
	mi := &file_question_answer_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Statistic.ProtoReflect.Descriptor instead.
func (*Statistic) Descriptor() ([]byte, []int) {
	return file_question_answer_proto_rawDescGZIP(), []int{2}
}

func (x *Statistic) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *Statistic) GetAverage() uint32 {
	if x != nil {
		return x.Average
	}
	return 0
}

type GetQuestionsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetQuestionsRequest) Reset() {
	*x = GetQuestionsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_question_answer_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetQuestionsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetQuestionsRequest) ProtoMessage() {}

func (x *GetQuestionsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_question_answer_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetQuestionsRequest.ProtoReflect.Descriptor instead.
func (*GetQuestionsRequest) Descriptor() ([]byte, []int) {
	return file_question_answer_proto_rawDescGZIP(), []int{3}
}

type GetQuestionsReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Questions []*Question `protobuf:"bytes,1,rep,name=questions,proto3" json:"questions,omitempty"`
}

func (x *GetQuestionsReply) Reset() {
	*x = GetQuestionsReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_question_answer_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetQuestionsReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetQuestionsReply) ProtoMessage() {}

func (x *GetQuestionsReply) ProtoReflect() protoreflect.Message {
	mi := &file_question_answer_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetQuestionsReply.ProtoReflect.Descriptor instead.
func (*GetQuestionsReply) Descriptor() ([]byte, []int) {
	return file_question_answer_proto_rawDescGZIP(), []int{4}
}

func (x *GetQuestionsReply) GetQuestions() []*Question {
	if x != nil {
		return x.Questions
	}
	return nil
}

type AddQuestionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Question *Question `protobuf:"bytes,1,opt,name=question,proto3" json:"question,omitempty"`
}

func (x *AddQuestionRequest) Reset() {
	*x = AddQuestionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_question_answer_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddQuestionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddQuestionRequest) ProtoMessage() {}

func (x *AddQuestionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_question_answer_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddQuestionRequest.ProtoReflect.Descriptor instead.
func (*AddQuestionRequest) Descriptor() ([]byte, []int) {
	return file_question_answer_proto_rawDescGZIP(), []int{5}
}

func (x *AddQuestionRequest) GetQuestion() *Question {
	if x != nil {
		return x.Question
	}
	return nil
}

type AddQuestionReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status bool `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *AddQuestionReply) Reset() {
	*x = AddQuestionReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_question_answer_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddQuestionReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddQuestionReply) ProtoMessage() {}

func (x *AddQuestionReply) ProtoReflect() protoreflect.Message {
	mi := &file_question_answer_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddQuestionReply.ProtoReflect.Descriptor instead.
func (*AddQuestionReply) Descriptor() ([]byte, []int) {
	return file_question_answer_proto_rawDescGZIP(), []int{6}
}

func (x *AddQuestionReply) GetStatus() bool {
	if x != nil {
		return x.Status
	}
	return false
}

type AddAnswerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Answer *Answer `protobuf:"bytes,1,opt,name=answer,proto3" json:"answer,omitempty"`
}

func (x *AddAnswerRequest) Reset() {
	*x = AddAnswerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_question_answer_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddAnswerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddAnswerRequest) ProtoMessage() {}

func (x *AddAnswerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_question_answer_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddAnswerRequest.ProtoReflect.Descriptor instead.
func (*AddAnswerRequest) Descriptor() ([]byte, []int) {
	return file_question_answer_proto_rawDescGZIP(), []int{7}
}

func (x *AddAnswerRequest) GetAnswer() *Answer {
	if x != nil {
		return x.Answer
	}
	return nil
}

type AddAnswerReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status bool `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *AddAnswerReply) Reset() {
	*x = AddAnswerReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_question_answer_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddAnswerReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddAnswerReply) ProtoMessage() {}

func (x *AddAnswerReply) ProtoReflect() protoreflect.Message {
	mi := &file_question_answer_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddAnswerReply.ProtoReflect.Descriptor instead.
func (*AddAnswerReply) Descriptor() ([]byte, []int) {
	return file_question_answer_proto_rawDescGZIP(), []int{8}
}

func (x *AddAnswerReply) GetStatus() bool {
	if x != nil {
		return x.Status
	}
	return false
}

type GetStatisticRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetStatisticRequest) Reset() {
	*x = GetStatisticRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_question_answer_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetStatisticRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetStatisticRequest) ProtoMessage() {}

func (x *GetStatisticRequest) ProtoReflect() protoreflect.Message {
	mi := &file_question_answer_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetStatisticRequest.ProtoReflect.Descriptor instead.
func (*GetStatisticRequest) Descriptor() ([]byte, []int) {
	return file_question_answer_proto_rawDescGZIP(), []int{9}
}

type GetStatisticReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Statistics []*Statistic `protobuf:"bytes,1,rep,name=statistics,proto3" json:"statistics,omitempty"`
}

func (x *GetStatisticReply) Reset() {
	*x = GetStatisticReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_question_answer_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetStatisticReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetStatisticReply) ProtoMessage() {}

func (x *GetStatisticReply) ProtoReflect() protoreflect.Message {
	mi := &file_question_answer_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetStatisticReply.ProtoReflect.Descriptor instead.
func (*GetStatisticReply) Descriptor() ([]byte, []int) {
	return file_question_answer_proto_rawDescGZIP(), []int{10}
}

func (x *GetStatisticReply) GetStatistics() []*Statistic {
	if x != nil {
		return x.Statistics
	}
	return nil
}

var File_question_answer_proto protoreflect.FileDescriptor

var file_question_answer_proto_rawDesc = []byte{
	0x0a, 0x15, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x2d, 0x61, 0x6e, 0x73, 0x77, 0x65,
	0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x64,
	0x0a, 0x08, 0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65,
	0x78, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x12, 0x19,
	0x0a, 0x08, 0x6d, 0x69, 0x6e, 0x5f, 0x74, 0x65, 0x78, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x6d, 0x69, 0x6e, 0x54, 0x65, 0x78, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x6d, 0x61, 0x78,
	0x5f, 0x74, 0x65, 0x78, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x61, 0x78,
	0x54, 0x65, 0x78, 0x74, 0x22, 0x63, 0x0a, 0x06, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1f,
	0x0a, 0x0b, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x0a, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12,
	0x14, 0x0a, 0x05, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x6d, 0x61, 0x72, 0x6b, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x04, 0x6d, 0x61, 0x72, 0x6b, 0x22, 0x39, 0x0a, 0x09, 0x53, 0x74, 0x61,
	0x74, 0x69, 0x73, 0x74, 0x69, 0x63, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x76,
	0x65, 0x72, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x61, 0x76, 0x65,
	0x72, 0x61, 0x67, 0x65, 0x22, 0x15, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x51, 0x75, 0x65, 0x73, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x42, 0x0a, 0x11, 0x47,
	0x65, 0x74, 0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79,
	0x12, 0x2d, 0x0a, 0x09, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x51, 0x75, 0x65, 0x73,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x09, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22,
	0x41, 0x0a, 0x12, 0x41, 0x64, 0x64, 0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2b, 0x0a, 0x08, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x08, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69,
	0x6f, 0x6e, 0x22, 0x2a, 0x0a, 0x10, 0x41, 0x64, 0x64, 0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x39,
	0x0a, 0x10, 0x41, 0x64, 0x64, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x25, 0x0a, 0x06, 0x61, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x6e, 0x73, 0x77, 0x65,
	0x72, 0x52, 0x06, 0x61, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x22, 0x28, 0x0a, 0x0e, 0x41, 0x64, 0x64,
	0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x22, 0x15, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x69, 0x73,
	0x74, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x45, 0x0a, 0x11, 0x47, 0x65,
	0x74, 0x53, 0x74, 0x61, 0x74, 0x69, 0x73, 0x74, 0x69, 0x63, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12,
	0x30, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x69, 0x73, 0x74, 0x69, 0x63, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x74, 0x61, 0x74,
	0x69, 0x73, 0x74, 0x69, 0x63, 0x52, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x69, 0x73, 0x74, 0x69, 0x63,
	0x73, 0x32, 0xa5, 0x02, 0x0a, 0x0f, 0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x46, 0x0a, 0x0c, 0x47, 0x65, 0x74, 0x51, 0x75, 0x65, 0x73,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65,
	0x74, 0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x51, 0x75, 0x65,
	0x73, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x43, 0x0a,
	0x0b, 0x41, 0x64, 0x64, 0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x19, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x64, 0x64, 0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x41, 0x64, 0x64, 0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x70, 0x6c, 0x79,
	0x22, 0x00, 0x12, 0x3d, 0x0a, 0x09, 0x41, 0x64, 0x64, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x12,
	0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x64, 0x64, 0x41, 0x6e, 0x73, 0x77, 0x65,
	0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x41, 0x64, 0x64, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22,
	0x00, 0x12, 0x46, 0x0a, 0x0c, 0x47, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x69, 0x73, 0x74, 0x69,
	0x63, 0x12, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x74, 0x61,
	0x74, 0x69, 0x73, 0x74, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x69, 0x73, 0x74,
	0x69, 0x63, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2f, 0x3b,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_question_answer_proto_rawDescOnce sync.Once
	file_question_answer_proto_rawDescData = file_question_answer_proto_rawDesc
)

func file_question_answer_proto_rawDescGZIP() []byte {
	file_question_answer_proto_rawDescOnce.Do(func() {
		file_question_answer_proto_rawDescData = protoimpl.X.CompressGZIP(file_question_answer_proto_rawDescData)
	})
	return file_question_answer_proto_rawDescData
}

var file_question_answer_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_question_answer_proto_goTypes = []interface{}{
	(*Question)(nil),            // 0: proto.Question
	(*Answer)(nil),              // 1: proto.Answer
	(*Statistic)(nil),           // 2: proto.Statistic
	(*GetQuestionsRequest)(nil), // 3: proto.GetQuestionsRequest
	(*GetQuestionsReply)(nil),   // 4: proto.GetQuestionsReply
	(*AddQuestionRequest)(nil),  // 5: proto.AddQuestionRequest
	(*AddQuestionReply)(nil),    // 6: proto.AddQuestionReply
	(*AddAnswerRequest)(nil),    // 7: proto.AddAnswerRequest
	(*AddAnswerReply)(nil),      // 8: proto.AddAnswerReply
	(*GetStatisticRequest)(nil), // 9: proto.GetStatisticRequest
	(*GetStatisticReply)(nil),   // 10: proto.GetStatisticReply
}
var file_question_answer_proto_depIdxs = []int32{
	0,  // 0: proto.GetQuestionsReply.questions:type_name -> proto.Question
	0,  // 1: proto.AddQuestionRequest.question:type_name -> proto.Question
	1,  // 2: proto.AddAnswerRequest.answer:type_name -> proto.Answer
	2,  // 3: proto.GetStatisticReply.statistics:type_name -> proto.Statistic
	3,  // 4: proto.QuestionService.GetQuestions:input_type -> proto.GetQuestionsRequest
	5,  // 5: proto.QuestionService.AddQuestion:input_type -> proto.AddQuestionRequest
	7,  // 6: proto.QuestionService.AddAnswer:input_type -> proto.AddAnswerRequest
	9,  // 7: proto.QuestionService.GetStatistic:input_type -> proto.GetStatisticRequest
	4,  // 8: proto.QuestionService.GetQuestions:output_type -> proto.GetQuestionsReply
	6,  // 9: proto.QuestionService.AddQuestion:output_type -> proto.AddQuestionReply
	8,  // 10: proto.QuestionService.AddAnswer:output_type -> proto.AddAnswerReply
	10, // 11: proto.QuestionService.GetStatistic:output_type -> proto.GetStatisticReply
	8,  // [8:12] is the sub-list for method output_type
	4,  // [4:8] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_question_answer_proto_init() }
func file_question_answer_proto_init() {
	if File_question_answer_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_question_answer_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Question); i {
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
		file_question_answer_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Answer); i {
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
		file_question_answer_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Statistic); i {
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
		file_question_answer_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetQuestionsRequest); i {
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
		file_question_answer_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetQuestionsReply); i {
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
		file_question_answer_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddQuestionRequest); i {
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
		file_question_answer_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddQuestionReply); i {
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
		file_question_answer_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddAnswerRequest); i {
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
		file_question_answer_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddAnswerReply); i {
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
		file_question_answer_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetStatisticRequest); i {
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
		file_question_answer_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetStatisticReply); i {
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
			RawDescriptor: file_question_answer_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_question_answer_proto_goTypes,
		DependencyIndexes: file_question_answer_proto_depIdxs,
		MessageInfos:      file_question_answer_proto_msgTypes,
	}.Build()
	File_question_answer_proto = out.File
	file_question_answer_proto_rawDesc = nil
	file_question_answer_proto_goTypes = nil
	file_question_answer_proto_depIdxs = nil
}
