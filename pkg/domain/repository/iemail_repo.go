package repository

import (
	emailCore "mail/pkg/domain/models"
)

// EmailRepository represents the interface for working with emails.
type EmailRepository interface {
	// GetAll returns all emails from the storage.
	GetAll() ([]*emailCore.Email, error)

	// GetByID returns the email by its unique identifier.
	GetByID(id uint64) (*emailCore.Email, error)

	// Add adds a new email to the storage and returns its assigned unique identifier.
	Add(email *emailCore.Email) (*emailCore.Email, error)

	// Update updates the information of an email in the storage based on the provided new email.
	Update(newEmail *emailCore.Email) (bool, error)

	// Delete removes the email from the storage by its unique identifier.
	Delete(id uint64) (bool, error)
}
