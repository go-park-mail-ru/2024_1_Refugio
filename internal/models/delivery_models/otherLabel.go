package delivery_models

// OtherLabel represents the information about a folder.
type OtherLabel struct {
	ID        string `json:"id,omitempty"`        // ID he unique ID of the folder in the database.
	ProfileId uint32 `json:"profileId,omitempty"` // ProfileId the unique identifier of the user who owns the folder.
	Name      string `json:"name"`                // Name the name of the folder.
}
