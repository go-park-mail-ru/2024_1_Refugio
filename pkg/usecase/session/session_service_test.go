package session

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"mail/pkg/domain/mock"
	"testing"

	domain "mail/pkg/domain/models"
)

func TestCreateNewSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockSessionRepository(ctrl)
	usecase := NewSessionUseCase(mockRepo)

	ID := uint32(1)
	userID := uint32(1)
	device := "testDevice"
	lifetime := 3600

	mockRepo.EXPECT().CreateSession(ID, userID, device, lifetime).Return(ID, nil)

	sessionID, err := usecase.CreateNewSession(ID, userID, device, lifetime)
	assert.NoError(t, err)
	assert.Equal(t, userID, sessionID)
}

func TestGetSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockSessionRepository(ctrl)
	usecase := NewSessionUseCase(mockRepo)

	expectedSession := &domain.Session{
		ID:       uint32(1),
		UserID:   uint32(100),
		Device:   "testDevice",
		LifeTime: 3600,
	}

	mockRepo.EXPECT().GetSessionByID(expectedSession.ID).Return(expectedSession, nil)

	session, err := usecase.GetSession(expectedSession.ID)
	assert.NoError(t, err)
	assert.Equal(t, expectedSession, session)
}

func TestDeleteSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockSessionRepository(ctrl)
	usecase := NewSessionUseCase(mockRepo)

	sessionID := uint32(1)

	mockRepo.EXPECT().DeleteSessionByID(sessionID).Return(nil)

	err := usecase.DeleteSession(sessionID)
	assert.NoError(t, err)
}

func TestCleanupExpiredSessions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockSessionRepository(ctrl)
	usecase := NewSessionUseCase(mockRepo)

	mockRepo.EXPECT().DeleteExpiredSessions().Return(nil)

	err := usecase.CleanupExpiredSessions()
	assert.NoError(t, err)
}
