package repository_models

import "time"

// Session represents a user's session information.
type Session struct {
	ID           string    `db:"id"`            // ID uniquely identifies the session.
	UserID       uint32    `db:"profile_id"`    // UserID specifies the ID of the user this session belongs to.
	CreationDate time.Time `db:"creation_date"` // CreationDate is the timestamp when the session was created.
	Device       string    `db:"device"`        // Device describes the device used to initiate the session, e.g., 'web', 'mobile'.
	LifeTime     int       `db:"life_time"`     // LifeTime indicates the duration (in seconds) for which the session is valid.
	CsrfToken    string    `db:"csrf_token"`    // CsrfToken represents the Cross-Site Request Forgery (CSRF) token associated with the session.
}
