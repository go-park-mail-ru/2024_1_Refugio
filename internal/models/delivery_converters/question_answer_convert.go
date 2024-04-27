package delivery_converters

import (
	core "mail/internal/microservice/models/domain_models"
	api "mail/internal/models/delivery_models"
)

func QuestionConvertCoreInApi(questionModelCore core.Question) *api.Question {
	return &api.Question{
		ID:          questionModelCore.ID,
		Text:        questionModelCore.Text,
		MinText:     questionModelCore.MinResult,
		MaxText:     questionModelCore.MaxResult,
		DopQuestion: questionModelCore.DopQuestion,
	}
}

func QuestionConvertApiInCore(questionModelApi api.Question) *core.Question {
	return &core.Question{
		ID:          questionModelApi.ID,
		Text:        questionModelApi.Text,
		MinResult:   questionModelApi.MinText,
		MaxResult:   questionModelApi.MaxText,
		DopQuestion: questionModelApi.DopQuestion,
	}
}

func AnswerConvertCoreInApi(answerModelCore core.Answer) *api.Answer {
	return &api.Answer{
		ID:         answerModelCore.ID,
		QuestionId: answerModelCore.QuestionID,
		Login:      answerModelCore.Login,
		Mark:       answerModelCore.Mark,
		Text:       answerModelCore.Text,
	}
}

func AnswerConvertApiInCore(answerModelApi api.Answer) *core.Answer {
	return &core.Answer{
		ID:         answerModelApi.ID,
		QuestionID: answerModelApi.QuestionId,
		Login:      answerModelApi.Login,
		Mark:       answerModelApi.Mark,
		Text:       answerModelApi.Text,
	}
}
