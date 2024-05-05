package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"mail/internal/microservice/auth/proto"

	session_mock "mail/internal/microservice/session/mock"
	session_proto "mail/internal/microservice/session/proto"
	user_mock "mail/internal/microservice/user/mock"
	user_proto "mail/internal/microservice/user/proto"
)

func TestAuthServer_Login_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{"requestID": "testID"}))

	mockUserServiceClient := user_mock.NewMockUserServiceClient(ctrl)
	mockSessionServiceClient := session_mock.NewMockSessionServiceClient(ctrl)

	server := NewAuthServer(mockSessionServiceClient, mockUserServiceClient)

	loginRequest := &proto.LoginRequest{Login: "user@mailhub.su", Password: "password123"}

	mockUser := &user_proto.User{Id: 123}

	mockUserServiceClient.EXPECT().GetUserByLogin(gomock.Any(), gomock.Any()).Return(&user_proto.GetUserByLoginReply{User: mockUser}, nil)
	mockSessionServiceClient.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Return(&session_proto.CreateSessionReply{SessionId: "10101010"}, nil)

	reply, err := server.Login(ctx, loginRequest)

	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.True(t, reply.LoginStatus)
}

func TestAuthServer_Login_InvalidLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{"requestID": "testID"}))

	mockUserServiceClient := user_mock.NewMockUserServiceClient(ctrl)
	mockSessionServiceClient := session_mock.NewMockSessionServiceClient(ctrl)

	server := NewAuthServer(mockSessionServiceClient, mockUserServiceClient)

	loginRequest := &proto.LoginRequest{Login: "user", Password: "password123"}

	reply, err := server.Login(ctx, loginRequest)

	assert.Error(t, err)
	assert.Nil(t, reply)
	assert.EqualError(t, err, "domain in the login is not suitable")
}

func TestAuthServer_Login_EmptyFields(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{"requestID": "testID"}))

	mockUserServiceClient := user_mock.NewMockUserServiceClient(ctrl)
	mockSessionServiceClient := session_mock.NewMockSessionServiceClient(ctrl)

	server := NewAuthServer(mockSessionServiceClient, mockUserServiceClient)

	loginRequest := &proto.LoginRequest{Login: "", Password: ""}

	reply, err := server.Login(ctx, loginRequest)

	assert.Error(t, err)
	assert.Nil(t, reply)
	assert.EqualError(t, err, "all fields must be filled in")
}

func TestAuthServer_Login_UserServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{"requestID": "testID"}))

	mockUserServiceClient := user_mock.NewMockUserServiceClient(ctrl)
	mockSessionServiceClient := session_mock.NewMockSessionServiceClient(ctrl)

	server := NewAuthServer(mockSessionServiceClient, mockUserServiceClient)

	loginRequest := &proto.LoginRequest{Login: "user@mailhub.su", Password: "password123"}

	mockUserServiceClient.EXPECT().GetUserByLogin(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("user service error"))

	reply, err := server.Login(ctx, loginRequest)

	assert.Error(t, err)
	assert.NotNil(t, reply)
	assert.False(t, reply.LoginStatus)
}

func TestAuthServer_Login_SessionServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{"requestID": "testID"}))

	mockUserServiceClient := user_mock.NewMockUserServiceClient(ctrl)
	mockSessionServiceClient := session_mock.NewMockSessionServiceClient(ctrl)

	server := NewAuthServer(mockSessionServiceClient, mockUserServiceClient)

	loginRequest := &proto.LoginRequest{Login: "user@mailhub.su", Password: "password123"}

	mockUser := &user_proto.User{Id: 123}

	mockUserServiceClient.EXPECT().GetUserByLogin(gomock.Any(), gomock.Any()).Return(&user_proto.GetUserByLoginReply{User: mockUser}, nil)
	mockSessionServiceClient.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("session service error"))

	reply, err := server.Login(ctx, loginRequest)

	assert.Error(t, err)
	assert.NotNil(t, reply)
	assert.False(t, reply.LoginStatus)
}

func TestAuthServer_Signup_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{"requestID": "testID"}))

	mockUserServiceClient := user_mock.NewMockUserServiceClient(ctrl)
	mockSessionServiceClient := session_mock.NewMockSessionServiceClient(ctrl)

	server := NewAuthServer(mockSessionServiceClient, mockUserServiceClient)

	signupRequest := &proto.SignupRequest{
		Login:       "user@mailhub.su",
		Password:    "password123",
		Firstname:   "John",
		Surname:     "Doe",
		Patronymic:  "Smith",
		PhoneNumber: "123456789",
		Gender:      "male",
	}

	mockUserServiceClient.EXPECT().IsLoginUnique(gomock.Any(), gomock.Any()).Return(&user_proto.IsLoginUniqueReply{}, nil)
	mockUserServiceClient.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(&user_proto.CreateUserReply{}, nil)

	reply, err := server.Signup(ctx, signupRequest)

	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.True(t, reply.SignupStatus)
}

func TestAuthServer_Signup_InvalidLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{"requestID": "testID"}))

	mockUserServiceClient := user_mock.NewMockUserServiceClient(ctrl)
	mockSessionServiceClient := session_mock.NewMockSessionServiceClient(ctrl)

	server := NewAuthServer(mockSessionServiceClient, mockUserServiceClient)

	signupRequest := &proto.SignupRequest{
		Login:       "invalid_email",
		Password:    "password123",
		Firstname:   "John",
		Surname:     "Doe",
		Patronymic:  "Smith",
		PhoneNumber: "123456789",
		Gender:      "male",
	}

	reply, err := server.Signup(ctx, signupRequest)

	assert.Error(t, err)
	assert.Nil(t, reply)
	assert.EqualError(t, err, "domain in the login is not suitable")
}

func TestAuthServer_Signup_EmptyFields(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{"requestID": "testID"}))

	mockUserServiceClient := user_mock.NewMockUserServiceClient(ctrl)
	mockSessionServiceClient := session_mock.NewMockSessionServiceClient(ctrl)

	server := NewAuthServer(mockSessionServiceClient, mockUserServiceClient)

	signupRequest := &proto.SignupRequest{}

	reply, err := server.Signup(ctx, signupRequest)

	assert.Error(t, err)
	assert.Nil(t, reply)
	assert.EqualError(t, err, "all fields must be filled in")
}

func TestAuthServer_Signup_DuplicateLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{"requestID": "testID"}))

	mockUserServiceClient := user_mock.NewMockUserServiceClient(ctrl)
	mockSessionServiceClient := session_mock.NewMockSessionServiceClient(ctrl)

	server := NewAuthServer(mockSessionServiceClient, mockUserServiceClient)

	signupRequest := &proto.SignupRequest{
		Login:       "user@mailhub.su",
		Password:    "password123",
		Firstname:   "John",
		Surname:     "Doe",
		Patronymic:  "Smith",
		PhoneNumber: "123456789",
		Gender:      "male",
	}

	mockUserServiceClient.EXPECT().IsLoginUnique(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("login already exists"))

	reply, err := server.Signup(ctx, signupRequest)

	assert.Error(t, err)
	assert.Nil(t, reply)
	assert.EqualError(t, err, "such a login already exists")
}

func TestAuthServer_Signup_UserServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{"requestID": "testID"}))

	mockUserServiceClient := user_mock.NewMockUserServiceClient(ctrl)
	mockSessionServiceClient := session_mock.NewMockSessionServiceClient(ctrl)

	server := NewAuthServer(mockSessionServiceClient, mockUserServiceClient)

	signupRequest := &proto.SignupRequest{
		Login:       "user@mailhub.su",
		Password:    "password123",
		Firstname:   "John",
		Surname:     "Doe",
		Patronymic:  "Smith",
		PhoneNumber: "123456789",
		Gender:      "male",
	}

	mockUserServiceClient.EXPECT().IsLoginUnique(gomock.Any(), gomock.Any()).Return(&user_proto.IsLoginUniqueReply{}, nil)
	mockUserServiceClient.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("user service error"))

	reply, err := server.Signup(ctx, signupRequest)

	assert.Error(t, err)
	assert.Nil(t, reply)
	assert.EqualError(t, err, "failed to add user")
}

func TestAuthServer_Logout_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{"requestID": "testID"}))

	mockSessionServiceClient := session_mock.NewMockSessionServiceClient(ctrl)

	server := NewAuthServer(mockSessionServiceClient, nil)

	logoutRequest := &proto.LogoutRequest{
		SessionId: "10101010",
	}

	mockSessionServiceClient.EXPECT().DeleteSession(gomock.Any(), gomock.Any()).Return(&session_proto.DeleteSessionReply{}, nil)

	reply, err := server.Logout(ctx, logoutRequest)

	assert.NoError(t, err)
	assert.NotNil(t, reply)
}

func TestAuthServer_Logout_EmptySessionId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{"requestID": "testID"}))

	mockSessionServiceClient := session_mock.NewMockSessionServiceClient(ctrl)

	server := NewAuthServer(mockSessionServiceClient, nil)

	logoutRequest := &proto.LogoutRequest{}

	reply, err := server.Logout(ctx, logoutRequest)

	assert.Error(t, err)
	assert.Nil(t, reply)
	assert.EqualError(t, err, "session id must be filled in")
}

func TestAuthServer_Logout_SessionServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{"requestID": "testID"}))

	mockSessionServiceClient := session_mock.NewMockSessionServiceClient(ctrl)

	server := NewAuthServer(mockSessionServiceClient, nil)

	logoutRequest := &proto.LogoutRequest{
		SessionId: "10101010",
	}

	mockSessionServiceClient.EXPECT().DeleteSession(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("session service error"))

	reply, err := server.Logout(ctx, logoutRequest)

	assert.Error(t, err)
	assert.Nil(t, reply)
	assert.EqualError(t, err, "delete session failed")
}
