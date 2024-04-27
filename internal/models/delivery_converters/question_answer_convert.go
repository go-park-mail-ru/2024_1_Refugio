package delivery_converters

import (
	core "mail/internal/microservice/models/domain_models"
	api "mail/internal/models/delivery_models"
)

func QuestionConvertCoreInApi(questionModelDb core.Session) *api.Question {
	return &api.Question{}
}

func QuestionConvertApiInCore(questionModelApi api.Question) *core.Session {
	return &core.Session{}
}

func AnswerConvertCoreInApi(answerModelDb core.Session) *api.Answer {
	return &api.Answer{}
}

func AnswerConvertApiInCore(answerModelApi api.Answer) *core.Session {
	return &core.Session{}
}
