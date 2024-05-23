package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"mail/internal/microservice/questionnaire/mock"
	"mail/internal/pkg/logger"
	"mail/internal/pkg/utils/constants"

	domain "mail/internal/microservice/models/domain_models"
)

func GetCTX() context.Context {
	ctx := context.WithValue(context.Background(), constants.LoggerKey, logger.InitializationBdLog(nil))
	ctx2 := context.WithValue(ctx, constants.RequestIDKey, []string{"testID"})

	return ctx2
}

func TestQuestionAnswerUseCase_GetQuestions_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockQuestionAnswerRepository(ctrl)
	useCase := NewQuestionAnswerUseCase(mockRepo)

	ctx := GetCTX()

	expectedQuestions := []*domain.Question{
		{ID: 1, Text: "Question 1", MinResult: "Min Text 1", MaxResult: "Max Text 1", DopQuestion: "Dop Question 1"},
		{ID: 2, Text: "Question 2", MinResult: "Min Text 2", MaxResult: "Max Text 2", DopQuestion: "Dop Question 2"},
	}

	mockRepo.EXPECT().GetAllQuestions(ctx).Return(expectedQuestions, nil)

	questions, err := useCase.GetQuestions(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedQuestions, questions)
}

func TestQuestionAnswerUseCase_GetQuestions_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockQuestionAnswerRepository(ctrl)
	useCase := NewQuestionAnswerUseCase(mockRepo)

	ctx := GetCTX()

	mockRepo.EXPECT().GetAllQuestions(ctx).Return(nil, errors.New("some error"))

	questions, err := useCase.GetQuestions(ctx)

	assert.Error(t, err)
	assert.Nil(t, questions)
}

func TestQuestionAnswerUseCase_AddQuestion_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockQuestionAnswerRepository(ctrl)
	useCase := NewQuestionAnswerUseCase(mockRepo)

	ctx := context.Background()
	newQuestion := &domain.Question{ID: 1, Text: "What is your favorite color?"}

	mockRepo.EXPECT().AddQuestion(newQuestion, ctx).Return(true, nil)

	success, err := useCase.AddQuestion(newQuestion, ctx)

	assert.NoError(t, err)
	assert.True(t, success)
}

func TestQuestionAnswerUseCase_AddQuestion_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockQuestionAnswerRepository(ctrl)
	useCase := NewQuestionAnswerUseCase(mockRepo)

	ctx := context.Background()
	newQuestion := &domain.Question{ID: 1, Text: "What is your favorite color?"}

	mockRepo.EXPECT().AddQuestion(newQuestion, ctx).Return(false, errors.New("failed to add question"))

	success, err := useCase.AddQuestion(newQuestion, ctx)

	assert.Error(t, err)
	assert.False(t, success)
}

func TestQuestionAnswerUseCase_AddAnswer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockQuestionAnswerRepository(ctrl)
	useCase := NewQuestionAnswerUseCase(mockRepo)

	ctx := context.Background()
	newAnswer := &domain.Answer{ID: 1, QuestionID: 1, Login: "user1", Mark: 5, Text: "This is an answer."}

	mockRepo.EXPECT().AddAnswer(newAnswer, ctx).Return(true, nil)

	success, err := useCase.AddAnswer(newAnswer, ctx)

	assert.NoError(t, err)
	assert.True(t, success)
}

func TestQuestionAnswerUseCase_AddAnswer_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockQuestionAnswerRepository(ctrl)
	useCase := NewQuestionAnswerUseCase(mockRepo)

	ctx := context.Background()
	newAnswer := &domain.Answer{ID: 1, QuestionID: 1, Login: "user1", Mark: 5, Text: "This is an answer."}

	mockRepo.EXPECT().AddAnswer(newAnswer, ctx).Return(false, errors.New("failed to add answer"))

	success, err := useCase.AddAnswer(newAnswer, ctx)

	assert.Error(t, err)
	assert.False(t, success)
}

func TestQuestionAnswerUseCase_GetStatistics_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockQuestionAnswerRepository(ctrl)
	useCase := NewQuestionAnswerUseCase(mockRepo)

	ctx := context.Background()
	answers := []*domain.Answer{
		{QuestionID: 1, Mark: 5, Text: "Answer 1"},
		{QuestionID: 1, Mark: 3, Text: "Answer 2"},
		{QuestionID: 2, Mark: 4, Text: "Answer 3"},
	}

	mockRepo.EXPECT().GetAllAnswers(ctx).Return(answers, nil)

	statistics, err := useCase.GetStatistics(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, statistics)
}

func TestQuestionAnswerUseCase_GetStatistics_NoAnswers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockQuestionAnswerRepository(ctrl)
	useCase := NewQuestionAnswerUseCase(mockRepo)

	ctx := context.Background()

	mockRepo.EXPECT().GetAllAnswers(ctx).Return(nil, nil)

	statistics, err := useCase.GetStatistics(ctx)

	assert.Error(t, err)
	assert.Nil(t, statistics)
}

func TestCalculatingStatistics_Success(t *testing.T) {
	answers := []*domain.Answer{
		{QuestionID: 1, Mark: 5, Text: "Answer 1"},
		{QuestionID: 1, Mark: 3, Text: "Answer 2"},
		{QuestionID: 2, Mark: 4, Text: "Answer 3"},
	}

	statistics, err := CalculatingStatistics(answers)

	assert.NoError(t, err)
	assert.NotNil(t, statistics)
}

func TestCalculatingStatistics_EmptyAnswers(t *testing.T) {
	answers := []*domain.Answer{}

	statistics, err := CalculatingStatistics(answers)

	assert.Error(t, err)
	assert.Nil(t, statistics)
}

func TestCalculatingStatistics_MaxQuestionIDZero(t *testing.T) {
	answers := []*domain.Answer{{QuestionID: 0, Mark: 5, Text: "Answer 1"}}

	statistics, err := CalculatingStatistics(answers)

	assert.Error(t, err)
	assert.Nil(t, statistics)
}

func TestCalculatingStatistics_ZeroQuestionsID(t *testing.T) {
	answers := []*domain.Answer{{QuestionID: 1, Mark: 5, Text: "Answer 1"}}

	statistics, err := CalculatingStatistics(answers)

	assert.NoError(t, err)
	assert.NotNil(t, statistics)

	expectedStatistics := []*domain.Statistics{
		{Text: "Answer 1", Average: 5},
	}
	assert.Equal(t, expectedStatistics, statistics)
}
