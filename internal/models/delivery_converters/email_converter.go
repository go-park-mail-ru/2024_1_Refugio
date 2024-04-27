package delivery_converters

import (
	emailCore "mail/internal/microservice/models/domain_models"
	emailApi "mail/internal/models/delivery_models"
)

// EmailConvertCoreInApi converts an email model from the core package to the API representation.
func EmailConvertCoreInApi(emailModelDb emailCore.Email) *emailApi.Email {
	return &emailApi.Email{
		ID:             emailModelDb.ID,
		Topic:          emailModelDb.Topic,
		Text:           emailModelDb.Text,
		PhotoID:        emailModelDb.PhotoID,
		ReadStatus:     emailModelDb.ReadStatus,
		Flag:           emailModelDb.Flag,
		Deleted:        emailModelDb.Deleted,
		DateOfDispatch: emailModelDb.DateOfDispatch,
		ReplyToEmailID: emailModelDb.ReplyToEmailID,
		DraftStatus:    emailModelDb.DraftStatus,
		SpamStatus:     emailModelDb.SpamStatus,
		SenderEmail:    emailModelDb.SenderEmail,
		RecipientEmail: emailModelDb.RecipientEmail,
	}
}

// EmailConvertApiInCore converts an email model from the API representation to the core package.
func EmailConvertApiInCore(emailModelApi emailApi.Email) *emailCore.Email {
	return &emailCore.Email{
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
		SpamStatus:     emailModelApi.SpamStatus,
		SenderEmail:    emailModelApi.SenderEmail,
		RecipientEmail: emailModelApi.RecipientEmail,
	}
}
