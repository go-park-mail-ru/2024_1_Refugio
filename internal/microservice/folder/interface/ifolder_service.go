//go:generate mockgen -source=./ifolder_service.go -destination=../mock/folder_service_mock.go -package=mock

package _interface

import (
	"context"

	folderCore "mail/internal/microservice/models/domain_models"
)

// FolderUseCase represents the use case for working with folders.
type FolderUseCase interface {
	// CreateFolder creates a new folder.
	CreateFolder(newFolder *folderCore.Folder, ctx context.Context) (uint32, *folderCore.Folder, error)

	// GetAllFolders get all folders as user.
	GetAllFolders(profileID uint32, offset, limit int64, ctx context.Context) ([]*folderCore.Folder, error)

	// DeleteFolder delete folder as user.
	DeleteFolder(folderID uint32, profileID uint32, ctx context.Context) (bool, error)

	// UpdateFolder update folder as user.
	UpdateFolder(newUpFolder *folderCore.Folder, ctx context.Context) (bool, error)

	// AddEmailInFolder add email in folder.
	AddEmailInFolder(folderID uint32, emailID uint32, ctx context.Context) (bool, error)

	// CheckFolderProfile checking that the folder belongs to the user.
	CheckFolderProfile(folderID uint32, profileID uint32, ctx context.Context) (bool, error)

	// CheckEmailProfile checking that the email belongs to the user.
	CheckEmailProfile(emailID uint32, profileID uint32, ctx context.Context) (bool, error)

	// GetAllEmailsInFolder get all emails in folder as user.
	GetAllEmailsInFolder(folderID, profileID, limit, offset uint32, login string, ctx context.Context) ([]*folderCore.Email, error)

	// DeleteEmailInFolder delete email in folder.
	DeleteEmailInFolder(folderID uint32, emailID uint32, ctx context.Context) (bool, error)

	// GetAllFolderName retrieves the names of all folders associated with a given email ID.
	GetAllFolderName(emailID uint32, ctx context.Context) ([]*folderCore.Folder, error)
}
