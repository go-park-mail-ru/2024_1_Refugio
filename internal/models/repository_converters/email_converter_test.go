package repository_converters

import (
	"github.com/stretchr/testify/assert"
	domain "mail/internal/models/domain_models"
	database "mail/internal/models/repository_models"
	"testing"
)

func TestEmailConvertDbInCore(t *testing.T) {
	t.Parallel()

	emailModelDb := database.Email{
		ID:             1,
		Topic:          "Test Email",
		Text:           "This is a test email.",
		PhotoID:        "url",
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
		PhotoID:        "url",
		ReadStatus:     false,
		Flag:           false,
		Deleted:        false,
		DraftStatus:    false,
		SenderEmail:    "sender@example.com",
		RecipientEmail: "recipient@example.com",
	}

	actual := EmailConvertDbInCore(emailModelDb)

	assert.Equal(t, expected, actual)
}

func TestEmailConvertCoreInDb(t *testing.T) {
	t.Parallel()

	emailModelCore := domain.Email{
		ID:             1,
		Topic:          "Test Email",
		Text:           "This is a test email.",
		PhotoID:        "url",
		ReadStatus:     false,
		Flag:           false,
		Deleted:        false,
		DraftStatus:    false,
		SenderEmail:    "sender@example.com",
		RecipientEmail: "recipient@example.com",
	}

	expected := &database.Email{
		ID:             1,
		Topic:          "Test Email",
		Text:           "This is a test email.",
		PhotoID:        "url",
		ReadStatus:     false,
		Flag:           false,
		Deleted:        false,
		DraftStatus:    false,
		SenderEmail:    "sender@example.com",
		RecipientEmail: "recipient@example.com",
	}

	actual := EmailConvertCoreInDb(emailModelCore)

	assert.Equal(t, expected, actual)
}
