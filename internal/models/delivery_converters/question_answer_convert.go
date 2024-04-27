package delivery_converters

import (
	core "mail/internal/microservice/models/domain_models"
	api "mail/internal/models/delivery_models"
)

func QuestionConvertCoreInApi(questionModelCore core.Question) *api.Question {
	return &api.Question{
		ID:      questionModelCore.ID,
		Text:    questionModelCore.Text,
		MinText: questionModelCore.MinResult,
		MaxText: questionModelCore.MaxResult,
	}
}

func QuestionConvertApiInCore(questionModelApi api.Question) *core.Question {
	return &core.Question{
		ID:        questionModelApi.ID,
		Text:      questionModelApi.Text,
		MinResult: questionModelApi.MinText,
		MaxResult: questionModelApi.MaxText,
	}
}

func AnswerConvertCoreInApi(answerModelCore core.Answer) *api.Answer {
	return &api.Answer{
		ID:         answerModelCore.ID,
		QuestionId: answerModelCore.QuestionId,
		Login:      answerModelCore.LoginUser,
		Mark:       answerModelCore.Mark,
	}
}

func AnswerConvertApiInCore(answerModelApi api.Answer) *core.Answer {
	return &core.Answer{
		ID:         answerModelApi.ID,
		QuestionId: answerModelApi.QuestionId,
		LoginUser:  answerModelApi.Login,
		Mark:       answerModelApi.Mark,
	}
}
