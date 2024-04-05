//go:generate mockgen -source=./isession_service.go -destination=../mock/session_service_mock.go -package=mock

package usecase

import domain "mail/pkg/domain/models"

// SessionUseCase represents the interface for session-related operations.
type SessionUseCase interface {
	// CreateNewSession initiates a new session for a user.
	CreateNewSession(userID uint32, device string, lifeTime int) (string, error)

	// GetSession fetches a session by its unique identifier.
	GetSession(sessionID string) (*domain.Session, error)

	GetLogin(sessionID string) (string, error)

	// DeleteSession terminates a session identified by its ID.
	DeleteSession(sessionID string) error

	// CleanupExpiredSessions removes sessions that have exceeded their lifetime.
	CleanupExpiredSessions() error
}
