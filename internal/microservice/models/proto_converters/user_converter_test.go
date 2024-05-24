package proto_converters

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	domain "mail/internal/microservice/models/domain_models"
	grpc "mail/internal/microservice/user/proto"
)

func TestUserConvertCoreInProto(t *testing.T) {
	birthday := time.Now()
	userModelCore := domain.User{
		ID:          123,
		Login:       "john_doe",
		Password:    "password",
		FirstName:   "John",
		Surname:     "Doe",
		Patronymic:  "Middle",
		Gender:      domain.Male,
		Birthday:    birthday,
		AvatarID:    "avatar123",
		PhoneNumber: "123456789",
		Description: "Some description",
	}

	expectedProto := &grpc.User{
		Id:          123,
		Login:       "john_doe",
		Password:    "password",
		Firstname:   "John",
		Surname:     "Doe",
		Patronymic:  "Middle",
		Gender:      "Male",
		Birthday:    timestamppb.New(birthday),
		Avatar:      "avatar123",
		PhoneNumber: "123456789",
		Description: "Some description",
	}

	actualProto := UserConvertCoreInProto(&userModelCore)
	assert.Equal(t, expectedProto, actualProto)
}

func TestUserConvertProtoInCore(t *testing.T) {
	userModelProto := grpc.User{
		Id:          123,
		Login:       "john_doe",
		Password:    "password",
		Firstname:   "John",
		Surname:     "Doe",
		Patronymic:  "Middle",
		Gender:      "Male",
		Avatar:      "avatar123",
		PhoneNumber: "123456789",
		Description: "Some description",
	}

	actualCore := UserConvertProtoInCore(&userModelProto)
	assert.NotNil(t, actualCore)
}
