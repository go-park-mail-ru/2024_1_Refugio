package usecase

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"mail/internal/pkg/logger"
	"mail/internal/pkg/utils/constants"

	mock_repository "mail/internal/microservice/folder/mock"
	domain "mail/internal/microservice/models/domain_models"
)

func GetCTX() context.Context {
	f, err := os.OpenFile("log_test.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create logfile" + "log.txt")
	}
	defer f.Close()

	ctx := context.WithValue(context.Background(), constants.LoggerKey, logger.InitializationBdLog(f))
	ctx2 := context.WithValue(ctx, constants.RequestIDKey, []string{"testID"})

	return ctx2
}

func TestNewFolderUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockFolderRepository(ctrl)

	ExpextedFolderUseCase := FolderUseCase{
		repo: mockRepo,
	}

	EmailUseCase := NewFolderUseCase(mockRepo)

	assert.Equal(t, ExpextedFolderUseCase, *EmailUseCase)
}

func TestCreateFolder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockFolderRepository(ctrl)
	useCase := NewFolderUseCase(mockRepo)

	expectedFolder := &domain.Folder{ID: 1, Name: "Test Folder", ProfileId: 1}
	ctx := GetCTX()
	id := uint32(1)

	mockRepo.EXPECT().Create(expectedFolder, ctx).Return(id, expectedFolder, nil)

	ID, folder, err := useCase.CreateFolder(expectedFolder, ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedFolder, folder)
	assert.Equal(t, id, ID)
}

func TestGetAllFolders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockFolderRepository(ctrl)
	useCase := NewFolderUseCase(mockRepo)

	expectedFolders := []*domain.Folder{
		{ID: 1, Name: "Test Folder", ProfileId: 1},
		{ID: 2, Name: "Test Folder", ProfileId: 2},
	}
	ctx := GetCTX()
	zero := int64(0)
	profile_id := uint32(1)

	mockRepo.EXPECT().GetAll(profile_id, zero, zero, ctx).Return(expectedFolders, nil)

	folders, err := useCase.GetAllFolders(profile_id, zero, zero, ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedFolders, folders)
}

func TestDeleteFolder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockFolderRepository(ctrl)
	useCase := NewFolderUseCase(mockRepo)

	ctx := GetCTX()
	folder_id := uint32(1)
	profile_id := uint32(1)

	mockRepo.EXPECT().Delete(folder_id, profile_id, ctx).Return(true, nil)

	status, err := useCase.DeleteFolder(folder_id, profile_id, ctx)

	assert.NoError(t, err)
	assert.True(t, status)
}

func TestUpdateFolder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockFolderRepository(ctrl)
	useCase := NewFolderUseCase(mockRepo)

	ctx := GetCTX()
	newUpFolder := &domain.Folder{ID: 1, Name: "Test Folder", ProfileId: 1}

	mockRepo.EXPECT().Update(newUpFolder, ctx).Return(true, nil)

	status, err := useCase.UpdateFolder(newUpFolder, ctx)

	assert.NoError(t, err)
	assert.True(t, status)
}

func TestAddEmailInFolder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockFolderRepository(ctrl)
	useCase := NewFolderUseCase(mockRepo)

	ctx := GetCTX()
	folder_id := uint32(1)
	email_id := uint32(1)

	mockRepo.EXPECT().AddEmailFolder(folder_id, email_id, ctx).Return(true, nil)

	status, err := useCase.AddEmailInFolder(folder_id, email_id, ctx)

	assert.NoError(t, err)
	assert.True(t, status)
}

func TestDeleteEmailInFolder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockFolderRepository(ctrl)
	useCase := NewFolderUseCase(mockRepo)

	ctx := GetCTX()
	folder_id := uint32(1)
	email_id := uint32(1)

	mockRepo.EXPECT().DeleteEmailFolder(folder_id, email_id, ctx).Return(true, nil)

	status, err := useCase.DeleteEmailInFolder(folder_id, email_id, ctx)

	assert.NoError(t, err)
	assert.True(t, status)
}

func TestCheckFolderProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockFolderRepository(ctrl)
	useCase := NewFolderUseCase(mockRepo)

	ctx := GetCTX()
	folder_id := uint32(1)
	profile_id := uint32(1)

	mockRepo.EXPECT().CheckFolder(folder_id, profile_id, ctx).Return(true, nil)

	status, err := useCase.CheckFolderProfile(folder_id, profile_id, ctx)

	assert.NoError(t, err)
	assert.True(t, status)
}

func TestCheckEmailProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockFolderRepository(ctrl)
	useCase := NewFolderUseCase(mockRepo)

	ctx := GetCTX()
	email_id := uint32(1)
	profile_id := uint32(1)

	mockRepo.EXPECT().CheckEmail(email_id, profile_id, ctx).Return(true, nil)

	status, err := useCase.CheckEmailProfile(email_id, profile_id, ctx)

	assert.NoError(t, err)
	assert.True(t, status)
}

func TestGetAllEmailsInFolder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockFolderRepository(ctrl)
	useCase := NewFolderUseCase(mockRepo)

	ctx := GetCTX()
	folder_id := uint32(1)
	profile_id := uint32(1)
	zero := uint32(0)
	login := "loginUser"

	expectedFolders := []*domain.Email{
		{ID: 1, Topic: "Test topic 1", Text: "Test text 1"},
		{ID: 2, Topic: "Test topic 2", Text: "Test text 2"},
	}

	mockRepo.EXPECT().GetAllEmails(folder_id, profile_id, zero, zero, ctx).Return(expectedFolders, nil)

	folders, err := useCase.GetAllEmailsInFolder(folder_id, profile_id, zero, zero, login, ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedFolders, folders)
}
