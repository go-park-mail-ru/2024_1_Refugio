package repository

import (
	userCore "mail/pkg/domain/models"
)

// UserRepository represents the interface for working with users.
type UserRepository interface {
	// GetAll returns all users from the storage.
	GetAll() ([]*userCore.User, error)

	// GetByID returns the user by its unique identifier.
	GetByID(id uint32) (*userCore.User, error)

	// GetUserByLogin returns the user by login.
	GetUserByLogin(login string, password string) (*userCore.User, error)

	// Add adds a new user to the storage and returns its assigned unique identifier.
	Add(user *userCore.User) (uint32, error)

	// Update updates the information of a user in the storage based on the provided new user.
	Update(newUser *userCore.User) (bool, error)

	// Delete removes the user from the storage by its unique identifier.
	Delete(id uint32) (bool, error)
}
