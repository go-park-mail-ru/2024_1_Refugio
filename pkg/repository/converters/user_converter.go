package converters

import (
	userCore "mail/pkg/domain/models"
	userDb "mail/pkg/repository/models"
)

func UserConvertDbInCore(userModelDb userDb.User) *userCore.User {
	return &userCore.User{
		ID:       userModelDb.ID,
		Login:    userModelDb.Login,
		Password: userModelDb.Password,
		Name:     userModelDb.Name,
		Surname:  userModelDb.Surname,
		AvatarID: userModelDb.AvatarID,
	}
}

func UserConvertCoreInDb(userModelCore userCore.User) *userDb.User {
	return &userDb.User{
		ID:       userModelCore.ID,
		Login:    userModelCore.Login,
		Password: userModelCore.Password,
		Name:     userModelCore.Name,
		Surname:  userModelCore.Surname,
		AvatarID: userModelCore.AvatarID,
	}
}
