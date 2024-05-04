//go:generate mockgen -source=./iuser.go -destination=../mock/user_mock.go -package=mock

package _interface

import (
	"context"

	"mail/internal/microservice/user/proto"
)

// UserServer represents the interface for working with users.
type UserServer interface {
	// GetUsers retrieves information about users.
	GetUsers(ctx context.Context, input *proto.GetUsersRequest) (*proto.GetUsersReply, error)

	// GetUser retrieves information about a user by their identifier.
	GetUser(ctx context.Context, input *proto.GetUserRequest) (*proto.GetUserReply, error)

	// GetUserByLogin retrieves information about a user by their login.
	GetUserByLogin(ctx context.Context, input *proto.GetUserByLoginRequest) (*proto.GetUserByLoginReply, error)

	// IsLoginUnique checks if the provided login is unique among all users.
	IsLoginUnique(ctx context.Context, input *proto.IsLoginUniqueRequest) (*proto.IsLoginUniqueReply, error)

	// UpdateUser updates user information.
	UpdateUser(ctx context.Context, input *proto.UpdateUserRequest) (*proto.UpdateUserReply, error)

	// DeleteUserById deletes a user by their identifier.
	DeleteUserById(ctx context.Context, input *proto.DeleteUserByIdRequest) (*proto.DeleteUserByIdReply, error)

	// UploadUserAvatar uploads a user's avatar.
	UploadUserAvatar(ctx context.Context, input *proto.UploadUserAvatarRequest) (*proto.UploadUserAvatarReply, error)

	// DeleteUserAvatar deletes a user's avatar.
	DeleteUserAvatar(ctx context.Context, input *proto.DeleteUserAvatarRequest) (*proto.DeleteUserAvatarReply, error)

	// CreateUser creates user.
	CreateUser(ctx context.Context, input *proto.CreateUserRequest) (*proto.CreateUserReply, error)
}
