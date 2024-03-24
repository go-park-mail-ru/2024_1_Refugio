package session

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	domain "mail/pkg/domain/models"
	"mail/pkg/repository/converters"
	database "mail/pkg/repository/models"
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

// CreateSession creates a new session and returns its ID.
func (repo *SessionRepository) CreateSession(ID uint32, userID uint32, device string, lifeTime int) (uint32, error) {
	query := `
		INSERT INTO sessions (id, user_id, device, creation_date, lifetime)
		VALUES ($1, $2, $3, $4, $5)
	`

	creationDate := time.Now()

	_, err := repo.DB.Exec(query, ID, userID, device, creationDate, lifeTime)
	if err != nil {
		return 0, fmt.Errorf("failed to create session: %v", err)
	}

	return ID, nil
}

// GetSessionByID retrieves a session by its ID.
func (repo *SessionRepository) GetSessionByID(sessionID uint32) (*domain.Session, error) {
	query := `SELECT * FROM sessions WHERE id = $1`

	var session database.Session
	err := repo.DB.Get(&session, query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %v", err)
	}

	return converters.SessionConvertDbInCore(session), nil
}

// DeleteSessionByID deletes a session by its ID.
func (repo *SessionRepository) DeleteSessionByID(sessionID uint32) error {
	query := "DELETE FROM sessions WHERE id = $1"

	_, err := repo.DB.Exec(query, sessionID)
	if err != nil {
		return fmt.Errorf("failed to delete session: %v", err)
	}

	return nil
}

// DeleteExpiredSessions removes all expired sessions.
func (repo *SessionRepository) DeleteExpiredSessions() error {
	query := "DELETE FROM sessions WHERE creation_date + lifetime * interval '1 second' < now()"

	_, err := repo.DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to delete expired sessions: %v", err)
	}

	return nil
}
