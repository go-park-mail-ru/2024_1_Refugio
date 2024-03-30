//go:generate mockgen -source=./isession_repo.go -destination=../mock/session_repository_mock.go -package=mock

package repository

import (
	domain "mail/pkg/domain/models"
)

// SessionRepository represents the interface for managing user sessions.
type SessionRepository interface {
	// CreateSession creates a new session and returns its ID.
	CreateSession(userID uint32, device string, lifeTime int) (string, error)

	// GetSessionByID retrieves a session by its ID.
	GetSessionByID(sessionID string) (*domain.Session, error)

	// DeleteSessionByID deletes a session by its ID.
	DeleteSessionByID(sessionID string) error

	// DeleteExpiredSessions removes all expired sessions.
	DeleteExpiredSessions() error
}
