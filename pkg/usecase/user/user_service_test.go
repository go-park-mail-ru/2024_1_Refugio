package user

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mock_repository "mail/pkg/domain/mock"
	domain "mail/pkg/domain/models"
	"testing"
	"time"
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

func TestIsLoginUnique_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	login := "testUser"
	mockUsers := []*domain.User{
		{ID: 1, Login: "user1"},
		{ID: 2, Login: "user2"},
		{ID: 3, Login: "user3"},
	}

	mockRepo.EXPECT().GetAll(0, 0).Return(mockUsers, nil)
	unique, err := useCase.IsLoginUnique(login)
	assert.NoError(t, err)
	assert.True(t, unique)
}

func TestIsLoginUnique_NonUnique(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	nonUniqueLogin := "user2"
	mockUsers := []*domain.User{
		{ID: 1, Login: "user1"},
		{ID: 2, Login: nonUniqueLogin},
		{ID: 3, Login: "user3"},
	}

	mockRepo.EXPECT().GetAll(0, 0).Return(mockUsers, nil)
	unique, err := useCase.IsLoginUnique(nonUniqueLogin)
	assert.NoError(t, err)
	assert.False(t, unique)
}

func TestIsLoginUnique_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	mockRepo.EXPECT().GetAll(0, 0).Return(nil, errors.New("repository error"))
	unique, err := useCase.IsLoginUnique("testUser")
	assert.Error(t, err)
	assert.False(t, unique)
}

func TestUpdateUserById_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	userNew := &domain.User{
		ID:          1,
		FirstName:   "John",
		Surname:     "Smith",
		Patronymic:  "William",
		Gender:      "Male",
		Birthday:    time.Date(1985, time.October, 20, 0, 0, 0, 0, time.UTC),
		Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
		PhoneNumber: "+1234567890",
	}

	userOld := &domain.User{
		ID:          1,
		FirstName:   "Doe",
		Surname:     "Johnson",
		Patronymic:  "Michael",
		Gender:      "Male",
		Birthday:    time.Date(1980, time.January, 15, 0, 0, 0, 0, time.UTC),
		Description: "Suspendisse potenti. Nulla facilisi.",
		PhoneNumber: "+9876543210",
	}

	mockRepo.EXPECT().GetByID(userNew.ID).Return(userOld, nil)

	mockRepo.EXPECT().Update(gomock.Any()).Return(true, nil)

	updatedUser, err := useCase.UpdateUser(userNew)

	assert.NoError(t, err)

	assert.Equal(t, userNew, updatedUser)
}

func TestUpdateUserById_FailureToUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	userNew := &domain.User{
		ID:          1,
		FirstName:   "John",
		Surname:     "Smith",
		Patronymic:  "William",
		Gender:      "Male",
		Birthday:    time.Date(1985, time.October, 20, 0, 0, 0, 0, time.UTC),
		Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
		PhoneNumber: "+1234567890",
	}

	userOld := &domain.User{
		ID:          1,
		FirstName:   "Doe",
		Surname:     "Johnson",
		Patronymic:  "Michael",
		Gender:      "Male",
		Birthday:    time.Date(1980, time.January, 15, 0, 0, 0, 0, time.UTC),
		Description: "Suspendisse potenti. Nulla facilisi.",
		PhoneNumber: "+9876543210",
	}

	mockRepo.EXPECT().GetByID(userNew.ID).Return(userOld, nil)

	mockRepo.EXPECT().Update(gomock.Any()).Return(false, errors.New("update failed"))

	_, err := useCase.UpdateUser(userNew)

	assert.Error(t, err)
}

func TestUpdateUserById_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	userNew := &domain.User{
		ID:          1,
		FirstName:   "John",
		Surname:     "Smith",
		Patronymic:  "William",
		Gender:      "Male",
		Birthday:    time.Date(1985, time.October, 20, 0, 0, 0, 0, time.UTC),
		Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
		PhoneNumber: "+1234567890",
	}

	mockRepo.EXPECT().GetByID(userNew.ID).Return(&domain.User{}, errors.New("repository error"))

	_, err := useCase.UpdateUser(userNew)

	assert.Error(t, err)
}

func TestDeleteUserByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	userID := uint32(1)
	mockRepo.EXPECT().Delete(userID).Return(true, nil)

	deleted, err := useCase.DeleteUserByID(userID)

	assert.NoError(t, err)
	assert.True(t, deleted)
}

func TestDeleteUserByID_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	userID := uint32(1)
	mockRepo.EXPECT().Delete(userID).Return(false, nil)

	deleted, err := useCase.DeleteUserByID(userID)

	assert.NoError(t, err)
	assert.False(t, deleted)
}

func TestDeleteUserByID_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	userID := uint32(1)
	mockRepo.EXPECT().Delete(userID).Return(false, errors.New("repository error"))

	deleted, err := useCase.DeleteUserByID(userID)

	assert.Error(t, err)
	assert.False(t, deleted)
}
