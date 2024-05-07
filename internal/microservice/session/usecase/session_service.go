package usecase

import (
	"context"

	domain "mail/internal/microservice/models/domain_models"
	repository "mail/internal/microservice/session/interface"
)

// SessionUseCase is a concrete implementation of the SessionUseCase interface.
type SessionUseCase struct {
	sessionRepo repository.SessionRepository
}

// NewSessionUseCase creates a new instance of a session use case with necessary dependencies.
func NewSessionUseCase(repo repository.SessionRepository) *SessionUseCase {
	return &SessionUseCase{
		sessionRepo: repo,
	}
}

// CreateNewSession initiates a new session for a user.
func (uc *SessionUseCase) CreateNewSession(userID uint32, device string, lifeTime int, ctx context.Context) (string, error) {
	return uc.sessionRepo.CreateSession(userID, device, lifeTime, ctx)
}

// GetSession fetches a session by its unique identifier.
func (uc *SessionUseCase) GetSession(sessionID string, ctx context.Context) (*domain.Session, error) {
	return uc.sessionRepo.GetSessionByID(sessionID, ctx)
}

// GetLogin retrieves the login associated with the provided session ID.
func (uc *SessionUseCase) GetLogin(sessionID string, ctx context.Context) (string, error) {
	return uc.sessionRepo.GetLoginBySessionID(sessionID, ctx)
}

// GetProfileID retrieves the login associated with the provided session ID.
func (uc *SessionUseCase) GetProfileID(sessionID string, ctx context.Context) (uint32, error) {
	return uc.sessionRepo.GetProfileIDBySessionID(sessionID, ctx)
}

// DeleteSession terminates a session identified by its ID.
func (uc *SessionUseCase) DeleteSession(sessionID string, ctx context.Context) error {
	return uc.sessionRepo.DeleteSessionByID(sessionID, ctx)
}

// CleanupExpiredSessions removes sessions that have exceeded their lifetime.
func (uc *SessionUseCase) CleanupExpiredSessions(ctx context.Context) error {
	return uc.sessionRepo.DeleteExpiredSessions(ctx)
}
