package user

import (
	userCore "mail/pkg/domain/models"
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
func (uc *UserUseCase) GetAllUsers() ([]*userCore.User, error) {
	return uc.repo.GetAll()
}

// GetUserByID returns the user by its ID.
func (uc *UserUseCase) GetUserByID(id uint32) (*userCore.User, error) {
	return uc.repo.GetByID(id)
}

// GetUserByLogin returns the user by login.
func (uc *UserUseCase) GetUserByLogin(login string, password string) (*userCore.User, error) {
	return uc.repo.GetUserByLogin(login, password)
}

// CreateUser creates a new user.
func (uc *UserUseCase) CreateUser(user *userCore.User) (uint32, error) {
	return uc.repo.Add(user)
}

// UpdateUser updates the user's information.
func (uc *UserUseCase) UpdateUser(user *userCore.User) (bool, error) {
	return uc.repo.Update(user)
}

// DeleteUser deletes the user.
func (uc *UserUseCase) DeleteUser(id uint32) (bool, error) {
	return uc.repo.Delete(id)
}
