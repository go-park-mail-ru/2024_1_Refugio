package proto_converters

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	grpc "mail/internal/microservice/folder/proto"
	domain "mail/internal/microservice/models/domain_models"
)

func TestFolderConvertCoreInProto(t *testing.T) {
	folderModelCore := domain.Folder{
		ID:        1,
		Name:      "Inbox",
		ProfileId: 123,
	}

	expectedProto := &grpc.Folder{
		Id:        1,
		Name:      "Inbox",
		ProfileId: 123,
	}

	actualProto := FolderConvertCoreInProto(&folderModelCore)
	assert.Equal(t, expectedProto, actualProto)
}

func TestFolderConvertProtoInCore(t *testing.T) {
	folderModelProto := grpc.Folder{
		Id:        1,
		Name:      "Inbox",
		ProfileId: 123,
	}

	expectedCore := &domain.Folder{
		ID:        1,
		Name:      "Inbox",
		ProfileId: 123,
	}

	actualCore := FolderConvertProtoInCore(&folderModelProto)
	assert.Equal(t, expectedCore, actualCore)
}

func TestFoldersConvertProtoInCore(t *testing.T) {
	folderModelProto := &grpc.Folders{
		Folders: []*grpc.Folder{
			{
				Id:        1,
				Name:      "Inbox",
				ProfileId: 123,
			},
			{
				Id:        2,
				Name:      "Sent",
				ProfileId: 123,
			},
		},
	}

	expectedCore := []*domain.Folder{
		{
			ID:        1,
			Name:      "Inbox",
			ProfileId: 123,
		},
		{
			ID:        2,
			Name:      "Sent",
			ProfileId: 123,
		},
	}

	actualCore := FoldersConvertProtoInCore(folderModelProto)
	assert.Equal(t, expectedCore, actualCore)
}

func TestObjectEmailConvertProtoInCore(t *testing.T) {
	folderEmailModelProto := grpc.ObjectEmail{
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
		DateOfDispatch: folderEmailModelProto.DateOfDispatch.AsTime(),
		ReplyToEmailID: 0,
		DraftStatus:    false,
		SpamStatus:     false,
		SenderEmail:    "sender@example.com",
		RecipientEmail: "recipient@example.com",
	}

	actualCore := ObjectEmailConvertProtoInCore(&folderEmailModelProto)
	assert.Equal(t, expectedCore, actualCore)
}

func TestObjectEmailConvertCoreInProto(t *testing.T) {
	folderEmailModelCore := domain.Email{
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

	expectedProto := &grpc.ObjectEmail{
		Id:             1,
		Topic:          "Test Email",
		Text:           "This is a test email.",
		PhotoID:        "photo123",
		ReadStatus:     true,
		Flag:           false,
		Deleted:        false,
		DateOfDispatch: timestamppb.New(folderEmailModelCore.DateOfDispatch),
		ReplyToEmailID: 0,
		DraftStatus:    false,
		SpamStatus:     false,
		SenderEmail:    "sender@example.com",
		RecipientEmail: "recipient@example.com",
	}

	actualProto := ObjectEmailConvertCoreInProto(&folderEmailModelCore)
	assert.Equal(t, expectedProto, actualProto)
}

func TestObjectsEmailConvertProtoInCore(t *testing.T) {
	folderEmailsModelProto := &grpc.ObjectsEmail{
		Emails: []*grpc.ObjectEmail{
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
			DateOfDispatch: folderEmailsModelProto.Emails[0].DateOfDispatch.AsTime(),
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
			DateOfDispatch: folderEmailsModelProto.Emails[1].DateOfDispatch.AsTime(),
			ReplyToEmailID: 0,
			DraftStatus:    true,
			SpamStatus:     false,
			SenderEmail:    "sender@example.com",
			RecipientEmail: "recipient@example.com",
		},
	}

	actualCore := ObjectsEmailConvertProtoInCore(folderEmailsModelProto)
	assert.Equal(t, expectedCore, actualCore)
}
