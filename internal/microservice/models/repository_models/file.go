package repository_models

// File represents information about a file.
type File struct {
	ID       uint64 `json:"id"`       // ID represents the unique identifier of the file in the database.
	FileId   string `json:"fileId"`   // FileId represents the identifier of the file.
	FileType string `json:"fileType"` // FileType represents the type of the file.
}
