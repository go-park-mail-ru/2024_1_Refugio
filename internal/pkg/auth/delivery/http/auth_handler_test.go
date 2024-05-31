package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"mail/internal/models/response"
	"mail/internal/pkg/utils/constants"

	auth_mock "mail/internal/microservice/auth/mock"
	auth_proto "mail/internal/microservice/auth/proto"
	user_mock "mail/internal/microservice/user/mock"
	session_mock "mail/internal/pkg/session/mock"
)

func TestAuthHandler_Login_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthServiceClient := auth_mock.NewMockAuthServiceClient(ctrl)
	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)

	ah := &AuthHandler{
		Sessions:          mockSessionsManager,
		AuthServiceClient: mockAuthServiceClient,
	}

	reqBody := `{"login": "user@mailhub.su", "password": "password123"}`

	req := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(reqBody))
	ctx := context.WithValue(req.Context(), interface{}(string(constants.RequestIDKey)), "testID")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	mockAuthServiceClient.EXPECT().
		Login(gomock.Any(), gomock.Any()).
		Return(&auth_proto.LoginReply{SessionId: "123"}, nil)

	mockSessionsManager.EXPECT().
		SetSession("123", w, req, req.Context()).
		Return(nil)

	ah.Login(w, req)

	resp := w.Result()

	var responseBody response.Response
	err := json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAuthHandler_Login_InvalidRequestJson(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthServiceClient := auth_mock.NewMockAuthServiceClient(ctrl)
	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)

	ah := &AuthHandler{
		Sessions:          mockSessionsManager,
		AuthServiceClient: mockAuthServiceClient,
	}

	reqBody := `{"invalid": json"}`

	req := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(reqBody))
	ctx := context.WithValue(req.Context(), constants.RequestIDKey, "testID")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	ah.Login(w, req)

	resp := w.Result()

	var responseBody response.ErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestAuthHandler_Login_InvalidRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthServiceClient := auth_mock.NewMockAuthServiceClient(ctrl)
	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)

	ah := &AuthHandler{
		Sessions:          mockSessionsManager,
		AuthServiceClient: mockAuthServiceClient,
	}

	reqBody := `{"invalid": "json"}`

	req := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(reqBody))
	ctx := context.WithValue(req.Context(), constants.RequestIDKey, "testID")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	ah.Login(w, req)

	resp := w.Result()

	var responseBody response.ErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestAuthHandler_Login_InvalidRequestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthServiceClient := auth_mock.NewMockAuthServiceClient(ctrl)
	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)

	ah := &AuthHandler{
		Sessions:          mockSessionsManager,
		AuthServiceClient: mockAuthServiceClient,
	}

	reqBody := `{"login": "user@mail.ru", "password": "password123"}`

	req := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(reqBody))
	ctx := context.WithValue(req.Context(), constants.RequestIDKey, "testID")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	ah.Login(w, req)

	resp := w.Result()

	var responseBody response.ErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestAuthHandler_Login_LoginFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthServiceClient := auth_mock.NewMockAuthServiceClient(ctrl)
	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)

	ah := &AuthHandler{
		Sessions:          mockSessionsManager,
		AuthServiceClient: mockAuthServiceClient,
	}

	reqBody := `{"login": "user@mailhub.su", "password": "password123"}`

	req := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(reqBody))
	ctx := context.WithValue(req.Context(), interface{}(string(constants.RequestIDKey)), "testID")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	mockAuthServiceClient.EXPECT().
		Login(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("login failed"))

	ah.Login(w, req)

	resp := w.Result()

	var responseBody response.ErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAuthHandler_Login_SessionFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthServiceClient := auth_mock.NewMockAuthServiceClient(ctrl)
	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)

	ah := &AuthHandler{
		Sessions:          mockSessionsManager,
		AuthServiceClient: mockAuthServiceClient,
	}

	reqBody := `{"login": "user@mailhub.su", "password": "password123"}`

	req := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(reqBody))
	ctx := context.WithValue(req.Context(), interface{}(string(constants.RequestIDKey)), "testID")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	mockAuthServiceClient.EXPECT().
		Login(gomock.Any(), gomock.Any()).
		Return(&auth_proto.LoginReply{SessionId: "123"}, nil)

	mockSessionsManager.EXPECT().
		SetSession("123", w, req, req.Context()).
		Return(errors.New("login failed"))

	ah.Login(w, req)

	resp := w.Result()

	var responseBody response.Response
	err := json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestAuthHandler_Signup_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthServiceClient := auth_mock.NewMockAuthServiceClient(ctrl)
	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockUserServiceClient := user_mock.NewMockUserServiceClient(ctrl)

	ah := &AuthHandler{
		Sessions:          mockSessionsManager,
		AuthServiceClient: mockAuthServiceClient,
		UserServiceClient: mockUserServiceClient,
	}

	reqBody := `{"login": "user@mailhub.su", "password": "password123", "firstname": "John", "surname": "Doe", "gender": "Male", "phoneNumber": "123456789"}`

	req := httptest.NewRequest("POST", "/api/v1/auth/signup", strings.NewReader(reqBody))
	ctx := context.WithValue(req.Context(), interface{}(string(constants.RequestIDKey)), "testID")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	mockUserServiceClient.EXPECT().
		GetUserByOnlyLogin(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("login failed"))

	mockAuthServiceClient.EXPECT().
		Signup(gomock.Any(), gomock.Any()).
		Return(nil, nil)

	ah.Signup(w, req)

	resp := w.Result()

	var responseBody response.Response
	err := json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAuthHandler_Signup_InvalidRequestBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthServiceClient := auth_mock.NewMockAuthServiceClient(ctrl)
	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)

	ah := &AuthHandler{
		Sessions:          mockSessionsManager,
		AuthServiceClient: mockAuthServiceClient,
	}

	reqBody := `{"invalid": json"}`

	req := httptest.NewRequest("POST", "/api/v1/auth/signup", strings.NewReader(reqBody))
	ctx := context.WithValue(req.Context(), constants.RequestIDKey, "testID")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	ah.Signup(w, req)

	resp := w.Result()

	var responseBody response.ErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestAuthHandler_Signup_InvalidRequestFields(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthServiceClient := auth_mock.NewMockAuthServiceClient(ctrl)
	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)

	ah := &AuthHandler{
		Sessions:          mockSessionsManager,
		AuthServiceClient: mockAuthServiceClient,
	}

	reqBody := `{"login": "", "password": "", "firstName": "", "surname": "", "gender": "", "phoneNumber": ""}`

	req := httptest.NewRequest("POST", "/api/v1/auth/signup", strings.NewReader(reqBody))
	ctx := context.WithValue(req.Context(), constants.RequestIDKey, "testID")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	ah.Signup(w, req)

	resp := w.Result()

	var responseBody response.ErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestAuthHandler_Signup_InvalidLoginFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthServiceClient := auth_mock.NewMockAuthServiceClient(ctrl)
	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)

	ah := &AuthHandler{
		Sessions:          mockSessionsManager,
		AuthServiceClient: mockAuthServiceClient,
	}

	reqBody := `{"login": "invalid_email", "password": "password123", "firstName": "John", "surname": "Doe", "gender": "Male", "phoneNumber": "123456789"}`

	req := httptest.NewRequest("POST", "/api/v1/auth/signup", strings.NewReader(reqBody))
	ctx := context.WithValue(req.Context(), constants.RequestIDKey, "testID")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	ah.Signup(w, req)

	resp := w.Result()

	var responseBody response.ErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestAuthHandler_Signup_FailedToAddUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthServiceClient := auth_mock.NewMockAuthServiceClient(ctrl)
	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockUserServiceClient := user_mock.NewMockUserServiceClient(ctrl)

	ah := &AuthHandler{
		Sessions:          mockSessionsManager,
		AuthServiceClient: mockAuthServiceClient,
		UserServiceClient: mockUserServiceClient,
	}

	reqBody := `{"login": "user@mailhub.su", "password": "password123", "firstname": "John", "surname": "Doe", "gender": "Male", "phoneNumber": "123456789"}`

	req := httptest.NewRequest("POST", "/api/v1/auth/signup", strings.NewReader(reqBody))
	ctx := context.WithValue(req.Context(), interface{}(string(constants.RequestIDKey)), "testID")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	mockUserServiceClient.EXPECT().
		GetUserByOnlyLogin(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("login failed"))

	mockAuthServiceClient.EXPECT().
		Signup(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("failed to add user"))

	ah.Signup(w, req)

	resp := w.Result()

	var responseBody response.ErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestAuthHandler_Logout_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)

	ah := &AuthHandler{
		Sessions: mockSessionsManager,
	}

	req := httptest.NewRequest("POST", "/api/v1/auth/logout", nil)
	ctx := context.WithValue(req.Context(), constants.RequestIDKey, "testID")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	mockSessionsManager.EXPECT().
		DestroyCurrent(w, req, req.Context()).
		Return(nil)

	ah.Logout(w, req)

	resp := w.Result()

	var responseBody response.Response
	err := json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAuthHandler_Logout_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)

	ah := &AuthHandler{
		Sessions: mockSessionsManager,
	}

	req := httptest.NewRequest("POST", "/api/v1/auth/logout", nil)
	ctx := context.WithValue(req.Context(), constants.RequestIDKey, "testID")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	mockSessionsManager.EXPECT().
		DestroyCurrent(w, req, req.Context()).
		Return(errors.New("not authorized"))

	ah.Logout(w, req)

	resp := w.Result()

	var responseBody response.ErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
