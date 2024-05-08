package domain_models

// Folder represents the information about a folder.
type Folder struct {
	ID        uint32 // ID unique id of the folder in the database.
	ProfileId uint32 // ProfileId unique identifier of the user who owns the folder.
	Name      string // Name of the folder.
}
