package usecase

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"mail/internal/microservice/session/mock"
	"mail/internal/pkg/logger"
	"mail/internal/pkg/utils/constants"

	domain "mail/internal/microservice/models/domain_models"
)

func GetCTX() context.Context {
	ctx := context.WithValue(context.Background(), constants.LoggerKey, logger.InitializationBdLog(nil))
	ctx2 := context.WithValue(ctx, constants.RequestIDKey, []string{"testID"})

	return ctx2
}

func TestCreateNewSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockSessionRepository(ctrl)
	usecase := NewSessionUseCase(mockRepo)

	ID := "10101010"
	userID := uint32(1)
	device := "testDevice"
	lifetime := 3600

	ctx := GetCTX()

	mockRepo.EXPECT().CreateSession(userID, device, lifetime, ctx).Return(ID, nil)

	sessionID, err := usecase.CreateNewSession(userID, device, lifetime, ctx)
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
	ctx := GetCTX()

	mockRepo.EXPECT().GetSessionByID(expectedSession.ID, ctx).Return(expectedSession, nil)

	session, err := usecase.GetSession(expectedSession.ID, ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedSession, session)
}

func TestDeleteSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockSessionRepository(ctrl)
	usecase := NewSessionUseCase(mockRepo)

	sessionID := "10101010"
	ctx := GetCTX()

	mockRepo.EXPECT().DeleteSessionByID(sessionID, ctx).Return(nil)

	err := usecase.DeleteSession(sessionID, ctx)
	assert.NoError(t, err)
}

func TestCleanupExpiredSessions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockSessionRepository(ctrl)
	usecase := NewSessionUseCase(mockRepo)

	ctx := GetCTX()

	mockRepo.EXPECT().DeleteExpiredSessions(ctx).Return(nil)

	err := usecase.CleanupExpiredSessions(ctx)
	assert.NoError(t, err)
}

func TestGetLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockSessionRepository(ctrl)
	usecase := NewSessionUseCase(mockRepo)

	sessionID := "10101010"

	ctx := GetCTX()

	expectedLogin := "testuser"

	mockRepo.EXPECT().GetLoginBySessionID(sessionID, ctx).Return(expectedLogin, nil)

	login, err := usecase.GetLogin(sessionID, ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedLogin, login)
}

func TestGetProfileID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockSessionRepository(ctrl)
	usecase := NewSessionUseCase(mockRepo)

	sessionID := "10101010"

	ctx := GetCTX()

	expectedProfileId := uint32(42)

	mockRepo.EXPECT().GetProfileIDBySessionID(sessionID, ctx).Return(expectedProfileId, nil)

	profileId, err := usecase.GetProfileID(sessionID, ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedProfileId, profileId)
}
