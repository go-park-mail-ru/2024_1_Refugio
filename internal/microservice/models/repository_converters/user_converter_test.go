package repository_converters

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	domain "mail/internal/microservice/models/domain_models"
	database "mail/internal/microservice/models/repository_models"
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
		PhoneNumber: "+1234567890",
		Description: "Test description",
	}

	userCore := UserConvertDbInCore(&userModelDb)

	expectedUserCore := &domain.User{
		ID:          userModelDb.ID,
		Login:       userModelDb.Login,
		FirstName:   userModelDb.FirstName,
		Surname:     userModelDb.Surname,
		Patronymic:  userModelDb.Patronymic,
		Gender:      userModelDb.Gender,
		Birthday:    userModelDb.Birthday,
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
		PhoneNumber: "+1234567890",
		Description: "Test description",
	}

	userDb := UserConvertCoreInDb(&userModelCore)

	assert.NotNil(t, userDb)
}
