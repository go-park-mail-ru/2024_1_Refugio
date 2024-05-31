package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	userMock "mail/internal/microservice/user/mock"
	userProto "mail/internal/microservice/user/proto"
	api "mail/internal/models/delivery_models"
	sessionMock "mail/internal/pkg/session/mock"
)

func TestVerifyAuth_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserServiceClient := userMock.NewMockUserServiceClient(ctrl)
	mockSessionsManager := sessionMock.NewMockSessionsManager(ctrl)

	userHandler := UserHandler{
		Sessions:          mockSessionsManager,
		UserServiceClient: mockUserServiceClient,
	}

	req, err := http.NewRequest("GET", "/api/v1/verify-auth", nil)
	assert.NoError(t, err)
	ctx := context.WithValue(req.Context(), "requestid", "testID")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	mockSessionsManager.EXPECT().GetSession(req, gomock.Any()).Return(&api.Session{CsrfToken: "csrf"})

	http.HandlerFunc(userHandler.VerifyAuth).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	csrfToken := rr.Header().Get("X-Csrf-Token")
	assert.NotEmpty(t, csrfToken)

	expectedResponseBody := `{"status":200,"body":{"Success":"OK"}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}

func TestGetUserBySession_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserServiceClient := userMock.NewMockUserServiceClient(ctrl)
	mockSessionsManager := sessionMock.NewMockSessionsManager(ctrl)

	userHandler := UserHandler{
		Sessions:          mockSessionsManager,
		UserServiceClient: mockUserServiceClient,
	}

	sessionData := &api.Session{
		UserID: 1,
	}

	req, err := http.NewRequest("GET", "/api/v1/user/get", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.WithValue(req.Context(), "requestid", "testID")
	req = req.WithContext(ctx)

	mockSessionsManager.EXPECT().GetSession(req, gomock.Any()).Return(sessionData)

	userData := &userProto.GetUserReply{User: &userProto.User{
		Id:          1,
		Login:       "test@mailhub.su",
		Firstname:   "John",
		Surname:     "Doe",
		Patronymic:  "Middle",
		Avatar:      "123",
		PhoneNumber: "1234567890",
		Description: "Test description",
	}}
	mockUserServiceClient.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(userData, nil)

	rr := httptest.NewRecorder()

	userHandler.GetUserBySession(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expectedResponseBody := `{"status":200,"body":{"user":{"id":1,"firstname":"John","surname":"Doe","middlename":"Middle","birthday":"1970-01-01T00:00:00Z","login":"test@mailhub.su","password":"","avatar":"123","phonenumber":"1234567890","description":"Test description"}}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}

func TestGetUserBySession_InternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserServiceClient := userMock.NewMockUserServiceClient(ctrl)
	mockSessionsManager := sessionMock.NewMockSessionsManager(ctrl)

	userHandler := UserHandler{
		Sessions:          mockSessionsManager,
		UserServiceClient: mockUserServiceClient,
	}

	sessionData := &api.Session{
		UserID: 1,
	}

	req, err := http.NewRequest("GET", "/api/v1/user/get", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.WithValue(req.Context(), "requestid", "testID")
	req = req.WithContext(ctx)

	mockSessionsManager.EXPECT().GetSession(req, gomock.Any()).Return(sessionData)

	mockUserServiceClient.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(nil, errors.New("user not found"))

	rr := httptest.NewRecorder()

	userHandler.GetUserBySession(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	expectedResponseBody := `{"status":500,"body":{"error":"Internal Server Error"}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}

func TestUpdateUserData_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserServiceClient := userMock.NewMockUserServiceClient(ctrl)
	mockSessionsManager := sessionMock.NewMockSessionsManager(ctrl)

	userHandler := UserHandler{
		Sessions:          mockSessionsManager,
		UserServiceClient: mockUserServiceClient,
	}

	requestBody := api.User{
		ID:          1,
		Login:       "updatedlogin@mailhub.su",
		FirstName:   "Updated First Name",
		Surname:     "Updated Surname",
		Patronymic:  "Updated Patronymic",
		AvatarID:    "updated-avatar-id",
		PhoneNumber: "1234567890",
		Description: "Updated Description",
	}
	requestBodyBytes, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("PUT", "/api/v1/user/update", bytes.NewReader(requestBodyBytes))
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.WithValue(req.Context(), "requestid", "testID")
	req = req.WithContext(ctx)

	sessionData := &api.Session{
		UserID: 1,
	}
	mockSessionsManager.EXPECT().GetSession(req, gomock.Any()).Return(sessionData)

	updatedUser := &userProto.UpdateUserReply{User: &userProto.User{
		Id:          1,
		Login:       "updatedlogin@mailhub.su",
		Firstname:   "Updated First Name",
		Surname:     "Updated Surname",
		Patronymic:  "Updated Patronymic",
		Avatar:      "updated-avatar-id",
		PhoneNumber: "1234567890",
		Description: "Updated Description",
	}}
	mockUserServiceClient.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(updatedUser, nil)

	rr := httptest.NewRecorder()

	userHandler.UpdateUserData(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expectedResponseBody := `{"status":200,"body":{"user":{"id":1,"firstname":"Updated First Name","surname":"Updated Surname","middlename":"Updated Patronymic","birthday":"1970-01-01T00:00:00Z","login":"updatedlogin@mailhub.su","password":"","avatar":"updated-avatar-id","phonenumber":"1234567890","description":"Updated Description"}}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}

func TestUpdateUserData_InvalidRequestBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserServiceClient := userMock.NewMockUserServiceClient(ctrl)
	mockSessionsManager := sessionMock.NewMockSessionsManager(ctrl)

	userHandler := UserHandler{
		Sessions:          mockSessionsManager,
		UserServiceClient: mockUserServiceClient,
	}

	requestBody := api.User{
		ID:          1,
		Login:       "updatedlogin@mailhub.su",
		FirstName:   "Updated First Name",
		Surname:     "Updated Surname",
		Patronymic:  "Updated Patronymic",
		AvatarID:    "updated-avatar-id",
		PhoneNumber: "1234567890",
		Description: "Updated Description",
	}
	requestBodyBytesMarshal, _ := json.Marshal(requestBody)
	fmt.Println(string(requestBodyBytesMarshal))
	requestBodyBytes := []byte(`{"id":1,"firstname":"Updated First Name""surname":"Updated Surname","middlename":"Updated Patronymic","birthday":"0001-01-01T00:00:00Z","login":"updatedlogin@mailhub.su","password":"","avatar":"updated-avatar-id","phonenumber":"1234567890","description":"Updated Description"}`)

	req, err := http.NewRequest("PUT", "/api/v1/user/update", bytes.NewReader(requestBodyBytes))
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.WithValue(req.Context(), "requestid", "testID")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	userHandler.UpdateUserData(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	expectedResponseBody := `{"status":400,"body":{"error":"Invalid request body"}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}

func TestUpdateUserData_NotAuthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserServiceClient := userMock.NewMockUserServiceClient(ctrl)
	mockSessionsManager := sessionMock.NewMockSessionsManager(ctrl)

	userHandler := UserHandler{
		Sessions:          mockSessionsManager,
		UserServiceClient: mockUserServiceClient,
	}

	requestBody := api.User{
		ID:          1,
		Login:       "updatedlogin@mailhub.su",
		FirstName:   "Updated First Name",
		Surname:     "Updated Surname",
		Patronymic:  "Updated Patronymic",
		AvatarID:    "updated-avatar-id",
		PhoneNumber: "1234567890",
		Description: "Updated Description",
	}
	requestBodyBytes, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("PUT", "/api/v1/user/update", bytes.NewReader(requestBodyBytes))
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.WithValue(req.Context(), "requestid", "testID")
	req = req.WithContext(ctx)

	sessionData := &api.Session{
		UserID: 2,
	}
	mockSessionsManager.EXPECT().GetSession(req, gomock.Any()).Return(sessionData)

	rr := httptest.NewRecorder()

	userHandler.UpdateUserData(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	expectedResponseBody := `{"status":401,"body":{"error":"Not authorized"}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}

func TestUpdateUserData_InternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserServiceClient := userMock.NewMockUserServiceClient(ctrl)
	mockSessionsManager := sessionMock.NewMockSessionsManager(ctrl)

	userHandler := UserHandler{
		Sessions:          mockSessionsManager,
		UserServiceClient: mockUserServiceClient,
	}

	requestBody := api.User{
		ID:          1,
		Login:       "updatedlogin@mailhub.su",
		FirstName:   "Updated First Name",
		Surname:     "Updated Surname",
		Patronymic:  "Updated Patronymic",
		AvatarID:    "updated-avatar-id",
		PhoneNumber: "1234567890",
		Description: "Updated Description",
	}
	requestBodyBytes, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("PUT", "/api/v1/user/update", bytes.NewReader(requestBodyBytes))
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.WithValue(req.Context(), "requestid", "testID")
	req = req.WithContext(ctx)

	sessionData := &api.Session{
		UserID: 1,
	}
	mockSessionsManager.EXPECT().GetSession(req, gomock.Any()).Return(sessionData)

	mockUserServiceClient.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(nil, errors.New("user no update"))

	rr := httptest.NewRecorder()

	userHandler.UpdateUserData(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	expectedResponseBody := `{"status":500,"body":{"error":"Internal Server Error"}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}

func TestDeleteUserData_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserServiceClient := userMock.NewMockUserServiceClient(ctrl)
	mockSessionsManager := sessionMock.NewMockSessionsManager(ctrl)

	userHandler := UserHandler{
		Sessions:          mockSessionsManager,
		UserServiceClient: mockUserServiceClient,
	}

	req, err := http.NewRequest("DELETE", "/api/v1/user/delete/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	ctx := context.WithValue(req.Context(), "requestid", "testID")
	req = req.WithContext(ctx)

	sessionData := &api.Session{
		UserID: 1,
	}
	mockSessionsManager.EXPECT().GetSession(req, gomock.Any()).Return(sessionData)

	mockUserServiceClient.EXPECT().DeleteUserById(gomock.Any(), gomock.Any()).Return(&userProto.DeleteUserByIdReply{Status: true}, nil)

	rr := httptest.NewRecorder()

	userHandler.DeleteUserData(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expectedResponseBody := `{"status":200,"body":{"message":"User data deleted successfully"}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}

func TestDeleteUserData_BadIdInRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserServiceClient := userMock.NewMockUserServiceClient(ctrl)
	mockSessionsManager := sessionMock.NewMockSessionsManager(ctrl)

	userHandler := UserHandler{
		Sessions:          mockSessionsManager,
		UserServiceClient: mockUserServiceClient,
	}

	req, err := http.NewRequest("DELETE", "/api/v1/user/delete/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": "i"})
	ctx := context.WithValue(req.Context(), "requestid", "testID")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	userHandler.DeleteUserData(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	expectedResponseBody := `{"status":400,"body":{"error":"Bad id in request"}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}

func TestDeleteUserData_NotAuthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserServiceClient := userMock.NewMockUserServiceClient(ctrl)
	mockSessionsManager := sessionMock.NewMockSessionsManager(ctrl)

	userHandler := UserHandler{
		Sessions:          mockSessionsManager,
		UserServiceClient: mockUserServiceClient,
	}

	req, err := http.NewRequest("DELETE", "/api/v1/user/delete/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	ctx := context.WithValue(req.Context(), "requestid", "testID")
	req = req.WithContext(ctx)

	sessionData := &api.Session{
		UserID: 2,
	}
	mockSessionsManager.EXPECT().GetSession(req, gomock.Any()).Return(sessionData)

	rr := httptest.NewRecorder()

	userHandler.DeleteUserData(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	expectedResponseBody := `{"status":401,"body":{"error":"Not authorized"}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}

func TestDeleteUserData_InternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserServiceClient := userMock.NewMockUserServiceClient(ctrl)
	mockSessionsManager := sessionMock.NewMockSessionsManager(ctrl)

	userHandler := UserHandler{
		Sessions:          mockSessionsManager,
		UserServiceClient: mockUserServiceClient,
	}

	req, err := http.NewRequest("DELETE", "/api/v1/user/delete/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	ctx := context.WithValue(req.Context(), "requestid", "testID")
	req = req.WithContext(ctx)

	sessionData := &api.Session{
		UserID: 1,
	}
	mockSessionsManager.EXPECT().GetSession(req, gomock.Any()).Return(sessionData)

	mockUserServiceClient.EXPECT().DeleteUserById(gomock.Any(), gomock.Any()).Return(&userProto.DeleteUserByIdReply{Status: false}, errors.New("fail with delete"))

	rr := httptest.NewRecorder()

	userHandler.DeleteUserData(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	expectedResponseBody := `{"status":500,"body":{"error":"Internal Server Error"}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}

func TestDeleteUserData_FailedToDeleteUserData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserServiceClient := userMock.NewMockUserServiceClient(ctrl)
	mockSessionsManager := sessionMock.NewMockSessionsManager(ctrl)

	userHandler := UserHandler{
		Sessions:          mockSessionsManager,
		UserServiceClient: mockUserServiceClient,
	}

	req, err := http.NewRequest("DELETE", "/api/v1/user/delete/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	ctx := context.WithValue(req.Context(), "requestid", "testID")
	req = req.WithContext(ctx)

	sessionData := &api.Session{
		UserID: 1,
	}
	mockSessionsManager.EXPECT().GetSession(req, gomock.Any()).Return(sessionData)

	mockUserServiceClient.EXPECT().DeleteUserById(gomock.Any(), gomock.Any()).Return(&userProto.DeleteUserByIdReply{Status: false}, nil)

	rr := httptest.NewRecorder()

	userHandler.DeleteUserData(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	expectedResponseBody := `{"status":500,"body":{"error":"Failed to delete user data"}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}

func TestDeleteUserAvatar_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserServiceClient := userMock.NewMockUserServiceClient(ctrl)
	mockSessionsManager := sessionMock.NewMockSessionsManager(ctrl)

	userHandler := UserHandler{
		Sessions:          mockSessionsManager,
		UserServiceClient: mockUserServiceClient,
	}

	req, err := http.NewRequest("DELETE", "/api/v1/user/delete/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.WithValue(req.Context(), "requestid", "testID")
	req = req.WithContext(ctx)

	sessionData := &api.Session{
		UserID: 1,
	}
	mockSessionsManager.EXPECT().GetSession(req, gomock.Any()).Return(sessionData)

	mockUserServiceClient.EXPECT().DeleteUserAvatar(gomock.Any(), gomock.Any()).Return(&userProto.DeleteUserAvatarReply{Status: true}, nil)

	rr := httptest.NewRecorder()

	userHandler.DeleteUserAvatar(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expectedResponseBody := `{"status":200,"body":{"message":"User avatar deleted successfully"}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}

func TestDeleteUserAvatar_InternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserServiceClient := userMock.NewMockUserServiceClient(ctrl)
	mockSessionsManager := sessionMock.NewMockSessionsManager(ctrl)

	userHandler := UserHandler{
		Sessions:          mockSessionsManager,
		UserServiceClient: mockUserServiceClient,
	}

	req, err := http.NewRequest("DELETE", "/api/v1/user/delete/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.WithValue(req.Context(), "requestid", "testID")
	req = req.WithContext(ctx)

	sessionData := &api.Session{
		UserID: 1,
	}
	mockSessionsManager.EXPECT().GetSession(req, gomock.Any()).Return(sessionData)

	mockUserServiceClient.EXPECT().DeleteUserAvatar(gomock.Any(), gomock.Any()).Return(&userProto.DeleteUserAvatarReply{Status: false}, errors.New("fail with delete"))

	rr := httptest.NewRecorder()

	userHandler.DeleteUserAvatar(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	expectedResponseBody := `{"status":500,"body":{"error":"Internal Server Error"}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}

func TestDeleteUserAvatar_FailedToDeleteUserData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserServiceClient := userMock.NewMockUserServiceClient(ctrl)
	mockSessionsManager := sessionMock.NewMockSessionsManager(ctrl)

	userHandler := UserHandler{
		Sessions:          mockSessionsManager,
		UserServiceClient: mockUserServiceClient,
	}

	req, err := http.NewRequest("DELETE", "/api/v1/user/delete/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.WithValue(req.Context(), "requestid", "testID")
	req = req.WithContext(ctx)

	sessionData := &api.Session{
		UserID: 1,
	}
	mockSessionsManager.EXPECT().GetSession(req, gomock.Any()).Return(sessionData)

	mockUserServiceClient.EXPECT().DeleteUserAvatar(gomock.Any(), gomock.Any()).Return(&userProto.DeleteUserAvatarReply{Status: false}, nil)

	rr := httptest.NewRecorder()

	userHandler.DeleteUserAvatar(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	expectedResponseBody := `{"status":500,"body":{"error":"Failed to delete user avatar"}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}
