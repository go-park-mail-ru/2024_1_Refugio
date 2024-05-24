//go:generate mockgen -source=./user.go -destination=../mock/user_mock.go -package=mock

package server

import (
	"context"
	"fmt"
	"strings"

	"mail/internal/microservice/user/proto"
	"mail/internal/pkg/utils/sanitize"

	domain "mail/internal/microservice/models/domain_models"
	converters "mail/internal/microservice/models/proto_converters"
	usecase "mail/internal/microservice/user/interface"
	validUtil "mail/internal/pkg/utils/validators"
)

// UserServer handles RPC calls for the UserService.
type UserServer struct {
	proto.UnimplementedUserServiceServer
	UserUseCase usecase.UserUseCase
}

// NewUserServer creates a new instance of UserServer.
func NewUserServer(userUseCase usecase.UserUseCase) *UserServer {
	return &UserServer{UserUseCase: userUseCase}
}

// GetUsers retrieves information about users.
func (us *UserServer) GetUsers(ctx context.Context, input *proto.GetUsersRequest) (*proto.GetUsersReply, error) {
	if input == nil {
		return &proto.GetUsersReply{Users: []*proto.User{}}, nil
	}

	users, err := us.UserUseCase.GetAllUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	usersProto := make([]*proto.User, 0, len(users))
	for _, user := range users {
		usersProto = append(usersProto, converters.UserConvertCoreInProto(user))
	}

	return &proto.GetUsersReply{Users: usersProto}, nil
}

// GetUser retrieves information about a user by their identifier.
func (us *UserServer) GetUser(ctx context.Context, input *proto.GetUserRequest) (*proto.GetUserReply, error) {
	if input.Id <= 0 {
		return nil, fmt.Errorf("invalid user id: %v", input.Id)
	}

	user, err := us.UserUseCase.GetUserByID(input.Id, ctx)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &proto.GetUserReply{User: converters.UserConvertCoreInProto(user)}, nil
}

// GetUserByLogin retrieves information about a user by their login.
func (us *UserServer) GetUserByLogin(ctx context.Context, input *proto.GetUserByLoginRequest) (*proto.GetUserByLoginReply, error) {
	if strings.TrimSpace(input.Login) == "" && strings.TrimSpace(input.Password) == "" {
		return nil, fmt.Errorf("invalid user login: %s", input.Login)
	}

	user, err := us.UserUseCase.GetUserByLogin(input.Login, input.Password, ctx)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &proto.GetUserByLoginReply{User: converters.UserConvertCoreInProto(user)}, nil
}

// IsLoginUnique checks if the provided login is unique among all users.
func (us *UserServer) IsLoginUnique(ctx context.Context, input *proto.IsLoginUniqueRequest) (*proto.IsLoginUniqueReply, error) {
	if strings.TrimSpace(input.Login) == "" {
		return nil, fmt.Errorf("invalid user login: %s", input.Login)
	}

	status, err := us.UserUseCase.IsLoginUnique(input.Login, ctx)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &proto.IsLoginUniqueReply{Status: status}, nil
}

// UpdateUser updates user information.
func (us *UserServer) UpdateUser(ctx context.Context, input *proto.UpdateUserRequest) (*proto.UpdateUserReply, error) {
	userDomain := converters.UserConvertProtoInCore(input.User)

	userDomain.Login = sanitize.SanitizeString(userDomain.Login)
	userDomain.FirstName = sanitize.SanitizeString(userDomain.FirstName)
	userDomain.Surname = sanitize.SanitizeString(userDomain.Surname)
	userDomain.Patronymic = sanitize.SanitizeString(userDomain.Patronymic)
	userDomain.AvatarID = sanitize.SanitizeString(userDomain.AvatarID)
	userDomain.PhoneNumber = sanitize.SanitizeString(userDomain.PhoneNumber)
	userDomain.Description = sanitize.SanitizeString(userDomain.Description)

	user, err := us.UserUseCase.UpdateUser(userDomain, ctx)
	if err != nil {
		return nil, fmt.Errorf("user with id %v update fail", userDomain.ID)
	}

	return &proto.UpdateUserReply{User: converters.UserConvertCoreInProto(user)}, nil
}

// DeleteUserById deletes a user by their identifier.
func (us *UserServer) DeleteUserById(ctx context.Context, input *proto.DeleteUserByIdRequest) (*proto.DeleteUserByIdReply, error) {
	if input.Id <= 0 {
		return nil, fmt.Errorf("invalid user id: %v", input.Id)
	}

	userStatus, err := us.UserUseCase.DeleteUserByID(input.Id, ctx)
	if err != nil {
		return nil, fmt.Errorf("user with id %v delete fail", input.Id)
	}

	return &proto.DeleteUserByIdReply{Status: userStatus}, nil
}

// UploadUserAvatar uploads a user's avatar.
func (us *UserServer) UploadUserAvatar(ctx context.Context, input *proto.UploadUserAvatarRequest) (*proto.UploadUserAvatarReply, error) {
	if input.Id <= 0 {
		return nil, fmt.Errorf("invalid user id: %v", input.Id)
	}

	input.Avatar = sanitize.SanitizeString(input.Avatar)
	if strings.TrimSpace(input.Avatar) == "" {
		return nil, fmt.Errorf("avatar has not been transferred")
	}

	_, errUpdate := us.UserUseCase.AddAvatar(input.Id, input.Avatar, ctx)
	if errUpdate != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &proto.UploadUserAvatarReply{}, nil
}

// DeleteUserAvatar deletes a user's avatar.
func (us *UserServer) DeleteUserAvatar(ctx context.Context, input *proto.DeleteUserAvatarRequest) (*proto.DeleteUserAvatarReply, error) {
	if input.Id <= 0 {
		return nil, fmt.Errorf("invalid user id: %v", input.Id)
	}

	user, err := us.UserUseCase.GetUserByID(input.Id, ctx)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	user.AvatarID = ""
	errUpdate := us.UserUseCase.DeleteAvatarByUserID(user.ID, ctx)
	if errUpdate != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &proto.DeleteUserAvatarReply{Status: true}, nil
}

// CreateUser creates user.
func (us *UserServer) CreateUser(ctx context.Context, input *proto.CreateUserRequest) (*proto.CreateUserReply, error) {
	userDomain := converters.UserConvertProtoInCore(input.User)

	userDomain.Login = sanitize.SanitizeString(userDomain.Login)
	userDomain.FirstName = sanitize.SanitizeString(userDomain.FirstName)
	userDomain.Surname = sanitize.SanitizeString(userDomain.Surname)
	userDomain.Patronymic = sanitize.SanitizeString(userDomain.Patronymic)
	userDomain.AvatarID = sanitize.SanitizeString(userDomain.AvatarID)
	userDomain.PhoneNumber = sanitize.SanitizeString(userDomain.PhoneNumber)
	userDomain.Description = sanitize.SanitizeString(userDomain.Description)

	if validUtil.IsEmpty(userDomain.Login) || validUtil.IsEmpty(userDomain.Password) || validUtil.IsEmpty(userDomain.FirstName) || validUtil.IsEmpty(userDomain.Surname) || !domain.IsValidGender(userDomain.Gender) {
		return nil, fmt.Errorf("all fields must be filled in")
	}

	if !validUtil.IsValidEmailFormat(userDomain.Login) {
		return nil, fmt.Errorf("domain in the login is not suitable")
	}

	user, err := us.UserUseCase.CreateUser(userDomain, ctx)
	if err != nil {
		return nil, fmt.Errorf("user with login %s create fail", userDomain.Login)
	}

	return &proto.CreateUserReply{User: converters.UserConvertCoreInProto(user)}, nil
}

// GetUserByVKId get vkId.
func (us *UserServer) GetUserByVKId(ctx context.Context, input *proto.GetUserVKIdRequest) (*proto.GetUserReply, error) {
	if input.VkId <= 0 {
		return nil, fmt.Errorf("bad vkId")
	}

	user, err := us.UserUseCase.GetUserVkID(input.VkId, ctx)
	if err != nil {
		fmt.Println("user with vkId create fail")
		return nil, fmt.Errorf("user with vkId create fail")
	}

	return &proto.GetUserReply{User: converters.UserConvertCoreInProto(user)}, nil
}

// GetUserByOnlyLogin retrieves information about a user by their login.
func (us *UserServer) GetUserByOnlyLogin(ctx context.Context, input *proto.GetUserByOnlyLoginRequest) (*proto.GetUserByOnlyLoginReply, error) {
	if strings.TrimSpace(input.Login) == "" {
		return nil, fmt.Errorf("invalid user login: %s", input.Login)
	}

	user, err := us.UserUseCase.GetUserByOnlyLogin(input.Login, ctx)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &proto.GetUserByOnlyLoginReply{User: converters.UserConvertCoreInProto(user)}, nil
}

// CreateUserOtherMail creates user.
func (us *UserServer) CreateUserOtherMail(ctx context.Context, input *proto.CreateUserRequest) (*proto.CreateUserReply, error) {
	userDomain := converters.UserConvertProtoInCore(input.User)

	userDomain.Login = sanitize.SanitizeString(userDomain.Login)
	userDomain.FirstName = sanitize.SanitizeString(userDomain.FirstName)
	userDomain.Surname = sanitize.SanitizeString(userDomain.Surname)
	userDomain.Patronymic = sanitize.SanitizeString(userDomain.Patronymic)
	userDomain.AvatarID = sanitize.SanitizeString(userDomain.AvatarID)
	userDomain.PhoneNumber = sanitize.SanitizeString(userDomain.PhoneNumber)
	userDomain.Description = sanitize.SanitizeString(userDomain.Description)

	if validUtil.IsEmpty(userDomain.Login) || validUtil.IsEmpty(userDomain.FirstName) || validUtil.IsEmpty(userDomain.Surname) {
		return nil, fmt.Errorf("all fields must be filled in")
	}

	user, err := us.UserUseCase.CreateUser(userDomain, ctx)
	if err != nil {
		return nil, fmt.Errorf("user with login %s create fail", userDomain.Login)
	}

	return &proto.CreateUserReply{User: converters.UserConvertCoreInProto(user)}, nil
}
