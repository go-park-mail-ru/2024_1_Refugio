package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"mail/internal/models/response"

	question_mock "mail/internal/microservice/questionnaire/mock"
	question_proto "mail/internal/microservice/questionnaire/proto"
	api "mail/internal/models/delivery_models"
	session_mock "mail/internal/pkg/session/mock"
)

func TestQuestionHandler_GetAllQuestions_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuestionServiceClient := question_mock.NewMockQuestionServiceClient(ctrl)

	qh := &QuestionHandler{
		QuestionServiceClient: mockQuestionServiceClient,
	}

	req := httptest.NewRequest("GET", "/api/v1/questions", nil)
	ctx := context.WithValue(req.Context(), "requestID", "testID")
	req = req.WithContext(ctx)
	req.Header.Set("X-Csrf-Token", "csrf_token")

	w := httptest.NewRecorder()

	mockQuestionServiceClient.EXPECT().
		GetQuestions(gomock.Any(), gomock.Any()).
		Return(&question_proto.GetQuestionsReply{
			Questions: []*question_proto.Question{
				{Id: 1, Text: "What is your name?"},
				{Id: 2, Text: "Where are you from?"},
			},
		}, nil)

	qh.GetAllQuestions(w, req)

	resp := w.Result()

	var responseBody map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body := responseBody["body"]
	bodyMap, _ := body.(map[string]interface{})

	questions, ok := bodyMap["questions"].([]interface{})

	assert.True(t, ok)
	assert.Equal(t, 2, len(questions))
}

func TestQuestionHandler_GetAllQuestions_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuestionServiceClient := question_mock.NewMockQuestionServiceClient(ctrl)

	qh := &QuestionHandler{
		QuestionServiceClient: mockQuestionServiceClient,
	}

	req := httptest.NewRequest("GET", "/api/v1/questions", nil)
	ctx := context.WithValue(req.Context(), "requestID", "testID")
	req = req.WithContext(ctx)
	req.Header.Set("X-Csrf-Token", "csrf_token")

	w := httptest.NewRecorder()

	mockQuestionServiceClient.EXPECT().
		GetQuestions(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("failed to get questions"))

	qh.GetAllQuestions(w, req)

	resp := w.Result()

	var responseBody response.ErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestQuestionHandler_AddQuestion_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockQuestionServiceClient := question_mock.NewMockQuestionServiceClient(ctrl)

	qh := &QuestionHandler{
		Sessions:              mockSessionsManager,
		QuestionServiceClient: mockQuestionServiceClient,
	}

	newQuestion := api.Question{
		Text:        "What is your favorite color?",
		MinText:     "Red",
		MaxText:     "Blue",
		DopQuestion: "Green",
	}

	requestBody, err := json.Marshal(newQuestion)
	assert.NoError(t, err)

	req := httptest.NewRequest("POST", "/api/v1/questions", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := req.Context()
	ctx = context.WithValue(ctx, "requestID", "testID")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	mockQuestionServiceClient.EXPECT().AddQuestion(gomock.Any(), gomock.Any()).Return(&question_proto.AddQuestionReply{Status: true}, nil)

	qh.AddQuestion(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestQuestionHandler_AddQuestion_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockQuestionServiceClient := question_mock.NewMockQuestionServiceClient(ctrl)

	qh := &QuestionHandler{
		Sessions:              mockSessionsManager,
		QuestionServiceClient: mockQuestionServiceClient,
	}

	requestBody := []byte(`{"Text": "What is your favorite color?, "MinText": "Red", "MaxText": "Blue", "DopQuestion": "Green"}`)

	req := httptest.NewRequest("POST", "/api/v1/questions", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := req.Context()
	ctx = context.WithValue(ctx, "requestID", "testID")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	qh.AddQuestion(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestQuestionHandler_AddQuestion_ServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockQuestionServiceClient := question_mock.NewMockQuestionServiceClient(ctrl)

	qh := &QuestionHandler{
		Sessions:              mockSessionsManager,
		QuestionServiceClient: mockQuestionServiceClient,
	}

	newQuestion := api.Question{
		Text:        "What is your favorite color?",
		MinText:     "Red",
		MaxText:     "Blue",
		DopQuestion: "Green",
	}

	requestBody, err := json.Marshal(newQuestion)
	assert.NoError(t, err)

	req := httptest.NewRequest("POST", "/api/v1/questions", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := req.Context()
	ctx = context.WithValue(ctx, "requestID", "testID")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	mockQuestionServiceClient.EXPECT().AddQuestion(gomock.Any(), gomock.Any()).Return(&question_proto.AddQuestionReply{Status: false}, errors.New("internal server error"))

	qh.AddQuestion(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestQuestionHandler_AddAnswer_ValidInput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockQuestionServiceClient := question_mock.NewMockQuestionServiceClient(ctrl)

	qh := &QuestionHandler{
		Sessions:              mockSessionsManager,
		QuestionServiceClient: mockQuestionServiceClient,
	}

	newAnswer := api.Answer{
		Text: "This is a sample answer.",
	}

	requestBody, err := json.Marshal(newAnswer)
	assert.NoError(t, err)

	req := httptest.NewRequest("POST", "/api/v1/answers", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := req.Context()
	ctx = context.WithValue(ctx, "requestID", "testID")
	req = req.WithContext(ctx)

	mockSessionsManager.EXPECT().GetLoginBySession(gomock.Any(), gomock.Any()).Return("testuser", nil)

	w := httptest.NewRecorder()

	mockQuestionServiceClient.EXPECT().AddAnswer(gomock.Any(), gomock.Any()).Return(&question_proto.AddAnswerReply{Status: true}, nil)

	qh.AddAnswer(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestQuestionHandler_AddAnswer_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockQuestionServiceClient := question_mock.NewMockQuestionServiceClient(ctrl)

	qh := &QuestionHandler{
		Sessions:              mockSessionsManager,
		QuestionServiceClient: mockQuestionServiceClient,
	}

	requestBody := []byte(`{"Text": "This is a sample answer.}`)

	req := httptest.NewRequest("POST", "/api/v1/answers", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := req.Context()
	ctx = context.WithValue(ctx, "requestID", "testID")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	qh.AddAnswer(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestQuestionHandler_AddAnswer_SessionServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockQuestionServiceClient := question_mock.NewMockQuestionServiceClient(ctrl)

	qh := &QuestionHandler{
		Sessions:              mockSessionsManager,
		QuestionServiceClient: mockQuestionServiceClient,
	}

	newAnswer := api.Answer{
		Text: "This is a sample answer.",
	}

	requestBody, err := json.Marshal(newAnswer)
	assert.NoError(t, err)

	req := httptest.NewRequest("POST", "/api/v1/answers", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := req.Context()
	ctx = context.WithValue(ctx, "requestID", "testID")
	req = req.WithContext(ctx)

	mockSessionsManager.EXPECT().GetLoginBySession(gomock.Any(), gomock.Any()).Return("", errors.New("internal server error"))

	w := httptest.NewRecorder()

	qh.AddAnswer(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestQuestionHandler_AddAnswer_AnswerServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockQuestionServiceClient := question_mock.NewMockQuestionServiceClient(ctrl)

	qh := &QuestionHandler{
		Sessions:              mockSessionsManager,
		QuestionServiceClient: mockQuestionServiceClient,
	}

	newAnswer := api.Answer{
		Text: "This is a sample answer.",
	}

	requestBody, err := json.Marshal(newAnswer)
	assert.NoError(t, err)

	req := httptest.NewRequest("POST", "/api/v1/answers", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := req.Context()
	ctx = context.WithValue(ctx, "requestID", "testID")
	req = req.WithContext(ctx)

	mockSessionsManager.EXPECT().GetLoginBySession(gomock.Any(), gomock.Any()).Return("testuser", nil)

	mockQuestionServiceClient.EXPECT().AddAnswer(gomock.Any(), gomock.Any()).Return(&question_proto.AddAnswerReply{Status: false}, errors.New("internal server error"))

	w := httptest.NewRecorder()

	qh.AddAnswer(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestQuestionHandler_GetStatistics_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockQuestionServiceClient := question_mock.NewMockQuestionServiceClient(ctrl)

	qh := &QuestionHandler{
		Sessions:              mockSessionsManager,
		QuestionServiceClient: mockQuestionServiceClient,
	}

	req := httptest.NewRequest("GET", "/api/v1/statistics", nil)
	req.Header.Set("Content-Type", "application/json")

	ctx := req.Context()
	ctx = context.WithValue(ctx, "requestID", "testID")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	mockStats := []*question_proto.Statistic{
		{Text: "100", Average: 200},
		{Text: "150", Average: 250},
	}

	mockQuestionServiceClient.EXPECT().GetStatistic(gomock.Any(), gomock.Any()).Return(&question_proto.GetStatisticReply{Statistics: mockStats}, nil)

	qh.GetStatistics(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestQuestionHandler_GetStatistics_ServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockQuestionServiceClient := question_mock.NewMockQuestionServiceClient(ctrl)

	qh := &QuestionHandler{
		Sessions:              mockSessionsManager,
		QuestionServiceClient: mockQuestionServiceClient,
	}

	req := httptest.NewRequest("GET", "/api/v1/statistics", nil)
	req.Header.Set("Content-Type", "application/json")

	ctx := req.Context()
	ctx = context.WithValue(ctx, "requestID", "testID")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	mockQuestionServiceClient.EXPECT().GetStatistic(gomock.Any(), gomock.Any()).Return(nil, errors.New("get statistics error"))

	qh.GetStatistics(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
