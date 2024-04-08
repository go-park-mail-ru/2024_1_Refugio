package auth

import (
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	api "mail/pkg/delivery/models"
	"mail/pkg/delivery/session"
	"mail/pkg/domain/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthHandler_Login_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)
	mockSessionsManager := session.NewSessionsManager()

	authHandler := AuthHandler{mockUserUseCase, mockSessionsManager}

	// Prepare request body
	credentials := api.User{Login: "test@example.com", Password: "password"}
	requestBody, _ := json.Marshal(credentials)

	// Mocking successful GetUserByLogin call
	mockUserUseCase.EXPECT().GetUserByLogin(credentials.Login, credentials.Password, gomock.Any()).Return(&models.User{ID: 1}, nil)

	// Mocking successful Create session call
	mockSessionsManager.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)

	// Create a request with JSON body
	req, err := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(authHandler.Login)

	// Serve the HTTP request to the ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	expected := `{"Success":"Login successful"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestAuthHandler_Login_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)
	mockSessionsManager := mock.NewMockSessionsManager(ctrl)

	authHandler := delivery.NewAuthHandler(mockUserUseCase, mockSessionsManager)

	// Prepare request body
	credentials := api.User{Login: "test@example.com", Password: "password"}
	requestBody, _ := json.Marshal(credentials)

	// Mocking unsuccessful GetUserByLogin call
	mockUserUseCase.EXPECT().GetUserByLogin(credentials.Login, credentials.Password, gomock.Any()).Return(nil, errors.New("user not found"))

	// Create a request with JSON body
	req, err := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(authHandler.Login)

	// Serve the HTTP request to the ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}

	// Check the response body
	expected := `{"Error":"Login failed"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
