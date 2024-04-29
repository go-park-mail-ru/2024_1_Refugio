package proto_converters

import (
	domain "mail/internal/microservice/models/domain_models"
	grpc "mail/internal/microservice/questionnaire/proto"
)

func QuestionConvertCoreInProto(questionModelCore domain.Question) *grpc.Question {
	return &grpc.Question{
		Id:          questionModelCore.ID,
		Text:        questionModelCore.Text,
		MinText:     questionModelCore.MinResult,
		MaxText:     questionModelCore.MaxResult,
		DopQuestion: questionModelCore.DopQuestion,
	}
}

func QuestionConvertProtoInCore(questionModelProto grpc.Question) *domain.Question {
	return &domain.Question{
		ID:          questionModelProto.Id,
		Text:        questionModelProto.Text,
		MinResult:   questionModelProto.MinText,
		MaxResult:   questionModelProto.MaxText,
		DopQuestion: questionModelProto.DopQuestion,
	}
}

func AnswerConvertCoreInProto(answerModelCore domain.Answer) *grpc.Answer {
	return &grpc.Answer{
		Id:         answerModelCore.ID,
		QuestionId: answerModelCore.QuestionID,
		Login:      answerModelCore.Login,
		Mark:       answerModelCore.Mark,
		Text:       answerModelCore.Text,
	}
}

func AnswerConvertProtoInCore(answerModelProto grpc.Answer) *domain.Answer {
	return &domain.Answer{
		ID:         answerModelProto.Id,
		QuestionID: answerModelProto.QuestionId,
		Login:      answerModelProto.Login,
		Mark:       answerModelProto.Mark,
		Text:       answerModelProto.Text,
	}
}

func StatisticConvertCoreInProto(statisticModelCore domain.Statistics) *grpc.Statistic {
	return &grpc.Statistic{
		Text:    statisticModelCore.Text,
		Average: statisticModelCore.Average,
	}
}
