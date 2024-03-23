package converters

import (
	emailCore "mail/pkg/domain/models"
	emailDb "mail/pkg/repository/models"
)

func EmailConvertDbInCore(emailModelDb emailDb.Email) *emailCore.Email {
	return &emailCore.Email{
		ID:             emailModelDb.ID,
		Topic:          emailModelDb.Topic,
		Text:           emailModelDb.Text,
		PhotoID:        emailModelDb.PhotoID,
		ReadStatus:     emailModelDb.ReadStatus,
		Mark:           emailModelDb.Mark,
		Deleted:        emailModelDb.Deleted,
		DateOfDispatch: emailModelDb.DateOfDispatch,
		ReplyToEmailID: emailModelDb.ReplyToEmailID,
		DraftStatus:    emailModelDb.DraftStatus,
	}
}

func EmailConvertCoreInDb(emailModelCore emailCore.Email) *emailDb.Email {
	return &emailDb.Email{
		ID:             emailModelCore.ID,
		Topic:          emailModelCore.Topic,
		Text:           emailModelCore.Text,
		PhotoID:        emailModelCore.PhotoID,
		ReadStatus:     emailModelCore.ReadStatus,
		Mark:           emailModelCore.Mark,
		Deleted:        emailModelCore.Deleted,
		DateOfDispatch: emailModelCore.DateOfDispatch,
		ReplyToEmailID: emailModelCore.ReplyToEmailID,
		DraftStatus:    emailModelCore.DraftStatus,
	}
}
