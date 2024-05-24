package proto_converters

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	grpc "mail/internal/microservice/folder/proto"
	domain "mail/internal/microservice/models/domain_models"
)

// FolderConvertCoreInProto converts a folder model from the application core to the gRPC format.
func FolderConvertCoreInProto(folderModelCore *domain.Folder) *grpc.Folder {
	return &grpc.Folder{
		Id:        folderModelCore.ID,
		Name:      folderModelCore.Name,
		ProfileId: folderModelCore.ProfileId,
	}
}

// FolderConvertProtoInCore converts a folder model from the gRPC format to the application core.
func FolderConvertProtoInCore(folderModelProto *grpc.Folder) *domain.Folder {
	return &domain.Folder{
		ID:        folderModelProto.Id,
		Name:      folderModelProto.Name,
		ProfileId: folderModelProto.ProfileId,
	}
}

// FoldersConvertProtoInCore converts a list of folder models from the gRPC format to the application core.
func FoldersConvertProtoInCore(folderModelProto *grpc.Folders) []*domain.Folder {
	foldersCore := make([]*domain.Folder, 0, len(folderModelProto.Folders))
	for _, folder := range folderModelProto.Folders {
		foldersCore = append(foldersCore, FolderConvertProtoInCore(folder))
	}
	return foldersCore
}

// ObjectEmailConvertProtoInCore converts an email model from the gRPC format to the application core.
func ObjectEmailConvertProtoInCore(folderEmailModelProto *grpc.ObjectEmail) *domain.Email {
	return &domain.Email{
		ID:             folderEmailModelProto.Id,
		Topic:          folderEmailModelProto.Topic,
		Text:           folderEmailModelProto.Text,
		PhotoID:        folderEmailModelProto.PhotoID,
		ReadStatus:     folderEmailModelProto.ReadStatus,
		Flag:           folderEmailModelProto.Flag,
		Deleted:        folderEmailModelProto.Deleted,
		DateOfDispatch: folderEmailModelProto.DateOfDispatch.AsTime(),
		ReplyToEmailID: folderEmailModelProto.ReplyToEmailID,
		DraftStatus:    folderEmailModelProto.DraftStatus,
		SpamStatus:     folderEmailModelProto.SpamStatus,
		SenderEmail:    folderEmailModelProto.SenderEmail,
		RecipientEmail: folderEmailModelProto.RecipientEmail,
	}
}

// ObjectEmailConvertCoreInProto converts an email model from the application core to the gRPC format.
func ObjectEmailConvertCoreInProto(folderEmailModelCore *domain.Email) *grpc.ObjectEmail {
	return &grpc.ObjectEmail{
		Id:             folderEmailModelCore.ID,
		Topic:          folderEmailModelCore.Topic,
		Text:           folderEmailModelCore.Text,
		PhotoID:        folderEmailModelCore.PhotoID,
		ReadStatus:     folderEmailModelCore.ReadStatus,
		Flag:           folderEmailModelCore.Flag,
		Deleted:        folderEmailModelCore.Deleted,
		DateOfDispatch: timestamppb.New(folderEmailModelCore.DateOfDispatch),
		ReplyToEmailID: folderEmailModelCore.ReplyToEmailID,
		DraftStatus:    folderEmailModelCore.DraftStatus,
		SpamStatus:     folderEmailModelCore.SpamStatus,
		SenderEmail:    folderEmailModelCore.SenderEmail,
		RecipientEmail: folderEmailModelCore.RecipientEmail,
	}
}

// ObjectsEmailConvertProtoInCore converts a list of email models from the gRPC format to the application core.
func ObjectsEmailConvertProtoInCore(folderEmailsModelProto *grpc.ObjectsEmail) []*domain.Email {
	emailsCore := make([]*domain.Email, 0, len(folderEmailsModelProto.Emails))
	for _, email := range folderEmailsModelProto.Emails {
		emailsCore = append(emailsCore, ObjectEmailConvertProtoInCore(email))
	}
	return emailsCore
}
