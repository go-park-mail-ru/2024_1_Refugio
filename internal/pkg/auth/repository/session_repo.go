package repository

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/jmoiron/sqlx"

	"mail/internal/logger"
	converters "mail/internal/models/repository_converters"

	domain "mail/internal/models/domain_models"
	database "mail/internal/models/repository_models"
)

// SessionRepository represents a PostgreSQL implementation of the SessionRepository interface.
type SessionRepository struct {
	DB *sqlx.DB
}

// Logger represents the logger used for logging database initialization.
var sessionLogger = logger.InitializationEmptyLog()

// NewSessionRepository creates a new instance of SessionRepository.
func NewSessionRepository(db *sqlx.DB) *SessionRepository {
	f, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create logfile in session_repo" + "log.txt")
	}
	Logger = logger.InitializationBdLog(f)
	return &SessionRepository{
		DB: db,
	}
}

// SessionGenerate is a function type used for generating session IDs.
type SessionGenerate func() string

// SessionGenerateRandomID generates a random session ID using cryptographic random numbers.
var SessionGenerateRandomID SessionGenerate = func() string {
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

	ID := SessionGenerateRandomID()
	csrfToken := SessionGenerateRandomID()
	creationDate := time.Now()

	args := []interface{}{ID, userID, device, creationDate, lifeTime, csrfToken}
	start := time.Now()

	_, err := repo.DB.Exec(query, ID, userID, device, creationDate, lifeTime, csrfToken)
	defer sessionLogger.DbLog(query, requestID, start, &err, args)
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}

	return ID, nil
}

// GetSessionByID retrieves a session by its ID.
func (repo *SessionRepository) GetSessionByID(sessionID, requestID string) (*domain.Session, error) {
	query := `SELECT * FROM session WHERE id = $1`

	args := []interface{}{sessionID}
	start := time.Now()

	var session database.Session
	err := repo.DB.Get(&session, query, sessionID)
	defer sessionLogger.DbLog(query, requestID, start, &err, args)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %v", err)
	}

	return converters.SessionConvertDbInCore(session), nil
}

// GetLoginBySessionID retrieves the login associated with the given session ID.
func (repo *SessionRepository) GetLoginBySessionID(sessionID, requestID string) (string, error) {
	query := `
		SELECT login FROM profile
		JOIN session ON session.profile_id = profile.id 
		WHERE session.id = $1
	`

	args := []interface{}{sessionID}
	start := time.Now()

	var login string
	err := repo.DB.Get(&login, query, sessionID)
	defer sessionLogger.DbLog(query, requestID, start, &err, args)
	if err != nil {
		return "", fmt.Errorf("failed to get session: %v", err)
	}

	return login, nil
}

// DeleteSessionByID deletes a session by its ID.
func (repo *SessionRepository) DeleteSessionByID(sessionID, requestID string) error {
	query := "DELETE FROM session WHERE id = $1"

	args := []interface{}{sessionID}
	start := time.Now()

	_, err := repo.DB.Exec(query, sessionID)
	defer sessionLogger.DbLog(query, requestID, start, &err, args)
	if err != nil {
		return fmt.Errorf("failed to delete session: %v", err)
	}

	return nil
}

// DeleteExpiredSessions removes all expired sessions.
func (repo *SessionRepository) DeleteExpiredSessions() error {
	query := "DELETE FROM session WHERE creation_date + life_time * interval '1 second' < now()"

	args := []interface{}{}
	start := time.Now()

	_, err := repo.DB.Exec(query)
	defer sessionLogger.DbLog(query, "DeleteExpiredSessionsNULL", start, &err, args)
	if err != nil {
		return fmt.Errorf("failed to delete expired sessions: %v", err)
	}

	return nil
}
