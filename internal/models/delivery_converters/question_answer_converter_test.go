package delivery_converters

import (
	domain "mail/internal/microservice/models/domain_models"
	api "mail/internal/models/delivery_models"
	"reflect"
	"testing"
)

func TestQuestionConvertCoreInApi(t *testing.T) {
	questionModelCore := domain.Question{
		ID:          1,
		Text:        "text",
		MinResult:   "text",
		MaxResult:   "text",
		DopQuestion: "text",
	}

	questionModelApi := QuestionConvertCoreInApi(questionModelCore)

	expectedQuestionModelApi := &api.Question{
		ID:          questionModelCore.ID,
		Text:        questionModelCore.Text,
		MinText:     questionModelCore.MinResult,
		MaxText:     questionModelCore.MaxResult,
		DopQuestion: questionModelCore.DopQuestion,
	}

	if !reflect.DeepEqual(questionModelApi, expectedQuestionModelApi) {
		t.Errorf("QuestionConvertCoreInApi() = %v, want %v", questionModelApi, expectedQuestionModelApi)
	}
}

func TestQuestionConvertApiInCore(t *testing.T) {
	questionModelApi := api.Question{
		ID:          1,
		Text:        "text",
		MinText:     "text",
		MaxText:     "text",
		DopQuestion: "text",
	}

	questionModelCore := QuestionConvertApiInCore(questionModelApi)

	expectedQuestionModelCore := &domain.Question{
		ID:          questionModelApi.ID,
		Text:        questionModelApi.Text,
		MinResult:   questionModelApi.MinText,
		MaxResult:   questionModelApi.MaxText,
		DopQuestion: questionModelApi.DopQuestion,
	}

	if !reflect.DeepEqual(questionModelCore, expectedQuestionModelCore) {
		t.Errorf("QuestionConvertApiInCore() = %v, want %v", questionModelCore, expectedQuestionModelCore)
	}
}

func TestAnswerConvertCoreInApi(t *testing.T) {
	answerModelCore := domain.Answer{
		ID:         1,
		QuestionID: 1,
		Text:       "text",
		Login:      "login",
		Mark:       5,
	}

	answerModelApi := AnswerConvertCoreInApi(answerModelCore)

	expectedAnswerModelApi := &api.Answer{
		ID:         answerModelCore.ID,
		QuestionId: answerModelCore.QuestionID,
		Text:       answerModelCore.Text,
		Login:      answerModelCore.Login,
		Mark:       answerModelCore.Mark,
	}

	if !reflect.DeepEqual(answerModelApi, expectedAnswerModelApi) {
		t.Errorf("AnswerConvertCoreInApi() = %v, want %v", answerModelApi, expectedAnswerModelApi)
	}
}

func TestAnswerConvertApiInCore(t *testing.T) {
	answerModelApi := api.Answer{
		ID:         1,
		QuestionId: 1,
		Text:       "text",
		Login:      "login",
		Mark:       5,
	}

	answerModelCore := AnswerConvertApiInCore(answerModelApi)

	expectedAnswerModelCore := &domain.Answer{
		ID:         answerModelApi.ID,
		QuestionID: answerModelApi.QuestionId,
		Text:       answerModelApi.Text,
		Login:      answerModelApi.Login,
		Mark:       answerModelApi.Mark,
	}

	if !reflect.DeepEqual(answerModelCore, expectedAnswerModelCore) {
		t.Errorf("AnswerConvertApiInCore() = %v, want %v", answerModelCore, expectedAnswerModelCore)
	}
}
