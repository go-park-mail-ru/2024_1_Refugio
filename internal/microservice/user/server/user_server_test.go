package server

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"mail/internal/microservice/models/domain_models"
	"mail/internal/microservice/user/mock"
	"mail/internal/microservice/user/proto"
	"mail/internal/pkg/logger"
	"mail/internal/pkg/utils/constants"
)

func GetCTX() context.Context {
	ctx := context.WithValue(context.Background(), constants.LoggerKey, logger.InitializationBdLog(nil))
	ctx2 := context.WithValue(ctx, constants.RequestIDKey, []string{"testID"})

	return ctx2
}

func TestGetUsers_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)

	server := NewUserServer(mockUserUseCase)

	ctx := GetCTX()

	expectedUsers := []*domain_models.User{
		{ID: 1, FirstName: "User1"},
		{ID: 2, FirstName: "User2"},
	}

	mockUserUseCase.EXPECT().GetAllUsers(ctx).Return(expectedUsers, nil)

	reply, err := server.GetUsers(ctx, &proto.GetUsersRequest{})

	assert.NoError(t, err)
	assert.NotNil(t, reply)
}

func TestGetUsers_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)

	server := NewUserServer(mockUserUseCase)

	ctx := GetCTX()

	mockUserUseCase.EXPECT().GetAllUsers(ctx).Return(nil, fmt.Errorf("user not found"))

	reply, err := server.GetUsers(ctx, &proto.GetUsersRequest{})

	assert.Error(t, err)
	assert.Nil(t, reply)
}

func TestGetUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)

	server := NewUserServer(mockUserUseCase)

	ctx := GetCTX()

	expectedUser := &domain_models.User{ID: 1, FirstName: "John"}

	mockUserUseCase.EXPECT().GetUserByID(expectedUser.ID, ctx).Return(expectedUser, nil)

	reply, err := server.GetUser(ctx, &proto.GetUserRequest{Id: expectedUser.ID})

	assert.NoError(t, err)
	assert.NotNil(t, reply)
}

func TestGetUser_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)

	server := NewUserServer(mockUserUseCase)

	ctx := GetCTX()

	reply, err := server.GetUser(ctx, &proto.GetUserRequest{Id: 0})

	assert.Error(t, err)
	assert.Nil(t, reply)
}

func TestGetUserByLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)

	server := NewUserServer(mockUserUseCase)

	ctx := GetCTX()

	expectedUser := &domain_models.User{ID: 1, FirstName: "John"}

	mockUserUseCase.EXPECT().GetUserByLogin("john_doe", "password", ctx).Return(expectedUser, nil)

	reply, err := server.GetUserByLogin(ctx, &proto.GetUserByLoginRequest{Login: "john_doe", Password: "password"})

	assert.NoError(t, err)
	assert.NotNil(t, reply)
}

func TestGetUserByLogin_InvalidLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)

	server := NewUserServer(mockUserUseCase)

	ctx := GetCTX()

	reply, err := server.GetUserByLogin(ctx, &proto.GetUserByLoginRequest{})

	assert.Error(t, err)
	assert.Nil(t, reply)
}

func TestIsLoginUnique_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)

	server := NewUserServer(mockUserUseCase)

	ctx := GetCTX()

	mockUserUseCase.EXPECT().IsLoginUnique("john_doe", ctx).Return(true, nil)

	reply, err := server.IsLoginUnique(ctx, &proto.IsLoginUniqueRequest{Login: "john_doe"})

	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.True(t, reply.Status)
}

func TestIsLoginUnique_InvalidLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)

	server := NewUserServer(mockUserUseCase)

	ctx := GetCTX()

	reply, err := server.IsLoginUnique(ctx, &proto.IsLoginUniqueRequest{})

	assert.Error(t, err)
	assert.Nil(t, reply)
}

func TestUpdateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)

	server := NewUserServer(mockUserUseCase)

	ctx := context.Background()

	mockUserUseCase.EXPECT().UpdateUser(gomock.Any(), ctx).Return(&domain_models.User{ID: 1, Login: "john_doe", FirstName: "John", Surname: "Doe", Gender: "male"}, nil)

	reply, err := server.UpdateUser(ctx, &proto.UpdateUserRequest{
		User: &proto.User{
			Id:          1,
			Login:       "john_doe",
			Firstname:   "John",
			Surname:     "Doe",
			Gender:      "male",
			Avatar:      "avatar123.jpg",
			PhoneNumber: "123456789",
			Description: "Some description",
		},
	})

	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.NotNil(t, reply.User)
}

func TestUpdateUser_UpdateFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)

	server := NewUserServer(mockUserUseCase)

	ctx := GetCTX()

	inputUser := &proto.User{Id: 1}

	mockUserUseCase.EXPECT().UpdateUser(gomock.Any(), ctx).Return(nil, fmt.Errorf("update failed"))

	reply, err := server.UpdateUser(ctx, &proto.UpdateUserRequest{User: inputUser})

	assert.Error(t, err)
	assert.Nil(t, reply)
}

func TestDeleteUserById_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)

	server := NewUserServer(mockUserUseCase)

	ctx := GetCTX()

	mockUserUseCase.EXPECT().DeleteUserByID(gomock.Any(), ctx).Return(true, nil)

	reply, err := server.DeleteUserById(ctx, &proto.DeleteUserByIdRequest{Id: 1})

	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Equal(t, true, reply.Status)
}

func TestDeleteUserById_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)

	server := NewUserServer(mockUserUseCase)

	ctx := GetCTX()

	mockUserUseCase.EXPECT().DeleteUserByID(gomock.Any(), ctx).Times(0)

	reply, err := server.DeleteUserById(ctx, &proto.DeleteUserByIdRequest{Id: 0})

	assert.Error(t, err)
	assert.Nil(t, reply)
}

func TestUploadUserAvatar_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)

	server := NewUserServer(mockUserUseCase)

	ctx := GetCTX()

	mockUserUseCase.EXPECT().AddAvatar(gomock.Any(), "avatar123.jpg", ctx).Return(true, nil)

	reply, err := server.UploadUserAvatar(ctx, &proto.UploadUserAvatarRequest{Id: 1, Avatar: "avatar123.jpg"})

	assert.NoError(t, err)
	assert.NotNil(t, reply)
}

func TestUploadUserAvatar_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)

	server := NewUserServer(mockUserUseCase)

	ctx := GetCTX()

	mockUserUseCase.EXPECT().AddAvatar(gomock.Any(), gomock.Any(), ctx).Times(0)

	reply, err := server.UploadUserAvatar(ctx, &proto.UploadUserAvatarRequest{Id: 0, Avatar: "avatar123.jpg"})

	assert.Error(t, err)
	assert.Nil(t, reply)
}

func TestUploadUserAvatar_EmptyAvatar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)

	server := NewUserServer(mockUserUseCase)

	ctx := GetCTX()

	mockUserUseCase.EXPECT().AddAvatar(gomock.Any(), gomock.Any(), ctx).Times(0)

	reply, err := server.UploadUserAvatar(ctx, &proto.UploadUserAvatarRequest{Id: 1, Avatar: ""})

	assert.Error(t, err)
	assert.Nil(t, reply)
}

func TestDeleteUserAvatar_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)

	server := NewUserServer(mockUserUseCase)

	ctx := GetCTX()

	mockUserUseCase.EXPECT().GetUserByID(uint32(1), ctx).Return(&domain_models.User{ID: uint32(1), AvatarID: "avatar123.jpg"}, nil)
	mockUserUseCase.EXPECT().DeleteAvatarByUserID(uint32(1), ctx).Return(nil)

	reply, err := server.DeleteUserAvatar(ctx, &proto.DeleteUserAvatarRequest{Id: 1})

	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.True(t, reply.Status)
}

func TestDeleteUserAvatar_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)

	server := NewUserServer(mockUserUseCase)

	ctx := GetCTX()

	mockUserUseCase.EXPECT().GetUserByID(gomock.Any(), ctx).Times(0)

	reply, err := server.DeleteUserAvatar(ctx, &proto.DeleteUserAvatarRequest{Id: 0})

	assert.Error(t, err)
	assert.Nil(t, reply)
}

func TestDeleteUserAvatar_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)

	server := NewUserServer(mockUserUseCase)

	ctx := GetCTX()

	mockUserUseCase.EXPECT().GetUserByID(uint32(1), ctx).Return(nil, fmt.Errorf("user not found")).Times(1)

	reply, err := server.DeleteUserAvatar(ctx, &proto.DeleteUserAvatarRequest{Id: 1})

	assert.Error(t, err)
	assert.Nil(t, reply)
}

func TestDeleteUserAvatar_DeleteAvatarError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)

	server := NewUserServer(mockUserUseCase)

	ctx := GetCTX()

	mockUserUseCase.EXPECT().GetUserByID(uint32(1), ctx).Return(&domain_models.User{ID: uint32(1), AvatarID: "avatar123.jpg"}, nil)
	mockUserUseCase.EXPECT().DeleteAvatarByUserID(uint32(1), ctx).Return(fmt.Errorf("delete avatar error"))

	reply, err := server.DeleteUserAvatar(ctx, &proto.DeleteUserAvatarRequest{Id: 1})

	assert.Error(t, err)
	assert.Nil(t, reply)
}

func TestCreateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)

	server := NewUserServer(mockUserUseCase)

	ctx := GetCTX()

	mockUserUseCase.EXPECT().CreateUser(gomock.Any(), ctx).Return(&domain_models.User{ID: uint32(1), Login: "john_doe", FirstName: "John", Surname: "Doe", Gender: "male"}, nil)

	reply, err := server.CreateUser(ctx, &proto.CreateUserRequest{
		User: &proto.User{
			Login:       "ivan@mailhub.su",
			Password:    "password",
			Firstname:   "John",
			Surname:     "Doe",
			Gender:      "Male",
			Avatar:      "avatar123.jpg",
			PhoneNumber: "123456789",
			Description: "Some description",
		},
	})

	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.NotNil(t, reply.User)
}

func TestCreateUser_InvalidInput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)

	server := NewUserServer(mockUserUseCase)

	ctx := GetCTX()

	reply, err := server.CreateUser(ctx, &proto.CreateUserRequest{
		User: &proto.User{
			Login:       "ivan@mailhub.su",
			Password:    "password",
			Firstname:   "",
			Surname:     "Doe",
			Gender:      "Male",
			Avatar:      "avatar123.jpg",
			PhoneNumber: "123456789",
			Description: "Some description",
		},
	})

	assert.Error(t, err)
	assert.Nil(t, reply)
}

func TestCreateUser_InvalidLoginFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)

	server := NewUserServer(mockUserUseCase)

	ctx := GetCTX()

	reply, err := server.CreateUser(ctx, &proto.CreateUserRequest{
		User: &proto.User{
			Login:       "john_doe@example",
			Password:    "password",
			Firstname:   "John",
			Surname:     "Doe",
			Gender:      "male",
			Avatar:      "avatar123.jpg",
			PhoneNumber: "123456789",
			Description: "Some description",
		},
	})

	assert.Error(t, err)
	assert.Nil(t, reply)
}

func TestCreateUser_UserCreationFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)

	server := NewUserServer(mockUserUseCase)

	ctx := GetCTX()

	mockUserUseCase.EXPECT().CreateUser(gomock.Any(), ctx).Return(nil, fmt.Errorf("user creation failed"))

	reply, err := server.CreateUser(ctx, &proto.CreateUserRequest{
		User: &proto.User{
			Login:       "ivan@mailhub.su",
			Password:    "password",
			Firstname:   "John",
			Surname:     "Doe",
			Gender:      "Male",
			Avatar:      "avatar123.jpg",
			PhoneNumber: "123456789",
			Description: "Some description",
		},
	})

	assert.Error(t, err)
	assert.Nil(t, reply)
}
