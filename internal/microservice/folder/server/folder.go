package server

import (
	"context"
	"fmt"

	"mail/internal/microservice/folder/proto"
	"mail/internal/microservice/folder/usecase"
	converters "mail/internal/microservice/models/proto_converters"
)

type FolderServer struct {
	proto.UnimplementedFolderServiceServer

	FolderUseCase *usecase.FolderUseCase
}

func NewFolderServer(folderUseCase *usecase.FolderUseCase) *FolderServer {
	return &FolderServer{FolderUseCase: folderUseCase}
}

func (es *FolderServer) CreateFolder(ctx context.Context, input *proto.Folder) (*proto.FolderWithID, error) {
	if input == nil {
		return nil, fmt.Errorf("invalid folder format: %s", input)
	}
	id, folder, err := es.FolderUseCase.CreateFolder(converters.FolderConvertProtoInCore(*input), ctx)
	if err != nil {
		return nil, fmt.Errorf("folder not found")
	}
	folderWithId := new(proto.FolderWithID)
	folderWithId.Id = id
	folderWithId.Folder = converters.FolderConvertCoreInProto(*folder)
	return folderWithId, nil
}

func (es *FolderServer) GetAllFolders(ctx context.Context, input *proto.GetAllFoldersData) (*proto.Folders, error) {
	if input == nil {
		return nil, fmt.Errorf("invalid folder format: %s", input)
	}

	if input.Id == 0 {
		return nil, fmt.Errorf("invalid profileID = %s", input.Id)
	}

	foldersCore, err := es.FolderUseCase.GetAllFolders(input.Id, input.Offset, input.Limit, ctx)
	if err != nil {
		return nil, fmt.Errorf("folder not found")
	}

	foldersProto := make([]*proto.Folder, len(foldersCore))
	for i, f := range foldersCore {
		foldersProto[i] = converters.FolderConvertCoreInProto(*f)
	}

	folderProto := new(proto.Folders)
	folderProto.Folders = foldersProto
	return folderProto, nil
}
