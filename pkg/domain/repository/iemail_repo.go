package repository

import (
	domain "mail/pkg/domain/models"
)

// EmailRepository represents the interface for working with emails.
type EmailRepository interface {
	// GetAll returns all emails incoming from the storage.
	GetAllIncoming(login, requestID string, offset, limit int) ([]*domain.Email, error)

	// GetAll returns all emails sent from the storage.
	GetAllSent(login, requestID string, offset, limit int) ([]*domain.Email, error)

	// GetByID returns the email by its unique identifier.
	GetByID(id uint64, login, requestID string) (*domain.Email, error)

	// Add adds a new email to the storage and returns its assigned unique identifier.
	Add(email *domain.Email, requestID string) (int64, *domain.Email, error)

	// Add adds a new profile_email
	AddProfileEmail(email_id int64, sender, recipient, requestID string) error

	// Update updates the information of an email in the storage based on the provided new email.
	Update(newEmail *domain.Email, requestID string) (bool, error)

	// Delete removes the email from the storage by its unique identifier.
	Delete(id uint64, login, requestID string) (bool, error)

	FindEmail(login, requestID string) error
}
