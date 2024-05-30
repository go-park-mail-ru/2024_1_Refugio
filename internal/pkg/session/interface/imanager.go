//go:generate mockgen -source=./imanager.go -destination=../mock/manager_mock.go -package=mock

package _interface

import (
	"context"
	"net/http"

	api "mail/internal/models/delivery_models"
)

// SessionsManager represents the interface for managing user sessions.
type SessionsManager interface {
	// SetSession set the session in the request.
	SetSession(sessionId string, w http.ResponseWriter, r *http.Request, ctx context.Context) error

	// GetSession retrieves the session from the request and sanitizes it.
	GetSession(r *http.Request, ctx context.Context) *api.Session

	// Check checks the validity of the session and CSRF token in the request.
	Check(r *http.Request, ctx context.Context) (*api.Session, error)

	// CheckLogin checks if the login associated with the session matches the provided login.
	CheckLogin(login string, r *http.Request, ctx context.Context) error

	// GetLoginBySession retrieves the login associated with the session from the request.
	GetLoginBySession(r *http.Request, ctx context.Context) (string, error)

	// GetProfileIDBySessionID retrieves the profile ID associated with the given session ID from the session service.
	GetProfileIDBySessionID(r *http.Request, ctx context.Context) (uint32, error)

	// Create creates a new session for the user and sets the session ID cookie in the response.
	Create(w http.ResponseWriter, userID uint32, ctx context.Context) (*api.Session, error)

	// DestroyCurrent destroys the current session by deleting the session ID cookie from the response.
	DestroyCurrent(w http.ResponseWriter, r *http.Request, ctx context.Context) error
}
