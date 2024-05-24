package proto_converters

import (
	"testing"

	"github.com/stretchr/testify/assert"

	domain "mail/internal/microservice/models/domain_models"
	grpc "mail/internal/microservice/questionnaire/proto"
)

func TestQuestionConvertCoreInProto(t *testing.T) {
	questionModelCore := domain.Question{
		ID:          1,
		Text:        "What is your age?",
		MinResult:   "18",
		MaxResult:   "99",
		DopQuestion: "What is your gender?",
	}

	expectedProto := &grpc.Question{
		Id:          1,
		Text:        "What is your age?",
		MinText:     "18",
		MaxText:     "99",
		DopQuestion: "What is your gender?",
	}

	actualProto := QuestionConvertCoreInProto(&questionModelCore)
	assert.Equal(t, expectedProto, actualProto)
}

func TestQuestionConvertProtoInCore(t *testing.T) {
	questionModelProto := grpc.Question{
		Id:          1,
		Text:        "What is your age?",
		MinText:     "18",
		MaxText:     "99",
		DopQuestion: "What is your gender?",
	}

	expectedCore := &domain.Question{
		ID:          1,
		Text:        "What is your age?",
		MinResult:   "18",
		MaxResult:   "99",
		DopQuestion: "What is your gender?",
	}

	actualCore := QuestionConvertProtoInCore(&questionModelProto)
	assert.Equal(t, expectedCore, actualCore)
}

func TestAnswerConvertCoreInProto(t *testing.T) {
	answerModelCore := domain.Answer{
		ID:         1,
		QuestionID: 1,
		Login:      "john_doe",
		Mark:       5,
		Text:       "Good job!",
	}

	expectedProto := &grpc.Answer{
		Id:         1,
		QuestionId: 1,
		Login:      "john_doe",
		Mark:       5,
		Text:       "Good job!",
	}

	actualProto := AnswerConvertCoreInProto(&answerModelCore)
	assert.Equal(t, expectedProto, actualProto)
}

func TestAnswerConvertProtoInCore(t *testing.T) {
	answerModelProto := grpc.Answer{
		Id:         1,
		QuestionId: 1,
		Login:      "john_doe",
		Mark:       5,
		Text:       "Good job!",
	}

	expectedCore := &domain.Answer{
		ID:         1,
		QuestionID: 1,
		Login:      "john_doe",
		Mark:       5,
		Text:       "Good job!",
	}

	actualCore := AnswerConvertProtoInCore(&answerModelProto)
	assert.Equal(t, expectedCore, actualCore)
}

func TestStatisticConvertCoreInProto(t *testing.T) {
	statisticModelCore := domain.Statistics{
		Text:    "Average Age",
		Average: 30.5,
	}

	expectedProto := &grpc.Statistic{
		Text:    "Average Age",
		Average: 30.5,
	}

	actualProto := StatisticConvertCoreInProto(&statisticModelCore)
	assert.Equal(t, expectedProto, actualProto)
}
