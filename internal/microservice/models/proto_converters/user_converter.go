package proto_converters

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	domain "mail/internal/microservice/models/domain_models"
	grpc "mail/internal/microservice/user/proto"
)

// UserConvertCoreInProto converts a user model from the application core to the gRPC format.
func UserConvertCoreInProto(userModelCore *domain.User) *grpc.User {
	return &grpc.User{
		Id:          userModelCore.ID,
		Login:       userModelCore.Login,
		Password:    userModelCore.Password,
		Firstname:   userModelCore.FirstName,
		Surname:     userModelCore.Surname,
		Patronymic:  userModelCore.Patronymic,
		Gender:      domain.GetGender(userModelCore.Gender),
		Birthday:    timestamppb.New(userModelCore.Birthday),
		Avatar:      userModelCore.AvatarID,
		PhoneNumber: userModelCore.PhoneNumber,
		Description: userModelCore.Description,
	}
}

// UserConvertProtoInCore converts a user model from the gRPC format to the application core.
func UserConvertProtoInCore(userModelProto *grpc.User) *domain.User {
	return &domain.User{
		ID:          userModelProto.Id,
		Login:       userModelProto.Login,
		Password:    userModelProto.Password,
		FirstName:   userModelProto.Firstname,
		Surname:     userModelProto.Surname,
		Patronymic:  userModelProto.Patronymic,
		Gender:      domain.UserGender(userModelProto.Gender),
		Birthday:    userModelProto.Birthday.AsTime(),
		AvatarID:    userModelProto.Avatar,
		PhoneNumber: userModelProto.PhoneNumber,
		Description: userModelProto.Description,
		VKId:        userModelProto.VkId,
	}
}
