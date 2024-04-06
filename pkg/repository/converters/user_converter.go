package converters

import (
	domain "mail/pkg/domain/models"
	database "mail/pkg/repository/models"
)

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
