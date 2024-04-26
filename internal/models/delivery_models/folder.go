package delivery_models

// Folder represents the information about an email.
type Folder struct {
	ID        uint32 `json:"id,omitempty"`        // ID he unique ID of the folder in the database.
	ProfileId uint32 `json:"profileId,omitempty"` // ProfileId the unique identifier of the user who owns the folder.
	Name      string `json:"name"`                // Name the name of the folder.
}
