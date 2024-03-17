package converters

import (
	userApi "mail/pkg/delivery/models"
	userCore "mail/pkg/domain/models"
)

func UserConvertCoreInApi(userModelCore userCore.User) *userApi.User {
	return &userApi.User{
		ID:       userModelCore.ID,
		Login:    userModelCore.Login,
		Name:     userModelCore.Name,
		Surname:  userModelCore.Surname,
		AvatarID: userModelCore.AvatarID,
	}
}

func UserConvertApiInCore(userModelApi userApi.User) *userCore.User {
	return &userCore.User{
		ID:       userModelApi.ID,
		Login:    userModelApi.Login,
		Password: userModelApi.Password,
		Name:     userModelApi.Name,
		Surname:  userModelApi.Surname,
		AvatarID: userModelApi.AvatarID,
	}
}
