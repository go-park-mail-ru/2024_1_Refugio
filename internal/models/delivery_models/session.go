package delivery_models

import "time"

// Session represents a user's session information.
type Session struct {
	ID           string    `json:"id,omitempty"`            // ID uniquely identifies the session.
	UserID       uint32    `json:"user-id,omitempty"`       // UserID specifies the ID of the user this session belongs to.
	CreationDate time.Time `json:"creation-date,omitempty"` // CreationDate is the timestamp when the session was created.
	Device       string    `json:"device,omitempty"`        // Device describes the device used to initiate the session, e.g., 'web', 'mobile'.
	LifeTime     int       `json:"life-time,omitempty"`     // LifeTime indicates the duration (in seconds) for which the session is valid.
	CsrfToken    string    `json:"csrf_token,omitempty"`    // CsrfToken represents the Cross-Site Request Forgery (CSRF) token associated with the session.
}
