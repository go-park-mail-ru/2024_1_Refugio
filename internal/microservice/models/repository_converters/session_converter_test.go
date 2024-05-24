package repository_converters

import (
	"reflect"
	"testing"
	"time"

	domain "mail/internal/microservice/models/domain_models"
	database "mail/internal/microservice/models/repository_models"
)

func TestSessionConvertDbInCore(t *testing.T) {
	sessionModelDb := database.Session{
		ID:           "session_123",
		UserID:       1,
		CreationDate: time.Now(),
		Device:       "test_device",
		LifeTime:     3600,
		CsrfToken:    "csrf_token_123",
	}

	sessionCore := SessionConvertDbInCore(&sessionModelDb)

	expectedSessionCore := &domain.Session{
		ID:           sessionModelDb.ID,
		UserID:       sessionModelDb.UserID,
		CreationDate: sessionModelDb.CreationDate,
		Device:       sessionModelDb.Device,
		LifeTime:     sessionModelDb.LifeTime,
		CsrfToken:    sessionModelDb.CsrfToken,
	}

	if !reflect.DeepEqual(sessionCore, expectedSessionCore) {
		t.Errorf("SessionConvertDbInCore() = %v, want %v", sessionCore, expectedSessionCore)
	}
}

func TestSessionConvertCoreInDb(t *testing.T) {
	sessionModelCore := domain.Session{
		ID:           "session_123",
		UserID:       1,
		CreationDate: time.Now(),
		Device:       "test_device",
		LifeTime:     3600,
		CsrfToken:    "csrf_token_123",
	}

	sessionDb := SessionConvertCoreInDb(&sessionModelCore)

	expectedSessionDb := &database.Session{
		ID:           sessionModelCore.ID,
		UserID:       sessionModelCore.UserID,
		CreationDate: sessionModelCore.CreationDate,
		Device:       sessionModelCore.Device,
		LifeTime:     sessionModelCore.LifeTime,
		CsrfToken:    sessionModelCore.CsrfToken,
	}

	if !reflect.DeepEqual(sessionDb, expectedSessionDb) {
		t.Errorf("SessionConvertCoreInDb() = %v, want %v", sessionDb, expectedSessionDb)
	}
}
