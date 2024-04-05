package session

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	domain "mail/pkg/domain/models"
	"mail/pkg/repository/converters"
	database "mail/pkg/repository/models"
	"math/rand"
	"time"
)

// SessionRepository represents a PostgreSQL implementation of the SessionRepository interface.
type SessionRepository struct {
	DB *sqlx.DB
}

// NewSessionRepository creates a new instance of SessionRepository.
func NewSessionRepository(db *sqlx.DB) *SessionRepository {
	return &SessionRepository{
		DB: db,
	}
}

type SessionGenerate func() string

var GenerateRandomID SessionGenerate = func() string {
	randID := make([]byte, 16)
	rand.Read(randID)

	return fmt.Sprintf("%x", randID)
}

// CreateSession creates a new session and returns its ID.
func (repo *SessionRepository) CreateSession(userID uint32, device string, lifeTime int) (string, error) {
	query := `
		INSERT INTO session (id, profile_id, device, creation_date, life_time, csrf_token)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	creationDate := time.Now()
	ID := GenerateRandomID()
	csrfToken := GenerateRandomID()

	_, err := repo.DB.Exec(query, ID, userID, device, creationDate, lifeTime, csrfToken)
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}

	return ID, nil
}

// GetSessionByID retrieves a session by its ID.
func (repo *SessionRepository) GetSessionByID(sessionID string) (*domain.Session, error) {
	query := `SELECT * FROM session WHERE id = $1`

	var session database.Session
	err := repo.DB.Get(&session, query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %v", err)
	}

	return converters.SessionConvertDbInCore(session), nil
}

func (repo *SessionRepository) GetLoginBySessionID(sessionID string) (string, error) {
	query := `
		SELECT login FROM profile
		JOIN session ON session.profile_id = profile.id 
		WHERE session.id = $1
	`

	var login string
	err := repo.DB.Get(&login, query, sessionID)
	if err != nil {
		return "", fmt.Errorf("failed to get session: %v", err)
	}

	return login, nil
}

// DeleteSessionByID deletes a session by its ID.
func (repo *SessionRepository) DeleteSessionByID(sessionID string) error {
	query := "DELETE FROM session WHERE id = $1"

	_, err := repo.DB.Exec(query, sessionID)
	if err != nil {
		return fmt.Errorf("failed to delete session: %v", err)
	}

	return nil
}

// DeleteExpiredSessions removes all expired sessions.
func (repo *SessionRepository) DeleteExpiredSessions() error {
	query := "DELETE FROM session WHERE creation_date + life_time * interval '1 second' < now()"

	_, err := repo.DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to delete expired sessions: %v", err)
	}

	return nil
}
