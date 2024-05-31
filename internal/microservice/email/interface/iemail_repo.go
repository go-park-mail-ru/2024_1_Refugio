//go:generate mockgen -source=./iemail_repo.go -destination=../mock/email_repository_mock.go -package=mock

package _interface

import (
	"context"

	domain "mail/internal/microservice/models/domain_models"
)

// EmailRepository represents the interface for working with emails.
type EmailRepository interface {
	// GetAllIncoming returns all emails incoming from the storage.
	GetAllIncoming(login string, offset, limit int64, ctx context.Context) ([]*domain.Email, error)

	// GetAllSent returns all emails sent from the storage.
	GetAllSent(login string, offset, limit int64, ctx context.Context) ([]*domain.Email, error)

	// GetAllDraft returns all draft emails from the storage.
	GetAllDraft(login string, offset, limit int64, ctx context.Context) ([]*domain.Email, error)

	// GetAllSpam returns all draft emails from the storage.
	GetAllSpam(login string, offset, limit int64, ctx context.Context) ([]*domain.Email, error)

	// GetByID returns the email by its unique identifier.
	GetByID(id uint64, login string, ctx context.Context) (*domain.Email, error)

	// GetAvatarFileIDByLogin getting an avatar by login.
	GetAvatarFileIDByLogin(login string, ctx context.Context) (string, error)

	// Add adds a new email to the storage and returns its assigned unique identifier.
	Add(emailModelCore *domain.Email, ctx context.Context) (uint64, *domain.Email, error)

	// AddProfileEmail links an email to one or more profiles based on sender and recipient information.
	AddProfileEmail(email_id uint64, sender, recipient string, ctx context.Context) error

	// AddProfileEmailMyself links an email to the profile corresponding to the sender (when sender and recipient are the same).
	AddProfileEmailMyself(email_id uint64, login string, ctx context.Context) error

	// Update updates the information of an email in the storage based on the provided new email.
	Update(newEmail *domain.Email, ctx context.Context) (bool, error)

	// Delete removes the email from the storage by its unique identifier.
	Delete(id uint64, login string, ctx context.Context) (bool, error)

	// FindEmail searches for a user in the database based on their login.
	FindEmail(login string, ctx context.Context) error

	// AddFile adds a file entry to the database with the provided file ID, file type, file name and file size.
	AddFile(fileID string, fileType string, fileName string, fileSize string, ctx context.Context) (uint64, error)

	// AddAttachment links a file to an email by inserting a record into the email_file table.
	AddAttachment(emailID uint64, fileID uint64, ctx context.Context) error

	// GetFileByID retrieves file information based on the provided file ID.
	GetFileByID(id uint64, ctx context.Context) (*domain.File, error)

	// GetFilesByEmailID retrieves all files associated with a given email ID.
	GetFilesByEmailID(emailID uint64, ctx context.Context) ([]*domain.File, error)

	// DeleteFileByID deletes a file entry from the database based on the provided file ID.
	DeleteFileByID(fileID uint64, ctx context.Context) error

	// UpdateFileByID updates the file ID, file type, file name and file size of a file entry in the database based on the provided file ID.
	UpdateFileByID(fileID uint64, newFileID string, newFileType string, newFileName string, newFileSize string, ctx context.Context) error
}
