package proto_converters

import (
	domain "mail/internal/microservice/models/domain_models"
	grpc "mail/internal/microservice/questionnaire/proto"
)

func QuestionConvertCoreInProto(questionModelCore domain.Question) *grpc.Question {
	return &grpc.Question{
		Id:      questionModelCore.ID,
		Text:    questionModelCore.Text,
		MinText: questionModelCore.MinResult,
		MaxText: questionModelCore.MaxResult,
	}
}

func QuestionConvertProtoInCore(questionModelProto grpc.Question) *domain.Question {
	return &domain.Question{
		ID:        questionModelProto.Id,
		Text:      questionModelProto.Text,
		MinResult: questionModelProto.MinText,
		MaxResult: questionModelProto.MaxText,
	}
}

func AnswerConvertCoreInProto(answerModelCore domain.Answer) *grpc.Answer {
	return &grpc.Answer{
		Id:         answerModelCore.ID,
		QuestionId: answerModelCore.QuestionId,
		Login:      answerModelCore.LoginUser,
		Mark:       answerModelCore.Mark,
	}
}

func AnswerConvertProtoInCore(answerModelProto grpc.Answer) *domain.Answer {
	return &domain.Answer{
		ID:         answerModelProto.Id,
		QuestionId: answerModelProto.QuestionId,
		LoginUser:  answerModelProto.Login,
		Mark:       answerModelProto.Mark,
	}
}
