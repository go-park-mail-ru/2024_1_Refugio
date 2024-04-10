//go:generate mockgen -source=./imanager.go -destination=../mock/manager_mock.go -package=mock
package delivery

import (
	"mail/pkg/delivery/models"
	"net/http"
)

// SessionsManager represents the interface for managing user sessions.
type SessionsManager interface {
	// GetSession retrieves the session from the request and sanitizes it.
	GetSession(r *http.Request, requestID string) *models.Session

	// Check checks the validity of the session and CSRF token in the request.
	Check(r *http.Request, requestID string) (*models.Session, error)

	// CheckLogin checks if the login associated with the session matches the provided login.
	CheckLogin(login, requestID string, r *http.Request) error

	// GetLoginBySession retrieves the login associated with the session from the request.
	GetLoginBySession(r *http.Request, requestID string) (string, error)

	// Create creates a new session for the user and sets the session ID cookie in the response.
	Create(w http.ResponseWriter, userID uint32, requestID string) (*models.Session, error)

	// DestroyCurrent destroys the current session by deleting the session ID cookie from the response.
	DestroyCurrent(w http.ResponseWriter, r *http.Request, requestID string) error
}

