//go:generate mockgen -source=./iuser_repo.go -destination=../mock/user_repository_mock.go -package=mock

package _interface

import (
	"context"

	domain "mail/internal/microservice/models/domain_models"
)

// UserRepository represents the interface for working with users.
type UserRepository interface {
	// GetAll returns all users from the storage.
	GetAll(offset, limit int, ctx context.Context) ([]*domain.User, error)

	// GetByID returns the user by its unique identifier.
	GetByID(id uint32, ctx context.Context) (*domain.User, error)

	// GetUserByLogin returns the user by login.
	GetUserByLogin(login, password string, ctx context.Context) (*domain.User, error)

	// Add adds a new user to the storage and returns its assigned unique identifier.
	Add(user *domain.User, ctx context.Context) (*domain.User, error)

	// Update updates the information of a user in the storage based on the provided new user.
	Update(newUser *domain.User, ctx context.Context) (bool, error)

	// Delete removes the user from the storage by its unique identifier.
	Delete(id uint32, ctx context.Context) (bool, error)

	// AddAvatar adds a new user avatar to the repository and associates it with the profile.
	AddAvatar(id uint32, fileID, fileType string, ctx context.Context) (bool, error)

	// DeleteAvatarByUserID deletes a user's photo and an entry from the file table by its ID in one request.
	DeleteAvatarByUserID(userID uint32, ctx context.Context) error

	// InitAvatar initializes the user's avatar by updating the corresponding entry in the database.
	InitAvatar(id uint32, fileID, fileType string, ctx context.Context) (bool, error)

	// GetByVKID returns the user by its unique identifier.
	GetByVKID(vkId uint32, ctx context.Context) (*domain.User, error)

	// GetByOnlyLogin returns the user by its unique identifier.
	GetByOnlyLogin(login string, ctx context.Context) (*domain.User, error)
}
