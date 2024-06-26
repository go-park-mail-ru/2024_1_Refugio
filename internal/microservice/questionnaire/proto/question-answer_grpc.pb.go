//go:generate mockgen -source=./question-answer_grpc.pb.go -destination=../mock/question-answer_grpc_mock.go -package=mock proto QuestionServiceClient

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.1
// source: question-answer.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// QuestionServiceClient is the client API for QuestionService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type QuestionServiceClient interface {
	GetQuestions(ctx context.Context, in *GetQuestionsRequest, opts ...grpc.CallOption) (*GetQuestionsReply, error)
	AddQuestion(ctx context.Context, in *AddQuestionRequest, opts ...grpc.CallOption) (*AddQuestionReply, error)
	AddAnswer(ctx context.Context, in *AddAnswerRequest, opts ...grpc.CallOption) (*AddAnswerReply, error)
	GetStatistic(ctx context.Context, in *GetStatisticRequest, opts ...grpc.CallOption) (*GetStatisticReply, error)
}

type questionServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewQuestionServiceClient(cc grpc.ClientConnInterface) QuestionServiceClient {
	return &questionServiceClient{cc}
}

func (c *questionServiceClient) GetQuestions(ctx context.Context, in *GetQuestionsRequest, opts ...grpc.CallOption) (*GetQuestionsReply, error) {
	out := new(GetQuestionsReply)
	err := c.cc.Invoke(ctx, "/proto.QuestionService/GetQuestions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *questionServiceClient) AddQuestion(ctx context.Context, in *AddQuestionRequest, opts ...grpc.CallOption) (*AddQuestionReply, error) {
	out := new(AddQuestionReply)
	err := c.cc.Invoke(ctx, "/proto.QuestionService/AddQuestion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *questionServiceClient) AddAnswer(ctx context.Context, in *AddAnswerRequest, opts ...grpc.CallOption) (*AddAnswerReply, error) {
	out := new(AddAnswerReply)
	err := c.cc.Invoke(ctx, "/proto.QuestionService/AddAnswer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *questionServiceClient) GetStatistic(ctx context.Context, in *GetStatisticRequest, opts ...grpc.CallOption) (*GetStatisticReply, error) {
	out := new(GetStatisticReply)
	err := c.cc.Invoke(ctx, "/proto.QuestionService/GetStatistic", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QuestionServiceServer is the server API for QuestionService service.
// All implementations must embed UnimplementedQuestionServiceServer
// for forward compatibility
type QuestionServiceServer interface {
	GetQuestions(context.Context, *GetQuestionsRequest) (*GetQuestionsReply, error)
	AddQuestion(context.Context, *AddQuestionRequest) (*AddQuestionReply, error)
	AddAnswer(context.Context, *AddAnswerRequest) (*AddAnswerReply, error)
	GetStatistic(context.Context, *GetStatisticRequest) (*GetStatisticReply, error)
	mustEmbedUnimplementedQuestionServiceServer()
}

// UnimplementedQuestionServiceServer must be embedded to have forward compatible implementations.
type UnimplementedQuestionServiceServer struct {
}

func (UnimplementedQuestionServiceServer) GetQuestions(context.Context, *GetQuestionsRequest) (*GetQuestionsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetQuestions not implemented")
}
func (UnimplementedQuestionServiceServer) AddQuestion(context.Context, *AddQuestionRequest) (*AddQuestionReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddQuestion not implemented")
}
func (UnimplementedQuestionServiceServer) AddAnswer(context.Context, *AddAnswerRequest) (*AddAnswerReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddAnswer not implemented")
}
func (UnimplementedQuestionServiceServer) GetStatistic(context.Context, *GetStatisticRequest) (*GetStatisticReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStatistic not implemented")
}
func (UnimplementedQuestionServiceServer) mustEmbedUnimplementedQuestionServiceServer() {}

// UnsafeQuestionServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to QuestionServiceServer will
// result in compilation errors.
type UnsafeQuestionServiceServer interface {
	mustEmbedUnimplementedQuestionServiceServer()
}

func RegisterQuestionServiceServer(s grpc.ServiceRegistrar, srv QuestionServiceServer) {
	s.RegisterService(&QuestionService_ServiceDesc, srv)
}

func _QuestionService_GetQuestions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetQuestionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QuestionServiceServer).GetQuestions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.QuestionService/GetQuestions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QuestionServiceServer).GetQuestions(ctx, req.(*GetQuestionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _QuestionService_AddQuestion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddQuestionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QuestionServiceServer).AddQuestion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.QuestionService/AddQuestion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QuestionServiceServer).AddQuestion(ctx, req.(*AddQuestionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _QuestionService_AddAnswer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddAnswerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QuestionServiceServer).AddAnswer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.QuestionService/AddAnswer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QuestionServiceServer).AddAnswer(ctx, req.(*AddAnswerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _QuestionService_GetStatistic_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStatisticRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QuestionServiceServer).GetStatistic(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.QuestionService/GetStatistic",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QuestionServiceServer).GetStatistic(ctx, req.(*GetStatisticRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// QuestionService_ServiceDesc is the grpc.ServiceDesc for QuestionService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var QuestionService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.QuestionService",
	HandlerType: (*QuestionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetQuestions",
			Handler:    _QuestionService_GetQuestions_Handler,
		},
		{
			MethodName: "AddQuestion",
			Handler:    _QuestionService_AddQuestion_Handler,
		},
		{
			MethodName: "AddAnswer",
			Handler:    _QuestionService_AddAnswer_Handler,
		},
		{
			MethodName: "GetStatistic",
			Handler:    _QuestionService_GetStatistic_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "question-answer.proto",
}
