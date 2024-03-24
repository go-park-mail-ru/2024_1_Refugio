//go:generate mockgen -source=./isession_repo.go -destination=../mock/session_repository_mock.go -package=mock

package repository

import (
	domain "mail/pkg/domain/models"
)

// SessionRepository represents the interface for managing user sessions.
type SessionRepository interface {
	// CreateSession creates a new session and returns its ID.
	CreateSession(ID uint32, userID uint32, device string, lifeTime int) (uint32, error)

	// GetSessionByID retrieves a session by its ID.
	GetSessionByID(sessionID uint32) (*domain.Session, error)

	// DeleteSessionByID deletes a session by its ID.
	DeleteSessionByID(sessionID uint32) error

	// DeleteExpiredSessions removes all expired sessions.
	DeleteExpiredSessions() error
}
