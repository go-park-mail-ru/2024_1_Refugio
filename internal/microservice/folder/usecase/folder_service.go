package usecase

import (
	"context"
	"mail/internal/pkg/utils/validators"

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
	return uc.repo.Create(newFolder, ctx)
}

// GetAllFolders list all folders.
func (uc *FolderUseCase) GetAllFolders(profileID uint32, offset, limit int64, ctx context.Context) ([]*domain.Folder, error) {
	return uc.repo.GetAll(profileID, offset, limit, ctx)
}

// DeleteFolder delete folder as user.
func (uc *FolderUseCase) DeleteFolder(folderID uint32, profileID uint32, ctx context.Context) (bool, error) {
	return uc.repo.Delete(folderID, profileID, ctx)
}

// UpdateFolder update folder as user.
func (uc *FolderUseCase) UpdateFolder(newUpFolder *domain.Folder, ctx context.Context) (bool, error) {
	return uc.repo.Update(newUpFolder, ctx)
}

// AddEmailInFolder add email in folder.
func (uc *FolderUseCase) AddEmailInFolder(folderID uint32, emailID uint32, ctx context.Context) (bool, error) {
	return uc.repo.AddEmailFolder(folderID, emailID, ctx)
}

// DeleteEmailInFolder delete email in folder.
func (uc *FolderUseCase) DeleteEmailInFolder(folderID uint32, emailID uint32, ctx context.Context) (bool, error) {
	return uc.repo.DeleteEmailFolder(folderID, emailID, ctx)
}

// CheckFolderProfile checking that the folder belongs to the user.
func (uc *FolderUseCase) CheckFolderProfile(folderID uint32, profileID uint32, ctx context.Context) (bool, error) {
	return uc.repo.CheckFolder(folderID, profileID, ctx)
}

// CheckEmailProfile checking that the email belongs to the user.
func (uc *FolderUseCase) CheckEmailProfile(emailID uint32, profileID uint32, ctx context.Context) (bool, error) {
	return uc.repo.CheckEmail(emailID, profileID, ctx)
}

// GetAllEmailsInFolder get all emails in folder as user.
func (uc *FolderUseCase) GetAllEmailsInFolder(folderID, profileID, limit, offset uint32, login string, ctx context.Context) ([]*domain.Email, error) {
	emails, err := uc.repo.GetAllEmails(folderID, profileID, limit, offset, ctx)
	if err != nil {
		return nil, err
	}

	for _, email := range emails {
		if validators.IsValidEmailFormat(email.SenderEmail) && validators.IsValidEmailFormat(email.RecipientEmail) {
			if email.SenderEmail == login {
				email.PhotoID, err = uc.repo.GetAvatarFileIDByLogin(email.RecipientEmail, ctx)
				if err != nil {
					return nil, err
				}
			} else if email.RecipientEmail == login {
				email.PhotoID, err = uc.repo.GetAvatarFileIDByLogin(email.SenderEmail, ctx)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return emails, nil
}

// GetAllFolderName retrieves the names of all folders associated with a given email ID.
func (uc *FolderUseCase) GetAllFolderName(emailID uint32, ctx context.Context) ([]*domain.Folder, error) {
	return uc.repo.GetAllFolderName(emailID, ctx)
}
