//go:generate mockgen -source=./iuser_service.go -destination=../mock/user_service_mock.go -package=mock

package _interface

import (
	"context"
	domain "mail/internal/microservice/models/domain_models"
)

// UserUseCase represents the use case for working with users.
type UserUseCase interface {
	// GetAllUsers returns all users.
	GetAllUsers(ctx context.Context) ([]*domain.User, error)

	// GetUserByID returns the user by its ID.
	GetUserByID(id uint32, ctx context.Context) (*domain.User, error)

	// GetUserByLogin returns the user by login.
	GetUserByLogin(login, password string, ctx context.Context) (*domain.User, error)

	// CreateUser creates a new user.
	CreateUser(user *domain.User, ctx context.Context) (*domain.User, error)

	// IsLoginUnique checks if the provided login is unique among all users.
	IsLoginUnique(login string, ctx context.Context) (bool, error)

	// UpdateUser updates user data based on the provided ID.
	UpdateUser(userNew *domain.User, ctx context.Context) (*domain.User, error)

	// DeleteUserAvatar updates user avatar based on the provided ID.
	DeleteUserAvatar(userNew *domain.User, ctx context.Context) (*domain.User, error)

	// DeleteUserByID deletes the user with the given ID.
	DeleteUserByID(id uint32, ctx context.Context) (bool, error)
}