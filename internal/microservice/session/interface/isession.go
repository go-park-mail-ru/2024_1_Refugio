//go:generate mockgen -source=./isession.go -destination=../mock/session_mock.go -package=mock

package _interface

import (
	"context"

	"mail/internal/microservice/session/proto"
)

// SessionServer represents the interface for working with users.
type SessionServer interface {
	// GetSession retrieves the session.
	GetSession(ctx context.Context, input *proto.GetSessionRequest) (*proto.GetSessionReply, error)

	// GetLoginBySession retrieves the login associated with the session.
	GetLoginBySession(ctx context.Context, input *proto.GetLoginBySessionRequest) (*proto.GetLoginBySessionReply, error)

	// GetProfileIDBySession retrieves the login associated with the session.
	GetProfileIDBySession(ctx context.Context, input *proto.GetLoginBySessionRequest) (*proto.GetProfileIDBySessionReply, error)

	// CreateSession creates a new session for the user.
	CreateSession(ctx context.Context, input *proto.CreateSessionRequest) (*proto.CreateSessionReply, error)

	// DeleteSession destroys the current session.
	DeleteSession(ctx context.Context, input *proto.DeleteSessionRequest) (*proto.DeleteSessionReply, error)

	// CleanupExpiredSessions destroys all current session.
	CleanupExpiredSessions(ctx context.Context, input *proto.CleanupExpiredSessionsRequest) (*proto.CleanupExpiredSessionsReply, error)
}
