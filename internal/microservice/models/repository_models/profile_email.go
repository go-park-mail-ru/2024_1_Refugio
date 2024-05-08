package repository_models

// ProfileEmail represents the information about a profile_email.
type ProfileEmail struct {
	ProfileID uint32 `db:"profile_id"` // ProfileID he unique ID of the profile in the database.
	EmailID   uint32 `db:"email_id"`   // EmailID he unique ID of the email in the database.
}
