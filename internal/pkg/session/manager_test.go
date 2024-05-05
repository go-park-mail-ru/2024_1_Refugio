package session_test

import (
	"errors"
	"fmt"
	domain "mail/internal/microservice/models/domain_models"
	mock "mail/internal/microservice/session/mock"
	converters "mail/internal/models/delivery_converters"
	"mail/internal/pkg/logger"
	session "mail/internal/pkg/session"
	"net/http/httptest"
	"os"

	"context"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func GetCTX() context.Context {
	requestID := "testID"

	f, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create logfile" + "log.txt")
	}
	defer f.Close()

	c := context.WithValue(context.Background(), "logger", logger.InitializationBdLog(f))
	ctx := context.WithValue(c, "requestID", requestID)
	return ctx
}

func TestSessionsManager_GetSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionUseCase := mock.NewMockSessionUseCase(ctrl)
	sessionManager := session.NewSessionsManager(mockSessionUseCase)

	expectedSession := &domain.Session{
		ID:           "session_id",
		UserID:       1,
		CreationDate: time.Now(),
		Device:       "desktop",
		LifeTime:     3600,
		CsrfToken:    "csrf_token",
	}
	ctx := GetCTX()

	mockSessionUseCase.EXPECT().
		GetSession("session_id", ctx).
		Return(expectedSession, nil).
		Times(1)

	req, _ := http.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "session_id"})

	sessionModel := sessionManager.GetSession(req, ctx)

	expectedSessionModel := converters.SessionConvertCoreInApi(*expectedSession)

	if sessionModel.ID != expectedSessionModel.ID {
		t.Errorf("Expected session ID %v, got %v", expectedSessionModel.ID, sessionModel.ID)
	}
}

func TestSessionsManager_GetSession_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionUseCase := mock.NewMockSessionUseCase(ctrl)
	sessionManager := session.NewSessionsManager(mockSessionUseCase)
	ctx := GetCTX()

	mockSessionUseCase.EXPECT().
		GetSession("invalid_session_id", ctx).
		Return(nil, errors.New("session not found")).
		Times(1)

	req, _ := http.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "invalid_session_id"})

	sessionModel := sessionManager.GetSession(req, ctx)

	if sessionModel != nil {
		t.Errorf("Expected nil session, got %v", sessionModel)
	}
}

func TestSessionsManager_Check(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionUseCase := mock.NewMockSessionUseCase(ctrl)
	sessionManager := session.NewSessionsManager(mockSessionUseCase)

	expectedSession := &domain.Session{
		ID:           "session_id",
		UserID:       1,
		CreationDate: time.Now(),
		Device:       "desktop",
		LifeTime:     3600,
		CsrfToken:    "csrf_token",
	}
	ctx := GetCTX()

	mockSessionUseCase.EXPECT().
		GetSession("session_id", ctx).
		Return(expectedSession, nil).
		Times(1)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-CSRF-Token", "csrf_token")
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "session_id"})

	sess, err := sessionManager.Check(req, ctx)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if sess == nil {
		t.Error("Expected non-nil session")
	}
}

func TestSessionsManager_Check_NoCSRFToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionUseCase := mock.NewMockSessionUseCase(ctrl)
	sessionManager := session.NewSessionsManager(mockSessionUseCase)
	ctx := GetCTX()

	req, _ := http.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "session_id"})

	sess, err := sessionManager.Check(req, ctx)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if sess != nil {
		t.Errorf("Expected nil session, got %v", sess)
	}
}

func TestSessionsManager_Check_NoSessionCookie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionUseCase := mock.NewMockSessionUseCase(ctrl)
	sessionManager := session.NewSessionsManager(mockSessionUseCase)
	ctx := GetCTX()

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-CSRF-Token", "csrf_token")

	sess, err := sessionManager.Check(req, ctx)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if sess != nil {
		t.Errorf("Expected nil session, got %v", sess)
	}
}

func TestSessionsManager_Check_SessionNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionUseCase := mock.NewMockSessionUseCase(ctrl)
	sessionManager := session.NewSessionsManager(mockSessionUseCase)
	ctx := GetCTX()

	mockSessionUseCase.EXPECT().
		GetSession("invalid_session_id", ctx).
		Return(nil, errors.New("no session found")).
		Times(1)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-CSRF-Token", "csrf_token")
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "invalid_session_id"})

	sess, err := sessionManager.Check(req, ctx)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if sess != nil {
		t.Errorf("Expected nil session, got %v", sess)
	}
}

func TestSessionsManager_Check_CSRFTokenMismatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionUseCase := mock.NewMockSessionUseCase(ctrl)
	sessionManager := session.NewSessionsManager(mockSessionUseCase)
	ctx := GetCTX()

	mockSessionUseCase.EXPECT().
		GetSession("session_id", ctx).
		Return(&domain.Session{CsrfToken: "csrf_token"}, nil).
		Times(1)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-CSRF-Token", "invalid_csrf_token")
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "session_id"})

	sess, err := sessionManager.Check(req, ctx)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if sess != nil {
		t.Errorf("Expected nil session, got %v", sess)
	}
}

func TestSessionsManager_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionUseCase := mock.NewMockSessionUseCase(ctrl)
	sessionManager := session.NewSessionsManager(mockSessionUseCase)

	expectedSession := &domain.Session{
		ID:        "session_id",
		UserID:    123,
		CsrfToken: "csrf_token",
	}
	ctx := GetCTX()

	mockSessionUseCase.EXPECT().
		CreateNewSession(uint32(123), gomock.Any(), gomock.Any(), gomock.Any()).
		Return("session_id", nil).
		Times(1)

	mockSessionUseCase.EXPECT().
		GetSession("session_id", ctx).
		Return(expectedSession, nil).
		Times(1)

	w := httptest.NewRecorder()

	sess, err := sessionManager.Create(w, 123, ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if sess == nil {
		t.Error("Expected session, got nil")
	}

	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Errorf("Expected 1 cookies, got %d", len(cookies))
	}

	var foundSessionCookie bool
	for _, cookie := range cookies {
		switch cookie.Name {
		case "session_id":
			if cookie.Value == "session_id" {
				foundSessionCookie = true
			}
		}
	}

	if !foundSessionCookie {
		t.Error("Expected session_cookie with name 'session_id' and value 'session_id'")
	}
}

func TestSessionsManager_Create_SessionAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionUseCase := mock.NewMockSessionUseCase(ctrl)
	sessionManager := session.NewSessionsManager(mockSessionUseCase)
	ctx := GetCTX()

	mockSessionUseCase.EXPECT().
		CreateNewSession(uint32(123), gomock.Any(), gomock.Any(), ctx).
		Return("", fmt.Errorf("session already exist")).
		Times(1)

	w := httptest.NewRecorder()

	sess, err := sessionManager.Create(w, 123, ctx)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if sess != nil {
		t.Errorf("Expected nil session, got %+v", sess)
	}

	cookies := w.Result().Cookies()
	if len(cookies) != 0 {
		t.Errorf("Expected 0 cookies, got %d", len(cookies))
	}
}

func TestSessionsManager_DestroyCurrent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionUseCase := mock.NewMockSessionUseCase(ctrl)
	sessionManager := session.NewSessionsManager(mockSessionUseCase)
	ctx := GetCTX()

	mockSessionUseCase.EXPECT().DeleteSession(gomock.Any(), gomock.Any()).Return(nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{Name: "session_id", Value: "session_id_value"})

	err := sessionManager.DestroyCurrent(w, r, ctx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	sessionCookie := w.Result().Cookies()[0]
	expectedExpires := time.Now().AddDate(0, 0, -1)
	if !sessionCookie.Expires.Before(expectedExpires) {
		t.Errorf("Expected session_id cookie to expire in the past, got %v", sessionCookie.Expires)
	}
}

func TestSessionsManager_DestroyCurrent_NoSessionIDCookie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionUseCase := mock.NewMockSessionUseCase(ctrl)
	sessionManager := session.NewSessionsManager(mockSessionUseCase)
	ctx := GetCTX()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	err := sessionManager.DestroyCurrent(w, r, ctx)

	if err == nil || err.Error() != "http: named cookie not present" {
		t.Errorf("Expected 'http: named cookie not present' error, got %v", err)
	}
}
