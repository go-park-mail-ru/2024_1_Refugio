package usecase

import (
	emailCore "mail/pkg/domain/models"
)

// EmailUseCase represents the use case for working with emails.
type EmailUseCase interface {
	// GetAllEmailsIncoming returns all emails incoming.
	GetAllEmailsIncoming(login, requestID string, offset, limit int) ([]*emailCore.Email, error)

	// GetAllEmailsSent returns all emails sent.
	GetAllEmailsSent(login, requestID string, offset, limit int) ([]*emailCore.Email, error)

	// GetEmailByID returns the email by its ID.
	GetEmailByID(id uint64, login, requestID string) (*emailCore.Email, error)

	// CreateEmail creates a new email.
	CreateEmail(newEmail *emailCore.Email, requestID string) (int64, *emailCore.Email, error)

	CreateProfileEmail(email_id int64, sender, recipient, requestID string) error

	// UpdateEmail updates the information of an email.
	UpdateEmail(updatedEmail *emailCore.Email, requestID string) (bool, error)

	// DeleteEmail deletes the email.
	DeleteEmail(id uint64, login, requestID string) (bool, error)

	CheckRecipientEmail(recipient, requestID string) error
}
