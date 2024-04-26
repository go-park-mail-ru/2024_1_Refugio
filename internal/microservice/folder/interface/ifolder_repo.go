package _interface

import (
	"context"
	domain "mail/internal/microservice/models/domain_models"
)

// FolderRepository represents the interface for working with folders.
type FolderRepository interface {
	// CreateFolder adds a new folder to the storage and returns its assigned unique identifier.
	CreateFolder(folder *domain.Folder, ctx context.Context) (uint32, *domain.Folder, error)

	// GetAll get list folder user.
	GetAll(profileID uint32, offset, limit int64, ctx context.Context) ([]*domain.Folder, error)
}
