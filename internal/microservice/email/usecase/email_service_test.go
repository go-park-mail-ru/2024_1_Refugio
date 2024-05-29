package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"mail/internal/pkg/logger"
	"mail/internal/pkg/utils/constants"

	mockRepository "mail/internal/microservice/email/mock"
	domain "mail/internal/microservice/models/domain_models"
)

func GetCTX() context.Context {
	ctx := context.WithValue(context.Background(), constants.LoggerKey, logger.InitializationBdLog(nil))
	ctx2 := context.WithValue(ctx, constants.RequestIDKey, []string{"testID"})

	return ctx2
}

func TestNewEmailUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)

	ExpectedEmailUseCase := EmailUseCase{
		repo: mockRepo,
	}

	EmailUseCase := NewEmailUseCase(mockRepo)

	assert.Equal(t, ExpectedEmailUseCase, *EmailUseCase)
}

func TestGetAllEmailsIncoming_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	login := "test@mailhub.su"
	expectedEmails := []*domain.Email{
		{Topic: "Topic 1", Text: "Text 1"},
		{Topic: "Topic 2", Text: "Text 2"},
	}
	ctx := GetCTX()

	zero := int64(0)

	mockRepo.EXPECT().GetAllIncoming(login, zero, zero, ctx).Return(expectedEmails, nil)

	emails, err := useCase.GetAllEmailsIncoming(login, zero, zero, ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedEmails, emails)
}

func TestAllEmailsIncoming_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	login := "test@mailhub.su"
	ctx := GetCTX()
	zero := int64(0)
	mockRepo.EXPECT().GetAllIncoming(login, zero, zero, ctx).Return(nil, errors.New("repository error"))

	emails, err := useCase.GetAllEmailsIncoming(login, zero, zero, ctx)

	assert.Error(t, err)
	assert.Nil(t, emails)
}

func TestGetAllEmailsSent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	login := "test@mailhub.su"
	expectedEmails := []*domain.Email{
		{Topic: "Topic 1", Text: "Text 1"},
		{Topic: "Topic 2", Text: "Text 2"},
	}
	ctx := GetCTX()
	zero := int64(0)
	mockRepo.EXPECT().GetAllSent(login, zero, zero, ctx).Return(expectedEmails, nil)

	emails, err := useCase.GetAllEmailsSent(login, zero, zero, ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedEmails, emails)
}

func TestGetAllEmailsSent_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	login := "test@mailhub.su"
	ctx := GetCTX()
	zero := int64(0)

	mockRepo.EXPECT().GetAllSent(login, zero, zero, ctx).Return(nil, errors.New("repository error"))

	emails, err := useCase.GetAllEmailsSent(login, zero, zero, ctx)

	assert.Error(t, err)
	assert.Nil(t, emails)
}

func TestGetAllEmailsDraft(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	login := "test@mailhub.su"
	expectedEmails := []*domain.Email{
		{Topic: "Topic 1", Text: "Text 1"},
		{Topic: "Topic 2", Text: "Text 2"},
	}
	ctx := GetCTX()
	zero := int64(0)

	mockRepo.EXPECT().GetAllDraft(login, zero, zero, ctx).Return(expectedEmails, nil)

	emails, err := useCase.GetAllDraftEmails(login, zero, zero, ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedEmails, emails)
}

func TestGetAllEmailsSpam(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	login := "test@mailhub.su"
	expectedEmails := []*domain.Email{
		{Topic: "Topic 1", Text: "Text 1"},
		{Topic: "Topic 2", Text: "Text 2"},
	}
	ctx := GetCTX()
	zero := int64(0)

	mockRepo.EXPECT().GetAllSpam(login, zero, zero, ctx).Return(expectedEmails, nil)

	emails, err := useCase.GetAllSpamEmails(login, zero, zero, ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedEmails, emails)
}

func TestGetEmailByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	login := "test@mailhub.su"
	ctx := GetCTX()

	expectedEmail := &domain.Email{ID: 1, Topic: "Topic 1", Text: "Text 1"}
	mockRepo.EXPECT().GetByID(uint64(1), login, ctx).Return(expectedEmail, nil)

	email, err := useCase.GetEmailByID(1, login, ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedEmail, email)
}

func TestGetEmailByID_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	login := "test@mailhub.su"
	ctx := GetCTX()

	mockRepo.EXPECT().GetByID(uint64(1), login, ctx).Return(nil, errors.New("repository error"))

	email, err := useCase.GetEmailByID(1, login, ctx)

	assert.Error(t, err)
	assert.Nil(t, email)
}

func TestCreateEmail_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	ctx := GetCTX()

	newEmail := &domain.Email{Topic: "Topic 1", Text: "Text 1"}
	mockRepo.EXPECT().Add(gomock.Any(), ctx).Return(uint64(1), newEmail, nil)

	id, emailRes, err := useCase.CreateEmail(newEmail, ctx)

	assert.Equal(t, uint64(1), id)
	assert.NoError(t, err)
	assert.Equal(t, newEmail, emailRes)
}

func TestCreateEmail_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)
	newEmail := &domain.Email{Topic: "Topic 1", Text: "Text 1"}

	ctx := GetCTX()

	mockRepo.EXPECT().Add(gomock.Any(), ctx).Return(uint64(1), newEmail, errors.New("repository error"))

	id, emailRes, err := useCase.CreateEmail(newEmail, ctx)

	assert.Equal(t, uint64(1), id)
	assert.Error(t, err)
	assert.Equal(t, newEmail, emailRes)
}

func TestCreateProfileEmail_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	emailId := uint64(1)
	sender := "test_sender@mailhub.su"
	recipient := "test_recipient@mailhub.su"
	ctx := GetCTX()

	mockRepo.EXPECT().AddProfileEmail(emailId, sender, recipient, ctx).Return(nil)

	err := useCase.CreateProfileEmail(emailId, sender, recipient, ctx)

	assert.NoError(t, err)
}

func TestCreateProfileEmail_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	emailId := uint64(1)
	sender := "test_sender@mailhub.su"
	recipient := "test_recipient@mailhub.su"
	ctx := GetCTX()

	mockRepo.EXPECT().AddProfileEmail(emailId, sender, recipient, ctx).Return(errors.New("repository error"))

	err := useCase.CreateProfileEmail(emailId, sender, recipient, ctx)

	assert.Error(t, err)
}

func TestCheckRecipientEmail_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	recipient := "test_recipient@mailhub.su"
	ctx := GetCTX()

	mockRepo.EXPECT().FindEmail(recipient, ctx).Return(nil)

	err := useCase.CheckRecipientEmail(recipient, ctx)

	assert.NoError(t, err)
}

func TestCheckRecipientEmail_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	recipient := "test_recipient@mailhub.su"
	ctx := GetCTX()

	mockRepo.EXPECT().FindEmail(recipient, ctx).Return(errors.New("repository error"))

	err := useCase.CheckRecipientEmail(recipient, ctx)

	assert.Error(t, err)
}

func TestUpdateEmail_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	newEmail := &domain.Email{ID: 1, Topic: "Topic 1", Text: "Text 1"}
	ctx := GetCTX()

	mockRepo.EXPECT().Update(gomock.Any(), ctx).Return(true, nil)

	emailRes, err := useCase.UpdateEmail(newEmail, ctx)

	assert.NoError(t, err)
	assert.Equal(t, true, emailRes)
}

func TestDeleteEmail_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	login := "test@mailhub.su"
	ctx := GetCTX()

	mockRepo.EXPECT().Delete(gomock.Any(), login, ctx).Return(true, nil)

	emailRes, err := useCase.DeleteEmail(uint64(1), login, ctx)

	assert.NoError(t, err)
	assert.Equal(t, true, emailRes)
}

func TestAddAttachment_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	fileID := "test_file_id"
	fileType := "pdf"
	fileName := "PDF"
	fileSize := "10101010"
	emailID := uint64(123)
	ctx := context.Background()

	mockRepo.EXPECT().AddFile(fileID, fileType, fileName, fileSize, ctx).Return(uint64(456), nil)
	mockRepo.EXPECT().AddAttachment(emailID, uint64(456), ctx).Return(nil)

	result, err := useCase.AddAttachment(fileID, fileType, fileName, fileSize, emailID, ctx)

	assert.NoError(t, err)
	assert.Equal(t, uint64(456), result)
}

func TestAddAttachment_EmptyFileID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	fileID := ""
	fileType := "pdf"
	fileName := "PDF"
	fileSize := "10101010"
	emailID := uint64(123)
	ctx := context.Background()

	result, err := useCase.AddAttachment(fileID, fileType, fileName, fileSize, emailID, ctx)

	assert.Error(t, err)
	assert.Zero(t, result)
}

func TestAddAttachment_EmptyFileType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	fileID := "test_file_id"
	fileType := ""
	fileName := "PDF"
	fileSize := "10101010"
	emailID := uint64(123)
	ctx := context.Background()

	result, err := useCase.AddAttachment(fileID, fileType, fileName, fileSize, emailID, ctx)

	assert.Error(t, err)
	assert.Zero(t, result)
}

func TestAddAttachment_InvalidEmailID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	fileID := "test_file_id"
	fileType := "pdf"
	fileName := "PDF"
	fileSize := "10101010"
	emailID := uint64(0)
	ctx := context.Background()

	result, err := useCase.AddAttachment(fileID, fileType, fileName, fileSize, emailID, ctx)

	assert.Error(t, err)
	assert.Zero(t, result)
}

func TestAddAttachment_ErrorAddingFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	fileID := "test_file_id"
	fileType := "pdf"
	fileName := "PDF"
	fileSize := "10101010"
	emailID := uint64(123)
	ctx := context.Background()

	mockRepo.EXPECT().AddFile(fileID, fileType, fileName, fileSize, ctx).Return(uint64(0), errors.New("file adding error"))

	result, err := useCase.AddAttachment(fileID, fileType, fileName, fileSize, emailID, ctx)

	assert.Error(t, err)
	assert.Zero(t, result)
}

func TestAddAttachment_ErrorAddingAttachment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	fileID := "test_file_id"
	fileType := "pdf"
	fileName := "PDF"
	fileSize := "10101010"
	emailID := uint64(123)
	ctx := context.Background()

	mockRepo.EXPECT().AddFile(fileID, fileType, fileName, fileSize, ctx).Return(uint64(456), nil)
	mockRepo.EXPECT().AddAttachment(emailID, uint64(456), ctx).Return(errors.New("attachment adding error"))

	result, err := useCase.AddAttachment(fileID, fileType, fileName, fileSize, emailID, ctx)

	assert.Error(t, err)
	assert.Zero(t, result)
}

func TestGetFileByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	fileID := uint64(123)
	expectedFile := &domain.File{ID: fileID, FileId: "test_file"}

	ctx := context.Background()

	mockRepo.EXPECT().GetFileByID(fileID, ctx).Return(expectedFile, nil)

	file, err := useCase.GetFileByID(fileID, ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedFile, file)
}

func TestGetFileByID_InvalidFileID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	fileID := uint64(0)
	ctx := context.Background()

	file, err := useCase.GetFileByID(fileID, ctx)

	assert.Error(t, err)
	assert.Nil(t, file)
}

func TestGetFileByID_ErrorGettingFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	fileID := uint64(123)
	ctx := context.Background()

	mockRepo.EXPECT().GetFileByID(fileID, ctx).Return(nil, errors.New("file retrieval error"))

	file, err := useCase.GetFileByID(fileID, ctx)

	assert.Error(t, err)
	assert.Nil(t, file)
}

func TestGetFilesByEmailID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	emailID := uint64(123)
	expectedFiles := []*domain.File{
		{ID: 1, FileId: "file1.pdf"},
		{ID: 2, FileId: "file2.pdf"},
	}

	ctx := context.Background()

	mockRepo.EXPECT().GetFilesByEmailID(emailID, ctx).Return(expectedFiles, nil)

	files, err := useCase.GetFilesByEmailID(emailID, ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedFiles, files)
}

func TestGetFilesByEmailID_InvalidEmailID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	emailID := uint64(0)
	ctx := context.Background()

	files, err := useCase.GetFilesByEmailID(emailID, ctx)

	assert.Error(t, err)
	assert.Nil(t, files)
}

func TestGetFilesByEmailID_ErrorGettingFiles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	emailID := uint64(123)
	ctx := context.Background()

	mockRepo.EXPECT().GetFilesByEmailID(emailID, ctx).Return(nil, errors.New("file retrieval error"))

	files, err := useCase.GetFilesByEmailID(emailID, ctx)

	assert.Error(t, err)
	assert.Nil(t, files)
}

func TestDeleteFileByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	fileID := uint64(123)
	ctx := context.Background()

	mockRepo.EXPECT().DeleteFileByID(fileID, ctx).Return(nil)

	deleted, err := useCase.DeleteFileByID(fileID, ctx)

	assert.NoError(t, err)
	assert.True(t, deleted)
}

func TestDeleteFileByID_InvalidFileID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	fileID := uint64(0)
	ctx := context.Background()

	deleted, err := useCase.DeleteFileByID(fileID, ctx)

	assert.Error(t, err)
	assert.False(t, deleted)
}

func TestDeleteFileByID_ErrorDeletingFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	fileID := uint64(123)
	ctx := context.Background()

	mockRepo.EXPECT().DeleteFileByID(fileID, ctx).Return(errors.New("file deletion error"))

	deleted, err := useCase.DeleteFileByID(fileID, ctx)

	assert.Error(t, err)
	assert.False(t, deleted)
}

func TestUpdateFileByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	fileID := uint64(123)
	newFileID := "new_file_id"
	newFileType := "pdf"
	newFileName := "PDF"
	newFileSize := "10101010"
	ctx := context.Background()

	oldFile := &domain.File{
		ID:       fileID,
		FileId:   "old_file_id",
		FileType: "old_type",
		FileName: "PDF",
		FileSize: "10101010",
	}

	mockRepo.EXPECT().GetFileByID(fileID, ctx).Return(oldFile, nil)
	mockRepo.EXPECT().UpdateFileByID(fileID, newFileID, newFileType, newFileName, newFileSize, ctx).Return(nil)

	updated, err := useCase.UpdateFileByID(fileID, newFileID, newFileType, newFileName, newFileSize, ctx)

	assert.NoError(t, err)
	assert.True(t, updated)
	assert.Equal(t, newFileID, oldFile.FileId)
	assert.Equal(t, newFileType, oldFile.FileType)
}

func TestUpdateFileByID_InvalidFileID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	fileID := uint64(0)
	newFileID := "new_file_id"
	newFileType := "pdf"
	newFileName := "PDF"
	newFileSize := "10101010"
	ctx := context.Background()

	updated, err := useCase.UpdateFileByID(fileID, newFileID, newFileType, newFileName, newFileSize, ctx)

	assert.Error(t, err)
	assert.False(t, updated)
}

func TestUpdateFileByID_EmptyNewFileID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	fileID := uint64(123)
	newFileID := ""
	newFileType := "pdf"
	newFileName := "PDF"
	newFileSize := "10101010"
	ctx := context.Background()

	updated, err := useCase.UpdateFileByID(fileID, newFileID, newFileType, newFileName, newFileSize, ctx)

	assert.Error(t, err)
	assert.False(t, updated)
}

func TestUpdateFileByID_EmptyNewFileType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	fileID := uint64(123)
	newFileID := "new_file_id"
	newFileType := ""
	newFileName := "PDF"
	newFileSize := "10101010"
	ctx := context.Background()

	updated, err := useCase.UpdateFileByID(fileID, newFileID, newFileType, newFileName, newFileSize, ctx)

	assert.Error(t, err)
	assert.False(t, updated)
}

func TestUpdateFileByID_ErrorGettingFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	fileID := uint64(123)
	newFileID := "new_file_id"
	newFileType := "pdf"
	newFileName := "PDF"
	newFileSize := "10101010"
	ctx := context.Background()

	mockRepo.EXPECT().GetFileByID(fileID, ctx).Return(nil, errors.New("file retrieval error"))

	updated, err := useCase.UpdateFileByID(fileID, newFileID, newFileType, newFileName, newFileSize, ctx)

	assert.Error(t, err)
	assert.False(t, updated)
}

func TestUpdateFileByID_ErrorUpdatingFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	fileID := uint64(123)
	newFileID := "new_file_id"
	newFileType := "pdf"
	newFileName := "PDF"
	newFileSize := "10101010"
	ctx := context.Background()

	oldFile := &domain.File{
		ID:       fileID,
		FileId:   "old_file_id",
		FileType: "old_type",
		FileName: "PDF",
		FileSize: "10101010",
	}

	mockRepo.EXPECT().GetFileByID(fileID, ctx).Return(oldFile, nil)
	mockRepo.EXPECT().UpdateFileByID(fileID, newFileID, newFileType, newFileName, newFileSize, ctx).Return(errors.New("file update error"))

	updated, err := useCase.UpdateFileByID(fileID, newFileID, newFileType, newFileName, newFileSize, ctx)

	assert.Error(t, err)
	assert.False(t, updated)
}

func TestAddFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	fileID := "test_file_id"
	fileType := "test_file_type"
	fileName := "test_file_name"
	fileSize := "test_file_size"
	ctx := context.Background()

	t.Run("AddFile_Success", func(t *testing.T) {
		expectedID := uint64(123)

		mockRepo.EXPECT().AddFile(fileID, fileType, fileName, fileSize, ctx).Return(expectedID, nil)

		returnedID, err := useCase.AddFile(fileID, fileType, fileName, fileSize, ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedID, returnedID)
	})

	t.Run("AddFile_EmptyFileID", func(t *testing.T) {
		_, err := useCase.AddFile("", fileType, fileName, fileSize, ctx)
		assert.Error(t, err)
	})

	t.Run("AddFile_EmptyFileType", func(t *testing.T) {
		_, err := useCase.AddFile(fileID, "", fileName, fileSize, ctx)
		assert.Error(t, err)
	})

	t.Run("AddFile_EmptyFileName", func(t *testing.T) {
		_, err := useCase.AddFile(fileID, fileType, "", fileSize, ctx)
		assert.Error(t, err)
	})

	t.Run("AddFile_EmptyFileSize", func(t *testing.T) {
		_, err := useCase.AddFile(fileID, fileType, fileName, "", ctx)
		assert.Error(t, err)
	})

	t.Run("AddFile_DBError", func(t *testing.T) {
		mockRepo.EXPECT().AddFile(fileID, fileType, fileName, fileSize, ctx).Return(uint64(0), errors.New("DB error"))

		_, err := useCase.AddFile(fileID, fileType, fileName, fileSize, ctx)
		assert.Error(t, err)
	})
}

func TestAddFileToEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	emailID := uint64(123)
	fileID := uint64(456)
	ctx := context.Background()

	t.Run("AddFileToEmail_Success", func(t *testing.T) {
		mockRepo.EXPECT().AddAttachment(emailID, fileID, ctx).Return(nil)

		err := useCase.AddFileToEmail(emailID, fileID, ctx)

		assert.NoError(t, err)
	})

	t.Run("AddFileToEmail_InvalidEmailID", func(t *testing.T) {
		err := useCase.AddFileToEmail(0, fileID, ctx)
		assert.Error(t, err)
		assert.Equal(t, "invalid file id", err.Error())
	})

	t.Run("AddFileToEmail_InvalidFileID", func(t *testing.T) {
		err := useCase.AddFileToEmail(emailID, 0, ctx)
		assert.Error(t, err)
		assert.Equal(t, "invalid file id", err.Error())
	})

	t.Run("AddFileToEmail_DBError", func(t *testing.T) {
		mockRepo.EXPECT().AddAttachment(emailID, fileID, ctx).Return(errors.New("DB error"))

		err := useCase.AddFileToEmail(emailID, fileID, ctx)

		assert.Error(t, err)
		assert.Equal(t, "failed to add attachment", err.Error())
	})
}
