package delivery_converters

import (
	sessionCore "mail/internal/microservice/models/domain_models"
	sessionApi "mail/internal/models/delivery_models"
)

// SessionConvertCoreInApi converts a session model from the core package to the API representation.
func SessionConvertCoreInApi(sessionModelCore *sessionCore.Session) *sessionApi.Session {
	return &sessionApi.Session{
		ID:           sessionModelCore.ID,
		UserID:       sessionModelCore.UserID,
		CreationDate: sessionModelCore.CreationDate,
		Device:       sessionModelCore.Device,
		LifeTime:     sessionModelCore.LifeTime,
		CsrfToken:    sessionModelCore.CsrfToken,
	}
}

// SessionConvertApiInCore converts a session model from the API representation to the core package.
func SessionConvertApiInCore(sessionModelApi *sessionApi.Session) *sessionCore.Session {
	return &sessionCore.Session{
		ID:           sessionModelApi.ID,
		UserID:       sessionModelApi.UserID,
		CreationDate: sessionModelApi.CreationDate,
		Device:       sessionModelApi.Device,
		LifeTime:     sessionModelApi.LifeTime,
		CsrfToken:    sessionModelApi.CsrfToken,
	}
}
