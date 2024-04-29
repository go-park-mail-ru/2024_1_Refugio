package delivery_converters

import (
	domain "mail/internal/microservice/models/domain_models"
	api "mail/internal/models/delivery_models"
	"reflect"
	"testing"
	"time"
)

func TestUserConvertCoreInApi(t *testing.T) {
	userModelCore := domain.User{
		ID:          1,
		Login:       "test_user",
		FirstName:   "John",
		Surname:     "Doe",
		Patronymic:  "Smith",
		Gender:      "Male",
		Birthday:    time.Date(1980, time.January, 1, 0, 0, 0, 0, time.UTC),
		AvatarID:    "avatar_123",
		PhoneNumber: "+1234567890",
		Description: "Test description",
	}

	userModelApi := UserConvertCoreInApi(userModelCore)

	expectedUserModelApi := &api.User{
		ID:          userModelCore.ID,
		Login:       userModelCore.Login,
		FirstName:   userModelCore.FirstName,
		Surname:     userModelCore.Surname,
		Patronymic:  userModelCore.Patronymic,
		Gender:      userModelCore.Gender,
		Birthday:    userModelCore.Birthday,
		AvatarID:    userModelCore.AvatarID,
		PhoneNumber: userModelCore.PhoneNumber,
		Description: userModelCore.Description,
	}

	if !reflect.DeepEqual(userModelApi, expectedUserModelApi) {
		t.Errorf("UserConvertCoreInApi() = %v, want %v", userModelApi, expectedUserModelApi)
	}
}

func TestUserConvertApiInCore(t *testing.T) {
	userModelApi := api.User{
		ID:          1,
		Login:       "test_user",
		Password:    "password",
		FirstName:   "John",
		Surname:     "Doe",
		Patronymic:  "Smith",
		Gender:      "Male",
		Birthday:    time.Date(1980, time.January, 1, 0, 0, 0, 0, time.UTC),
		AvatarID:    "avatar_123",
		PhoneNumber: "+1234567890",
		Description: "Test description",
	}

	userModelCore := UserConvertApiInCore(userModelApi)

	expectedUserModelCore := &domain.User{
		ID:          userModelApi.ID,
		Login:       userModelApi.Login,
		Password:    userModelApi.Password,
		FirstName:   userModelApi.FirstName,
		Surname:     userModelApi.Surname,
		Patronymic:  userModelApi.Patronymic,
		Gender:      userModelApi.Gender,
		Birthday:    userModelApi.Birthday,
		AvatarID:    userModelApi.AvatarID,
		PhoneNumber: userModelApi.PhoneNumber,
		Description: userModelApi.Description,
	}

	if !reflect.DeepEqual(userModelCore, expectedUserModelCore) {
		t.Errorf("UserConvertApiInCore() = %v, want %v", userModelCore, expectedUserModelCore)
	}
}
