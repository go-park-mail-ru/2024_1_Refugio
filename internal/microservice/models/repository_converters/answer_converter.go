package repository_converters

import (
	domain "mail/internal/microservice/models/domain_models"
	database "mail/internal/microservice/models/repository_models"
)

// AnswerConvertDbInCore converts an answer model from the database format to the application core format.
func AnswerConvertDbInCore(answerModelDb *database.Answer) *domain.Answer {
	return &domain.Answer{
		ID:         answerModelDb.ID,
		QuestionID: answerModelDb.QuestionID,
		Login:      answerModelDb.Login,
		Mark:       answerModelDb.Mark,
		Text:       answerModelDb.Text,
	}
}

// AnswerConvertCoreInDb converts an answer model from the application core format to the database format.
func AnswerConvertCoreInDb(answerModelDb *domain.Answer) *database.Answer {
	return &database.Answer{
		ID:         answerModelDb.ID,
		QuestionID: answerModelDb.QuestionID,
		Login:      answerModelDb.Login,
		Mark:       answerModelDb.Mark,
		Text:       answerModelDb.Text,
	}
}
