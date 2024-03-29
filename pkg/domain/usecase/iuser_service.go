//go:generate mockgen -source=./iuser_service.go -destination=../mock/user_service_mock.go -package=mock

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
	CreateUser(user *userCore.User) (*userCore.User, error)

	// IsLoginUnique checks if the provided login is unique among all users.
	IsLoginUnique(login string) (bool, error)
}
