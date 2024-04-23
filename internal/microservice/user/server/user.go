package server

import (
	"context"
	"fmt"
	"strings"

	converters "mail/internal/microservice/models/proto_converters"
	usecase "mail/internal/microservice/user/interface"
	"mail/internal/microservice/user/proto"
	"mail/internal/pkg/utils/sanitize"
)

type UserServer struct {
	proto.UnimplementedUserServiceServer
	UserUseCase usecase.UserUseCase
}

func NewUserServer(userUseCase usecase.UserUseCase) *UserServer {
	return &UserServer{UserUseCase: userUseCase}
}

// GetUsers retrieves information about users.
func (us *UserServer) GetUsers(ctx context.Context, input *proto.Empty) (*proto.Users, error) {
	users, err := us.UserUseCase.GetAllUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	usersProto := make([]*proto.User, 0, len(users))
	for _, user := range users {
		usersProto = append(usersProto, converters.UserConvertCoreInProto(*user))
	}

	return &proto.Users{Users: usersProto}, nil
}

// GetUser retrieves information about a user by their identifier.
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

// GetUserByLogin retrieves information about a user by their login.
func (us *UserServer) GetUserByLogin(ctx context.Context, input *proto.UserLogin) (*proto.User, error) {
	if strings.TrimSpace(input.Login) == "" && strings.TrimSpace(input.Password) == "" {
		return nil, fmt.Errorf("invalid user login: %s", input.Login)
	}

	user, err := us.UserUseCase.GetUserByLogin(input.Login, input.Password, ctx)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return converters.UserConvertCoreInProto(*user), nil
}

// IsLoginUnique checks if the provided login is unique among all users.
func (us *UserServer) IsLoginUnique(ctx context.Context, input *proto.Login) (*proto.Status, error) {
	if strings.TrimSpace(input.Login) == "" {
		return nil, fmt.Errorf("invalid user login: %s", input.Login)
	}

	status, err := us.UserUseCase.IsLoginUnique(input.Login, ctx)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &proto.Status{Status: status}, nil
}

// UpdateUser updates user information.
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

// DeleteUserById deletes a user by their identifier.
func (us *UserServer) DeleteUserById(ctx context.Context, input *proto.UserId) (*proto.Status, error) {
	if input.Id <= 0 {
		return nil, fmt.Errorf("invalid user id: %s", input.Id)
	}

	userStatus, err := us.UserUseCase.DeleteUserByID(input.Id, ctx)
	if err != nil {
		return nil, fmt.Errorf("user with id %s delete fail", input.Id)
	}

	return &proto.Status{Status: userStatus}, nil
}

// UploadUserAvatar uploads a user's avatar.
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

// DeleteUserAvatar deletes a user's avatar.
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

	return &proto.Status{Status: true}, nil
}

// CreateUser creates user.
func (us *UserServer) CreateUser(ctx context.Context, input *proto.User) (*proto.User, error) {
	userDomain := converters.UserConvertProtoInCore(*input)

	userDomain.Login = sanitize.SanitizeString(userDomain.Login)
	userDomain.FirstName = sanitize.SanitizeString(userDomain.FirstName)
	userDomain.Surname = sanitize.SanitizeString(userDomain.Surname)
	userDomain.Patronymic = sanitize.SanitizeString(userDomain.Patronymic)
	userDomain.AvatarID = sanitize.SanitizeString(userDomain.AvatarID)
	userDomain.PhoneNumber = sanitize.SanitizeString(userDomain.PhoneNumber)
	userDomain.Description = sanitize.SanitizeString(userDomain.Description)

	user, err := us.UserUseCase.CreateUser(userDomain, ctx)
	if err != nil {
		return nil, fmt.Errorf("user with login %s create fail", userDomain.Login)
	}

	return converters.UserConvertCoreInProto(*user), nil
}
