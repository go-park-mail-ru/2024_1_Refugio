package user

import (
	"mail/pkg/domain/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) GetAll() ([]*models.User, error) {
	args := m.Called()
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *mockUserRepository) GetByID(id uint32) (*models.User, error) {
	args := m.Called(id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *mockUserRepository) GetUserByLogin(login string, password string) (*models.User, error) {
	args := m.Called(login, password)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *mockUserRepository) Add(user *models.User) (uint32, error) {
	args := m.Called(user)
	return uint32(args.Int(0)), args.Error(1)
}

func (m *mockUserRepository) Update(user *models.User) (bool, error) {
	args := m.Called(user)
	return args.Bool(0), args.Error(1)
}

func (m *mockUserRepository) Delete(id uint32) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}

func TestUserUseCase_GetAllUsers(t *testing.T) {
	repo := new(mockUserRepository)
	expectedUsers := []*models.User{{ID: 1, Name: "John"}, {ID: 2, Name: "Jane"}}
	repo.On("GetAll").Return(expectedUsers, nil)

	uc := NewUserUseCase(repo)
	users, err := uc.GetAllUsers()
	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
	repo.AssertExpectations(t)
}

func TestUserUseCase_GetUserByID(t *testing.T) {
	repo := new(mockUserRepository)
	expectedUser := &models.User{ID: 1, Name: "John"}
	repo.On("GetByID", uint32(1)).Return(expectedUser, nil)

	uc := NewUserUseCase(repo)
	foundUser, err := uc.GetUserByID(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, foundUser)
	repo.AssertExpectations(t)
}

func TestUserUseCase_GetUserByLogin(t *testing.T) {
	repo := new(mockUserRepository)
	expectedUser := &models.User{ID: 1, Name: "John"}
	repo.On("GetUserByLogin", "john", "password").Return(expectedUser, nil)

	uc := NewUserUseCase(repo)
	foundUser, err := uc.GetUserByLogin("john", "password")
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, foundUser)
	repo.AssertExpectations(t)
}

func TestUserUseCase_CreateUser(t *testing.T) {
	repo := new(mockUserRepository)
	newUser := &models.User{Name: "John"}
	repo.On("Add", newUser).Return(1, nil)

	uc := NewUserUseCase(repo)
	id, err := uc.CreateUser(newUser)
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), id)
	repo.AssertExpectations(t)
}
