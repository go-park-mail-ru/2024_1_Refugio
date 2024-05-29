package server

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"mail/internal/microservice/folder/mock"
	"mail/internal/microservice/folder/proto"
	"mail/internal/microservice/models/domain_models"
	"mail/internal/pkg/logger"
	"mail/internal/pkg/utils/constants"

	converters "mail/internal/microservice/models/proto_converters"
)

func GetCTX() context.Context {
	ctx := context.WithValue(context.Background(), constants.LoggerKey, logger.InitializationBdLog(nil))
	ctx2 := context.WithValue(ctx, constants.RequestIDKey, []string{"testID"})

	return ctx2
}

func TestNewFolderServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFolderUseCase := mock.NewMockFolderUseCase(ctrl)

	server := NewFolderServer(mockFolderUseCase)

	expectedServer := FolderServer{
		FolderUseCase: mockFolderUseCase,
	}

	assert.Equal(t, expectedServer, *server)
}

func TestCreateFolder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFolderUseCase := mock.NewMockFolderUseCase(ctrl)

	server := NewFolderServer(mockFolderUseCase)

	ctx := GetCTX()

	domainFolder := &domain_models.Folder{ID: 1, ProfileId: 1, Name: "Test"}
	folderProto := converters.FolderConvertCoreInProto(domainFolder)
	id := uint32(1)
	expectedFolder := &proto.FolderWithID{Folder: folderProto, Id: id}

	t.Run("CreateFolderSuccessfully", func(t *testing.T) {
		mockFolderUseCase.EXPECT().CreateFolder(domainFolder, ctx).Return(id, domainFolder, nil)

		folderWithID, err := server.CreateFolder(ctx, folderProto)

		assert.NoError(t, err)
		assert.Equal(t, expectedFolder, folderWithID)
	})

	t.Run("CreateFolder invalid folder format", func(t *testing.T) {
		_, err := server.CreateFolder(ctx, nil)

		assert.Error(t, err)
	})

	t.Run("CreateFolder failed create folder", func(t *testing.T) {
		mockFolderUseCase.EXPECT().CreateFolder(domainFolder, ctx).Return(id, domainFolder, fmt.Errorf("failed create folder"))

		_, err := server.CreateFolder(ctx, folderProto)

		assert.Error(t, err)
		assert.Equal(t, fmt.Errorf("failed create folder"), err)
	})
}

func TestGetAllFolders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFolderUseCase := mock.NewMockFolderUseCase(ctrl)

	server := NewFolderServer(mockFolderUseCase)

	ctx := GetCTX()

	domainFolders := []*domain_models.Folder{
		{ID: 1, ProfileId: 1, Name: "Test 2"},
		{ID: 2, ProfileId: 2, Name: "Test 2"},
	}
	foldersProto := make([]*proto.Folder, len(domainFolders))
	for i, e := range domainFolders {
		foldersProto[i] = converters.FolderConvertCoreInProto(e)
	}
	folderProto := new(proto.Folders)
	folderProto.Folders = foldersProto

	foldersData := &proto.GetAllFoldersData{Id: 1, Limit: 0, Offset: 0}
	t.Run("GetAllFoldersSuccessfully", func(t *testing.T) {
		mockFolderUseCase.EXPECT().GetAllFolders(foldersData.Id, foldersData.Offset, foldersData.Limit, ctx).Return(domainFolders, nil)

		folders, err := server.GetAllFolders(ctx, foldersData)

		assert.NoError(t, err)
		assert.Equal(t, folderProto, folders)
	})

	t.Run("CreateFolder invalid folder format", func(t *testing.T) {
		_, err := server.GetAllFolders(ctx, nil)

		assert.Error(t, err)
	})

	t.Run("CreateFolder failed create folder", func(t *testing.T) {
		mockFolderUseCase.EXPECT().GetAllFolders(foldersData.Id, foldersData.Offset, foldersData.Limit, ctx).Return(domainFolders, fmt.Errorf("folder not found"))

		_, err := server.GetAllFolders(ctx, foldersData)

		assert.Error(t, err)
		assert.Equal(t, fmt.Errorf("folder not found"), err)
	})
}

func TestDeleteFolder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFolderUseCase := mock.NewMockFolderUseCase(ctrl)

	server := NewFolderServer(mockFolderUseCase)

	ctx := GetCTX()

	domainFolders := []*domain_models.Folder{
		{ID: 1, ProfileId: 1, Name: "Test 2"},
		{ID: 2, ProfileId: 2, Name: "Test 2"},
	}
	foldersProto := make([]*proto.Folder, len(domainFolders))
	for i, e := range domainFolders {
		foldersProto[i] = converters.FolderConvertCoreInProto(e)
	}
	folderProto := new(proto.Folders)
	folderProto.Folders = foldersProto

	foldersData := &proto.DeleteFolderData{FolderID: 1, ProfileID: 1}
	t.Run("DeleteFolderSuccessfully", func(t *testing.T) {
		mockFolderUseCase.EXPECT().DeleteFolder(foldersData.FolderID, foldersData.ProfileID, ctx).Return(true, nil)

		status, err := server.DeleteFolder(ctx, foldersData)

		assert.NoError(t, err)
		assert.True(t, status.Status)
	})

	t.Run("CreateFolder invalid folder format", func(t *testing.T) {
		_, err := server.DeleteFolder(ctx, nil)

		assert.Error(t, err)
	})

	t.Run("CreateFolder failed create folder", func(t *testing.T) {
		mockFolderUseCase.EXPECT().DeleteFolder(foldersData.FolderID, foldersData.ProfileID, ctx).Return(false, fmt.Errorf("folder not found"))

		status, err := server.DeleteFolder(ctx, foldersData)

		assert.Error(t, err)
		assert.Equal(t, fmt.Errorf("folder not found"), err)
		assert.False(t, status.Status)
	})
}

func TestUpdateFolder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFolderUseCase := mock.NewMockFolderUseCase(ctrl)

	server := NewFolderServer(mockFolderUseCase)

	ctx := GetCTX()

	domainFolders := []*domain_models.Folder{
		{ID: 1, ProfileId: 1, Name: "Test 2"},
		{ID: 2, ProfileId: 2, Name: "Test 2"},
	}
	foldersProto := make([]*proto.Folder, len(domainFolders))
	for i, e := range domainFolders {
		foldersProto[i] = converters.FolderConvertCoreInProto(e)
	}
	folderProto := new(proto.Folders)
	folderProto.Folders = foldersProto

	t.Run("UpdateFolderSuccessfully", func(t *testing.T) {
		mockFolderUseCase.EXPECT().UpdateFolder(domainFolders[0], ctx).Return(true, nil)

		status, err := server.UpdateFolder(ctx, foldersProto[0])

		assert.NoError(t, err)
		assert.True(t, status.Status)
	})

	t.Run("UpdateFolder invalid folder format", func(t *testing.T) {
		_, err := server.UpdateFolder(ctx, nil)

		assert.Error(t, err)
	})

	t.Run("UpdateFolder failed create folder", func(t *testing.T) {
		mockFolderUseCase.EXPECT().UpdateFolder(domainFolders[0], ctx).Return(false, fmt.Errorf("folder not found"))

		status, err := server.UpdateFolder(ctx, foldersProto[0])

		assert.Error(t, err)
		assert.Equal(t, fmt.Errorf("folder not found"), err)
		assert.False(t, status.Status)
	})

}

func TestAddEmailInFolder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFolderUseCase := mock.NewMockFolderUseCase(ctrl)

	server := NewFolderServer(mockFolderUseCase)

	ctx := GetCTX()

	data := &proto.FolderEmail{FolderID: 1, EmailID: 1}

	t.Run("AddEmailInFolderSuccessfully", func(t *testing.T) {
		mockFolderUseCase.EXPECT().AddEmailInFolder(data.FolderID, data.EmailID, ctx).Return(true, nil)

		status, err := server.AddEmailInFolder(ctx, data)

		assert.NoError(t, err)
		assert.True(t, status.Status)
	})

	t.Run("UpdateFolder invalid folder format", func(t *testing.T) {
		_, err := server.AddEmailInFolder(ctx, nil)

		assert.Error(t, err)
	})

	t.Run("UpdateFolder failed create folder", func(t *testing.T) {
		mockFolderUseCase.EXPECT().AddEmailInFolder(data.FolderID, data.EmailID, ctx).Return(false, fmt.Errorf("folder or email not found"))

		status, err := server.AddEmailInFolder(ctx, data)

		assert.Error(t, err)
		assert.Equal(t, fmt.Errorf("folder or email not found"), err)
		assert.Nil(t, status)
	})

}

func TestDeleteEmailInFolder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFolderUseCase := mock.NewMockFolderUseCase(ctrl)

	server := NewFolderServer(mockFolderUseCase)

	ctx := GetCTX()

	data := &proto.FolderEmail{FolderID: 1, EmailID: 1}

	t.Run("DeleteEmailInFolderSuccessfully", func(t *testing.T) {
		mockFolderUseCase.EXPECT().DeleteEmailInFolder(data.FolderID, data.EmailID, ctx).Return(true, nil)

		status, err := server.DeleteEmailInFolder(ctx, data)

		assert.NoError(t, err)
		assert.True(t, status.Status)
	})

	t.Run("DeleteEmailInFolder invalid folder format", func(t *testing.T) {
		_, err := server.DeleteEmailInFolder(ctx, nil)

		assert.Error(t, err)
	})

	t.Run("DeleteEmailInFolder failed create folder", func(t *testing.T) {
		mockFolderUseCase.EXPECT().DeleteEmailInFolder(data.FolderID, data.EmailID, ctx).Return(false, fmt.Errorf("folder or email not found"))

		status, err := server.DeleteEmailInFolder(ctx, data)

		assert.Error(t, err)
		assert.Equal(t, fmt.Errorf("folder or email not found"), err)
		assert.False(t, status.Status)
	})
}

func TestCheckFolderProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFolderUseCase := mock.NewMockFolderUseCase(ctrl)

	server := NewFolderServer(mockFolderUseCase)

	ctx := GetCTX()

	data := &proto.FolderProfile{FolderID: 1, ProfileID: 1}

	t.Run("CheckFolderProfileSuccessfully", func(t *testing.T) {
		mockFolderUseCase.EXPECT().CheckFolderProfile(data.FolderID, data.ProfileID, ctx).Return(true, nil)

		status, err := server.CheckFolderProfile(ctx, data)

		assert.NoError(t, err)
		assert.True(t, status.Status)
	})

	t.Run("CheckFolderProfile invalid folder format", func(t *testing.T) {
		_, err := server.CheckFolderProfile(ctx, nil)

		assert.Error(t, err)
	})

	t.Run("CheckFolderProfile failed create folder", func(t *testing.T) {
		mockFolderUseCase.EXPECT().CheckFolderProfile(data.FolderID, data.ProfileID, ctx).Return(false, fmt.Errorf("folder and profile not found"))

		status, err := server.CheckFolderProfile(ctx, data)

		assert.Error(t, err)
		assert.Equal(t, fmt.Errorf("folder and profile not found"), err)
		assert.Nil(t, status)
	})
}

func TestCheckEmailProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFolderUseCase := mock.NewMockFolderUseCase(ctrl)

	server := NewFolderServer(mockFolderUseCase)

	ctx := GetCTX()

	data := &proto.EmailProfile{EmailID: 1, ProfileID: 1}

	t.Run("CheckEmailProfileSuccessfully", func(t *testing.T) {
		mockFolderUseCase.EXPECT().CheckEmailProfile(data.EmailID, data.ProfileID, ctx).Return(true, nil)

		status, err := server.CheckEmailProfile(ctx, data)

		assert.NoError(t, err)
		assert.True(t, status.Status)
	})

	t.Run("CheckEmailProfile invalid folder format", func(t *testing.T) {
		_, err := server.CheckEmailProfile(ctx, nil)

		assert.Error(t, err)
	})

	t.Run("CheckEmailProfile failed create folder", func(t *testing.T) {
		mockFolderUseCase.EXPECT().CheckEmailProfile(data.EmailID, data.ProfileID, ctx).Return(false, fmt.Errorf("email and profile not found"))

		status, err := server.CheckEmailProfile(ctx, data)

		assert.Error(t, err)
		assert.Equal(t, fmt.Errorf("email and profile not found"), err)
		assert.Nil(t, status)
	})
}

func TestGetAllEmailsInFolder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFolderUseCase := mock.NewMockFolderUseCase(ctrl)

	server := NewFolderServer(mockFolderUseCase)

	ctx := GetCTX()

	data := &proto.GetAllEmailsInFolderData{FolderID: 1, ProfileID: 1, Limit: 0, Offset: 0, Login: "loginUser"}
	emails := []*domain_models.Email{
		{ID: 1, Topic: "Topic 1", Text: "Text 1"},
		{ID: 2, Topic: "Topic 2", Text: "Text 2"},
		{ID: 3, Topic: "Topic 3", Text: "Text 3"},
	}

	emailObArr := []*proto.ObjectEmail{
		{Id: 1, Topic: "Topic 1", Text: "Text 1"},
		{Id: 2, Topic: "Topic 2", Text: "Text 2"},
		{Id: 3, Topic: "Topic 3", Text: "Text 3"},
	}

	expectedObjectsEmail := &proto.ObjectsEmail{Emails: emailObArr}

	t.Run("GetAllEmailsInFolderSuccessfully", func(t *testing.T) {
		mockFolderUseCase.EXPECT().GetAllEmailsInFolder(data.FolderID, data.ProfileID, data.Limit, data.Offset, data.Login, ctx).Return(emails, nil)

		objectEmail, err := server.GetAllEmailsInFolder(ctx, data)

		assert.NoError(t, err)
		for i := range objectEmail.Emails {
			objectEmail.Emails[i].DateOfDispatch = nil
			assert.Equal(t, expectedObjectsEmail.Emails[i], objectEmail.Emails[i])
		}
	})

	t.Run("GetAllEmailsInFolder invalid folder format", func(t *testing.T) {
		_, err := server.GetAllEmailsInFolder(ctx, nil)

		assert.Error(t, err)
	})

	t.Run("CheckEmailProfile failed create folder", func(t *testing.T) {
		mockFolderUseCase.EXPECT().GetAllEmailsInFolder(data.FolderID, data.ProfileID, data.Limit, data.Offset, data.Login, ctx).Return(emails, fmt.Errorf("emails not found"))

		objectEmail, err := server.GetAllEmailsInFolder(ctx, data)

		assert.Error(t, err)
		assert.Equal(t, fmt.Errorf("emails not found"), err)
		assert.Nil(t, objectEmail)
	})
}

func TestGetAllFolder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFolderUseCase := mock.NewMockFolderUseCase(ctrl)

	server := NewFolderServer(mockFolderUseCase)

	ctx := GetCTX()

	data := &proto.GetAllFoldersData{Id: 1, Offset: 0, Limit: 0}
	folders := []*domain_models.Folder{
		{ID: 1, Name: "Folder 1"},
		{ID: 2, Name: "Folder 2"},
		{ID: 3, Name: "Folder 3"},
	}

	folderArr := []*proto.Folder{
		{Id: 1, Name: "Folder 1"},
		{Id: 2, Name: "Folder 2"},
		{Id: 3, Name: "Folder 3"},
	}

	expectedFolders := &proto.Folders{Folders: folderArr}

	t.Run("GetAllFolders_Successfully", func(t *testing.T) {
		mockFolderUseCase.EXPECT().GetAllFolders(data.Id, data.Offset, data.Limit, ctx).Return(folders, nil)

		foldersProto, err := server.GetAllFolders(ctx, data)

		assert.NoError(t, err)
		assert.Equal(t, expectedFolders, foldersProto)
	})

	t.Run("GetAllFolders_InvalidFolderFormat", func(t *testing.T) {
		_, err := server.GetAllFolders(ctx, nil)

		assert.Error(t, err)
	})

	t.Run("GetAllFolders_NotFound", func(t *testing.T) {
		mockFolderUseCase.EXPECT().GetAllFolders(data.Id, data.Offset, data.Limit, ctx).Return(folders, fmt.Errorf("folders not found"))

		foldersProto, err := server.GetAllFolders(ctx, data)

		assert.Error(t, err)
		assert.Nil(t, foldersProto)
	})
}
