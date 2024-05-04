//go:generate mockgen -source=./iquestion_service.go -destination=../mock/question_service_mock.go -package=mock

package _interface

import (
	"context"

	emailCore "mail/internal/microservice/models/domain_models"
)

// QuestionAnswerUseCase defines the methods for managing questions and answers.
type QuestionAnswerUseCase interface {
	// GetQuestions retrieves all questions.
	GetQuestions(ctx context.Context) ([]*emailCore.Question, error)

	// GetStatistics retrieves statistics related to questions and answers.
	GetStatistics(ctx context.Context) ([]*emailCore.Statistics, error)

	// AddQuestion adds a new question.
	AddQuestion(newQuestion *emailCore.Question, ctx context.Context) (bool, error)

	// AddAnswer adds a new answer.
	AddAnswer(newAnswer *emailCore.Answer, ctx context.Context) (bool, error)
}
