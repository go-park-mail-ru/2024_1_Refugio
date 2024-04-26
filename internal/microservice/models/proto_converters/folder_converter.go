package proto_converters

import (
	grpc "mail/internal/microservice/folder/proto"
	domain "mail/internal/microservice/models/domain_models"
)

func FolderConvertCoreInProto(folderModelCore domain.Folder) *grpc.Folder {
	return &grpc.Folder{
		Id:        folderModelCore.ID,
		Name:      folderModelCore.Name,
		ProfileId: folderModelCore.ProfileId,
	}
}

func FolderConvertProtoInCore(folderModelProto grpc.Folder) *domain.Folder {
	return &domain.Folder{
		ID:        folderModelProto.Id,
		Name:      folderModelProto.Name,
		ProfileId: folderModelProto.ProfileId,
	}
}

func FoldersConvertProtoInCore(folderModelProto *grpc.Folders) []*domain.Folder {
	foldersCore := make([]*domain.Folder, 0, len(folderModelProto.Folders))
	for _, folder := range folderModelProto.Folders {
		foldersCore = append(foldersCore, FolderConvertProtoInCore(*folder))
	}
	return foldersCore
}
