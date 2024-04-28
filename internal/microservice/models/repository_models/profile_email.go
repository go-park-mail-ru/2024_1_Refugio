package repository_models

// ProfileEmail represents the information about an profile_email.
type ProfileEmail struct {
	ProfileID uint32 `db:"profile_id"` // ID he unique ID of the profile in the database.
	EmailID   uint32 `db:"email_id"`   // ID he unique ID of the email in the database.
}
