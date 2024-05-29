package repository_converters

import (
	"testing"

	"github.com/stretchr/testify/assert"

	domain "mail/internal/microservice/models/domain_models"
	database "mail/internal/microservice/models/repository_models"
)

func TestEmailConvertDbInCore(t *testing.T) {
	t.Parallel()

	emailModelDb := database.Email{
		ID:             1,
		Topic:          "Test Email",
		Text:           "This is a test email.",
		ReadStatus:     false,
		Flag:           false,
		Deleted:        false,
		DraftStatus:    false,
		SenderEmail:    "sender@example.com",
		RecipientEmail: "recipient@example.com",
	}

	expected := &domain.Email{
		ID:             1,
		Topic:          "Test Email",
		Text:           "This is a test email.",
		ReadStatus:     false,
		Flag:           false,
		Deleted:        false,
		DraftStatus:    false,
		SenderEmail:    "sender@example.com",
		RecipientEmail: "recipient@example.com",
	}

	actual := EmailConvertDbInCore(&emailModelDb)

	assert.Equal(t, expected, actual)
}

func TestEmailConvertCoreInDb(t *testing.T) {
	t.Parallel()

	emailModelCore := domain.Email{
		ID:             1,
		Topic:          "Test Email",
		Text:           "This is a test email.",
		ReadStatus:     false,
		Flag:           false,
		Deleted:        false,
		DraftStatus:    false,
		SenderEmail:    "sender@example.com",
		RecipientEmail: "recipient@example.com",
	}

	actual := EmailConvertCoreInDb(&emailModelCore)

	assert.NotNil(t, actual)
}
