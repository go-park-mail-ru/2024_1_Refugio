package usecase

import domain "mail/pkg/domain/models"

// SessionUseCase represents the interface for session-related operations.
type SessionUseCase interface {
	// CreateNewSession initiates a new session for a user.
	CreateNewSession(ID uint32, userID uint32, device string, lifeTime int) (uint32, error)

	// GetSession fetches a session by its unique identifier.
	GetSession(sessionID uint32) (*domain.Session, error)

	// DeleteSession terminates a session identified by its ID.
	DeleteSession(sessionID uint32) error

	// CleanupExpiredSessions removes sessions that have exceeded their lifetime.
	CleanupExpiredSessions() error
}
