//go:generate mockgen -source=./isession_repo.go -destination=../mock/session_repository_mock.go -package=mock

package _interface

import (
	"context"

	domain "mail/internal/microservice/models/domain_models"
)

// SessionRepository represents the interface for managing user sessions.
type SessionRepository interface {
	// CreateSession creates a new session and returns its ID.
	CreateSession(userID uint32, device string, lifeTime int, ctx context.Context) (string, error)

	// GetSessionByID retrieves a session by its ID.
	GetSessionByID(sessionID string, ctx context.Context) (*domain.Session, error)

	// GetLoginBySessionID retrieves the login associated with the given session ID.
	GetLoginBySessionID(sessionID string, ctx context.Context) (string, error)

	// GetProfileIDBySessionID retrieves the profile id associated with the given session ID.
	GetProfileIDBySessionID(sessionID string, ctx context.Context) (uint32, error)

	// DeleteSessionByID deletes a session by its ID.
	DeleteSessionByID(sessionID string, ctx context.Context) error

	// DeleteExpiredSessions removes all expired sessions.
	DeleteExpiredSessions(ctx context.Context) error
}
