package _interface

import (
	"context"
	domain "mail/internal/microservice/models/domain_models"
)

// EmailRepository represents the interface for working with emails.
type FolderRepository interface {
	// Add adds a new email to the storage and returns its assigned unique identifier.
	CreateFolder(email *domain.Folder, ctx context.Context) (uint64, *domain.Folder, error)
}
