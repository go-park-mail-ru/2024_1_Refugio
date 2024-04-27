package _interface

import (
	"context"
	domain "mail/internal/microservice/models/domain_models"
)

// EmailRepository represents the interface for working with emails.
type QuestionRepository interface {
	GetAllQuestions(ctx context.Context) ([]*domain.Question, error)

	/*
		AddQuestion(newQuestion *domain.Question, ctx context.Context) (bool, error)

		GetAllAnswers(ctx context.Context) ([]*domain.Answer, error)

		AddAnswer(newAnswer *domain.Answer, ctx context.Context) (bool, error)
	*/
}
