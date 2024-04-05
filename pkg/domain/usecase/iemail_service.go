package usecase

import (
	emailCore "mail/pkg/domain/models"
)

// EmailUseCase represents the use case for working with emails.
type EmailUseCase interface {
	// GetAllEmails returns all emails.
	GetAllEmails(login string, offset, limit int) ([]*emailCore.Email, error)

	// GetEmailByID returns the email by its ID.
	GetEmailByID(id uint64, login string) (*emailCore.Email, error)

	// CreateEmail creates a new email.
	CreateEmail(newEmail *emailCore.Email) (int64, *emailCore.Email, error)

	CreateProfileEmail(email_id int64, sender, recipient string) error

	// UpdateEmail updates the information of an email.
	UpdateEmail(updatedEmail *emailCore.Email) (bool, error)

	// DeleteEmail deletes the email.
	DeleteEmail(id uint64, login string) (bool, error)

	CheckRecipientEmail(recipient string) error
}
