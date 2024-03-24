package user

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mock_repository "mail/pkg/domain/mock"
	domain "mail/pkg/domain/models"
	"testing"
)

func TestGetAllUsers_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	expectedUsers := []*domain.User{
		{ID: 1, FirstName: "User 1"},
		{ID: 2, FirstName: "User 2"},
	}
	mockRepo.EXPECT().GetAll(0, 0).Return(expectedUsers, nil)

	users, err := useCase.GetAllUsers()

	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
}

func TestGetAllUsers_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	mockRepo.EXPECT().GetAll(0, 0).Return(nil, errors.New("repository error"))

	users, err := useCase.GetAllUsers()

	assert.Error(t, err)
	assert.Nil(t, users)
}

func TestGetUserByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	expectedUser := &domain.User{ID: 1, FirstName: "User 1"}
	mockRepo.EXPECT().GetByID(uint32(1)).Return(expectedUser, nil)

	user, err := useCase.GetUserByID(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
}

func TestGetUserByID_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	mockRepo.EXPECT().GetByID(uint32(1)).Return(nil, errors.New("repository error"))

	user, err := useCase.GetUserByID(1)

	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestGetUserByLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	expectedUser := &domain.User{ID: 1, FirstName: "User 1"}
	mockRepo.EXPECT().GetUserByLogin("username", "password").Return(expectedUser, nil)

	user, err := useCase.GetUserByLogin("username", "password")

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
}

func TestGetUserByLogin_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	mockRepo.EXPECT().GetUserByLogin("username", "password").Return(nil, errors.New("repository error"))

	user, err := useCase.GetUserByLogin("username", "password")

	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestCreateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	expectedUserID := uint32(1)
	mockRepo.EXPECT().Add(gomock.Any()).Return(expectedUserID, nil)

	newUser := &domain.User{FirstName: "New User"}

	userID, err := useCase.CreateUser(newUser)

	assert.NoError(t, err)
	assert.Equal(t, expectedUserID, userID)
}

func TestCreateUser_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	mockRepo.EXPECT().Add(gomock.Any()).Return(uint32(0), errors.New("repository error"))

	newUser := &domain.User{FirstName: "New User"}

	userID, err := useCase.CreateUser(newUser)

	assert.Error(t, err)
	assert.Zero(t, userID)
}
