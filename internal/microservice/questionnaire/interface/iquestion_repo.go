package _interface

import (
	"context"
	domain "mail/internal/microservice/models/domain_models"
)

// EmailRepository represents the interface for working with emails.
type QuestionAnswerRepository interface {
	GetAllQuestions(ctx context.Context) ([]*domain.Question, error)

	GetAllAnswers(ctx context.Context) ([]*domain.Answer, error)

	AddQuestion(newQuestion *domain.Question, ctx context.Context) (bool, error)

	AddAnswer(newAnswer *domain.Answer, ctx context.Context) (bool, error)
}
