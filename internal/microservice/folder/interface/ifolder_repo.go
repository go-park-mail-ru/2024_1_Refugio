//go:generate mockgen -source=./ifolder_repo.go -destination=../mock/folder_repo_mock.go -package=mock

package _interface

import (
	"context"

	domain "mail/internal/microservice/models/domain_models"
)

// FolderRepository represents the interface for working with folders.
type FolderRepository interface {
	// Create adds a new folder to the storage and returns its assigned unique identifier.
	Create(folder *domain.Folder, ctx context.Context) (uint32, *domain.Folder, error)

	// GetAll get list folder user.
	GetAll(profileID uint32, offset, limit int64, ctx context.Context) ([]*domain.Folder, error)

	// Delete delete folder as user.
	Delete(folderID uint32, profileID uint32, ctx context.Context) (bool, error)

	// Update folder as user.
	Update(newUpFolder *domain.Folder, ctx context.Context) (bool, error)

	// AddEmailFolder adds a new email in folder to the storage and returns its assigned unique identifier.
	AddEmailFolder(folderID uint32, emailID uint32, ctx context.Context) (bool, error)

	// DeleteEmailFolder adds a new email in folder to the storage and returns its assigned unique identifier.
	DeleteEmailFolder(folderID uint32, emailID uint32, ctx context.Context) (bool, error)

	// CheckFolder checking that the folder belongs to the user.
	CheckFolder(folderID uint32, profileID uint32, ctx context.Context) (bool, error)

	// CheckEmail checking that the email belongs to the user.
	CheckEmail(emailID uint32, profileID uint32, ctx context.Context) (bool, error)

	// GetAllEmails get list emails folder user.
	GetAllEmails(folderID, profileId, limit, offset uint32, ctx context.Context) ([]*domain.Email, error)

	// GetAvatarFileIDByLogin getting an avatar by login.
	GetAvatarFileIDByLogin(login string, ctx context.Context) (string, error)

	// GetAllFolderName retrieves the names of all folders associated with a given email ID.
	GetAllFolderName(emailID uint32, ctx context.Context) ([]*domain.Folder, error)
}
