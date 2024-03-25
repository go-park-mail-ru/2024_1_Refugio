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

func (sm *SessionsManager) GetSession(r *http.Request) *models.Session {
	sessionCookie, _ := r.Cookie("session_id")

	sess, _ := sm.sessionUseCase.GetSession(sessionCookie.Value)

	return converters.SessionConvertCoreInApi(*sess)
}

func (sm *SessionsManager) Check(r *http.Request) (*models.Session, error) {
	sessionCookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		return nil, fmt.Errorf("no session found")
	}

	sess, ok := sm.sessionUseCase.GetSession(sessionCookie.Value)
	if ok != nil {
		return nil, fmt.Errorf("no session found")
	}

	return converters.SessionConvertCoreInApi(*sess), nil
}

func (sm *SessionsManager) Create(w http.ResponseWriter, userID uint32) (*models.Session, error) {
	sessionID, _ := sm.sessionUseCase.CreateNewSession(userID, "", 60*60*24*7)

	if sessionID == "" {
		return nil, fmt.Errorf("session already exist")
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(90 * 24 * time.Hour),
		Path:     "/",
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)

	sess, _ := sm.sessionUseCase.GetSession(sessionID)

	return converters.SessionConvertCoreInApi(*sess), nil
}

func (sm *SessionsManager) DestroyCurrent(w http.ResponseWriter, r *http.Request) error {
	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		return err
	}

	ok := sm.sessionUseCase.DeleteSession(sessionCookie.Value)
	if ok != nil {
		return fmt.Errorf("no session found")
	}

	cookie := http.Cookie{
		Name:    "session_id",
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    "/",
	}

	http.SetCookie(w, &cookie)

	return nil
}
