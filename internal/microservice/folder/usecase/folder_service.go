package usecase

import (
	"context"
	repository "mail/internal/microservice/folder/interface"
	domain "mail/internal/microservice/models/domain_models"
)

// FolderUseCase represents the use case for working with folders.
type FolderUseCase struct {
	repo repository.FolderRepository
}

// NewFolderUseCase creates a new instance of FolderUseCase.
func NewFolderUseCase(repo repository.FolderRepository) *FolderUseCase {
	return &FolderUseCase{
		repo: repo,
	}
}

// CreateFolder new folder.
func (uc *FolderUseCase) CreateFolder(newFolder *domain.Folder, ctx context.Context) (uint32, *domain.Folder, error) {
	return uc.repo.CreateFolder(newFolder, ctx)
}

// GetAllFolders list all folders.
func (uc *FolderUseCase) GetAllFolders(profileID uint32, offset, limit int64, ctx context.Context) ([]*domain.Folder, error) {
	return uc.repo.GetAll(profileID, offset, limit, ctx)
}
