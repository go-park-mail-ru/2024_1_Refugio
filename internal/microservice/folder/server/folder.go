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
