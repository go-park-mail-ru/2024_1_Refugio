package server

import (
	"context"
	"fmt"
	"strings"

	converters "mail/internal/microservice/models/proto_converters"
	"mail/internal/microservice/user/proto"
	"mail/internal/microservice/user/usecase"
	"mail/internal/pkg/utils/sanitize"
)

type UserServer struct {
	proto.UnimplementedUserServiceServer

	UserUseCase *usecase.UserUseCase
}

func NewUserServer(userUseCase *usecase.UserUseCase) *UserServer {
	return &UserServer{UserUseCase: userUseCase}
}

func (us *UserServer) GetUser(ctx context.Context, input *proto.UserId) (*proto.User, error) {
	if input.Id <= 0 {
		return nil, fmt.Errorf("invalid user id: %s", input.Id)
	}

	user, err := us.UserUseCase.GetUserByID(input.Id, ctx)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return converters.UserConvertCoreInProto(*user), nil
}

func (us *UserServer) UpdateUser(ctx context.Context, input *proto.User) (*proto.User, error) {
	userDomain := converters.UserConvertProtoInCore(*input)

	userDomain.Login = sanitize.SanitizeString(userDomain.Login)
	userDomain.FirstName = sanitize.SanitizeString(userDomain.FirstName)
	userDomain.Surname = sanitize.SanitizeString(userDomain.Surname)
	userDomain.Patronymic = sanitize.SanitizeString(userDomain.Patronymic)
	userDomain.AvatarID = sanitize.SanitizeString(userDomain.AvatarID)
	userDomain.PhoneNumber = sanitize.SanitizeString(userDomain.PhoneNumber)
	userDomain.Description = sanitize.SanitizeString(userDomain.Description)

	user, err := us.UserUseCase.UpdateUser(userDomain, ctx)
	if err != nil {
		return nil, fmt.Errorf("user with id %s update fail", userDomain.ID)
	}

	return converters.UserConvertCoreInProto(*user), nil
}

func (us *UserServer) DeleteUserById(ctx context.Context, input *proto.UserId) (*proto.Status, error) {
	if input.Id <= 0 {
		return nil, fmt.Errorf("invalid user id: %s", input.Id)
	}

	userStatus, err := us.UserUseCase.DeleteUserByID(input.Id, ctx)
	if err != nil {
		return nil, fmt.Errorf("user with id %s delete fail", input.Id)
	}

	return &proto.Status{DeleteStatus: userStatus}, nil
}

func (us *UserServer) UploadUserAvatar(ctx context.Context, input *proto.UserAvatar) (*proto.Empty, error) {
	if input.Id <= 0 {
		return nil, fmt.Errorf("invalid user id: %s", input.Id)
	}

	input.Avatar = sanitize.SanitizeString(input.Avatar)
	if strings.TrimSpace(input.Avatar) == "" {
		return nil, fmt.Errorf("avatar has not been transferred")
	}

	user, err := us.UserUseCase.GetUserByID(input.Id, ctx)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	user.AvatarID = input.Avatar
	_, errUpdate := us.UserUseCase.UpdateUser(user, ctx)
	if errUpdate != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &proto.Empty{}, nil
}

func (us *UserServer) DeleteUserAvatar(ctx context.Context, input *proto.UserId) (*proto.Status, error) {
	if input.Id <= 0 {
		return nil, fmt.Errorf("invalid user id: %s", input.Id)
	}

	user, err := us.UserUseCase.GetUserByID(input.Id, ctx)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	user.AvatarID = ""
	_, errUpdate := us.UserUseCase.DeleteUserAvatar(user, ctx)
	if errUpdate != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &proto.Status{DeleteStatus: true}, nil
}
