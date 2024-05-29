package domain_models

// File represents information about a file.
type File struct {
	ID       uint64 // ID represents the unique identifier of the file in the database.
	FileId   string // FileId represents the identifier of the file.
	FileType string // FileType represents the type of the file.
	FileName string // FileName represents the name of the file.
	FileSize string // FileSize represents the size of the file.
}
