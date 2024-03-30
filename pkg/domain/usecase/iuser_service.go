//go:generate mockgen -source=./iuser_service.go -destination=../mock/user_service_mock.go -package=mock

package usecase

import (
	domain "mail/pkg/domain/models"
)

// UserUseCase represents the use case for working with users.
type UserUseCase interface {
	// GetAllUsers returns all users.
	GetAllUsers() ([]*domain.User, error)

	// GetUserByID returns the user by its ID.
	GetUserByID(id uint32) (*domain.User, error)

	// GetUserByLogin returns the user by login.
	GetUserByLogin(login string, password string) (*domain.User, error)

	// CreateUser creates a new user.
	CreateUser(user *domain.User) (*domain.User, error)

	// IsLoginUnique checks if the provided login is unique among all users.
	IsLoginUnique(login string) (bool, error)

	// UpdateUser updates user data based on the provided ID.
	UpdateUser(userNew *domain.User) (*domain.User, error)

	// DeleteUserByID deletes the user with the given ID.
	DeleteUserByID(id uint32) (bool, error)
}
