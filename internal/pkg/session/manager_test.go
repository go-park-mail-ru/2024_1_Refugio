package session_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"mail/internal/microservice/session/mock"
	"mail/internal/pkg/session"

	session_proto "mail/internal/microservice/session/proto"
)

func TestSessionsManager_SetSession_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionServiceClient := mock.NewMockSessionServiceClient(ctrl)

	sm := session.NewSessionsManager(mockSessionServiceClient)

	mockSessionServiceClient.EXPECT().GetSession(gomock.Any(), gomock.Any()).
		Return(&session_proto.GetSessionReply{
			Session: &session_proto.Session{
				SessionId: "123",
				CsrfToken: "csrfToken",
			},
		}, nil)

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	err = sm.SetSession("123", rr, req, context.WithValue(context.Background(), "requestID", "testID"))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if rr.Header().Get("X-Csrf-Token") != "csrfToken" {
		t.Errorf("Unexpected CSRF token in response header: got %s, want csrfToken", rr.Header().Get("X-Csrf-Token"))
	}

	cookie := rr.Result().Cookies()[0]
	if cookie.Name != "session_id" {
		t.Errorf("Unexpected cookie name: got %s, want session_id", cookie.Name)
	}
	if cookie.Value != "123" {
		t.Errorf("Unexpected cookie value: got %s, want 123", cookie.Value)
	}
}

func TestSessionsManager_SetSession_SessionServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionServiceClient := mock.NewMockSessionServiceClient(ctrl)

	sm := session.NewSessionsManager(mockSessionServiceClient)

	mockSessionServiceClient.EXPECT().GetSession(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("session not found"))

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	err = sm.SetSession("123", rr, req, context.WithValue(context.Background(), "requestID", "testID"))

	assert.Error(t, err)
}

func TestSessionsManager_GetSession_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionServiceClient := mock.NewMockSessionServiceClient(ctrl)

	sm := session.NewSessionsManager(mockSessionServiceClient)

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "123"})

	mockSessionServiceClient.EXPECT().GetSession(gomock.Any(), gomock.Any()).
		Return(&session_proto.GetSessionReply{
			Session: &session_proto.Session{
				SessionId: "123",
				CsrfToken: "csrfToken",
			},
		}, nil)

	sessionTest := sm.GetSession(req, context.WithValue(context.Background(), "requestID", "testID"))

	if sessionTest == nil {
		t.Error("Expected session, got nil")
	}
	if sessionTest != nil && sessionTest.ID != "123" {
		t.Errorf("Unexpected session ID: got %s, want 123", sessionTest.ID)
	}
	if sessionTest != nil && sessionTest.CsrfToken != "csrfToken" {
		t.Errorf("Unexpected CSRF token: got %s, want csrfToken", sessionTest.CsrfToken)
	}
}

func TestSessionsManager_GetSession_SessionNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionServiceClient := mock.NewMockSessionServiceClient(ctrl)

	sm := session.NewSessionsManager(mockSessionServiceClient)

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "123"})

	mockSessionServiceClient.EXPECT().GetSession(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("session not found"))

	sessionTest := sm.GetSession(req, context.WithValue(context.Background(), "requestID", "testID"))

	if sessionTest != nil {
		t.Error("Expected nil session, got non-nil session")
	}
}

func TestSessionsManager_Check_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionServiceClient := mock.NewMockSessionServiceClient(ctrl)

	sm := session.NewSessionsManager(mockSessionServiceClient)

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "123"})
	req.Header.Set("X-Csrf-Token", "csrfToken")

	mockSessionServiceClient.EXPECT().GetSession(gomock.Any(), gomock.Any()).
		Return(&session_proto.GetSessionReply{
			Session: &session_proto.Session{
				SessionId: "123",
				CsrfToken: "csrfToken",
			},
		}, nil)

	sessionTest, err := sm.Check(req, context.WithValue(context.Background(), "requestID", "testID"))

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if sessionTest == nil {
		t.Error("Expected session, got nil")
	}
	if sessionTest != nil && sessionTest.ID != "123" {
		t.Errorf("Unexpected session ID: got %s, want 123", sessionTest.ID)
	}
}

func TestSessionsManager_Check_NoCsrf(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionServiceClient := mock.NewMockSessionServiceClient(ctrl)

	sm := session.NewSessionsManager(mockSessionServiceClient)

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	sessionTest, err := sm.Check(req, context.WithValue(context.Background(), "requestID", "testID"))

	assert.Error(t, err)
	assert.Nil(t, sessionTest)
}

func TestSessionsManager_Check_NoSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionServiceClient := mock.NewMockSessionServiceClient(ctrl)

	sm := session.NewSessionsManager(mockSessionServiceClient)

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-Csrf-Token", "csrfToken")

	sessionTest, err := sm.Check(req, context.WithValue(context.Background(), "requestID", "testID"))

	assert.Error(t, err)
	assert.Nil(t, sessionTest)
}

func TestSessionsManager_Check_ErrorNotSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionServiceClient := mock.NewMockSessionServiceClient(ctrl)

	sm := session.NewSessionsManager(mockSessionServiceClient)

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "123"})
	req.Header.Set("X-Csrf-Token", "csrfToken")

	mockSessionServiceClient.EXPECT().
		GetSession(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("session not found"))

	sessionTest, err := sm.Check(req, context.WithValue(context.Background(), "requestID", "testID"))

	assert.Error(t, err)
	assert.Nil(t, sessionTest)
}

func TestSessionsManager_Check_ErrorCsrf(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionServiceClient := mock.NewMockSessionServiceClient(ctrl)

	sm := session.NewSessionsManager(mockSessionServiceClient)

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "123"})
	req.Header.Set("X-Csrf-Token", "csrfTokenFail")

	mockSessionServiceClient.EXPECT().GetSession(gomock.Any(), gomock.Any()).
		Return(&session_proto.GetSessionReply{
			Session: &session_proto.Session{
				SessionId: "123",
				CsrfToken: "csrfToken",
			},
		}, nil)

	sessionTest, err := sm.Check(req, context.WithValue(context.Background(), "requestID", "testID"))

	assert.Error(t, err)
	assert.Nil(t, sessionTest)
}

func TestSessionsManager_CheckLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionServiceClient := mock.NewMockSessionServiceClient(ctrl)

	sm := session.NewSessionsManager(mockSessionServiceClient)

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "123"})

	mockSessionServiceClient.EXPECT().GetLoginBySession(gomock.Any(), gomock.Any()).
		Return(&session_proto.GetLoginBySessionReply{Login: "user@mailhub.su"}, nil)

	err = sm.CheckLogin("user@mailhub.su", req, context.WithValue(context.Background(), "requestID", "testID"))

	assert.NoError(t, err)
}

func TestSessionsManager_CheckLogin_WrongLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionServiceClient := mock.NewMockSessionServiceClient(ctrl)

	sm := session.NewSessionsManager(mockSessionServiceClient)

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "123"})

	mockSessionServiceClient.EXPECT().GetLoginBySession(gomock.Any(), gomock.Any()).
		Return(&session_proto.GetLoginBySessionReply{Login: "anotheruser@mailhub.su"}, nil)

	err = sm.CheckLogin("user@mailhub.su", req, context.WithValue(context.Background(), "requestID", "testID"))

	assert.Error(t, err)
}

func TestSessionsManager_GetLoginBySession_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionServiceClient := mock.NewMockSessionServiceClient(ctrl)

	sm := session.NewSessionsManager(mockSessionServiceClient)

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "123"})

	mockSessionServiceClient.EXPECT().
		GetLoginBySession(gomock.Any(), gomock.Any()).
		Return(&session_proto.GetLoginBySessionReply{Login: "user@mailhub.su"}, nil)

	login, err := sm.GetLoginBySession(req, context.WithValue(context.Background(), "requestID", "testID"))

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if login != "user@mailhub.su" {
		t.Errorf("Unexpected login: got %s, want user@mailhub.su", login)
	}
}

func TestSessionsManager_GetLoginBySession_NoSessionCookie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionServiceClient := mock.NewMockSessionServiceClient(ctrl)

	sm := session.NewSessionsManager(mockSessionServiceClient)

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "123"})

	mockSessionServiceClient.EXPECT().GetLoginBySession(gomock.Any(), gomock.Any()).
		Return(&session_proto.GetLoginBySessionReply{}, fmt.Errorf("no session found"))

	login, err := sm.GetLoginBySession(req, context.WithValue(context.Background(), "requestID", "testID"))

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if login != "" {
		t.Errorf("Unexpected login: got %s, want empty string", login)
	}
}

func TestSessionsManager_GetProfileIDBySessionID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionServiceClient := mock.NewMockSessionServiceClient(ctrl)

	sm := session.NewSessionsManager(mockSessionServiceClient)

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "123"})

	mockSessionServiceClient.EXPECT().GetProfileIDBySession(gomock.Any(), gomock.Any()).
		Return(&session_proto.GetProfileIDBySessionReply{Id: 123}, nil)

	profileID, err := sm.GetProfileIDBySessionID(req, context.WithValue(context.Background(), "requestID", "testID"))

	assert.NoError(t, err)
	assert.Equal(t, uint32(123), profileID)
}

func TestSessionsManager_GetProfileIDBySessionID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionServiceClient := mock.NewMockSessionServiceClient(ctrl)

	sm := session.NewSessionsManager(mockSessionServiceClient)

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "123"})

	mockSessionServiceClient.EXPECT().GetProfileIDBySession(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("session not found"))

	profileID, err := sm.GetProfileIDBySessionID(req, context.WithValue(context.Background(), "requestID", "testID"))

	assert.Error(t, err)
	assert.Zero(t, profileID)
}

func TestSessionsManager_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionServiceClient := mock.NewMockSessionServiceClient(ctrl)

	sm := session.NewSessionsManager(mockSessionServiceClient)

	mockSessionServiceClient.EXPECT().
		CreateSession(gomock.Any(), gomock.Any()).
		Return(&session_proto.CreateSessionReply{
			SessionId: "123",
		}, nil)

	mockSessionServiceClient.EXPECT().
		GetSession(gomock.Any(), gomock.Any()).
		Return(&session_proto.GetSessionReply{
			Session: &session_proto.Session{
				SessionId: "123",
				CsrfToken: "csrfToken",
			},
		}, nil)

	w := httptest.NewRecorder()

	ctx := context.WithValue(context.Background(), "requestID", "testID")

	sessionTest, err := sm.Create(w, 123, ctx)

	assert.NoError(t, err)

	assert.Equal(t, "csrfToken", w.Header().Get("X-Csrf-Token"))

	cookie := w.Result().Cookies()[0]
	assert.Equal(t, "session_id", cookie.Name)
	assert.Equal(t, "123", cookie.Value)
	assert.WithinDuration(t, time.Now().Add(24*time.Hour), cookie.Expires, 2*time.Second)

	assert.Equal(t, "123", sessionTest.ID)
	assert.Equal(t, "csrfToken", sessionTest.CsrfToken)
}

func TestSessionsManager_Create_SessionCreateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionServiceClient := mock.NewMockSessionServiceClient(ctrl)

	sm := session.NewSessionsManager(mockSessionServiceClient)

	mockSessionServiceClient.EXPECT().CreateSession(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("session already exists"))

	w := httptest.NewRecorder()

	ctx := context.WithValue(context.Background(), "requestID", "testID")

	sessionTest, err := sm.Create(w, 123, ctx)

	assert.Error(t, err)
	assert.Nil(t, sessionTest)
}

func TestSessionsManager_Create_SessionGetError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionServiceClient := mock.NewMockSessionServiceClient(ctrl)

	sm := session.NewSessionsManager(mockSessionServiceClient)

	mockSessionServiceClient.EXPECT().
		CreateSession(gomock.Any(), gomock.Any()).
		Return(&session_proto.CreateSessionReply{
			SessionId: "123",
		}, nil)

	mockSessionServiceClient.EXPECT().
		GetSession(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("session not found"))

	w := httptest.NewRecorder()

	ctx := context.WithValue(context.Background(), "requestID", "testID")

	sessionTest, err := sm.Create(w, 123, ctx)

	assert.Error(t, err)
	assert.Nil(t, sessionTest)
}

func TestSessionsManager_DestroyCurrent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionServiceClient := mock.NewMockSessionServiceClient(ctrl)

	sm := session.NewSessionsManager(mockSessionServiceClient)

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "123"})

	w := httptest.NewRecorder()

	mockSessionServiceClient.EXPECT().
		DeleteSession(gomock.Any(), gomock.Any()).
		Return(&session_proto.DeleteSessionReply{Status: true}, nil)

	err := sm.DestroyCurrent(w, req, context.WithValue(context.Background(), "requestID", "testID"))

	assert.NoError(t, err)

	cookie := w.Result().Cookies()[0]
	assert.Equal(t, "session_id", cookie.Name)
	assert.Equal(t, "", cookie.Value)
	assert.WithinDuration(t, time.Now().AddDate(0, 0, -1), cookie.Expires, 2*time.Second)
}

func TestSessionsManager_DestroyCurrent_NoSessionCookie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionServiceClient := mock.NewMockSessionServiceClient(ctrl)

	sm := session.NewSessionsManager(mockSessionServiceClient)

	req := httptest.NewRequest("GET", "/", nil)

	w := httptest.NewRecorder()

	err := sm.DestroyCurrent(w, req, context.WithValue(context.Background(), "requestID", "testID"))

	assert.Error(t, err)
	assert.EqualError(t, err, "http: named cookie not present")
}
