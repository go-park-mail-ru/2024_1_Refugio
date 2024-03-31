package converters

import (
	domain "mail/pkg/domain/models"
	emailDb "mail/pkg/repository/models"
	"reflect"
	"testing"
)

func TestEmailConvertDbInCore(t *testing.T) {
	emailModelDb := emailDb.Email{
		ID:             1,
		Topic:          "Тема письма",
		Text:           "Это тестовое письмо",
		PhotoID:        "photo_id",
		ReadStatus:     false,
		Flag:           true,
		Deleted:        false,
		DraftStatus:    false,
		SenderEmail:    "sender@example.com",
		RecipientEmail: "recipient@example.com",
	}

	emailModelCore := EmailConvertDbInCore(emailModelDb)

	expectedEmailModelCore := &domain.Email{
		ID:             emailModelDb.ID,
		Topic:          emailModelDb.Topic,
		Text:           emailModelDb.Text,
		PhotoID:        emailModelDb.PhotoID,
		ReadStatus:     emailModelDb.ReadStatus,
		Flag:           emailModelDb.Flag,
		Deleted:        emailModelDb.Deleted,
		DateOfDispatch: emailModelDb.DateOfDispatch,
		DraftStatus:    emailModelDb.DraftStatus,
		SenderEmail:    emailModelDb.SenderEmail,
		RecipientEmail: emailModelDb.RecipientEmail,
	}

	if !reflect.DeepEqual(emailModelCore, expectedEmailModelCore) {
		t.Errorf("EmailConvertDbInCore() = %v, ожидалось %v", emailModelCore, expectedEmailModelCore)
	}
}

func TestEmailConvertCoreInDb(t *testing.T) {
	emailModelCore := domain.Email{
		ID:             1,
		Topic:          "Тема письма",
		Text:           "Это тестовое письмо",
		PhotoID:        "photo_id",
		ReadStatus:     false,
		Flag:           true,
		Deleted:        false,
		DraftStatus:    false,
		SenderEmail:    "sender@example.com",
		RecipientEmail: "recipient@example.com",
	}

	emailModelDb := EmailConvertCoreInDb(emailModelCore)

	expectedEmailModelDb := &emailDb.Email{
		ID:             emailModelCore.ID,
		Topic:          emailModelCore.Topic,
		Text:           emailModelCore.Text,
		PhotoID:        emailModelCore.PhotoID,
		ReadStatus:     emailModelCore.ReadStatus,
		Flag:           emailModelCore.Flag,
		Deleted:        emailModelCore.Deleted,
		DateOfDispatch: emailModelCore.DateOfDispatch,
		DraftStatus:    emailModelCore.DraftStatus,
		SenderEmail:    emailModelCore.SenderEmail,
		RecipientEmail: emailModelCore.RecipientEmail,
	}

	if !reflect.DeepEqual(emailModelDb, expectedEmailModelDb) {
		t.Errorf("EmailConvertCoreInDb() = %v, ожидалось %v", emailModelDb, expectedEmailModelDb)
	}
}
