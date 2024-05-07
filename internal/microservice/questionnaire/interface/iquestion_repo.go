//go:generate mockgen -source=./iquestion_repo.go -destination=../mock/question_repository_mock.go -package=mock

package _interface

import (
	"context"

	domain "mail/internal/microservice/models/domain_models"
)

// QuestionAnswerRepository defines the methods for interacting with question and answer data.
type QuestionAnswerRepository interface {
	// GetAllQuestions retrieves all questions from the repository.
	GetAllQuestions(ctx context.Context) ([]*domain.Question, error)

	// GetAllAnswers retrieves all answers from the repository.
	GetAllAnswers(ctx context.Context) ([]*domain.Answer, error)

	// AddQuestion adds a new question to the repository.
	AddQuestion(newQuestion *domain.Question, ctx context.Context) (bool, error)

	// AddAnswer adds a new answer to the repository.
	AddAnswer(newAnswer *domain.Answer, ctx context.Context) (bool, error)
}
