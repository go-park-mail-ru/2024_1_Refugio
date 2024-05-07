//go:generate mockgen -source=./isession_service.go -destination=../mock/session_service_mock.go -package=mock

package _interface

import (
	"context"

	domain "mail/internal/microservice/models/domain_models"
)

// SessionUseCase represents the interface for session-related operations.
type SessionUseCase interface {
	// CreateNewSession initiates a new session for a user.
	CreateNewSession(userID uint32, device string, lifeTime int, ctx context.Context) (string, error)

	// GetSession fetches a session by its unique identifier.
	GetSession(sessionID string, ctx context.Context) (*domain.Session, error)

	// GetLogin retrieves the login associated with the provided session ID.
	GetLogin(sessionID string, ctx context.Context) (string, error)

	// GetProfileID retrieves the profile id associated with the given session ID.
	GetProfileID(sessionID string, ctx context.Context) (uint32, error)

	// DeleteSession terminates a session identified by its ID.
	DeleteSession(sessionID string, ctx context.Context) error

	// CleanupExpiredSessions removes sessions that have exceeded their lifetime.
	CleanupExpiredSessions(ctx context.Context) error
}
