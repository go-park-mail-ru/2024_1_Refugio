package server

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"mail/internal/microservice/email/mock"
	"mail/internal/microservice/email/proto"
	"mail/internal/microservice/models/domain_models"
	converters "mail/internal/microservice/models/proto_converters"
	"mail/internal/pkg/logger"
	"os"
	"testing"
)

func GetCTX() context.Context {
	f, err := os.OpenFile("log_test.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create logfile" + "log.txt")
	}
	defer f.Close()

	ctx := context.WithValue(context.Background(), "logger", logger.InitializationBdLog(f))
	ctx2 := context.WithValue(ctx, "requestID", []string{"testID"})

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
	emailProto := converters.EmailConvertCoreInProto(*domainEmails)

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
		emailsProto[i] = converters.EmailConvertCoreInProto(*e)
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
		emailsProto[i] = converters.EmailConvertCoreInProto(*e)
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
		emailsProto[i] = converters.EmailConvertCoreInProto(*e)
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
		emailsProto[i] = converters.EmailConvertCoreInProto(*e)
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
	emailProto := converters.EmailConvertCoreInProto(*domainEmail)

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
	emailProto := converters.EmailConvertCoreInProto(*domainEmail)

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
		mockEmailUseCase.EXPECT().CheckRecipientEmail(recipient.Recipient, ctx).Return(fmt.Errorf("Recipient login not found"))

		_, err := server.CheckRecipientEmail(ctx, recipient)

		assert.Error(t, err)
		assert.Equal(t, fmt.Errorf("Recipient login not found"), err)
	})
}
