//go:generate mockgen -source=./iuser_repo.go -destination=../mock/user_repository_mock.go -package=mock

package _interface

import (
	domain "mail/internal/models/domain_models"
)

// UserRepository represents the interface for working with users.
type UserRepository interface {
	// GetAll returns all users from the storage.
	GetAll(offset, limit int, requestID string) ([]*domain.User, error)

	// GetByID returns the user by its unique identifier.
	GetByID(id uint32, requestID string) (*domain.User, error)

	// GetUserByLogin returns the user by login.
	GetUserByLogin(login, password, requestID string) (*domain.User, error)

	// Add adds a new user to the storage and returns its assigned unique identifier.
	Add(user *domain.User, requestID string) (*domain.User, error)

	// Update updates the information of a user in the storage based on the provided new user.
	Update(newUser *domain.User, requestID string) (bool, error)

	// Delete removes the user from the storage by its unique identifier.
	Delete(id uint32, requestID string) (bool, error)
}
