package server

import (
	"context"
	"fmt"
	converters "mail/internal/microservice/models/proto_converters"
	"mail/internal/microservice/user/proto"
	"mail/internal/microservice/user/usecase"
)

type UserServer struct {
	proto.UnimplementedUserServiceServer

	UserUseCase *usecase.UserUseCase
}

func NewUserServer(userUseCase *usecase.UserUseCase) *UserServer {
	return &UserServer{UserUseCase: userUseCase}
}

func (us *UserServer) GetUser(ctx context.Context, input *proto.UserId) (*proto.User, error) {
	user, err := us.UserUseCase.GetUserByID(input.Id, ctx)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return converters.UserConvertCoreInProto(*user), nil
}
