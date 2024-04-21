package usecase

import (
	"context"
	"fmt"
	"mail/internal/microservice/models/domain_models"
	repository "mail/internal/microservice/user/interface"
	"strings"
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
	return uc.repo.Add(user, ctx)
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
	if strings.TrimSpace(userNew.Patronymic) != "" && userNew.Patronymic != userOld.Patronymic {
		userOld.Patronymic = userNew.Patronymic
	}
	if domain_models.IsValidGender(userNew.Gender) && userNew.Gender != userOld.Gender {
		userOld.Gender = userNew.Gender
	}
	if !userNew.Birthday.Equal(userOld.Birthday) {
		userOld.Birthday = userNew.Birthday
	}
	if strings.TrimSpace(userNew.Description) != "" && userNew.Description != userOld.Description {
		userOld.Description = userNew.Description
	}
	if strings.TrimSpace(userNew.PhoneNumber) != "" && userNew.PhoneNumber != userOld.PhoneNumber {
		userOld.PhoneNumber = userNew.PhoneNumber
	}
	if strings.TrimSpace(userNew.AvatarID) != "" && userNew.AvatarID != userOld.AvatarID {
		userOld.AvatarID = userNew.AvatarID
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
