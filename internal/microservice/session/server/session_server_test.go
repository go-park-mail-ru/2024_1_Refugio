package server

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"mail/internal/microservice/models/domain_models"
	"mail/internal/microservice/session/mock"
	"mail/internal/microservice/session/proto"
	"mail/internal/pkg/logger"
	"mail/internal/pkg/utils/constants"
)

func GetCTX() context.Context {
	ctx := context.WithValue(context.Background(), constants.LoggerKey, logger.InitializationBdLog(nil))
	ctx2 := context.WithValue(ctx, constants.RequestIDKey, []string{"testID"})

	return ctx2
}

func TestGetSession_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionUseCase := mock.NewMockSessionUseCase(ctrl)

	server := NewSessionServer(mockSessionUseCase)

	ctx := GetCTX()

	expectedSession := &domain_models.Session{ID: "10101010", UserID: uint32(32), CsrfToken: "01010101"}

	mockSessionUseCase.EXPECT().GetSession(gomock.Any(), ctx).Return(expectedSession, nil)

	session, err := server.GetSession(ctx, &proto.GetSessionRequest{SessionId: "some_session_id"})

	assert.NoError(t, err)
	assert.NotNil(t, session)
}

func TestGetSession_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionUseCase := mock.NewMockSessionUseCase(ctrl)

	server := NewSessionServer(mockSessionUseCase)

	ctx := GetCTX()

	mockSessionUseCase.EXPECT().GetSession(gomock.Any(), ctx).Return(nil, fmt.Errorf("session not found"))

	session, err := server.GetSession(ctx, &proto.GetSessionRequest{SessionId: "some_session_id"})

	assert.Error(t, err)
	assert.Nil(t, session)
}

func TestGetLoginBySession_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionUseCase := mock.NewMockSessionUseCase(ctrl)

	server := NewSessionServer(mockSessionUseCase)

	ctx := GetCTX()

	expectedLogin := "example_login"

	mockSessionUseCase.EXPECT().GetLogin("some_session_id", ctx).Return(expectedLogin, nil)

	reply, err := server.GetLoginBySession(ctx, &proto.GetLoginBySessionRequest{SessionId: "some_session_id"})

	assert.NoError(t, err)
	assert.Equal(t, expectedLogin, reply.Login)
}

func TestGetLoginBySession_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionUseCase := mock.NewMockSessionUseCase(ctrl)

	server := NewSessionServer(mockSessionUseCase)

	ctx := GetCTX()

	mockSessionUseCase.EXPECT().GetLogin("some_session_id", ctx).Return("", fmt.Errorf("session not found"))

	reply, err := server.GetLoginBySession(ctx, &proto.GetLoginBySessionRequest{SessionId: "some_session_id"})

	assert.Error(t, err)
	assert.Nil(t, reply)
}

func TestGetProfileIDBySession_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionUseCase := mock.NewMockSessionUseCase(ctrl)

	server := NewSessionServer(mockSessionUseCase)

	ctx := GetCTX()

	expectedProfileID := uint32(42)

	mockSessionUseCase.EXPECT().GetProfileID("some_session_id", ctx).Return(expectedProfileID, nil)

	reply, err := server.GetProfileIDBySession(ctx, &proto.GetLoginBySessionRequest{SessionId: "some_session_id"})

	assert.NoError(t, err)
	assert.Equal(t, expectedProfileID, reply.Id)
}

func TestGetProfileIDBySession_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionUseCase := mock.NewMockSessionUseCase(ctrl)

	server := NewSessionServer(mockSessionUseCase)

	ctx := context.Background()

	mockSessionUseCase.EXPECT().GetProfileID("some_session_id", ctx).Return(uint32(0), fmt.Errorf("session not found"))

	reply, err := server.GetProfileIDBySession(ctx, &proto.GetLoginBySessionRequest{SessionId: "some_session_id"})

	assert.Error(t, err)
	assert.Nil(t, reply)
}

func TestCreateSession_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionUseCase := mock.NewMockSessionUseCase(ctrl)

	server := NewSessionServer(mockSessionUseCase)

	ctx := GetCTX()

	expectedSessionID := "abc123"

	mockSessionUseCase.EXPECT().CreateNewSession(uint32(42), "device", 3600, ctx).Return(expectedSessionID, nil)

	reply, err := server.CreateSession(ctx, &proto.CreateSessionRequest{
		Session: &proto.Session{
			UserId:   42,
			Device:   "device",
			LifeTime: 3600,
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, expectedSessionID, reply.SessionId)
}

func TestCreateSession_UserIdNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionUseCase := mock.NewMockSessionUseCase(ctrl)

	server := NewSessionServer(mockSessionUseCase)

	ctx := GetCTX()

	reply, err := server.CreateSession(ctx, &proto.CreateSessionRequest{
		Session: &proto.Session{
			UserId:   0,
			Device:   "device",
			LifeTime: 3600,
		},
	})

	assert.Error(t, err)
	assert.Nil(t, reply)
}

func TestCreateSession_LifeTimeError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionUseCase := mock.NewMockSessionUseCase(ctrl)

	server := NewSessionServer(mockSessionUseCase)

	ctx := GetCTX()

	reply, err := server.CreateSession(ctx, &proto.CreateSessionRequest{
		Session: &proto.Session{
			UserId:   42,
			Device:   "device",
			LifeTime: 0,
		},
	})

	assert.Error(t, err)
	assert.Nil(t, reply)
}

func TestCreateSession_SessionNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionUseCase := mock.NewMockSessionUseCase(ctrl)

	server := NewSessionServer(mockSessionUseCase)

	ctx := GetCTX()

	mockSessionUseCase.EXPECT().CreateNewSession(uint32(42), "device", 3600, ctx).Return("", fmt.Errorf("session not found"))

	reply, err := server.CreateSession(ctx, &proto.CreateSessionRequest{
		Session: &proto.Session{
			UserId:   42,
			Device:   "device",
			LifeTime: 3600,
		},
	})

	assert.Error(t, err)
	assert.Nil(t, reply)
}

func TestDeleteSession_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionUseCase := mock.NewMockSessionUseCase(ctrl)

	server := NewSessionServer(mockSessionUseCase)

	ctx := GetCTX()

	mockSessionUseCase.EXPECT().DeleteSession("some_session_id", ctx).Return(nil)

	reply, err := server.DeleteSession(ctx, &proto.DeleteSessionRequest{SessionId: "some_session_id"})

	assert.NoError(t, err)
	assert.True(t, reply.Status)
}

func TestDeleteSession_SessionNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionUseCase := mock.NewMockSessionUseCase(ctrl)

	server := NewSessionServer(mockSessionUseCase)

	ctx := GetCTX()

	mockSessionUseCase.EXPECT().DeleteSession("some_session_id", ctx).Return(fmt.Errorf("session not found"))

	reply, err := server.DeleteSession(ctx, &proto.DeleteSessionRequest{SessionId: "some_session_id"})

	assert.Error(t, err)
	assert.False(t, reply.Status)
}

func TestDeleteSession_SessionIdNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionUseCase := mock.NewMockSessionUseCase(ctrl)

	server := NewSessionServer(mockSessionUseCase)

	ctx := GetCTX()

	reply, err := server.DeleteSession(ctx, &proto.DeleteSessionRequest{SessionId: ""})

	assert.Error(t, err)
	assert.Nil(t, reply)
}

func TestCleanupExpiredSessions_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionUseCase := mock.NewMockSessionUseCase(ctrl)

	server := NewSessionServer(mockSessionUseCase)

	ctx := GetCTX()

	mockSessionUseCase.EXPECT().CleanupExpiredSessions(ctx).Return(nil)

	reply, err := server.CleanupExpiredSessions(ctx, &proto.CleanupExpiredSessionsRequest{})

	assert.NoError(t, err)
	assert.NotNil(t, reply)
}

func TestCleanupExpiredSessions_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionUseCase := mock.NewMockSessionUseCase(ctrl)

	server := NewSessionServer(mockSessionUseCase)

	ctx := GetCTX()

	mockSessionUseCase.EXPECT().CleanupExpiredSessions(ctx).Return(fmt.Errorf("session cleanup expired"))

	reply, err := server.CleanupExpiredSessions(ctx, &proto.CleanupExpiredSessionsRequest{})

	assert.Error(t, err)
	assert.NotNil(t, reply)
}
