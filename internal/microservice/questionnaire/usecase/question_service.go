package usecase

import (
	"context"
	"errors"
	domain "mail/internal/microservice/models/domain_models"
	repository "mail/internal/microservice/questionnaire/interface"
)

// QuestionAnswerUseCase represents the use case for working questionnaire.
type QuestionAnswerUseCase struct {
	repo repository.QuestionAnswerRepository
}

// NewQuestionAnswerUseCase creates a new instance of QuestionAnswerUseCase.
func NewQuestionAnswerUseCase(repo repository.QuestionAnswerRepository) *QuestionAnswerUseCase {
	return &QuestionAnswerUseCase{
		repo: repo,
	}
}

// GetQuestions returns all question.
func (uc *QuestionAnswerUseCase) GetQuestions(ctx context.Context) ([]*domain.Question, error) {
	return uc.repo.GetAllQuestions(ctx)
}

// GetAnswers returns all answer.
func (uc *QuestionAnswerUseCase) GetStatistics(ctx context.Context) ([]*domain.Statistics, error) {
	answers, err := uc.repo.GetAllAnswers(ctx)
	if err != nil {
		return nil, err
	}

	statistics, err := CalculatingStatistics(answers)
	if err != nil {
		return nil, err
	}

	return statistics, nil
}

// AddQuestion add question.
func (uc *QuestionAnswerUseCase) AddQuestion(newQuestion *domain.Question, ctx context.Context) (bool, error) {
	return uc.repo.AddQuestion(newQuestion, ctx)
}

// AddAnswer add answer.
func (uc *QuestionAnswerUseCase) AddAnswer(newQuestion *domain.Answer, ctx context.Context) (bool, error) {
	return uc.repo.AddAnswer(newQuestion, ctx)
}

func CalculatingStatistics(answers []*domain.Answer) ([]*domain.Statistics, error) {
	var maxQuestionId uint32 = 0
	for _, a := range answers {
		if a.QuestionID > maxQuestionId {
			maxQuestionId = a.QuestionID
		}
	}

	if maxQuestionId == 0 {
		return nil, errors.New("max question id is 0")
	}

	var questionsID = make([]uint32, maxQuestionId+1)
	var sumMark = make([]uint32, maxQuestionId+1)

	var statistics = make([]*domain.Statistics, maxQuestionId)
	for i, a := range answers {
		statistics[i] = new(domain.Statistics)
		statistics[i].Average = sumMark[i+1] / questionsID[i+1]
		statistics[i].Text = a.Text
	}

	return statistics, nil
}
