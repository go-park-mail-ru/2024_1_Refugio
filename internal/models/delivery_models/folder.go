package delivery_models

// Folder represents the information about a folder.
type Folder struct {
	ID        uint32 `json:"id,omitempty"`        // ID he unique ID of the folder in the database.
	ProfileId uint32 `json:"profileId,omitempty"` // ProfileId the unique identifier of the user who owns the folder.
	Name      string `json:"name"`                // Name the name of the folder.
}

// FolderEmail represents the information about an folderID and emailID.
type FolderEmail struct {
	FolderID uint32 `json:"folderId"` // FolderID he unique ID of the folder in the database.
	EmailID  uint32 `json:"emailId"`  // EmailID he unique ID of the email in the database.
}
