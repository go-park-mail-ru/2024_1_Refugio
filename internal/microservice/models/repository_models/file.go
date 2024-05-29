package repository_models

// File represents information about a file.
type File struct {
	ID       uint64 `db:"id"`       // ID represents the unique identifier of the file in the database.
	FileId   string `db:"fileId"`   // FileId represents the identifier of the file.
	FileType string `db:"fileType"` // FileType represents the type of the file.
	FileName string `db:"fileName"` // FileName represents the name of the file.
	FileSize string `db:"fileSize"` // FileSize represents the size of the file.
}
