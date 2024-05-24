package repository_converters

import (
	"testing"

	"github.com/stretchr/testify/assert"

	domain "mail/internal/microservice/models/domain_models"
	database "mail/internal/microservice/models/repository_models"
)

func TestQuestionConvertDbInCore(t *testing.T) {
	questionModelDb := database.Question{
		ID:          123,
		Text:        "What is your favorite color?",
		MinResult:   "Red",
		MaxResult:   "Blue",
		DopQuestion: "Why do you like it?",
	}

	expectedCore := &domain.Question{
		ID:          123,
		Text:        "What is your favorite color?",
		MinResult:   "Red",
		MaxResult:   "Blue",
		DopQuestion: "Why do you like it?",
	}

	actualCore := QuestionConvertDbInCore(&questionModelDb)
	assert.Equal(t, expectedCore, actualCore)
}

func TestQuestionConvertCoreInDb(t *testing.T) {
	questionModelCore := domain.Question{
		ID:          123,
		Text:        "What is your favorite color?",
		MinResult:   "Red",
		MaxResult:   "Blue",
		DopQuestion: "Why do you like it?",
	}

	expectedDb := &database.Question{
		ID:          123,
		Text:        "What is your favorite color?",
		MinResult:   "Red",
		MaxResult:   "Blue",
		DopQuestion: "Why do you like it?",
	}

	actualDb := QuestionConvertCoreInDb(&questionModelCore)
	assert.Equal(t, expectedDb, actualDb)
}
