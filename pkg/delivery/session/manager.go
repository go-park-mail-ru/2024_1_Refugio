package session

import (
	"fmt"
	"mail/pkg/delivery/converters"
	"mail/pkg/delivery/models"
	domain "mail/pkg/domain/usecase"
	"net/http"
	"time"
)

var (
	GlobalSeaaionManager = &SessionsManager{}
)

type SessionsManager struct {
	sessionUseCase domain.SessionUseCase
}

func InitializationGlobalSeaaionManager(sessionManager *SessionsManager) {
	GlobalSeaaionManager = sessionManager
}

func NewSessionsManager(sessionUc domain.SessionUseCase) *SessionsManager {
	return &SessionsManager{
		sessionUseCase: sessionUc,
	}
}

func (sm *SessionsManager) GetSession(r *http.Request, requestID string) *models.Session {
	sessionCookie, _ := r.Cookie("session_id")

	sess, _ := sm.sessionUseCase.GetSession(sessionCookie.Value, requestID)
	if sess == nil {
		return nil
	}

	return converters.SessionConvertCoreInApi(*sess)
}

func (sm *SessionsManager) Check(r *http.Request, requestID string) (*models.Session, error) {
	csrfToken := r.Header.Get("X-CSRF-Token")
	if csrfToken == "" {
		return nil, fmt.Errorf("CSRF token not found in request headers")
	}

	sessionCookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		return nil, fmt.Errorf("no session found")
	}

	sess, ok := sm.sessionUseCase.GetSession(sessionCookie.Value, requestID)
	if ok != nil {
		return nil, fmt.Errorf("no session found")
	}
	if sess.CsrfToken != csrfToken {
		return nil, fmt.Errorf("CSRF token mismatch")
	}

	return converters.SessionConvertCoreInApi(*sess), nil
}

func (sm *SessionsManager) ChekLogin(login, requestID string, r *http.Request) error {
	sessionCookie, _ := r.Cookie("session_id")
	LoginBd, _ := sm.sessionUseCase.GetLogin(sessionCookie.Value, requestID)
	if LoginBd != login {
		return fmt.Errorf("No right sender email")
	}

	return nil
}

func (sm *SessionsManager) Create(w http.ResponseWriter, userID uint32, requestID string) (*models.Session, error) {
	sessionID, err := sm.sessionUseCase.CreateNewSession(userID, "", requestID, 60*60*24)

	if err != nil {
		return nil, fmt.Errorf("session already exist")
	}

	sess, _ := sm.sessionUseCase.GetSession(sessionID, requestID)

	csrfCookie := &http.Cookie{
		Name:     "csrf_token",
		Value:    sess.CsrfToken,
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, csrfCookie)
	fmt.Println("CSRFTOKEN: ", csrfCookie.Value)
	sessionCookie := &http.Cookie{
		Name:     "session_id",
		Value:    sess.ID,
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, sessionCookie)

	return converters.SessionConvertCoreInApi(*sess), nil
}

func (sm *SessionsManager) DestroyCurrent(w http.ResponseWriter, r *http.Request, requestID string) error {
	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		return err
	}

	ok := sm.sessionUseCase.DeleteSession(sessionCookie.Value, requestID)
	if ok != nil {
		return fmt.Errorf("no session found")
	}

	sessionCookieToDelete := http.Cookie{
		Name:    "session_id",
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    "/",
	}
	http.SetCookie(w, &sessionCookieToDelete)

	csrfCookieToDelete := http.Cookie{
		Name:    "csrf_token",
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    "/",
	}
	http.SetCookie(w, &csrfCookieToDelete)

	return nil
}
