package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"mail/internal/microservice/auth/proto"
	domain "mail/internal/microservice/models/domain_models"
	session_proto "mail/internal/microservice/session/proto"
	user_proto "mail/internal/microservice/user/proto"
	"mail/internal/models/microservice_ports"
	"mail/internal/pkg/utils/connect_microservice"
	"mail/internal/pkg/utils/sanitize"
	validUtil "mail/internal/pkg/utils/validators"
)

type AuthServer struct {
	proto.UnimplementedAuthServiceServer
}

func NewAuthServer() *AuthServer {
	return &AuthServer{}
}

// Login handles user login.
func (as *AuthServer) Login(ctx context.Context, input *proto.LoginRequest) (*proto.LoginReply, error) {
	input.Login = sanitize.SanitizeString(input.Login)
	input.Password = sanitize.SanitizeString(input.Password)

	if validUtil.IsEmpty(input.Login) || validUtil.IsEmpty(input.Password) {
		return nil, fmt.Errorf("all fields must be filled in")
	}

	if !validUtil.IsValidEmailFormat(input.Login) {
		return nil, fmt.Errorf("domain in the login is not suitable")
	}

	conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.UserService))
	if err != nil {
		return nil, fmt.Errorf("invalid connection")
	}
	defer conn.Close()

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("metadata error")
	}
	value := md.Get("requestID")

	userServiceClient := user_proto.NewUserServiceClient(conn)
	user, errLogin := userServiceClient.GetUserByLogin(
		metadata.NewOutgoingContext(ctx,
			metadata.New(map[string]string{"requestID": value[0]})),
		&user_proto.GetUserByLoginRequest{Login: input.Login, Password: input.Password},
	)
	if errLogin != nil {
		return &proto.LoginReply{LoginStatus: false}, fmt.Errorf("login failed")
	}

	conn2, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.SessionService))
	if err != nil {
		return nil, fmt.Errorf("invalid connection")
	}
	defer conn2.Close()

	sessionServiceClient := session_proto.NewSessionServiceClient(conn2)
	session, errStatus := sessionServiceClient.CreateSession(
		metadata.NewOutgoingContext(ctx,
			metadata.New(map[string]string{"requestID": value[0]})),
		&session_proto.CreateSessionRequest{Session: &session_proto.Session{UserId: user.User.Id,
			Device:   "",
			LifeTime: 60 * 60 * 24},
		},
	)
	if errStatus != nil {
		return &proto.LoginReply{LoginStatus: false}, fmt.Errorf("create session failed")
	}

	return &proto.LoginReply{LoginStatus: true, SessionId: session.SessionId}, nil
}

// Signup handles user signup.
func (as *AuthServer) Signup(ctx context.Context, input *proto.SignupRequest) (*proto.SignupReply, error) {
	input.Login = sanitize.SanitizeString(input.Login)
	input.Password = sanitize.SanitizeString(input.Password)
	input.Firstname = sanitize.SanitizeString(input.Firstname)
	input.Surname = sanitize.SanitizeString(input.Surname)
	input.Patronymic = sanitize.SanitizeString(input.Patronymic)
	input.PhoneNumber = sanitize.SanitizeString(input.PhoneNumber)
	input.Description = sanitize.SanitizeString(input.Description)
	input.Avatar = sanitize.SanitizeString(input.Avatar)

	if validUtil.IsEmpty(input.Login) || validUtil.IsEmpty(input.Password) || validUtil.IsEmpty(input.Firstname) || validUtil.IsEmpty(input.Surname) || !domain.IsValidGender(domain.GetGenderType(input.Gender)) {
		return nil, fmt.Errorf("all fields must be filled in")
	}

	if !validUtil.IsValidEmailFormat(input.Login) {
		return nil, fmt.Errorf("domain in the login is not suitable")
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("metadata error")
	}
	value := md.Get("requestID")

	conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.UserService))
	if err != nil {
		return nil, fmt.Errorf("invalid connection")
	}
	defer conn.Close()

	userServiceClient := user_proto.NewUserServiceClient(conn)
	_, errLogin := userServiceClient.IsLoginUnique(
		metadata.NewOutgoingContext(ctx,
			metadata.New(map[string]string{"requestID": value[0]})),
		&user_proto.IsLoginUniqueRequest{Login: input.Login},
	)
	if errLogin != nil {
		return nil, fmt.Errorf("such a login already exists")
	}

	_, errCreate := userServiceClient.CreateUser(
		metadata.NewOutgoingContext(ctx,
			metadata.New(map[string]string{"requestID": value[0]})),
		&user_proto.CreateUserRequest{User: &user_proto.User{
			Login:       input.Login,
			Password:    input.Password,
			Firstname:   input.Firstname,
			Surname:     input.Surname,
			Patronymic:  input.Patronymic,
			Birthday:    input.Birthday,
			Gender:      input.Gender,
			Avatar:      input.Avatar,
			PhoneNumber: input.PhoneNumber,
			Description: input.Description,
		}},
	)
	if errCreate != nil {
		return nil, fmt.Errorf("failed to add user")
	}

	return &proto.SignupReply{SignupStatus: true}, nil
}

// Logout handles user logout.
func (as *AuthServer) Logout(ctx context.Context, input *proto.LogoutRequest) (*proto.LogoutReply, error) {
	input.SessionId = sanitize.SanitizeString(input.SessionId)

	if validUtil.IsEmpty(input.SessionId) {
		return nil, fmt.Errorf("session id must be filled in")
	}

	conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.SessionService))
	if err != nil {
		return nil, fmt.Errorf("invalid connection")
	}
	defer conn.Close()

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("metadata error")
	}
	value := md.Get("requestID")

	sessionServiceClient := session_proto.NewSessionServiceClient(conn)
	_, errStatus := sessionServiceClient.DeleteSession(
		metadata.NewOutgoingContext(ctx,
			metadata.New(map[string]string{"requestID": value[0]})),
		&session_proto.DeleteSessionRequest{SessionId: input.SessionId},
	)
	if errStatus != nil {
		return nil, fmt.Errorf("delete session failed")
	}

	return &proto.LogoutReply{}, nil
}
