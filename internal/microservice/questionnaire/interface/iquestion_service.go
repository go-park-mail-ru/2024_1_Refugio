package _interface

import (
	"context"
	emailCore "mail/internal/microservice/models/domain_models"
)

// EmailRepository represents the interface for working with emails.
type QuestionAnswerUseCase interface {
	GetQuestions(ctx context.Context) ([]*emailCore.Question, error)

	GetStatistics(ctx context.Context) ([]*emailCore.Statistics, error)

	AddQuestion(newQuestion *emailCore.Question, ctx context.Context) (bool, error)

	AddAnswer(newAnswer *emailCore.Answer, ctx context.Context) (bool, error)
}
