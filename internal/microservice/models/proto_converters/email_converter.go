package proto_converters

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	grpc "mail/internal/microservice/email/proto"
	domain "mail/internal/microservice/models/domain_models"
)

// EmailConvertCoreInProto converts an email model from the application core to the gRPC format.
func EmailConvertCoreInProto(emailModelCore *domain.Email) *grpc.Email {
	return &grpc.Email{
		Id:             emailModelCore.ID,
		Topic:          emailModelCore.Topic,
		Text:           emailModelCore.Text,
		PhotoID:        emailModelCore.PhotoID,
		ReadStatus:     emailModelCore.ReadStatus,
		Flag:           emailModelCore.Flag,
		Deleted:        emailModelCore.Deleted,
		DateOfDispatch: timestamppb.New(emailModelCore.DateOfDispatch),
		ReplyToEmailID: emailModelCore.ReplyToEmailID,
		DraftStatus:    emailModelCore.DraftStatus,
		SpamStatus:     emailModelCore.SpamStatus,
		SenderEmail:    emailModelCore.SenderEmail,
		RecipientEmail: emailModelCore.RecipientEmail,
	}
}

// EmailConvertProtoInCore converts an email model from the gRPC format to the application core.
func EmailConvertProtoInCore(emailModelProto *grpc.Email) *domain.Email {
	return &domain.Email{
		ID:             emailModelProto.Id,
		Topic:          emailModelProto.Topic,
		Text:           emailModelProto.Text,
		PhotoID:        emailModelProto.PhotoID,
		ReadStatus:     emailModelProto.ReadStatus,
		Flag:           emailModelProto.Flag,
		Deleted:        emailModelProto.Deleted,
		DateOfDispatch: emailModelProto.DateOfDispatch.AsTime(),
		ReplyToEmailID: emailModelProto.ReplyToEmailID,
		DraftStatus:    emailModelProto.DraftStatus,
		SpamStatus:     emailModelProto.SpamStatus,
		SenderEmail:    emailModelProto.SenderEmail,
		RecipientEmail: emailModelProto.RecipientEmail,
	}
}

// EmailsConvertProtoInCore converts a list of email models from the gRPC format to the application core.
func EmailsConvertProtoInCore(emailModelProto *grpc.Emails) []*domain.Email {
	emailsCore := make([]*domain.Email, 0, len(emailModelProto.Emails))
	for _, email := range emailModelProto.Emails {
		emailsCore = append(emailsCore, EmailConvertProtoInCore(email))
	}
	return emailsCore
}
