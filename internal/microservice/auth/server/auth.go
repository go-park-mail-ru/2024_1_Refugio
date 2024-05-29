package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"

	"mail/internal/microservice/auth/proto"
	"mail/internal/models/microservice_ports"
	"mail/internal/pkg/utils/connect_microservice"
	"mail/internal/pkg/utils/sanitize"

	domain "mail/internal/microservice/models/domain_models"
	session_proto "mail/internal/microservice/session/proto"
	user_proto "mail/internal/microservice/user/proto"
	validUtil "mail/internal/pkg/utils/validators"
)

// AuthServer handles RPC calls for the AuthService.
type AuthServer struct {
	proto.UnimplementedAuthServiceServer
	sessionServiceClient session_proto.SessionServiceClient
	userServiceClient    user_proto.UserServiceClient
}

// NewAuthServer creates a new instance of AuthServer.
func NewAuthServer(sessionClient session_proto.SessionServiceClient, userClient user_proto.UserServiceClient) *AuthServer {
	return &AuthServer{
		sessionServiceClient: sessionClient,
		userServiceClient:    userClient,
	}
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

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("metadata error")
	}
	value := md.Get("requestID")

	user, errLogin := as.userServiceClient.GetUserByLogin(
		metadata.NewOutgoingContext(ctx,
			metadata.New(map[string]string{"requestID": value[0]})),
		&user_proto.GetUserByLoginRequest{Login: input.Login, Password: input.Password},
	)
	if errLogin != nil {
		return &proto.LoginReply{LoginStatus: false}, fmt.Errorf("login failed")
	}

	session, errStatus := as.sessionServiceClient.CreateSession(
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

// LoginVK handles user login.
func (as *AuthServer) LoginVK(ctx context.Context, input *proto.LoginVKRequest) (*proto.LoginReply, error) {
	if input.VkId <= 0 {
		return nil, fmt.Errorf("bad vkId = %d", input.VkId)
	}

	conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.UserService))
	if err != nil {
		return nil, fmt.Errorf("invalid connection")
	}
	defer conn.Close()

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		fmt.Println("metadata error")
		return nil, fmt.Errorf("metadata error")
	}
	value := md.Get("requestID")

	userServiceClient := user_proto.NewUserServiceClient(conn)
	user, errId := userServiceClient.GetUserByVKId(
		metadata.NewOutgoingContext(ctx,
			metadata.New(map[string]string{"requestID": value[0]})),
		&user_proto.GetUserVKIdRequest{VkId: input.VkId},
	)
	if errId != nil {
		fmt.Println("vkId failed")
		return &proto.LoginReply{LoginStatus: false}, fmt.Errorf("vkId failed")
	}
	if user == nil {
		fmt.Println("User with vkId not found")
		return &proto.LoginReply{LoginStatus: false}, fmt.Errorf("User with vkId not found")
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
		fmt.Println("create session failed\"")
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

	_, errLogin := as.userServiceClient.IsLoginUnique(
		metadata.NewOutgoingContext(ctx,
			metadata.New(map[string]string{"requestID": value[0]})),
		&user_proto.IsLoginUniqueRequest{Login: input.Login},
	)
	if errLogin != nil {
		return nil, fmt.Errorf("such a login already exists")
	}

	_, errCreate := as.userServiceClient.CreateUser(
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

// SignupVK handles user signup.
func (as *AuthServer) SignupVK(ctx context.Context, input *proto.SignupVKRequest) (*proto.SignupReply, error) {
	input.Login = sanitize.SanitizeString(input.Login)
	input.Firstname = sanitize.SanitizeString(input.Firstname)
	input.Surname = sanitize.SanitizeString(input.Surname)

	if validUtil.IsEmpty(input.Login) || validUtil.IsEmpty(input.Firstname) || validUtil.IsEmpty(input.Surname) || !domain.IsValidGender(domain.GetGenderType(input.Gender)) {
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

	userVk, _ := userServiceClient.GetUserByVKId(
		metadata.NewOutgoingContext(ctx,
			metadata.New(map[string]string{"requestID": value[0]})),
		&user_proto.GetUserVKIdRequest{VkId: input.VkId},
	)
	if userVk != nil {
		return nil, fmt.Errorf("A user with this VKId has already been registered")
	}

	_, errCreate := userServiceClient.CreateUserOtherMail(
		metadata.NewOutgoingContext(ctx,
			metadata.New(map[string]string{"requestID": value[0]})),
		&user_proto.CreateUserRequest{User: &user_proto.User{
			Login:     input.Login,
			Firstname: input.Firstname,
			Surname:   input.Surname,
			Birthday:  input.Birthday,
			Gender:    input.Gender,
			VkId:      input.VkId,
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

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("metadata error")
	}
	value := md.Get("requestID")

	_, errStatus := as.sessionServiceClient.DeleteSession(
		metadata.NewOutgoingContext(ctx,
			metadata.New(map[string]string{"requestID": value[0]})),
		&session_proto.DeleteSessionRequest{SessionId: input.SessionId},
	)
	if errStatus != nil {
		return nil, fmt.Errorf("delete session failed")
	}

	return &proto.LogoutReply{}, nil
}

// SignupOtherMail handles user signup.
func (as *AuthServer) SignupOtherMail(ctx context.Context, input *proto.SignupOtherMailRequest) (*proto.SignupReply, error) {
	input.Login = sanitize.SanitizeString(input.Login)
	input.Firstname = sanitize.SanitizeString(input.Firstname)
	input.Surname = sanitize.SanitizeString(input.Surname)
	input.Patronymic = sanitize.SanitizeString(input.Patronymic)
	input.PhoneNumber = sanitize.SanitizeString(input.PhoneNumber)
	input.Description = sanitize.SanitizeString(input.Description)
	input.Avatar = sanitize.SanitizeString(input.Avatar)

	if validUtil.IsEmpty(input.Login) || validUtil.IsEmpty(input.Firstname) || validUtil.IsEmpty(input.Surname) {
		return nil, fmt.Errorf("all fields must be filled in")
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("metadata error")
	}
	value := md.Get("requestID")

	_, errLogin := as.userServiceClient.IsLoginUnique(
		metadata.NewOutgoingContext(ctx,
			metadata.New(map[string]string{"requestID": value[0]})),
		&user_proto.IsLoginUniqueRequest{Login: input.Login},
	)
	if errLogin != nil {
		return nil, fmt.Errorf("such a login already exists")
	}

	_, errCreate := as.userServiceClient.CreateUserOtherMail(
		metadata.NewOutgoingContext(ctx,
			metadata.New(map[string]string{"requestID": value[0]})),
		&user_proto.CreateUserRequest{User: &user_proto.User{
			Login:       input.Login,
			Password:    "password",
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

// LoginOtherMail handles user login.
func (as *AuthServer) LoginOtherMail(ctx context.Context, input *proto.LoginOtherMailRequest) (*proto.LoginReply, error) {
	if input.Id <= 0 {
		return nil, fmt.Errorf("bad Id = %d", input.Id)
	}
	conn2, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.SessionService))
	if err != nil {
		return nil, fmt.Errorf("invalid connection")
	}
	defer conn2.Close()

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		fmt.Println("metadata error")
		return nil, fmt.Errorf("metadata error")
	}
	value := md.Get("requestID")

	sessionServiceClient := session_proto.NewSessionServiceClient(conn2)
	session, errStatus := sessionServiceClient.CreateSession(
		metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{"requestID": value[0]})),
		&session_proto.CreateSessionRequest{Session: &session_proto.Session{UserId: input.Id,
			Device:   "",
			LifeTime: 60 * 60 * 24},
		},
	)
	if errStatus != nil {
		fmt.Println("create session failed\"")
		return &proto.LoginReply{LoginStatus: false}, fmt.Errorf("create session failed")
	}

	return &proto.LoginReply{LoginStatus: true, SessionId: session.SessionId}, nil
}
