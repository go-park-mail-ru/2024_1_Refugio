package models

import "time"

// Session represents a user's session information.
type Session struct {
	ID           uint32    `db:"id"`            // ID uniquely identifies the session.
	UserID       uint32    `db:"user_id"`       // UserID specifies the ID of the user this session belongs to.
	CreationDate time.Time `db:"creation_date"` // CreationDate is the timestamp when the session was created.
	Device       string    `db:"device"`        // Device describes the device used to initiate the session, e.g., 'web', 'mobile'.
	LifeTime     int       `db:"lifetime"`      // LifeTime indicates the duration (in seconds) for which the session is valid.
}
