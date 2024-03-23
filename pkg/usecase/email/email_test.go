package email

import (
	"github.com/stretchr/testify/assert"
	emailCore "mail/pkg/domain/models"
	"mail/pkg/repository/converters"
	mockRepo "mail/pkg/repository/email"
	"sort"
	"testing"
)

func TestEmailUseCase_GetAllEmails(t *testing.T) {
	emails := make([]*emailCore.Email, 0, len(mockRepo.FakeEmails))
	for _, email := range mockRepo.FakeEmails {
		emails = append(emails, converters.EmailConvertDbInCore(*email))
	}
	emailRepo := mockRepo.NewEmailMemoryRepository()
	useCase := NewEmailUseCase(emailRepo)
	result, err := useCase.GetAllEmails()
	SortEmailsByID(emails)
	SortEmailsByID(result)

	assert.NoError(t, err)
	assert.Equal(t, emails, result)
}

func TestEmailUseCase_GetEmailByID(t *testing.T) {
	emails := make([]*emailCore.Email, 0, len(mockRepo.FakeEmails))
	for _, email := range mockRepo.FakeEmails {
		emails = append(emails, converters.EmailConvertDbInCore(*email))
	}
	SortEmailsByID(emails)
	emailRepo := mockRepo.NewEmailMemoryRepository()
	useCase := NewEmailUseCase(emailRepo)
	result, err := useCase.GetEmailByID(1)

	assert.NoError(t, err)
	assert.Equal(t, emails[0], result)
}

func TestEmailUseCase_CreateEmail(t *testing.T) {
	newEmail := converters.EmailConvertDbInCore(*mockRepo.FakeEmails[1])
	emailRepo := mockRepo.NewEmptyInMemoryEmailRepository()
	useCase := NewEmailUseCase(emailRepo)
	result, err := useCase.CreateEmail(newEmail)

	assert.NoError(t, err)
	assert.Equal(t, newEmail, result)
}

func TestEmailUseCase_UpdateEmail(t *testing.T) {
	updatedEmail := converters.EmailConvertDbInCore(*mockRepo.FakeEmails[1])
	updatedEmail.ReadStatus = true
	emailRepo := mockRepo.NewEmailMemoryRepository()
	useCase := NewEmailUseCase(emailRepo)
	result, err := useCase.UpdateEmail(updatedEmail)

	assert.NoError(t, err)
	assert.True(t, result)
}

func TestEmailUseCase_DeleteEmail(t *testing.T) {
	emailRepo := mockRepo.NewEmailMemoryRepository()
	useCase := NewEmailUseCase(emailRepo)
	result, err := useCase.DeleteEmail(1)

	assert.NoError(t, err)
	assert.True(t, result)
}

func SortEmailsByID(emails []*emailCore.Email) {
	sort.Slice(emails, func(i, j int) bool {
		return emails[i].ID < (emails[j].ID)
	})
}
