package repository

import (
	"context"
	"fmt"
	domain "mail/internal/microservice/models/domain_models"
	converters "mail/internal/microservice/models/repository_converters"
	database "mail/internal/microservice/models/repository_models"
	"mail/internal/pkg/logger"
	"math/rand"
	"time"

	"github.com/jmoiron/sqlx"
)

var requestIDContextKey interface{} = "requestID"

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

// SessionGenerate is a function type used for generating session IDs.
type SessionGenerate func() string

// SessionGenerateRandomID generates a random session ID using cryptographic random numbers.
var SessionGenerateRandomID SessionGenerate = func() string {
	randID := make([]byte, 16)
	rand.Read(randID)

	return fmt.Sprintf("%x", randID)
}

// CreateSession creates a new session and returns its ID.
func (repo *SessionRepository) CreateSession(userID uint32, device string, lifeTime int, ctx context.Context) (string, error) {
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
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, ctx.Value(requestIDContextKey).(string), start, &err, args)
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}

	return ID, nil
}

// GetSessionByID retrieves a session by its ID.
func (repo *SessionRepository) GetSessionByID(sessionID string, ctx context.Context) (*domain.Session, error) {
	query := `SELECT * FROM session WHERE id = $1`

	args := []interface{}{sessionID}
	start := time.Now()

	var session database.Session
	err := repo.DB.Get(&session, query, sessionID)
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, ctx.Value(requestIDContextKey).(string), start, &err, args)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %v", err)
	}

	return converters.SessionConvertDbInCore(session), nil
}

// GetLoginBySessionID retrieves the login associated with the given session ID.
func (repo *SessionRepository) GetLoginBySessionID(sessionID string, ctx context.Context) (string, error) {
	query := `
		SELECT login FROM profile
		JOIN session ON session.profile_id = profile.id 
		WHERE session.id = $1
	`

	args := []interface{}{sessionID}
	start := time.Now()

	var login string
	err := repo.DB.Get(&login, query, sessionID)
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, ctx.Value(requestIDContextKey).(string), start, &err, args)
	if err != nil {
		return "", fmt.Errorf("failed to get session: %v", err)
	}

	return login, nil
}

// DeleteSessionByID deletes a session by its ID.
func (repo *SessionRepository) DeleteSessionByID(sessionID string, ctx context.Context) error {
	query := "DELETE FROM session WHERE id = $1"

	args := []interface{}{sessionID}
	start := time.Now()

	_, err := repo.DB.Exec(query, sessionID)
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, ctx.Value(requestIDContextKey).(string), start, &err, args)
	if err != nil {
		return fmt.Errorf("failed to delete session: %v", err)
	}

	return nil
}

// DeleteExpiredSessions removes all expired sessions.
func (repo *SessionRepository) DeleteExpiredSessions(ctx context.Context) error {
	query := "DELETE FROM session WHERE creation_date + life_time * interval '1 second' < now()"

	args := []interface{}{}
	start := time.Now()

	_, err := repo.DB.Exec(query)
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, ctx.Value(requestIDContextKey).(string), start, &err, args)
	if err != nil {
		return fmt.Errorf("failed to delete expired sessions: %v", err)
	}

	return nil
}
