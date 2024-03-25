package converters

import (
	sessionApi "mail/pkg/delivery/models"
	sessionCore "mail/pkg/domain/models"
)

func SessionConvertCoreInApi(sessionModelCore sessionCore.Session) *sessionApi.Session {
	return &sessionApi.Session{
		ID:           sessionModelCore.ID,
		UserID:       sessionModelCore.UserID,
		CreationDate: sessionModelCore.CreationDate,
		Device:       sessionModelCore.Device,
		LifeTime:     sessionModelCore.LifeTime,
		CsrfToken:    sessionModelCore.CsrfToken,
	}
}

func SessionConvertApiInCore(sessionModelApi sessionApi.Session) *sessionCore.Session {
	return &sessionCore.Session{
		ID:           sessionModelApi.ID,
		UserID:       sessionModelApi.UserID,
		CreationDate: sessionModelApi.CreationDate,
		Device:       sessionModelApi.Device,
		LifeTime:     sessionModelApi.LifeTime,
		CsrfToken:    sessionModelApi.CsrfToken,
	}
}
