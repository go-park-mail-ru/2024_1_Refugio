package _interface

import (
	"context"
	folderCore "mail/internal/microservice/models/domain_models"
)

// FolderUseCase represents the use case for working with folders.
type FolderUseCase interface {
	// CreateFolder creates a new folder.
	CreateFolder(newFolder *folderCore.Folder, ctx context.Context) (uint32, *folderCore.Folder, error)
}
