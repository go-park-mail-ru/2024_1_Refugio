package repository_converters

import (
	domain "mail/internal/microservice/models/domain_models"
	database "mail/internal/microservice/models/repository_models"
)

func QuestionConvertDbInCore(questionModelDb database.Question) *domain.Question {
	return &domain.Question{
		ID:          questionModelDb.ID,
		Text:        questionModelDb.Text,
		MinResult:   questionModelDb.MinResult,
		MaxResult:   questionModelDb.MaxResult,
		DopQuestion: questionModelDb.DopQuestion,
	}
}

func QuestionConvertCoreInDb(questionModelCore domain.Question) *database.Question {
	return &database.Question{
		ID:          questionModelCore.ID,
		Text:        questionModelCore.Text,
		MinResult:   questionModelCore.MinResult,
		MaxResult:   questionModelCore.MaxResult,
		DopQuestion: questionModelCore.DopQuestion,
	}
}
