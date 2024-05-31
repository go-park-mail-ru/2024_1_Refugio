package server

import (
	"context"
	"fmt"
	"strconv"

	"mail/internal/microservice/folder/proto"
	"mail/internal/pkg/utils/validators"

	usecase "mail/internal/microservice/folder/interface"
	converters "mail/internal/microservice/models/proto_converters"
)

type FolderServer struct {
	proto.UnimplementedFolderServiceServer
	FolderUseCase usecase.FolderUseCase
}

func NewFolderServer(folderUseCase usecase.FolderUseCase) *FolderServer {
	return &FolderServer{FolderUseCase: folderUseCase}
}

func (es *FolderServer) CreateFolder(ctx context.Context, input *proto.Folder) (*proto.FolderWithID, error) {
	if input == nil {
		return nil, fmt.Errorf("invalid folder format: %s", input)
	}

	id, folder, err := es.FolderUseCase.CreateFolder(converters.FolderConvertProtoInCore(input), ctx)
	if err != nil {
		return nil, fmt.Errorf("failed create folder")
	}

	folderWithId := new(proto.FolderWithID)
	folderWithId.Id = id
	folderWithId.Folder = converters.FolderConvertCoreInProto(folder)
	return folderWithId, nil
}

func (es *FolderServer) GetAllFolders(ctx context.Context, input *proto.GetAllFoldersData) (*proto.Folders, error) {
	if input == nil {
		return nil, fmt.Errorf("invalid folder format: %s", input)
	}

	if input.Id <= 0 {
		return nil, fmt.Errorf("invalid profileID = %s", strconv.Itoa(int(input.Id)))
	}

	foldersCore, err := es.FolderUseCase.GetAllFolders(input.Id, input.Offset, input.Limit, ctx)
	if err != nil {
		return nil, fmt.Errorf("folder not found")
	}

	foldersProto := make([]*proto.Folder, len(foldersCore))
	for i, f := range foldersCore {
		foldersProto[i] = converters.FolderConvertCoreInProto(f)
	}

	folderProto := new(proto.Folders)
	folderProto.Folders = foldersProto
	return folderProto, nil
}

func (es *FolderServer) DeleteFolder(ctx context.Context, input *proto.DeleteFolderData) (*proto.FolderStatus, error) {
	if input == nil {
		return nil, fmt.Errorf("invalid folder format: %s", input)
	}

	if input.FolderID <= 0 || input.ProfileID <= 0 {
		return nil, fmt.Errorf("invalid folderID = %s or profileID = %s", strconv.Itoa(int(input.FolderID)), strconv.Itoa(int(input.ProfileID)))
	}

	foldersCore, err := es.FolderUseCase.DeleteFolder(input.FolderID, input.ProfileID, ctx)
	folderProto := new(proto.FolderStatus)
	folderProto.Status = foldersCore
	if err != nil || !folderProto.Status {
		return folderProto, fmt.Errorf("folder not found")
	}

	return folderProto, nil
}

func (es *FolderServer) UpdateFolder(ctx context.Context, input *proto.Folder) (*proto.FolderStatus, error) {
	if input == nil {
		return nil, fmt.Errorf("invalid folder format: %s", input)
	}

	if input.Id <= 0 || input.ProfileId <= 0 || input.Name == "" {
		return nil, fmt.Errorf("invalid folderID = %s or ProfileId = %s or Name = %s", strconv.Itoa(int(input.Id)), strconv.Itoa(int(input.ProfileId)), input.Name)
	}

	foldersCore, err := es.FolderUseCase.UpdateFolder(converters.FolderConvertProtoInCore(input), ctx)
	folderProto := new(proto.FolderStatus)
	folderProto.Status = foldersCore
	if err != nil || !folderProto.Status {
		return folderProto, fmt.Errorf("folder not found")
	}

	return folderProto, nil
}

func (es *FolderServer) AddEmailInFolder(ctx context.Context, input *proto.FolderEmail) (*proto.FolderEmailStatus, error) {
	if input == nil {
		return nil, fmt.Errorf("invalid folder format: %s", input)
	}

	if input.EmailID <= 0 || input.FolderID <= 0 {
		return nil, fmt.Errorf("invalid EmailID = %s or FolderID = %s", strconv.Itoa(int(input.EmailID)), strconv.Itoa(int(input.FolderID)))
	}

	status, err := es.FolderUseCase.AddEmailInFolder(input.FolderID, input.EmailID, ctx)
	if err != nil || !status {
		return nil, fmt.Errorf("folder or email not found")
	}
	folderStatus := new(proto.FolderEmailStatus)
	folderStatus.Status = status
	return folderStatus, nil
}

func (es *FolderServer) DeleteEmailInFolder(ctx context.Context, input *proto.FolderEmail) (*proto.FolderEmailStatus, error) {
	if input == nil {
		return nil, fmt.Errorf("invalid folder format: %s", input)
	}

	if input.EmailID <= 0 || input.FolderID <= 0 {
		return nil, fmt.Errorf("invalid EmailID = %s or FolderID = %s", strconv.Itoa(int(input.EmailID)), strconv.Itoa(int(input.FolderID)))
	}

	status, err := es.FolderUseCase.DeleteEmailInFolder(input.FolderID, input.EmailID, ctx)
	folderStatus := new(proto.FolderEmailStatus)
	folderStatus.Status = status

	if err != nil || !status {
		return folderStatus, fmt.Errorf("folder or email not found")
	}

	return folderStatus, nil
}

func (es *FolderServer) CheckFolderProfile(ctx context.Context, input *proto.FolderProfile) (*proto.FolderEmailStatus, error) {
	if input == nil {
		return nil, fmt.Errorf("invalid folder format: %s", input)
	}

	if input.ProfileID <= 0 || input.FolderID <= 0 {
		return nil, fmt.Errorf("invalid ProfileID = %s or FolderID = %s", strconv.Itoa(int(input.ProfileID)), strconv.Itoa(int(input.FolderID)))
	}

	status, err := es.FolderUseCase.CheckFolderProfile(input.FolderID, input.ProfileID, ctx)
	if err != nil || !status {
		return nil, fmt.Errorf("folder and profile not found")
	}

	folderStatus := new(proto.FolderEmailStatus)
	folderStatus.Status = status
	return folderStatus, nil
}

func (es *FolderServer) CheckEmailProfile(ctx context.Context, input *proto.EmailProfile) (*proto.FolderEmailStatus, error) {
	if input == nil {
		return nil, fmt.Errorf("invalid folder format: %s", input)
	}

	if input.ProfileID <= 0 || input.EmailID <= 0 {
		return nil, fmt.Errorf("invalid ProfileID = %s or EmailID = %s", strconv.Itoa(int(input.ProfileID)), strconv.Itoa(int(input.EmailID)))
	}

	status, err := es.FolderUseCase.CheckEmailProfile(input.EmailID, input.ProfileID, ctx)
	if err != nil || !status {
		return nil, fmt.Errorf("email and profile not found")
	}

	folderStatus := new(proto.FolderEmailStatus)
	folderStatus.Status = status
	return folderStatus, nil
}

func (es *FolderServer) GetAllEmailsInFolder(ctx context.Context, input *proto.GetAllEmailsInFolderData) (*proto.ObjectsEmail, error) {
	if input == nil {
		return nil, fmt.Errorf("invalid folder format: %s", input)
	}

	if input.FolderID <= 0 || input.ProfileID <= 0 {
		return nil, fmt.Errorf("invalid folderID=%s or profileID=%s or limit=%s or offset=%s", strconv.Itoa(int(input.FolderID)), strconv.Itoa(int(input.ProfileID)), strconv.Itoa(int(input.Limit)), strconv.Itoa(int(input.Offset)))
	}

	if validators.IsEmpty(input.Login) {
		return nil, fmt.Errorf("invalid login")
	}

	emailsCore, err := es.FolderUseCase.GetAllEmailsInFolder(input.FolderID, input.ProfileID, input.Limit, input.Offset, input.Login, ctx)
	if err != nil {
		return nil, fmt.Errorf("emails not found")
	}

	emailsProto := make([]*proto.ObjectEmail, len(emailsCore))
	for i, e := range emailsCore {
		emailsProto[i] = converters.ObjectEmailConvertCoreInProto(e)
	}

	emailProto := new(proto.ObjectsEmail)
	emailProto.Emails = emailsProto
	return emailProto, nil
}

func (es *FolderServer) GetAllNameFolders(ctx context.Context, input *proto.GetAllNameFoldersRequest) (*proto.Folders, error) {
	if input == nil {
		return nil, fmt.Errorf("invalid folder format: %s", input)
	}

	if input.EmailId <= 0 {
		return nil, fmt.Errorf("invalid emailID=%s", strconv.Itoa(int(input.EmailId)))
	}

	foldersCore, err := es.FolderUseCase.GetAllFolderName(input.EmailId, ctx)
	if err != nil {
		return nil, fmt.Errorf("emails not found")
	}

	foldersProto := make([]*proto.Folder, len(foldersCore))
	for i, f := range foldersCore {
		foldersProto[i] = converters.FolderConvertCoreInProto(f)
	}

	folderProto := new(proto.Folders)
	folderProto.Folders = foldersProto
	return folderProto, nil
}
