package session

import (
	"context"
	"fmt"
	converters "mail/internal/models/delivery_converters"
	api "mail/internal/models/delivery_models"
	domain "mail/internal/pkg/session/interface"
	"net/http"
	"time"
)

// GlobalSessionManager is a global instance of SessionsManager.
var (
	GlobalSessionManager = &SessionsManager{}
)

// SessionsManager manages user sessions.
type SessionsManager struct {
	sessionUseCase domain.SessionUseCase
}

// InitializationGlobalSeaaionManager initializes the global session manager.
func InitializationGlobalSeaaionManager(sessionManager *SessionsManager) {
	GlobalSessionManager = sessionManager
}

// NewSessionsManager creates a new instance of SessionsManager.
func NewSessionsManager(sessionUc domain.SessionUseCase) *SessionsManager {
	return &SessionsManager{
		sessionUseCase: sessionUc,
	}
}

// GetSession retrieves the session from the request and sanitizes it.
func (sm *SessionsManager) GetSession(r *http.Request, ctx context.Context) *api.Session {
	sessionCookie, _ := r.Cookie("session_id")

	sess, _ := sm.sessionUseCase.GetSession(sessionCookie.Value, ctx)
	if sess == nil {
		return nil
	}

	return converters.SessionConvertCoreInApi(*sess)
}

// Check checks the validity of the session and CSRF token in the request.
func (sm *SessionsManager) Check(r *http.Request, ctx context.Context) (*api.Session, error) {
	csrfToken := r.Header.Get("X-Csrf-Token")
	if r.URL.Path != "/api/v1/verify-auth" && csrfToken == "" {
		return nil, fmt.Errorf("CSRF token not found in request headers")
	}

	sessionCookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		return nil, fmt.Errorf("no session found")
	}

	sess, ok := sm.sessionUseCase.GetSession(sessionCookie.Value, ctx)
	if ok != nil {
		return nil, fmt.Errorf("no session found")
	}
	if r.URL.Path != "/api/v1/verify-auth" && sess.CsrfToken != csrfToken {
		return nil, fmt.Errorf("CSRF token mismatch")
	}

	return converters.SessionConvertCoreInApi(*sess), nil
}

// CheckLogin checks if the login associated with the session matches the provided login.
func (sm *SessionsManager) CheckLogin(login string, r *http.Request, ctx context.Context) error {
	sessionCookie, _ := r.Cookie("session_id")
	LoginBd, err := sm.sessionUseCase.GetLogin(sessionCookie.Value, ctx)
	if err != nil {
		return err
	}
	if LoginBd != login {
		return fmt.Errorf("No right sender email")
	}

	return nil
}

// GetLoginBySession retrieves the login associated with the session from the request.
func (sm SessionsManager) GetLoginBySession(r *http.Request, ctx context.Context) (string, error) {
	sessionCookie, _ := r.Cookie("session_id")
	Login, err := sm.sessionUseCase.GetLogin(sessionCookie.Value, ctx)
	if err != nil {
		return "", err
	}

	return Login, nil
}

// Create creates a new session for the user and sets the session ID cookie in the response.
func (sm *SessionsManager) Create(w http.ResponseWriter, userID uint32, ctx context.Context) (*api.Session, error) {
	sessionID, err := sm.sessionUseCase.CreateNewSession(userID, "", 60*60*24, ctx)

	if err != nil {
		return nil, fmt.Errorf("session already exist")
	}

	sess, _ := sm.sessionUseCase.GetSession(sessionID, ctx)

	w.Header().Set("X-Csrf-Token", sess.CsrfToken)

	sessionCookie := &http.Cookie{
		Name:     "session_id",
		Value:    sess.ID,
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
		HttpOnly: true,
		// SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(w, sessionCookie)

	return converters.SessionConvertCoreInApi(*sess), nil
}

// DestroyCurrent destroys the current session by deleting the session ID cookie from the response.
func (sm *SessionsManager) DestroyCurrent(w http.ResponseWriter, r *http.Request, ctx context.Context) error {
	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		return err
	}

	ok := sm.sessionUseCase.DeleteSession(sessionCookie.Value, ctx)
	if ok != nil {
		return fmt.Errorf("no session found")
	}

	sessionCookieToDelete := http.Cookie{
		Name:    "session_id",
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    "/",
	}
	http.SetCookie(w, &sessionCookieToDelete)

	return nil
}
