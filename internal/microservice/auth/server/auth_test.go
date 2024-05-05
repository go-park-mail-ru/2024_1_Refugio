package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"mail/internal/pkg/logger"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"mail/internal/microservice/auth/proto"
	session_mock "mail/internal/microservice/session/mock"
	session_proto "mail/internal/microservice/session/proto"
	user_mock "mail/internal/microservice/user/mock"
	user_proto "mail/internal/microservice/user/proto"
)

func GetCTX() context.Context {
	f, err := os.OpenFile("log_test.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create logfile" + "log.txt")
	}
	defer f.Close()

	ctx := context.WithValue(context.Background(), "logger", logger.InitializationBdLog(f))
	ctx2 := context.WithValue(ctx, "requestID", []string{"testID"})

	return ctx2
}

func TestAuthServer_Login_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserServiceClient := user_mock.NewMockUserServiceClient(ctrl)
	mockSessionServiceClient := session_mock.NewMockSessionServiceClient(ctrl)

	server := NewAuthServer(mockSessionServiceClient, mockUserServiceClient)

	ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{"requestID": "testID"}))

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

	server := &AuthServer{}

	ctx := GetCTX()

	loginRequest := &proto.LoginRequest{Login: "user", Password: "password123"}

	reply, err := server.Login(ctx, loginRequest)

	assert.Error(t, err)
	assert.Nil(t, reply)
	assert.EqualError(t, err, "domain in the login is not suitable")
}

func TestAuthServer_Login_EmptyFields(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	server := &AuthServer{}

	ctx := GetCTX()

	loginRequest := &proto.LoginRequest{Login: "", Password: ""}

	reply, err := server.Login(ctx, loginRequest)

	assert.Error(t, err)
	assert.Nil(t, reply)
	assert.EqualError(t, err, "all fields must be filled in")
}

/*
func TestAuthServer_Login_UserServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserServiceClient := user_mock.NewMockUserServer(ctrl)

	server := &AuthServer{}

	ctx := context.Background()

	loginRequest := &proto.LoginRequest{Login: "user@mailhub.su", Password: "password123"}

	mockUserServiceClient.EXPECT().GetUserByLogin(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("user service error"))

	reply, err := server.Login(ctx, loginRequest)

	assert.Error(t, err)
	assert.Nil(t, reply)
	assert.EqualError(t, err, "login failed")
}
*/

/*
func TestAuthServer_Login_SessionServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserServiceClient := user_mock.NewMockUserServer(ctrl)
	mockSessionServiceClient := session_mock.NewMockSessionServer(ctrl)

	server := &AuthServer{}

	ctx := GetCTX()

	loginRequest := &proto.LoginRequest{Login: "user@mailhub.su", Password: "password123"}

	mockUser := &user_proto.User{Id: 123}

	mockUserServiceClient.EXPECT().GetUserByLogin(gomock.Any(), gomock.Any()).Return(&user_proto.GetUserByLoginReply{User: mockUser}, nil)
	mockSessionServiceClient.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("session service error"))

	reply, err := server.Login(ctx, loginRequest)

	assert.Error(t, err)
	assert.Nil(t, reply)
	assert.EqualError(t, err, "create session failed")
}
*/
