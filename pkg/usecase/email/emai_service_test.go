package email

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mock_repository "mail/pkg/domain/mock"
	domain "mail/pkg/domain/models"
	"testing"
)

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

	requestID := "test_request"
	login := "test@mailhub.su"
	expectedEmails := []*domain.Email{
		{Topic: "Topic 1", Text: "Text 1"},
		{Topic: "Topic 2", Text: "Text 2"},
	}
	mockRepo.EXPECT().GetAllIncoming(login, requestID, 0, 0).Return(expectedEmails, nil)

	emails, err := useCase.GetAllEmailsIncoming(login, requestID, 0, 0)

	assert.NoError(t, err)
	assert.Equal(t, expectedEmails, emails)
}

func TestAllEmailsIncoming_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	requestID := "test_request"
	login := "test@mailhub.su"

	//mockRepo.EXPECT().GetAllIncoming(0, 0).Return(nil, errors.New("repository error"))
	mockRepo.EXPECT().GetAllIncoming(login, requestID, 0, 0).Return(nil, errors.New("repository error"))

	emails, err := useCase.GetAllEmailsIncoming(login, requestID, 0, 0)

	assert.Error(t, err)
	assert.Nil(t, emails)
}

func TestGetAllEmailsSent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	requestID := "test_request"
	login := "test@mailhub.su"
	expectedEmails := []*domain.Email{
		{Topic: "Topic 1", Text: "Text 1"},
		{Topic: "Topic 2", Text: "Text 2"},
	}
	mockRepo.EXPECT().GetAllSent(login, requestID, 0, 0).Return(expectedEmails, nil)

	emails, err := useCase.GetAllEmailsSent(login, requestID, 0, 0)

	assert.NoError(t, err)
	assert.Equal(t, expectedEmails, emails)
}

func TestGetAllEmailsSent_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	requestID := "test_request"
	login := "test@mailhub.su"

	//mockRepo.EXPECT().GetAllIncoming(0, 0).Return(nil, errors.New("repository error"))
	mockRepo.EXPECT().GetAllSent(login, requestID, 0, 0).Return(nil, errors.New("repository error"))

	emails, err := useCase.GetAllEmailsSent(login, requestID, 0, 0)

	assert.Error(t, err)
	assert.Nil(t, emails)
}

func TestGetEmailByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	requestID := "test_request"
	login := "test@mailhub.su"

	expectedEmail := &domain.Email{ID: 1, Topic: "Topic 1", Text: "Text 1"}
	mockRepo.EXPECT().GetByID(uint64(1), login, requestID).Return(expectedEmail, nil)

	email, err := useCase.GetEmailByID(1, login, requestID)

	assert.NoError(t, err)
	assert.Equal(t, expectedEmail, email)
}

func TestGetEmailByID_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	requestID := "test_request"
	login := "test@mailhub.su"

	mockRepo.EXPECT().GetByID(uint64(1), login, requestID).Return(nil, errors.New("repository error"))

	email, err := useCase.GetEmailByID(1, login, requestID)

	assert.Error(t, err)
	assert.Nil(t, email)
}

func TestCreateEmail_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	requestID := "test_request"

	newEmail := &domain.Email{Topic: "Topic 1", Text: "Text 1"}
	mockRepo.EXPECT().Add(gomock.Any(), requestID).Return(int64(1), newEmail, nil)

	id, emailRes, err := useCase.CreateEmail(newEmail, requestID)

	assert.Equal(t, int64(1), id)
	assert.NoError(t, err)
	assert.Equal(t, newEmail, emailRes)
}

func TestCreateEmail_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)
	newEmail := &domain.Email{Topic: "Topic 1", Text: "Text 1"}

	requestID := "test_request"

	mockRepo.EXPECT().Add(gomock.Any(), requestID).Return(int64(1), newEmail, errors.New("repository error"))

	id, emailRes, err := useCase.CreateEmail(newEmail, requestID)

	assert.Equal(t, int64(1), id)
	assert.Error(t, err)
	assert.Equal(t, newEmail, emailRes)
}

func TestCreateProfileEmail_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	email_id := int64(1)
	sender := "test_sender@mailhub.su"
	recipient := "test_recipient@mailhub.su"
	requestID := "test_request"

	mockRepo.EXPECT().AddProfileEmail(email_id, sender, recipient, requestID).Return(nil)

	err := useCase.CreateProfileEmail(email_id, sender, recipient, requestID)

	assert.NoError(t, err)
}

func TestCreateProfileEmail_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	email_id := int64(1)
	sender := "test_sender@mailhub.su"
	recipient := "test_recipient@mailhub.su"
	requestID := "test_request"

	mockRepo.EXPECT().AddProfileEmail(email_id, sender, recipient, requestID).Return(errors.New("repository error"))

	err := useCase.CreateProfileEmail(email_id, sender, recipient, requestID)

	assert.Error(t, err)
}

func TestCheckRecipientEmail_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	recipient := "test_recipient@mailhub.su"
	requestID := "test_request"

	mockRepo.EXPECT().FindEmail(recipient, requestID).Return(nil)

	err := useCase.CheckRecipientEmail(recipient, requestID)

	assert.NoError(t, err)
}

func TestCheckRecipientEmail_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	recipient := "test_recipient@mailhub.su"
	requestID := "test_request"

	mockRepo.EXPECT().FindEmail(recipient, requestID).Return(errors.New("repository error"))

	err := useCase.CheckRecipientEmail(recipient, requestID)

	assert.Error(t, err)
}

func TestUpdateEmail_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	requestID := "test_request"
	newEmail := &domain.Email{ID: 1, Topic: "Topic 1", Text: "Text 1"}

	mockRepo.EXPECT().Update(gomock.Any(), requestID).Return(true, nil)

	emailRes, err := useCase.UpdateEmail(newEmail, requestID)

	assert.NoError(t, err)
	assert.Equal(t, true, emailRes)
}

func TestDeleteEmail_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	login := "test@mailhub.su"
	requestID := "test_request"

	mockRepo.EXPECT().Delete(gomock.Any(), login, requestID).Return(true, nil)

	emailRes, err := useCase.DeleteEmail(uint64(1), login, requestID)

	assert.NoError(t, err)
	assert.Equal(t, true, emailRes)
}
