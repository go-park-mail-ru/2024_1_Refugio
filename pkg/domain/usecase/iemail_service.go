package usecase

import (
	emailCore "mail/pkg/domain/models"
)

// EmailUseCase represents the use case for working with emails.
type EmailUseCase interface {
	// GetAllEmails returns all emails.
	GetAllEmails() ([]*emailCore.Email, error)

	// GetEmailByID returns the email by its ID.
	GetEmailByID(id uint64) (*emailCore.Email, error)

	// CreateEmail creates a new email.
	CreateEmail(newEmail *emailCore.Email) (*emailCore.Email, error)

	// UpdateEmail updates the information of an email.
	UpdateEmail(updatedEmail *emailCore.Email) (bool, error)

	// DeleteEmail deletes the email.
	DeleteEmail(id uint64) (bool, error)
}
