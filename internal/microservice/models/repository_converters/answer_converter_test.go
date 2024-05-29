package repository_converters

import (
	"testing"

	"github.com/stretchr/testify/assert"

	domain "mail/internal/microservice/models/domain_models"
	database "mail/internal/microservice/models/repository_models"
)

func TestAnswerConvertDbInCore(t *testing.T) {
	answerModelDb := database.Answer{
		ID:         123,
		QuestionID: 456,
		Login:      "john_doe",
		Mark:       5,
		Text:       "Some text",
	}

	expectedCore := &domain.Answer{
		ID:         123,
		QuestionID: 456,
		Login:      "john_doe",
		Mark:       5,
		Text:       "Some text",
	}

	actualCore := AnswerConvertDbInCore(&answerModelDb)
	assert.Equal(t, expectedCore, actualCore)
}

func TestAnswerConvertCoreInDb(t *testing.T) {
	answerModelCore := domain.Answer{
		ID:         123,
		QuestionID: 456,
		Login:      "john_doe",
		Mark:       5,
		Text:       "Some text",
	}

	expectedDb := &database.Answer{
		ID:         123,
		QuestionID: 456,
		Login:      "john_doe",
		Mark:       5,
		Text:       "Some text",
	}

	actualDb := AnswerConvertCoreInDb(&answerModelCore)
	assert.Equal(t, expectedDb, actualDb)
}
