package server

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"mail/internal/microservice/questionnaire/mock"
	"mail/internal/microservice/questionnaire/proto"
	"mail/internal/pkg/logger"
	"mail/internal/pkg/utils/constants"

	domain "mail/internal/microservice/models/domain_models"
)

func GetCTX() context.Context {
	ctx := context.WithValue(context.Background(), constants.LoggerKey, logger.InitializationBdLog(nil))
	ctx2 := context.WithValue(ctx, constants.RequestIDKey, []string{"testID"})

	return ctx2
}

func TestQuestionAnswerServer_GetQuestions_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuestionUseCase := mock.NewMockQuestionAnswerUseCase(ctrl)

	server := NewQuestionAnswerServer(mockQuestionUseCase)

	ctx := GetCTX()

	expectedQuestions := []*domain.Question{
		{ID: 1, Text: "Question 1", MinResult: "Min Text 1", MaxResult: "Max Text 1", DopQuestion: "Dop Question 1"},
		{ID: 2, Text: "Question 2", MinResult: "Min Text 2", MaxResult: "Max Text 2", DopQuestion: "Dop Question 2"},
	}

	mockQuestionUseCase.EXPECT().GetQuestions(ctx).Return(expectedQuestions, nil)

	reply, err := server.GetQuestions(ctx, &proto.GetQuestionsRequest{})

	assert.NoError(t, err)
	assert.NotNil(t, reply)
}

func TestQuestionAnswerServer_GetQuestions_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuestionUseCase := mock.NewMockQuestionAnswerUseCase(ctrl)

	server := NewQuestionAnswerServer(mockQuestionUseCase)

	ctx := GetCTX()

	mockQuestionUseCase.EXPECT().GetQuestions(ctx).Return(nil, errors.New("question not found"))

	reply, err := server.GetQuestions(ctx, &proto.GetQuestionsRequest{})

	assert.Error(t, err)
	assert.Nil(t, reply)
}

func TestQuestionAnswerServer_AddQuestion_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuestionUseCase := mock.NewMockQuestionAnswerUseCase(ctrl)
	server := NewQuestionAnswerServer(mockQuestionUseCase)
	ctx := context.Background()

	inputQuestion := &proto.AddQuestionRequest{
		Question: &proto.Question{
			Id:          1,
			Text:        "Test Question",
			MinText:     "Test Min Result",
			MaxText:     "Test Max Result",
			DopQuestion: "Test Dop Question",
		},
	}

	mockQuestionUseCase.EXPECT().AddQuestion(gomock.Any(), ctx).Return(true, nil)

	reply, err := server.AddQuestion(ctx, inputQuestion)

	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.True(t, reply.Status)
}

func TestQuestionAnswerServer_AddQuestion_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuestionUseCase := mock.NewMockQuestionAnswerUseCase(ctrl)
	server := NewQuestionAnswerServer(mockQuestionUseCase)
	ctx := context.Background()

	inputQuestion := &proto.AddQuestionRequest{
		Question: &proto.Question{
			Id:          1,
			Text:        "Test Question",
			MinText:     "Test Min Result",
			MaxText:     "Test Max Result",
			DopQuestion: "Test Dop Question",
		},
	}

	mockQuestionUseCase.EXPECT().AddQuestion(gomock.Any(), ctx).Return(false, fmt.Errorf("failed to add question"))

	reply, err := server.AddQuestion(ctx, inputQuestion)

	assert.Error(t, err)
	assert.NotNil(t, reply)
	assert.False(t, reply.Status)
}

func TestQuestionAnswerServer_AddAnswer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuestionUseCase := mock.NewMockQuestionAnswerUseCase(ctrl)
	server := NewQuestionAnswerServer(mockQuestionUseCase)
	ctx := context.Background()

	inputAnswer := &proto.AddAnswerRequest{
		Answer: &proto.Answer{
			Id:         1,
			QuestionId: 1,
			Login:      "test_user",
			Mark:       5,
			Text:       "Test Answer",
		},
	}

	mockQuestionUseCase.EXPECT().AddAnswer(gomock.Any(), ctx).Return(true, nil)

	reply, err := server.AddAnswer(ctx, inputAnswer)

	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.True(t, reply.Status)
}

func TestQuestionAnswerServer_AddAnswer_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuestionUseCase := mock.NewMockQuestionAnswerUseCase(ctrl)
	server := NewQuestionAnswerServer(mockQuestionUseCase)
	ctx := context.Background()

	inputAnswer := &proto.AddAnswerRequest{
		Answer: &proto.Answer{
			Id:         1,
			QuestionId: 1,
			Login:      "test_user",
			Mark:       5,
			Text:       "Test Answer",
		},
	}

	mockQuestionUseCase.EXPECT().AddAnswer(gomock.Any(), ctx).Return(false, fmt.Errorf("failed to add answer"))

	reply, err := server.AddAnswer(ctx, inputAnswer)

	assert.Error(t, err)
	assert.NotNil(t, reply)
	assert.False(t, reply.Status)
}

func TestQuestionAnswerServer_GetStatistic_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuestionUseCase := mock.NewMockQuestionAnswerUseCase(ctrl)
	server := NewQuestionAnswerServer(mockQuestionUseCase)
	ctx := context.Background()

	mockStatistics := []*domain.Statistics{
		{Text: "Statistic 1", Average: 4.5},
		{Text: "Statistic 2", Average: 3.8},
	}

	mockQuestionUseCase.EXPECT().GetStatistics(gomock.Any()).Return(mockStatistics, nil)

	reply, err := server.GetStatistic(ctx, &proto.GetStatisticRequest{})

	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Statistics, 2)
	assert.Equal(t, "Statistic 1", reply.Statistics[0].Text)
	assert.Equal(t, float32(4.5), reply.Statistics[0].Average)
	assert.Equal(t, "Statistic 2", reply.Statistics[1].Text)
	assert.Equal(t, float32(3.8), reply.Statistics[1].Average)
}

func TestQuestionAnswerServer_GetStatistic_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuestionUseCase := mock.NewMockQuestionAnswerUseCase(ctrl)
	server := NewQuestionAnswerServer(mockQuestionUseCase)
	ctx := context.Background()

	mockQuestionUseCase.EXPECT().GetStatistics(gomock.Any()).Return(nil, fmt.Errorf("failed to fetch statistics"))

	reply, err := server.GetStatistic(ctx, &proto.GetStatisticRequest{})

	assert.Error(t, err)
	assert.Nil(t, reply)
}
