package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"mail/internal/pkg/logger"
	"mail/internal/pkg/utils/constants"

	domain "mail/internal/microservice/models/domain_models"
	mock_repository "mail/internal/microservice/user/mock"
)

func GetCTX() context.Context {
	ctx := context.WithValue(context.Background(), constants.LoggerKey, logger.InitializationBdLog(nil))
	ctx2 := context.WithValue(ctx, constants.RequestIDKey, []string{"testID"})

	return ctx2
}

func TestGetAllUsers_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	expectedUsers := []*domain.User{
		{ID: 1, FirstName: "User 1"},
		{ID: 2, FirstName: "User 2"},
	}
	ctx := GetCTX()

	mockRepo.EXPECT().GetAll(0, 0, ctx).Return(expectedUsers, nil)

	users, err := useCase.GetAllUsers(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
}

func TestGetAllUsers_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	ctx := GetCTX()

	mockRepo.EXPECT().GetAll(0, 0, ctx).Return(nil, errors.New("repository error"))

	users, err := useCase.GetAllUsers(ctx)

	assert.Error(t, err)
	assert.Nil(t, users)
}

func TestGetUserByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	expectedUser := &domain.User{ID: 1, FirstName: "User 1"}
	ctx := GetCTX()

	mockRepo.EXPECT().GetByID(uint32(1), ctx).Return(expectedUser, nil)

	user, err := useCase.GetUserByID(1, ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
}

func TestGetUserByID_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	ctx := GetCTX()
	mockRepo.EXPECT().GetByID(uint32(1), ctx).Return(nil, errors.New("repository error"))

	user, err := useCase.GetUserByID(1, ctx)

	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestGetUserByLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	ctx := GetCTX()
	expectedUser := &domain.User{ID: 1, FirstName: "User 1"}
	mockRepo.EXPECT().GetUserByLogin("username", "password", ctx).Return(expectedUser, nil)

	user, err := useCase.GetUserByLogin("username", "password", ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
}

func TestGetUserByLogin_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	ctx := GetCTX()
	mockRepo.EXPECT().GetUserByLogin("username", "password", ctx).Return(nil, errors.New("repository error"))

	user, err := useCase.GetUserByLogin("username", "password", ctx)

	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestCreateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	ctx := GetCTX()
	newUser := &domain.User{FirstName: "New User", Login: "login", Password: "password"}
	mockRepo.EXPECT().Add(newUser, ctx).Return(newUser, nil)
	mockRepo.EXPECT().GetUserByLogin(newUser.Login, newUser.Password, ctx).Return(newUser, nil)
	mockRepo.EXPECT().InitAvatar(newUser.ID, "", "PHOTO", ctx).Return(true, nil)

	userRes, err := useCase.CreateUser(newUser, ctx)

	assert.NoError(t, err)
	assert.Equal(t, newUser, userRes)
}

func TestCreateUser_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)
	ctx := GetCTX()
	newUser := &domain.User{FirstName: "New User", Login: "login"}
	mockRepo.EXPECT().Add(newUser, ctx).Return(nil, errors.New("repository error"))

	userRes, err := useCase.CreateUser(newUser, ctx)

	assert.Error(t, err)
	assert.Nil(t, userRes)
}

func TestCreateUser_ErrorGetUserByLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)
	ctx := GetCTX()
	newUser := &domain.User{FirstName: "New User", Login: "login", Password: "password"}
	mockRepo.EXPECT().Add(newUser, ctx).Return(newUser, nil)
	mockRepo.EXPECT().GetUserByLogin(newUser.Login, newUser.Password, ctx).Return(nil, errors.New("repository error"))

	userRes, err := useCase.CreateUser(newUser, ctx)

	assert.Error(t, err)
	assert.Nil(t, userRes)
}

func TestCreateUser_ErrorInitUserAvatar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)
	ctx := GetCTX()
	newUser := &domain.User{FirstName: "New User", Login: "login", Password: "password"}
	mockRepo.EXPECT().Add(newUser, ctx).Return(newUser, nil)
	mockRepo.EXPECT().GetUserByLogin(newUser.Login, newUser.Password, ctx).Return(newUser, nil)
	mockRepo.EXPECT().InitAvatar(newUser.ID, "", "PHOTO", ctx).Return(false, errors.New("repository error"))

	userRes, err := useCase.CreateUser(newUser, ctx)

	assert.Error(t, err)
	assert.Nil(t, userRes)
}

func TestIsLoginUnique_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	ctx := GetCTX()
	login := "testUser"
	mockUsers := []*domain.User{
		{ID: 1, Login: "user1"},
		{ID: 2, Login: "user2"},
		{ID: 3, Login: "user3"},
	}

	mockRepo.EXPECT().GetAll(0, 0, ctx).Return(mockUsers, nil)
	unique, err := useCase.IsLoginUnique(login, ctx)
	assert.NoError(t, err)
	assert.True(t, unique)
}

func TestIsLoginUnique_NonUnique(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	ctx := GetCTX()
	nonUniqueLogin := "user2"
	mockUsers := []*domain.User{
		{ID: 1, Login: "user1"},
		{ID: 2, Login: nonUniqueLogin},
		{ID: 3, Login: "user3"},
	}

	mockRepo.EXPECT().GetAll(0, 0, ctx).Return(mockUsers, nil)
	unique, err := useCase.IsLoginUnique(nonUniqueLogin, ctx)
	assert.NoError(t, err)
	assert.False(t, unique)
}

func TestIsLoginUnique_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	ctx := GetCTX()
	mockRepo.EXPECT().GetAll(0, 0, ctx).Return(nil, errors.New("repository error"))
	unique, err := useCase.IsLoginUnique("testUser", ctx)
	assert.Error(t, err)
	assert.False(t, unique)
}

func TestUpdateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	ctx := GetCTX()
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

	mockRepo.EXPECT().GetByID(userNew.ID, ctx).Return(userOld, nil)

	mockRepo.EXPECT().Update(gomock.Any(), ctx).Return(true, nil)

	updatedUser, err := useCase.UpdateUser(userNew, ctx)

	assert.NoError(t, err)
	assert.Equal(t, userNew, updatedUser)
}

func TestUpdateUser_FailureToUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	ctx := GetCTX()
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

	mockRepo.EXPECT().GetByID(userNew.ID, ctx).Return(userOld, nil)

	mockRepo.EXPECT().Update(gomock.Any(), ctx).Return(false, errors.New("update failed"))

	_, err := useCase.UpdateUser(userNew, ctx)

	assert.Error(t, err)
}

func TestUpdateUser_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	ctx := GetCTX()
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

	mockRepo.EXPECT().GetByID(userNew.ID, ctx).Return(&domain.User{}, errors.New("repository error"))

	_, err := useCase.UpdateUser(userNew, ctx)

	assert.Error(t, err)
}

func TestDeleteUserByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	ctx := GetCTX()
	userID := uint32(1)
	mockRepo.EXPECT().Delete(userID, ctx).Return(true, nil)

	deleted, err := useCase.DeleteUserByID(userID, ctx)

	assert.NoError(t, err)
	assert.True(t, deleted)
}

func TestDeleteUserByID_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	ctx := GetCTX()
	userID := uint32(1)
	mockRepo.EXPECT().Delete(userID, ctx).Return(false, nil)

	deleted, err := useCase.DeleteUserByID(userID, ctx)

	assert.NoError(t, err)
	assert.False(t, deleted)
}

func TestDeleteUserByID_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	ctx := GetCTX()
	userID := uint32(1)
	mockRepo.EXPECT().Delete(userID, ctx).Return(false, errors.New("repository error"))

	deleted, err := useCase.DeleteUserByID(userID, ctx)

	assert.Error(t, err)
	assert.False(t, deleted)
}

func TestAddAvatar_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	ctx := GetCTX()
	userID := uint32(1)
	mockRepo.EXPECT().AddAvatar(userID, "", "PHOTO", ctx).Return(true, nil)

	added, err := useCase.AddAvatar(userID, "", ctx)

	assert.NoError(t, err)
	assert.True(t, added)
}

func TestAddAvatar_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	ctx := GetCTX()
	userID := uint32(1)
	mockRepo.EXPECT().AddAvatar(userID, "", "PHOTO", ctx).Return(false, errors.New("repository error"))

	added, err := useCase.AddAvatar(userID, "", ctx)

	assert.Error(t, err)
	assert.False(t, added)
}

func TestDeleteAvatarByUserID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	ctx := GetCTX()
	userID := uint32(1)
	mockRepo.EXPECT().DeleteAvatarByUserID(userID, ctx).Return(nil)

	err := useCase.DeleteAvatarByUserID(userID, ctx)

	assert.NoError(t, err)
}

func TestDeleteAvatarByUserID_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	ctx := GetCTX()
	userID := uint32(1)
	mockRepo.EXPECT().DeleteAvatarByUserID(userID, ctx).Return(errors.New("repository error"))

	err := useCase.DeleteAvatarByUserID(userID, ctx)

	assert.Error(t, err)
}
