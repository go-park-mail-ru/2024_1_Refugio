package repository_converters

import (
	domain "mail/internal/microservice/models/domain_models"
	database "mail/internal/microservice/models/repository_models"
)

// EmailConvertDbInCore converts an email model from database representation to core domain representation.
func EmailConvertDbInCore(emailModelDb *database.Email) *domain.Email {
	var ReplyToEmailID uint64

	if emailModelDb.ReplyToEmailID == nil {
		ReplyToEmailID = 0
	} else {
		ReplyToEmailID = uint64(emailModelDb.ReplyToEmailID.(int64))
	}

	var avatar string
	if emailModelDb.PhotoID != nil {
		avatar = *emailModelDb.PhotoID
	}

	return &domain.Email{
		ID:             emailModelDb.ID,
		Topic:          emailModelDb.Topic,
		Text:           emailModelDb.Text,
		PhotoID:        avatar,
		ReadStatus:     emailModelDb.ReadStatus,
		Flag:           emailModelDb.Flag,
		Deleted:        emailModelDb.Deleted,
		DateOfDispatch: emailModelDb.DateOfDispatch,
		ReplyToEmailID: ReplyToEmailID,
		DraftStatus:    emailModelDb.DraftStatus,
		SpamStatus:     emailModelDb.SpamStatus,
		SenderEmail:    emailModelDb.SenderEmail,
		RecipientEmail: emailModelDb.RecipientEmail,
	}
}

// EmailConvertCoreInDb converts an email model from core domain representation to database representation.
func EmailConvertCoreInDb(emailModelCore *domain.Email) *database.Email {
	emailDB := &database.Email{
		ID:             emailModelCore.ID,
		Topic:          emailModelCore.Topic,
		Text:           emailModelCore.Text,
		PhotoID:        &emailModelCore.PhotoID,
		ReadStatus:     emailModelCore.ReadStatus,
		Flag:           emailModelCore.Flag,
		Deleted:        emailModelCore.Deleted,
		DateOfDispatch: emailModelCore.DateOfDispatch,
		ReplyToEmailID: nil,
		DraftStatus:    emailModelCore.DraftStatus,
		SpamStatus:     emailModelCore.SpamStatus,
		SenderEmail:    emailModelCore.SenderEmail,
		RecipientEmail: emailModelCore.RecipientEmail,
	}

	if emailModelCore.ReplyToEmailID != 0 {
		emailDB.ReplyToEmailID = emailModelCore.ReplyToEmailID
	}

	if emailModelCore.PhotoID != "" {
		photoId := &emailModelCore.PhotoID
		emailDB.PhotoID = photoId
	}

	return emailDB
}
