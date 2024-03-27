package user

import (
	"fmt"
	domain "mail/pkg/domain/models"
	"mail/pkg/domain/repository"
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
func (uc *UserUseCase) GetAllUsers() ([]*domain.User, error) {
	return uc.repo.GetAll(0, 0)
}

// GetUserByID returns the user by its ID.
func (uc *UserUseCase) GetUserByID(id uint32) (*domain.User, error) {
	return uc.repo.GetByID(id)
}

// GetUserByLogin returns the user by login.
func (uc *UserUseCase) GetUserByLogin(login string, password string) (*domain.User, error) {
	return uc.repo.GetUserByLogin(login, password)
}

// CreateUser creates a new user.
func (uc *UserUseCase) CreateUser(user *domain.User) (uint32, error) {
	return uc.repo.Add(user)
}

// IsLoginUnique checks if the provided login is unique among all users.
func (uc *UserUseCase) IsLoginUnique(login string) (bool, error) {
	users, err := uc.repo.GetAll(0, 0)
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

// UpdateUserById updates user data based on the provided ID.
func (uc *UserUseCase) UpdateUser(userNew *domain.User) (*domain.User, error) {
	userOld, err := uc.repo.GetByID(userNew.ID)
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
	if domain.IsValidGender(userNew.Gender) && userNew.Gender != userOld.Gender {
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

	status, err := uc.repo.Update(userOld)
	if err != nil {
		return nil, err
	}
	if !status {
		return nil, fmt.Errorf("failed to update user")
	}

	return userOld, nil
}

// DeleteUserByID deletes the user with the given ID.
func (uc *UserUseCase) DeleteUserByID(id uint32) (bool, error) {
	return uc.repo.Delete(id)
}
