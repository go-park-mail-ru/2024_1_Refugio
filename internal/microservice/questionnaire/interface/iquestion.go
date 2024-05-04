//go:generate mockgen -source=./iquestion.go -destination=../mock/question_mock.go -package=mock

package _interface

import (
	"context"

	"mail/internal/microservice/questionnaire/proto"
)

// QuestionServer represents the interface for working with users.
type QuestionServer interface {
	// GetQuestions retrieves questions via RPC call.
	GetQuestions(ctx context.Context, input *proto.GetQuestionsRequest) (*proto.GetQuestionsReply, error)

	// AddQuestion adds a question via RPC call.
	AddQuestion(ctx context.Context, input *proto.AddQuestionRequest) (*proto.AddQuestionReply, error)

	// AddAnswer adds an answer via RPC call.
	AddAnswer(ctx context.Context, input *proto.AddAnswerRequest) (*proto.AddAnswerReply, error)

	// GetStatistic retrieves statistics via RPC call.
	GetStatistic(ctx context.Context, input *proto.GetStatisticRequest) (*proto.GetStatisticReply, error)
}
