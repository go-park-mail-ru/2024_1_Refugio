package _interface

import (
	"context"
	emailCore "mail/internal/models/domain_models"
)

// EmailUseCase represents the use case for working with emails.
type EmailUseCase interface {
	// GetAllEmailsIncoming returns all emails incoming.
	GetAllEmailsIncoming(login string, offset, limit int, ctx context.Context) ([]*emailCore.Email, error)

	// GetAllEmailsSent returns all emails sent.
	GetAllEmailsSent(login string, offset, limit int, ctx context.Context) ([]*emailCore.Email, error)

	// GetEmailByID returns the email by its ID.
	GetEmailByID(id uint64, login string, ctx context.Context) (*emailCore.Email, error)

	// CreateEmail creates a new email.
	CreateEmail(newEmail *emailCore.Email, ctx context.Context) (int64, *emailCore.Email, error)

	CreateProfileEmail(email_id int64, sender, recipient string, ctx context.Context) error

	// UpdateEmail updates the information of an email.
	UpdateEmail(updatedEmail *emailCore.Email, ctx context.Context) (bool, error)

	// DeleteEmail deletes the email.
	DeleteEmail(id uint64, login string, ctx context.Context) (bool, error)

	CheckRecipientEmail(recipient string, ctx context.Context) error
}
