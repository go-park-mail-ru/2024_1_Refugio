package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"mail/internal/models/microservice_ports"

	"mail/internal/microservice/auth/proto"
	user_proto "mail/internal/microservice/user/proto"
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

func (us *AuthServer) Login(ctx context.Context, input *proto.UserLogin) (*proto.StatusLogin, error) {
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

	userServiceClient := user_proto.NewUserServiceClient(conn)
	_, errLogin := userServiceClient.GetUserByLogin(
		metadata.NewOutgoingContext(ctx,
			metadata.New(map[string]string{"requestID": ctx.Value("requestID").(string)})),
		&user_proto.UserLogin{Login: input.Login, Password: input.Password},
	)
	if errLogin != nil {
		return nil, fmt.Errorf("login failed")
	}

	return &proto.StatusLogin{LoginStatus: true}, nil
}
