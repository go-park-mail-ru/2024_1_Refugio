package usecase

import (
	"context"
	"fmt"
	"strings"

	"mail/internal/microservice/models/domain_models"

	repository "mail/internal/microservice/user/interface"
)

// UserUseCase represents the use case for working with users.
type UserUseCase struct {
	repo repository.UserRepository
}

// NewUserUseCase creates a new instance of UserUseCase.
func NewUserUseCase(repo repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		repo: repo,
	}
}

// GetAllUsers returns all users.
func (uc *UserUseCase) GetAllUsers(ctx context.Context) ([]*domain_models.User, error) {
	return uc.repo.GetAll(0, 0, ctx)
}

// GetUserByID returns the user by its ID.
func (uc *UserUseCase) GetUserByID(id uint32, ctx context.Context) (*domain_models.User, error) {
	return uc.repo.GetByID(id, ctx)
}

// GetUserByLogin returns the user by login.
func (uc *UserUseCase) GetUserByLogin(login, password string, ctx context.Context) (*domain_models.User, error) {
	return uc.repo.GetUserByLogin(login, password, ctx)
}

// CreateUser creates a new user.
func (uc *UserUseCase) CreateUser(user *domain_models.User, ctx context.Context) (*domain_models.User, error) {
	_, err := uc.repo.Add(user, ctx)
	if err != nil {
		return nil, fmt.Errorf("user with login %s not create", user.Login)
	}

	userByLogin, err := uc.repo.GetUserByLogin(user.Login, user.Password, ctx)
	if err != nil {
		return nil, fmt.Errorf("user with login %s not found", user.Login)
	}

	_, errAva := uc.repo.InitAvatar(userByLogin.ID, "", "PHOTO", ctx)
	if errAva != nil {
		return nil, fmt.Errorf("user avatar fail")
	}

	return userByLogin, nil
}

// IsLoginUnique checks if the provided login is unique among all users.
func (uc *UserUseCase) IsLoginUnique(login string, ctx context.Context) (bool, error) {
	users, err := uc.repo.GetAll(0, 0, ctx)
	if err != nil {
		return false, err
	}

	for _, u := range users {
		if u.Login == login {
			return false, nil
		}
	}

	return true, nil
}

// UpdateUser updates user data based on the provided ID.
func (uc *UserUseCase) UpdateUser(userNew *domain_models.User, ctx context.Context) (*domain_models.User, error) {
	userOld, err := uc.repo.GetByID(userNew.ID, ctx)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(userNew.FirstName) != "" && userNew.FirstName != userOld.FirstName {
		userOld.FirstName = userNew.FirstName
	}
	if strings.TrimSpace(userNew.Surname) != "" && userNew.Surname != userOld.Surname {
		userOld.Surname = userNew.Surname
	}
	if userNew.Patronymic != userOld.Patronymic {
		userOld.Patronymic = userNew.Patronymic
	}
	if domain_models.IsValidGender(userNew.Gender) && userNew.Gender != userOld.Gender {
		userOld.Gender = userNew.Gender
	}
	if !userNew.Birthday.Equal(userOld.Birthday) {
		userOld.Birthday = userNew.Birthday
	}
	if userNew.Description != userOld.Description {
		userOld.Description = userNew.Description
	}
	if userNew.PhoneNumber != userOld.PhoneNumber {
		userOld.PhoneNumber = userNew.PhoneNumber
	}

	status, err := uc.repo.Update(userOld, ctx)
	if err != nil {
		return nil, err
	}
	if !status {
		return nil, fmt.Errorf("failed to update user")
	}

	return userOld, nil
}

// DeleteUserByID deletes the user with the given ID.
func (uc *UserUseCase) DeleteUserByID(id uint32, ctx context.Context) (bool, error) {
	return uc.repo.Delete(id, ctx)
}

// AddAvatar adds a new user avatar.
func (uc *UserUseCase) AddAvatar(id uint32, fileID string, ctx context.Context) (bool, error) {
	return uc.repo.AddAvatar(id, fileID, "PHOTO", ctx)
}

// DeleteAvatarByUserID deletes a user's photo.
func (uc *UserUseCase) DeleteAvatarByUserID(userID uint32, ctx context.Context) error {
	return uc.repo.DeleteAvatarByUserID(userID, ctx)
}

// GetUserVkID get user by VK Id.
func (uc *UserUseCase) GetUserVkID(vkId uint32, ctx context.Context) (*domain_models.User, error) {
	return uc.repo.GetByVKID(vkId, ctx)
}

// GetUserByOnlyLogin get user by login.
func (uc *UserUseCase) GetUserByOnlyLogin(login string, ctx context.Context) (*domain_models.User, error) {
	return uc.repo.GetByOnlyLogin(login, ctx)
}
