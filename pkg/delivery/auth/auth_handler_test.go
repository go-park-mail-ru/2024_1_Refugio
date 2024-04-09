package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	api "mail/pkg/delivery/models"
	"mail/pkg/domain/mock"
	domain "mail/pkg/domain/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)
	mockSessionsManager := mock.NewMockSessionsManager(ctrl)

	authHandler := AuthHandler{
		UserUseCase: mockUserUseCase,
		Sessions:    mockSessionsManager,
	}

	requestBody := api.User{
		Login:    "test@mailhub.su",
		Password: "password",
	}
	requestBodyBytes, _ := json.Marshal(requestBody)

	mockUserUseCase.EXPECT().GetUserByLogin("test@mailhub.su", "password", gomock.Any()).Return(&domain.User{ID: 1}, nil)
	mockSessionsManager.EXPECT().Create(gomock.Any(), uint32(1), gomock.Any()).Return(nil, nil)

	req, err := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(requestBodyBytes))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	http.HandlerFunc(authHandler.Login).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expectedResponseBody := `{"status":200,"body":{"Success":"Login successful"}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}

func TestLogin_InvalidRequestBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)
	mockSessionsManager := mock.NewMockSessionsManager(ctrl)

	authHandler := AuthHandler{
		UserUseCase: mockUserUseCase,
		Sessions:    mockSessionsManager,
	}

	requestBody := api.User{
		Login:    "test@mailhub.su",
		Password: "password",
	}
	requestBodyBytesMarshal, _ := json.Marshal(requestBody)
	fmt.Println(string(requestBodyBytesMarshal))
	requestBodyBytes := []byte(`{"birthday":"0001-01-01T00:00:00Z","login":"test@mailhub.su""password":"password"}`)

	req, err := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(requestBodyBytes))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	http.HandlerFunc(authHandler.Login).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	expectedResponseBody := `{"status":400,"body":{"error":"Invalid request body"}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}

func TestLogin_AllFieldsMustBeFilledIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)
	mockSessionsManager := mock.NewMockSessionsManager(ctrl)

	authHandler := AuthHandler{
		UserUseCase: mockUserUseCase,
		Sessions:    mockSessionsManager,
	}

	requestBody := api.User{
		Login:    "test@mailhub.su",
		Password: "",
	}
	requestBodyBytes, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(requestBodyBytes))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	http.HandlerFunc(authHandler.Login).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	expectedResponseBody := `{"status":500,"body":{"error":"All fields must be filled in"}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}

func TestLogin_DomainInLoginIsNotSuitable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)
	mockSessionsManager := mock.NewMockSessionsManager(ctrl)

	authHandler := AuthHandler{
		UserUseCase: mockUserUseCase,
		Sessions:    mockSessionsManager,
	}

	requestBody := api.User{
		Login:    "test@mailhub.ru",
		Password: "password",
	}
	requestBodyBytes, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(requestBodyBytes))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	http.HandlerFunc(authHandler.Login).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	expectedResponseBody := `{"status":400,"body":{"error":"Domain in the login is not suitable"}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}

func TestLogin_UserWithLoginNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)
	mockSessionsManager := mock.NewMockSessionsManager(ctrl)

	authHandler := AuthHandler{
		UserUseCase: mockUserUseCase,
		Sessions:    mockSessionsManager,
	}

	requestBody := api.User{
		Login:    "test@mailhub.su",
		Password: "password",
	}
	requestBodyBytes, _ := json.Marshal(requestBody)

	mockUserUseCase.EXPECT().GetUserByLogin("test@mailhub.su", "password", gomock.Any()).Return(nil, errors.New("user with login test@mailhub.su not found"))

	req, err := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(requestBodyBytes))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	http.HandlerFunc(authHandler.Login).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	expectedResponseBody := `{"status":401,"body":{"error":"Login failed"}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}

func TestLogin_FailedToCreateSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)
	mockSessionsManager := mock.NewMockSessionsManager(ctrl)

	authHandler := AuthHandler{
		UserUseCase: mockUserUseCase,
		Sessions:    mockSessionsManager,
	}

	requestBody := api.User{
		Login:    "test@mailhub.su",
		Password: "password",
	}
	requestBodyBytes, _ := json.Marshal(requestBody)

	mockUserUseCase.EXPECT().GetUserByLogin("test@mailhub.su", "password", gomock.Any()).Return(&domain.User{ID: 1}, nil)
	mockSessionsManager.EXPECT().Create(gomock.Any(), uint32(1), gomock.Any()).Return(nil, errors.New("session already exist"))

	req, err := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(requestBodyBytes))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	http.HandlerFunc(authHandler.Login).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	expectedResponseBody := `{"status":500,"body":{"error":"Failed to create session"}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}

func TestSignup_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)
	mockSessionsManager := mock.NewMockSessionsManager(ctrl)

	authHandler := AuthHandler{
		UserUseCase: mockUserUseCase,
		Sessions:    mockSessionsManager,
	}

	requestBody := api.User{
		Login:     "test@mailhub.su",
		Password:  "password",
		FirstName: "John",
		Surname:   "Doe",
		Gender:    "Male",
	}
	requestBodyBytes, _ := json.Marshal(requestBody)

	mockUserUseCase.EXPECT().IsLoginUnique("test@mailhub.su", gomock.Any()).Return(true, nil)
	mockUserUseCase.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(&domain.User{ID: 1}, nil)

	req, err := http.NewRequest("POST", "/api/v1/auth/signup", bytes.NewReader(requestBodyBytes))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	http.HandlerFunc(authHandler.Signup).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expectedResponseBody := `{"status":200,"body":{"Success":"Signup successful"}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}

func TestSignup_InvalidRequestBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)
	mockSessionsManager := mock.NewMockSessionsManager(ctrl)

	authHandler := AuthHandler{
		UserUseCase: mockUserUseCase,
		Sessions:    mockSessionsManager,
	}

	requestBody := api.User{
		Login:     "test@mailhub.su",
		Password:  "password",
		FirstName: "John",
		Surname:   "Doe",
		Gender:    "Male",
	}
	requestBodyBytesMarshal, _ := json.Marshal(requestBody)
	fmt.Println(string(requestBodyBytesMarshal))
	requestBodyBytes := []byte(`{"firstname":"John","surname":"Doe","gender":"Male","birthday":"0001-01-01T00:00:00Z""login":"test@mailhub.su","password":"password"}`)

	req, err := http.NewRequest("POST", "/api/v1/auth/signup", bytes.NewReader(requestBodyBytes))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	http.HandlerFunc(authHandler.Signup).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	expectedResponseBody := `{"status":400,"body":{"error":"Invalid request body"}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}

func TestSignup_AllFieldsMustBeFilledIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)
	mockSessionsManager := mock.NewMockSessionsManager(ctrl)

	authHandler := AuthHandler{
		UserUseCase: mockUserUseCase,
		Sessions:    mockSessionsManager,
	}

	requestBody := api.User{
		Login:     "test@mailhub.su",
		Password:  "password",
		FirstName: "John",
		Surname:   "Doe",
		Gender:    "male",
	}
	requestBodyBytes, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "/api/v1/auth/signup", bytes.NewReader(requestBodyBytes))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	http.HandlerFunc(authHandler.Signup).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	expectedResponseBody := `{"status":400,"body":{"error":"All fields must be filled in"}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}

func TestSignup_DomainInLoginIsNotSuitable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)
	mockSessionsManager := mock.NewMockSessionsManager(ctrl)

	authHandler := AuthHandler{
		UserUseCase: mockUserUseCase,
		Sessions:    mockSessionsManager,
	}

	requestBody := api.User{
		Login:     "test@mailhub.ru",
		Password:  "password",
		FirstName: "John",
		Surname:   "Doe",
		Gender:    "Male",
	}
	requestBodyBytes, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "/api/v1/auth/signup", bytes.NewReader(requestBodyBytes))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	http.HandlerFunc(authHandler.Signup).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	expectedResponseBody := `{"status":400,"body":{"error":"Domain in the login is not suitable"}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}

func TestSignup_SuchLoginAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)
	mockSessionsManager := mock.NewMockSessionsManager(ctrl)

	authHandler := AuthHandler{
		UserUseCase: mockUserUseCase,
		Sessions:    mockSessionsManager,
	}

	requestBody := api.User{
		Login:     "test@mailhub.su",
		Password:  "password",
		FirstName: "John",
		Surname:   "Doe",
		Gender:    "Male",
	}
	requestBodyBytes, _ := json.Marshal(requestBody)

	mockUserUseCase.EXPECT().IsLoginUnique("test@mailhub.su", gomock.Any()).Return(false, nil)

	req, err := http.NewRequest("POST", "/api/v1/auth/signup", bytes.NewReader(requestBodyBytes))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	http.HandlerFunc(authHandler.Signup).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	expectedResponseBody := `{"status":400,"body":{"error":"Such a login already exists"}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}

func TestSignup_FailedToAddUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)
	mockSessionsManager := mock.NewMockSessionsManager(ctrl)

	authHandler := AuthHandler{
		UserUseCase: mockUserUseCase,
		Sessions:    mockSessionsManager,
	}

	requestBody := api.User{
		Login:     "test@mailhub.su",
		Password:  "password",
		FirstName: "John",
		Surname:   "Doe",
		Gender:    "Male",
	}
	requestBodyBytes, _ := json.Marshal(requestBody)

	mockUserUseCase.EXPECT().IsLoginUnique("test@mailhub.su", gomock.Any()).Return(true, nil)
	mockUserUseCase.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil, errors.New("user with login test@mailhub.su fail"))

	req, err := http.NewRequest("POST", "/api/v1/auth/signup", bytes.NewReader(requestBodyBytes))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	http.HandlerFunc(authHandler.Signup).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	expectedResponseBody := `{"status":500,"body":{"error":"Failed to add user"}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}

func TestLogout_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := mock.NewMockSessionsManager(ctrl)

	authHandler := AuthHandler{
		Sessions: mockSessionsManager,
	}

	req, err := http.NewRequest("POST", "/api/v1/auth/logout", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	mockSessionsManager.EXPECT().DestroyCurrent(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	http.HandlerFunc(authHandler.Logout).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expectedResponseBody := `{"status":200,"body":{"Success":"Logout successful"}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}

func TestLogout_NotAuthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := mock.NewMockSessionsManager(ctrl)

	authHandler := AuthHandler{
		Sessions: mockSessionsManager,
	}

	req, err := http.NewRequest("POST", "/api/v1/auth/logout", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	mockSessionsManager.EXPECT().DestroyCurrent(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("Not Authorized"))

	http.HandlerFunc(authHandler.Logout).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	expectedResponseBody := `{"status":401,"body":{"error":"Not Authorized"}}` + "\n"
	assert.Equal(t, expectedResponseBody, rr.Body.String())
}

func TestSanitizeString(t *testing.T) {
	input := "<script>alert('Hello, World!');</script>"

	expected := ""

	output := sanitizeString(input)
	if output != expected {
		t.Errorf("Expected: %s, Got: %s", expected, output)
	}
}

func TestIsEmpty(t *testing.T) {
	emptyStr := ""
	if !isEmpty(emptyStr) {
		t.Error("Expected empty string")
	}

	nonEmptyStr := "Hello, World!"
	if isEmpty(nonEmptyStr) {
		t.Error("Expected non-empty string")
	}
}

func TestIsValidEmailFormat(t *testing.T) {
	validEmail := "test@mailhub.su"
	if !isValidEmailFormat(validEmail) {
		t.Error("Expected valid email format")
	}

	invalidEmail := "invalid.email@example.com"
	if isValidEmailFormat(invalidEmail) {
		t.Error("Expected invalid email format")
	}
}
