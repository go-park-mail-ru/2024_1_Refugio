package usecase

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mock "mail/internal/pkg/auth/mocks"
	"testing"

	domain "mail/internal/models/domain_models"
)

func TestCreateNewSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockSessionRepository(ctrl)
	usecase := NewSessionUseCase(mockRepo)

	ID := "10101010"
	userID := uint32(1)
	device := "testDevice"
	lifetime := 3600
	requestID := "testRequestID"

	mockRepo.EXPECT().CreateSession(userID, device, requestID, lifetime).Return(ID, nil)

	sessionID, err := usecase.CreateNewSession(userID, device, requestID, lifetime)
	assert.NoError(t, err)
	assert.Equal(t, ID, sessionID)
}

func TestGetSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockSessionRepository(ctrl)
	usecase := NewSessionUseCase(mockRepo)

	expectedSession := &domain.Session{
		ID:       "10101010",
		UserID:   uint32(100),
		Device:   "testDevice",
		LifeTime: 3600,
	}
	requestID := "testRequestID"

	mockRepo.EXPECT().GetSessionByID(expectedSession.ID, requestID).Return(expectedSession, nil)

	session, err := usecase.GetSession(expectedSession.ID, requestID)
	assert.NoError(t, err)
	assert.Equal(t, expectedSession, session)
}

func TestDeleteSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockSessionRepository(ctrl)
	usecase := NewSessionUseCase(mockRepo)

	sessionID := "10101010"
	requestID := "testRequestID"

	mockRepo.EXPECT().DeleteSessionByID(sessionID, requestID).Return(nil)

	err := usecase.DeleteSession(sessionID, requestID)
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

func TestGetLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockSessionRepository(ctrl)
	usecase := NewSessionUseCase(mockRepo)

	sessionID := "10101010"
	requestID := "test_request"

	expectedLogin := "testuser"

	mockRepo.EXPECT().GetLoginBySessionID(sessionID, requestID).Return(expectedLogin, nil)

	login, err := usecase.GetLogin(sessionID, requestID)
	assert.NoError(t, err)
	assert.Equal(t, expectedLogin, login)
}
