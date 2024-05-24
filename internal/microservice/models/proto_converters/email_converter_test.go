package proto_converters

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	grpc "mail/internal/microservice/email/proto"
	domain "mail/internal/microservice/models/domain_models"
)

func TestEmailConvertCoreInProto(t *testing.T) {
	emailModelCore := domain.Email{
		ID:             1,
		Topic:          "Test Email",
		Text:           "This is a test email.",
		PhotoID:        "photo123",
		ReadStatus:     true,
		Flag:           false,
		Deleted:        false,
		DateOfDispatch: time.Now(),
		ReplyToEmailID: 0,
		DraftStatus:    false,
		SpamStatus:     false,
		SenderEmail:    "sender@example.com",
		RecipientEmail: "recipient@example.com",
	}

	expectedProto := &grpc.Email{
		Id:             1,
		Topic:          "Test Email",
		Text:           "This is a test email.",
		PhotoID:        "photo123",
		ReadStatus:     true,
		Flag:           false,
		Deleted:        false,
		DateOfDispatch: timestamppb.New(emailModelCore.DateOfDispatch),
		ReplyToEmailID: 0,
		DraftStatus:    false,
		SpamStatus:     false,
		SenderEmail:    "sender@example.com",
		RecipientEmail: "recipient@example.com",
	}

	actualProto := EmailConvertCoreInProto(&emailModelCore)
	assert.Equal(t, expectedProto, actualProto)
}

func TestEmailConvertProtoInCore(t *testing.T) {
	emailModelProto := grpc.Email{
		Id:             1,
		Topic:          "Test Email",
		Text:           "This is a test email.",
		PhotoID:        "photo123",
		ReadStatus:     true,
		Flag:           false,
		Deleted:        false,
		DateOfDispatch: timestamppb.Now(),
		ReplyToEmailID: 0,
		DraftStatus:    false,
		SpamStatus:     false,
		SenderEmail:    "sender@example.com",
		RecipientEmail: "recipient@example.com",
	}

	expectedCore := &domain.Email{
		ID:             1,
		Topic:          "Test Email",
		Text:           "This is a test email.",
		PhotoID:        "photo123",
		ReadStatus:     true,
		Flag:           false,
		Deleted:        false,
		DateOfDispatch: emailModelProto.DateOfDispatch.AsTime(),
		ReplyToEmailID: 0,
		DraftStatus:    false,
		SpamStatus:     false,
		SenderEmail:    "sender@example.com",
		RecipientEmail: "recipient@example.com",
	}

	actualCore := EmailConvertProtoInCore(&emailModelProto)
	assert.Equal(t, expectedCore, actualCore)
}

func TestEmailsConvertProtoInCore(t *testing.T) {
	emailModelProto := &grpc.Emails{
		Emails: []*grpc.Email{
			{
				Id:             1,
				Topic:          "Test Email 1",
				Text:           "This is test email 1.",
				PhotoID:        "photo123",
				ReadStatus:     true,
				Flag:           false,
				Deleted:        false,
				DateOfDispatch: timestamppb.Now(),
				ReplyToEmailID: 0,
				DraftStatus:    false,
				SpamStatus:     false,
				SenderEmail:    "sender@example.com",
				RecipientEmail: "recipient@example.com",
			},
			{
				Id:             2,
				Topic:          "Test Email 2",
				Text:           "This is test email 2.",
				PhotoID:        "photo456",
				ReadStatus:     false,
				Flag:           true,
				Deleted:        false,
				DateOfDispatch: timestamppb.Now(),
				ReplyToEmailID: 0,
				DraftStatus:    true,
				SpamStatus:     false,
				SenderEmail:    "sender@example.com",
				RecipientEmail: "recipient@example.com",
			},
		},
	}

	expectedCore := []*domain.Email{
		{
			ID:             1,
			Topic:          "Test Email 1",
			Text:           "This is test email 1.",
			PhotoID:        "photo123",
			ReadStatus:     true,
			Flag:           false,
			Deleted:        false,
			DateOfDispatch: emailModelProto.Emails[0].DateOfDispatch.AsTime(),
			ReplyToEmailID: 0,
			DraftStatus:    false,
			SpamStatus:     false,
			SenderEmail:    "sender@example.com",
			RecipientEmail: "recipient@example.com",
		},
		{
			ID:             2,
			Topic:          "Test Email 2",
			Text:           "This is test email 2.",
			PhotoID:        "photo456",
			ReadStatus:     false,
			Flag:           true,
			Deleted:        false,
			DateOfDispatch: emailModelProto.Emails[1].DateOfDispatch.AsTime(),
			ReplyToEmailID: 0,
			DraftStatus:    true,
			SpamStatus:     false,
			SenderEmail:    "sender@example.com",
			RecipientEmail: "recipient@example.com",
		},
	}

	actualCore := EmailsConvertProtoInCore(emailModelProto)
	assert.Equal(t, expectedCore, actualCore)
}
