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

	requestID := "test_request"
	expectedUsers := []*domain.User{
		{ID: 1, FirstName: "User 1"},
		{ID: 2, FirstName: "User 2"},
	}
	mockRepo.EXPECT().GetAll(0, 0, requestID).Return(expectedUsers, nil)

	users, err := useCase.GetAllUsers(requestID)

	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
}

func TestGetAllUsers_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	requestID := "test_request"

	mockRepo.EXPECT().GetAll(0, 0, requestID).Return(nil, errors.New("repository error"))

	users, err := useCase.GetAllUsers(requestID)

	assert.Error(t, err)
	assert.Nil(t, users)
}

func TestGetUserByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	requestID := "test_request"
	expectedUser := &domain.User{ID: 1, FirstName: "User 1"}
	mockRepo.EXPECT().GetByID(uint32(1), requestID).Return(expectedUser, nil)

	user, err := useCase.GetUserByID(1, requestID)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
}

func TestGetUserByID_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	requestID := "test_request"
	mockRepo.EXPECT().GetByID(uint32(1), requestID).Return(nil, errors.New("repository error"))

	user, err := useCase.GetUserByID(1, requestID)

	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestGetUserByLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	requestID := "test_request"
	expectedUser := &domain.User{ID: 1, FirstName: "User 1"}
	mockRepo.EXPECT().GetUserByLogin("username", "password", requestID).Return(expectedUser, nil)

	user, err := useCase.GetUserByLogin("username", "password", requestID)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
}

func TestGetUserByLogin_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	requestID := "test_request"
	mockRepo.EXPECT().GetUserByLogin("username", "password", requestID).Return(nil, errors.New("repository error"))

	user, err := useCase.GetUserByLogin("username", "password", requestID)

	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestCreateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	requestID := "test_request"
	newUser := &domain.User{FirstName: "New User"}
	mockRepo.EXPECT().Add(newUser, requestID).Return(newUser, nil)

	userRes, err := useCase.CreateUser(newUser, requestID)

	assert.NoError(t, err)
	assert.Equal(t, newUser, userRes)
}

func TestCreateUser_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)
	requestID := "test_request"
	newUser := &domain.User{FirstName: "New User"}
	mockRepo.EXPECT().Add(newUser, requestID).Return(newUser, errors.New("repository error"))

	userRes, err := useCase.CreateUser(newUser, requestID)

	assert.Error(t, err)
	assert.Equal(t, newUser, userRes)
}

func TestIsLoginUnique_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	requestID := "test_request"
	login := "testUser"
	mockUsers := []*domain.User{
		{ID: 1, Login: "user1"},
		{ID: 2, Login: "user2"},
		{ID: 3, Login: "user3"},
	}

	mockRepo.EXPECT().GetAll(0, 0, requestID).Return(mockUsers, nil)
	unique, err := useCase.IsLoginUnique(login, requestID)
	assert.NoError(t, err)
	assert.True(t, unique)
}

func TestIsLoginUnique_NonUnique(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	requestID := "test_request"
	nonUniqueLogin := "user2"
	mockUsers := []*domain.User{
		{ID: 1, Login: "user1"},
		{ID: 2, Login: nonUniqueLogin},
		{ID: 3, Login: "user3"},
	}

	mockRepo.EXPECT().GetAll(0, 0, requestID).Return(mockUsers, nil)
	unique, err := useCase.IsLoginUnique(nonUniqueLogin, requestID)
	assert.NoError(t, err)
	assert.False(t, unique)
}

func TestIsLoginUnique_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	requestID := "test_request"
	mockRepo.EXPECT().GetAll(0, 0, requestID).Return(nil, errors.New("repository error"))
	unique, err := useCase.IsLoginUnique("testUser", requestID)
	assert.Error(t, err)
	assert.False(t, unique)
}

func TestUpdateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	requestID := "test_request"
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

	mockRepo.EXPECT().GetByID(userNew.ID, requestID).Return(userOld, nil)

	mockRepo.EXPECT().Update(gomock.Any(), requestID).Return(true, nil)

	updatedUser, err := useCase.UpdateUser(userNew, requestID)

	assert.NoError(t, err)
	assert.Equal(t, userNew, updatedUser)
}

func TestUpdateUser_FailureToUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	requestID := "test_request"
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

	mockRepo.EXPECT().GetByID(userNew.ID, requestID).Return(userOld, nil)

	mockRepo.EXPECT().Update(gomock.Any(), requestID).Return(false, errors.New("update failed"))

	_, err := useCase.UpdateUser(userNew, requestID)

	assert.Error(t, err)
}

func TestUpdateUser_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	requestID := "test_request"
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

	mockRepo.EXPECT().GetByID(userNew.ID, requestID).Return(&domain.User{}, errors.New("repository error"))

	_, err := useCase.UpdateUser(userNew, requestID)

	assert.Error(t, err)
}

func TestDeleteUserByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	requestID := "test_request"
	userID := uint32(1)
	mockRepo.EXPECT().Delete(userID, requestID).Return(true, nil)

	deleted, err := useCase.DeleteUserByID(userID, requestID)

	assert.NoError(t, err)
	assert.True(t, deleted)
}

func TestDeleteUserByID_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	requestID := "test_request"
	userID := uint32(1)
	mockRepo.EXPECT().Delete(userID, requestID).Return(false, nil)

	deleted, err := useCase.DeleteUserByID(userID, requestID)

	assert.NoError(t, err)
	assert.False(t, deleted)
}

func TestDeleteUserByID_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	requestID := "test_request"
	userID := uint32(1)
	mockRepo.EXPECT().Delete(userID, requestID).Return(false, errors.New("repository error"))

	deleted, err := useCase.DeleteUserByID(userID, requestID)

	assert.Error(t, err)
	assert.False(t, deleted)
}
