package session

import (
	domain "mail/pkg/domain/models"
	repository "mail/pkg/domain/repository"
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
func (uc *SessionUseCase) CreateNewSession(userID uint32, device string, lifeTime int) (string, error) {
	return uc.sessionRepo.CreateSession(userID, device, lifeTime)
}

// GetSession fetches a session by its unique identifier.
func (uc *SessionUseCase) GetSession(sessionID string) (*domain.Session, error) {
	return uc.sessionRepo.GetSessionByID(sessionID)
}

func (uc *SessionUseCase) GetLogin(sessionID string) (string, error) {
	return uc.sessionRepo.GetLoginBySessionID(sessionID)
}

// DeleteSession terminates a session identified by its ID.
func (uc *SessionUseCase) DeleteSession(sessionID string) error {
	return uc.sessionRepo.DeleteSessionByID(sessionID)
}

// CleanupExpiredSessions removes sessions that have exceeded their lifetime.
func (uc *SessionUseCase) CleanupExpiredSessions() error {
	return uc.sessionRepo.DeleteExpiredSessions()
}
