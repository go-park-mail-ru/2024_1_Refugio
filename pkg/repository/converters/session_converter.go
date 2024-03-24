package converters

import (
	sessionCore "mail/pkg/domain/models"
	sessionDb "mail/pkg/repository/models"
)

func SessionConvertDbInCore(sessionModelDb sessionDb.Session) *sessionCore.Session {
	return &sessionCore.Session{
		ID:           sessionModelDb.ID,
		UserID:       sessionModelDb.UserID,
		CreationDate: sessionModelDb.CreationDate,
		Device:       sessionModelDb.Device,
		LifeTime:     sessionModelDb.LifeTime,
	}
}

func SessionConvertCoreInDb(sessionModelCore sessionCore.Session) *sessionDb.Session {
	return &sessionDb.Session{
		ID:           sessionModelCore.ID,
		UserID:       sessionModelCore.UserID,
		CreationDate: sessionModelCore.CreationDate,
		Device:       sessionModelCore.Device,
		LifeTime:     sessionModelCore.LifeTime,
	}
}
