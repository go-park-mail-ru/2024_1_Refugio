package session

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"mail/pkg/domain/logger"
	domain "mail/pkg/domain/models"
	"mail/pkg/repository/converters"
	database "mail/pkg/repository/models"
	"math/rand"
	"time"
)

var Logger = logger.InitializationBdLog()

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
func (repo *SessionRepository) CreateSession(userID uint32, device, requestID string, lifeTime int) (string, error) {
	query := `
		INSERT INTO session (id, profile_id, device, creation_date, life_time, csrf_token)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	creationDate := time.Now()
	ID := GenerateRandomID()
	csrfToken := GenerateRandomID()
	args := []interface{}{ID, userID, device, creationDate, lifeTime, csrfToken}

	start := time.Now()
	_, err := repo.DB.Exec(query, ID, userID, device, creationDate, lifeTime, csrfToken)
	if err != nil {
		Logger.DbLog(query, requestID, 500, start, err, args)
		return "", fmt.Errorf("failed to create session: %v", err)
	}

	Logger.DbLog(query, requestID, 200, start, nil, args)
	return ID, nil
}

// GetSessionByID retrieves a session by its ID.
func (repo *SessionRepository) GetSessionByID(sessionID, requestID string) (*domain.Session, error) {
	query := `SELECT * FROM session WHERE id = $1`
	args := []interface{}{sessionID}

	var session database.Session
	start := time.Now()
	err := repo.DB.Get(&session, query, sessionID)
	if err != nil {
		Logger.DbLog(query, requestID, 500, start, err, args)
		return nil, fmt.Errorf("failed to get session: %v", err)
	}

	Logger.DbLog(query, requestID, 200, start, nil, args)
	return converters.SessionConvertDbInCore(session), nil
}

func (repo *SessionRepository) GetLoginBySessionID(sessionID, requestID string) (string, error) {
	query := `
		SELECT login FROM profile
		JOIN session ON session.profile_id = profile.id 
		WHERE session.id = $1
	`
	args := []interface{}{sessionID}

	var login string
	start := time.Now()
	err := repo.DB.Get(&login, query, sessionID)
	if err != nil {
		Logger.DbLog(query, requestID, 500, start, err, args)
		return "", fmt.Errorf("failed to get session: %v", err)
	}

	Logger.DbLog(query, requestID, 200, start, nil, args)
	return login, nil
}

// DeleteSessionByID deletes a session by its ID.
func (repo *SessionRepository) DeleteSessionByID(sessionID, requestID string) error {
	query := "DELETE FROM session WHERE id = $1"
	args := []interface{}{sessionID}

	start := time.Now()
	_, err := repo.DB.Exec(query, sessionID)
	if err != nil {
		Logger.DbLog(query, requestID, 500, start, err, args)
		return fmt.Errorf("failed to delete session: %v", err)
	}

	Logger.DbLog(query, requestID, 200, start, nil, args)
	return nil
}

// DeleteExpiredSessions removes all expired sessions.
func (repo *SessionRepository) DeleteExpiredSessions() error {
	query := "DELETE FROM session WHERE creation_date + life_time * interval '1 second' < now()"
	args := []interface{}{}

	start := time.Now()
	_, err := repo.DB.Exec(query)
	if err != nil {
		Logger.DbLog(query, "DeleteExpiredSessionsNULL", 500, start, err, args)
		return fmt.Errorf("failed to delete expired sessions: %v", err)
	}

	Logger.DbLog(query, "DeleteExpiredSessionsNULL", 200, start, nil, args)
	return nil
}
