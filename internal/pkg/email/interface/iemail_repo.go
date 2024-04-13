package _interface

import (
	"context"
	domain "mail/internal/models/domain_models"
)

// EmailRepository represents the interface for working with emails.
type EmailRepository interface {
	// GetAll returns all emails incoming from the storage.
	GetAllIncoming(login string, offset, limit int, ctx context.Context) ([]*domain.Email, error)

	// GetAll returns all emails sent from the storage.
	GetAllSent(login string, offset, limit int, ctx context.Context) ([]*domain.Email, error)

	// GetByID returns the email by its unique identifier.
	GetByID(id uint64, login string, ctx context.Context) (*domain.Email, error)

	// Add adds a new email to the storage and returns its assigned unique identifier.
	Add(email *domain.Email, ctx context.Context) (int64, *domain.Email, error)

	// Add adds a new profile_email
	AddProfileEmail(email_id int64, sender, recipient string, ctx context.Context) error

	// Update updates the information of an email in the storage based on the provided new email.
	Update(newEmail *domain.Email, ctx context.Context) (bool, error)

	// Delete removes the email from the storage by its unique identifier.
	Delete(id uint64, login string, ctx context.Context) (bool, error)

	FindEmail(login string, ctx context.Context) error
}
