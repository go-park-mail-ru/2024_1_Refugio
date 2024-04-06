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
func (uc *SessionUseCase) CreateNewSession(userID uint32, device, requestID string, lifeTime int) (string, error) {
	return uc.sessionRepo.CreateSession(userID, device, requestID, lifeTime)
}

// GetSession fetches a session by its unique identifier.
func (uc *SessionUseCase) GetSession(sessionID, requestID string) (*domain.Session, error) {
	return uc.sessionRepo.GetSessionByID(sessionID, requestID)
}

func (uc *SessionUseCase) GetLogin(sessionID, requestID string) (string, error) {
	return uc.sessionRepo.GetLoginBySessionID(sessionID, requestID)
}

// DeleteSession terminates a session identified by its ID.
func (uc *SessionUseCase) DeleteSession(sessionID, requestID string) error {
	return uc.sessionRepo.DeleteSessionByID(sessionID, requestID)
}

// CleanupExpiredSessions removes sessions that have exceeded their lifetime.
func (uc *SessionUseCase) CleanupExpiredSessions() error {
	return uc.sessionRepo.DeleteExpiredSessions()
}
