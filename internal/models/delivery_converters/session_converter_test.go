package delivery_converters

import (
	domain "mail/internal/microservice/models/domain_models"
	api "mail/internal/models/delivery_models"
	"reflect"
	"testing"
	"time"
)

func TestSessionConvertCoreInApi(t *testing.T) {
	sessionModelCore := domain.Session{
		ID:           "session_id",
		UserID:       1,
		CreationDate: time.Now(),
		Device:       "desktop",
		LifeTime:     3600,
		CsrfToken:    "csrf_token",
	}

	sessionModelApi := SessionConvertCoreInApi(&sessionModelCore)

	expectedSessionModelApi := &api.Session{
		ID:           sessionModelCore.ID,
		UserID:       sessionModelCore.UserID,
		CreationDate: sessionModelCore.CreationDate,
		Device:       sessionModelCore.Device,
		LifeTime:     sessionModelCore.LifeTime,
		CsrfToken:    sessionModelCore.CsrfToken,
	}

	if !reflect.DeepEqual(sessionModelApi, expectedSessionModelApi) {
		t.Errorf("SessionConvertCoreInApi() = %v, want %v", sessionModelApi, expectedSessionModelApi)
	}
}

func TestSessionConvertApiInCore(t *testing.T) {
	sessionModelApi := api.Session{
		ID:           "session_id",
		UserID:       1,
		CreationDate: time.Now(),
		Device:       "desktop",
		LifeTime:     3600,
		CsrfToken:    "csrf_token",
	}

	sessionModelCore := SessionConvertApiInCore(&sessionModelApi)

	expectedSessionModelCore := &domain.Session{
		ID:           sessionModelApi.ID,
		UserID:       sessionModelApi.UserID,
		CreationDate: sessionModelApi.CreationDate,
		Device:       sessionModelApi.Device,
		LifeTime:     sessionModelApi.LifeTime,
		CsrfToken:    sessionModelApi.CsrfToken,
	}

	if !reflect.DeepEqual(sessionModelCore, expectedSessionModelCore) {
		t.Errorf("SessionConvertApiInCore() = %v, want %v", sessionModelCore, expectedSessionModelCore)
	}
}
