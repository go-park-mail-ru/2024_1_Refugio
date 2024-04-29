package domain_models

// Email represents the information about an email.
type Folder struct {
	ID        uint32 // The unique ID of the folder in the database.
	ProfileId uint32 // The unique identifier of the user who owns the folder.
	Name      string // The name of the folder.
}
