//go:generate mockgen -source=./iemail_service.go -destination=../mock/email_service_mock.go -package=mock

package _interface

import (
	"context"

	emailCore "mail/internal/microservice/models/domain_models"
)

// EmailUseCase represents the use case for working with emails.
type EmailUseCase interface {
	// GetAllEmailsIncoming returns all incoming emails for the specified user.
	GetAllEmailsIncoming(login string, offset, limit int64, ctx context.Context) ([]*emailCore.Email, error)

	// GetAllEmailsSent returns all sent emails for the specified user.
	GetAllEmailsSent(login string, offset, limit int64, ctx context.Context) ([]*emailCore.Email, error)

	// GetAllDraftEmails returns all draft emails for the specified user.
	GetAllDraftEmails(login string, offset, limit int64, ctx context.Context) ([]*emailCore.Email, error)

	// GetAllSpamEmails returns all spam emails for the specified user.
	GetAllSpamEmails(login string, offset, limit int64, ctx context.Context) ([]*emailCore.Email, error)

	// GetEmailByID returns the email with the specified ID for the specified user.
	GetEmailByID(id uint64, login string, ctx context.Context) (*emailCore.Email, error)

	// CreateEmail creates a new email.
	CreateEmail(newEmail *emailCore.Email, ctx context.Context) (uint64, *emailCore.Email, error)

	// CreateProfileEmail creates a new profile email.
	CreateProfileEmail(emailId uint64, sender, recipient string, ctx context.Context) error

	// UpdateEmail updates the information of the specified email.
	UpdateEmail(updatedEmail *emailCore.Email, ctx context.Context) (bool, error)

	// DeleteEmail deletes the email with the specified ID for the specified user.
	DeleteEmail(id uint64, login string, ctx context.Context) (bool, error)

	// CheckRecipientEmail checks if the recipient email address is valid.
	CheckRecipientEmail(recipient string, ctx context.Context) error

	// AddAttachment adds an attachment to the specified email.
	AddAttachment(fileID, fileType, fileName, fileSize string, emailID uint64, ctx context.Context) (uint64, error)

	// GetFileByID returns the file with the specified ID.
	GetFileByID(fileID uint64, ctx context.Context) (*emailCore.File, error)

	// GetFilesByEmailID returns all files attached to the specified email.
	GetFilesByEmailID(emailID uint64, ctx context.Context) ([]*emailCore.File, error)

	// DeleteFileByID deletes the file with the specified ID.
	DeleteFileByID(fileID uint64, ctx context.Context) (bool, error)

	// UpdateFileByID updates the information of the specified file.
	UpdateFileByID(fileID uint64, newFileID, newFileType, newFileName, newFileSize string, ctx context.Context) (bool, error)

	// AddFile add an file to database.
	AddFile(fileID, fileType, fileName, fileSize string, ctx context.Context) (uint64, error)

	// AddFileToEmail add a file to an email.
	AddFileToEmail(emailID uint64, fileID uint64, ctx context.Context) error
}
