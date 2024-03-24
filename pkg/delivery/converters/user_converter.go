package converters

import (
	api "mail/pkg/delivery/models"
	domain "mail/pkg/domain/models"
)

func UserConvertCoreInApi(userModelCore domain.User) *api.User {
	return &api.User{
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
}

func UserConvertApiInCore(userModelApi api.User) *domain.User {
	return &domain.User{
		ID:          userModelApi.ID,
		Login:       userModelApi.Login,
		FirstName:   userModelApi.FirstName,
		Surname:     userModelApi.Surname,
		Patronymic:  userModelApi.Patronymic,
		Gender:      userModelApi.Gender,
		Birthday:    userModelApi.Birthday,
		AvatarID:    userModelApi.AvatarID,
		PhoneNumber: userModelApi.PhoneNumber,
		Description: userModelApi.Description,
	}
}
