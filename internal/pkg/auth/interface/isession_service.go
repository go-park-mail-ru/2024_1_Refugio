//go:generate mockgen -source=./isession_service.go -destination=../mock/session_service_mock.go -package=mock

package _interface

import domain "mail/internal/models/domain_models"

// SessionUseCase represents the interface for session-related operations.
type SessionUseCase interface {
	// CreateNewSession initiates a new session for a user.
	CreateNewSession(userID uint32, device, requestID string, lifeTime int) (string, error)

	// GetSession fetches a session by its unique identifier.
	GetSession(sessionID, requestID string) (*domain.Session, error)

	// GetLogin retrieves the login associated with the provided session ID.
	GetLogin(sessionID, requestID string) (string, error)

	// DeleteSession terminates a session identified by its ID.
	DeleteSession(sessionID, requestID string) error

	// CleanupExpiredSessions removes sessions that have exceeded their lifetime.
	CleanupExpiredSessions() error
}
