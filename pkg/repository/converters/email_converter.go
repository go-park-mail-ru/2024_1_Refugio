package converters

import (
	domain "mail/pkg/domain/models"
	emailDb "mail/pkg/repository/models"
)

func EmailConvertDbInCore(emailModelDb emailDb.Email) *domain.Email {
	var ReplyToEmailID uint64
	if emailModelDb.ReplyToEmailID == nil {
		ReplyToEmailID = 0
	} else {
		ReplyToEmailID = emailModelDb.ReplyToEmailID.(uint64)
	}
	return &domain.Email{
		ID:             emailModelDb.ID,
		Topic:          emailModelDb.Topic,
		Text:           emailModelDb.Text,
		PhotoID:        emailModelDb.PhotoID,
		ReadStatus:     emailModelDb.ReadStatus,
		Flag:           emailModelDb.Flag,
		Deleted:        emailModelDb.Deleted,
		DateOfDispatch: emailModelDb.DateOfDispatch,
		ReplyToEmailID: ReplyToEmailID,
		DraftStatus:    emailModelDb.DraftStatus,
		SenderID:       emailModelDb.SenderID,
		RecipientID:    emailModelDb.RecipientID,
	}
}

func EmailConvertCoreInDb(emailModelCore domain.Email) *emailDb.Email {
	emailDB := &emailDb.Email{
		ID:             emailModelCore.ID,
		Topic:          emailModelCore.Topic,
		Text:           emailModelCore.Text,
		PhotoID:        emailModelCore.PhotoID,
		ReadStatus:     emailModelCore.ReadStatus,
		Flag:           emailModelCore.Flag,
		Deleted:        emailModelCore.Deleted,
		DateOfDispatch: emailModelCore.DateOfDispatch,
		ReplyToEmailID: nil,
		DraftStatus:    emailModelCore.DraftStatus,
		SenderID:       emailModelCore.SenderID,
		RecipientID:    emailModelCore.RecipientID,
	}
	if emailModelCore.ReplyToEmailID != 0 {
		emailDB.ReplyToEmailID = emailModelCore.ReplyToEmailID //interface{emailModelCore.ReplyToEmailID}
	}
	return emailDB
}
