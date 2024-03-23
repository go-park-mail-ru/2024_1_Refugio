package usecase

import (
	userCore "mail/pkg/domain/models"
)

// UserUseCase represents the use case for working with users.
type UserUseCase interface {
	// GetAllUsers returns all users.
	GetAllUsers() ([]*userCore.User, error)

	// GetUserByID returns the user by its ID.
	GetUserByID(id uint32) (*userCore.User, error)

	// GetUserByLogin returns the user by login.
	GetUserByLogin(login string, password string) (*userCore.User, error)

	// CreateUser creates a new user.
	CreateUser(user *userCore.User) (uint32, error)

	// UpdateUser updates the user's information.
	UpdateUser(user *userCore.User) (bool, error)

	// DeleteUser deletes the user.
	DeleteUser(id uint32) (bool, error)
}
