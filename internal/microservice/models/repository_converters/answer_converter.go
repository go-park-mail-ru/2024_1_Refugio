package repository_converters

import (
	domain "mail/internal/microservice/models/domain_models"
	database "mail/internal/microservice/models/repository_models"
)

func AnswerConvertDbInCore(answerModelDb database.Answer) *domain.Answer {
	return &domain.Answer{
		ID:         answerModelDb.ID,
		QuestionID: answerModelDb.QuestionID,
		Login:      answerModelDb.Login,
		Mark:       answerModelDb.Mark,
		Text:       answerModelDb.Text,
	}
}

func AnswerConvertCoreInDb(answerModelDb domain.Answer) *database.Answer {
	return &database.Answer{
		ID:         answerModelDb.ID,
		QuestionID: answerModelDb.QuestionID,
		Login:      answerModelDb.Login,
		Mark:       answerModelDb.Mark,
		Text:       answerModelDb.Text,
	}
}
