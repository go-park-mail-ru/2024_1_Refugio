//go:generate mockgen -source=./iuser_repo.go -destination=../mock/user_repository_mock.go -package=mock

package repository

import (
	domain "mail/pkg/domain/models"
)

// UserRepository represents the interface for working with users.
type UserRepository interface {
	// GetAll returns all users from the storage.
	GetAll(offset, limit int) ([]*domain.User, error)

	// GetByID returns the user by its unique identifier.
	GetByID(id uint32) (*domain.User, error)

	// GetUserByLogin returns the user by login.
	GetUserByLogin(login string, password string) (*domain.User, error)

	// Add adds a new user to the storage and returns its assigned unique identifier.
	Add(user *domain.User) (uint32, error)

	// Update updates the information of a user in the storage based on the provided new user.
	Update(newUser *domain.User) (bool, error)

	// Delete removes the user from the storage by its unique identifier.
	Delete(id uint32) (bool, error)
}
