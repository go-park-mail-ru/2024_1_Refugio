package session

import (
	"encoding/binary"
	"fmt"
	"mail/pkg/delivery/converters"
	"mail/pkg/delivery/models"
	domain "mail/pkg/domain/usecase"
	"math/rand"
	"net/http"
	"strconv"
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

	sessionID, err := strconv.Atoi(sessionCookie.Value)
	if err != nil {
		return nil
	}
	sess, _ := sm.sessionUseCase.GetSession(uint32(sessionID))

	return converters.SessionConvertCoreInApi(*sess)
}

func (sm *SessionsManager) Check(r *http.Request) (*models.Session, error) {
	sessionCookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		return nil, fmt.Errorf("no session found")
	}

	sessionID, err := strconv.ParseUint(sessionCookie.Value, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("no session found")
	}

	sess, ok := sm.sessionUseCase.GetSession(uint32(sessionID))
	if ok != nil {
		return nil, fmt.Errorf("no session found")
	}

	return converters.SessionConvertCoreInApi(*sess), nil
}

func (sm *SessionsManager) Create(w http.ResponseWriter, userID uint32) (*models.Session, error) {
	sessionID, _ := sm.sessionUseCase.CreateNewSession(GenerateRandomID(userID), userID, "", 60*60*24*7)

	if sessionID == 0 {
		return nil, fmt.Errorf("session already exist")
	}

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   strconv.FormatUint(uint64(sessionID), 10),
		Expires: time.Now().Add(90 * 24 * time.Hour),
		Path:    "/",
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

	sessionID, err := strconv.ParseUint(sessionCookie.Value, 10, 32)
	if err != nil {
		return fmt.Errorf("no session found")
	}

	ok := sm.sessionUseCase.DeleteSession(uint32(sessionID))
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

func GenerateRandomID(userID uint32) uint32 {
	rand.Seed(time.Now().UnixNano())

	randBytes := make([]byte, 4)
	rand.Read(randBytes)

	randID := binary.BigEndian.Uint32(randBytes)

	return randID
}
