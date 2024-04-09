package converters

import (
	domain "mail/pkg/domain/models"
	database "mail/pkg/repository/models"
)

// EmailConvertDbInCore converts an email model from database representation to core domain representation.
func EmailConvertDbInCore(emailModelDb database.Email) *domain.Email {
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
		SenderEmail:    emailModelDb.SenderEmail,
		RecipientEmail: emailModelDb.RecipientEmail,
	}
}

// EmailConvertCoreInDb converts an email model from core domain representation to database representation.
func EmailConvertCoreInDb(emailModelCore domain.Email) *database.Email {
	emailDB := &database.Email{
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
		SenderEmail:    emailModelCore.SenderEmail,
		RecipientEmail: emailModelCore.RecipientEmail,
	}

	if emailModelCore.ReplyToEmailID != 0 {
		emailDB.ReplyToEmailID = emailModelCore.ReplyToEmailID
	}

	return emailDB
}
