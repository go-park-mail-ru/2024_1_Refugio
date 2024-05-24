package server

import (
	"context"
	"fmt"

	"mail/internal/microservice/questionnaire/interface"
	"mail/internal/microservice/questionnaire/proto"

	converters "mail/internal/microservice/models/proto_converters"
)

// QuestionAnswerServer handles RPC calls for the QuestionService.
type QuestionAnswerServer struct {
	proto.UnimplementedQuestionServiceServer
	QuestionAnswerUseCase _interface.QuestionAnswerUseCase
}

// NewQuestionAnswerServer creates a new instance of QuestionAnswerServer.
func NewQuestionAnswerServer(questionUseCase _interface.QuestionAnswerUseCase) *QuestionAnswerServer {
	return &QuestionAnswerServer{QuestionAnswerUseCase: questionUseCase}
}

// GetQuestions retrieves questions via RPC call.
func (es *QuestionAnswerServer) GetQuestions(ctx context.Context, input *proto.GetQuestionsRequest) (*proto.GetQuestionsReply, error) {
	questionsCore, err := es.QuestionAnswerUseCase.GetQuestions(ctx)
	if err != nil {
		return nil, fmt.Errorf("question not found")
	}

	questionsProto := make([]*proto.Question, len(questionsCore))
	for i, q := range questionsCore {
		questionsProto[i] = converters.QuestionConvertCoreInProto(q)
	}

	questionProto := new(proto.GetQuestionsReply)
	questionProto.Questions = questionsProto

	return questionProto, nil
}

// AddQuestion adds a question via RPC call.
func (es *QuestionAnswerServer) AddQuestion(ctx context.Context, input *proto.AddQuestionRequest) (*proto.AddQuestionReply, error) {
	if input == nil {
		return nil, fmt.Errorf("Question bad request")
	}

	status, err := es.QuestionAnswerUseCase.AddQuestion(converters.QuestionConvertProtoInCore(input.Question), ctx)

	statusProto := new(proto.AddQuestionReply)
	statusProto.Status = status
	if err != nil || !status {
		return statusProto, fmt.Errorf("Question no add")
	}

	return statusProto, nil
}

// AddAnswer adds an answer via RPC call.
func (es *QuestionAnswerServer) AddAnswer(ctx context.Context, input *proto.AddAnswerRequest) (*proto.AddAnswerReply, error) {
	if input == nil {
		return nil, fmt.Errorf("Answer bad request")
	}

	status, err := es.QuestionAnswerUseCase.AddAnswer(converters.AnswerConvertProtoInCore(input.Answer), ctx)

	statusProto := new(proto.AddAnswerReply)
	statusProto.Status = status
	if err != nil || !status {
		return statusProto, fmt.Errorf("Answer no add")
	}

	return statusProto, nil
}

// GetStatistic retrieves statistics via RPC call.
func (es *QuestionAnswerServer) GetStatistic(ctx context.Context, input *proto.GetStatisticRequest) (*proto.GetStatisticReply, error) {
	statisticsCore, err := es.QuestionAnswerUseCase.GetStatistics(ctx)
	if err != nil {
		return nil, fmt.Errorf("statistics not found")
	}

	statisticsProto := make([]*proto.Statistic, len(statisticsCore))
	for i, s := range statisticsCore {
		if s == nil {
			continue
		}
		statisticsProto[i] = converters.StatisticConvertCoreInProto(s)
	}

	statisticProto := new(proto.GetStatisticReply)
	statisticProto.Statistics = statisticsProto
	return statisticProto, nil
}
