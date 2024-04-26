package usecase

import (
	"context"
	repository "mail/internal/microservice/folder/interface"
	domain "mail/internal/microservice/models/domain_models"
)

// FolderUseCase represents the use case for working with emails.
type FolderUseCase struct {
	repo repository.FolderRepository
}

// NewFolderUseCase creates a new instance of EmailUseCase.
func NewFolderUseCase(repo repository.FolderRepository) *FolderUseCase {
	return &FolderUseCase{
		repo: repo,
	}
}

// GetAllEmailsIncoming returns all emails incoming.
func (uc *FolderUseCase) CreateFolder(newFolder *domain.Folder, ctx context.Context) (uint64, *domain.Folder, error) {
	return uc.repo.CreateFolder(newFolder, ctx)
}
