//go:generate mockgen -source=./isession_repo.go -destination=../mock/session_repository_mock.go -package=mock

package _interface

import (
	domain "mail/internal/models/domain_models"
)

// SessionRepository represents the interface for managing user sessions.
type SessionRepository interface {
	// CreateSession creates a new session and returns its ID.
	CreateSession(userID uint32, device, requestID string, lifeTime int) (string, error)

	// GetSessionByID retrieves a session by its ID.
	GetSessionByID(sessionID, requestID string) (*domain.Session, error)

	// GetLoginBySessionID retrieves the login associated with the given session ID.
	GetLoginBySessionID(sessionID, requestID string) (string, error)

	// DeleteSessionByID deletes a session by its ID.
	DeleteSessionByID(sessionID, requestID string) error

	// DeleteExpiredSessions removes all expired sessions.
	DeleteExpiredSessions() error
}
