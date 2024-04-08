package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	api "mail/pkg/delivery/models"
	"mail/pkg/domain/mock"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	domain "mail/pkg/domain/models"
)

func TestGetUserBySession_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)
	mockSessionManager := mock.NewMockSessionsManager(ctrl)

	userHandler := UserHandler{
		UserUseCase: mockUserUseCase,
		Sessions:    mockSessionManager,
	}

	sessionData := &api.Session{
		UserID: 1,
	}

	req, err := http.NewRequest("GET", "/api/v1/user/get", nil)
	if err != nil {
		t.Fatal(err)
	}

	mockSessionManager.EXPECT().GetSession(req, gomock.Any()).Return(sessionData)

	userData := &domain.User{
		ID:          1,
		Login:       "test@mailhub.su",
		FirstName:   "John",
		Surname:     "Doe",
		Patronymic:  "Middle",
		AvatarID:    "123",
		PhoneNumber: "1234567890",
		Description: "Test description",
	}
	mockUserUseCase.EXPECT().GetUserByID(sessionData.UserID, gomock.Any()).Return(userData, nil)

	rr := httptest.NewRecorder()

	userHandler.GetUserBySession(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expectedResponseBody := `{"status":200,"body":{"user":{"id":1,"firstname":"John","surname":"Doe","middlename":"Middle","birthday":"0001-01-01T00:00:00Z","login":"test@mailhub.su","password":"","avatar":"123","phonenumber":"1234567890","description":"Test description"}}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}

func TestGetUserBySession_InternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)
	mockSessionManager := mock.NewMockSessionsManager(ctrl)

	userHandler := UserHandler{
		UserUseCase: mockUserUseCase,
		Sessions:    mockSessionManager,
	}

	sessionData := &api.Session{
		UserID: 1,
	}

	req, err := http.NewRequest("GET", "/api/v1/user/get", nil)
	if err != nil {
		t.Fatal(err)
	}

	mockSessionManager.EXPECT().GetSession(req, gomock.Any()).Return(sessionData)

	mockUserUseCase.EXPECT().GetUserByID(sessionData.UserID, gomock.Any()).Return(nil, errors.New("user no found"))

	rr := httptest.NewRecorder()

	userHandler.GetUserBySession(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	expectedResponseBody := `{"status":500,"body":{"error":"Internal Server Error"}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}

func TestUpdateUserData_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)
	mockSessionManager := mock.NewMockSessionsManager(ctrl)

	userHandler := UserHandler{
		UserUseCase: mockUserUseCase,
		Sessions:    mockSessionManager,
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

	sessionData := &api.Session{
		UserID: 1,
	}
	mockSessionManager.EXPECT().GetSession(req, gomock.Any()).Return(sessionData)

	updatedUser := &domain.User{
		ID:          1,
		Login:       "updatedlogin@mailhub.su",
		FirstName:   "Updated First Name",
		Surname:     "Updated Surname",
		Patronymic:  "Updated Patronymic",
		AvatarID:    "updated-avatar-id",
		PhoneNumber: "1234567890",
		Description: "Updated Description",
	}
	mockUserUseCase.EXPECT().UpdateUser(updatedUser, gomock.Any()).Return(updatedUser, nil)

	rr := httptest.NewRecorder()

	userHandler.UpdateUserData(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expectedResponseBody := `{"status":200,"body":{"user":{"id":1,"firstname":"Updated First Name","surname":"Updated Surname","middlename":"Updated Patronymic","birthday":"0001-01-01T00:00:00Z","login":"updatedlogin@mailhub.su","password":"","avatar":"updated-avatar-id","phonenumber":"1234567890","description":"Updated Description"}}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}

func TestUpdateUserData_InvalidRequestBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)
	mockSessionManager := mock.NewMockSessionsManager(ctrl)

	userHandler := UserHandler{
		UserUseCase: mockUserUseCase,
		Sessions:    mockSessionManager,
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

	rr := httptest.NewRecorder()

	userHandler.UpdateUserData(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	expectedResponseBody := `{"status":400,"body":{"error":"Invalid request body"}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}

func TestUpdateUserData_NotAuthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)
	mockSessionManager := mock.NewMockSessionsManager(ctrl)

	userHandler := UserHandler{
		UserUseCase: mockUserUseCase,
		Sessions:    mockSessionManager,
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

	sessionData := &api.Session{
		UserID: 2,
	}
	mockSessionManager.EXPECT().GetSession(req, gomock.Any()).Return(sessionData)

	rr := httptest.NewRecorder()

	userHandler.UpdateUserData(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	expectedResponseBody := `{"status":401,"body":{"error":"Not authorized"}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}

func TestUpdateUserData_InternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)
	mockSessionManager := mock.NewMockSessionsManager(ctrl)

	userHandler := UserHandler{
		UserUseCase: mockUserUseCase,
		Sessions:    mockSessionManager,
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

	sessionData := &api.Session{
		UserID: 1,
	}
	mockSessionManager.EXPECT().GetSession(req, gomock.Any()).Return(sessionData)

	updatedUser := &domain.User{
		ID:          1,
		Login:       "updatedlogin@mailhub.su",
		FirstName:   "Updated First Name",
		Surname:     "Updated Surname",
		Patronymic:  "Updated Patronymic",
		AvatarID:    "updated-avatar-id",
		PhoneNumber: "1234567890",
		Description: "Updated Description",
	}
	mockUserUseCase.EXPECT().UpdateUser(updatedUser, gomock.Any()).Return(nil, errors.New("user no update"))

	rr := httptest.NewRecorder()

	userHandler.UpdateUserData(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	expectedResponseBody := `{"status":500,"body":{"error":"Internal Server Error"}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}

func TestDeleteUserData_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)
	mockSessionManager := mock.NewMockSessionsManager(ctrl)

	userHandler := UserHandler{
		UserUseCase: mockUserUseCase,
		Sessions:    mockSessionManager,
	}

	req, err := http.NewRequest("DELETE", "/api/v1/user/delete/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	sessionData := &api.Session{
		UserID: 1,
	}
	mockSessionManager.EXPECT().GetSession(req, gomock.Any()).Return(sessionData)

	mockUserUseCase.EXPECT().DeleteUserByID(uint32(1), gomock.Any()).Return(true, nil)

	rr := httptest.NewRecorder()

	userHandler.DeleteUserData(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expectedResponseBody := `{"status":200,"body":{"message":"User data deleted successfully"}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}

func TestDeleteUserData_BadIdInRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)
	mockSessionManager := mock.NewMockSessionsManager(ctrl)

	userHandler := UserHandler{
		UserUseCase: mockUserUseCase,
		Sessions:    mockSessionManager,
	}

	req, err := http.NewRequest("DELETE", "/api/v1/user/delete/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": "i"})

	rr := httptest.NewRecorder()

	userHandler.DeleteUserData(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	expectedResponseBody := `{"status":400,"body":{"error":"Bad id in request"}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}

func TestDeleteUserData_NotAuthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)
	mockSessionManager := mock.NewMockSessionsManager(ctrl)

	userHandler := UserHandler{
		UserUseCase: mockUserUseCase,
		Sessions:    mockSessionManager,
	}

	req, err := http.NewRequest("DELETE", "/api/v1/user/delete/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	sessionData := &api.Session{
		UserID: 2,
	}
	mockSessionManager.EXPECT().GetSession(req, gomock.Any()).Return(sessionData)

	rr := httptest.NewRecorder()

	userHandler.DeleteUserData(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	expectedResponseBody := `{"status":401,"body":{"error":"Not authorized"}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}

func TestDeleteUserData_InternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)
	mockSessionManager := mock.NewMockSessionsManager(ctrl)

	userHandler := UserHandler{
		UserUseCase: mockUserUseCase,
		Sessions:    mockSessionManager,
	}

	req, err := http.NewRequest("DELETE", "/api/v1/user/delete/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	sessionData := &api.Session{
		UserID: 1,
	}
	mockSessionManager.EXPECT().GetSession(req, gomock.Any()).Return(sessionData)

	mockUserUseCase.EXPECT().DeleteUserByID(uint32(1), gomock.Any()).Return(false, errors.New("fail with delete"))

	rr := httptest.NewRecorder()

	userHandler.DeleteUserData(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	expectedResponseBody := `{"status":500,"body":{"error":"Internal Server Error"}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}

func TestDeleteUserData_FailedToDeleteUserData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)
	mockSessionManager := mock.NewMockSessionsManager(ctrl)

	userHandler := UserHandler{
		UserUseCase: mockUserUseCase,
		Sessions:    mockSessionManager,
	}

	req, err := http.NewRequest("DELETE", "/api/v1/user/delete/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	sessionData := &api.Session{
		UserID: 1,
	}
	mockSessionManager.EXPECT().GetSession(req, gomock.Any()).Return(sessionData)

	mockUserUseCase.EXPECT().DeleteUserByID(uint32(1), gomock.Any()).Return(false, nil)

	rr := httptest.NewRecorder()

	userHandler.DeleteUserData(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	expectedResponseBody := `{"status":500,"body":{"error":"Failed to delete user data"}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}

func TestGenerateUniqueFileName(t *testing.T) {
	tests := []struct {
		format string
	}{
		{"_test.txt"},
		{"_data.csv"},
		{"_output.json"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Format_%s", tt.format), func(t *testing.T) {
			uniqueFileName := generateUniqueFileName(tt.format)

			assert.Contains(t, uniqueFileName, tt.format)

			currentTime := time.Now().Format("20060102_150405")
			assert.Contains(t, uniqueFileName, currentTime)

			randomNumStr := uniqueFileName[len(currentTime)+1 : len(currentTime)+4]
			randomNum, err := strconv.Atoi(randomNumStr)
			assert.NoError(t, err)
			assert.True(t, randomNum >= 0 && randomNum <= 999)
		})
	}
}
