package repository_converters

import (
	domain "mail/internal/models/domain_models"
	database "mail/internal/models/repository_models"
)

// UserConvertDbInCore converts a user model from database representation to core domain representation.
func UserConvertDbInCore(userModelDb database.User) *domain.User {
	return &domain.User{
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
}

// UserConvertCoreInDb converts a user model from core domain representation to database representation.
func UserConvertCoreInDb(userModelCore domain.User) *database.User {
	return &database.User{
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
}
