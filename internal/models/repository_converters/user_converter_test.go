package repository_converters

import (
	domain "mail/internal/models/domain_models"
	database "mail/internal/models/repository_models"
	"reflect"
	"testing"
	"time"
)

func TestUserConvertDbInCore(t *testing.T) {
	userModelDb := database.User{
		ID:          1,
		Login:       "test_user",
		FirstName:   "John",
		Surname:     "Doe",
		Patronymic:  "Smith",
		Gender:      "Male",
		Birthday:    time.Now(),
		AvatarID:    "avatar_123",
		PhoneNumber: "+1234567890",
		Description: "Test description",
	}

	userCore := UserConvertDbInCore(userModelDb)

	expectedUserCore := &domain.User{
		ID:          userModelDb.ID,
		Login:       userModelDb.Login,
		FirstName:   userModelDb.FirstName,
		Surname:     userModelDb.Surname,
		Patronymic:  userModelDb.Patronymic,
		Gender:      userModelDb.Gender,
		Birthday:    userModelDb.Birthday,
		AvatarID:    userModelDb.AvatarID,
		PhoneNumber: userModelDb.PhoneNumber,
		Description: userModelDb.Description,
	}

	if !reflect.DeepEqual(userCore, expectedUserCore) {
		t.Errorf("UserConvertDbInCore() = %v, want %v", userCore, expectedUserCore)
	}
}

func TestUserConvertCoreInDb(t *testing.T) {
	userModelCore := domain.User{
		ID:          1,
		Login:       "test_user",
		Password:    "password123",
		FirstName:   "John",
		Surname:     "Doe",
		Patronymic:  "Smith",
		Gender:      "Male",
		Birthday:    time.Now(),
		AvatarID:    "avatar_123",
		PhoneNumber: "+1234567890",
		Description: "Test description",
	}

	userDb := UserConvertCoreInDb(userModelCore)

	expectedUserDb := &database.User{
		ID:          userModelCore.ID,
		Login:       userModelCore.Login,
		Password:    userModelCore.Password,
		FirstName:   userModelCore.FirstName,
		Surname:     userModelCore.Surname,
		Patronymic:  userModelCore.Patronymic,
		Gender:      userModelCore.Gender,
		Birthday:    userModelCore.Birthday,
		AvatarID:    userModelCore.AvatarID,
		PhoneNumber: userModelCore.PhoneNumber,
		Description: userModelCore.Description,
	}

	if !reflect.DeepEqual(userDb, expectedUserDb) {
		t.Errorf("UserConvertCoreInDb() = %v, want %v", userDb, expectedUserDb)
	}
}
