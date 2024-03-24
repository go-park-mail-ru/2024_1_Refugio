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
func (uc *SessionUseCase) CreateNewSession(ID uint32, userID uint32, device string, lifeTime int) (uint32, error) {
	return uc.sessionRepo.CreateSession(ID, userID, device, lifeTime)
}

// GetSession fetches a session by its unique identifier.
func (uc *SessionUseCase) GetSession(sessionID uint32) (*domain.Session, error) {
	return uc.sessionRepo.GetSessionByID(sessionID)
}

// DeleteSession terminates a session identified by its ID.
func (uc *SessionUseCase) DeleteSession(sessionID uint32) error {
	return uc.sessionRepo.DeleteSessionByID(sessionID)
}

// CleanupExpiredSessions removes sessions that have exceeded their lifetime.
func (uc *SessionUseCase) CleanupExpiredSessions() error {
	return uc.sessionRepo.DeleteExpiredSessions()
}
