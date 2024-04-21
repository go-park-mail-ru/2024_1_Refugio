package domain_models

import "time"

// Session represents a user's session information.
type Session struct {
	ID           string    // ID uniquely identifies the session.
	UserID       uint32    // UserID specifies the ID of the user this session belongs to.
	CreationDate time.Time // CreationDate is the timestamp when the session was created.
	Device       string    // Device describes the device used to initiate the session, e.g., 'web', 'mobile'.
	LifeTime     int       // LifeTime indicates the duration (in seconds) for which the session is valid.
	CsrfToken    string    // CsrfToken represents the Cross-Site Request Forgery (CSRF) token associated with the session.
}
