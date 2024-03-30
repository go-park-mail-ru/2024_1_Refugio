package email

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mock_repository "mail/pkg/domain/mock"
	domain "mail/pkg/domain/models"
	"testing"
)

func TestGetAllEmails_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	expectedEmails := []*domain.Email{
		{Topic: "Topic 1", Text: "Text 1"},
		{Topic: "Topic 2", Text: "Text 2"},
	}
	mockRepo.EXPECT().GetAll(0, 0).Return(expectedEmails, nil)

	emails, err := useCase.GetAllEmails(0, 0)

	assert.NoError(t, err)
	assert.Equal(t, expectedEmails, emails)
}

func TestGetAllEmails_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	mockRepo.EXPECT().GetAll(0, 0).Return(nil, errors.New("repository error"))

	emails, err := useCase.GetAllEmails(0, 0)

	assert.Error(t, err)
	assert.Nil(t, emails)
}

func TestGetEmailByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	expectedEmail := &domain.Email{ID: 1, Topic: "Topic 1", Text: "Text 1"}
	mockRepo.EXPECT().GetByID(uint64(1)).Return(expectedEmail, nil)

	email, err := useCase.GetEmailByID(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedEmail, email)
}

func TestGetEmailByID_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	mockRepo.EXPECT().GetByID(uint64(1)).Return(nil, errors.New("repository error"))

	email, err := useCase.GetEmailByID(1)

	assert.Error(t, err)
	assert.Nil(t, email)
}

func TestCreateEmail_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)

	newEmail := &domain.Email{Topic: "Topic 1", Text: "Text 1"}
	mockRepo.EXPECT().Add(gomock.Any()).Return(newEmail, nil)

	emailRes, err := useCase.CreateEmail(newEmail)

	assert.NoError(t, err)
	assert.Equal(t, newEmail, emailRes)
}

func TestCreateEmail_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)
	newEmail := &domain.Email{Topic: "Topic 1", Text: "Text 1"}
	mockRepo.EXPECT().Add(gomock.Any()).Return(newEmail, errors.New("repository error"))

	emailRes, err := useCase.CreateEmail(newEmail)

	assert.Error(t, err)
	assert.Equal(t, newEmail, emailRes)
}

func TestUpdateEmail_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)
	newEmail := &domain.Email{ID: 1, Topic: "Topic 1", Text: "Text 1"}
	mockRepo.EXPECT().Update(gomock.Any()).Return(true, nil)

	emailRes, err := useCase.UpdateEmail(newEmail)

	assert.NoError(t, err)
	assert.Equal(t, true, emailRes)
}

func TestDeleteEmail_ErrorFromRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockEmailRepository(ctrl)
	useCase := NewEmailUseCase(mockRepo)
	mockRepo.EXPECT().Delete(gomock.Any()).Return(true, nil)

	emailRes, err := useCase.DeleteEmail(uint64(1))

	assert.NoError(t, err)
	assert.Equal(t, true, emailRes)
}
