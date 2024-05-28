package server

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"mail/internal/microservice/email/mock"
	"mail/internal/microservice/email/proto"
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

func TestNewEmailServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mock.NewMockEmailUseCase(ctrl)

	server := NewEmailServer(mockEmailUseCase)

	expectedServer := EmailServer{
		EmailUseCase: mockEmailUseCase,
	}

	assert.Equal(t, expectedServer, *server)
}

func TestGetEmailByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mock.NewMockEmailUseCase(ctrl)

	server := NewEmailServer(mockEmailUseCase)

	ctx := GetCTX()
	login := "test@mailhub.su"
	id := uint64(1)

	domainEmails := &domain_models.Email{ID: 1, Topic: "Topic 1", Text: "Text 1"}
	emailProto := converters.EmailConvertCoreInProto(domainEmails)

	t.Run("GetEmailByIDSuccessfully", func(t *testing.T) {
		mockEmailUseCase.EXPECT().GetEmailByID(id, login, ctx).Return(domainEmails, nil)

		email, err := server.GetEmailByID(ctx, &proto.EmailIdAndLogin{Id: id, Login: login})

		assert.NoError(t, err)
		assert.Equal(t, emailProto, email)
	})

	t.Run("GetEmailByIDFail invalid invalid email id", func(t *testing.T) {
		_, err := server.GetEmailByID(ctx, &proto.EmailIdAndLogin{Id: uint64(0), Login: login})
		assert.Error(t, err)
	})

	t.Run("GetEmailByIDFail email not found", func(t *testing.T) {
		mockEmailUseCase.EXPECT().GetEmailByID(id, login, ctx).Return(nil, fmt.Errorf("email not found"))

		_, err := server.GetEmailByID(ctx, &proto.EmailIdAndLogin{Id: id, Login: login})

		assert.Error(t, err)
		assert.Equal(t, fmt.Errorf("email not found"), err)
	})
}

func TestGetAllIncoming(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mock.NewMockEmailUseCase(ctrl)

	server := NewEmailServer(mockEmailUseCase)

	ctx := GetCTX()
	login := "test@mailhub.su"
	offset := int64(0)
	limit := int64(0)

	domainEmails := []*domain_models.Email{
		{ID: 1, Topic: "Topic 1", Text: "Text 1"},
		{ID: 2, Topic: "Topic 2", Text: "Text 2"},
	}

	emailsProto := make([]*proto.Email, len(domainEmails))
	for i, e := range domainEmails {
		emailsProto[i] = converters.EmailConvertCoreInProto(e)
	}
	emailProto := new(proto.Emails)
	emailProto.Emails = emailsProto

	t.Run("EmailGetAllIncomingSuccessfully", func(t *testing.T) {
		mockEmailUseCase.EXPECT().GetAllEmailsIncoming(login, offset, limit, ctx).Return(domainEmails, nil)

		emails, err := server.GetAllIncoming(ctx, &proto.LoginOffsetLimit{Login: login, Offset: offset, Limit: limit})

		assert.NoError(t, err)
		assert.Equal(t, emailProto, emails)
	})

	t.Run("EmailGetAllIncomingFail invalid email login", func(t *testing.T) {
		_, err := server.GetAllIncoming(ctx, &proto.LoginOffsetLimit{Login: "", Offset: offset, Limit: limit})
		assert.Error(t, err)
	})

	t.Run("EmailGetAllIncomingFail email not found", func(t *testing.T) {
		mockEmailUseCase.EXPECT().GetAllEmailsIncoming(login, offset, limit, ctx).Return(nil, fmt.Errorf("email not found"))

		_, err := server.GetAllIncoming(ctx, &proto.LoginOffsetLimit{Login: login, Offset: offset, Limit: limit})

		assert.Error(t, err)
		assert.Equal(t, fmt.Errorf("email not found"), err)
	})
}

func TestGetAllSent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mock.NewMockEmailUseCase(ctrl)

	server := NewEmailServer(mockEmailUseCase)

	ctx := GetCTX()
	login := "test@mailhub.su"
	offset := int64(0)
	limit := int64(0)

	domainEmails := []*domain_models.Email{
		{ID: 1, Topic: "Topic 1", Text: "Text 1"},
		{ID: 2, Topic: "Topic 2", Text: "Text 2"},
	}

	emailsProto := make([]*proto.Email, len(domainEmails))
	for i, e := range domainEmails {
		emailsProto[i] = converters.EmailConvertCoreInProto(e)
	}
	emailProto := new(proto.Emails)
	emailProto.Emails = emailsProto

	t.Run("EmailGetAllSentSuccessfully", func(t *testing.T) {
		mockEmailUseCase.EXPECT().GetAllEmailsSent(login, offset, limit, ctx).Return(domainEmails, nil)

		emails, err := server.GetAllSent(ctx, &proto.LoginOffsetLimit{Login: login, Offset: offset, Limit: limit})

		assert.NoError(t, err)
		assert.Equal(t, emailProto, emails)
	})

	t.Run("EmailGetAllSentFail invalid email login", func(t *testing.T) {
		_, err := server.GetAllSent(ctx, &proto.LoginOffsetLimit{Login: "", Offset: offset, Limit: limit})
		assert.Error(t, err)
	})

	t.Run("EmailGetAllSentFail email not found", func(t *testing.T) {
		mockEmailUseCase.EXPECT().GetAllEmailsSent(login, offset, limit, ctx).Return(nil, fmt.Errorf("email not found"))

		_, err := server.GetAllSent(ctx, &proto.LoginOffsetLimit{Login: login, Offset: offset, Limit: limit})

		assert.Error(t, err)
		assert.Equal(t, fmt.Errorf("email not found"), err)
	})
}

func TestGetAllDraft(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mock.NewMockEmailUseCase(ctrl)

	server := NewEmailServer(mockEmailUseCase)

	ctx := GetCTX()
	login := "test@mailhub.su"
	offset := int64(0)
	limit := int64(0)

	domainEmails := []*domain_models.Email{
		{ID: 1, Topic: "Topic 1", Text: "Text 1"},
		{ID: 2, Topic: "Topic 2", Text: "Text 2"},
	}

	emailsProto := make([]*proto.Email, len(domainEmails))
	for i, e := range domainEmails {
		emailsProto[i] = converters.EmailConvertCoreInProto(e)
	}
	emailProto := new(proto.Emails)
	emailProto.Emails = emailsProto

	t.Run("EmailGetAllDraftSuccessfully", func(t *testing.T) {
		mockEmailUseCase.EXPECT().GetAllDraftEmails(login, offset, limit, ctx).Return(domainEmails, nil)

		emails, err := server.GetDraftEmails(ctx, &proto.LoginOffsetLimit{Login: login, Offset: offset, Limit: limit})

		assert.NoError(t, err)
		assert.Equal(t, emailProto, emails)
	})

	t.Run("EmailGetAllDraftFail invalid email login", func(t *testing.T) {
		_, err := server.GetDraftEmails(ctx, &proto.LoginOffsetLimit{Login: "", Offset: offset, Limit: limit})
		assert.Error(t, err)
	})

	t.Run("EmailGetAllDraftFail email not found", func(t *testing.T) {
		mockEmailUseCase.EXPECT().GetAllDraftEmails(login, offset, limit, ctx).Return(nil, fmt.Errorf("email not found"))

		_, err := server.GetDraftEmails(ctx, &proto.LoginOffsetLimit{Login: login, Offset: offset, Limit: limit})

		assert.Error(t, err)
		assert.Equal(t, fmt.Errorf("email not found"), err)
	})
}

func TestGetAllSpam(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mock.NewMockEmailUseCase(ctrl)

	server := NewEmailServer(mockEmailUseCase)

	ctx := GetCTX()
	login := "test@mailhub.su"
	offset := int64(0)
	limit := int64(0)

	domainEmails := []*domain_models.Email{
		{ID: 1, Topic: "Topic 1", Text: "Text 1"},
		{ID: 2, Topic: "Topic 2", Text: "Text 2"},
	}

	emailsProto := make([]*proto.Email, len(domainEmails))
	for i, e := range domainEmails {
		emailsProto[i] = converters.EmailConvertCoreInProto(e)
	}
	emailProto := new(proto.Emails)
	emailProto.Emails = emailsProto

	t.Run("EmailGetAllSpamSuccessfully", func(t *testing.T) {
		mockEmailUseCase.EXPECT().GetAllSpamEmails(login, offset, limit, ctx).Return(domainEmails, nil)

		emails, err := server.GetSpamEmails(ctx, &proto.LoginOffsetLimit{Login: login, Offset: offset, Limit: limit})

		assert.NoError(t, err)
		assert.Equal(t, emailProto, emails)
	})

	t.Run("EmailGetAllSpamFail invalid email login", func(t *testing.T) {
		_, err := server.GetSpamEmails(ctx, &proto.LoginOffsetLimit{Login: "", Offset: offset, Limit: limit})
		assert.Error(t, err)
	})

	t.Run("EmailGetAllSpamFail email not found", func(t *testing.T) {
		mockEmailUseCase.EXPECT().GetAllSpamEmails(login, offset, limit, ctx).Return(nil, fmt.Errorf("email not found"))

		_, err := server.GetSpamEmails(ctx, &proto.LoginOffsetLimit{Login: login, Offset: offset, Limit: limit})

		assert.Error(t, err)
		assert.Equal(t, fmt.Errorf("email not found"), err)
	})
}

func TestCreateEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mock.NewMockEmailUseCase(ctrl)

	server := NewEmailServer(mockEmailUseCase)

	ctx := GetCTX()
	id := uint64(1)

	domainEmail := &domain_models.Email{ID: 1, Topic: "Topic 1", Text: "Text 1"}
	emailProto := converters.EmailConvertCoreInProto(domainEmail)

	emailWithID := &proto.EmailWithID{Email: emailProto, Id: id}

	t.Run("CreateEmailSuccessfully", func(t *testing.T) {
		mockEmailUseCase.EXPECT().CreateEmail(domainEmail, ctx).Return(id, domainEmail, nil)

		email, err := server.CreateEmail(ctx, emailProto)

		assert.NoError(t, err)
		assert.Equal(t, emailWithID, email)
	})

	t.Run("CreateEmailFail invalid email format", func(t *testing.T) {
		_, err := server.CreateEmail(ctx, nil)

		assert.Error(t, err)
	})

	t.Run("CreateEmailFail failed create email", func(t *testing.T) {
		mockEmailUseCase.EXPECT().CreateEmail(domainEmail, ctx).Return(id, domainEmail, fmt.Errorf("failed create email"))

		_, err := server.CreateEmail(ctx, emailProto)

		assert.Error(t, err)
		assert.Equal(t, fmt.Errorf("failed create email"), err)
	})
}

func TestUpdateEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mock.NewMockEmailUseCase(ctrl)

	server := NewEmailServer(mockEmailUseCase)

	ctx := GetCTX()

	domainEmail := &domain_models.Email{ID: 1, Topic: "Topic 1", Text: "Text 1"}
	emailProto := converters.EmailConvertCoreInProto(domainEmail)

	t.Run("UpdateEmailSuccessfully", func(t *testing.T) {
		mockEmailUseCase.EXPECT().UpdateEmail(domainEmail, ctx).Return(true, nil)

		emailStatus, err := server.UpdateEmail(ctx, emailProto)

		assert.NoError(t, err)
		assert.True(t, emailStatus.Status)
	})

	t.Run("UpdateEmailFail invalid email format", func(t *testing.T) {
		_, err := server.UpdateEmail(ctx, nil)

		assert.Error(t, err)
	})

	t.Run("UpdateEmailFail email not found", func(t *testing.T) {
		mockEmailUseCase.EXPECT().UpdateEmail(domainEmail, ctx).Return(false, fmt.Errorf("email not found"))

		_, err := server.UpdateEmail(ctx, emailProto)

		assert.Error(t, err)
		assert.Equal(t, fmt.Errorf("email not found"), err)
	})

}

func TestDeleteEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mock.NewMockEmailUseCase(ctrl)

	server := NewEmailServer(mockEmailUseCase)

	ctx := GetCTX()

	loginWithId := &proto.LoginWithID{Id: uint64(1), Login: "test@mailhub.su"}

	t.Run("DeleteEmailSuccessfully", func(t *testing.T) {
		mockEmailUseCase.EXPECT().DeleteEmail(loginWithId.Id, loginWithId.Login, ctx).Return(true, nil)

		emailStatus, err := server.DeleteEmail(ctx, loginWithId)

		assert.NoError(t, err)
		assert.True(t, emailStatus.Status)
	})

	t.Run("DeleteEmailFail invalid input data", func(t *testing.T) {
		_, err := server.DeleteEmail(ctx, &proto.LoginWithID{Id: uint64(0), Login: "test@mailhub.su"})

		assert.Error(t, err)
		assert.Equal(t, fmt.Errorf("invalid input data"), err)
	})

	t.Run("DeleteEmailFail email not found", func(t *testing.T) {
		mockEmailUseCase.EXPECT().DeleteEmail(loginWithId.Id, loginWithId.Login, ctx).Return(false, fmt.Errorf("email not found"))

		_, err := server.DeleteEmail(ctx, loginWithId)

		assert.Error(t, err)
		assert.Equal(t, fmt.Errorf("email not found"), err)
	})

}

func TestCreateProfileEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mock.NewMockEmailUseCase(ctrl)

	server := NewEmailServer(mockEmailUseCase)

	ctx := GetCTX()

	idSenderRecipient := &proto.IdSenderRecipient{Id: uint64(1), Sender: "test@mailhub.su", Recipient: "test@mailhub.su"}

	t.Run("CreateProfileEmailSuccessfully", func(t *testing.T) {
		mockEmailUseCase.EXPECT().CreateProfileEmail(idSenderRecipient.Id, idSenderRecipient.Sender, idSenderRecipient.Recipient, ctx).Return(nil)

		_, err := server.CreateProfileEmail(ctx, idSenderRecipient)

		assert.NoError(t, err)
	})

	t.Run("CreateProfileEmailFail invalid email id", func(t *testing.T) {
		_, err := server.CreateProfileEmail(ctx, &proto.IdSenderRecipient{Id: uint64(0), Sender: "test@mailhub.su", Recipient: "test@mailhub.su"})

		assert.Error(t, err)
	})

	t.Run("CreateProfileEmailFail invalid email sender or recipient login", func(t *testing.T) {
		_, err := server.CreateProfileEmail(ctx, &proto.IdSenderRecipient{Id: uint64(1), Sender: "", Recipient: "test@mailhub.su"})

		assert.Error(t, err)
	})

	t.Run("CreateProfileEmailFail", func(t *testing.T) {
		mockEmailUseCase.EXPECT().CreateProfileEmail(idSenderRecipient.Id, idSenderRecipient.Sender, idSenderRecipient.Recipient, ctx).Return(fmt.Errorf("sender, recipient or id not found"))

		_, err := server.CreateProfileEmail(ctx, idSenderRecipient)

		assert.Error(t, err)
		assert.Equal(t, fmt.Errorf("sender, recipient or id not found"), err)
	})

}

func TestCheckRecipientEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mock.NewMockEmailUseCase(ctrl)

	server := NewEmailServer(mockEmailUseCase)

	ctx := GetCTX()

	recipient := &proto.Recipient{Recipient: "test@mailhub.su"}

	t.Run("CheckRecipientEmailSuccessfully", func(t *testing.T) {
		mockEmailUseCase.EXPECT().CheckRecipientEmail(recipient.Recipient, ctx).Return(nil)

		_, err := server.CheckRecipientEmail(ctx, recipient)

		assert.NoError(t, err)
	})

	t.Run("CheckRecipientEmailFail invalid recipient login", func(t *testing.T) {
		_, err := server.CheckRecipientEmail(ctx, &proto.Recipient{Recipient: ""})

		assert.Error(t, err)
	})

	t.Run("CheckRecipientEmailFail Recipient login not found", func(t *testing.T) {
		mockEmailUseCase.EXPECT().CheckRecipientEmail(recipient.Recipient, ctx).Return(fmt.Errorf("recipient login not found"))

		_, err := server.CheckRecipientEmail(ctx, recipient)

		assert.Error(t, err)
		assert.Equal(t, fmt.Errorf("recipient login not found"), err)
	})
}

func TestAddAttachment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mock.NewMockEmailUseCase(ctrl)

	server := NewEmailServer(mockEmailUseCase)

	ctx := GetCTX()

	fileID := "test_file_id"
	fileType := "pdf"
	fileName := "PDF"
	fileSize := "10101010"
	emailID := uint64(123)

	t.Run("AddAttachment_Success", func(t *testing.T) {
		mockEmailUseCase.EXPECT().AddAttachment(fileID, fileType, fileName, fileSize, emailID, ctx).Return(uint64(456), nil)

		request := &proto.AddAttachmentRequest{
			FileId:   fileID,
			FileType: fileType,
			FileName: fileName,
			FileSize: fileSize,
			EmailId:  emailID,
		}

		reply, err := server.AddAttachment(ctx, request)

		assert.NoError(t, err)
		assert.Equal(t, uint64(456), reply.FileId)
	})

	t.Run("AddAttachment_NilInput", func(t *testing.T) {
		request := (*proto.AddAttachmentRequest)(nil)

		reply, err := server.AddAttachment(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, reply)
	})

	t.Run("AddAttachment_EmptyFileID", func(t *testing.T) {
		request := &proto.AddAttachmentRequest{
			FileId:   "",
			FileType: fileType,
			EmailId:  emailID,
		}

		reply, err := server.AddAttachment(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, reply)
	})

	t.Run("AddAttachment_EmptyFileType", func(t *testing.T) {
		request := &proto.AddAttachmentRequest{
			FileId:   fileID,
			FileType: "",
			EmailId:  emailID,
		}

		reply, err := server.AddAttachment(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, reply)
	})

	t.Run("AddAttachment_InvalidEmailID", func(t *testing.T) {
		request := &proto.AddAttachmentRequest{
			FileId:   fileID,
			FileType: fileType,
			EmailId:  uint64(0),
		}

		reply, err := server.AddAttachment(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, reply)
	})

	t.Run("AddAttachment_FailedToAddAttachment", func(t *testing.T) {
		mockEmailUseCase.EXPECT().AddAttachment(fileID, fileType, fileName, fileSize, emailID, ctx).Return(uint64(0), fmt.Errorf("failed to add attachment"))

		request := &proto.AddAttachmentRequest{
			FileId:   fileID,
			FileType: fileType,
			FileName: fileName,
			FileSize: fileSize,
			EmailId:  emailID,
		}

		reply, err := server.AddAttachment(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, reply)
	})
}

func TestGetFileByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mock.NewMockEmailUseCase(ctrl)

	server := NewEmailServer(mockEmailUseCase)

	ctx := GetCTX()

	fileID := uint64(123)

	t.Run("GetFileByID_Success", func(t *testing.T) {
		mockEmailUseCase.EXPECT().GetFileByID(fileID, ctx).Return(&domain_models.File{
			ID:       fileID,
			FileId:   "test_file_id",
			FileType: "pdf",
		}, nil)

		request := &proto.GetFileByIDRequest{
			FileId: fileID,
		}

		reply, err := server.GetFileByID(ctx, request)

		assert.NoError(t, err)
		assert.NotNil(t, reply)
	})

	t.Run("GetFileByID_NilInput", func(t *testing.T) {
		request := (*proto.GetFileByIDRequest)(nil)

		reply, err := server.GetFileByID(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, reply)
	})

	t.Run("GetFileByID_InvalidFileID", func(t *testing.T) {
		request := &proto.GetFileByIDRequest{
			FileId: uint64(0),
		}

		reply, err := server.GetFileByID(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, reply)
	})

	t.Run("GetFileByID_FailedToGetFile", func(t *testing.T) {
		mockEmailUseCase.EXPECT().GetFileByID(fileID, ctx).Return(nil, fmt.Errorf("failed to get file"))

		request := &proto.GetFileByIDRequest{
			FileId: fileID,
		}

		reply, err := server.GetFileByID(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, reply)
	})
}

func TestGetFilesByEmailID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mock.NewMockEmailUseCase(ctrl)

	server := NewEmailServer(mockEmailUseCase)

	ctx := GetCTX()

	emailID := uint64(123)

	t.Run("GetFilesByEmailID_Success", func(t *testing.T) {
		mockEmailUseCase.EXPECT().GetFilesByEmailID(emailID, ctx).Return([]*domain_models.File{
			{
				ID:       1,
				FileId:   "file_id_1",
				FileType: "pdf",
			},
			{
				ID:       2,
				FileId:   "file_id_2",
				FileType: "docx",
			},
		}, nil)

		request := &proto.GetFilesByEmailIDRequest{
			EmailId: emailID,
		}

		reply, err := server.GetFilesByEmailID(ctx, request)

		assert.NoError(t, err)
		assert.NotNil(t, reply)
		assert.Len(t, reply.Files, 2)
	})

	t.Run("GetFilesByEmailID_NilInput", func(t *testing.T) {
		request := (*proto.GetFilesByEmailIDRequest)(nil)

		reply, err := server.GetFilesByEmailID(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, reply)
	})

	t.Run("GetFilesByEmailID_InvalidEmailID", func(t *testing.T) {
		request := &proto.GetFilesByEmailIDRequest{
			EmailId: 0,
		}

		reply, err := server.GetFilesByEmailID(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, reply)
	})

	t.Run("GetFilesByEmailID_FailedToGetFiles", func(t *testing.T) {
		mockEmailUseCase.EXPECT().GetFilesByEmailID(emailID, ctx).Return(nil, fmt.Errorf("failed to get files"))

		request := &proto.GetFilesByEmailIDRequest{
			EmailId: emailID,
		}

		reply, err := server.GetFilesByEmailID(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, reply)
	})
}

func TestDeleteFileByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mock.NewMockEmailUseCase(ctrl)

	server := NewEmailServer(mockEmailUseCase)

	ctx := GetCTX()

	fileID := uint64(123)

	t.Run("DeleteFileByID_Success", func(t *testing.T) {
		mockEmailUseCase.EXPECT().DeleteFileByID(fileID, ctx).Return(true, nil)

		request := &proto.DeleteFileByIDRequest{
			FileId: fileID,
		}

		reply, err := server.DeleteFileByID(ctx, request)

		assert.NoError(t, err)
		assert.NotNil(t, reply)
		assert.True(t, reply.Status)
	})

	t.Run("DeleteFileByID_NilInput", func(t *testing.T) {
		request := (*proto.DeleteFileByIDRequest)(nil)

		reply, err := server.DeleteFileByID(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, reply)
	})

	t.Run("DeleteFileByID_InvalidFileID", func(t *testing.T) {
		request := &proto.DeleteFileByIDRequest{
			FileId: uint64(0),
		}

		reply, err := server.DeleteFileByID(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, reply)
	})

	t.Run("DeleteFileByID_FailedToDeleteFile", func(t *testing.T) {
		mockEmailUseCase.EXPECT().DeleteFileByID(fileID, ctx).Return(false, nil)

		request := &proto.DeleteFileByIDRequest{
			FileId: fileID,
		}

		reply, err := server.DeleteFileByID(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, reply)
	})

	t.Run("DeleteFileByID_FileNotDeleted", func(t *testing.T) {
		mockEmailUseCase.EXPECT().DeleteFileByID(fileID, ctx).Return(false, nil)

		request := &proto.DeleteFileByIDRequest{
			FileId: fileID,
		}

		reply, err := server.DeleteFileByID(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, reply)
	})
}

func TestUpdateFileByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mock.NewMockEmailUseCase(ctrl)

	server := NewEmailServer(mockEmailUseCase)

	ctx := context.Background()

	fileID := uint64(123)
	newFileID := "new_file_id"
	newFileType := "pdf"
	newFileName := "PDF"
	newFileSize := "10101010"

	t.Run("UpdateFileByID_Success", func(t *testing.T) {
		mockEmailUseCase.EXPECT().UpdateFileByID(fileID, newFileID, newFileType, newFileName, newFileSize, ctx).Return(true, nil)

		request := &proto.UpdateFileByIDRequest{
			Id:          fileID,
			NewFileId:   newFileID,
			NewFileType: newFileType,
			NewFileName: newFileName,
			NewFileSize: newFileSize,
		}

		reply, err := server.UpdateFileByID(ctx, request)

		assert.NoError(t, err)
		assert.NotNil(t, reply)
		assert.True(t, reply.Status)
	})

	t.Run("UpdateFileByID_NilInput", func(t *testing.T) {
		request := (*proto.UpdateFileByIDRequest)(nil)

		reply, err := server.UpdateFileByID(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, reply)
	})

	t.Run("UpdateFileByID_InvalidFileID", func(t *testing.T) {
		request := &proto.UpdateFileByIDRequest{
			Id:          uint64(0),
			NewFileId:   newFileID,
			NewFileType: newFileType,
			NewFileSize: newFileSize,
			NewFileName: newFileName,
		}

		reply, err := server.UpdateFileByID(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, reply)
	})

	t.Run("UpdateFileByID_EmptyNewFileId", func(t *testing.T) {
		request := &proto.UpdateFileByIDRequest{
			Id:          fileID,
			NewFileId:   "",
			NewFileType: newFileType,
		}

		reply, err := server.UpdateFileByID(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, reply)
	})

	t.Run("UpdateFileByID_EmptyNewFileType", func(t *testing.T) {
		request := &proto.UpdateFileByIDRequest{
			Id:          fileID,
			NewFileId:   newFileID,
			NewFileType: "",
		}

		reply, err := server.UpdateFileByID(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, reply)
	})

	t.Run("UpdateFileByID_FailedToUpdateFile", func(t *testing.T) {
		mockEmailUseCase.EXPECT().UpdateFileByID(fileID, newFileID, newFileType, newFileName, newFileSize, ctx).Return(false, nil)

		request := &proto.UpdateFileByIDRequest{
			Id:          fileID,
			NewFileId:   newFileID,
			NewFileType: newFileType,
			NewFileName: newFileName,
			NewFileSize: newFileSize,
		}

		reply, err := server.UpdateFileByID(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, reply)
	})

	t.Run("UpdateFileByID_FileNotUpdated", func(t *testing.T) {
		mockEmailUseCase.EXPECT().UpdateFileByID(fileID, newFileID, newFileType, newFileName, newFileSize, ctx).Return(false, nil)

		request := &proto.UpdateFileByIDRequest{
			Id:          fileID,
			NewFileId:   newFileID,
			NewFileType: newFileType,
			NewFileName: newFileName,
			NewFileSize: newFileSize,
		}

		reply, err := server.UpdateFileByID(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, reply)
	})
}

func TestAddEmailDraft(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mock.NewMockEmailUseCase(ctrl)
	server := NewEmailServer(mockEmailUseCase)
	ctx := context.Background()

	emailID := uint64(123)
	createdEmail := &proto.Email{
		SenderEmail:    "sender",
		RecipientEmail: "recipient",
		Topic:          "topic",
		Text:           "text",
	}

	t.Run("AddEmailDraft_Success", func(t *testing.T) {
		mockEmailUseCase.EXPECT().CreateEmail(converters.EmailConvertProtoInCore(createdEmail), ctx).Return(emailID, converters.EmailConvertProtoInCore(createdEmail), nil)
		mockEmailUseCase.EXPECT().CreateProfileEmail(emailID, converters.EmailConvertProtoInCore(createdEmail).SenderEmail, "", ctx).Return(nil)

		request := &proto.Email{
			SenderEmail:    "sender",
			RecipientEmail: "recipient",
			Topic:          "topic",
			Text:           "text",
		}

		reply, err := server.AddEmailDraft(ctx, request)

		assert.NoError(t, err)
		assert.NotNil(t, reply)
		assert.Equal(t, emailID, reply.Id)
	})

	t.Run("AddEmailDraft_NilInput", func(t *testing.T) {
		reply, err := server.AddEmailDraft(ctx, nil)

		assert.Error(t, err)
		assert.Nil(t, reply)
	})

	t.Run("AddEmailDraft_FailedToCreateEmail", func(t *testing.T) {
		mockEmailUseCase.EXPECT().CreateEmail(converters.EmailConvertProtoInCore(createdEmail), ctx).Return(uint64(0), nil, errors.New("failed create email"))

		request := &proto.Email{
			SenderEmail:    "sender",
			RecipientEmail: "recipient",
			Topic:          "topic",
			Text:           "text",
		}

		reply, err := server.AddEmailDraft(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, reply)
	})

	t.Run("AddEmailDraft_FailedToCreateProfileEmail", func(t *testing.T) {
		mockEmailUseCase.EXPECT().CreateEmail(converters.EmailConvertProtoInCore(createdEmail), ctx).Return(emailID, converters.EmailConvertProtoInCore(createdEmail), nil)
		mockEmailUseCase.EXPECT().CreateProfileEmail(emailID, converters.EmailConvertProtoInCore(createdEmail).SenderEmail, "", ctx).Return(errors.New("failed create profile email"))

		request := &proto.Email{
			SenderEmail:    "sender",
			RecipientEmail: "recipient",
			Topic:          "topic",
			Text:           "text",
		}

		reply, err := server.AddEmailDraft(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, reply)
	})
}

func TestAddFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mock.NewMockEmailUseCase(ctrl)
	server := NewEmailServer(mockEmailUseCase)
	ctx := context.Background()

	fileId := "123"
	fileType := "pdf"
	fileName := "name"
	fileSize := "size"

	t.Run("AddFile_Success", func(t *testing.T) {
		mockEmailUseCase.EXPECT().AddFile(fileId, fileType, fileName, fileSize, ctx).Return(uint64(1), nil)

		request := &proto.AddFileRequest{
			FileId:   fileId,
			FileType: fileType,
			FileName: fileName,
			FileSize: fileSize,
		}

		reply, err := server.AddFile(ctx, request)

		assert.NoError(t, err)
		assert.NotNil(t, reply)
	})

	t.Run("AddFile_NilInput", func(t *testing.T) {
		reply, err := server.AddFile(ctx, nil)

		assert.Error(t, err)
		assert.Nil(t, reply)
	})

	t.Run("AddFile_EmptyFields", func(t *testing.T) {
		request := &proto.AddFileRequest{
			FileId:   "",
			FileType: "",
			FileName: "",
			FileSize: "",
		}

		reply, err := server.AddFile(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, reply)
	})

	t.Run("AddFile_FailedToAddFile", func(t *testing.T) {
		mockEmailUseCase.EXPECT().AddFile(fileId, fileType, fileName, fileSize, ctx).Return(uint64(1), fmt.Errorf("failed to add file"))

		request := &proto.AddFileRequest{
			FileId:   fileId,
			FileType: fileType,
			FileName: fileName,
			FileSize: fileSize,
		}

		reply, err := server.AddFile(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, reply)
	})
}

func TestAddFileToEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mock.NewMockEmailUseCase(ctrl)
	server := NewEmailServer(mockEmailUseCase)
	ctx := context.Background()

	emailId := uint64(123)
	fileId := uint64(234)

	t.Run("AddFileToEmail_Success", func(t *testing.T) {
		mockEmailUseCase.EXPECT().AddFileToEmail(emailId, fileId, ctx).Return(nil)

		request := &proto.AddFileToEmailRequest{
			EmailId: emailId,
			FileId:  fileId,
		}

		reply, err := server.AddFileToEmail(ctx, request)

		assert.NoError(t, err)
		assert.NotNil(t, reply)
		assert.Equal(t, true, reply.Status)
	})

	t.Run("AddFileToEmail_NilInput", func(t *testing.T) {
		reply, err := server.AddFileToEmail(ctx, nil)

		assert.Error(t, err)
		assert.Nil(t, reply)
	})

	t.Run("AddFileToEmail_InvalidEmailId", func(t *testing.T) {
		request := &proto.AddFileToEmailRequest{
			EmailId: 0,
			FileId:  fileId,
		}

		reply, err := server.AddFileToEmail(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, reply)
	})

	t.Run("AddFileToEmail_InvalidFileId", func(t *testing.T) {
		request := &proto.AddFileToEmailRequest{
			EmailId: emailId,
			FileId:  0,
		}

		reply, err := server.AddFileToEmail(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, reply)
	})

	t.Run("AddFileToEmail_FailedToAddFileToEmail", func(t *testing.T) {
		mockEmailUseCase.EXPECT().AddFileToEmail(emailId, fileId, ctx).Return(fmt.Errorf("failed to add file to email"))

		request := &proto.AddFileToEmailRequest{
			EmailId: emailId,
			FileId:  fileId,
		}

		reply, err := server.AddFileToEmail(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, reply)
	})
}
