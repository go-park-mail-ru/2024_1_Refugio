package user

import (
	domain "mail/pkg/domain/models"
	"mail/pkg/domain/repository"
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
func (uc *UserUseCase) CreateUser(user *domain.User) (*domain.User, error) {
	return uc.repo.Add(user)
}

// IsLoginUnique checks if the provided login is unique among all users.
func (uh *UserUseCase) IsLoginUnique(login string) (bool, error) {
	users, err := uh.repo.GetAll(0, 0)
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
