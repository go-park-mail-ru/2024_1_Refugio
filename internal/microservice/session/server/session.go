package server

import (
	"context"
	"fmt"

	"mail/internal/microservice/models/proto_converters"
	"mail/internal/microservice/session/proto"

	usecase "mail/internal/microservice/session/interface"
	validUtil "mail/internal/pkg/utils/validators"
)

// SessionServer handles RPC calls for the SessionService.
type SessionServer struct {
	proto.UnimplementedSessionServiceServer
	SessionUseCase usecase.SessionUseCase
}

// NewSessionServer creates a new instance of SessionServer.
func NewSessionServer(sessionUseCase usecase.SessionUseCase) *SessionServer {
	return &SessionServer{SessionUseCase: sessionUseCase}
}

// GetSession retrieves the session.
func (ss *SessionServer) GetSession(ctx context.Context, input *proto.GetSessionRequest) (*proto.GetSessionReply, error) {
	if validUtil.IsEmpty(input.SessionId) {
		return nil, fmt.Errorf("session not found")
	}

	sessionCore, err := ss.SessionUseCase.GetSession(input.SessionId, ctx)
	if err != nil {
		return nil, fmt.Errorf("session not found")
	}

	return &proto.GetSessionReply{Session: proto_converters.SessionConvertCoreInProto(sessionCore)}, nil
}

// GetLoginBySession retrieves the login associated with the session.
func (ss *SessionServer) GetLoginBySession(ctx context.Context, input *proto.GetLoginBySessionRequest) (*proto.GetLoginBySessionReply, error) {
	if validUtil.IsEmpty(input.SessionId) {
		return nil, fmt.Errorf("session not found")
	}

	login, err := ss.SessionUseCase.GetLogin(input.SessionId, ctx)
	if err != nil {
		return nil, fmt.Errorf("session not found")
	}

	return &proto.GetLoginBySessionReply{Login: login}, nil
}

// GetProfileIDBySession retrieves the login associated with the session.
func (ss *SessionServer) GetProfileIDBySession(ctx context.Context, input *proto.GetLoginBySessionRequest) (*proto.GetProfileIDBySessionReply, error) {
	if validUtil.IsEmpty(input.SessionId) {
		return nil, fmt.Errorf("session not found")
	}

	profileId, err := ss.SessionUseCase.GetProfileID(input.SessionId, ctx)
	if err != nil {
		return nil, fmt.Errorf("session not found")
	}

	return &proto.GetProfileIDBySessionReply{Id: profileId}, nil
}

// CreateSession creates a new session for the user.
func (ss *SessionServer) CreateSession(ctx context.Context, input *proto.CreateSessionRequest) (*proto.CreateSessionReply, error) {
	if input.Session.UserId <= 0 {
		return nil, fmt.Errorf("user id not found")
	}

	if input.Session.LifeTime <= 0 {
		return nil, fmt.Errorf("life time session error")
	}

	sessionId, err := ss.SessionUseCase.CreateNewSession(input.Session.UserId, input.Session.Device, int(input.Session.LifeTime), ctx)
	if err != nil {
		return nil, fmt.Errorf("session not found")
	}

	return &proto.CreateSessionReply{SessionId: sessionId}, nil
}

// DeleteSession destroys the current session.
func (ss *SessionServer) DeleteSession(ctx context.Context, input *proto.DeleteSessionRequest) (*proto.DeleteSessionReply, error) {
	if validUtil.IsEmpty(input.SessionId) {
		return nil, fmt.Errorf("session not found")
	}

	err := ss.SessionUseCase.DeleteSession(input.SessionId, ctx)
	if err != nil {
		return &proto.DeleteSessionReply{Status: false}, fmt.Errorf("session not found")
	}

	return &proto.DeleteSessionReply{Status: true}, nil
}

// CleanupExpiredSessions destroys all current session.
func (ss *SessionServer) CleanupExpiredSessions(ctx context.Context, input *proto.CleanupExpiredSessionsRequest) (*proto.CleanupExpiredSessionsReply, error) {
	err := ss.SessionUseCase.CleanupExpiredSessions(ctx)
	if err != nil {
		return &proto.CleanupExpiredSessionsReply{}, fmt.Errorf("session cleanup expired")
	}

	return &proto.CleanupExpiredSessionsReply{}, nil
}
