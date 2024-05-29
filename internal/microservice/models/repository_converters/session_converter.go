package repository_converters

import (
	sessionCore "mail/internal/microservice/models/domain_models"
	sessionDb "mail/internal/microservice/models/repository_models"
)

// SessionConvertDbInCore converts a session model from database representation to core domain representation.
func SessionConvertDbInCore(sessionModelDb *sessionDb.Session) *sessionCore.Session {
	return &sessionCore.Session{
		ID:           sessionModelDb.ID,
		UserID:       sessionModelDb.UserID,
		CreationDate: sessionModelDb.CreationDate,
		Device:       sessionModelDb.Device,
		LifeTime:     sessionModelDb.LifeTime,
		CsrfToken:    sessionModelDb.CsrfToken,
	}
}

// SessionConvertCoreInDb converts a session model from core domain representation to database representation.
func SessionConvertCoreInDb(sessionModelCore *sessionCore.Session) *sessionDb.Session {
	return &sessionDb.Session{
		ID:           sessionModelCore.ID,
		UserID:       sessionModelCore.UserID,
		CreationDate: sessionModelCore.CreationDate,
		Device:       sessionModelCore.Device,
		LifeTime:     sessionModelCore.LifeTime,
		CsrfToken:    sessionModelCore.CsrfToken,
	}
}
