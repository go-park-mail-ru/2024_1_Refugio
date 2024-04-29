package delivery_converters

import (
	emailCore "mail/internal/microservice/models/domain_models"
	emailApi "mail/internal/models/delivery_models"
	"reflect"
	"testing"
)

func TestEmailConvertCoreInApi(t *testing.T) {
	emailModelCore := emailCore.Email{
		ID:             1,
		Topic:          "Тема письма",
		Text:           "Это тестовое письмо",
		PhotoID:        "photo_id",
		ReadStatus:     false,
		Flag:           true,
		Deleted:        false,
		ReplyToEmailID: 2,
		DraftStatus:    false,
		SenderEmail:    "sender@example.com",
		RecipientEmail: "recipient@example.com",
	}

	emailModelApi := EmailConvertCoreInApi(emailModelCore)

	expectedEmailModelApi := &emailApi.Email{
		ID:             emailModelCore.ID,
		Topic:          emailModelCore.Topic,
		Text:           emailModelCore.Text,
		PhotoID:        emailModelCore.PhotoID,
		ReadStatus:     emailModelCore.ReadStatus,
		Flag:           emailModelCore.Flag,
		Deleted:        emailModelCore.Deleted,
		DateOfDispatch: emailModelCore.DateOfDispatch,
		ReplyToEmailID: emailModelCore.ReplyToEmailID,
		DraftStatus:    emailModelCore.DraftStatus,
		SenderEmail:    emailModelCore.SenderEmail,
		RecipientEmail: emailModelCore.RecipientEmail,
	}

	if !reflect.DeepEqual(emailModelApi, expectedEmailModelApi) {
		t.Errorf("EmailConvertCoreInApi() = %v, ожидалось %v", emailModelApi, expectedEmailModelApi)
	}
}

func TestEmailConvertApiInCore(t *testing.T) {
	emailModelApi := emailApi.Email{
		ID:             1,
		Topic:          "Тема письма",
		Text:           "Это тестовое письмо",
		PhotoID:        "photo_id",
		ReadStatus:     false,
		Flag:           true,
		Deleted:        false,
		ReplyToEmailID: 2,
		DraftStatus:    false,
		SenderEmail:    "sender@example.com",
		RecipientEmail: "recipient@example.com",
	}

	emailModelCore := EmailConvertApiInCore(emailModelApi)

	expectedEmailModelCore := &emailCore.Email{
		ID:             emailModelApi.ID,
		Topic:          emailModelApi.Topic,
		Text:           emailModelApi.Text,
		PhotoID:        emailModelApi.PhotoID,
		ReadStatus:     emailModelApi.ReadStatus,
		Flag:           emailModelApi.Flag,
		Deleted:        emailModelApi.Deleted,
		DateOfDispatch: emailModelApi.DateOfDispatch,
		ReplyToEmailID: emailModelApi.ReplyToEmailID,
		DraftStatus:    emailModelApi.DraftStatus,
		SenderEmail:    emailModelApi.SenderEmail,
		RecipientEmail: emailModelApi.RecipientEmail,
	}

	if !reflect.DeepEqual(emailModelCore, expectedEmailModelCore) {
		t.Errorf("EmailConvertApiInCore() = %v, ожидалось %v", emailModelCore, expectedEmailModelCore)
	}
}
