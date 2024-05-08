package repository_models

// Folder represents the information about an folder.
type Folder struct {
	ID        uint32 `db:"id"`         // ID he unique id of the folder in the database.
	ProfileId uint32 `db:"profile_id"` // ProfileId the unique identifier of the user who owns the folder.
	Name      string `db:"name"`       // Name the name of the folder.
}
