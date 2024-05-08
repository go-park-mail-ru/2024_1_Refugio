package delivery_converters

import (
	domain "mail/internal/microservice/models/domain_models"
	api "mail/internal/models/delivery_models"
)

// UserConvertCoreInApi converts a user model from the domain package to the API representation.
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

// UserConvertApiInCore converts a user model from the API representation to the domain package.
func UserConvertApiInCore(userModelApi api.User) *domain.User {
	return &domain.User{
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
}
