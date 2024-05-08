package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mock_repository "mail/internal/microservice/email/mock"
	domain "mail/internal/microservice/models/domain_models"
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

func TestNewEmailUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockEmailRepository(ctrl)

	ExpextedEmailUseCase := EmailUseCase{
		repo: mockRepo,
	}

	EmailUseCase := NewEmailUseCase(mockRepo)

	assert.Equal(t, ExpextedEmailUseCase, *EmailUseCase)
}

func TestGetAllEmailsIncoming_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	login := "test@mailhub.su"
	expectedEmails := []*domain.Email{
		{Topic: "Topic 1", Text: "Text 1"},
		{Topic: "Topic 2", Text: "Text 2"},
	}
	ctx := GetCTX()

	sero := int64(0)

	mockRepo.EXPECT().GetAllIncoming(login, sero, sero, ctx).Return(expectedEmails, nil)

	emails, err := useCase.GetAllEmailsIncoming(login, sero, sero, ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedEmails, emails)
}

func TestAllEmailsIncoming_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	login := "test@mailhub.su"
	ctx := GetCTX()
	sero := int64(0)
	//mockRepo.EXPECT().GetAllIncoming(0, 0).Return(nil, errors.New("repository error"))
	mockRepo.EXPECT().GetAllIncoming(login, sero, sero, ctx).Return(nil, errors.New("repository error"))

	emails, err := useCase.GetAllEmailsIncoming(login, sero, sero, ctx)

	assert.Error(t, err)
	assert.Nil(t, emails)
}

func TestGetAllEmailsSent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	login := "test@mailhub.su"
	expectedEmails := []*domain.Email{
		{Topic: "Topic 1", Text: "Text 1"},
		{Topic: "Topic 2", Text: "Text 2"},
	}
	ctx := GetCTX()
	sero := int64(0)
	mockRepo.EXPECT().GetAllSent(login, sero, sero, ctx).Return(expectedEmails, nil)

	emails, err := useCase.GetAllEmailsSent(login, sero, sero, ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedEmails, emails)
}

func TestGetAllEmailsSent_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	login := "test@mailhub.su"
	ctx := GetCTX()
	sero := int64(0)

	mockRepo.EXPECT().GetAllSent(login, sero, sero, ctx).Return(nil, errors.New("repository error"))

	emails, err := useCase.GetAllEmailsSent(login, sero, sero, ctx)

	assert.Error(t, err)
	assert.Nil(t, emails)
}

func TestGetAllEmailsDraft(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	login := "test@mailhub.su"
	expectedEmails := []*domain.Email{
		{Topic: "Topic 1", Text: "Text 1"},
		{Topic: "Topic 2", Text: "Text 2"},
	}
	ctx := GetCTX()
	sero := int64(0)

	mockRepo.EXPECT().GetAllDraft(login, sero, sero, ctx).Return(expectedEmails, nil)

	emails, err := useCase.GetAllDraftEmails(login, sero, sero, ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedEmails, emails)
}

func TestGetAllEmailsSpam(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	login := "test@mailhub.su"
	expectedEmails := []*domain.Email{
		{Topic: "Topic 1", Text: "Text 1"},
		{Topic: "Topic 2", Text: "Text 2"},
	}
	ctx := GetCTX()
	sero := int64(0)

	mockRepo.EXPECT().GetAllSpam(login, sero, sero, ctx).Return(expectedEmails, nil)

	emails, err := useCase.GetAllSpamEmails(login, sero, sero, ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedEmails, emails)
}

func TestGetEmailByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockEmailRepository(ctrl)
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

	mockRepo := mock_repository.NewMockEmailRepository(ctrl)
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

	mockRepo := mock_repository.NewMockEmailRepository(ctrl)
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

	mockRepo := mock_repository.NewMockEmailRepository(ctrl)
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

	mockRepo := mock_repository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	email_id := uint64(1)
	sender := "test_sender@mailhub.su"
	recipient := "test_recipient@mailhub.su"
	ctx := GetCTX()

	mockRepo.EXPECT().AddProfileEmail(email_id, sender, recipient, ctx).Return(nil)

	err := useCase.CreateProfileEmail(email_id, sender, recipient, ctx)

	assert.NoError(t, err)
}

func TestCreateProfileEmail_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	email_id := uint64(1)
	sender := "test_sender@mailhub.su"
	recipient := "test_recipient@mailhub.su"
	ctx := GetCTX()

	mockRepo.EXPECT().AddProfileEmail(email_id, sender, recipient, ctx).Return(errors.New("repository error"))

	err := useCase.CreateProfileEmail(email_id, sender, recipient, ctx)

	assert.Error(t, err)
}

func TestCheckRecipientEmail_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockEmailRepository(ctrl)
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

	mockRepo := mock_repository.NewMockEmailRepository(ctrl)
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

	mockRepo := mock_repository.NewMockEmailRepository(ctrl)
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

	mockRepo := mock_repository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	login := "test@mailhub.su"
	ctx := GetCTX()

	mockRepo.EXPECT().Delete(gomock.Any(), login, ctx).Return(true, nil)

	emailRes, err := useCase.DeleteEmail(uint64(1), login, ctx)

	assert.NoError(t, err)
	assert.Equal(t, true, emailRes)
}
