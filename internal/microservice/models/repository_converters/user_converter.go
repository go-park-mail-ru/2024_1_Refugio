package repository_converters

import (
	domain "mail/internal/microservice/models/domain_models"
	database "mail/internal/microservice/models/repository_models"
)

// UserConvertDbInCore converts a user model from database representation to core domain representation.
func UserConvertDbInCore(userModelDb *database.User) *domain.User {
	avatar := ""
	if userModelDb.AvatarID != nil {
		avatar = *userModelDb.AvatarID
	}

	return &domain.User{
		ID:          userModelDb.ID,
		Login:       userModelDb.Login,
		FirstName:   userModelDb.FirstName,
		Surname:     userModelDb.Surname,
		Patronymic:  userModelDb.Patronymic,
		Gender:      userModelDb.Gender,
		Birthday:    userModelDb.Birthday,
		AvatarID:    avatar,
		PhoneNumber: userModelDb.PhoneNumber,
		Description: userModelDb.Description,
		VKId:        userModelDb.VKId,
	}
}

// UserConvertCoreInDb converts a user model from core domain representation to database representation.
func UserConvertCoreInDb(userModelCore *domain.User) *database.User {
	return &database.User{
		ID:          userModelCore.ID,
		Login:       userModelCore.Login,
		Password:    userModelCore.Password,
		FirstName:   userModelCore.FirstName,
		Surname:     userModelCore.Surname,
		Patronymic:  userModelCore.Patronymic,
		Gender:      userModelCore.Gender,
		Birthday:    userModelCore.Birthday,
		AvatarID:    &userModelCore.AvatarID,
		PhoneNumber: userModelCore.PhoneNumber,
		Description: userModelCore.Description,
		VKId:        userModelCore.VKId,
	}
}
